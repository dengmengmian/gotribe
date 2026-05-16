# 本文件定义本地开发最常用的构建、格式化、测试和运行命令。
APP_NAME := gotribe-api
ADMIN_NAME := gotribe-admin
GOPROXY ?=
-include .env
DB_CONTAINER ?= gotribe-postgres
DB_USER ?= $(POSTGRES_USER)
DB_PASSWORD ?= $(POSTGRES_PASSWORD)
DB_NAME ?= $(POSTGRES_DB)
DB_SCHEMA ?= public
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X 'gotribe/internal/buildinfo.Version=$(VERSION)' -X 'gotribe/internal/buildinfo.Commit=$(COMMIT)' -X 'gotribe/internal/buildinfo.BuildTime=$(BUILD_TIME)'

.PHONY: run run-api run-admin run-all tidy fmt test-unit pre-deploy build build-api build-admin build-admin-web lint vet bench bench-http bench-mem dev-db-reset

# 本地直接启动 API 应用。
run:
	GOPROXY=$(GOPROXY) go run -ldflags="$(LDFLAGS)" ./cmd/api

# run-api 是 run 的别名，方便明确区分两个服务。
run-api: run

# 本地直接启动 Admin 后台应用。
run-admin:
	GOPROXY=$(GOPROXY) go run -ldflags="$(LDFLAGS)" ./cmd/admin

# 同时启动 API 和 Admin 两个服务（含前端构建）。
run-all: build-admin-web
	GOPROXY=$(GOPROXY) go run -ldflags="$(LDFLAGS)" ./cmd/api & \
	GOPROXY=$(GOPROXY) go run -ldflags="$(LDFLAGS)" ./cmd/admin & \
	wait

# 同步并整理 Go 依赖。
tidy:
	GOPROXY=$(GOPROXY) go mod tidy

# 统一格式化所有 Go 源文件。
fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

# 运行单元测试（不含集成测试，速度快）。
test-unit:
	GOPROXY=$(GOPROXY) go test ./... -short -count=1 -skip 'TestIntegrationSuite'

# 运行集成测试（需要 Docker，自动启动 PostgreSQL + Redis 容器）。
integration-test:
	GOPROXY=$(GOPROXY) go test ./tests/integration/... -v -count=1

# 运行全部测试（单元 + 集成）。
test: test-unit integration-test

# 上线前全量检查：格式化 + 代码检查 + 单元测试 + 集成测试。
pre-deploy: fmt vet test-unit integration-test
	@echo "✅ all checks passed, ready to deploy"

# 静态代码检查。
vet:
	GOPROXY=$(GOPROXY) go vet ./...

# 运行 golangci-lint（需先安装: brew install golangci-lint）。
lint:
	golangci-lint run ./...

# 运行 golangci-lint 并自动修复部分问题。
lint-fix:
	golangci-lint run ./... --fix

# 构建 API 可执行文件到 bin 目录（当前平台）。
build:
	GOPROXY=$(GOPROXY) go build -ldflags="$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/api

# build-api 是 build 的别名，方便明确区分两个服务。
build-api: build

# 构建 Admin 后台可执行文件到 bin 目录（当前平台）。
build-admin:
	GOPROXY=$(GOPROXY) go build -ldflags="$(LDFLAGS)" -o bin/$(ADMIN_NAME) ./cmd/admin

# 构建 admin 前端静态资源（需要 pnpm）。
build-admin-web:
	test -f web/admin/.env || cp web/admin/.env.example web/admin/.env
	cd web/admin && pnpm install && pnpm build

# 本地一键启动全套依赖（PostgreSQL + Redis + API + Admin）。
dev-up:
	cp -n .env.example .env 2>/dev/null || true
	docker-compose up -d --build

# 停止全套本地依赖。
dev-down:
	docker-compose down

# 查看本地服务日志。
dev-logs:
	docker-compose logs -f

# 清理本地数据卷（会删除 PostgreSQL 和 Redis 数据，谨慎使用）。
dev-clean:
	docker-compose down -v

# 重置本地开发数据库 schema。用于反复清库重跑 migration + seed。
# 默认使用 docker-compose 的 gotribe-postgres；如使用公共容器可覆盖：
# make dev-db-reset DB_CONTAINER=common-postgres DB_USER=develop DB_NAME=develop
dev-db-reset:
	@test -n "$(DB_CONTAINER)" || (echo "DB_CONTAINER is required" && exit 1)
	@test -n "$(DB_USER)" || (echo "DB_USER is required" && exit 1)
	@test -n "$(DB_NAME)" || (echo "DB_NAME is required" && exit 1)
	docker exec -e PGPASSWORD="$(DB_PASSWORD)" $(DB_CONTAINER) psql -U "$(DB_USER)" -d "$(DB_NAME)" -v ON_ERROR_STOP=1 -c "DROP SCHEMA IF EXISTS $(DB_SCHEMA) CASCADE; CREATE SCHEMA $(DB_SCHEMA);"
	@echo "✅ reset database schema $(DB_NAME).$(DB_SCHEMA); restart admin to run migrations and seed data"


# 交叉编译 Linux AMD64 版本（用于容器或服务器部署）。
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=$(GOPROXY) go build -ldflags="$(LDFLAGS) -s -w" -o bin/$(APP_NAME)-linux-amd64 ./cmd/api

# 交叉编译 Linux ARM64 版本（用于 ARM 服务器部署）。
build-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOPROXY=$(GOPROXY) go build -ldflags="$(LDFLAGS) -s -w" -o bin/$(APP_NAME)-linux-arm64 ./cmd/api

# 运行全部 benchmark 测试。
bench:
	GOPROXY=$(GOPROXY) go test ./tests/bench/... -bench=. -benchmem -run='^$$'

# 运行 HTTP handler 级别 benchmark。
bench-http:
	GOPROXY=$(GOPROXY) go test ./tests/bench/... -bench=BenchmarkHTTP_ -benchmem -run='^$$'

# 生成测试覆盖率报告。
coverage:
	@mkdir -p tests/coverage
	GOPROXY=$(GOPROXY) go test ./... -coverprofile=tests/coverage/coverage.out -skip 'TestIntegrationSuite'
	go tool cover -html=tests/coverage/coverage.out -o tests/coverage/coverage.html
	@echo "Coverage report: tests/coverage/coverage.html"
	@go tool cover -func=tests/coverage/coverage.out | tail -1

# 生成覆盖率报告（包含集成测试，需要 Docker）。
coverage-full:
	@mkdir -p tests/coverage
	GOPROXY=$(GOPROXY) go test ./... -coverprofile=tests/coverage/coverage-full.out
	go tool cover -html=tests/coverage/coverage-full.out -o tests/coverage/coverage-full.html
	@echo "Full coverage report: tests/coverage/coverage-full.html"
	@go tool cover -func=tests/coverage/coverage-full.out | tail -1

# 运行带 CPU 和内存分析的 benchmark（输出到 tests/bench/profiles/）。
bench-mem:
	@mkdir -p tests/bench/profiles
	GOPROXY=$(GOPROXY) go test ./tests/bench/... -bench=. -benchmem -memprofile=tests/bench/profiles/mem.out -cpuprofile=tests/bench/profiles/cpu.out -run='^$$'
	@echo "Profiles saved to tests/bench/profiles/"
	@echo "View CPU: go tool pprof tests/bench/profiles/cpu.out"
	@echo "View Mem: go tool pprof tests/bench/profiles/mem.out"

# 一键运行全部性能测试（Go benchmark + K6 smoke）并生成报告。
perf-report:
	@bash tests/run_perf.sh
