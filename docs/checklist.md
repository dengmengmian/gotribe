# 提交前检查清单

每次提交代码前，建议按此清单自查，确保代码质量和文档一致性。

## 快速检查（必须）

```bash
# 一键执行全部检查
make pre-deploy
```

如果 `make pre-deploy` 通过，以下手工检查项可快速浏览确认。

## 1. 代码质量

- [ ] `make fmt` — 代码已格式化
- [ ] `go vet ./...` — 无静态分析警告
- [ ] `make test-unit` — 单元测试通过
- [ ] `make integration-test` — 集成测试通过（需要 Docker）
- [ ] `make coverage` — 新增代码有配套测试覆盖
- [ ] `.dockerignore` — 如有新增不应打入镜像的文件

## 2. 编译与构建

- [ ] `go build ./...` — 全量编译通过
- [ ] `make build` — 可执行文件构建成功
- [ ] 交叉编译（如修改构建脚本）：`make build-linux`

## 3. 文档同步

- [ ] `README.md` — 如修改了核心功能或接口，同步更新
- [ ] `docs/development-guide.md` — 如修改了开发流程或 Makefile 命令
- [ ] `docs/testing-guide.md` — 如新增或修改了测试策略
- [ ] `docs/deployment-guide.md` — 如修改了部署配置
- [ ] `docs/architecture.md` — 如修改了架构设计
- [ ] `api/openapi.yaml` — 如新增或修改了接口
- [ ] 新增模块：是否补充了 `docs/example-module.md` 级别的参考说明

## 4. 配置与密钥

- [ ] 未提交敏感信息（数据库密码、JWT 密钥等）
- [ ] 新配置文件已加入 `.gitignore`
- [ ] 配置示例文件已更新（`*.example` / `*.example.yaml`）

## 5. 新增模块检查（如适用）

- [ ] 目录结构符合规范：`dto/`, `handler/`, `model/`, `repository/`, `service/`
- [ ] handler 中注册了路由（`handler/routes.go`）
- [ ] provider 中装配了依赖（`internal/bootstrap/module_builders.go`）
- [ ] 包注释已添加（`// Package xxx ...`）
- [ ] 导出标识符有注释
- [ ] 单元测试或集成测试已补充

## 6. 性能影响（如适用）

- [ ] 修改了中间件/日志/限流：`make bench` 通过
- [ ] 修改了数据库查询：确认无 N+1、有索引
- [ ] 修改了缓存逻辑：确认缓存 key 设计合理

## 7. 提交信息

建议格式：

```
<type>(<scope>): <subject>

<body>
```

类型（type）：

| 类型 | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | 修复 |
| `docs` | 文档 |
| `style` | 格式调整（不影响代码逻辑）|
| `refactor` | 重构 |
| `test` | 测试 |
| `chore` | 构建/工具/配置 |

示例：

```
feat(post): add keyword search to post list

docs(api): update openapi.yaml with new endpoints

chore(makefile): add perf-report target
```

## 8. 最终确认

- [ ] `git status` — 只提交必要的文件
- [ ] 无意外的大文件提交
- [ ] 无调试代码（如 `fmt.Println`、`debugger`）残留
