import { ref } from 'vue';
import { useLoading } from './useLoading';

interface DetailOptions<T> {
  fetcher: (id: string) => Promise<T>;
  onError?: (error: Error) => void;
  immediate?: boolean;
}

interface DetailResult<T> {
  data: ReturnType<typeof ref<T | null>>;
  loading: ReturnType<typeof ref<boolean>>;
  error: ReturnType<typeof ref<Error | null>>;
  loadData: (id: string) => Promise<void>;
  reload: () => Promise<void>;
}

/**
 * 详情页数据管理 Composable
 * 统一管理详情页的加载、错误处理和数据获取
 * 
 * @example
 * ```ts
 * // model-detail.vue
 * const {
 *   data: model,
 *   loading,
 *   error,
 *   loadData
 * } = useDetail<Model>({
 *   fetcher: async (id) => {
 *     const res = await modelService.getModelById(id);
 *     if (!res.success) throw new Error(res.message);
 *     return res.data;
 *   }
 * });
 * 
 * // 加载详情
 * const id = getPageId(); // 获取页面参数
 * loadData(id);
 * ```
 * 
 * @example
 * ```ts
 * // 带错误处理的用法
 * const { data, loading, loadData } = useDetail<Paper>({
 *   fetcher: paperService.getPaperById,
 *   onError: (error) => {
 *     uni.showToast({ title: error.message, icon: 'none' });
 *   }
 * });
 * ```
 */
export function useDetail<T>(options: DetailOptions<T>): DetailResult<T> {
  const { fetcher, onError, immediate = false } = options;
  const { loading: loadingState, withLoading } = useLoading();

  const data = ref<T | null>(null) as ReturnType<typeof ref<T | null>>;
  const loading = loadingState;
  const error = ref<Error | null>(null);
  const currentId = ref<string>('');

  /**
   * 加载详情数据
   */
  const loadData = async (id: string): Promise<void> => {
    if (!id) {
      error.value = new Error('ID 不能为空');
      return;
    }

    currentId.value = id;
    error.value = null;

    try {
      const result = await withLoading(async () => {
        return await fetcher(id);
      });
      if (result) {
        data.value = result;
      }
    } catch (err) {
      error.value = err as Error;
      onError?.(err as Error);
    }
  };

  /**
   * 重新加载当前数据
   */
  const reload = async (): Promise<void> => {
    if (currentId.value) {
      await loadData(currentId.value);
    }
  };

  return {
    data,
    loading,
    error,
    loadData,
    reload
  };
}

/**
 * 获取当前页面 URL 参数中的 ID
 * 适用于 uni-app 的页面参数获取
 * 
 * @example
 * ```ts
 * const id = getPageId();
 * if (id) {
 *   await loadData(id);
 * }
 * ```
 */
export function getPageId(paramName: string = 'id'): string | undefined {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1];
  
  // 兼容不同 uni-app 版本的参数获取方式
  const options = (currentPage as any).options || 
                  (currentPage as any).$page?.options || 
                  (currentPage as any).$route?.query;
  
  return options?.[paramName];
}

/**
 * 获取当前页面多个 URL 参数
 * 
 * @example
 * ```ts
 * const { id, type } = getPageParams(['id', 'type']);
 * ```
 */
export function getPageParams(paramNames: string[]): Record<string, string | undefined> {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1];
  
  const options = (currentPage as any).options || 
                  (currentPage as any).$page?.options || 
                  (currentPage as any).$route?.query;
  
  const result: Record<string, string | undefined> = {};
  paramNames.forEach(name => {
    result[name] = options?.[name];
  });
  
  return result;
}
