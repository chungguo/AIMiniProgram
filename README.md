# AI模型与论文 微信小程序

一个用于对比大模型参数和阅读AI领域最新论文的微信小程序。

## 功能特性

### 🚀 大模型对比
- 展示主流大模型（GPT-4o、Claude 3.5、Gemini、DeepSeek、Qwen等）的详细参数
- 支持多模型并排对比（类似Apple官网的iPhone对比）
- 对比维度包括：
  - 基本信息（家族、发布日期、知识截止日期）
  - 能力特性（推理、工具调用、附件支持等）
  - 模态支持（输入/输出：文本、图片、音频、视频）
  - 上下文窗口与限制
  - 定价信息（含缓存、音频等细分价格）

### 📚 AI论文阅读
- 最新AI领域论文展示
- 英文原文 + 中文翻译
- 论文分类浏览
- 关键词标签

## 技术栈

### 后端 (Go)
- **框架**: Gin
- **数据**: PostgreSQL
- **API**: RESTful API

### 前端 (uni-app)
- **框架**: Vue 3 + uni-app
- **样式**: 原生 CSS
- **构建**: Vite

## 项目结构

```
ai-model-papers-miniapp/
├── backend/           # Go 后端
│   ├── handlers/      # HTTP处理器
│   ├── middleware/    # 中间件（限流、安全等）
│   ├── models/        # 数据模型
│   ├── repository/    # 数据仓库（PostgreSQL）
│   ├── routers/       # 路由配置
│   └── main.go        # 入口文件
├── docs/              # 文档
│   └── schema-mapping.md
├── frontend/          # uni-app 前端
│   └── src/
│       ├── pages/     # 页面
│       ├── services/  # API服务
│       ├── types/     # TypeScript类型
│       └── utils/     # 工具函数
└── README.md
```

## 快速开始

### 1. 启动后端

**要求**: PostgreSQL 数据库

```bash
# 设置数据库连接字符串
export DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"

cd backend

# 安装依赖
go mod tidy

# 运行服务
go run main.go

# 或编译运行
go build -o server main.go
./server
```

后端服务将在 `http://localhost:3000` 启动。

### 2. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 开发模式（H5）
npm run dev:h5

# 编译微信小程序
npm run dev:mp-weixin
```

## 数据库表结构

### model 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | varchar(100) | 主键 |
| name | varchar(100) | 模型名称 |
| family | varchar(100) | 模型家族（如 OpenAI、Anthropic） |
| attachment | boolean | 是否支持附件 |
| reasoning | boolean | 是否支持推理 |
| tool_call | boolean | 是否支持工具调用 |
| structured_output | boolean | 是否支持结构化输出 |
| temperature | boolean | 是否支持温度调节 |
| knowledge | varchar(24) | 知识截止日期 |
| release_date | varchar(24) | 发布日期 |
| last_updated | varchar(24) | 最后更新日期 |
| modalities_input | modality[] | 输入模态数组（text/image/audio/video/file） |
| modalities_output | modality[] | 输出模态数组 |
| open_weights | boolean | 权重是否开源 |
| cost_input | numeric(10,6) | 输入价格（$/1M tokens） |
| cost_output | numeric(10,6) | 输出价格（$/1M tokens） |
| cost_reasoning | numeric(10,6) | 推理价格 |
| cost_cache_read | numeric(10,6) | 缓存读取价格 |
| cost_cache_write | numeric(10,6) | 缓存写入价格 |
| cost_input_audio | numeric(10,6) | 音频输入价格 |
| cost_output_audio | numeric(10,6) | 音频输出价格 |
| limit_context | integer | 最大上下文窗口（tokens） |
| limit_input | integer | 最大输入限制（tokens） |
| limit_output | integer | 最大输出限制（tokens） |
| interleaved_field | varchar(50) | 推理内容字段名 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

### paper 表（论文）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | varchar(100) | 主键 |
| title | text | 英文标题 |
| title_cn | text | 中文标题 |
| abstract | text | 英文摘要 |
| abstract_cn | text | 中文摘要 |
| authors | text[] | 作者列表 |
| institutions | text[] | 机构列表 |
| publish_date | varchar(24) | 发布日期 |
| arxiv_url | text | arXiv链接 |
| pdf_url | text | PDF链接 |
| categories | text[] | 分类标签 |
| keywords | text[] | 关键词 |
| read_time | integer | 阅读时间（分钟） |
| language | varchar(10) | 语言 |

## API 文档

### 大模型接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/models | 获取所有模型（支持筛选、排序、分页） |
| GET | /api/models?search=keyword | 搜索模型 |
| GET | /api/models?family=OpenAI | 按家族筛选 |
| GET | /api/models/detail/:id | 获取单个模型详情 |
| GET | /api/models/families | 获取所有模型家族 |
| GET | /api/models/family/:family | 获取指定家族的模型 |
| POST | /api/models/compare | 对比多个模型 |
| GET | /api/models/meta/comparison-categories | 获取对比类别 |

#### 查询参数

- `family` - 按家族筛选
- `hasAttachment` - 是否支持附件（true/false）
- `hasReasoning` - 是否支持推理（true/false）
- `hasToolCall` - 是否支持工具调用（true/false）
- `openWeights` - 权重是否开源（true/false）
- `minContext` - 最小上下文窗口
- `maxCostInput` - 最大输入价格
- `sortBy` - 排序字段（name/family/costInput/costOutput/limitContext/releaseDate）
- `sortOrder` - 排序方向（asc/desc）
- `page` - 页码（默认1）
- `limit` - 每页数量（默认20）
- `search` - 搜索关键词

### 论文接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/papers | 获取所有论文（分页） |
| GET | /api/papers/detail/:id | 获取单篇论文 |
| GET | /api/papers/categories | 获取论文分类 |
| GET | /api/papers/featured/latest | 获取最新论文 |

## 环境变量

| 变量 | 说明 | 必需 |
|------|------|------|
| `DATABASE_URL` | PostgreSQL 连接字符串 | ✅ |
| `PORT` | 服务端口号 | 否（默认3000） |

## 后续优化方向

1. **数据更新**: 接入arXiv API自动获取最新论文
2. **搜索功能**: 添加全文搜索
3. **用户系统**: 收藏、历史记录
4. **评论互动**: 用户对模型和论文的评价
5. **图表展示**: 性能对比图表可视化

## License

MIT
