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
  type: string;
  unit?: string;
}

export interface ComparisonCategory {
  key: string;
  name: string;
  items: ComparisonItem[];
}

// 论文类型 - 严格匹配数据库表 arxiv_cs_ai 结构
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
  id: string;
  slug: string;
  model_creator: string | null;
  artificial_analysis_intelligence_index: number | null;
  artificial_analysis_coding_index: number | null;
  artificial_analysis_math_index: number | null;
  mmlu_pro: number | null;
  gpqa: number | null;
  hle: number | null;
  livecodebench: number | null;
  scicode: number | null;
  math_500: number | null;
  aime: number | null;
  aime_25: number | null;
  ifbench: number | null;
  lcr: number | null;
  terminalbench_hard: number | null;
  tau2: number | null;
  price_1m_blended_3_to_1: number | null;
  price_1m_input_tokens: number | null;
  price_1m_output_tokens: number | null;
  median_output_tokens_per_second: number | null;
  median_time_to_first_token_seconds: number | null;
  median_time_to_first_answer_token: number | null;
  created_at: string;
  updated_at: string;
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
