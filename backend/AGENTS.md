# Backend

Go + Gin 后端服务。

## 职责

- 提供 RESTful API
- PostgreSQL 数据访问
- 业务逻辑处理

## 核心模块

| 模块 | 说明 |
|------|------|
| `handlers/` | HTTP 处理器 |
| `models/` | 数据模型定义 |
| `repository/` | 数据库访问层 |
| `routers/` | 路由配置 |
| `middleware/` | 中间件（限流、安全） |

## 数据表

- `model` - 大模型信息
- `arxiv_cs_ai` - 论文数据
- `artificialanalysis` - 评测数据
