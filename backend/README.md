# AI Model & Papers Backend

Go 后端服务，为微信小程序提供 API 支持。

## 项目结构

```
backend/
├── data/              # JSON 数据文件
│   ├── models.json    # 大模型数据
│   └── papers.json    # 论文数据
├── handlers/          # HTTP 处理器
│   ├── base.go        # 基础功能
│   ├── models.go      # 模型相关
│   ├── papers.go      # 论文相关
│   └── health.go      # 健康检查
├── models/            # 数据模型定义
│   └── types.go
├── routers/           # 路由配置
│   └── router.go
├── main.go            # 入口文件
└── go.mod             # Go 模块配置
```

## API 接口

### 健康检查
- `GET /api/health`

### 大模型接口
- `GET /api/models` - 获取所有模型（支持筛选和搜索）
- `GET /api/models/:id` - 获取单个模型详情
- `GET /api/models/providers` - 获取所有提供商
- `POST /api/models/compare` - 对比多个模型
- `GET /api/models/meta/comparison-categories` - 获取对比类别

### 论文接口
- `GET /api/papers` - 获取所有论文（支持分页和筛选）
- `GET /api/papers/:id` - 获取单篇论文详情
- `GET /api/papers/categories` - 获取论文分类
- `GET /api/papers/featured/latest` - 获取最新论文

## 运行

```bash
# 安装依赖
go mod tidy

# 运行
go run main.go

# 编译
 go build -o server main.go
```

服务器将在 `http://localhost:3000` 启动。
