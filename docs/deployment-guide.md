# Deployment Guide

## 目标

这份文档说明如何将 `gotribe-api` 构建为容器镜像并部署到不同环境。

项目提供三种部署方式：

| 方式 | 适用场景 | 复杂度 |
|------|---------|--------|
| **本地二进制** | 本地开发、调试 | 低 |
| **Docker Compose** | 本地/测试环境快速启动 | 低 |
| **Kubernetes** | 生产环境、云原生部署 | 中 |

## 环境要求

无论哪种部署方式，都需要提前准备：

- PostgreSQL 16+
- Redis 7+
- JWT 密钥（至少 32 位随机字符串）

## 本地二进制部署

### 交叉编译 Linux 版本

```bash
make build-linux
```

输出：`bin/gotribe-api-linux-amd64`

### 配置环境变量后启动

```bash
export GOTRIBE_APP_ENV=production
export GOTRIBE_APP_DEFAULT_PROJECT_ID=your-project-id
export GOTRIBE_AUTH_SECRET=your-long-random-secret
export GOTRIBE_DATABASE_HOST=your-db-host
export GOTRIBE_DATABASE_PORT=5432
export GOTRIBE_DATABASE_USERNAME=your-db-user
export GOTRIBE_DATABASE_PASSWORD=your-db-password
export GOTRIBE_DATABASE_DATABASE=your-db-name
export GOTRIBE_REDIS_ADDR=your-redis-host:6379

./bin/gotribe-api-linux-amd64
```

### 健康检查

```bash
curl http://localhost:8080/version
curl http://localhost:8080/livez
curl http://localhost:8080/readyz
curl http://localhost:8080/metrics
```

## Docker 部署

### 构建镜像

```bash
docker build -t gotribe-api:latest .
```

### 使用 Docker Compose 启动

修改 `docker-compose.yml` 中的环境变量（至少替换 `GOTRIBE_APP_DEFAULT_PROJECT_ID`、`GOTRIBE_DATABASE_USERNAME`、`GOTRIBE_DATABASE_PASSWORD`、`GOTRIBE_DATABASE_DATABASE`、`GOTRIBE_AUTH_SECRET`）：

```bash
docker compose up --build
```

### 配置说明

Docker Compose 中关键环境变量：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `GOTRIBE_DATABASE_HOST` | 数据库地址 | `host.docker.internal`（指向宿主机） |
| `GOTRIBE_REDIS_ADDR` | Redis 地址 | `host.docker.internal:6379` |
| `GOTRIBE_AUTH_SECRET` | JWT 密钥 | **必须替换** |

> 如果 PostgreSQL 和 Redis 也在 Docker 中运行，将 `host.docker.internal` 改为容器网络内的服务名。

## Kubernetes 部署

### 文件清单

```text
deployments/k8s/
├── configmap.yaml          # 非敏感配置
├── secret.example.yaml     # 敏感配置示例
├── deployment.yaml         # 应用部署
├── pdb.yaml                # PodDisruptionBudget
├── service.yaml            # 集群内服务暴露
├── ingress.yaml            # 外部入口
└── hpa.yaml                # 自动扩缩容
```

### 部署步骤

**Step 1**：基于 `secret.example.yaml` 创建真实的 Secret

```bash
cp deployments/k8s/secret.example.yaml deployments/k8s/secret.yaml
# 编辑 secret.yaml，替换 GOTRIBE_DATABASE_PASSWORD 和 GOTRIBE_AUTH_SECRET
kubectl apply -f deployments/k8s/secret.yaml
```

> **警告**：`secret.yaml` 包含真实密钥，已加入 `.gitignore`，请勿提交到仓库。

**Step 2**：修改 ConfigMap 中的配置

编辑 `deployments/k8s/configmap.yaml`，替换：

- `GOTRIBE_APP_DEFAULT_PROJECT_ID`
- `GOTRIBE_DATABASE_USERNAME`
- `GOTRIBE_DATABASE_DATABASE`
- 数据库 host（如果 PG 不在同一 namespace）
- Redis 地址

**Step 3**：应用所有资源

```bash
kubectl apply -f deployments/k8s/configmap.yaml
kubectl apply -f deployments/k8s/deployment.yaml
kubectl apply -f deployments/k8s/pdb.yaml
kubectl apply -f deployments/k8s/service.yaml
kubectl apply -f deployments/k8s/ingress.yaml
kubectl apply -f deployments/k8s/hpa.yaml
```

### 部署配置说明

#### Deployment

- **副本数**：默认 2，配合 HPA 自动扩缩
- **镜像**：`gotribe-api:latest`，实际部署时应替换为具体版本 tag
- **滚动发布**：`maxSurge=1`、`maxUnavailable=0`
- **资源限制**：已提供默认 `requests/limits`
- **安全上下文**：默认非 root、禁提权、只读根文件系统
- **健康探针**：
  - `livenessProbe` → `/livez`（进程存活检查）
  - `readinessProbe` → `/readyz`（依赖就绪检查，含 DB + Redis）
  - `startupProbe` → `/livez`（启动阶段保护）
- **metrics 抓取**：Pod annotation 已暴露 `/metrics`

#### ConfigMap vs Secret

| 类型 | 存放内容 |
|------|---------|
| ConfigMap | 环境标识、端口、数据库 host、非敏感连接信息 |
| Secret | 数据库密码、JWT 密钥 |

#### Ingress

- 示例域名：`gotribe.local`
- 实际部署时替换为真实域名
- 需要集群内已安装 Ingress Controller（如 nginx-ingress）

#### HPA（自动扩缩容）

- **最小副本**：2
- **最大副本**：5
- **触发条件**：CPU 使用率 > 70%

```bash
# 查看 HPA 状态
kubectl get hpa gotribe-api
```

#### PDB（PodDisruptionBudget）

- 默认 `minAvailable: 1`
- 用于节点维护、驱逐等场景下尽量保证服务可用

### 验证部署

```bash
# 查看 Pod 状态
kubectl get pods -l app=gotribe-api

# 查看服务
kubectl get svc gotribe-api

# 查看日志
kubectl logs -l app=gotribe-api --tail=100

# 端口转发到本地调试验证
kubectl port-forward svc/gotribe-api 8080:80
curl http://localhost:8080/version
curl http://localhost:8080/readyz
curl http://localhost:8080/metrics
```

## 生产部署检查清单

上线前确认：

- [ ] `GOTRIBE_AUTH_SECRET` 已替换为真实长随机串（≥32 位）
- [ ] 数据库连接信息正确且网络可达
- [ ] Redis 连接信息正确且网络可达
- [ ] Ingress 域名已替换为真实域名
- [ ] 镜像 tag 已固定为具体版本（非 `latest`）
- [ ] 构建版本号已通过 `ldflags` 注入
- [ ] `make pre-deploy` 全量测试通过
- [ ] `/version` 返回当前构建元数据
- [ ] `/readyz` 返回 ready 状态
- [ ] `/livez` 返回 ok 状态
- [ ] `/metrics` 可抓取
- [ ] 至少一个接口联调通过（如登录 → 查资料）
