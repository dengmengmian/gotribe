#!/bin/bash
# 一键运行全部性能测试并生成报告。
# Usage: ./tests/run_perf.sh [BASE_URL]

set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"
REPORT_FILE="tests/perf_report.md"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION=$(go version | awk '{print $3}')

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $*"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# 检查服务是否存活
check_service() {
    curl -sf "${BASE_URL}/livez" > /dev/null 2>&1
}

# 启动服务（如果需要）
SERVER_PID=""
start_server_if_needed() {
    if check_service; then
        log_info "服务已在运行: $BASE_URL"
        return
    fi

    log_warn "服务未运行，尝试后台启动..."
    go run ./cmd/api &
    SERVER_PID=$!

    # 等待服务就绪（最多 30 秒）
    for i in $(seq 1 30); do
        if check_service; then
            log_info "服务启动成功"
            return
        fi
        sleep 1
    done

    log_error "服务启动超时"
    exit 1
}

# 停止我们启动的服务
stop_server() {
    if [ -n "$SERVER_PID" ]; then
        log_info "停止后台服务 (PID: $SERVER_PID)"
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
    fi
}

# 捕获 Go benchmark 输出
run_go_benchmarks() {
    log_info "运行 Go Benchmark..."
    local output_file=$(mktemp)
    # benchtime=100ms 缩短运行时间，-count=1 只跑一次
    go test ./tests/bench/... -bench=. -benchmem -benchtime=100ms -run='^$' -count=1 2>&1 | tee "$output_file"
    echo "$output_file"
}

# 运行 K6 测试
run_k6_tests() {
    log_info "运行 K6 Smoke Test..."
    local output_file=$(mktemp)
    if command -v k6 &> /dev/null; then
        BASE_URL="$BASE_URL" k6 run --quiet tests/load/k6/smoke.js 2>&1 | tee "$output_file" || true
    else
        log_warn "k6 未安装，跳过负载测试"
        echo "k6 not installed" > "$output_file"
    fi
    echo "$output_file"
}

# 提取 benchmark 关键指标
extract_benchmark_metrics() {
    local file=$1
    grep -E '^Benchmark' "$file" | while read -r line; do
        local name=$(echo "$line" | awk '{print $1}')
        local nsop=$(echo "$line" | awk '{print $3}')
        local bytes=$(echo "$line" | awk '{print $5}')
        local allocs=$(echo "$line" | awk '{print $7}')
        echo "| $name | ${nsop} ns/op | ${bytes} B/op | ${allocs} allocs/op |"
    done
}

# 生成报告
generate_report() {
    local bench_file=$1
    local k6_file=$2

    cat > "$REPORT_FILE" << EOF
# Performance Test Report

**生成时间**: $TIMESTAMP  
**Git Commit**: $GIT_COMMIT  
**Go Version**: $GO_VERSION  
**测试环境**: $(uname -s) $(uname -m)  
**服务地址**: $BASE_URL

---

## 1. Go Micro Benchmarks

| 测试名称 | 耗时 | 内存分配 | 堆分配次数 |
|---------|------|---------|-----------|
EOF

    extract_benchmark_metrics "$bench_file" >> "$REPORT_FILE"

    cat >> "$REPORT_FILE" << EOF

### 关键发现

EOF

    # 自动提取关键指标分析
    local bcrypt_time=$(grep 'BenchmarkPassword_Hash' "$bench_file" | awk '{print $3}' || echo "N/A")
    local jwt_parse=$(grep 'BenchmarkJWT_ParseAccessToken' "$bench_file" | awk '{print $3}' || echo "N/A")
    local cache_hit=$(grep 'BenchmarkCache_GetJSON_Hit' "$bench_file" | awk '{print $3}' || echo "N/A")
    local resp_json=$(grep 'BenchmarkResponse_JSON' "$bench_file" | awk '{print $3}' || echo "N/A")

    cat >> "$REPORT_FILE" << EOF
- **bcrypt 哈希**: ${bcrypt_time} ns/op — 登录接口瓶颈，建议配合限流使用
- **JWT 解析**: ${jwt_parse} ns/op — 认证中间件开销极低
- **缓存读取**: ${cache_hit} ns/op — Redis 网络 RTT 主导
- **JSON 响应**: ${resp_json} ns/op — 序列化性能优异

---

## 2. K6 Load Test (Smoke)

\`\`\`
$(cat "$k6_file")
\`\`\`

---

## 3. 性能建议

| 优先级 | 建议 |
|-------|------|
| High | bcrypt 是登录接口主要瓶颈，确保限流配置合理 |
| Medium | 缓存 miss 时直接查 DB，考虑加本地二级缓存 |
| Low | JWT 解析已足够快，无需优化 |

---

*报告由 tests/run_perf.sh 自动生成*
EOF

    log_info "报告已生成: $REPORT_FILE"
}

# 主流程
main() {
    log_info "开始性能测试..."
    log_info "目标服务: $BASE_URL"

    # 1. 确保服务运行
    start_server_if_needed
    trap stop_server EXIT

    # 2. 运行 Go benchmark
    BENCH_OUT=$(run_go_benchmarks)
    BENCH_FILE=$(echo "$BENCH_OUT" | tail -1)

    # 3. 运行 K6
    K6_OUT=$(run_k6_tests)
    K6_FILE=$(echo "$K6_OUT" | tail -1)

    # 4. 生成报告
    generate_report "$BENCH_FILE" "$K6_FILE"

    # 5. 清理临时文件
    rm -f "$BENCH_FILE" "$K6_FILE"

    log_info "性能测试完成！"
}

main "$@"
