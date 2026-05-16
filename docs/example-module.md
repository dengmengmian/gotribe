# Example Module Guide

## 目标

这份文档把 `internal/example` 作为标准参考模块，给团队成员一个可以直接照着写的新业务模块模板。

重点回答这些问题：

- 一个完整业务模块应该长什么样
- 新增模块时推荐的开发顺序是什么
- 当前用户信息应该怎么拿
- 一个 service 调另一个 service 应该怎么接线
- 事务应该放在哪里
- 哪些写法是错误的，为什么错

如果你只想记一句话，请记这个：

- `handler` 只处理 HTTP
- `service` 只处理业务
- `repository` 只处理数据访问
- 模块依赖在 `internal/bootstrap/provider.go` / `internal/bootstrap/module_builders.go` 装配
- 路由在模块自己的 `handler/routes.go` 注册

## 先看哪几个文件

推荐按这个顺序阅读 `example` 模块：

1. `internal/example/handler/routes.go`
2. `internal/example/handler/example.go`
3. `internal/example/dto/example.go`
4. `internal/example/view/example.go`
5. `internal/example/service/example.go`
6. `internal/example/repository/example.go`
7. `internal/example/model/example.go`
8. `internal/bootstrap/provider.go`
9. `internal/bootstrap/module_builders.go`
10. `internal/bootstrap/router.go`

理解顺序是：

- 先看暴露了哪些接口
- 再看 handler 怎么收参数和回响应
- 再看 service 怎么组织业务
- 再看 repository 怎么访问数据库
- 最后看整个模块是怎么接到应用里的

## 标准目录结构

新增业务模块时，推荐至少有这些目录：

```text
internal/example/
├── dto/
├── handler/
├── model/
├── repository/
├── service/
└── view/
```

每层职责如下：

- `dto`
  只放 HTTP 请求和响应结构

- `handler`
  只做参数绑定、调用 service、输出响应

- `model`
  只做数据库表结构映射

- `repository`
  只做数据库读写

- `service`
  只做业务逻辑、事务边界、跨模块协作

- `view`
  放模块内部返回结构

`view` 很重要。它的作用是把“模块内部结果”和“HTTP 响应 DTO”隔开，避免 service 直接依赖对外响应结构。

## 标准开发流程

新增一个完整业务模块，建议按下面顺序写：

1. 先确认它是不是一个独立业务域
2. 设计接口和 OpenAPI
3. 定义 `dto`
4. 定义 `view`
5. 定义 `model`
6. 写 `repository`
7. 写 `service`
8. 写 `handler`
9. 在模块里增加 `RegisterRoutes(...)`
10. 在 `internal/bootstrap/provider.go` / `internal/bootstrap/module_builders.go` 装配依赖
11. 在 `internal/bootstrap/router.go` 挂接路由
12. 补测试并执行 `make test`

推荐原因是：

- 先把边界想清楚，再写实现
- repository 和 service 的接口会更稳定
- 不容易把 DTO、数据库 model、业务 view 混在一起

## `example` 模块做了什么

`internal/example` 这个模块完整展示了下面几类常见需求：

- 增删改查
- 当前用户信息获取
- service 调别的模块 service
- 一个业务动作里同时写主表和关联表
- 列表查询 + 分页
- handler 层 DTO 映射

它提供的接口包括：

- `POST /api/v1/examples`
- `GET /api/v1/examples`
- `GET /api/v1/examples/{exampleID}`
- `PATCH /api/v1/examples/{exampleID}`
- `DELETE /api/v1/examples/{exampleID}`

这些接口都挂在 `currentUser` 路由层，因为它依赖完整当前用户对象。

## handler 应该怎么写

看 `internal/example/handler/example.go`，handler 只做 4 件事：

1. 从上下文拿当前用户
2. 绑定请求参数
3. 调用 service
4. 把 service 返回的 `view` 转成响应 DTO

也就是说，handler 不负责：

- 事务
- GORM
- Redis
- 拼接 SQL
- 直接 new 其他模块

标准形状大概就是：

```go
func (h *Handler) Create(c *gin.Context) {
    actor, err := actorFromContext(c)
    if err != nil {
        response.Error(c, err)
        return
    }

    var req dto.CreateRequest
    if err := request.BindJSON(c, &req); err != nil {
        response.Error(c, err)
        return
    }

    data, err := h.service.Create(c.Request.Context(), middleware.GetProjectID(c), actor, req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Created(c, toResponse(*data))
}
```

这个写法的价值是：

- HTTP 细节留在 handler
- 业务细节留在 service
- 错误出口统一

## 当前用户应该怎么拿

如果接口依赖完整当前用户对象，先把路由挂到 `currentUser` 组，然后在 handler 里通过：

```go
currentUser, ok := middleware.GetCurrentUser(c)
```

如果只需要 `user_id`，优先挂在 `secured`，再通过：

```go
userID, ok := middleware.GetUserID(c)
```

不要为了偷懒，把所有登录后接口都挂到 `currentUser`。因为 `CurrentUser` 中间件会额外查一次用户信息。

## service 应该怎么写

`internal/example/service/example.go` 里展示的是标准 service 形状：

- 收 `context.Context`
- 收业务需要的输入参数
- 做参数整理和业务校验
- 调 repository
- 需要时调别的模块能力
- 需要时定义事务边界
- 返回模块内部 `view`

比如 `Create()` 做了这些事：

1. 清理和校验输入
2. 校验引用的 `post` 是否存在
3. 组装主表 model
4. 在事务里写主表和关联表
5. 返回内部 `view`

`Update()` 和 `Delete()` 也都把事务边界放在 service 里控制。

这就是推荐的原则：

- 一个业务动作涉及多个写操作时，事务边界一定由 service 来定义

## repository 应该怎么写

`internal/example/repository/example.go` 展示的是标准 repository：

- 只关心表结构和查询条件
- 不关心 HTTP 参数长什么样
- 不返回 HTTP 语义
- 不跨模块自己乱查别人表

比如列表查询接收的是自己的 `ListFilter`，而不是 `dto.ListQuery`。

这个边界非常重要：

- `dto.ListQuery` 是传输层
- `repository.ListFilter` 是数据访问层

以后接口参数名改了、校验规则改了，repository 不应该被一起拖着改。

## 模块内部 `view` 为什么要单独存在

`internal/example/view/example.go` 和 `internal/profile/view/me.go` 的存在，都是为了解决一个问题：

- service 返回值不应该直接等于 HTTP 响应 DTO

推荐做法是：

- service 返回模块内部 `view`
- handler 再把 `view` 转成响应 DTO

这样做的好处：

- middleware、其他 service、后台任务都可以复用这个 `view`
- 接口响应结构调整时，不会把 service 一起绑死
- 模块边界更干净

## 一个 service 调另一个 service，标准写法是什么

这是团队最容易写歪的地方。

以 `example -> post` 为例，正确理解是：

- `example service` 只声明“我需要帖子读取能力”
- bootstrap 装配层负责把真正的 `post service` 注入进去

### 正确写法

先在调用方 service 里定义自己需要的能力接口：

```go
type PostSummaryReader interface {
    GetSummaries(ctx context.Context, projectID string, postIDs []string) (map[string]postview.Summary, error)
}
```

然后 service 只依赖这个接口：

```go
type Service struct {
    repo  *Repository
    tx    *database.TransactionManager
    posts PostSummaryReader
}
```

构造函数也只接这个能力：

```go
func NewService(repo *Repository, tx *database.TransactionManager, posts PostSummaryReader) *Service
```

最后在 `internal/bootstrap/module_builders.go` 装配：

```go
post := buildPostModule(infra)
example := buildExampleModule(infra, post.Service)
```

这里的关键是：

- `example` 只知道自己要一个 `PostSummaryReader`
- 装配层决定由 `post.Service` 来提供这个能力

### 为什么一定要这么做

因为这样才能做到：

- 依赖关系集中可见
- service 本身更容易测试
- 后续替换实现时改动更小
- 模块之间不会偷偷耦合初始化细节

## 事务标准写法

这个项目里，事务统一由 `database.TransactionManager` 管理，事务边界统一放在 service。

推荐写法：

```go
if err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
    if err := s.repo.Create(txCtx, entity); err != nil {
        return errs.Internal("create example", err)
    }
    if err := s.repo.CreatePosts(txCtx, postRows); err != nil {
        return errs.Internal("create example posts", err)
    }
    return nil
}); err != nil {
    return nil, err
}
```

为什么是 service 开事务，而不是 repository：

- repository 只知道单次读写
- service 才知道一个业务动作到底包含几步
- 只有 service 知道哪些步骤必须一起成功或一起失败

## provider 和 router 应该怎么接

新增模块后，统一走这两个文件：

- `internal/bootstrap/provider.go`
- `internal/bootstrap/module_builders.go`
- `internal/bootstrap/router.go`

### 在 `provider.go` 里做什么

只做“创建基础设施 + 汇总模块”：

1. 初始化 `Infra`
2. 调用模块 builder
3. 返回 `Providers{Infra, Modules}`

### 在 `module_builders.go` 里做什么

只做“每个模块的 repo / service / handler 装配”：

1. new repository
2. new service
3. new handler
4. 把模块收进 `Modules`

### 在 `router.go` 里做什么

只做“挂接路由”：

1. 选择路由层级
2. 调模块自己的 `RegisterRoutes(...)`

不要把模块内的具体业务路由又展开写回总路由文件。

## 错误写法

下面这些写法需要明确禁止。

### 错误 1：在 service 里自己 new 别的模块

错误示例：

```go
func NewService(repo *Repository, tx *database.TransactionManager) *Service {
    postRepo := postrepo.NewRepository(tx)
    postSvc := postservice.NewService(postRepo, nil, nil)
    return &Service{
        repo: repo,
        tx:   tx,
        post: postSvc,
    }
}
```

为什么错：

- 依赖关系藏进业务代码里了
- 初始化逻辑散掉了
- 测试时很难替换成 mock
- 以后 `post` 的构造函数一改，这里也得跟着改

正确做法：

- 在 service 里定义能力接口
- 在 bootstrap 装配层注入实现

### 错误 2：repository 直接接收 HTTP DTO

错误示例：

```go
func (r *Repository) List(ctx context.Context, query dto.ListQuery) ([]model.Example, error)
```

为什么错：

- repository 被 HTTP 传输层污染
- `form`、`json`、`binding` 这些语义不该进入数据访问层
- 接口契约变化会直接冲击 repository

正确做法：

- handler 绑定 `dto`
- service 映射到 repository 自己的 `ListFilter`
- repository 只接 `ListFilter`

### 错误 3：service 直接返回 HTTP 响应 DTO

错误示例：

```go
func (s *Service) Detail(...) (*dto.ExampleResponse, error)
```

为什么错：

- service 和对外响应结构绑死
- middleware 和其他 service 很难复用
- 后续改接口字段时容易牵一发动全身

正确做法：

- service 返回 `view`
- handler 负责 `view -> dto`

### 错误 4：在 handler 里写业务逻辑或开事务

错误示例：

```go
func (h *Handler) Create(c *gin.Context) {
    tx := h.db.Begin()
    ...
}
```

为什么错：

- HTTP 层开始承担业务职责
- 事务边界和业务规则混进传输层
- 逻辑不可复用，也难测试

正确做法：

- handler 只收参和回参
- service 定义事务边界

### 错误 5：service 依赖 Gin

错误示例：

```go
func (s *Service) Create(c *gin.Context, req dto.CreateRequest) error
```

为什么错：

- service 被 HTTP 框架绑死
- 无法在测试、任务、其他入口中复用

正确做法：

- service 只接 `context.Context`
- handler 把需要的信息拆出来传进去

### 错误 6：需要完整当前用户时，在 handler 里自己查库

错误示例：

```go
userID, _ := middleware.GetUserID(c)
user, _ := h.userRepo.GetByID(c, userID)
```

为什么错：

- 当前用户加载逻辑分散
- 破坏统一中间件约定
- 以后字段变化容易到处改

正确做法：

- 需要完整当前用户时挂 `CurrentUser`
- 在 handler 里调用 `middleware.GetCurrentUser(c)`

### 错误 7：跨模块直接改别人表

错误示例：

- `example repository` 直接写 `post` 表
- `order repository` 直接改 `user` 表里的别的业务字段

为什么错：

- 领域边界被打穿
- 后面很难知道谁在维护这张表的业务规则

正确做法：

- 通过对方模块暴露的 service 能力协作
- 自己模块只写自己负责的数据

## 团队开发检查清单

写完一个新模块后，至少自己对照一遍：

- 是否有清晰的 `dto / view / model / repository / service / handler`
- handler 是否没有直接查数据库
- service 是否没有依赖 Gin
- repository 是否没有依赖 HTTP DTO
- 跨模块依赖是否通过 bootstrap 装配层注入
- 事务是否放在 service
- 路由是否在模块自己的 `routes.go`
- 是否更新了 `api/openapi.yaml`
- 是否补了必要测试
- `make test` 是否通过

## 最后建议

以后团队新增业务模块，优先参考：

- `internal/example`
- `internal/profile`
- `internal/post`

其中：

- `example` 用来看完整标准流程
- `profile` 用来看当前用户和事务
- `post` 用来看列表读取和跨模块能力输出

如果一个新模块写出来之后，已经开始出现这些现象：

- handler 越来越胖
- repository 开始吃 DTO
- service 直接返回 HTTP 响应
- service 自己 new 别的模块

那就说明边界已经开始歪了，应该及时收回来。
