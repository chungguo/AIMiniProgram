import type { 
  Capability, 
  Model, 
  ComparisonItem,
  CapabilityInfo 
} from '@/types/api';

/**
 * 获取能力图标列表
 */
export function getCapabilityList(capabilities: Capability): CapabilityInfo[] {
  const list: CapabilityInfo[] = [];
  
  const capabilityMap: Record<keyof Capability, CapabilityInfo> = {
    text: { icon: '📝', name: '文本', desc: '支持文本输入输出' },
    image: { icon: '🖼️', name: '图片', desc: '支持图像理解' },
    audio: { icon: '🔊', name: '音频', desc: '支持语音处理' },
    video: { icon: '🎬', name: '视频', desc: '支持视频理解' },
    file: { icon: '📄', name: '文件', desc: '支持文件上传' }
  };

  (Object.keys(capabilities) as Array<keyof Capability>).forEach((key) => {
    if (capabilities[key]) {
      list.push(capabilityMap[key]);
    }
  });

  return list;
}

/**
 * 格式化数值
 */
export function formatNumber(value: number, unit?: string): string {
  if (value >= 1000000) {
    return `${(value / 1000000).toFixed(1)}M${unit ? ' ' + unit : ''}`;
  }
  if (value >= 1000) {
    return `${(value / 1000).toFixed(0)}K${unit ? ' ' + unit : ''}`;
  }
  return `${value}${unit ? ' ' + unit : ''}`;
}

/**
 * 从对象中获取嵌套值
 */
export function getNestedValue<T>(obj: Record<string, unknown>, path: string): T | undefined {
  const keys = path.split('.');
  let value: unknown = obj;
  
  for (const key of keys) {
    if (value === null || value === undefined) {
      return undefined;
    }
    value = (value as Record<string, unknown>)[key];
  }
  
  return value as T;
}

/**
 * 格式化对比值
 */
export function formatComparisonValue(value: unknown, type: string, unit?: string): string {
  if (value === undefined || value === null) return '-';
  
  switch (type) {
    case 'boolean':
      return value ? '✓' : '✗';
    case 'percentage':
      return `${value}%`;
    case 'currency':
      return `$${value}`;
    case 'number':
      return `${Number(value).toLocaleString()}${unit ? ' ' + unit : ''}`;
    case 'date':
      return String(value);
    default:
      return String(value);
  }
}

/**
 * 获取值样式类
 */
export function getValueClass(value: unknown, type: string): string {
  if (type === 'boolean') {
    return value ? 'value-yes' : 'value-no';
  }
  return '';
}

/**
 * 计算最优值索引
 */
export function calculateBestValues(
  models: Model[], 
  itemKey: string, 
  isHigherBetter: boolean = true
): (string | undefined)[] {
  const values = models.map((model) => ({
    value: getNestedValue<number>(model as Record<string, unknown>, itemKey),
    index: models.indexOf(model)
  }));

  const validValues = values.filter((v): v is { value: number; index: number } => 
    typeof v.value === 'number'
  );

  if (validValues.length === 0) {
    return new Array(models.length).fill(undefined);
  }

  const bestValue = isHigherBetter
    ? Math.max(...validValues.map((v) => v.value))
    : Math.min(...validValues.map((v) => v.value));

  const result: (string | undefined)[] = new Array(models.length).fill(undefined);
  
  validValues.forEach(({ value, index }) => {
    if (value === bestValue) {
      result[index] = 'best-value';
    }
  });

  return result;
}

/**
 * 判断数值类型是否越高越好
 */
export function isHigherBetterMetric(itemKey: string): boolean {
  const lowerBetterKeys = ['pricing.inputPrice', 'pricing.outputPrice', 'latency'];
  return !lowerBetterKeys.some((key) => itemKey.includes(key));
}
