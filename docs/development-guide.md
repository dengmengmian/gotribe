# Development Guide

## 目标

这份文档面向刚接手项目的新同学，重点回答这些问题：

- 本地怎么把项目跑起来
- 新需求应该从哪里开始写
- 新增接口和新增模块的标准步骤是什么
- 开发完成后要怎么自测和提交

这份文档偏实操，不重复讲太多架构原理。

如果你想先理解“为什么这么设计”，请先看：

- `docs/architecture.md`

## 第一天应该怎么开始

推荐按下面顺序完成：

1. 阅读 `README.md`
2. 阅读 `docs/architecture.md`
3. 看 `internal/bootstrap/router.go`
4. 看 `internal/bootstrap/provider.go`
5. 看 `internal/bootstrap/module_builders.go`
6. 选一个现有模块通读一遍，推荐 `internal/profile`
7. 本地跑通 `make test`
8. 本地跑起服务并用接口测试工具调一次登录和 `/me`

这样能最快建立整体认知。

## 本地开发准备

### 依赖

本地需要准备：

- Go `1.25.x`
- PostgreSQL
- Redis

### 配置

项目使用 `Viper + 环境变量`。

环境变量前缀统一是：

```bash
GOTRIBE_
```

开发时至少要确认这些配置项：

```bash
GOTRIBE_APP_ENV=development
GOTRIBE_APP_DEFAULT_PROJECT_ID=your-project-id
GOTRIBE_AUTH_SECRET=your-long-random-secret
GOTRIBE_DATABASE_HOST=localhost
GOTRIBE_DATABASE_PORT=5432
GOTRIBE_DATABASE_USERNAME=develop
GOTRIBE_DATABASE_PASSWORD=your-password
GOTRIBE_DATABASE_DATABASE=develop
GOTRIBE_REDIS_ADDR=127.0.0.1:6379
```

注意：

- `auth.secret` 不能使用占位值
- `/api/v1/*` 请求需要项目隔离信息，除非已经配置默认项目 ID
- 项目不在本服务内做 migration，数据库 schema 由外部统一维护

### 常用命令

```bash
# 依赖管理
make tidy

# 代码格式化
make fmt

# 只跑单元测试（快，3-5 秒）
make test-unit

# 跑集成测试（需要 Docker，自动起 PG + Redis）
make integration-test

# 跑全部测试（单元 + 集成）
make test

# 上线前全量检查（格式化 + 代码检查 + 全量测试）
make pre-deploy

# 本地构建
make build

# 交叉编译 Linux 版本
make build-linux

# 本地启动服务
make run
```

推荐开发习惯：

1. 拉代码后先执行 `make test-unit`
2. 开发过程中随手执行 `make fmt`
3. 涉及 Docker 全链路时执行 `make integration-test`
4. 提交前至少执行一次 `make test`
5. 上线前执行 `make pre-deploy`

如果你需要本地确认观测能力，也可以顺手检查：

- `GET /metrics`
- 响应头里的 `X-Trace-ID`
- 日志里的 `trace_id`

## 本地启动与验证

### 启动服务

```bash
make run
```

### 基础健康检查

服务启动后优先检查：

- `GET /version`
- `GET /livez`
- `GET /readyz`
- `GET /metrics`

如果 `readyz` 失败，通常优先排查：

- PostgreSQL 是否能连通
- Redis 是否能连通
- 配置是否正确

### 接口联调建议顺序

推荐按这个顺序调接口：

1. `POST /api/v1/auth/login`
2. `GET /api/v1/me`
3. `PATCH /api/v1/me`
4. `GET /api/v1/posts`
5. `GET /api/v1/posts/{postID}`
6. `POST /api/v1/user-events`

这样可以顺便把：

- JWT
- 项目隔离
- Redis
- 当前用户
- 帖子读取

这些核心链路都验证一遍。

## 看懂一个模块的推荐方法

以 `profile` 模块为例，推荐按这个顺序看：

1. `internal/profile/handler/routes.go`
2. `internal/profile/handler/profile.go`
3. `internal/profile/dto/*.go`
4. `internal/profile/service/profile.go`
5. `internal/profile/repository/profile.go`
6. `internal/profile/model/user_profile.go`

理解顺序是：

- 先看暴露了哪些接口
- 再看 handler 怎么收参数和回响应
- 再看 service 怎么处理业务
- 最后看 repository 怎么落到数据库

## 新增一个接口怎么写

这里以“在已有模块里新增一个接口”为例。

### 第一步：确认接口属于哪个模块

先问自己：

- 这个接口属于哪个业务域
- 是放到现有模块，还是应该新建模块

例如：

- 当前用户资料相关，属于 `profile`
- 帖子读取相关，属于 `post`
- 用户行为上报，属于 `user_event`

不要把不同领域的逻辑塞进一个现有 service。

### 第二步：先定义 DTO

在模块的 `dto/` 中新增请求和响应结构。

例如：

```text
internal/profile/dto/
```

约定：

- 请求结构和响应结构分开定义
- 参数校验优先用 `binding` tag
- 不要直接把数据库 model 暴露给 API

### 第三步：写 handler

在模块的 `handler/` 中新增 HTTP 处理器。

约定：

- handler 只做请求绑定、调用 service、输出响应
- 不直接操作 GORM
- 不直接操作 Redis
- 参数绑定统一走 `internal/request`

例如：

```go
var req dto.UpdateSomethingRequest
if err := request.BindJSON(c, &req); err != nil {
    response.Error(c, err)
    return
}
```

### 第四步：写 service

业务逻辑统一写在 `service/`。

约定：

- service 不依赖 Gin
- service 负责业务规则
- service 负责缓存策略
- 需要事务时，由 service 控制事务边界

### 第五步：写 repository

数据访问统一写在 `repository/`。

约定：

- repository 只负责数据读写
- repository 不返回 HTTP 语义
- repository 不跨模块直接改别人表

### 第六步：注册路由

把新接口注册到模块自己的：

```text
handler/routes.go
```

然后由 `internal/bootstrap/router.go` 统一挂接模块路由。

不要把业务路由继续堆回总路由文件。

### 第七步：更新 OpenAPI

接口改动完成后，要同步更新：

```text
api/openapi.yaml
```

接口契约必须和实际行为保持一致。

## 新增一个模块怎么写

如果需求已经不是现有模块的自然扩展，就新建模块。

完整标准示例请优先参考：

- `internal/example`
- `docs/example-module.md`

例如以后新增 `order`，推荐目录：

```text
internal/order/
├── dto/
├── handler/
├── model/
├── repository/
└── service/
```

标准步骤：

1. 新建模块目录结构
2. 定义 `dto`
3. 编写 `handler`
4. 编写 `service`
5. 编写 `repository`
6. 在 `handler/routes.go` 里提供 `RegisterRoutes(...)`
7. 在 `internal/bootstrap/module_builders.go` 增加模块 builder
8. 由 `internal/bootstrap/provider.go` 汇总到 `Providers`
9. 在 `internal/bootstrap/router.go` 挂接路由
10. 更新 `api/openapi.yaml`

## 路由应该挂在哪一层

项目目前有三层路由：

- `public`
- `secured`
- `currentUser`

挂载原则：

- 公共读接口放 `public`
- 只需要 JWT 身份信息的接口放 `secured`
- 依赖完整当前用户对象的接口放 `currentUser`

默认优先选 `secured`，不要一上来就挂 `currentUser`。

原因是：

- 大多数接口只需要 `user_id`
- `currentUser` 会额外查一次用户信息
- 生产上不应该把每个登录接口都变成“默认查库”

## 当前用户相关开发怎么做

如果 handler 只需要 `user_id`，直接用：

```go
userID, ok := middleware.GetUserID(c)
```

如果 handler 需要完整当前用户对象，并且路由已经挂了 `CurrentUser`，就用：

```go
currentUser, ok := middleware.GetCurrentUser(c)
```

不要在 handler 里自己重复解析 JWT。

## 日志应该怎么写

业务日志统一走：

```go
logger.Debug(ctx, ...)
logger.Info(ctx, ...)
logger.Warn(ctx, ...)
logger.Error(ctx, ...)
```

推荐场景：

- `Info`
  关键流程开始、结束、重要状态变化

- `Warn`
  可疑但还能继续的情况

- `Error`
  失败、依赖错误、不可恢复问题

- `Debug`
  只用于开发排查

不要在业务代码里继续扩散 `log.Printf(...)`。

如果请求已经经过 tracing 中间件，日志会自动带：

- `trace_id`
- `span_id`

## 参数校验怎么写

项目统一使用：

- Gin binding
- `internal/request`
- validator
- 中英文翻译

所以：

- JSON 参数用 `request.BindJSON`
- Query 参数用 `request.BindQuery`

不要在 handler 里继续手写一堆 `strconv.Atoi`、`if query == ""` 这类散逻辑。

## 缓存怎么写

当前缓存统一走 `internal/core/cache.Store`。

开发时注意：

- cache key 统一通过 `KeyBuilder` 构造
- 不要在业务代码里手拼 Redis key
- 缓存只是优化，不是业务真相源
- 列表缓存的 key 要包含所有查询参数，且分页参数需先 `NormalizePagination` 再生成 key

### 当前已实现的缓存策略

| 场景 | 缓存 TTL | 失效触发 |
|------|---------|---------|
| 文章详情 (`post:detail`) | `cacheTTL` 配置（默认 5 分钟） | 无自动失效，靠 TTL 过期 |
| 文章列表 (`post:list`) | `cacheTTL` 配置（默认 5 分钟） | Admin 端 Create/Update/Delete/Publish 后主动清除 |
| 用户资料 (`profile`) | `cacheTTL` 配置 | Profile Update 后主动清除 |

列表缓存使用 `DeleteByPattern`（基于 `SCAN`）批量清除，避免 `KEYS` 阻塞 Redis。清除失败只记 warn 日志，不影响主流程。

## 数据库怎么写

项目使用 GORM，但有几个明确约束：

- 不在本服务里做 migration
- model 只做表结构映射
- repository 只做数据访问
- 事务边界放在 service

如果一次业务动作涉及多个写操作，需要原子性，就在 service 层通过事务管理器控制。

不要：

- 在 handler 里开事务
- 在 repository 里偷偷开事务

## 提交流程建议

开发完成后，建议至少做这些检查：

1. 执行 `make fmt`
2. 执行 `make test`
3. 自己过一遍 OpenAPI 是否需要更新
4. 自己检查 `/metrics`、请求日志和 trace 头是否符合预期
5. 自己检查日志和错误响应是否符合约定
6. 自己检查路由是不是挂在正确层级

如果改动涉及：

- 新接口
- 参数变化
- 返回结构变化

就一定要同步更新 `api/openapi.yaml`。

## 自测清单

提交前建议至少自查：

- 路由是否注册在正确模块
- handler 是否只做 HTTP 输入输出
- service 是否没有依赖 Gin
- repository 是否没有跨模块乱查表
- 是否补了必要的 DTO
- 是否使用统一绑定层
- 是否使用统一日志层
- 是否更新了 OpenAPI
- `make test` 是否通过

## 常见错误

新人最容易犯的几个问题：

- 把业务逻辑写进 handler
- 直接在 handler 里查数据库
- 为了方便把 GORM model 当 API response 返回
- 在业务代码里手拼 Redis key
- 每个登录接口都额外查一次完整用户信息
- 改了接口但没更新 OpenAPI
- 在总路由文件里继续堆具体业务路由

## 建议的开发节奏

一个比较稳的节奏是：

1. 先看清楚接口应该归属哪个模块
2. 先定义 DTO 和接口契约
3. 再写 handler / service / repository
4. 本地跑通并自测
5. 更新 OpenAPI
6. 最后再提交

这样后续返工会少很多。

## 文档分工

当前项目的文档分工建议固定成这样：

- `README.md`
  首页、项目定位、快速开始

- `docs/architecture.md`
  架构设计、目录说明、核心原则

- `docs/development-guide.md`
  新人开发手册、实操流程、提交流程

- `docs/example-module.md`
  标准业务参考模块说明，给后续业务直接抄写

## 标准参考模块

如果你需要一个“可以直接照着写”的完整样例，请直接看：

- `internal/example/`
- `docs/example-module.md`

这个示例模块刻意展示了这些标准动作：

- `currentUser` middleware 中获取当前用户
- `handler` 只做请求绑定和响应映射
- `service` 调用别的模块 service contract
- 先做跨模块只读校验，再开事务写本模块
- 一个业务模块里同时维护主表和关联表

以后团队里有人问“新增业务到底按什么标准写”，默认就按这个示例来。
