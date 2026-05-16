# Architecture Guide

## 目标

这份文档用于帮助新人快速理解 `gotribe` 的整体架构、设计边界和日常开发方式。

这个项目不是后台管理系统，也不是微服务平台，而是一个：

- 面向用户端的 Go REST API 单体服务
- 以 `Gin + GORM + Viper + Redis + JWT` 为基础技术栈
- 以云原生部署为前提设计
- 以生产环境可维护性为优先目标

## 为什么这么设计

这套结构有几个明确目标：

- 让新人进入项目后，能快速找到入口、业务、数据访问和基础设施代码
- 让业务模块自然扩展，而不是把所有 handler、service、repository 堆在公共目录里
- 让登录态、限流、缓存、日志、配置这些横切能力统一管理
- 让路由、业务逻辑、数据库访问各自职责清晰，避免后续代码发散
- 让未来新增业务模块时尽量少改已有代码

项目遵循的核心原则是：

- 按业务域组织代码，而不是只按技术层堆目录
- `handler -> service -> repository` 单向依赖
- HTTP 框架只停留在 handler 和 middleware 层
- service 层不依赖 Gin
- repository 只负责数据访问，不写 HTTP 语义
- 跨模块协作优先通过 service contract 或显式能力接口完成

## 顶层目录

```text
gotribe/
├── api/
├── cmd/
├── configs/
├── deployments/
├── docs/
├── internal/
├── pkg/
└── tests/
```

各目录职责如下：

- `api/`
  放 OpenAPI 契约文档，接口联调以这里为准。

- `cmd/api/`
  程序入口，只负责加载配置、启动应用、优雅退出。

- `configs/`
  配置样例文件。运行期优先使用环境变量。

- `deployments/`
  Docker、Compose、Kubernetes 等部署文件。

- `docs/`
  项目架构文档、开发约定、运行说明等长期文档。

- `internal/`
  项目私有代码。业务模块和基础设施都放在这里。

- `pkg/`
  极少量跨模块通用工具。这里要克制使用，不能变成杂货铺。

- `tests/`
  集成测试和端到端测试。

## internal 目录设计

```text
internal/
├── auth/
├── bootstrap/
├── cache/
├── config/
├── database/
├── errs/
├── health/
├── logger/
├── middleware/
├── post/
├── profile/
├── request/
├── response/
├── tag/
└── user_event/
```

### 基础设施目录

- `bootstrap/`
  应用装配层。负责初始化基础设施、路由、中间件和各模块依赖。

- `config/`
  配置加载层。Viper 只允许在这里出现，其他代码只使用 `Config struct`。

- `database/`
  GORM 初始化、连接池配置、事务上下文封装。

- `cache/`
  Redis 客户端、key 构造、缓存存储封装。

- `middleware/`
  Gin 中间件集合，包括 request id、日志、recover、JWT、项目隔离、限流、当前用户加载等。

- `logger/`
  统一业务日志入口。业务代码不要直接散落 `log.Printf`。

- `errs/`
  统一业务错误码和错误包装。

- `request/`
  统一请求绑定和参数校验，包括中英文 validator 翻译。

- `response/`
  统一响应输出格式。

### 业务模块目录

- `auth/`
  登录、刷新 token、登出、JWT、refresh token 管理。

- `profile/`
  当前用户资料读取、更新、改密码。

- `post/`
  帖子列表和详情读取。

- `tag/`
  标签模型和标签仓储。当前不暴露独立 API，但作为独立实体维护。

- `user_event/`
  用户行为上报。

- `health/`
  `/livez` 和 `/readyz` 探针。

## 一次请求是怎么走的

典型请求链路如下：

```text
HTTP Request
-> Gin Router
-> Global Middleware
-> Route Middleware
-> Handler
-> Service
-> Repository
-> DB / Redis
```

每层职责：

- `Router`
  负责路由分组，不写业务逻辑。

- `Middleware`
  负责横切逻辑，比如请求追踪、日志、JWT、限流、recover。

- `Handler`
  负责请求绑定、响应返回、把 HTTP 参数传给 service。

- `Service`
  负责业务规则、缓存策略、事务边界、跨模块协作。

- `Repository`
  负责数据库读写和查询拼装。

## 路由设计

总路由文件在：

- `internal/bootstrap/router.go`

它只保留三类核心职责：

- 初始化全局中间件
- 定义路由层级
- 调用各模块自己的 `RegisterRoutes`

当前路由层级分成：

- `public`
  不需要登录，适合公共读接口。

- `secured`
  需要 JWT，默认只透传轻量身份信息。

- `currentUser`
  在 `secured` 基础上进一步加载当前完整用户对象，只给确实需要完整用户资料的接口使用。

这套设计的原因：

- 大多数已登录接口只需要 `user_id`，没必要默认查库
- 少量依赖完整用户信息的接口，可以按需挂 `CurrentUser` 中间件
- 路由权限边界清楚，后续模块扩展容易

## JWT 和当前用户设计

JWT 中间件只做这些事情：

- 校验 access token
- 提取 `user_id`、`username`、`project_id`
- 写入 Gin context
- 写入 `request.Context()` 供日志和 service 使用

也就是说，JWT 中间件默认只透传轻量身份，不直接查数据库。

如果接口需要完整用户资料，再挂 `CurrentUser` 中间件：

- 先从 JWT claims 中获取身份
- 再调用 profile service 加载当前用户
- 把完整用户对象放进 context

这样能兼顾性能和可维护性。

## 数据访问设计

repository 不负责业务逻辑，只负责数据访问。

这里有几个关键约定：

- repository 不直接返回 HTTP 错误
- repository 不依赖 Gin
- repository 只操作自己的实体和表
- repository 使用统一的 `TransactionManager.DB(ctx)` 取 DB 句柄

注意：

- `TransactionManager.DB(ctx)` 不代表一定开了事务
- 它的含义是“取当前上下文对应的 DB 或事务句柄”
- 真正显式开事务的是 `WithinTransaction(...)`

这能保证：

- 普通读请求不强制走事务
- 真正需要事务的写流程可以共享一套 repository 代码

## 为什么不在本服务做 migration

数据库 schema 由外部后台系统统一维护，本服务只维护 GORM model 映射。

这样做的原因：

- schema 变更职责集中
- 避免多个服务同时维护一份数据库迁移
- 本服务只专注数据读写和 API 逻辑

因此：

- `model` 的职责是表结构映射
- `repository` 的职责是数据访问
- schema 变更不在本仓库执行

## 缓存和 Redis 设计

Redis 在本项目里承担这些职责：

- refresh token 存储
- 接口限流
- 用户资料缓存
- 帖子详情缓存

关键设计点：

- Redis key 统一通过 `KeyBuilder` 构造
- key 前缀来自 `app.name`
- key builder 是实例依赖，不用全局变量
- refresh token 的索引带 `project_id + user_id` 维度

这样做是为了：

- 避免全局可变状态污染
- 适配多项目隔离
- 便于未来拆分 worker 或多入口程序

## 日志设计

项目有三类日志：

- HTTP 请求日志
- GORM SQL 日志
- 业务日志

日志策略：

- `development`
  输出完整请求日志、完整 SQL 日志、`Debug/Info/Warn/Error`

- `production`
  HTTP 只记录错误请求
  GORM 只输出 error 级 SQL
  业务日志只保留 `Error`

业务代码统一使用：

- `logger.Debug(ctx, ...)`
- `logger.Info(ctx, ...)`
- `logger.Warn(ctx, ...)`
- `logger.Error(ctx, ...)`

不要在业务代码里到处直接写 `log.Printf(...)`。

## 配置设计

配置统一从 `internal/config` 加载。

原则：

- 使用 `Config struct`
- 环境变量优先
- 允许本地配置文件辅助开发
- 业务代码不能直接调用 Viper

这样能保证：

- 配置来源清晰
- 测试更容易构造
- 配置项变更更容易收敛

## 参数校验设计

参数校验统一使用：

- Gin binding
- go-playground validator
- universal-translator
- `Accept-Language` / `X-Language`

也就是说：

- JSON 请求优先 `request.BindJSON`
- Query 参数优先 `request.BindQuery`
- 不鼓励在 handler 里散写 `strconv.Atoi` 或手工解析 query

这样能保证：

- 校验风格统一
- 错误结构统一
- 支持中英文校验提示

## 现在的业务模块为什么这样拆

### auth

独立拆出认证模块，是因为认证通常有单独的生命周期和风险边界：

- token
- refresh token
- 密码校验
- 登出
- 第三方登录扩展

### profile

用户资料和认证是两个不同领域：

- auth 关心“你是谁”
- profile 关心“你的资料是什么”

这两者不应该混成一个大 service。

### post + tag

`tag` 虽然当前不暴露独立 API，但它是独立实体，不应该塞在 `post model` 里。

当前做法是：

- `post` 模块负责帖子 API
- `tag` 模块负责标签模型和标签查询
- `post` 通过标签仓储协作，而不是自己长期吸收 tag 职责

这样后面如果 banner、resource、专题等也要关联 tag，就不会重复定义。

## 新人开发应该怎么开发

### 新增接口的标准步骤

1. 先确定接口属于哪个业务模块。
2. 在模块的 `dto/` 中定义请求和响应结构。
3. 在模块的 `handler/` 中编写 handler。
4. 在模块的 `service/` 中编写业务逻辑。
5. 在模块的 `repository/` 中编写数据访问。
6. 在模块的 `handler/routes.go` 中注册路由。
7. 更新 `api/openapi.yaml`。
8. 补最少必要的测试。

### 新增模块的标准步骤

例如以后新增 `order` 模块，建议目录直接按模块组织：

```text
internal/order/
├── dto/
├── handler/
├── model/
├── repository/
└── service/
```

然后：

- 在 `bootstrap/module_builders.go` 装配模块依赖
- 在 `bootstrap/provider.go` 汇总 `Infra + Modules`
- 在 `bootstrap/router.go` 调用 `order handler` 的 `RegisterRoutes`
- 不要把订单逻辑塞进 `profile`、`auth` 或 `post`

### 什么时候挂 `CurrentUser`

默认不要挂。

只有接口确实需要完整当前用户对象时才挂，比如：

- `GET /me`
- 某些强依赖当前资料快照的能力

只需要 `user_id` 的接口，直接走 `secured` 即可。

### 什么时候开事务

只有在一次业务动作涉及多个写操作、必须保证原子性时，才在 service 层开事务。

比如未来订单场景：

- 创建订单
- 预占积分
- 写订单流水

这类流程应该由 service 控制事务边界。

不要：

- 在 handler 里开事务
- 在 repository 里自己偷偷开事务

## 开发时要遵守的几条硬规则

- handler 不直接操作 GORM 或 Redis
- service 不依赖 Gin
- repository 不跨模块直接改别人的表
- 配置读取只经过 `internal/config`
- 业务日志只经过 `internal/logger`
- 参数绑定优先走 `internal/request`
- 路由注册只放在模块 `handler/routes.go`
- 需要完整当前用户时才挂 `CurrentUser`

## 适合新人先看的文件

建议新人按这个顺序看：

1. `/gotribe/README.md`
2. `/gotribe/internal/bootstrap/router.go`
3. `/gotribe/internal/bootstrap/provider.go`
4. `/gotribe/internal/bootstrap/module_builders.go`
5. `/gotribe/internal/middleware/jwt.go`
6. `/gotribe/internal/profile/handler/profile.go`
7. `/gotribe/internal/profile/service/profile.go`
8. `/gotribe/internal/profile/repository/profile.go`

这样最容易理解整个请求从路由到数据库是怎么走的。

## API 文档策略

本项目对 ToC API 和 Admin 接口采用**双轨文档策略**，二者刻意不统一。

| 端 | 策略 | 文档源 |
|---|---|---|
| ToC API（`cmd/api`）| contract-first | `api/openapi.yaml` 手写 OpenAPI 3.0 |
| Admin（`cmd/admin`）| code-first | handler 上 `swaggo/swag` 注解，`swag init` 生成 `docs/admin/swagger` |

### 为什么双轨

- **ToC API 对外**，要给前端、第三方、SDK 生成器用。手写 yaml 让契约**先于代码**，不被实现细节绑架；变更要慎重，要有审计。
- **Admin 只给自家前端用**，前后端可以一起改。注解就近写更轻，改代码 + `swag init` 一步到位。
- 这种"双轨"在 Stripe（公开 API vs 内部 Dashboard）、GitHub Enterprise（REST API vs admin panel）、Cloudflare（control plane vs edge）等同时维护对外契约和内部后台的项目里是常见做法。

### 改 ToC API 的流程

1. 先改 `api/openapi.yaml` 的契约
2. 让 handler / service / repository 实现追上契约
3. `go test ./tests/integration/...` 跑通

ToC handler 不要加 swaggo 注解，避免 yaml 与注解互相漂移。

### 改 Admin 接口的流程

1. 改 handler，更新或新增 `@Summary` / `@Tags` / `@Param` / `@Success` / `@Router` 注解
2. 跑 `swag init`（参见 `Makefile`）重新生成 `docs/admin/swagger`
3. `go test ./tests/admin/integration/...` 跑通

Admin 不要新增 OpenAPI yaml，避免与 swag 输出互为事实源。

## 当前阶段的结论

这套架构适合当前项目的原因是：

- 业务规模还在可控范围内，单体 API 是合适的
- 目录按业务域拆分，后面扩展不会太痛
- JWT、限流、缓存、日志、配置已经统一收口
- 生产和开发环境的行为已经做了区分
- 新人进入项目后能较快建立整体认知

如果后面继续扩展业务，优先按现在这套约定自然演进，不要过早引入重型 DDD、CQRS 全家桶或微服务拆分。
