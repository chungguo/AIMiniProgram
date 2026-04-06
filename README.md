# ModelLens

AI 大模型评测与论文研读平台 - 帮助开发者和研究者洞察 AI 模型能力

## 项目结构

```
ModelLens/
├── proto/                           # Protocol Buffers 定义（前后端共享）
│   └── modellens/
│       ├── common/v1/common.proto
│       ├── model/v1/model.proto
│       ├── paper/v1/paper.proto
│       └── analysis/v1/analysis.proto
├── backend/                         # Go 后端服务（tRPC-Go + PowerWeChat）
│   ├── cmd/server/main.go           # 服务入口
│   ├── internal/
│   │   ├── handler/                 # HTTP 处理器
│   │   ├── repository/              # 数据访问层
│   │   └── wechat/                  # PowerWeChat 客户端
│   ├── bin/server                   # 编译后的可执行文件
│   ├── Dockerfile
│   └── Dockerfile.prod
└── frontend/                        # uni-app 前端
    ├── src/composables/             # Vue3 Composition API 复用逻辑
    └── src/pages/                   # 页面
```

## 快速启动

### 方式一：Docker 部署（推荐）

```bash
# 1. 构建并启动
make docker-build
make docker-run

# 2. 测试服务
curl http://localhost:8000/api/health

# 3. 查看日志
make docker-logs

# 4. 停止服务
make docker-stop
```

或使用 docker-compose 直接启动：

```bash
# 开发环境
docker-compose up -d

# 生产环境
docker-compose -f docker-compose.prod.yml up -d
```

### 方式二：本地运行

#### 1. 启动后端服务

```bash
# 一键启动
./start-server.sh

# 或手动启动
cd backend
go run cmd/server/main.go
```

服务将在 `http://localhost:8000` 启动。

#### 2. 测试 API

```bash
# 运行测试脚本
./test-api.sh

# 或手动测试
curl http://localhost:8000/api/health
curl "http://localhost:8000/api/models?page=1&limit=5"
```

#### 3. 启动前端

```bash
cd frontend
npm install
npm run dev:h5      # H5 模式
npm run dev:mp-weixin  # 微信小程序模式
```

## API 端点

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | /api/health | 健康检查 |
| GET | /api/models | 获取模型列表（支持分页、筛选） |
| GET | /api/models/families | 获取模型家族列表 |
| GET | /api/models/detail/:id | 获取模型详情 |
| POST | /api/models/compare | 对比多个模型 |
| GET | /api/papers | 获取论文列表 |
| GET | /api/papers/detail/:id | 获取论文详情 |
| GET | /api/analysis/artificialanalysis | 获取评测数据 |

## 前端 Composables

项目中实现了多个可复用的 Vue3 Composables：

- **useLoading** - 加载状态管理
- **useListManager** - 列表数据管理（分页、刷新、加载更多）
- **useSimpleList** - 简单列表（无分页）
- **useTabFilter** - Tab 筛选管理
- **useDetail** - 详情页数据管理

使用示例：

```typescript
import { useListManager, useDetail } from '@/composables';

// 列表管理
const { list, loading, refresh, loadMore } = useListManager<Model>({
  fetcher: async (params) => {
    const res = await modelService.getModels(params);
    return { data: res.data, pagination: res.pagination };
  },
  pageSize: 10
});

// 详情管理
const { data: model, loadData } = useDetail<Model>({
  fetcher: modelService.getModelById
});
```

## Proto 代码生成

当修改 proto 文件后，需要重新生成代码：

```bash
# 安装 protoc 工具
# macOS: brew install protobuf
# Ubuntu: apt-get install -y protobuf-compiler

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install trpc.group/trpc-go/trpc-cmdline/protoc-gen-go-trpc@latest

# 生成代码
./scripts/generate-proto.sh
```

## Docker 部署

### 镜像构建

```bash
# 开发镜像（基于 Alpine，带 shell）
make docker-build

# 生产镜像（基于 Scratch，最小体积）
make docker-build-prod

# 手动构建
docker build -t modellens/backend:latest -f backend-trpc/Dockerfile .
docker build -t modellens/backend:prod -f backend-trpc/Dockerfile.prod ./backend-trpc
```

### 容器运行

```bash
# 启动开发环境
docker-compose up -d

# 启动生产环境
docker-compose -f docker-compose.prod.yml up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 镜像信息

- **开发镜像**: ~20MB（Alpine + Go binary）
- **生产镜像**: ~10MB（Scratch + Go binary）
- **基础镜像**: golang:1.21-alpine（构建阶段）

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| PORT | 服务端口 | 8000 |
| GIN_MODE | Gin 运行模式 | release |

## 技术栈

### 后端
- **tRPC-Go** - 微服务框架（当前使用 Gin 兼容实现）
- **Protocol Buffers** - 接口定义和数据序列化
- **Gin** - HTTP Web 框架
- **JSON** - 数据存储（支持 PostgreSQL 扩展）
- **Docker** - 容器化部署

### 前端
- **uni-app** - 跨端框架
- **Vue 3** - 前端框架
- **TypeScript** - 类型系统
- **TDesign UniApp** - UI 组件库
- **UnoCSS** - 原子化 CSS

## 开发规范

1. **Proto 先行** - 新增接口时先定义 proto 文件
2. **Composables 优先** - 复用逻辑抽离为 composables
3. **类型安全** - 全项目使用 TypeScript/Go 类型
4. **前后端共享** - 类型定义通过 proto 共享

## 许可证

MIT
