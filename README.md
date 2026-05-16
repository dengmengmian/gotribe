# gotribe

一个支持多项目的 Go CMS 工程，包含用户端 API 和管理后台，基于：

- `Gin`
- `GORM`
- `Viper`
- `Redis`
- `JWT` + Casbin RBAC

## 项目定位

gotribe 支持多项目运行，同一实例可按 `project_id` 隔离数据。包含两套服务，共享同一数据库和基础设施：

| 服务 | 入口 | 端口 | 说明 |
|--- |--- |--- |--- |
| ToC API | `cmd/api` | 8080 | 面向用户端的 REST API（登录、内容浏览、用户行为） |
| Admin | `cmd/admin` | 8088 | 后台管理系统（CRUD、权限、定时任务、资源管理）+ React SPA 前端 |

当前已落地：

- ToC API：认证、用户资料、帖子列表/详情、用户事件上报、示例 CRUD 模块
- Admin：30+ 模块的完整后台管理（文章、分类、标签、评论、广告、菜单、角色权限、配置等）
- Admin 前端：基于 React + Shadcn 的管理界面（`web/admin/`）
- 基础设施：PostgreSQL + Redis、JWT 双域认证、Casbin RBAC、限流、日志/Trace、Docker/K8s 部署模板

## 文档入口

- 架构设计：`docs/architecture.md`
- 新人开发手册：`docs/development-guide.md`
- 业务参考模块：`docs/example-module.md`
- 测试说明：`docs/testing-guide.md`
- 部署说明：`docs/deployment-guide.md`
- 提交检查清单：`docs/checklist.md`
- ToC API 接口契约：`api/openapi.yaml`
- Admin Swagger：启动后访问 `http://localhost:8088/swagger/index.html`

建议阅读顺序：

1. `README.md`
2. `docs/architecture.md`
3. `docs/development-guide.md`

## 快速开始

### 依赖

- PostgreSQL 16+
- Redis 7+
- Go 1.25+
- pnpm（仅构建 Admin 前端时需要）

### 本地配置

复制配置模板，按实际环境填写：

```bash
cp configs/config.yaml.example configs/config.yaml
cp .env.example .env
```

关键配置项：

- `database.*` — PostgreSQL 连接信息
- `redis.*` — Redis 连接信息
- `auth.secret` — ToC API JWT 密钥（至少 32 位）
- `app.default_project_id` — 默认项目 ID
- `admin.jwt.key` — Admin JWT 密钥（至少 32 位）

### 启动（本地 Go 直接运行）

```bash
# ToC API
make run

# Admin 后台（另一个终端）
make run-admin
```

### 启动（Docker Compose 一键全套）

```bash
cp .env.example .env
# 编辑 .env，设置 JWT 密钥
make dev-up
```

启动后访问：

| 服务 | 地址 |
|--- |--- |
| ToC API | http://localhost:8080 |
| Admin Web | http://localhost:8088/admin |
| Admin Swagger | http://localhost:8088/swagger/index.html |

### 常用命令

```bash
make tidy          # 整理依赖
make fmt           # 格式化代码
make vet           # 静态检查
make test-unit     # 单元测试
make build         # 构建 ToC API 二进制
make build-admin   # 构建 Admin 二进制
make build-admin-web  # 构建 Admin 前端
```

## 核心目录

```text
gotribe/
├── api/                     # OpenAPI 契约
│   └── openapi.yaml
├── cmd/
│   ├── api/                 # ToC API 入口
│   └── admin/               # Admin 后台入口
├── configs/                 # 配置文件与 RBAC 模型
├── deployments/             # K8s 部署模板
├── docs/                    # 项目文档
├── internal/
│   ├── admin/               # Admin 模块（30+ 业务域）
│   │   ├── ad/              #   广告管理
│   │   ├── admin_user/      #   管理员管理
│   │   ├── auth/            #   Admin 登录认证
│   │   ├── bootstrap/       #   Admin 应用装配
│   │   ├── category/        #   分类管理
│   │   ├── column/          #   栏目管理
│   │   ├── comment/         #   评论管理
│   │   ├── common/          #   Admin 公共组件、Seeder
│   │   ├── config/          #   系统配置
│   │   ├── feedback/        #   反馈管理
│   │   ├── index/           #   仪表盘
│   │   ├── job/             #   定时任务管理
│   │   ├── jobs/            #   定时任务引擎
│   │   ├── menu/            #   菜单管理
│   │   ├── middleware/      #   Admin 中间件（JWT/Casbin/操作日志）
│   │   ├── migration/       #   数据库迁移
│   │   ├── operation_log/   #   操作日志
│   │   ├── point/           #   积分管理
│   │   ├── post/            #   文章管理
│   │   ├── project/         #   项目管理
│   │   ├── resource/        #   资源管理
│   │   ├── role/            #   角色权限
│   │   ├── routes/          #   Admin 路由注册
│   │   ├── system_config/   #   系统设置
│   │   ├── tag/             #   标签管理
│   │   └── user/            #   用户管理
│   ├── auth/                # ToC 认证模块
│   ├── bootstrap/           # ToC API 应用装配
│   ├── core/                # 公用基础设施
│   │   ├── cache/           #   Redis 缓存
│   │   ├── config/          #   配置加载（Viper）
│   │   ├── constant/        #   公共常量
│   │   ├── database/        #   GORM 与事务
│   │   ├── errs/            #   统一错误定义
│   │   ├── logger/          #   日志
│   │   ├── middleware/       #   公用中间件
│   │   ├── response/        #   统一返回格式
│   │   └── util/            #   公用工具
│   ├── example/             # 标准业务参考模块
│   ├── health/              # 健康检查
│   ├── model/               # GORM Model 定义
│   ├── post/                # ToC 帖子模块
│   ├── profile/             # ToC 用户资料模块
│   ├── request/             # 请求参数绑定
│   ├── tag/                 # ToC 标签模块
│   └── user_event/          # ToC 用户事件模块
├── migrations/              # 数据库迁移 SQL
├── pkg/                     # 已废弃（内容已迁移至 internal/core）
├── tests/                   # 基准测试
└── web/
    └── admin/               # Admin React SPA 前端
```

## ToC API 主要接口

| 方法 | 路径 | 说明 |
|--- |--- |--- |
| GET | `/livez` | 存活检查 |
| GET | `/readyz` | 就绪检查 |
| GET | `/version` | 版本信息 |
| POST | `/api/v1/auth/login` | 登录 |
| POST | `/api/v1/auth/refresh` | 刷新 token |
| POST | `/api/v1/auth/logout` | 登出 |
| GET | `/api/v1/me` | 当前用户资料 |
| PATCH | `/api/v1/me` | 修改资料 |
| POST | `/api/v1/me/password` | 修改密码 |
| GET | `/api/v1/posts` | 帖子列表 |
| GET | `/api/v1/posts/:id` | 帖子详情 |
| POST | `/api/v1/user-events` | 用户事件上报 |

### 使用示例

```bash
# 健康检查
curl http://localhost:8080/livez

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Project-ID: 1" \
  -d '{"identity":"user","password":"pass"}'

# 获取当前用户
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer <access_token>"

# 帖子列表
curl "http://localhost:8080/api/v1/posts?project_id=1&page=1&per_page=10"
```

## 开发约定

- `handler` 只处理 HTTP 输入输出，不操作数据库或缓存
- `service` 负责业务逻辑、缓存策略、事务边界
- `repository` 只负责数据访问
- 参数绑定通过 `internal/request`
- 配置统一通过 `internal/core/config`
- 日志统一通过 `internal/core/logger`
- 每个模块的路由注册在各自的 `handler/routes.go`
- Admin 模块统一采用 `dto/request.go` + `dto/response.go` 分层
- 模块间调用通过 service 层，禁止跨模块直接访问 repository

## 环境行为

- `development`：完整请求日志、SQL 日志、Debug 级业务日志
- `production`：仅错误日志、Error 级 SQL、Error 级业务日志

## 扩展新模块

按业务域新增目录，以 Admin 模块为例：

```text
internal/admin/<module>/
├── dto/
│   ├── request.go
│   └── response.go
├── handler/
│   ├── <module>.go
│   └── routes.go
├── repository/
│   └── <module>.go
└── service/
    ├── <module>.go
    └── <module>_test.go
```

ToC API 模块类似，去掉 `dto/` 下的 `request.go`/`response.go`，改用 `internal/<module>/dto/`。

开发步骤详见 `docs/development-guide.md` 和 `docs/example-module.md`。
