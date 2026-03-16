# AI模型与论文 微信小程序

一个用于对比大模型参数和阅读AI领域最新论文的微信小程序。

## 功能特性

### 🚀 大模型对比
- 展示主流大模型（GPT-4o、Claude 3.5、Gemini、DeepSeek、Qwen等）的详细参数
- 支持多模型并排对比（类似Apple官网的iPhone对比）
- 对比维度包括：
  - 基本信息（提供商、发布日期、架构）
  - 能力支持（文本、图片、音频、视频、文件）
  - 上下文窗口与最大输出
  - 定价信息
  - 性能基准（MMLU、HumanEval、MT-Bench）

### 📚 AI论文阅读
- 最新AI领域论文展示
- 英文原文 + 中文翻译
- 论文分类浏览
- 关键词标签

## 技术栈

### 后端 (Go)
- **框架**: Gin
- **数据**: JSON 文件存储
- **API**: RESTful API

### 前端 (uni-app)
- **框架**: Vue 3 + uni-app
- **样式**: 原生 CSS
- **构建**: Vite

## 项目结构

```
ai-model-papers-miniapp/
├── backend/           # Go 后端
│   ├── data/          # JSON数据文件
│   ├── handlers/      # HTTP处理器
│   ├── models/        # 数据模型
│   ├── routers/       # 路由配置
│   └── main.go        # 入口文件
├── frontend/          # uni-app 前端
│   ├── src/
│   │   ├── pages/     # 页面
│   │   ├── types/     # TypeScript类型
│   │   └── App.vue
│   └── package.json
└── README.md
```

## 快速开始

### 1. 启动后端

```bash
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

## API 文档

### 大模型接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/models | 获取所有模型 |
| GET | /api/models/:id | 获取单个模型 |
| GET | /api/models/providers | 获取所有提供商 |
| POST | /api/models/compare | 对比多个模型 |

### 论文接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/papers | 获取所有论文 |
| GET | /api/papers/:id | 获取单篇论文 |
| GET | /api/papers/categories | 获取论文分类 |

## 数据说明

### 大模型数据包含
- **基本信息**: 名称、提供商、描述
- **能力支持**: 文本、图片、音频、视频、文件
- **参数规格**: 上下文窗口、最大输出token
- **定价**: 输入/输出价格（每百万token）
- **性能**: MMLU、HumanEval、MT-Bench等基准分数

### 论文数据包含
- 英文标题与摘要
- 中文翻译
- 作者与机构
- 关键词与分类
- arXiv链接

## 后续优化方向

1. **数据更新**: 接入arXiv API自动获取最新论文
2. **搜索功能**: 添加全文搜索
3. **用户系统**: 收藏、历史记录
4. **评论互动**: 用户对模型和论文的评价
5. **图表展示**: 性能对比图表可视化

## License

MIT
