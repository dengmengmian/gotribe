# 本文件定义应用的容器镜像构建方式，使用多阶段构建生成精简运行镜像。
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /app
# 可通过 --build-arg GOPROXY=https://goproxy.cn,direct 指定代理。
ARG GOPROXY
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w -X 'gotribe/internal/buildinfo.Version=${VERSION}' -X 'gotribe/internal/buildinfo.Commit=${COMMIT}' -X 'gotribe/internal/buildinfo.BuildTime=${BUILD_TIME}'" \
    -o /out/gotribe-api ./cmd/api

# 运行阶段只保留最小运行时环境。
FROM alpine:3.22

WORKDIR /app
RUN apk add --no-cache tzdata
RUN adduser -D -g '' appuser

COPY --from=builder /out/gotribe-api /app/gotribe-api
# 保留示例配置，方便容器内排查配置项。
COPY configs/config.yaml.example /app/configs/config.yaml.example

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/gotribe-api"]
