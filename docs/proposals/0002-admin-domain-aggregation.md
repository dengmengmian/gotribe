# Proposal 0002: Admin 模块按业务域聚合

## 元数据

| 项 | 值 |
|---|---|
| 状态 | Draft |
| 类别 | 架构（CLAUDE.md §2.3）|
| 创建 | 2026-05-10 |
| 触发任务 | ARCHITECTURE_REPORT.md 中度问题 #4、第二轮评审 #9 |

## 摘要

把 `internal/admin/` 下 28 个平铺业务包按业务域聚合为 7 个域目录，与前端 `web/admin/src/features/` 域划分对齐。前后端同构（isomorphic monorepo）是 mono-repo 业内标准做法（参考 Next.js / Nx workspace / Bazel monorepo 实践）。

**不考虑历史兼容**——大批量 `git mv` + import 路径重写，但 API 路径、数据库、配置全部不变。零用户感知。

## 问题

| 现象 | 证据 |
|---|---|
| 28 个平铺业务包 | `ls internal/admin/`（不算 `bootstrap/`、`routes/`、`common/`、`middleware/`、`migration/`、`jobs/` 等基础设施）|
| import 别名块 50+ 行 | `internal/admin/routes/admin_modules.go` |
| 强关联模块独立平铺 | `ad/` + `ad_scene/`；`config/` + `system_config/`；`job/`（HTTP）+ `jobs/`（基础设施，包名冲突）|
| dto 命名风格漂移 | 部分用 `request.go`/`response.go`，部分散布 |
| 与前端不对齐 | 前端 `features/` 9 域，后端 28 包 |

## 目标

1. 后端业务模块对齐前端 7 个有后端对应的域：`auth / content / dashboard / operation / promotion / settings / system`
2. 强关联模块（如 ad + ad_scene）合并到同一域目录
3. import 别名块从 50+ 行降至 ~15 行
4. 解决 `job/` vs `jobs/` 包名冲突
5. 路由注册由"逐模块"改为"逐域"

## 非目标

- 不改 handler / service / repository 三层结构
- 不改业务行为
- 不改数据模型 / migrations
- 不改 swagger 路由 path
- 不重命名公共 API URL
- 不改前端代码
- 不动基础设施目录（bootstrap、routes、common、middleware、migration、jobs）

## 范围

### 改

| 路径 | 动作 |
|---|---|
| `internal/admin/<28 业务模块>/` | `git mv` 到对应域 |
| `internal/admin/routes/admin_modules.go` | import 路径与构造函数调用更新 |
| `internal/admin/routes/routes.go` | 路由注册按域分组 |
| `internal/admin/bootstrap/app.go` | 引用更新 |
| swagger 注解 | 路径前缀如有引用旧路径需更新 |
| 测试文件 | import 路径更新 |

### 不改

- API URL（前端调用方零感知）
- 数据库 schema、migrations
- 业务 service / repository / dto 内部代码
- `internal/admin/common/`、`bootstrap/`、`routes/`、`middleware/`、`migration/`、`jobs/`（基础设施保留原位）
- `internal/admin/casbin/`（独立的 RBAC 模块，留原位）

## 设计

### 域映射表

| 业务域 | 现有模块 | 新路径 |
|---|---|---|
| auth | `auth` + `admin_user` + `role` + `menu` + `api` | `internal/admin/auth/`（注：本 proposal 0002 落地时若 0001 未先合并，此处 `auth/` 与 0001 设计冲突，需协调）|
| content | `post` + `category` + `tag` + `column` + `comment` | `internal/admin/content/` |
| promotion | `ad` + `ad_scene` + `point` | `internal/admin/promotion/` |
| operation | `operation_log` + `job`（jobs 基础设施留原位）| `internal/admin/operation/` |
| settings | `config` + `system_config` + `project` + `feedback` | `internal/admin/settings/` |
| system | `resource` + `util/upload` | `internal/admin/system/` |
| dashboard | `index` | `internal/admin/dashboard/` |

总计 7 域聚合 27 个业务模块。剩 `casbin/` 留原位（RBAC 基础设施）。

### 域内结构

每个域内部按子模块继续保持 `handler/service/repository/dto`：

```
internal/admin/content/
├── post/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── dto/
├── category/
│   └── ...
├── tag/
├── column/
└── comment/
```

域顶层不放代码，只是命名空间。

### 路由注册简化

每个域提供一个 `BuildModules(infra) DomainModules` 函数：

```go
// internal/admin/content/modules.go
package content

type Modules struct {
    Post     *post.Handler
    Category *category.Handler
    Tag      *tag.Handler
    Column   *column.Handler
    Comment  *comment.Handler
}

func BuildModules(infra *Infra) *Modules { ... }

func (m *Modules) RegisterRoutes(r gin.IRoutes, auth gin.HandlerFunc) {
    m.Post.RegisterRoutes(r, auth)
    m.Category.RegisterRoutes(r, auth)
    // ...
}
```

`admin_modules.go` 简化后：

```go
// 旧（50+ 行 import + 散落构造）
adHandler := adhandler.NewHandler(adService)
adSceneHandler := adscenehandler.NewHandler(adSceneService)
// ...

// 新（~15 行 import）
type AdminModules struct {
    Auth       *auth.Modules
    Content    *content.Modules
    Promotion  *promotion.Modules
    Operation  *operation.Modules
    Settings   *settings.Modules
    System     *system.Modules
    Dashboard  *dashboard.Modules
}

func BuildAdminModules(infra *Infra) *AdminModules {
    return &AdminModules{
        Auth:      auth.BuildModules(infra),
        Content:   content.BuildModules(infra),
        // ...
    }
}
```

### 迁移策略

按域逐个迁移，每个域操作：

1. 创建域目录 `internal/admin/{domain}/`
2. `git mv` 子模块到域目录
3. 改 `package` 声明（保留原包名，path 改即可）
4. `gofmt -r` + `goimports` 批量改 import
5. `go build ./...` 通过
6. 跑相关 service/handler 测试
7. 写域级 `BuildModules`
8. 改 `admin_modules.go` 调用方
9. 跑 `tests/admin/integration`

7 域全部完成后：
- 全量 test
- 重新生成 swagger（`swag init`）
- 更新 `docs/architecture.md` 目录树

### 域顺序

按依赖低到高：

1. `dashboard`（依赖最少，先练手）
2. `system`（resource + upload）
3. `settings`（config + system_config + project + feedback）
4. `operation`（operation_log + job）
5. `promotion`（ad + ad_scene + point）
6. `content`（post + category + tag + column + comment）
7. `auth`（admin_user + role + menu + api，**等 0001 先合并**）

## 影响

### 代码层面

- 影响 ~28 个目录
- ~200+ 文件 import path 改动
- routes 注册代码减少 ~30 行
- import 别名减少 ~40 行
- git history：大量 rename（git 自动识别为 rename，diff 易读）

### 用户层面

- API URL 不变
- 配置不变
- 行为不变
- **零用户感知**

### 工具链层面

- swagger 注解需重新生成（`make swag` 或 `swag init`）
- `golangci-lint` path-based 规则需检查
- IDE 索引会重建

## 验证

### 自动化

| 检查 | 通过条件 |
|---|---|
| `go build ./...` | 无错 |
| `go vet ./...` | 无错 |
| `go test ./...` | 全绿，含 `tests/admin/integration` 与 `tests/integration` |
| swagger 生成 | `swag init` 无错 |

### 手动

- 启动 admin 端，每个聚合域至少调一个接口确认可达：
  - content：列文章
  - promotion：列广告
  - settings：读配置
  - system：列资源
  - operation：列操作日志
  - dashboard：仪表盘加载
  - auth：登录
- swagger UI 全部接口可见、可调

## 回滚

- `git revert` 单个聚合 commit
- git 自动识别 rename，回滚干净

## 备选方案

### A: 保留 28 包，只改 import 别名

- ✅ 改动小
- ❌ 治标不治本
- ❌ 模块爆炸问题没解决

**决定：拒绝。**

### B: DDD bounded context（不参考前端）

- ✅ 纯后端建模
- ❌ 前后端语义漂移
- ❌ 后续新增功能两边都要思考边界

**决定：拒绝。**前后端同构是 mono-repo 业内标准。

### C: 拆成多个 service 进程

- ✅ 彻底解耦
- ❌ 远超本任务范围，运维代价高

**决定：拒绝。**延后。

### D: 不动后端，让前端改成跟后端一致

- ❌ 前端 `features/` 已经是合理结构，不应倒退
- ❌ 不解决根本问题

**决定：拒绝。**

## 开放问题

1. **`internal/admin/jobs/` 归属**：
   - 决定：留原位（`internal/admin/jobs/` 是 cron 调度基础设施）；HTTP 接口 `job/` 归入 `operation/`
2. **`internal/admin/casbin/` 归属**：
   - 决定：留原位（RBAC 基础设施，与 `common/` 同级）
3. **`internal/admin/util/upload/` 归属**：
   - 决定：归入 `system/`
4. **0001 与 0002 顺序**：
   - 决定：**先 0001 后 0002**。0001 中 `internal/admin/auth/` 会被删除，0002 的 `auth` 域只聚合 `admin_user / role / menu / api`，避免冲突。

## 实现完成定义

- [ ] 7 个域目录全部建立
- [ ] 27 个业务模块全部 `git mv` 到位
- [ ] `admin_modules.go` import 块 ≤ 15 行
- [ ] `routes.go` 路由注册按域分组
- [ ] swagger 重新生成无错
- [ ] 全量测试通过
- [ ] `docs/architecture.md` 目录树更新
- [ ] `web/admin/` 前端零修改
