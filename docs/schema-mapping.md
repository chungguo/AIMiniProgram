# 数据库字段映射文档

## Model 表字段映射

| 数据库字段 | Go 字段 | TypeScript 字段 | 类型 | 说明 |
|-----------|---------|-----------------|------|------|
| `id` | `ID` | `id` | string | 主键 |
| `name` | `Name` | `name` | string | 模型名称 |
| `family` | `Family` | `family` | string | 模型家族 |
| `attachment` | `Attachment` | `attachment` | boolean | 支持附件 |
| `reasoning` | `Reasoning` | `reasoning` | boolean | 支持推理 |
| `tool_call` | `ToolCall` | `toolCall` | boolean | 支持工具调用 |
| `structured_output` | `StructuredOutput` | `structuredOutput` | boolean | 支持结构化输出 |
| `temperature` | `Temperature` | `temperature` | boolean | 支持温度调节 |
| `knowledge` | `Knowledge` | `knowledge` | string | 知识截止日期 |
| `release_date` | `ReleaseDate` | `releaseDate` | string | 发布日期 |
| `last_updated` | `LastUpdated` | `lastUpdated` | string | 最后更新 |
| `modalities_input` | `ModalitiesInput` | `modalitiesInput` | []Modality | 输入模态数组 |
| `modalities_output` | `ModalitiesOutput` | `modalitiesOutput` | []Modality | 输出模态数组 |
| `open_weights` | `OpenWeights` | `openWeights` | boolean | 开源权重 |
| `cost_input` | `CostInput` | `costInput` | number | 输入价格($/1M) |
| `cost_output` | `CostOutput` | `costOutput` | number | 输出价格($/1M) |
| `cost_reasoning` | `CostReasoning` | `costReasoning` | number | 推理价格($/1M) |
| `cost_cache_read` | `CostCacheRead` | `costCacheRead` | number | 缓存读取价格($/1M) |
| `cost_cache_write` | `CostCacheWrite` | `costCacheWrite` | number | 缓存写入价格($/1M) |
| `cost_input_audio` | `CostInputAudio` | `costInputAudio` | number | 音频输入价格($/1M) |
| `cost_output_audio` | `CostOutputAudio` | `costOutputAudio` | number | 音频输出价格($/1M) |
| `limit_context` | `LimitContext` | `limitContext` | number | 最大上下文窗口 |
| `limit_output` | `LimitOutput` | `limitOutput` | number | 最大输出tokens |
| `limit_input` | `LimitInput` | `limitInput` | number | 最大输入tokens |
| `interleaved_field` | `InterleavedField` | `interleavedField` | string | 推理内容字段名 |
| `created_at` | `CreatedAt` | `createdAt` | timestamp | 创建时间 |
| `updated_at` | `UpdatedAt` | `updatedAt` | timestamp | 更新时间 |

## Modality 枚举

```sql
CREATE TYPE modality AS ENUM ('text', 'image', 'audio', 'video', 'file');
```

对应 Go 类型：
```go
type Modality string
const (
    ModalityText  Modality = "text"
    ModalityImage Modality = "image"
    ModalityAudio Modality = "audio"
    ModalityVideo Modality = "video"
    ModalityFile  Modality = "file"
)
```

对应 TypeScript 类型：
```typescript
type Modality = 'text' | 'image' | 'audio' | 'video' | 'file';
```

## API 端点

| 端点 | 说明 |
|------|------|
| `GET /api/models` | 获取模型列表（支持筛选/排序/分页） |
| `GET /api/models/detail/:id` | 获取单个模型详情 |
| `GET /api/models/families` | 获取所有模型家族 |
| `GET /api/models/family/:family` | 获取指定家族的模型 |
| `POST /api/models/compare` | 对比多个模型 |
| `GET /api/models/meta/comparison-categories` | 获取对比类别 |

## 查询参数

- `family` - 按家族筛选
- `hasAttachment` - 是否支持附件 (true/false)
- `hasReasoning` - 是否支持推理 (true/false)
- `hasToolCall` - 是否支持工具调用 (true/false)
- `openWeights` - 是否开源 (true/false)
- `minContext` - 最小上下文窗口
- `maxCostInput` - 最大输入价格
- `sortBy` - 排序字段 (name/family/costInput/costOutput/limitContext/releaseDate)
- `sortOrder` - 排序方向 (asc/desc)
- `page` - 页码
- `limit` - 每页数量
- `search` - 搜索关键词
