# AIMiniProgram

AI 模型与论文微信小程序项目根目录。

## 子项目

| 目录 | 说明 | 技术栈 |
|------|------|--------|
| `backend/` | API 服务 | Go + Gin + PostgreSQL |
| `frontend/` | 小程序前端 | Vue 3 + uni-app |

## 快速开始

```bash
# 后端
cd backend && go run main.go

# 前端  
cd frontend && npm run dev:h5
```

## 数据库

PostgreSQL 必需，连接字符串通过 `DATABASE_URL` 环境变量设置。
