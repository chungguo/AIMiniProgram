// API 服务接口定义
export interface IModelService {
  getModels(params?: ModelQueryParams): Promise<PaginatedResponse<Model[]>>;
  getModelById(id: string): Promise<APIResponse<Model>>;
  getFamilies(): Promise<APIResponse<string[]>>;
  getFamilyModels(family: string): Promise<APIResponse<Model[]>>;
  compareModels(ids: string[]): Promise<APIResponse<CompareResult>>;
  getComparisonCategories(): Promise<APIResponse<ComparisonCategory[]>>;
}

export interface IPaperService {
  getPapers(params?: PaperQueryParams): Promise<PaginatedResponse<Paper[]>>;
  getPaperById(id: string): Promise<APIResponse<Paper>>;
  getCategories(): Promise<APIResponse<PaperCategory[]>>;
}

export interface IHttpClient {
  get<T>(url: string, params?: Record<string, unknown>): Promise<T>;
  post<T>(url: string, data?: unknown): Promise<T>;
}

// 查询参数类型
export interface ModelQueryParams {
  family?: string;
  hasAttachment?: boolean;
  hasReasoning?: boolean;
  hasToolCall?: boolean;
  openWeights?: boolean;
  minContext?: number;
  maxCostInput?: number;
  sortBy?: 'name' | 'family' | 'costInput' | 'costOutput' | 'limitContext' | 'releaseDate';
  sortOrder?: 'asc' | 'desc';
  search?: string;
  page?: number;
  limit?: number;
}

export interface PaperQueryParams {
  category?: string;
  search?: string;
  page?: number;
  limit?: number;
}

// 响应类型
export interface APIResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  code?: string;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface CompareResult {
  models: Model[];
  comparisonCategories: ComparisonCategory[];
}

// 模态类型
export type Modality = 'text' | 'image' | 'audio' | 'video' | 'file';

// 模型相关类型 - 匹配数据库表结构
export interface Model {
  // 主键
  id: string;
  name: string;
  
  // 家族和分类
  family: string;        // 模型家族（替代 provider）
  
  // 能力特性（布尔标志）
  attachment: boolean;        // 支持附件
  reasoning: boolean;         // 支持推理
  toolCall: boolean;          // 支持工具调用
  structuredOutput: boolean;  // 支持结构化输出
  temperature: boolean;       // 支持温度调节
  
  // 时间信息
  knowledge: string;      // 知识截止日期
  releaseDate: string;    // 发布日期
  lastUpdated: string;    // 最后更新日期
  
  // 模态支持
  modalitiesInput: Modality[];   // 输入模态
  modalitiesOutput: Modality[];  // 输出模态
  
  // 开源状态
  openWeights: boolean;   // 权重是否开源
  
  // 定价（每百万 tokens，USD）
  costInput: number;        // 输入价格
  costOutput: number;       // 输出价格
  costReasoning: number;    // 推理价格
  costCacheRead: number;    // 缓存读取价格
  costCacheWrite: number;   // 缓存写入价格
  costInputAudio: number;   // 音频输入价格
  costOutputAudio: number;  // 音频输出价格
  
  // 限制
  limitContext: number;  // 最大上下文窗口
  limitInput: number;    // 最大输入 tokens
  limitOutput: number;   // 最大输出 tokens
  
  // 推理内容字段名
  interleavedField: string;  // "reasoning_content" 或 "reasoning_details"
  
  // 时间戳
  createdAt: string;
  updatedAt: string;
  
  // 扩展字段（JSON 兼容）
  description?: string;
  features?: string[];
  logo?: string;
  architecture?: string;
  
  // 旧版兼容字段
  provider?: string;
  providerId?: string;
  contextWindow?: number;
  maxTokens?: number;
}

// 兼容旧代码的 Capability 结构
export interface Capability {
  text: boolean;
  image: boolean;
  audio: boolean;
  video: boolean;
  file: boolean;
}

// 兼容旧代码的 Pricing 结构
export interface Pricing {
  inputPrice: number;
  outputPrice: number;
  currency: string;
  unit: string;
}

// 家族（替代 Provider）
export interface Family {
  id: string;
  name: string;
  logo?: string;
  website?: string;
  apiUrl?: string;
}

// 兼容旧代码的 Provider
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
  type: 'text' | 'date' | 'boolean' | 'number' | 'currency' | 'percentage' | 'score' | 'array';
  unit?: string;
}

export interface ComparisonCategory {
  key: string;
  name: string;
  items: ComparisonItem[];
}

// 论文相关类型
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

// 通用工具类型
export type ModalityType = Modality;

export interface ModalityInfo {
  icon: string;
  name: string;
  desc: string;
}
