# Proposal 0001: 双端认证统一

## 元数据

| 项 | 值 |
|---|---|
| 状态 | Draft |
| 类别 | 架构（CLAUDE.md §2.3）|
| 创建 | 2026-05-10 |
| 触发任务 | ARCHITECTURE_REPORT.md 严重问题 #2、第二轮评审 #6 |

## 摘要

把 API 端 `internal/auth/` 与 Admin 端 `internal/admin/auth/` 两套独立的认证实现，统一到单一 `internal/auth/core/`。终端用户与管理员属于不同 audience，token 不互通，但 JWT 签发 / 校验 / 刷新 / 黑名单 / 密码哈希 / 限流键全部共用一份代码。

按业内标准（OAuth 2.0 / RFC 7519）做 audience 隔离的 multi-tenant auth 模式。**不考虑历史兼容**——所有现有 token 失效，所有用户重新登录。

## 问题

| 现象 | 文件证据 |
|---|---|
| 两套 handler / service / repository | `internal/auth/handler/` vs `internal/admin/auth/handler/` |
| 两套 JWT TTL 配置且类型不一致 | `cfg.Auth.AccessTokenTTL()` 是 `time.Duration`；`cfg.Admin.JWT.Timeout` 是 `int` 单位"小时"（`internal/admin/bootstrap/app.go:91`）|
| Token 生命周期管理双份 | `internal/auth/tokenstore/` vs Admin 自己的 redis 操作 |
| 密码哈希双份 | `internal/auth/password/` vs admin/auth 独立实现 |
| middleware 双份 | `internal/core/middleware/jwt.go` vs `internal/admin/middleware/auth_middleware.go` |
| 一处安全修复要改两处 | 任何 JWT / refresh / blacklist 改动都要复制到对端 |

后果：

- 安全审计成本 ×2
- 修复不对称：可能 patch 一端漏另一端
- 配置 schema 已经事实漂移（int vs time.Duration）

## 目标

1. 单一 `internal/auth/core/` 提供 JWT 签发 / 校验 / refresh / blacklist / 密码哈希
2. Admin 与 User 是不同 audience，token 互不通用，但走同一套 manager
3. middleware 参数化 audience，一份实现两端复用
4. 配置 schema 统一为 `auth.{user, admin}.*`，删除 `cfg.Admin.JWT`
5. 共用 `auth.secret`，不再两套密钥

## 非目标

- 不引入 SSO / OAuth2 第三方提供者
- 不合并 `users` 与 `admins` 数据表（独立 proposal）
- 不引入 token rotation / device fingerprint 等高级特性
- 不更换 JWT 库（继续 `golang-jwt/jwt`）
- 不改 Casbin 策略

## 范围

### 改

| 路径 | 动作 |
|---|---|
| `internal/auth/` | 重构为 `core/` + `user/` 子目录 |
| `internal/admin/auth/` | 删除，逻辑迁入 `internal/auth/admin/` |
| `internal/admin/middleware/auth_middleware.go` | 删除，改用 core middleware |
| `internal/core/middleware/jwt.go` | 改造为 thin wrapper 调 core |
| `internal/core/config/config.go` | schema 调整 |
| `internal/admin/bootstrap/app.go` | manager 构造改 |
| `internal/bootstrap/provider.go`、`module_builders.go` | manager 构造改 |
| `internal/admin/routes/routes.go`、`admin_modules.go` | import 与注册改 |
| `configs/*.yaml` | 配置 schema 迁移 |

### 不改

- `internal/admin/admin_user/`（管理员用户数据访问）
- `internal/profile/`、`internal/post/` 等业务模块
- 数据库 schema、migrations
- Casbin 模型与策略

## 设计

### 包结构

```
internal/auth/
├── core/
│   ├── manager.go          # JWT 签发与校验，audience 参数化
│   ├── tokenstore.go       # refresh / blacklist (Redis)；key 按 audience 隔离
│   ├── claims.go           # Claims 定义（sub / aud / iss / exp / iat / jti）
│   ├── middleware.go       # auth.JWTMiddleware(manager, audience, store)
│   ├── password.go         # 密码哈希
│   └── errors.go           # 认证错误类型
├── user/                   # End-user audience
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── dto/
└── admin/                  # Admin audience
    ├── handler/
    ├── service/
    ├── repository/
    └── dto/
```

### 核心接口

```go
// core/manager.go
type Manager struct {
    issuer    string
    secret    []byte
    audiences map[string]AudienceConfig
}

type AudienceConfig struct {
    Audience        string
    AccessTokenTTL  time.Duration
    RefreshTokenTTL time.Duration
}

func NewManager(issuer string, secret string, audiences map[string]AudienceConfig) *Manager
func (m *Manager) Sign(audience string, subject string, extra map[string]any) (accessToken, refreshToken string, err error)
func (m *Manager) Verify(audience string, tokenString string) (*Claims, error)
func (m *Manager) Refresh(audience string, refreshToken string) (newAccess, newRefresh string, err error)
```

```go
// core/middleware.go
func JWTMiddleware(m *Manager, audience string, store *TokenStore) gin.HandlerFunc
```

注册示例：

```go
// API 端 (internal/bootstrap/router.go)
secured.Use(auth.JWTMiddleware(authManager, "gotribe.user", tokenStore))

// Admin 端 (internal/admin/routes/routes.go)
protected.Use(auth.JWTMiddleware(authManager, "gotribe.admin", tokenStore))
```

### 配置 schema

新增结构：

```yaml
auth:
  issuer: gotribe
  secret: <strong-secret-min-32-chars>
  user:
    audience: gotribe.user
    access_token_ttl_minutes: 120
    refresh_token_ttl_hours: 168
  admin:
    audience: gotribe.admin
    access_token_ttl_minutes: 60
    refresh_token_ttl_hours: 24
```

废弃字段（启动时 config validate 直接报错指引）：

| 旧字段 | 新字段 |
|---|---|
| `auth.access_token_ttl_minutes` | `auth.user.access_token_ttl_minutes` |
| `auth.refresh_token_ttl_hours` | `auth.user.refresh_token_ttl_hours` |
| `admin.jwt.realm` | `auth.admin.audience` |
| `admin.jwt.key` | 删除（共用 `auth.secret`）|
| `admin.jwt.timeout` | `auth.admin.access_token_ttl_minutes`（注意单位从小时改成分钟）|
| `admin.jwt.max_refresh` | `auth.admin.refresh_token_ttl_hours` |
| `admin.jwt.token_lookup` | 删除（统一 Bearer header）|

### TokenStore Redis Key 规范

```
{appName}:auth:{audience}:refresh:{userID}:{jti}
{appName}:auth:{audience}:blacklist:{jti}
```

audience 作为 key 前缀的一部分，确保两边互不干扰。

### 迁移步骤

1. 新建 `internal/auth/core/`，整合 `jwt/`、`tokenstore/`、`password/`
2. 把 `internal/auth/handler|service|repository|dto/` 移到 `internal/auth/user/`
3. 把 `internal/admin/auth/handler|service/` 移到 `internal/auth/admin/`，改造 service 让其依赖 `core.Manager`
4. `internal/core/middleware/jwt.go` 改造为通用 middleware
5. 删除 `internal/admin/middleware/auth_middleware.go`
6. 改 `internal/core/config/config.go` schema + viper defaults + validate
7. 改两端 bootstrap：构造一个 `core.Manager` + 两个 audience config
8. 改两端 routes：注册时传 audience
9. 迁移 `configs/*.yaml` 文件
10. 删除 `internal/admin/auth/` 与 `internal/auth/{handler,service,repository,dto}/`（已迁走）
11. 跑全量测试 + 手动 smoke

## 影响

### 用户层面

- **所有 admin 用户必须重新登录**（旧 token 无 `aud` claim）
- **所有 api 用户也必须重新登录**（同上）
- 配置 schema 变更：现有部署 `config.yaml` 必须更新

### 代码层面

预估改动：

| 类别 | 文件数 |
|---|---|
| `internal/auth/` 重构 | ~15 |
| `internal/admin/auth/` 删除 + 迁移 | ~6 |
| middleware 统一 | ~3 + 测试 |
| config schema | 1 + 测试 |
| bootstrap (api + admin) | 2 |
| routes (api + admin) | 2 |
| 配置 yaml | ~3 |
| 总计 | ~50 |

### 配置层面

- `auth.secret` 变成共用密钥
- Admin JWT 单位从「小时」改成「分钟」（与 User 端一致）
- 新增 `validate()` 检查老字段是否仍存在，存在则报错指引迁移

## 验证

### 自动化

| 测试 | 覆盖 |
|---|---|
| `auth/core/manager_test.go` | 签发 / 校验 / audience mismatch / 过期 / 错签名 |
| `auth/core/tokenstore_test.go` | refresh / blacklist / audience key 隔离 |
| `auth/core/middleware_test.go` | 通过 / 拒绝 / audience 不匹配 |
| `auth/user/service_test.go` | login / refresh / change_password |
| `auth/admin/service_test.go` | 同上 |
| `tests/integration/auth_test.go` | API 端登录 → 受保护接口 200 |
| `tests/admin/integration/auth_test.go` | Admin 端登录 → 受保护接口 200 |
| 跨端 token 拒绝 | API token 投到 Admin → 401（aud mismatch）|

### 手动

- `go build ./...`
- `go vet ./...`
- `go test ./...`
- 启动两端，分别 login 并访问保护接口
- 把一端 token 投到另一端，确认 401

## 回滚

- `git revert` 实现 commit
- 配置文件还原（保留旧字段）
- 重启两端

回滚后所有用户再次重新登录（带 aud 的新 token 失效）。

## 备选方案

### A: 完全合并 admin/api 为单一用户体系

合并 `users` 与 `admins` 表，用 role 区分。

- ✅ 彻底零冗余
- ❌ 数据模型大改，业务逻辑耦合
- ❌ 超出认证统一范围

**决定：拒绝。**本 proposal 只统一认证 core，不动用户数据模型。

### B: 保留两个 Manager 实例，仅抽接口

- ✅ 迁移成本最低
- ❌ 仍有两份运行时状态，配置 schema 仍漂移
- ❌ 不符合"业内标准"

**决定：拒绝。**

### C: 引入 OAuth2 / OIDC（ory/fosite 或 hydra）

- ✅ 完全标准
- ❌ 引入新依赖，运维代价大
- ❌ 远超本任务范围

**决定：拒绝。**延后到独立 proposal。

## 开放问题

1. **密码哈希成本**：admin 是否设置更高的 bcrypt cost？
   - 提议：admin cost = 12，user cost = 10
2. **Refresh token rotation**：是否在本次顺带启用？
   - 提议：本次只做存储隔离，rotation 延后
3. **Token blacklist 是否引入**：
   - 提议：本次顺带做，否则 logout 无效

待 proposal 通过时一并明确。

## 实现完成定义

- [ ] 所有"改"清单内文件改完
- [ ] 所有"不改"清单内文件未被触碰
- [ ] 全部测试通过
- [ ] swagger 重新生成
- [ ] `configs/config.example.yaml` 更新为新 schema
- [ ] `docs/architecture.md` 认证部分更新
- [ ] `lessons.md` 记录此次架构变更原因
