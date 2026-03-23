import type { 
  IHttpClient, 
  IPaperService, 
  APIResponse, 
  Paper, 
  PaperQueryParams,
  PaginatedResponse
} from '@/types/api';
import { httpClient } from './modelService';

// PaperService 实现 - 依赖注入模式
class PaperService implements IPaperService {
  private httpClient: IHttpClient;

  constructor(httpClient: IHttpClient) {
    this.httpClient = httpClient;
  }

  async getPapers(params?: PaperQueryParams): Promise<PaginatedResponse<Paper[]>> {
    const defaultParams: PaperQueryParams = {
      page: 1,
      limit: 10,
      ...params
    };
    return this.httpClient.get<PaginatedResponse<Paper[]>>('/papers', defaultParams as Record<string, unknown>);
  }

  async getPaperById(id: string): Promise<APIResponse<Paper>> {
    return this.httpClient.get<APIResponse<Paper>>(`/papers/detail/${id}`);
  }
}

// 工厂函数
export function createPaperService(httpClient?: IHttpClient): IPaperService {
  const client = httpClient || httpClient;
  return new PaperService(client);
}

// 默认导出
export const paperService = createPaperService(httpClient);
