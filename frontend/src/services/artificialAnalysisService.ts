import type { 
  IHttpClient, 
  IArtificialAnalysisService, 
  APIResponse, 
  ArtificialAnalysis
} from '@/types/api';
import { httpClient } from './modelService';

// ArtificialAnalysisService 实现 - 依赖注入模式
class ArtificialAnalysisService implements IArtificialAnalysisService {
  private httpClient: IHttpClient;

  constructor(httpClient: IHttpClient) {
    this.httpClient = httpClient;
  }

  async getAll(): Promise<APIResponse<ArtificialAnalysis[]>> {
    return this.httpClient.get<APIResponse<ArtificialAnalysis[]>>('/analysis/artificialanalysis');
  }

  async getBySlug(slug: string): Promise<APIResponse<ArtificialAnalysis>> {
    return this.httpClient.get<APIResponse<ArtificialAnalysis>>(`/analysis/artificialanalysis/${slug}`);
  }

  async getModelWithAnalysis(modelId: string): Promise<APIResponse<{ model: unknown; analysis: ArtificialAnalysis | null }>> {
    return this.httpClient.get<APIResponse<{ model: unknown; analysis: ArtificialAnalysis | null }>>(`/models/analysis/${modelId}`);
  }
}

// 工厂函数
export function createArtificialAnalysisService(httpClient?: IHttpClient): IArtificialAnalysisService {
  const client = httpClient || httpClient;
  return new ArtificialAnalysisService(client);
}

// 默认导出
export const artificialAnalysisService = createArtificialAnalysisService(httpClient);
