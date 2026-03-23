// API 服务接口定义 - 简化版：直接使用数据类型，不再包装
export interface IModelService {
  getModels(params?: ModelQueryParams): Promise<PaginatedData<Model[]>>;
  getModelById(id: string): Promise<Model>;
  getFamilies(): Promise<string[]>;
  getFamilyModels(family: string): Promise<Model[]>;
  compareModels(ids: string[]): Promise<CompareResult>;
  getComparisonCategories(): Promise<ComparisonCategory[]>;
}

export interface IPaperService {
  getPapers(params?: PaperQueryParams): Promise<PaginatedData<Paper[]>>;
  getPaperById(id: string): Promise<Paper>;
}

export interface IArtificialAnalysisService {
  getAll(): Promise<ArtificialAnalysis[]>;
  getBySlug(slug: string): Promise<ArtificialAnalysis>;
  getModelWithAnalysis(modelId: string): Promise<{ model: unknown; analysis: ArtificialAnalysis | null }>;
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
  search?: string;
  page?: number;
  limit?: number;
}

// 分页数据结构 - 不再包装 success
export interface PaginatedData<T> {
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

// 论文相关类型 - 严格匹配数据库表 arxiv_cs_ai 结构
export interface Paper {
  id: string;                    // varchar(50) PRIMARY KEY
  title: string;                 // text
  author: string;                // text (逗号分隔的作者字符串)
  abstract: string;              // text
  title_cn: string | null;       // text (数据库 snake_case)
  abstract_cn: string | null;    // text (数据库 snake_case)
  submit_at: string;             // date -> ISO 日期字符串 (YYYY-MM-DD)
  created_at: string;            // timestamp -> ISO 字符串
  update_at: string;             // timestamp -> ISO 字符串 (注意: 数据库字段名是 update_at)
}

// ArtificialAnalysis 评测数据类型 - 匹配 PostgreSQL 表结构
export interface ArtificialAnalysis {
  id: string;                                        // uuid
  slug: string;                                      // varchar(100) UNIQUE, 关联 model.name
  model_creator: string | null;                      // varchar(100)
  artificial_analysis_intelligence_index: number | null;  // numeric(5,2)
  artificial_analysis_coding_index: number | null;        // numeric(5,2)
  artificial_analysis_math_index: number | null;          // numeric(5,2)
  mmlu_pro: number | null;                           // numeric(5,3)
  gpqa: number | null;                               // numeric(5,3)
  hle: number | null;                                // numeric(5,3)
  livecodebench: number | null;                      // numeric(5,3)
  scicode: number | null;                            // numeric(5,3)
  math_500: number | null;                           // numeric(5,3)
  aime: number | null;                               // numeric(5,3)
  aime_25: number | null;                            // numeric(5,3)
  ifbench: number | null;                            // numeric(5,3)
  lcr: number | null;                                // numeric(5,3)
  terminalbench_hard: number | null;                 // numeric(5,3)
  tau2: number | null;                               // numeric(5,3)
  price_1m_blended_3_to_1: number | null;            // numeric(10,6)
  price_1m_input_tokens: number | null;              // numeric(10,6)
  price_1m_output_tokens: number | null;             // numeric(10,6)
  median_output_tokens_per_second: number | null;    // numeric(10,3)
  median_time_to_first_token_seconds: number | null; // numeric(10,3)
  median_time_to_first_answer_token: number | null;  // numeric(10,3)
  created_at: string;                                // timestamp
  updated_at: string;                                // timestamp
}

// 通用工具类型
export interface ModalityInfo {
  icon: string;
  name: string;
  desc: string;
}
