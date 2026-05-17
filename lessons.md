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
