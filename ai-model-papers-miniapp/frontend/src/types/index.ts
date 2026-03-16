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
