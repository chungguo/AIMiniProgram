# Backend

Go API 服务。

## 结构

```
handlers/    # HTTP 处理器
models/      # 数据模型
repository/  # 数据仓库
routers/     # 路由
middleware/  # 中间件
```

## 运行

```bash
export DATABASE_URL="postgres://user:pass@localhost/dbname"
go run main.go
```

## API

- `GET /api/models` - 模型列表
- `GET /api/papers` - 论文列表  
- `GET /api/analysis/artificialanalysis` - 评测数据
- `GET /api/health` - 健康检查
