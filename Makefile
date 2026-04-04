# AIMiniProgram tRPC-Go Makefile

.PHONY: help proto proto-go proto-ts clean build run docker-build docker-run docker-stop docker-clean

# 默认目标
help:
	@echo "Available targets:"
	@echo ""
	@echo "  Development:"
	@echo "    proto     - Generate protobuf code for all languages"
	@echo "    proto-go  - Generate Go code from protobuf"
	@echo "    proto-ts  - Generate TypeScript code from protobuf"
	@echo "    build     - Build the Go server"
	@echo "    run       - Run the Go server"
	@echo "    clean     - Clean generated files"
	@echo ""
	@echo "  Docker:"
	@echo "    docker-build    - Build Docker image"
	@echo "    docker-run      - Run with docker-compose"
	@echo "    docker-run-prod - Run production build"
	@echo "    docker-stop     - Stop containers"
	@echo "    docker-clean    - Clean Docker resources"
	@echo "    docker-logs     - View container logs"

# Proto 文件路径
PROTO_DIR := proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/aiminiprogram/*/v1/*.proto)

# ==================== Development Commands ====================

# 生成所有代码
proto: proto-go proto-ts

# 生成 Go 代码
proto-go:
	@echo "Generating Go code from protobuf..."
	@mkdir -p backend-trpc/internal/pb
	@protoc \
		--go_out=backend-trpc/internal/pb \
		--go_opt=paths=source_relative \
		--go-trpc_out=backend-trpc/internal/pb \
		--go-trpc_opt=paths=source_relative \
		--proto_path=proto \
		$(PROTO_FILES)
	@echo "Go code generated successfully!"

# 生成 TypeScript 代码
proto-ts:
	@echo "Generating TypeScript code from protobuf..."
	@mkdir -p frontend/src/proto
	@protoc \
		--ts_out=frontend/src/proto \
		--ts_opt=long_type_string,client_grpc1,server_grpc1 \
		--proto_path=proto \
		$(PROTO_FILES)
	@echo "TypeScript code generated successfully!"

# 构建 Go 服务
build:
	@echo "Building Go server..."
	@cd backend-trpc && go build -o bin/server cmd/server/main.go
	@echo "Build complete: backend-trpc/bin/server"

# 运行 Go 服务
run: build
	@./backend-trpc/bin/server

# 清理生成文件
clean:
	@echo "Cleaning generated files..."
	@rm -rf backend-trpc/internal/pb
	@rm -rf frontend/src/proto
	@rm -rf backend-trpc/bin
	@echo "Clean complete!"

# 安装依赖
deps:
	@echo "Installing Go dependencies..."
	@cd backend-trpc && go mod tidy
	@echo "Installing protobuf tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install trpc.group/trpc-go/trpc-cmdline/protoc-gen-go-trpc@latest
	@echo "Dependencies installed!"

# 检查工具
check:
	@echo "Checking required tools..."
	@which protoc || echo "ERROR: protoc not found"
	@which protoc-gen-go || echo "ERROR: protoc-gen-go not found"
	@which protoc-gen-go-trpc || echo "ERROR: protoc-gen-go-trpc not found"

# ==================== Docker Commands ====================

# 构建 Docker 镜像
docker-build:
	@echo "Building Docker image..."
	@docker build -t aiminiprogram/backend:latest -f backend-trpc/Dockerfile .
	@echo "Docker image built: aiminiprogram/backend:latest"

# 构建生产镜像
docker-build-prod:
	@echo "Building production Docker image..."
	@docker build -t aiminiprogram/backend:prod -f backend-trpc/Dockerfile.prod ./backend-trpc
	@echo "Production Docker image built: aiminiprogram/backend:prod"

# 运行开发环境
docker-run:
	@echo "Starting development containers..."
	@docker-compose up -d
	@echo "Services started:"
	@echo "  Backend: http://localhost:8000"
	@echo "  Health:  http://localhost:8000/api/health"

# 运行生产环境
docker-run-prod:
	@echo "Starting production containers..."
	@docker-compose -f docker-compose.prod.yml up -d
	@echo "Production services started:"
	@echo "  Backend: http://localhost:8000"

# 停止容器
docker-stop:
	@echo "Stopping containers..."
	@docker-compose down
	@docker-compose -f docker-compose.prod.yml down 2>/dev/null || true
	@echo "Containers stopped"

# 清理 Docker 资源
docker-clean:
	@echo "Cleaning Docker resources..."
	@docker-compose down -v --remove-orphans
	@docker rmi aiminiprogram/backend:latest 2>/dev/null || true
	@docker rmi aiminiprogram/backend:prod 2>/dev/null || true
	@echo "Docker resources cleaned"

# 查看日志
docker-logs:
	@docker-compose logs -f backend

# 进入容器
docker-shell:
	@docker-compose exec backend /bin/sh

# 推送镜像到仓库（需要配置 DOCKER_REGISTRY）
docker-push: docker-build-prod
	@echo "Pushing image to registry..."
	@docker tag aiminiprogram/backend:prod $(DOCKER_REGISTRY)/aiminiprogram/backend:$(VERSION)
	@docker push $(DOCKER_REGISTRY)/aiminiprogram/backend:$(VERSION)
	@echo "Image pushed to $(DOCKER_REGISTRY)/aiminiprogram/backend:$(VERSION)"
