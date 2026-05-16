# ToC 消费者端主题插件系统 — 实施方案

## Context

GotribeCMS 目前只有 Admin 后台（React SPA）+ 公共 API（`cmd/api`），缺少消费者-facing 网站。用户需要在后台可选择/切换的主题系统来渲染 ToC 前端页面，主题可以是纯 HTML+JS 静态包，也可以是 Next.js SSR 应用。

## 目标

- 后台新增主题管理 — 注册、预览、一键切换主题，**立即生效，不重启**
- 支持两类主题：`static`（纯 HTML+JS）和 `nextjs`（SSR 应用）
- 主题以"插件"形式组织，每个主题是独立目录，放入即用
- 两类主题均可独立部署、独立扩展

---

## 总体架构

```
Browser
  │
  ▼
Nginx (统一入口)
  │  /          → static主题: Go Admin serve 文件  /  nextjs主题: proxy → toc:3000
  │  /admin     → Go Admin (cmd/admin, 端口 8089)
  │  /api/v1    → Go Public API (cmd/api, 端口 8080)
  │
  ▼
PostgreSQL ◄── 所有服务
Redis      ◄── 缓存
```

**两种主题类型的对比：**

| | static (HTML+JS) | nextjs |
|---|---|---|
| 运行时 | 无，纯文件 | Node.js 进程 |
| 渲染位置 | 浏览器 CSR | 服务端 SSR |
| 谁 serve | Go Admin（或 CDN/Nginx 独立部署） | Next.js standalone server |
| 切换生效 | **立即**（换个目录路径） | **秒级**（运行时读 API，有缓存 TTL） |
| 可分离部署 | 是（推到 OSS/CDN 即可） | 是（本就是独立进程） |
| 适合场景 | 简单展示站、活动页 | SEO 敏感、复杂交互 |

**关键数据流（运行时切换，零重启）：**

```
Admin 切换主题 → DB system_config.theme = "xxx"
                    │
static 主题:    Go Admin 每次请求读 DB → serve 对应目录 → 立即生效
nextjs 主题:    Next.js 内存缓存(TTL 60s)过期 → 读 API → 渲染新主题
                或 Admin 切主题时调 Next.js /api/theme/refresh → 立即刷新缓存
```

---

## Phase 1：后台主题管理

### 1.1 数据库变更

**`internal/model/system_config.go`** — 新增字段：
```go
Theme string `gorm:"type:varchar(100);default:default;comment:当前激活的ToC主题名" json:"theme"`
```

**新建 `internal/model/theme.go`** — 主题注册表：
```go
type Theme struct {
    model.Model
    Name        string `gorm:"type:varchar(100);uniqueIndex;comment:主题目录名"`
    Type        string `gorm:"type:varchar(20);default:static;comment:static或nextjs"`
    DisplayName string `gorm:"type:varchar(255);comment:显示名称"`
    Description string `gorm:"type:text;comment:描述"`
    Version     string `gorm:"type:varchar(50);comment:版本"`
    Thumbnail   string `gorm:"type:varchar(500);comment:缩略图URL"`
    Author      string `gorm:"type:varchar(255);comment:作者"`
    Enabled     bool   `gorm:"default:true;comment:是否启用"`
}
```

### 1.2 新建 `internal/admin/theme/` 模块

按项目标准四层模式：

| 层 | 文件 | 职责 |
|----|------|------|
| dto | `request.go`, `response.go` | CreateThemeRequest, ThemeListRequest, ThemeResponse |
| handler | `theme.go`, `routes.go` | Detail/List/Create/Update/Delete/Activate |
| service | `theme.go`, `theme_test.go` | CRUD + 激活时更新 system_config + 通知刷新 |
| repository | `theme.go` | DB 操作 |

**路由：**
```
GET    /api/themes             → List
GET    /api/themes/:id         → Detail
POST   /api/themes             → Create
PUT    /api/themes/:id         → Update
DELETE /api/themes/:id         → Delete
PUT    /api/themes/:id/activate → Activate
```

**Activate 逻辑（service 层）：**
1. 校验 theme 存在且 enabled
2. 调用 SystemConfigRepo 更新 `theme` = 目标 name
3. 如果是 nextjs 类型 → 调 `POST http://toc:3000/api/theme/refresh` 清缓存（失败不影响主流程）

### 1.3 修改现有文件

| 文件 | 改动 |
|------|------|
| `internal/model/system_config.go` | 新增 `Theme` 字段 |
| `internal/admin/system_config/dto/request.go` | 新增 `Theme` 字段 |
| `internal/admin/system_config/dto/response.go` | 新增 `Theme` 字段 + 映射 |
| `internal/admin/system_config/service/system_config.go` | Update 方法赋值 Theme |
| `internal/admin/routes/admin_modules.go` | 新增 Theme handler |
| `internal/admin/routes/routes.go` | 注册 theme 路由 |
| `internal/model/migrate/migrate.go` | 新增 `&model.Theme{}` |
| `internal/admin/common/seeder/` | 种子数据：default 主题 |
| `internal/admin/common/init_database_data.go` | Casbin 策略 |

### 1.4 Admin 前端 — 主题管理页面

| 文件 | 用途 |
|------|------|
| `routes/_authenticated/system/theme.tsx` | 路由 |
| `features/system/theme/theme-list.tsx` | 卡片列表（缩略图/名称/类型/版本，当前激活高亮，一键切换） |
| `features/system/theme/theme-form.tsx` | 新增/编辑表单（选 type: static/nextjs） |
| `features/system/theme/service.ts` | TanStack Query hooks |

---

## Phase 2：主题运行时

### 2.1 static 主题（纯 HTML+JS）

**目录结构：**
```
web/toc/
  themes/
    default/               ← 一个 static 主题就是一个目录
      index.html           ← 入口
      assets/
        style.css
        app.js
      pages/
        post.html          ← 文章详情页模板
        category.html      ← 分类页模板
```

**Go Admin 服务方式：**

在 `routes.go` 中新增路由：
```go
// ToC 前端页面 — 根据 system_config.theme 动态 serve
tocGroup := router.Group("/")
{
    tocGroup.Use(TOCMiddleware(systemConfigRepo))  // 中间件读 DB 拿当前 theme 名
    tocGroup.Static("/", "./web/toc/themes/{theme}/")  // 动态路径
    tocGroup.GET("/post/:slug", tocHandler.PostDetail)
    tocGroup.GET("/category/:slug", tocHandler.CategoryList)
}
```

静态文件直接 serve，HTML 页面里的 JS 调 `/api/v1/posts` 获取数据渲染。切换主题后下一个请求立即走新目录。

**独立部署 static 主题：**
```
今天：  Go Admin serve  web/toc/themes/default/
明天：  推到 OSS，Nginx 配 location / { proxy_pass https://oss-cdn/xxx/; }
       或者 Nginx location / { root /opt/themes/current/; }
```
主题代码零改动，因为 static 主题不依赖 Go 运行时，只依赖 `/api/v1/*` 公共 API。

### 2.2 nextjs 主题

**目录结构：**
```
web/toc/
  next.config.ts
  tsconfig.json
  package.json
  src/
    app/
      layout.tsx                # 根布局 → 从 theme 目录加载 Layout
      page.tsx                  # 首页 → theme.HomePage
      posts/[slug]/page.tsx     # 文章详情 → theme.PostDetail
      categories/[slug]/page.tsx
      api/theme/refresh/route.ts  # 内部接口：刷新主题缓存
    themes/
      default/
        index.ts                # export { layout, pages, components }
        layout.tsx
        pages/
          home.tsx
          post-detail.tsx
          category-page.tsx
        components/
          header.tsx
          footer.tsx
          post-card.tsx
      theme-2/
        index.ts
        ...
    lib/
      api.ts                    # 公共 API 客户端
      theme-loader.ts           # 运行时动态加载主题
      theme-config.ts           # 从 Admin API 读取当前主题名（带缓存）
```

**运行时主题加载（不重启）：**

```typescript
// lib/theme-config.ts — Next.js 运行时定期从 Admin API 读当前主题
let cachedTheme = 'default';
let lastFetch = 0;

export async function getActiveTheme(): Promise<string> {
  if (Date.now() - lastFetch < 60_000) return cachedTheme;  // TTL 60s
  try {
    const res = await fetch(`${ADMIN_URL}/base/config`);
    const data = await res.json();
    cachedTheme = data.data?.theme || 'default';
    lastFetch = Date.now();
  } catch { /* 网络失败用缓存 */ }
  return cachedTheme;
}

// lib/theme-loader.ts — 动态 import 主题
const themes: Record<string, () => Promise<any>> = {
  default: () => import('@/themes/default'),
  'theme-2': () => import('@/themes/theme-2'),
};

export async function loadTheme(name: string) {
  const loader = themes[name] || themes['default'];
  return loader();
}
```

```typescript
// app/page.tsx — 路由只做转发
export default async function HomePage() {
  const themeName = await getActiveTheme();
  const theme = await loadTheme(themeName);
  const ThemeHome = theme.pages.home;
  return <ThemeHome />;
}
```

**刷新缓存接口：**
```typescript
// app/api/theme/refresh/route.ts
export async function POST() {
  // 清掉 getActiveTheme() 的内存缓存，强制下次请求重新读 API
  clearThemeCache();
  return Response.json({ ok: true });
}
```

Admin 切主题时调这个接口 → 缓存清零 → 下一个请求立即渲染新主题。即使调失败，60s TTL 后自动生效。

### 2.3 需扩展的公共 API

ToC 页面基础需求，`cmd/api` 侧需新增：

```
GET /api/v1/categories        → 公开分类列表
GET /api/v1/categories/:slug  → 分类详情
GET /api/v1/tags              → 公开标签列表
```

按项目现有模块模式实现（`internal/category/`, `internal/tag/` 各自新增 handler/service，注册到公共路由）。

---

## Phase 3：部署

### 3.1 开发环境

每天开发时只启动 Go 服务 + Admin 前端，static 主题 Go 直接 serve，无需额外服务。nextjs 主题需单独 `pnpm dev`。

### 3.2 Docker Compose

```yaml
# static 主题模式：api + admin 即可，无需额外容器
# nextjs 主题模式：加一个 toc 服务
toc:
  build:
    context: .
    dockerfile: Dockerfile.toc
  ports:
    - "3000:3000"
  environment:
    - ADMIN_URL=http://admin:8089
    - API_URL=http://api:8080
  depends_on:
    - api
    - admin
```

### 3.3 生产部署 & 水平扩展

**核心前提：所有服务无状态（状态在共享 PostgreSQL + Redis），天然支持水平扩展。**

```
                         Nginx / LB
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
     /           /admin,/api   /api/v1
              │             │             │
    toc-1      admin-1      api-1
    toc-2      admin-2      api-2
    (nextjs)   admin-3      api-3
              │             │             │
              └─────────────┼─────────────┘
                            │
                    ┌───────┴───────┐
                    ▼               ▼
               PostgreSQL         Redis
               (共享)             (共享)
```

**三类服务多机器部署分析：**

| 服务 | 能否多机 | 前提 | 注意事项 |
|------|---------|------|---------|
| **static 主题** | ✅ | 主题文件每台 copy 一份，或用共享卷/NFS，或推到 CDN 统一 serve | 无运行时状态，纯文件 |
| **nextjs 主题** (`toc`) | ✅ | 无状态进程，LB 轮询即可 | redirect 多实例时 `POST /api/theme/refresh` 需广播或通过 Redis pub/sub |
| **Public API** (`cmd/api`) | ✅ | JWT 无状态 + Redis 共享（已满足） | 无 |
| **Admin** (`cmd/admin`) | ✅ | DB + Redis + Casbin(DB adapter) 共享 | ① 定时任务需单例（已有 `jobs` 模块，加分布式锁即可）② 文件上传需共享存储（已用 OSS） |

**多机器部署时 Nginx 动态路由：**

问题：Nginx 需要知道当前激活的是 static 还是 nextjs，才能决定 `/` 打到哪个 upstream。

解决：Admin 暴露一个轻量端点：

```
GET /api/toc-route  →  { "type": "static" }
                    或  { "type": "nextjs", "upstream": "toc" }
```

Nginx 配置两种方式选一：

```nginx
# 方式 A：简单粗暴 — 两套 upstream 都配好，Nginx 用 map + 变量切换
map $toc_type $toc_upstream {
    static    "http://admin:8089";
    nextjs    "http://toc:3000";
}
# $toc_type 通过 auth_request 或 nginx-lua 定期从 /api/toc-route 拉取

# 方式 B：手动切换 — 改 upstream 配置后 nginx -s reload（适用于不频繁切换）
```

> static 主题多机器：每台机器都放主题文件（或用共享卷/NFS），Nginx 轮询。或者推到 CDN，Nginx 直接 `proxy_pass` CDN，Admin 机器完全不需要存主题文件。

**从单机到多机的演进路径：**

```
单机                          多机                             完全分离
─────────────────────────────────────────────────────────────────────
docker-compose           K8s 多副本                      CDN + K8s 多集群
1 台跑全部               3 台各跑 api/admin/toc           按服务粒度独立扩缩
                        共享 DB + Redis                  api x5 / admin x2 / toc x10
                                                        不同 region 各自部署
```

三类服务从一开始就是解耦的（独立进程/独立目录/独立 Dockerfile），所以什么时候想拆机器都行，不用改代码。

---

## 实施顺序

| # | 步骤 | 文件数 |
|---|------|--------|
| 1 | SystemConfig 新增 Theme 字段 | 4 |
| 2 | 创建 Theme model + 迁移注册 | 2 |
| 3 | 创建 theme 模块（dto/handler/service/repo） | 5 |
| 4 | 注册路由和模块，Casbin 策略 | 3 |
| 5 | 种子数据（default static 主题） | 1 |
| 6 | Admin 前端主题管理页面 | 4 |
| 7 | Go Admin serve static 主题（TOCMiddleware + 路由） | 2 |
| 8 | 创建 default static 主题（HTML+JS 样板） | ~6 |
| 9 | 扩展公共 API（categories/tags） | ~6 |
| 10 | 创建 `web/toc/` Next.js 项目 + default nextjs 主题 | ~15 |
| 11 | Docker 部署配置 | 2 |
| 12 | 端到端验证 | — |

---

## 验证方式

1. `go build ./...` — 编译通过
2. `go test ./...` — 测试通过
3. `cd web/admin && pnpm build` — Admin 前端构建通过
4. static 主题：`curl http://localhost:8089/` → 返回主题 HTML
5. 后台切换主题 → `curl http://localhost:8089/` → 立即返回新主题 HTML
6. nextjs 主题：`curl http://localhost:3000/` → SSR 渲染的首页
7. 后台切换 nextjs 主题 → 等 60s 或触发 refresh → 新主题生效

## 风险

| 风险 | 缓解 |
|------|------|
| static 主题动态路由（/post/:slug）需 Go 处理 | 用 Gin 路由 + TOCMiddleware 注入当前主题路径，HTML 模板渲染 |
| nextjs 主题动态加载在 App Router 下有编译限制 | 主题需预注册到 theme-loader 的 import map；后续可改为 Webpack/Turbopack 动态 alias |
| 公共 API 缺 categories/tags 端点 | Phase 2 同步补齐，按项目现有模块模式扩展 |
| 两类主题共存时 Nginx 路由需要知道当前是哪种 | `system_config.theme` + theme.Type 决定；可加一个统一的 `/api/toc-type` 端点给 Nginx 做路由决策 |
