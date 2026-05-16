# 贡献指南

感谢你对 gotribe 的关注！

## 问题反馈

- 通过 GitHub Issues 提交 Bug 报告和功能建议
- 提交前请先搜索已有 issue，避免重复
- Bug 报告应包含：环境信息（Go 版本、数据库版本）、复现步骤、预期行为与实际行为

## 提交代码

1. Fork 本仓库
2. 基于 `main` 分支创建你的 feature 分支：`git checkout -b feature/your-feature`
3. 遵守项目的开发约定，详见 `README.md` 和 `docs/development-guide.md`
4. 确保测试通过：`make test-unit`
5. 提交前运行静态检查：`make vet`
6. 使用清晰的中文或英文提交信息，描述本次变更做了什么、为什么
7. 推送到你的 fork 并发起 Pull Request

## Pull Request 要求

- PR 标题简洁明确，描述本次变更的目的
- 关联相关 Issue（如有）
- 确保 CI 所有检查通过
- 新增功能应有对应的单元测试
- 变更 API 接口需同步更新 `api/openapi.yaml` 或 Swagger 文档

## 开发规范

项目采用标准的 handler → service → repository 三层架构：

- `handler`：HTTP 输入输出，不操作数据库
- `service`：业务逻辑、事务边界、缓存策略
- `repository`：数据访问

详细规范请参考：
- `docs/development-guide.md`：新人开发手册
- `docs/example-module.md`：标准模块参考实现
- `docs/testing-guide.md`：测试说明
- `docs/architecture.md`：架构设计

## 行为准则

本项目遵循 [Contributor Covenant 行为准则](CODE_OF_CONDUCT.md)。
