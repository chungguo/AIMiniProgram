# tRPC-Go 架构改造方案

## 1. 架构设计

```
AIMiniProgram/
├── proto/                           # 共享的 protobuf 定义
│   └── aiminiprogram/
│       ├── common/v1/common.proto   # 公共类型
│       ├── model/v1/model.proto     # 模型相关协议
│       ├── paper/v1/paper.proto     # 论文相关协议
│       └── analysis/v1/analysis.proto  # 评测数据协议
├── backend-trpc/                    # tRPC-Go 服务
│   ├── cmd/server/
│   │   └── main.go                  # 服务入口
│   ├── internal/
│   │   ├── service/                 # 业务逻辑层
│   │   │   ├── model_service.go
│   │   │   ├── paper_service.go
│   │   │   └── analysis_service.go
│   │   └── repository/              # 数据访问层
│   │       ├── model_repository.go
│   │       ├── paper_repository.go
│   │       └── analysis_repository.go
│   ├── go.mod
│   └── trpc_go.yaml                 # tRPC 配置
├── backend/                         # 原 Gin 服务（保留）
└── frontend/                        # 前端
    └── src/
        └── proto/                   # 生成的 TypeScript 代码
```

## 2. Proto 定义规范

### 2.1 命名规范
- 包名：`aiminiprogram.{service}.v1`
- 服务名：`{Service}Service`
- 方法名：动词 + 名词，如 `GetModel`, `ListPapers`
- 消息名：名词，如 `Model`, `Paper`, `ListModelsRequest`

### 2.2 字段编号规则
- 1-15：常用字段（单字节编码）
- 16-2047：普通字段
- 保留字段：19000-19999（protobuf 预留）

### 2.3 响应格式统一
```protobuf
message Response {
  bool success = 1;
  string message = 2;
  google.protobuf.Any data = 3;
}

// 或使用泛型风格
message PaginatedResponse {
  bool success = 1;
  string message = 2;
  repeated T data = 3;  // T 为具体类型
  Pagination pagination = 4;
}
```

## 3. 前后端共享方案

### 方案 A：Git Submodule（推荐）
```bash
# 创建独立的 proto 仓库
git submodule add https://github.com/chungguo/aiminiprogram-proto.git proto/
```

### 方案 B：Monorepo 内部共享
```
proto/                      # 根目录 proto/
├── go/                     # Go 生成的代码
├── ts/                     # TypeScript 生成的代码
└── *.proto                 # 原始 proto 文件
```

### 方案 C：NPM/Go Module 发布
- 前端：通过 npm 包引入
- 后端：通过 go module 引入

## 4. trpc-go 服务实现

### 4.1 服务注册
```go
import (
    "trpc.group/trpc-go/trpc-go/server"
    pb "aiminiprogram/proto/model/v1"
)

func main() {
    s := server.New(
        server.WithServiceName("aiminiprogram.model"),
    )
    
    pb.RegisterModelService(s, &modelServiceImpl{})
    
    s.Serve()
}
```

### 4.2 拦截器配置
```go
// 全局拦截器
server.Use(
    recovery.Interceptor(),      //  panic 恢复
    logging.Interceptor(),       // 日志
    ratelimit.Interceptor(),     // 限流
    auth.Interceptor(),          // 认证
)

// 服务级拦截器
service := pb.RegisterModelService(s, &impl{},
    server.WithFilter(middleware.ValidateRequest),
)
```

## 5. 前端调用方式

### 5.1 TypeScript 生成代码
```typescript
// 使用 protobuf-ts 或 ts-proto 生成
import { ModelServiceClient } from '@/proto/model/v1/service.client';
import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';

const transport = new GrpcWebFetchTransport({
  baseUrl: '/api',
  format: 'binary',
});

const client = new ModelServiceClient(transport);

// 调用
const { response } = await client.getModel({ id: 'gpt-4' });
```

### 5.2 或使用 HTTP 网关模式
```typescript
// trpc 自动生成 HTTP 网关
// 前端保持现有 HTTP 调用方式
fetch('/api/v1/models?id=gpt-4')
  .then(res => res.json())
  .then(data => console.log(data));
```

## 6. 迁移步骤

### 阶段 1：准备（1-2 天）
1. 创建 proto 定义文件
2. 设置 proto 生成工具链
3. 定义接口契约

### 阶段 2：后端改造（3-5 天）
1. 初始化 trpc-go 项目
2. 实现服务接口
3. 迁移业务逻辑
4. 配置拦截器

### 阶段 3：前端适配（2-3 天）
1. 生成 TypeScript 代码
2. 更新 API 调用层
3. 测试联调

### 阶段 4：切换（1 天）
1. 灰度发布
2. 监控指标
3. 全量切换

## 7. 工具链

### 7.1 Protobuf 编译器
```bash
# Go
protoc --go_out=. --go_opt=paths=source_relative \
       --go-trpc_out=. --go-trpc_opt=paths=source_relative \
       proto/model.proto

# TypeScript
protoc --ts_out=frontend/src/proto \
       --ts_opt=long_type_string \
       proto/model.proto
```

### 7.2 Makefile 自动化
```makefile
.PHONY: proto proto-go proto-ts

proto: proto-go proto-ts

proto-go:
	protoc --go_out=backend --go_opt=paths=source_relative \
	       --go-trpc_out=backend --go-trpc_opt=paths=source_relative \
	       proto/*.proto

proto-ts:
	protoc --ts_out=frontend/src/proto \
	       --ts_opt=long_type_string \
	       proto/*.proto
```

## 8. 性能优化

### 8.1 连接管理
- 使用连接池
- 启用 keepalive
- 配置超时策略

### 8.2 序列化优化
- 使用 protobuf 二进制格式
- 大数据分页传输
- 启用压缩（gzip/snappy）

## 9. 监控与治理

### 9.1 指标采集
```go
// trpc-go 内置指标
server.WithStat(
    stat.WithMetrics(metrics),
)
```

### 9.2 链路追踪
```go
server.WithFilter(
    opentracing.ServerFilter(),
)
```

## 10. 参考资源

- [trpc-go 官方文档](https://github.com/trpc-group/trpc-go)
- [Protocol Buffers 指南](https://developers.google.com/protocol-buffers)
- [API 设计最佳实践](https://cloud.google.com/apis/design)
