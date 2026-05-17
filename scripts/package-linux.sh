#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

ARCH="${1:-amd64}"
case "$ARCH" in
  amd64)
    API_BIN="gotribe-api-linux-amd64"
    ADMIN_BIN="gotribe-admin-linux-amd64"
    MAKE_TARGETS=(build-linux build-admin-linux)
    PACKAGE_NAME="gotribe-linux-amd64.tar.gz"
    ;;
  arm64)
    API_BIN="gotribe-api-linux-arm64"
    ADMIN_BIN="gotribe-admin-linux-arm64"
    MAKE_TARGETS=(build-linux-arm64 build-admin-linux-arm64)
    PACKAGE_NAME="gotribe-linux-arm64.tar.gz"
    ;;
  *)
    echo "unsupported arch: $ARCH (use amd64 or arm64)" >&2
    exit 1
    ;;
esac

RELEASE_DIR="$ROOT_DIR/release"
STAGE_DIR="$RELEASE_DIR/gotribe"
PACKAGE_PATH="$RELEASE_DIR/$PACKAGE_NAME"

echo "==> build linux binaries ($ARCH)"
make "${MAKE_TARGETS[@]}"

echo "==> stage package"
rm -rf "$STAGE_DIR"
mkdir -p "$STAGE_DIR/bin" "$STAGE_DIR/configs"

cp "$ROOT_DIR/bin/$API_BIN" "$STAGE_DIR/bin/gotribe-api"
cp "$ROOT_DIR/bin/$ADMIN_BIN" "$STAGE_DIR/bin/gotribe-admin"
chmod +x "$STAGE_DIR/bin/gotribe-api" "$STAGE_DIR/bin/gotribe-admin"

cp "$ROOT_DIR/configs/rbac_model.conf" "$STAGE_DIR/configs/rbac_model.conf"
cp -R "$ROOT_DIR/migrations" "$STAGE_DIR/migrations"

cat > "$STAGE_DIR/README.deploy.txt" <<'EOF'
This package intentionally does not include configs/config.yaml.
Keep the existing server config.yaml in place and only replace binaries plus
support files such as migrations/ and configs/rbac_model.conf when needed.

Typical server commands:
  cd /home/mafan/workspace/deploy/gotribe
  tar -xzf /path/to/gotribe-linux-amd64.tar.gz --strip-components=1
  chmod +x bin/gotribe-api bin/gotribe-admin
  pm2 restart gotribe-api
  pm2 restart gotribe-admin
  pm2 save
EOF

echo "==> create archive"
rm -f "$PACKAGE_PATH"
COPYFILE_DISABLE=1 tar --no-xattrs -C "$RELEASE_DIR" -czf "$PACKAGE_PATH" gotribe

echo "package created: $PACKAGE_PATH"
