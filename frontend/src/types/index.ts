// 模态类型
export type Modality = 'text' | 'image' | 'audio' | 'video' | 'file';

// 模型类型 - 匹配数据库表结构
export interface Model {
  // 主键
  id: string;
  name: string;
  
  // 家族
  family: string;
  
  // 能力特性
  attachment: boolean;
  reasoning: boolean;
  toolCall: boolean;
  structuredOutput: boolean;
  temperature: boolean;
  
  // 时间信息
  knowledge: string;
  releaseDate: string;
  lastUpdated: string;
  
  // 模态支持
  modalitiesInput: Modality[];
  modalitiesOutput: Modality[];
  
  // 开源状态
  openWeights: boolean;
  
  // 定价
  costInput: number;
  costOutput: number;
  costReasoning: number;
  costCacheRead: number;
  costCacheWrite: number;
  costInputAudio: number;
  costOutputAudio: number;
  
  // 限制
  limitContext: number;
  limitInput: number;
  limitOutput: number;
  
  // 推理字段
  interleavedField: string;
  
  // 时间戳
  createdAt: string;
  updatedAt: string;
  
  // 扩展字段
  description?: string;
  features?: string[];
  logo?: string;
  architecture?: string;
  
  // 旧版兼容
  provider?: string;
  providerId?: string;
  contextWindow?: number;
  maxTokens?: number;
}

// 兼容旧代码
export interface Capability {
  text: boolean;
  image: boolean;
  audio: boolean;
  video: boolean;
  file: boolean;
}

export interface Pricing {
  inputPrice: number;
  outputPrice: number;
  currency: string;
  unit: string;
}

export interface Speed {
  latency: string;
  tokensPerSecond: number;
}

export interface Quality {
  mmlu?: number;
  humanEval?: number;
  mtBench?: number;
  math?: number;
}

export interface Family {
  id: string;
  name: string;
  logo?: string;
  website?: string;
  apiUrl?: string;
}

export interface Provider {
  id: string;
  name: string;
  logo: string;
  website: string;
  apiUrl: string;
}

export interface ComparisonItem {
  key: string;
  name: string;
  type: string;
  unit?: string;
}

export interface ComparisonCategory {
  key: string;
  name: string;
  items: ComparisonItem[];
}

export interface Paper {
  id: string;
  title: string;
  titleCN: string;
  abstract: string;
  abstractCN: string;
  authors: string[];
  institutions: string[];
  publishDate: string;
  arxivUrl: string;
  pdfUrl: string;
  categories: string[];
  keywords: string[];
  readTime: number;
  language: string;
}

export interface PaperCategory {
  id: string;
  name: string;
  nameEn: string;
}

export interface APIResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface Pagination {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T;
  pagination: Pagination;
}
