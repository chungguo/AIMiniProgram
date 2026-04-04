# Composables 使用指南

本项目从 MP 小程序项目引入了以下前端优秀实践，统一使用 Vue 3 Composition API 风格的 composables 来管理状态和逻辑。

## 📦 已提供的 Composables

### 1. useLoading - 加载状态管理

统一管理页面加载状态，支持普通加载、下拉刷新、触底加载等场景。

```typescript
import { useLoading, useListLoading } from '@/composables';

// ========== 基础用法 ==========
const { loading, showLoading, hideLoading, withLoading } = useLoading();

// 自动管理加载状态
await withLoading(async () => {
  const data = await fetchData();
  list.value = data;
});

// 手动控制
showLoading('加载中...');
await doSomething();
hideLoading();

// ========== 列表场景用法 ==========
const { 
  loading, 
  refreshing, 
  loadingMore,
  withRefresh, 
  withLoadMore 
} = useListLoading();

// 下拉刷新
await withRefresh(async () => {
  await refreshData();
});

// 加载更多
await withLoadMore(async () => {
  await loadMoreData();
});
```

### 2. useListManager - 列表数据管理

专门为列表场景设计，内置分页、下拉刷新、触底加载、空状态等完整功能。

```typescript
import { useListManager } from '@/composables';
import type { Paper } from '@/types/api';

const {
  list,           // 列表数据
  loading,        // 加载状态
  refreshing,     // 下拉刷新状态
  loadingMore,    // 加载更多状态
  finished,       // 是否已加载完
  error,          // 错误信息
  currentPage,    // 当前页码
  total,          // 总数量
  refresh,        // 刷新方法
  loadMore,       // 加载更多方法
  reset           // 重置方法
} = useListManager<Paper>({
  fetcher: async (params) => {
    // 调用 API
    const res = await paperService.getPapers({
      page: params.page,
      limit: params.pageSize
    });
    
    // 返回标准格式
    return {
      data: res.data,
      pagination: {
        current: res.pagination.page,
        pageSize: res.pagination.limit,
        total: res.pagination.total
      }
    };
  },
  pageSize: 10,      // 每页数量
  immediate: true    // 是否立即加载
});
```

**模板中使用：**

```vue
<template>
  <scroll-view
    scroll-y
    @scrolltolower="loadMore"
    refresher-enabled
    :refresher-triggered="refreshing"
    @refresherrefresh="refresh"
  >
    <!-- 列表内容 -->
    <view v-for="item in list" :key="item.id">
      {{ item.name }}
    </view>
    
    <!-- 加载状态 -->
    <view v-if="loadingMore">加载中...</view>
    <view v-else-if="finished">没有更多了</view>
  </scroll-view>
</template>
```

### 3. useSimpleList - 简单列表（无分页）

适用于数据量较小、无需分页的场景。

```typescript
import { useSimpleList } from '@/composables';

const { list, loading, loadData } = useSimpleList<Model>(
  async () => {
    const res = await modelService.getModels();
    return res.data;
  }
);

// 加载数据
await loadData();
```

### 4. useTabFilter - Tab 筛选管理

统一管理 Tab 切换状态，支持回调和数据源切换。

```typescript
import { useTabFilter, useTabDataManager } from '@/composables';

// ========== 基础用法 ==========
const { 
  activeTab,      // 当前激活的 tab
  options,        // tab 选项配置
  activeLabel,    // 当前激活的标签名
  onTabChange     // 切换回调
} = useTabFilter({
  options: [
    { label: '全部', value: 'all' },
    { label: '进行中', value: 'active' },
    { label: '已完成', value: 'completed' }
  ],
  defaultValue: 'all',
  onChange: (value) => {
    console.log('切换到:', value);
    // 触发数据刷新
    refreshData(value);
  }
});
```

**模板中使用：**

```vue
<template>
  <t-tabs :value="activeTab" @change="onTabChange">
    <t-tab-panel 
      v-for="tab in options" 
      :key="tab.value"
      :label="tab.label" 
      :value="tab.value"
    >
      <!-- 内容 -->
    </t-tab-panel>
  </t-tabs>
</template>
```

## 🎯 代码改进示例

### 改进前（papers.vue）

```typescript
const papers = ref<Paper[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const hasMore = ref(true);

async function loadPapers(refresh = false) {
  if (loading.value) return;
  
  if (refresh) {
    currentPage.value = 1;
    hasMore.value = true;
    papers.value = [];
  }
  
  try {
    loading.value = true;
    const res = await paperService.getPapers({
      page: currentPage.value,
      limit: PAGE_SIZE
    });
    
    if (refresh) papers.value = res.data;
    else papers.value.push(...res.data);
    hasMore.value = res.data.length === PAGE_SIZE;
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
}
```

### 改进后（papers.vue）

```typescript
const {
  list: papers,
  refreshing,
  loadingMore,
  finished,
  refresh,
  loadMore
} = useListManager<Paper>({
  fetcher: async (params) => {
    const res = await paperService.getPapers({
      page: params.page,
      limit: params.pageSize
    });
    return {
      data: res.data,
      pagination: {
        current: res.pagination.page,
        pageSize: res.pagination.limit,
        total: res.pagination.total
      }
    };
  },
  pageSize: 10
});
```

**代码量减少约 60%，逻辑更清晰！**

## 🔐 Auth 增强

新增了 `services/auth.ts`，提供更完善的 Token 管理：

```typescript
import { 
  getCurrentToken,    // 获取当前 token
  setToken,           // 设置 token
  clearToken,         // 清除 token
  getValidToken,      // 获取有效 token（自动刷新）
  setupEnhancedInterceptors,  // 设置增强拦截器
  retryRequest        // 重试失败的请求
} from '@/services';

// 使用增强拦截器
setupEnhancedInterceptors();

// 重试请求
const data = await retryRequest(async () => {
  return await httpClient.get('/api/data');
}, 3); // 最多重试 3 次
```

## 📁 文件结构

```
src/
├── composables/
│   ├── index.ts           # 统一导出
│   ├── useLoading.ts      # 加载状态管理
│   ├── useListManager.ts  # 列表数据管理
│   └── useTabFilter.ts    # Tab 筛选管理
├── services/
│   ├── auth.ts            # Token 管理增强
│   ├── interceptor.ts     # 拦截器
│   └── ...
└── ...
```

## 💡 最佳实践

1. **列表场景**：优先使用 `useListManager`，内置分页、刷新、加载更多
2. **简单加载**：使用 `useLoading` 或 `useListLoading` 统一管理状态
3. **Tab 切换**：使用 `useTabFilter` 统一管理状态，避免重复代码
4. **错误处理**：在 `fetcher` 中统一处理，或通过 `onError` 回调
5. **类型安全**：所有 composables 都支持泛型，确保类型推断正确

## 📝 迁移建议

现有页面可以逐步迁移：

1. 新页面直接使用 composables
2. 旧页面在重构时替换
3. 保持业务代码不动，提取重复逻辑到 composables
