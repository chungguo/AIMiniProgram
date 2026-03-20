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

// 模型类型 - 匹配数据库表结构
export interface Model {
  id: string;
  name: string;
  family: string;
  attachment: boolean;
  reasoning: boolean;
  toolCall: boolean;
  structuredOutput: boolean;
  temperature: boolean;
  knowledge: string;
  releaseDate: string;
  lastUpdated: string;
  modalitiesInput: Modality[];
  modalitiesOutput: Modality[];
  openWeights: boolean;
  costInput: number;
  costOutput: number;
  costReasoning: number;
  costCacheRead: number;
  costCacheWrite: number;
  costInputAudio: number;
  costOutputAudio: number;
  limitContext: number;
  limitInput: number;
  limitOutput: number;
  interleavedField: string;
  createdAt: string;
  updatedAt: string;
}

// Capability 辅助结构
export interface Capability {
  text: boolean;
  image: boolean;
  audio: boolean;
  video: boolean;
  file: boolean;
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
export interface ModalityInfo {
  icon: string;
  name: string;
  desc: string;
}
