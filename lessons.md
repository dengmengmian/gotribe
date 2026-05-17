# Lessons

本文件记录在与 AI 编程 Agent 合作过程中识别出的错误、遗漏、误判或返工点。
按 `CLAUDE.md` 第 7 节要求维护:每次新任务开始前必须读取,不允许重复同类错误。

---

## L001 — Proposal 设计字段时凭印象,未先核实代码事实

**问题**
2026-05-17,在编写 `proposals/001-gotribe-mcp.md` 时,Agent 把 `GOTRIBE_REFRESH_TOKEN`
列为 MCP 服务器的环境变量,并设计了"401 时用 refresh_token 调刷新接口"的流程。
M2 实施阶段实际读 `internal/auth/admin/handler/auth.go:84-116` 才发现:gotribecms 的
刷新流程是**用过期的 access_token 调 `POST /api/base/refreshToken`** —— 后端用
`VerifyAccessTokenWithoutExpiry` 只校验签名不校验过期,**根本不消费独立的 refresh_token**。
后端虽然内部有 refresh_token 概念(Redis 里 256-bit 随机串),但没暴露给 admin 前端使用。
导致 proposal §4.2 和 §4.5 写错,实施时返工修正。

**原因**
前置调研由 Explore subagent 完成,报告里 Token 刷新章节的"鉴权"字段写了 "需 Bearer token(即使已过期也可用于刷新)" —— 这句话其实**已经隐含**了"用 access_token 刷新而非 refresh_token",但 Agent 在写 proposal 时把它跟"常见 OAuth refresh_token 模式"混在一起,凭印象补了 refresh_token env 变量,没有再回头核实代码。

**规则**
在 proposal / 设计文档里**每一个具体字段、每一个具体流程**,必须**直接对应源码行号**,不允许凭"通用模式""一般来说"补充。
具体做法:
1. 设计文档草稿写完后,做一次"字段 ↔ 源码"反向核对 —— 每个字段是否在源码中实际存在 / 实际被消费;
2. 涉及鉴权 / refresh / token 这种"行业有常见模式"的部分,**最容易凭印象写错**,要逐字段确认;
3. Explore subagent 的报告里如有"即使...也可以"这类暗示性描述,要主动追问"那是不是就不需要 X 了?",不要默认补 X。

**适用范围**
- 写 proposal / OpenSpec 设计文档时
- 跨模块调研后的方案输出环节
- 任何"通用模式跟项目实际实现不一定一致"的领域(鉴权、缓存、错误处理、ID 生成、时间格式等)

---

## L002 — 见到"特殊格式"先查仓库里是否已有现成的转换代码

**问题**
2026-05-17,实现 `gotribe-mcp` 的 `create_post` tool 时,Agent 看到文章 `content` 字段是 Slate.js 结构化 JSON 而非 Markdown,**第一反应是"自己写一个 markdown → Slate 转换器"**,甚至计划好了支持哪些 markdown 语法、不支持哪些。用户提醒"用 Slate.js 转,不要自己转"才回头去搜,3 分钟内就在 `web/admin/src/lib/slate-markdown.ts` 找到一份 360 行、已经在生产环境运行、完整覆盖所有 editor schema 的现成实现(`markdownToSlate` / `slateToMarkdown` / `slateContentToHtml` / `slateToPlainText`)。如果按 Agent 原计划手写,会出现:
- 不支持后台已支持的 `block-quote` / `code-block` / `image` / `link` / `table` / `math` / `check-list` 等节点
- HTML 渲染规则跟前端不一致(width / align / 转义细节)
- 未来后台 editor 改 schema 后,MCP 也要跟着改,要维护两套规则

**原因**
看到 `content` 是 Slate JSON 这个"非标准"格式时,Agent 把它当成"未知问题"启动了从零设计模式,而没有先做一步反向搜索:"项目里有没有人已经处理过这个格式 → 一定有,因为 admin web 必须能编辑这些文章"。

**规则**
碰到任何"看起来很复杂、需要自己实现"的转换/序列化/适配逻辑,**先在仓库里 grep 一遍关键字**(库名、关键函数名、目标格式特征),确认没有现成实现再考虑自己写。常见场景:
- 富文本编辑器格式(Slate / Tiptap / Lexical / Quill / ProseMirror)
- 自定义数据格式(项目专属 DSL / 协议 buffer)
- 任何"业务里已经在用、但你刚接触"的格式

具体做法:
1. `grep -rn '<关键字>' --include='*.ts' --include='*.go'`
2. 看 `web/admin` 或 admin 前端目录(因为编辑器最可能用到转换)
3. 看 `package.json` 是否依赖相关库,顺着库名再 grep 用法
4. 找到现成实现后,**逐字复制并标注同步责任**(README/lessons 写明"如果上游变了 MCP 也要跟着变")

**适用范围**
- 实现"涉及业务特定数据格式"的客户端/SDK/MCP server 时
- 跨项目/跨语言共享数据 schema 的场景
- 任何"我以为没人做过、其实业务系统天天在用"的领域
