# Proposal 001 — `@gotribe/mcp`:gotribecms 的统一 MCP 接入层

| 字段 | 内容 |
|---|---|
| 提案 ID | 001 |
| 状态 | **待用户确认** |
| 创建日期 | 2026-05-17 |
| 作者 | Claude(Opus 4.7)+ 麻凡 |
| 关联仓库 | 新建 `gotribe-mcp`(独立仓库,与本仓库平层) |
| npm 包名 | `@gotribe/mcp` |

---

## 1. 问题(Problem)

`gotribecms` 已经实现了完整的发文、读文、分类、上传等后台 API,但目前**没有一种跨 AI Agent 通用的接入方式**让外部 agent(Claude Code / Codex / OpenCode 等)能通过自然语言驱动这些能力。

现状下三条可走的路都不够好:

| 方案 | 缺陷 |
|---|---|
| 每个 agent 单独写 skill / AGENTS.md | 三家格式各不相同,要写并维护 3+ 份;skill 只是"说明书",让 agent 自己拼 curl 调 API,参数容易错、鉴权易丢 |
| 让 agent 直接读 OpenAPI 文档调 API | token 会写进 agent 上下文,泄露风险高;agent 拼请求容易拼错;错误码不会翻译,体验差 |
| 让用户每次手工 curl | 完全失去 agent 自动化的意义 |

需要一个**协议层的统一接入方案**,满足:跨 agent 通用、自然语言触发、token 不进 agent 上下文、字段结构化校验。

## 2. 目标(Goals)

| # | 目标 | 衡量标准 |
|---|---|---|
| G1 | **一份代码,三家 agent 通用** | Claude Code / Codex / OpenCode 用同一份配置即可接入 |
| G2 | **自然语言触发** | 用户说"在 gotribe 发一篇标题为 X 的草稿",agent 自动调正确 tool |
| G3 | **零鉴权泄露** | access_token 永远不进 agent 的对话上下文,仅 MCP 进程持有 |
| G4 | **覆盖核心发文/读文场景** | P0 暴露 8 个 tool(见 §4.3),覆盖列表、详情、创建、更新、发布、删除、分类、上传 |
| G5 | **极低接入成本** | 用户三步完成接入:① 注入 token 到 env;② 复制 mcpServers 配置;③ 重启 agent |
| G6 | **错误透明,不掩盖** | 后端错误码翻译成 agent 可读消息,不静默吞错,不假装成功 |

## 3. 范围(Scope)

### 3.1 在范围

- 新建独立仓库 `gotribe-mcp`(平层目录:`/Users/mengmian/Develop/app/gotribe/gotribe-mcp`)
- **Node.js 18+ / TypeScript** 实现
- 基于 `@modelcontextprotocol/sdk` 官方 SDK,**stdio transport**
- 暴露 P0 工具(8 个,见 §4.3)
- 内部封装:
  - HTTP 客户端(`fetch`)
  - Token 自动 refresh
  - 后端错误码 → MCP 错误消息翻译
  - Markdown → HTML 自动渲染(免去 agent 同时传两份)
  - 上传后 `domain` + `key` 拼接为完整 URL
- npm 包名 `@gotribe/mcp`,`bin` 字段使其支持 `npx -y @gotribe/mcp`
- 完整 README:三家 agent 接入示例 + 首次 token 获取流程
- CI:GitHub Actions(build + test + lint),tag 触发自动 publish

### 3.2 不在范围(第一版不做)

| 项 | 推迟到 |
|---|---|
| `list_tags` / 标签管理 | P1 |
| 资源列表 / 详情 / 删除(只做 `upload_resource`) | P1 |
| 评论 / 用户 / 项目管理 | P2+ |
| HTTP / SSE transport | 暂不需要,stdio 满足三家客户端 |
| 多账号切换 | 暂不需要,env 配一个账号即可 |
| Web UI / 配置面板 | 不做 |
| 复用 gotribecms Go 代码 | 跨语言不可行,也不应耦合 |

## 4. 设计(Design)

### 4.1 仓库结构

```
gotribe-mcp/
├── src/
│   ├── index.ts                # MCP server 入口(stdio)
│   ├── server.ts               # tool 注册、请求分发
│   ├── client/
│   │   ├── http.ts             # fetch 封装 + 错误码翻译
│   │   ├── auth.ts             # token 缓存、refresh 流程
│   │   └── types.ts            # 共享类型(Post/Category/Resource 等)
│   ├── tools/
│   │   ├── list-posts.ts
│   │   ├── get-post.ts
│   │   ├── create-post.ts
│   │   ├── update-post.ts
│   │   ├── publish-post.ts
│   │   ├── delete-posts.ts
│   │   ├── list-categories.ts
│   │   └── upload-resource.ts
│   ├── config.ts               # 环境变量加载 + zod 校验
│   ├── markdown.ts             # markdown-it 渲染
│   └── errors.ts               # 错误码翻译表
├── tests/                      # vitest
├── package.json
├── tsconfig.json
├── README.md
├── LICENSE                     # MIT
└── .github/workflows/
    ├── ci.yml
    └── publish.yml
```

### 4.2 环境变量

| 变量 | 必填 | 说明 |
|---|---|---|
| `GOTRIBE_API_BASE_URL` | ✓ | 后端地址,如 `https://cms.gotribe.cn` |
| `GOTRIBE_API_PREFIX` | | 路由前缀,默认 `/api`(对应 `admin.url_path_prefix`) |
| `GOTRIBE_ACCESS_TOKEN` | ✓ | 已登录拿到的 access_token。**M2 实施时发现:gotribecms refresh 流程是"用过期 access_token 调 refresh 接口"换新 token,不需要独立 refresh_token,因此移除原计划的 `GOTRIBE_REFRESH_TOKEN` env** |
| `GOTRIBE_DEFAULT_PROJECT_ID` | ✓ | 创建文章必带 |
| `GOTRIBE_DEFAULT_USER_ID` | ✓ | 创建文章必带 |
| `GOTRIBE_DEFAULT_AUTHOR` | ✓ | 创建文章必带 |
| `GOTRIBE_CDN_DOMAIN_OVERRIDE` | | 上传返回 `domain` 为空时的兜底 |
| `GOTRIBE_LOG_LEVEL` | | `debug` / `info` / `warn` / `error`,默认 `info`,写 stderr |

### 4.3 P0 Tools 清单

| Tool | 后端端点 | 关键封装 |
|---|---|---|
| `list_posts` | `GET /api/post` | 默认走后台接口(可见草稿);支持 `keyword/status/category_id/page/per_page` |
| `get_post` | `GET /api/post/:id` | 返回完整文章(含 markdown 原文 + html) |
| `create_post` | `POST /api/post` | **MCP 内部 markdown→html 渲染**;`project_id/user_id/author` 自动从 env 注入 |
| `update_post` | `PATCH /api/post/:id` | 同 create;但只传入差异字段(MCP 先 get 拿原文,合并后整体 PATCH) |
| `publish_post` | `PUT /api/post/:id` | 草稿 → 发布(触发百度推送) |
| `delete_posts` | `DELETE /api/post` | 批量,`post_ids: [int64]` |
| `list_categories` | `GET /api/category/tree` | 返回树形,方便 agent 给用户选 category_id |
| `upload_resource` | `POST /api/resource/upload` | 接收 base64 + 文件名;**MCP 拼接 `{domain}/{key}` 返回完整 URL** |

### 4.4 `create_post` Tool Schema 示例

```typescript
{
  name: "create_post",
  description: "在 gotribecms 创建一篇新文章。content_markdown 由 MCP 自动渲染为 HTML,无需手动传 html_content。category_id 必填,如不知道请先调用 list_categories。",
  inputSchema: {
    type: "object",
    properties: {
      title:            { type: "string",  minLength: 2, maxLength: 60 },
      description:      { type: "string",  minLength: 2, maxLength: 300 },
      content_markdown: { type: "string",  description: "Markdown 原文" },
      category_id:      { type: "integer", description: "先调 list_categories 查" },
      type:             { type: "integer", enum: [1, 2, 3], default: 1, description: "1=普通,2=页面,3=短文" },
      status:           { type: "integer", enum: [1, 2],    default: 1, description: "1=草稿,2=发布。建议先草稿后用 publish_post 发布" },
      slug:             { type: "string",  description: "URL slug,不填自动生成" },
      tag:              { type: "string",  description: "逗号分隔的 tag slug,如 'golang,web'" },
      icon:             { type: "string",  description: "封面图 URL(先用 upload_resource 上传)" },
      images:           { type: "array", items: { type: "string" } }
    },
    required: ["title", "description", "content_markdown", "category_id"]
  }
}
```

**关键设计决策**:

| 决策 | 理由 |
|---|---|
| 不暴露 `html_content` 入参 | 让 agent 同时维护 markdown + html 心智负担过大,**MCP 用 `markdown-it` 内部渲染**,agent 只关心 markdown |
| 不暴露 `project_id/user_id/author` | agent 不可能知道这些值,从 env 默认值取 |
| `category_id` 必填但描述指引"先 list" | 不做"分类名 → ID"自动映射(避免幻觉匹配错分类) |
| `status` 默认 1(草稿) | 安全默认值,鼓励两步走"先创建草稿、确认后再 publish",防误发 |

### 4.5 鉴权与 Token 流程

**为什么不让 MCP 持有账号密码自动登录?**

后端最近加了 `admin.totp.required` 强制 TOTP 开关(见近期 commit `ebbf3ba`)。开启后:
- 未绑定 TOTP 的账号 → 登录返回 `stage=bind_required`,**需要扫码绑定**,纯后端流程跑不通
- 已绑定的账号 → 登录返回 `stage=totp_required`,**需要人输 6 位动态码**

因此 MCP **不能持有用户名密码做自动登录**。改为:

| 阶段 | 谁来做 |
|---|---|
| 首次登录(含 TOTP) | 用户在 gotribecms 前端手动登录一次 |
| 拷贝 token | 用户从浏览器开发者工具拷贝 `access_token`(**只需要这一个**,gotribecms refresh 接口不消费独立 refresh_token) |
| 配置到 MCP | 写到 `GOTRIBE_ACCESS_TOKEN` env |
| 后续 refresh | **MCP 自动**:401 时把当前(过期的)access_token 放 Authorization Bearer 调 `POST /api/base/refreshToken`(后端用 `VerifyAccessTokenWithoutExpiry` 校验签名),拿到新 access_token,重试原请求一次 |
| refresh 也失败 | MCP 报错"请重新登录获取 token",用户重走上述流程,**不静默** |

**Token 持久化**:
- MCP 进程内存中保存最新的 access_token(refresh 后会更新)
- **不写本地文件**(避免在用户机器上落地 token,降低泄露面)
- 进程退出 = 内存 token 丢失,下次启动从 env 读初始值

### 4.6 错误处理

后端统一格式:`{ code, message, details, request_id }`(见 `internal/core/errs/code.go`)。

MCP 错误翻译表(摘要):

| 后端 code | MCP 返回给 agent 的消息 |
|---|---|
| `unauthorized` | "认证失败,token 已失效。请刷新 token 或重新登录后更新环境变量" |
| `forbidden` | "权限不足:此账号无 post 管理权限,请确认 Casbin 角色配置" |
| `not_found` | "资源不存在:{资源类型} id={id}" |
| `conflict` | "{字段} 已存在(如 slug 重复)" |
| `bad_request` | "参数错误:{details.field} — {message}" |
| `rate_limit_exceeded` | "请求过快,请稍后重试" |

**严守 CLAUDE.md 第 4 节**:不静默吞错,不假装成功,不兜底掩盖。

### 4.7 上传 tool 适配

后端 `POST /api/resource/upload` 返回(见 `internal/admin/resource/dto/response.go:60-65`):

```json
{ "upload": { "file_ext": "jpg", "key": "uploads/2026/xxx.jpg", "domain": "https://cdn...", "file_type": 1 } }
```

**问题**:`domain` 来自 `cdnDomain` 注入,**可能为空字符串**;`key` 是相对路径。

**MCP 处理**:
1. 拼接完整 URL:`${domain || GOTRIBE_CDN_DOMAIN_OVERRIDE}/${key}`
2. 两个都为空 → 报错"未配置 CDN 域名,请设置 `GOTRIBE_CDN_DOMAIN_OVERRIDE`",**不返回不完整 URL**
3. 返回给 agent:`{ url: "https://...", file_ext, file_type }`

### 4.8 用户接入示例

**Claude Code** (`~/.claude/settings.json` 或项目 `.claude/settings.json`):

```json
{
  "mcpServers": {
    "gotribe": {
      "command": "npx",
      "args": ["-y", "@gotribe/mcp"],
      "env": {
        "GOTRIBE_API_BASE_URL": "https://cms.gotribe.cn",
        "GOTRIBE_ACCESS_TOKEN": "eyJ...",
        "GOTRIBE_REFRESH_TOKEN": "eyJ...",
        "GOTRIBE_DEFAULT_PROJECT_ID": "1",
        "GOTRIBE_DEFAULT_USER_ID": "1",
        "GOTRIBE_DEFAULT_AUTHOR": "麻凡"
      }
    }
  }
}
```

**Codex / OpenCode**:配置语法略有不同,但都支持同样的 `command + args + env` 三段式。README 给三家具体示例。

## 5. 影响(Impact)

| 维度 | 影响 |
|---|---|
| **gotribecms 主项目代码** | **零改动** |
| **gotribecms 主项目文档** | README 加一节"AI Agent 接入(MCP)",链接到 gotribe-mcp 仓库 |
| **API 协议** | 零变更。MCP 是消费者,严格按现有 API 走 |
| **数据库** | 零变更 |
| **部署** | 主项目不需重新部署。MCP 在用户侧运行 |
| **认证体系** | 不引入新认证方式,复用现有 JWT + refresh |
| **未来兼容** | 主项目 API 字段如有变更,MCP 需同步并 bump 版本。建议主项目改 API 时打 `mcp-breaking` 标签 |
| **新增依赖(仅 MCP 仓库)** | `@modelcontextprotocol/sdk`、`markdown-it`、`zod`、`undici`(fetch),~4 个生产依赖 |
| **安全面** | 用户机器上的 env 变量持有 token。建议 README 提醒用户用 secret manager / 1Password CLI 注入 |

## 6. 验证(Validation)

| 类型 | 方式 |
|---|---|
| **类型 / 编译** | `tsc --noEmit` 通过 |
| **Lint** | `eslint` + `prettier` 通过 |
| **单元测试** | `vitest`:tool schema 校验、markdown 渲染、错误码翻译、URL 拼接 |
| **集成测试** | 本地起一个 gotribecms 实例(可选 docker-compose),跑 `list_posts → create_post → publish_post → get_post → delete_posts` 完整链路 |
| **客户端联调** | 在 Claude Code 实际配置 MCP,自然语言指令验证 8 个 tool 全部可达 |
| **错误路径** | 故意配错 token / 故意传错 category_id / 故意上传超大文件,验证报错清晰、不静默 |
| **三家 agent 兼容** | Claude Code 必测;Codex、OpenCode 各跑一次"创建草稿"流程作为冒烟 |

## 7. 回滚(Rollback)

| 场景 | 回滚 |
|---|---|
| 发包后发现严重 bug | 72h 内可 `npm unpublish`;否则发布 patch 修复版本(用户 `npx -y` 自动拉最新) |
| 用户配置后出现安全问题 | 在 gotribecms 后台撤销对应账号 token,主服务端立即拒绝 |
| 项目彻底放弃 MCP 方向 | 仓库归档,主项目零影响。用户改回手工 / curl 流程 |

## 8. 验收标准(Definition of Done)

- [ ] `gotribe-mcp` 仓库已建立(平层目录、独立 git repo)
- [ ] 8 个 P0 tool 全部实现
- [ ] 单测覆盖 schema / markdown / 错误翻译 / URL 拼接,关键路径覆盖率 ≥ 70%
- [ ] `npx -y @gotribe/mcp` 本地能启动(stdio 模式)
- [ ] README 给出 Claude Code / Codex / OpenCode 三家配置示例 + 首次 token 获取流程
- [ ] 在 Claude Code 实际跑通"创建草稿 → 发布 → 列表 → 删除"全链路
- [ ] 错误路径(token 失效、参数错、网络断、上传 domain 为空)都有清晰报错,不静默不兜底
- [ ] CI(build + test + lint)通过
- [ ] npm 包发布到 `@gotribe/mcp`,版本 `0.1.0`

## 9. 实施分阶段(Milestones)

| 阶段 | 内容 | 预估 |
|---|---|---|
| **M1 骨架** | 仓库初始化、tsconfig、ESLint、MCP SDK 接入、stdio 跑通 hello-world tool | 0.5 天 |
| **M2 HTTP + Auth** | http.ts、auth.ts、token refresh、错误翻译 | 0.5 天 |
| **M3 P0 Tools** | 8 个 tool 实现 + 单测 | 1.5 天 |
| **M4 联调 + README** | Claude Code 实测、写文档、修 bug | 1 天 |
| **M5 发包** | npm publish `@gotribe/mcp@0.1.0` | 0.5 天 |
| **合计** | | **~4 天** |

## 10. 待用户确认事项

| # | 问题 | 当前默认 |
|---|---|---|
| Q1 | npm `@gotribe` org **已确认拿到** ✓ | — |
| Q2 | 仓库路径 `/Users/mengmian/Develop/app/gotribe/gotribe-mcp` 是否 OK | 用户已同意"平层目录" |
| Q3 | `GOTRIBE_DEFAULT_PROJECT_ID/USER_ID/AUTHOR` 用 env 配 vs. MCP 启动时自动查 | **建议先用 env**(简单),后续如有多 project 切换需求再加 tool |
| Q4 | `markdown-it` 默认配置是否够用,要不要支持 GFM / 数学公式 / 代码高亮 | **建议默认开 GFM + 代码高亮**(`@shikijs/markdown-it`?),不开数学公式 |
| Q5 | LICENSE 用 MIT 吗 | 默认 MIT |
| Q6 | CI 用 GitHub Actions 吗 | 默认 GitHub Actions |
| Q7 | 是否需要在 gotribecms 主项目 README 也加 MCP 引导 | **建议加**,但属于本提案之外的小改动 |

---

## 用户确认

请用户在以下选项打勾或回复:

- [ ] **批准提案,按 §9 分阶段实施**
- [ ] **批准提案,但需要先调整**(请指出哪一节)
- [ ] **不批准 / 需要重新讨论**(请说明)

确认前 Claude 不会创建任何 `gotribe-mcp/` 仓库内文件,也不会写任何 MCP 代码。
