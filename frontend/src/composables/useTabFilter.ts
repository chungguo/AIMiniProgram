import { ref, computed, watch } from 'vue';

interface TabOption<T = string> {
  label: string;
  value: T;
  icon?: string;
}

interface TabFilterOptions<T> {
  options: TabOption<T>[];
  defaultValue?: T;
  onChange?: (value: T) => void;
}

interface TabFilterResult<T> {
  activeTab: ReturnType<typeof ref<T>>;
  options: TabOption<T>[];
  activeLabel: ReturnType<typeof computed<string>>;
  activeOption: ReturnType<typeof computed<TabOption<T> | undefined>>;
  onTabChange: (value: T) => void;
  setTab: (value: T) => void;
}

/**
 * Tab 筛选管理 Composable
 * 统一管理 Tab 切换状态，支持单数据源和多数据源场景
 * 
 * @example
 * ```ts
 * // 基础用法 - 单数据源 Tab 筛选
 * const { activeTab, onTabChange, options } = useTabFilter({
 *   options: [
 *     { label: '全部', value: 'all' },
 *     { label: '进行中', value: 'active' },
 *     { label: '已完成', value: 'completed' }
 *   ],
 *   defaultValue: 'all',
 *   onChange: (value) => {
 *     console.log('切换到:', value);
 *     // 触发数据刷新
 *     refreshData(value);
 *   }
 * });
 * 
 * // 模板中使用
 * <t-tabs :value="activeTab" @change="onTabChange">
 *   <t-tab-panel 
 *     v-for="tab in options" 
 *     :key="tab.value"
 *     :label="tab.label" 
 *     :value="tab.value"
 *   >
 *     <!-- 内容 -->
 *   </t-tab-panel>
 * </t-tabs>
 * ```
 * 
 * @example
 * ```ts
 * // 多数据源场景 - 每个 Tab 对应不同数据
 * const tabConfigs = {
 *   models: {
 *     fetcher: modelService.getModels,
 *     columns: ['名称', '价格', '厂商']
 *   },
 *   papers: {
 *     fetcher: paperService.getPapers,
 *     columns: ['标题', '作者', '日期']
 *   }
 * };
 * 
 * const { activeTab } = useTabFilter({
 *   options: [
 *     { label: '模型', value: 'models' },
 *     { label: '论文', value: 'papers' }
 *   ],
 *   onChange: (value) => {
 *     const config = tabConfigs[value];
 *     // 切换数据源
 *     loadData(config.fetcher);
 *   }
 * });
 * ```
 */
export function useTabFilter<T = string>(options: TabFilterOptions<T>): TabFilterResult<T> {
  const { options: tabOptions, defaultValue, onChange } = options;

  const activeTab = ref<T>(defaultValue ?? tabOptions[0]?.value) as ReturnType<typeof ref<T>>;

  // 当前激活的选项
  const activeOption = computed(() => {
    return tabOptions.find(opt => opt.value === activeTab.value);
  });

  // 当前激活的标签名
  const activeLabel = computed(() => {
    return activeOption.value?.label ?? '';
  });

  /**
   * Tab 切换处理
   */
  const onTabChange = (value: T): void => {
    activeTab.value = value;
    onChange?.(value);
  };

  /**
   * 程序化设置 Tab
   */
  const setTab = (value: T): void => {
    if (tabOptions.some(opt => opt.value === value)) {
      onTabChange(value);
    }
  };

  return {
    activeTab,
    options: tabOptions,
    activeLabel,
    activeOption,
    onTabChange,
    setTab
  };
}

/**
 * 带数据源的 Tab 筛选（高级用法）
 * 每个 Tab 对应独立的数据源和列表管理
 */
interface TabDataSource<T, R> {
  key: T;
  label: string;
  fetcher: () => Promise<R[]>;
}

interface TabDataManagerResult<T, R> {
  activeTab: ReturnType<typeof ref<T>>;
  activeData: ReturnType<typeof ref<R[]>>;
  loading: ReturnType<typeof ref<boolean>>;
  error: ReturnType<typeof ref<Error | null>>;
  onTabChange: (value: T) => void;
  refresh: () => Promise<void>;
}

export function useTabDataManager<T extends string, R>(
  dataSources: TabDataSource<T, R>[],
  defaultKey?: T
): TabDataManagerResult<T, R> {
  const activeTab = ref<T>(defaultKey ?? dataSources[0]?.key);
  const activeData = ref<R[]>([]);
  const loading = ref(false);
  const error = ref<Error | null>(null);

  // 缓存每个 Tab 的数据
  const cache = new Map<T, R[]>();

  /**
   * 加载指定 Tab 的数据
   */
  const loadTabData = async (key: T): Promise<void> => {
    const dataSource = dataSources.find(ds => ds.key === key);
    if (!dataSource) return;

    // 优先从缓存读取
    if (cache.has(key)) {
      activeData.value = cache.get(key)!;
      return;
    }

    try {
      loading.value = true;
      error.value = null;
      const data = await dataSource.fetcher();
      cache.set(key, data);
      activeData.value = data;
    } catch (err) {
      error.value = err as Error;
    } finally {
      loading.value = false;
    }
  };

  /**
   * Tab 切换
   */
  const onTabChange = (value: T): void => {
    activeTab.value = value;
    loadTabData(value);
  };

  /**
   * 刷新当前 Tab 数据
   */
  const refresh = async (): Promise<void> => {
    cache.delete(activeTab.value);
    await loadTabData(activeTab.value);
  };

  // 初始加载
  if (activeTab.value) {
    loadTabData(activeTab.value);
  }

  return {
    activeTab,
    activeData,
    loading,
    error,
    onTabChange,
    refresh
  };
}
