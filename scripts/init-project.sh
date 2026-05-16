#!/usr/bin/env bash

# 本脚本用于基于当前模板快速初始化一个新的内部项目。
# 建议仅在私有复制仓库里使用，不建议在公开文档中暴露。
# 它会同步替换 Go module、import 路径、应用名、环境变量前缀、
# 指标前缀、K8s 资源名、默认域名以及相关文档/配置中的项目标识。
# 示例：
#   bash scripts/init-project.sh blog
#   bash scripts/init-project.sh github.com/acme/blog

set -euo pipefail

usage() {
	cat <<'EOF'
用法:
  bash scripts/init-project.sh <new-module> [--force] [--skip-test]

示例:
  bash scripts/init-project.sh blog
  bash scripts/init-project.sh github.com/acme/blog

说明:
  - <new-module> 可以是简单模块名，也可以是完整 Go module path。
  - 脚本会自动同步替换：
    - go.mod module
    - Go import 路径
    - Makefile / Dockerfile 中的 ldflags 路径
    - app 名称（例如 gotribe-api -> blog-api）
    - 环境变量前缀（例如 GOTRIBE_ -> BLOG_）
    - 指标前缀、文档、OpenAPI、K8s 资源名、默认域名
  - 默认会执行 go mod tidy 和 go test ./... 做一次回归验证。
  - 工作区不干净时默认拒绝执行，可加 --force 跳过。
EOF
}

require_cmd() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "缺少依赖命令: $1" >&2
		exit 1
	fi
}

FORCE=0
SKIP_TEST=0
NEW_MODULE=""

while [[ $# -gt 0 ]]; do
	case "$1" in
	--force)
		FORCE=1
		;;
	--skip-test)
		SKIP_TEST=1
		;;
	-h | --help)
		usage
		exit 0
		;;
	-*)
		echo "不支持的参数: $1" >&2
		usage
		exit 1
		;;
	*)
		if [[ -n "$NEW_MODULE" ]]; then
			echo "只能传入一个 new-module 参数" >&2
			usage
			exit 1
		fi
		NEW_MODULE="$1"
		;;
	esac
	shift
done

if [[ -z "$NEW_MODULE" ]]; then
	usage
	exit 1
fi

if [[ ! "$NEW_MODULE" =~ ^[A-Za-z0-9._/-]+$ ]]; then
	echo "new-module 只能包含字母、数字、点、下划线、中划线和斜杠" >&2
	exit 1
fi

NEW_SLUG="${NEW_MODULE##*/}"
if [[ ! "$NEW_SLUG" =~ ^[a-z][a-z0-9-]*$ ]]; then
	echo "模块名最后一段必须是小写 slug，例如 blog 或 my-blog" >&2
	exit 1
fi

require_cmd git
require_cmd go
require_cmd perl
require_cmd rg

ROOT="$(git rev-parse --show-toplevel 2>/dev/null || true)"
if [[ -z "$ROOT" ]]; then
	echo "当前目录不是 git 仓库，无法安全执行初始化" >&2
	exit 1
fi

cd "$ROOT"

if [[ "$FORCE" -ne 1 ]] && [[ -n "$(git status --porcelain)" ]]; then
	echo "工作区存在未提交变更，请先提交或暂存后再执行。需要强制执行可加 --force" >&2
	exit 1
fi

OLD_MODULE="$(go list -m -f '{{.Path}}' 2>/dev/null || awk '/^module /{print $2; exit}' go.mod)"
OLD_SLUG="${OLD_MODULE##*/}"
OLD_UPPER="$(printf '%s' "$OLD_SLUG" | tr '[:lower:]-' '[:upper:]_')"
NEW_UPPER="$(printf '%s' "$NEW_SLUG" | tr '[:lower:]-' '[:upper:]_')"

if [[ "$OLD_MODULE" == "$NEW_MODULE" ]]; then
	echo "当前 module 已经是 ${NEW_MODULE}，无需替换"
	exit 0
fi

OLD_APP_NAME="$(awk '/^APP_NAME[[:space:]]*:=[[:space:]]*/{print $3; exit}' Makefile 2>/dev/null || true)"
if [[ -z "$OLD_APP_NAME" ]]; then
	OLD_APP_NAME="${OLD_SLUG}-api"
fi
NEW_APP_NAME="${NEW_SLUG}-api"

OLD_HOST="${OLD_SLUG}.local"
NEW_HOST="${NEW_SLUG}.local"

echo "开始初始化项目名称..."
echo "  module: ${OLD_MODULE} -> ${NEW_MODULE}"
echo "  app:    ${OLD_APP_NAME} -> ${NEW_APP_NAME}"
echo "  env:    ${OLD_UPPER} -> ${NEW_UPPER}"

while IFS= read -r -d '' file; do
	[[ -f "$file" ]] || continue
	[[ "$file" == .git/* ]] && continue
	[[ "$file" == bin/* ]] && continue
	[[ "$file" == vendor/* ]] && continue
	[[ "$file" == .DS_Store ]] && continue

	if [[ -s "$file" ]] && ! grep -Iq . "$file"; then
		continue
	fi

	OLD_MODULE="$OLD_MODULE" \
	NEW_MODULE="$NEW_MODULE" \
	OLD_SLUG="$OLD_SLUG" \
	NEW_SLUG="$NEW_SLUG" \
	OLD_APP_NAME="$OLD_APP_NAME" \
	NEW_APP_NAME="$NEW_APP_NAME" \
	OLD_UPPER="$OLD_UPPER" \
	NEW_UPPER="$NEW_UPPER" \
	OLD_HOST="$OLD_HOST" \
	NEW_HOST="$NEW_HOST" \
	perl -0pi -e '
		s/\Q$ENV{OLD_MODULE}\E/$ENV{NEW_MODULE}/g;
		s/\Q$ENV{OLD_APP_NAME}\E/$ENV{NEW_APP_NAME}/g;
		s/\Q$ENV{OLD_HOST}\E/$ENV{NEW_HOST}/g;
		s/\Q$ENV{OLD_UPPER}\E/$ENV{NEW_UPPER}/g;
		s/\Q$ENV{OLD_SLUG}\E/$ENV{NEW_SLUG}/g;
	' "$file"
done < <(git ls-files -z)

go mod edit -module "$NEW_MODULE"
go mod tidy

if [[ "$SKIP_TEST" -ne 1 ]]; then
	go test ./...
fi

echo
echo "项目初始化完成。"
echo "建议你再确认这些内容："
echo "  1. git 远程仓库地址是否需要改名"
echo "  2. 当前目录名是否需要手动调整"
echo "  3. configs/ 和 deployments/ 里的示例域名、镜像名、项目 ID 是否符合新项目"
