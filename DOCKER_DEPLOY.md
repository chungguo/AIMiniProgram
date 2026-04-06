# Docker 部署指南

## 快速开始

```bash
# 1. 克隆项目
git clone https://github.com/chungguo/ModelLens.git
cd ModelLens

# 2. 启动服务（一键部署）
make docker-run

# 3. 验证服务
curl http://localhost:8000/api/health
```

## 目录结构

```
ModelLens/
├── backend-trpc/
│   ├── Dockerfile           # 开发环境 Dockerfile
│   ├── Dockerfile.prod      # 生产环境 Dockerfile
│   └── .dockerignore        # Docker 构建忽略文件
├── docker-compose.yml       # 开发环境编排
├── docker-compose.prod.yml  # 生产环境编排
└── Makefile                 # Docker 命令快捷方式
```

## 镜像说明

### 开发镜像 (Dockerfile)

- 基于 `golang:1.21-alpine`
- 包含完整的 Alpine Linux 环境
- 支持 shell 调试
- 适合开发和测试

```bash
docker build -t modellens/backend:dev -f backend-trpc/Dockerfile .
```

### 生产镜像 (Dockerfile.prod)

- 基于 `scratch`（空镜像）
- 仅包含编译后的二进制文件
- 镜像体积最小（约 10MB）
- 非 root 用户运行
- 适合生产部署

```bash
docker build -t modellens/backend:prod -f backend-trpc/Dockerfile.prod ./backend-trpc
```

## 常用命令

### 使用 Makefile

```bash
# 构建镜像
make docker-build          # 开发镜像
make docker-build-prod     # 生产镜像

# 运行容器
make docker-run            # 开发环境
make docker-run-prod       # 生产环境

# 管理容器
make docker-stop           # 停止容器
make docker-clean          # 清理资源
make docker-logs           # 查看日志
make docker-shell          # 进入容器
```

### 使用 docker-compose

```bash
# 开发环境
docker-compose up -d              # 后台启动
docker-compose down               # 停止并删除
docker-compose logs -f backend    # 查看日志
docker-compose ps                 # 查看状态

# 生产环境
docker-compose -f docker-compose.prod.yml up -d
docker-compose -f docker-compose.prod.yml down
```

### 使用原生 Docker

```bash
# 运行容器
docker run -d \
  --name modellens-backend \
  -p 8000:8000 \
  -e GIN_MODE=release \
  -v $(pwd)/backend/data:/root/data:ro \
  modellens/backend:latest

# 查看日志
docker logs -f modellens-backend

# 停止容器
docker stop modellens-backend
docker rm modellens-backend
```

## 生产部署建议

### 1. 使用反向代理（Nginx/Caddy）

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. 使用 Docker Swarm

```bash
# 初始化 Swarm
docker swarm init

# 部署服务
docker stack deploy -c docker-compose.prod.yml modellens

# 查看服务
docker service ls
docker service logs modellens_backend
```

### 3. 使用 Kubernetes

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: modellens-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: modellens/backend:prod
        ports:
        - containerPort: 8000
        resources:
          limits:
            memory: "512Mi"
            cpu: "1000m"
```

### 4. CI/CD 集成

```yaml
# .github/workflows/docker.yml
name: Docker Build and Push

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker Image
        run: |
          docker build -t modellens/backend:${{ github.ref_name }} \
            -f backend-trpc/Dockerfile.prod ./backend-trpc
      
      - name: Push to Registry
        run: |
          docker push modellens/backend:${{ github.ref_name }}
```

## 监控与日志

### 查看容器日志

```bash
# 实时日志
docker-compose logs -f backend

# 最近 100 行
docker-compose logs --tail=100 backend

# 带时间戳
docker-compose logs -f --timestamps backend
```

### 资源监控

```bash
# 容器资源使用
docker stats

# 容器详情
docker inspect modellens-backend
```

## 故障排查

### 服务无法启动

```bash
# 检查日志
docker-compose logs backend

# 检查端口占用
netstat -tlnp | grep 8000

# 重启服务
docker-compose restart backend
```

### 数据不更新

```bash
# 检查数据卷挂载
docker inspect -f '{{ .Mounts }}' modellens-backend

# 重新加载数据
docker-compose restart backend
```

## 安全建议

1. **使用非 root 用户运行**（已在 Dockerfile.prod 中配置）
2. **限制容器资源**（CPU、内存）
3. **使用只读数据卷**
4. **定期更新基础镜像**
5. **扫描镜像漏洞**

```bash
# 镜像安全扫描
docker scan modellens/backend:latest
```
