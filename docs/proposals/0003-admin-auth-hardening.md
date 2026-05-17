# Proposal 0003: 后台登录加固（失败锁定 + TOTP 2FA，opt-in）

## 元数据

| 项 | 值 |
|---|---|
| 状态 | In Progress |
| 类别 | 架构（CLAUDE.md §2.3：改认证模型 / 改接口协议 / 改数据库 / 引新依赖）|
| 创建 | 2026-05-17 |
| 触发 | 上线前安全加固讨论 |

## 摘要

后台登录加两道防线：① 密码失败按账户+IP 双维度计数并临时锁定；② 可选 TOTP 2FA（与 Google / Microsoft Authenticator / 1Password 等所有 RFC 6238 客户端兼容）。**TOTP 为 opt-in**：未绑用户登录后弹窗提示但可关闭继续工作；一旦绑定则后续登录强制二次校验。Redis 故障时 fail-open（允许登录 + 告警日志）。

## 问题

| 现象 | 证据 |
|---|---|
| 知道用户名后可无限次试密码 | `internal/auth/admin/service/auth.go` 的 `Login` 不计失败次数 |
| 全局限流不防针对性爆破 | `internal/admin/routes/routes.go:79` 全 engine 共享一个 token bucket |
| 密码泄露即沦陷，无第二因子 | 同上文件，登录仅验密码即签 JWT |

## 目标

1. 登录失败按「账户 + IP」双维度计数，达阈值后短时锁定
2. 提供可选 TOTP 绑定能力，与所有 RFC 6238 客户端兼容
3. 一次性备份码兜底（防 TOTP 设备丢失）
4. 超管可重置他人 TOTP
5. 现存 admin 用户无感升级（不强制绑定）

## 非目标

- 不做短信 / 邮件 OTP（成本+被劫持风险）
- 不做 WebAuthn / Passkey（成本太高，下阶段考虑）
- 不动前台 user 登录（仅 admin）
- 不在用户已绑 TOTP 后提供「跳过 TOTP」开关（一旦绑定即强制；不绑则可继续工作）
- 不做「永远不再提示」记忆——未绑用户每次登录都提示，直到绑定

## 设计

### 登录流程

| 用户状态 | 流程 |
|---|---|
| 触发锁定 | `POST /api/base/login` → 返 `{ error: "locked", locked_until }` |
| 未绑 TOTP | `POST /api/base/login` → 返 `{ access_token, expires, mfa_reminder: true }`（前端展示弹窗）|
| 已绑 TOTP | `POST /api/base/login` → 返 `{ stage: "totp_required", step_token }`（5 分钟有效）<br/>→ `POST /api/base/totp/verify` 提交 6 位 code 或备份码 → `{ access_token, expires }` |

### TOTP 管理（登录态下）

| 端点 | 鉴权 | 用途 |
|---|---|---|
| `GET  /api/base/totp/status` | access_token | 返回当前账户是否已绑 + 剩余备份码数量 |
| `POST /api/base/totp/bind` | access_token | 生成新 secret + `otpauth://` URL + 10 个明文备份码（只此一次返回）|
| `POST /api/base/totp/confirm` | access_token | 提交 6 位 code 完成绑定，激活记录 |
| `DELETE /api/base/totp` | access_token + 当前 6 位 code | 自助解绑 |
| `POST /api/base/admin/:id/totp/reset` | 超管 | 强制重置他人 TOTP（用于设备丢失救援）|

### 失败锁定

| 项 | 设计 |
|---|---|
| 存储 | Redis（项目已有 `internal/core/cache`） |
| Key | `auth:fail:user:<name>` 与 `auth:fail:ip:<ip>` |
| 计数 | 密码错或 TOTP 错均 +1；登录成功（含 TOTP 验证通过）清零 |
| 阈值（默认） | 账户：5 次 → 锁 15 分钟<br/>IP：20 次 → 锁 1 小时（防单 IP 扫账户）|
| 错误响应 | 包含 `remaining_attempts`、`locked_until` |
| Redis 故障 | fail-open：允许登录 + 写 warn 日志 + 计数 metric 上报 |
| 配置 | `auth.lockout.*` 在 `internal/core/config` 中可调；可整体关闭 |

### TOTP 弹窗（前端）

| 项 | 设计 |
|---|---|
| 触发 | 登录响应含 `mfa_reminder: true` 时显示 |
| 操作 | 「去绑定」跳安全设置页 / 「下次再说」关闭 |
| 频率 | 每次登录都展示（未绑期间）|
| 视觉 | shadcn `<Dialog>`，warning 色调，非阻塞 |

### 数据库

新增表 `admin_totp`：

| 字段 | 类型 | 说明 |
|---|---|---|
| admin_id | bigint PK FK | 关联 admin.id |
| secret | varbinary(96) | TOTP 共享密钥（AES-GCM 加密存储，密钥取自配置）|
| enabled | tinyint(1) | 是否完成绑定（bind 创建时 0，confirm 后置 1）|
| recovery_codes | json | 10 个备份码的 bcrypt hash 数组，每项 `{hash, used_at}` |
| created_at | datetime | 标准 |
| updated_at | datetime | 标准 |

**不**在 `admin` 表加 `mfa_enabled` 字段——通过 `admin_totp` 行存在 + `enabled=1` 判断，避免双写一致性问题。

up migration 建表；down migration 直接 `DROP TABLE admin_totp`。

### step_token 设计

短期签发的 JWT（5 分钟），claims 含 `admin_id`、`purpose: "totp_verify"`、`jti`。
验证后 `jti` 写入 Redis 黑名单到过期为止，防重放。Redis 故障时 fail-open（允许验证一次，写日志）。

### 依赖

| 依赖 | 用途 |
|---|---|
| `github.com/pquerna/otp` | TOTP 生成 / 验证（RFC 6238 标准实现）|
| 前端 `qrcode` 类库（~3KB） | 前端用 `otpauth://` URL 渲染 QR 图 |

### 配置

```yaml
auth:
  lockout:
    enabled: true
    account_max_fails: 5
    account_lock_minutes: 15
    ip_max_fails: 20
    ip_lock_minutes: 60
  totp:
    issuer: "GoTribe 管理后台"
    secret_encryption_key: "${TOTP_SECRET_KEY}"  # 32 字节随机，缺失则启动失败
    recovery_codes_count: 10
    step_token_ttl_seconds: 300
    period: 30        # 时间步长
    digits: 6         # 验证码位数
    skew: 1           # 允许 ±1 步时钟漂移（30s）
```

### 前端改动

| 文件 / 模块 | 改动 |
|---|---|
| `features/auth/sign-in/components/user-auth-form.tsx` | 判定登录响应 stage，已绑→跳 TOTP 步骤，未绑→正常进入 + 触发弹窗 |
| `features/auth/otp/` | 替换 stub `onSubmit`，接入 `/totp/verify`，支持「使用备份码」切换 |
| 新增 `features/auth/mfa-reminder-dialog.tsx` | 弹窗组件 |
| 新增 `features/settings/security/` | 安全设置页：绑定 / 解绑 TOTP，显示备份码状态 |
| `service/index.ts` `AUTH_API_PATHS` | 增加 `/api/base/totp/...`、`/api/base/admin/:id/totp/reset` |
| 新增 i18n：锁定提示 / 剩余次数 / TOTP 错误 / 备份码提示 / 弹窗文案 | zh.json + en.json |

## 影响

| 范围 | 影响 |
|---|---|
| 现存 admin 用户 | 无感升级，登录方式不变；首次登录后看到一次弹窗，可关 |
| 数据库 | 新增 1 张表 `admin_totp`，up / down 齐备 |
| API 协议 | `/api/base/login` 返回结构变更（新增 `stage` / `mfa_reminder` 字段，原有字段保留）|
| 部署 | 必须新增环境变量 `TOTP_SECRET_KEY`（32 字节随机），缺失则启动失败 |
| 紧急恢复 | 备份码 + 超管重置 + 配置项关锁定，三条恢复路径 |

## 验证

- **单测**：
  - `auth/admin/service`：失败计数与锁定边界、清零逻辑
  - TOTP service：生成/验证/skew、备份码消费、step_token 签发与黑名单
  - secret AES-GCM 加解密往返
- **集成测**：完整两步登录、未绑直登 + 弹窗 flag、锁定恢复、备份码用尽后回退
- **手动验证**（浏览器全流程）：
  - 用 Google Authenticator 扫码绑定
  - 用 Microsoft Authenticator 扫同一 QR 验证（确认互通）
  - 用 1Password 内置 TOTP 验证（确认互通）
  - 故意输错 5 次触发锁定
  - 用备份码登录 + 备份码用尽后的提示
  - 解绑 + 重绑

## 回滚

| 场景 | 操作 |
|---|---|
| 仅想关锁定 | 改 `auth.lockout.enabled=false`，运行时生效，无需重发 |
| 仅想停 TOTP 校验 | 不能热关。TOTP 流程已内嵌登录路径；如需回滚需 deploy 旧 binary |
| 完整回滚 | 切回旧 binary + 跑 down migration（仅删 `admin_totp` 表，admin 表无变更）|

## 风险与已决

| 项 | 决策 |
|---|---|
| 备份码明文一次性返回 | 用户截图泄露——文案强提示「仅显示一次，请立即保存」+ 提供「重新生成」按钮 |
| Redis 不可用 | **fail-open**：允许登录 + 写 warn 日志 + 计数 metric |
| 时钟漂移 | 允许 ±1 步（30s）容差 |
| step_token 重放 | 单次使用，jti 写 Redis 黑名单到过期 |
| Casbin 授权 | 超管重置接口需在 RBAC seed 加新 API path 条目 |
| 强制 vs 可选 | **可选**，未绑用户登录后弹窗提示，可关闭继续工作 |

## 工作量预估

| 阶段 | 预估 |
|---|---|
| 后端：migration、lockout 中间件、TOTP service、handler、单测 | 1.5 天 |
| 前端：弹窗、TOTP 输入接 API、绑定页、安全设置页、i18n | 1 天 |
| 集成测 + 手动验证 + 文档 | 0.5 天 |
| **合计** | **3 天** |

## 后续（不在本 proposal）

- 登录历史 / 设备审计页
- WebAuthn / Passkey 二次升级
- 后台超管的强制 IP 白名单
