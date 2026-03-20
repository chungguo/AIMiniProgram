import type { 
  Model, 
  ComparisonItem,
  ModalityInfo,
  Modality
} from '@/types/api';

// 模态图标映射
const MODALITY_ICONS: Record<Modality, string> = {
  text: '📝',
  image: '🖼️',
  audio: '🔊',
  video: '🎬',
  file: '📄'
};

const MODALITY_NAMES: Record<Modality, string> = {
  text: '文本',
  image: '图片',
  audio: '音频',
  video: '视频',
  file: '文件'
};

/**
 * 获取模态图标列表
 */
export function getModalityList(model: Model): ModalityInfo[] {
  return model.modalitiesInput?.map((mod) => ({
    icon: MODALITY_ICONS[mod] || '❓',
    name: MODALITY_NAMES[mod] || mod,
    desc: `支持${MODALITY_NAMES[mod] || mod}输入`
  })) || [];
}

/**
 * 获取模型家族显示名
 */
export function getFamilyName(model: Model): string {
  return model.family || 'Unknown';
}

/**
 * 格式化价格
 */
export function formatPrice(value: number): string {
  if (value === 0) return '免费';
  if (value < 0.01) return '<$0.01';
  return `$${value.toFixed(2)}`;
}

/**
 * 格式化数值
 */
export function formatNumber(value: number, unit?: string): string {
  if (value === undefined || value === null) return '-';
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
    case 'array':
      if (Array.isArray(value)) {
        return value.join(', ');
      }
      return String(value);
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
  const lowerBetterKeys = [
    'costInput', 'costOutput', 'costReasoning', 
    'costCacheRead', 'costCacheWrite', 
    'costInputAudio', 'costOutputAudio'
  ];
  return !lowerBetterKeys.some((key) => itemKey.includes(key));
}

/**
 * 获取模型核心特性列表
 */
export function getModelFeatures(model: Model): { icon: string; label: string; value: string }[] {
  const features: { icon: string; label: string; value: string }[] = [];
  
  if (model.reasoning) {
    features.push({ icon: '🧠', label: '推理', value: '支持' });
  }
  if (model.toolCall) {
    features.push({ icon: '🔧', label: '工具', value: '支持' });
  }
  if (model.structuredOutput) {
    features.push({ icon: '📋', label: '结构化', value: '支持' });
  }
  if (model.openWeights) {
    features.push({ icon: '🔓', label: '开源', value: '是' });
  }
  
  return features;
}
