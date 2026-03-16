// API 服务接口定义
export interface IModelService {
  getModels(params?: ModelQueryParams): Promise<APIResponse<Model[]>>;
  getModelById(id: string): Promise<APIResponse<Model>>;
  getProviders(): Promise<APIResponse<Provider[]>>;
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
  provider?: string;
  capability?: string;
  search?: string;
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

// 模型相关类型
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

export interface Model {
  id: string;
  name: string;
  provider: string;
  providerId: string;
  logo: string;
  description: string;
  capabilities: Capability;
  contextWindow: number;
  maxTokens: number;
  pricing: Pricing;
  speed: Speed;
  quality: Quality;
  releaseDate: string;
  modelType: string;
  architecture: string;
  features: string[];
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
  type: 'text' | 'date' | 'boolean' | 'number' | 'currency' | 'percentage' | 'score';
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
export type CapabilityType = keyof Capability;

export interface CapabilityInfo {
  icon: string;
  name: string;
  desc: string;
}
