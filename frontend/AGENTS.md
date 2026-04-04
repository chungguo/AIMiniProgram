# Frontend

uni-app 微信小程序前端。

## 技术栈

- Vue 3 + TypeScript
- uni-app (微信小程序/H5/App)
- UnoCSS (原子化 CSS)
- TDesign UniApp (UI 组件库)

## 结构

```
src/
├── composables/  # Vue 3 Composition API 逻辑复用
│   ├── useLoading.ts      # 加载状态管理
│   ├── useListManager.ts  # 列表数据管理
│   └── useTabFilter.ts    # Tab 筛选管理
├── pages/        # 页面
├── services/     # API 服务
├── types/        # TypeScript 类型
└── utils/        # 工具函数
```

## 运行

```bash
npm run dev:h5        # H5 开发
npm run dev:mp-weixin # 微信小程序
```

## 服务层

- `modelService` - 模型 API
- `paperService` - 论文 API
- `artificialAnalysisService` - 评测数据 API
- `interceptorManager` - HTTP 拦截器
- `auth.ts` - Token 管理增强

## Composables 使用规范（必读）

项目已引入 MP 项目的优秀实践，统一使用 composables 管理状态和逻辑。

### 1. 加载状态 - useLoading

```typescript
import { useLoading } from '@/composables';

// ✅ 推荐：自动管理加载状态
const { withLoading } = useLoading();
await withLoading(async () => {
  const data = await fetchData();
});

// ❌ 避免：手动管理 loading
const loading = ref(false);
loading.value = true;
try {
  await fetchData();
} finally {
  loading.value = false;
}
```

### 2. 列表场景 - useListManager

```typescript
import { useListManager } from '@/composables';

// ✅ 推荐：使用 useListManager
const { list, loading, refresh, loadMore, finished } = useListManager<Item>({
  fetcher: async (params) => {
    const res = await api.getList({ page: params.page, limit: params.pageSize });
    return { data: res.data, pagination: res.pagination };
  },
  pageSize: 10
});

// ❌ 避免：手动管理分页、刷新、加载更多
```

### 3. Tab 筛选 - useTabFilter

```typescript
import { useTabFilter } from '@/composables';

// ✅ 推荐：使用 useTabFilter
const { activeTab, onTabChange, options } = useTabFilter({
  options: [
    { label: '全部', value: 'all' },
    { label: '进行中', value: 'active' }
  ],
  onChange: (value) => refreshData(value)
});
```

### 4. 代码组织原则

- **保持业务代码稳定**：新增功能优先使用 composables，不动现有页面逻辑
- **提取重复逻辑**：发现重复代码时，提取为 composable
- **类型安全**：所有 composables 都支持泛型，确保类型推断

## 迁移状态

| 页面 | 状态 | 使用的 Composables |
|------|------|-------------------|
| index.vue | ✅ 已迁移 | useLoading |
| models.vue | ✅ 已迁移 | useLoading, useSimpleList |
| papers.vue | ✅ 已迁移 | useListManager |
| model-detail.vue | ⏳ 待迁移 | - |
| paper-detail.vue | ⏳ 待迁移 | - |
| compare.vue | ⏳ 待迁移 | - |

## 参考文档

- [Composables 使用指南](./src/composables/README.md)
