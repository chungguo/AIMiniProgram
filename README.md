# AIMiniProgram

AI 模型与论文微信小程序。

## 功能

- 大模型参数对比（GPT-4o、Claude、Gemini、DeepSeek 等）
- AI 论文阅读（arXiv CS.AI 每日更新）
- ArtificialAnalysis 评测数据集成

## 技术栈

**后端**: Go + Gin + PostgreSQL  
**前端**: Vue 3 + uni-app + TypeScript

## 目录

```
├── backend/     # Go 后端 API
├── frontend/    # uni-app 前端
└── docs/        # 文档（已归档）
```

## 启动

```bash
# 后端
cd backend
go run main.go

# 前端
cd frontend
npm run dev:h5
```

## API

- `GET /api/models` - 模型列表
- `GET /api/papers` - 论文列表
- `GET /api/analysis/artificialanalysis` - 评测数据

## 环境变量

`DATABASE_URL` - PostgreSQL 连接字符串
