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
