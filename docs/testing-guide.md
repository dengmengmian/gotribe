# Testing Guide

## 目标

这份文档说明项目中的测试策略、测试分层、以及如何为新功能编写测试。

项目测试以**单元测试**为主，辅以**基准测试**衡量关键路径性能。

## 测试目录结构

```text
gotribe/
├── internal/
│   ├── example/service/example_test.go   # service 层单元测试（mock）
│   ├── health/handler/health_test.go     # handler 层单元测试（httptest）
│   └── ...                               # 其他单元测试
└── tests/bench/                          # 基准测试
    ├── jwt_bench_test.go
    ├── password_bench_test.go
    ├── cache_bench_test.go
    └── http_handler_bench_test.go
```

## 单元测试

### 什么时候写单元测试

- service 层包含业务逻辑分支（如条件判断、错误处理）
- 纯函数或工具函数（如格式化、校验、转换）
- 不需要数据库或 Redis 即可验证的逻辑

### 单元测试标准

| 项目 | 要求 |
|------|------|
| 依赖 | 不允许连接真实数据库或 Redis |
| 外部依赖 | 通过接口 mock（如 `testify/mock`） |
| HTTP 测试 | 使用 `httptest.NewRecorder` + `gin.TestMode` |
| 运行速度 | 单个测试应在 1 秒内完成 |

### 示例：mock 外部服务

参考 `internal/example/service/example_test.go`：

```go
type mockPostReader struct {
    mock.Mock
}

func (m *mockPostReader) GetSummaries(ctx context.Context, projectID string, postIDs []string) (map[string]postview.Summary, error) {
    args := m.Called(ctx, postIDs)
    return args.Get(0).(map[string]postview.Summary), args.Error(1)
}
```

### 示例：handler HTTP 单元测试

参考 `internal/health/handler/health_test.go`：

```go
w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/livez", nil)
engine.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
```

## 为新接口编写测试

在对应模块的 `service/*_test.go` 或 `handler/*_test.go` 中添加用例。

## Makefile 快捷命令

```bash
# 只跑单元测试（3-5 秒）
make test-unit

# 构建
make build
make build-admin
```

## 性能测试

项目提供 Go Benchmark 用于测量关键路径的纯计算开销，无需外部依赖。

```bash
# 运行全部 benchmark
make bench

# 只跑 HTTP handler 级别
make bench-http

# 带 CPU / 内存分析
make bench-mem
```

测试文件位于 `tests/bench/`：

| 文件 | 测试内容 |
|------|---------|
| `jwt_bench_test.go` | JWT 生成 / 解析 |
| `password_bench_test.go` | bcrypt 哈希 / 比对 |
| `cache_bench_test.go` | Redis 缓存读写（hit/miss） |
| `http_handler_bench_test.go` | 完整 handler 链路（health/post list/auth） |

示例输出解读：

```
BenchmarkJWT_GenerateAccessToken-8    500000    2341 ns/op    512 B/op    8 allocs/op
```

- `2341 ns/op` — 单次操作耗时 2.3 微秒
- `512 B/op` — 单次操作分配 512 字节
- `8 allocs/op` — 单次操作 8 次堆分配

### pprof 性能分析

生成 CPU / 内存 profile：

```bash
# 运行 benchmark 同时生成 profile
make bench-mem

# 交互式查看热点
go tool pprof tests/bench/profiles/cpu.out
# (pprof) top10
# (pprof) list YourFunction
# (pprof) web          # 生成 SVG 火焰图

# 内存分配分析
go tool pprof tests/bench/profiles/mem.out
```

也可对运行中的服务实时采集：

```bash
# 开启 pprof 端口（在 main.go 中添加 _ "net/http/pprof"）
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof
```

### 已有基准

如果修改了中间件、日志、限流或可观测性，建议跑一次对比：

```bash
go test ./internal/observability -bench BenchmarkHTTP_ -benchmem -run '^$'
```

对比 `NoObservability` vs `WithObservability` 的 ns/op 和 allocs/op，评估中间件额外开销。

## CI/CD 集成

### GitHub Actions 示例

```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Unit Tests
        run: go test ./... -short -count=1
```
