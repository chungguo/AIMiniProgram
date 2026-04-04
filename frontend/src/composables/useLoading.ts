import { ref } from 'vue';

interface LoadingOptions {
  title?: string;
  mask?: boolean;
}

/**
 * 加载状态管理 Composable
 * 统一处理页面加载状态，支持普通加载、带文字加载、全屏加载
 * 
 * @example
 * ```ts
 * // 基础用法
 * const { loading, showLoading, hideLoading, withLoading } = useLoading();
 * 
 * // 自动管理加载状态
 * await withLoading(async () => {
 *   const data = await fetchData();
 *   return data;
 * });
 * 
 * // 手动控制
 * showLoading('加载中...');
 * await doSomething();
 * hideLoading();
 * ```
 */
export function useLoading() {
  const loading = ref(false);
  const loadingText = ref('加载中...');

  /**
   * 显示加载
   */
  const showLoading = (options: LoadingOptions = {}) => {
    const { title = '加载中...', mask = true } = options;
    loading.value = true;
    loadingText.value = title;
    uni.showLoading({ title, mask });
  };

  /**
   * 隐藏加载
   */
  const hideLoading = () => {
    loading.value = false;
    uni.hideLoading();
  };

  /**
   * 执行异步操作并自动管理加载状态
   */
  const withLoading = async <T>(
    fn: () => Promise<T>,
    options: LoadingOptions = {}
  ): Promise<T | undefined> => {
    try {
      showLoading(options);
      return await fn();
    } finally {
      hideLoading();
    }
  };

  return {
    loading,
    loadingText,
    showLoading,
    hideLoading,
    withLoading
  };
}

/**
 * 列表加载状态管理
 * 专门用于列表场景，支持下拉刷新和触底加载
 */
export function useListLoading() {
  const { loading, showLoading, hideLoading, withLoading } = useLoading();
  const refreshing = ref(false);
  const loadingMore = ref(false);

  /**
   * 执行下拉刷新
   */
  const withRefresh = async <T>(fn: () => Promise<T>): Promise<T | undefined> => {
    try {
      refreshing.value = true;
      return await fn();
    } finally {
      refreshing.value = false;
      uni.stopPullDownRefresh();
    }
  };

  /**
   * 执行加载更多
   */
  const withLoadMore = async <T>(fn: () => Promise<T>): Promise<T | undefined> => {
    if (loadingMore.value) return;
    try {
      loadingMore.value = true;
      return await fn();
    } finally {
      loadingMore.value = false;
    }
  };

  return {
    loading,
    refreshing,
    loadingMore,
    showLoading,
    hideLoading,
    withLoading,
    withRefresh,
    withLoadMore
  };
}
