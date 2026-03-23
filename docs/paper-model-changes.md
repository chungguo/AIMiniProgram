# Paper 数据模型修正记录

## 数据库真源 (Source of Truth)

表名: `public.arxiv_cs_ai`

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(50) | PRIMARY KEY |
| title | text | 英文标题 |
| author | text | 作者（逗号分隔的字符串）|
| abstract | text | 英文摘要 |
| update_at | timestamp | 更新时间（注意不是 updated_at）|
| created_at | timestamp | 创建时间 |
| submit_at | date | 提交日期 |
| title_cn | text | 中文标题（可为null）|
| abstract_cn | text | 中文摘要（可为null）|

## 修正范围

### 1. 前端类型定义

**文件:**
- `/frontend/src/types/api.ts`
- `/frontend/src/types/index.ts`

**变更:**
```typescript
// 修改前
interface Paper {
  id: string;
  title: string;
  titleCN: string;          // ❌ 驼峰命名
  abstract: string;
  abstractCN: string;       // ❌ 驼峰命名
  authors: string[];        // ❌ 数组类型
  institutions: string[];   // ❌ 字段不存在
  publishDate: string;      // ❌ 命名不一致
  arxivUrl: string;         // ❌ 字段不存在
  pdfUrl: string;           // ❌ 字段不存在
  categories: string[];     // ❌ 字段不存在
  keywords: string[];       // ❌ 字段不存在
  readTime: number;         // ❌ 字段不存在
  language: string;         // ❌ 字段不存在
}

// 修改后
interface Paper {
  id: string;
  title: string;
  author: string;           // ✅ 字符串（逗号分隔）
  abstract: string;
  title_cn: string | null;  // ✅ snake_case，可为null
  abstract_cn: string | null;
  submit_at: string;        // ✅ 与数据库一致
  created_at: string;
  update_at: string;        // ✅ 注意不是 updated_at
}
```

### 2. 前端页面组件

**文件:**
- `/frontend/src/pages/papers/papers.vue`
- `/frontend/src/pages/paper-detail/paper-detail.vue`

**主要变更:**
- `paper.titleCN` → `paper.title_cn || paper.title`
- `paper.abstractCN` → `paper.abstract_cn || paper.abstract`
- `paper.authors` (数组) → `paper.author.split(',')`
- `paper.publishDate` → `paper.submit_at`
- 移除了 `readTime`, `categories`, `keywords`, `institutions` 等字段
- arxiv/pdf 链接改为根据 id 动态构造

### 3. 后端模型定义

**文件:** `/backend/models/types.go`

**变更:**
```go
// 修改前
type Paper struct {
    ID           string   `json:"id"`
    Title        string   `json:"title"`
    TitleCN      string   `json:"titleCN"`
    Abstract     string   `json:"abstract"`
    AbstractCN   string   `json:"abstractCN"`
    Authors      []string `json:"authors"`
    Institutions []string `json:"institutions"`
    PublishDate  string   `json:"publishDate"`
    ArxivURL     string   `json:"arxivUrl"`
    PDFURL       string   `json:"pdfUrl"`
    Categories   []string `json:"categories"`
    Keywords     []string `json:"keywords"`
    ReadTime     int      `json:"readTime"`
    Language     string   `json:"language"`
}

// 修改后
type Paper struct {
    ID         string    `json:"id" db:"id"`
    Title      string    `json:"title" db:"title"`
    Author     string    `json:"author" db:"author"`
    Abstract   string    `json:"abstract" db:"abstract"`
    TitleCn    string    `json:"title_cn" db:"title_cn"`
    AbstractCn string    `json:"abstract_cn" db:"abstract_cn"`
    SubmitAt   string    `json:"submit_at" db:"submit_at"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdateAt   time.Time `json:"update_at" db:"update_at"`
}
```

### 4. 后端 Repository

**文件:** `/backend/repository/postgres_paper_repository.go`

**变更:**
- 表名: `paper` → `arxiv_cs_ai`
- 移除了 PostgreSQL 数组类型 (`pq.Array`) 的使用
- 移除了不存在的字段查询
- `GetCategories()` 返回固定分类（因为数据库表无 categories 字段）

## 注意事项

1. **字段名 `update_at`**: 数据库中不是常见的 `updated_at`，修正时特别注意
2. **作者字段**: 数据库中是逗号分隔的字符串，不是数组
3. **分类功能**: 数据库表没有 categories 字段，如需此功能需要扩展表结构
4. **链接构造**: arxiv 链接根据 id 动态构造（`https://arxiv.org/abs/{id}`）

### 5. 移除分类功能

由于数据库表 `arxiv_cs_ai` 没有 `categories` 字段，已移除以下分类相关代码：

**前端:**
- `papers.vue`: 移除分类筛选 UI 和 `loadCategories()`
- `types/api.ts`: 移除 `PaperCategory` 类型和 `getCategories()` 方法
- `types/index.ts`: 移除 `PaperCategory` 类型
- `paperService.ts`: 移除 `getCategories()` 方法
- `PaperQueryParams`: 移除 `category` 参数

**后端:**
- `models/types.go`: 移除 `PaperCategory` 结构体
- `repository/interfaces.go`: 移除 `GetByCategory()` 和 `GetCategories()` 接口方法
- `postgres_paper_repository.go`: 移除 `GetByCategory()` 和 `GetCategories()` 实现
- `handlers/papers.go`: 移除 `GetPaperCategories()` 处理器，简化 `GetPapers()`（移除 category 参数处理）
- `routers/router.go`: 移除 `/papers/categories` 路由

## 当前 API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/papers` | 获取论文列表（支持 search/page/limit）|
| GET | `/api/papers/detail/:id` | 获取单篇论文详情 |
| GET | `/api/papers/featured/latest` | 获取最新论文 |

## 后续建议

如需恢复分类、关键词等功能，建议：

1. 扩展数据库表结构，添加新字段
2. 或使用关联表存储多值属性
3. 或使用 JSONB 字段存储灵活结构

---

## 6. ArtificialAnalysis 评测数据模块

基于 `public.artificialanalysis` 表创建了新模块，通过 `slug` 与 `model.name` 关联。

### 数据库表结构

```sql
CREATE TABLE IF NOT EXISTS public.artificialanalysis (
  id uuid PRIMARY KEY,
  slug varchar(100) UNIQUE NOT NULL,  -- 关联 model.name
  model_creator varchar(100),
  artificial_analysis_intelligence_index numeric(5,2),
  artificial_analysis_coding_index numeric(5,2),
  artificial_analysis_math_index numeric(5,2),
  mmlu_pro, gpqa, hle, livecodebench, scicode, math_500, aime, aime_25,
  ifbench, lcr, terminalbench_hard, tau2 numeric(5,3),
  price_1m_blended_3_to_1, price_1m_input_tokens, price_1m_output_tokens numeric(10,6),
  median_output_tokens_per_second, median_time_to_first_token_seconds, 
  median_time_to_first_answer_token numeric(10,3),
  created_at, updated_at timestamp
);
```

### 新增 API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/analysis/artificialanalysis` | 获取所有评测数据 |
| GET | `/api/analysis/artificialanalysis/:slug` | 根据 slug 获取评测数据 |
| GET | `/api/models/analysis/:id` | 获取模型及其评测数据 |

### 新增文件

**后端:**
- `models/types.go`: 添加 `ArtificialAnalysis` 和 `ModelWithAnalysis` 结构体
- `repository/interfaces.go`: 添加 `ArtificialAnalysisRepository` 接口
- `repository/postgres_artificial_analysis_repository.go`: Repository 实现
- `repository/factory.go`: 添加仓库创建和 getter
- `handlers/base.go`: 添加 `artificialAnalysisRepo` 变量
- `handlers/artificial_analysis.go`: 3 个 handler 实现
- `routers/router.go`: 添加 3 条路由

**前端:**
- `types/api.ts`: 添加 `ArtificialAnalysis` 接口和 `IArtificialAnalysisService`
- `types/index.ts`: 导出 `ArtificialAnalysis` 类型
- `services/artificialAnalysisService.ts`: Service 实现
- `services/index.ts`: 导出 service

### 关联方式

```
Model.name (varchar) <-- slug --> ArtificialAnalysis.slug (varchar)
    "gpt-4o"         <-- slug -->         "gpt-4o"
```

获取模型评测数据流程：
1. 获取模型信息: `GET /api/models/detail/:id`
2. 使用 `model.name` 作为 slug 查询评测: `GET /api/analysis/artificialanalysis/:slug`
3. 或者直接使用: `GET /api/models/analysis/:id`

---

## 7. API 响应格式简化

移除了嵌套的 `success`/`data`/`message` 包装，使用 HTTP status code 表示状态。

### 变更前

```json
{
  "success": true,
  "data": { ... },
  "message": ""
}
```

错误时：
```json
{
  "success": false,
  "data": null,
  "message": "错误信息"
}
```

### 变更后

成功时（HTTP 200）：
```json
{ ... }  // 直接返回数据
```

分页数据：
```json
{
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "totalPages": 10
  }
}
```

错误时（HTTP 4xx/5xx）：
- 响应头：`X-Error-Message: 错误信息`
- 响应体：`{}` 或空

### 修改的文件

**后端:**
- `handlers/base.go`: 添加 `ErrorResponse()` 函数，错误信息放入响应头
- `handlers/models.go`: 移除所有 `success`/`data` 包装
- `handlers/papers.go`: 移除所有 `success`/`data` 包装
- `handlers/artificial_analysis.go`: 移除所有 `success`/`data` 包装

**前端:**
- `types/api.ts`: 
  - 移除 `APIResponse<T>` 类型
  - `PaginatedResponse` → `PaginatedData`，移除 `success` 字段
  - 更新所有服务接口返回类型
- `types/index.ts`: 同步更新类型定义
- `services/modelService.ts`: 
  - 更新返回类型
  - 错误处理读取 `X-Error-Message` 响应头
- `services/paperService.ts`: 更新返回类型
- `services/artificialAnalysisService.ts`: 更新返回类型
- `pages/papers/papers.vue`: 移除 `res.success` 判断
- `pages/paper-detail/paper-detail.vue`: 移除 `res.success` 判断
