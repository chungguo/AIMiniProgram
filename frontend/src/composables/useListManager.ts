import { ref, computed } from 'vue';
import { useListLoading } from './useLoading';

interface PaginationParams {
  page: number;
  limit: number;
}

interface PaginationInfo {
  current: number;
  pageSize: number;
  total: number;
}

interface ListManagerOptions<T> {
  fetcher: (params: PaginationParams) => Promise<{
    data: T[];
    pagination: PaginationInfo;
  }>;
  pageSize?: number;
  immediate?: boolean;
  onError?: (error: Error) => void;
}

interface ListManagerResult<T> {
  list: ReturnType<typeof ref<T[]>>;
  loading: ReturnType<typeof ref<boolean>>;
  refreshing: ReturnType<typeof ref<boolean>>;
  loadingMore: ReturnType<typeof ref<boolean>>;
  finished: ReturnType<typeof ref<boolean>>;
  error: ReturnType<typeof ref<Error | null>>;
  currentPage: ReturnType<typeof ref<number>>;
  total: ReturnType<typeof ref<number>>;
  loadData: (page?: number) => Promise<void>;
  refresh: () => Promise<void>;
  loadMore: () => Promise<void>;
  reset: () => void;
}

/**
 * 列表数据管理 Composable
 * 统一处理列表分页、下拉刷新、触底加载等场景
 * 
 * @example
 * ```ts
 * // 基础用法
 * const { list, loading, finished, refresh, loadMore } = useListManager({
 *   fetcher: async (params) => {
 *     const res = await paperService.getPapers({
 *       page: params.page,
 *       limit: params.pageSize
 *     });
 *     return {
 *       data: res.data,
 *       pagination: {
 *         current: res.pagination.page,
 *         pageSize: res.pagination.limit,
 *         total: res.pagination.total
 *       }
 *     };
 *   },
 *   pageSize: 10
 * });
 * 
 * // 在模板中使用
 * <scroll-view
 *   scroll-y
 *   @scrolltolower="loadMore"
 *   refresher-enabled
 *   :refresher-triggered="refreshing"
 *   @refresherrefresh="refresh"
 * >
 *   <view v-for="item in list" :key="item.id">...</view>
 *   <view v-if="finished">没有更多了</view>
 * </scroll-view>
 * ```
 */
export function useListManager<T>(options: ListManagerOptions<T>): ListManagerResult<T> {
  const { fetcher, pageSize = 10, immediate = true, onError } = options;

  const list = ref<T[]>([]) as ReturnType<typeof ref<T[]>>;
  const currentPage = ref(1);
  const total = ref(0);
  const error = ref<Error | null>(null);
  
  const { loading, refreshing, loadingMore, withRefresh, withLoadMore } = useListLoading();

  // 是否已加载完所有数据
  const finished = computed(() => {
    if (total.value === 0) return false;
    return list.value.length >= total.value;
  });

  /**
   * 加载数据
   */
  const loadData = async (page: number = 1): Promise<void> => {
    try {
      error.value = null;
      const result = await fetcher({
        page,
        limit: pageSize
      });

      if (page === 1) {
        list.value = result.data;
      } else {
        list.value.push(...result.data);
      }
      
      total.value = result.pagination.total;
      currentPage.value = page;
    } catch (err) {
      error.value = err as Error;
      onError?.(err as Error);
      throw err;
    }
  };

  /**
   * 刷新列表（回到第一页）
   */
  const refresh = async (): Promise<void> => {
    await withRefresh(async () => {
      await loadData(1);
    });
  };

  /**
   * 加载更多（下一页）
   */
  const loadMore = async (): Promise<void> => {
    if (finished.value || loadingMore.value) return;
    
    await withLoadMore(async () => {
      await loadData(currentPage.value + 1);
    });
  };

  /**
   * 重置列表状态
   */
  const reset = (): void => {
    list.value = [];
    currentPage.value = 1;
    total.value = 0;
    error.value = null;
  };

  // 立即加载
  if (immediate) {
    loadData(1);
  }

  return {
    list,
    loading,
    refreshing,
    loadingMore,
    finished,
    error,
    currentPage,
    total,
    loadData,
    refresh,
    loadMore,
    reset
  };
}

/**
 * 简单列表管理（无分页）
 */
export function useSimpleList<T>(fetcher: () => Promise<T[]>) {
  const list = ref<T[]>([]);
  const { loading, withLoading } = useListLoading();
  const error = ref<Error | null>(null);

  const loadData = async (): Promise<void> => {
    try {
      error.value = null;
      const result = await withLoading(fetcher);
      if (result) {
        list.value = result;
      }
    } catch (err) {
      error.value = err as Error;
      throw err;
    }
  };

  return {
    list,
    loading,
    error,
    loadData
  };
}
