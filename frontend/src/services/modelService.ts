import type { 
  IHttpClient, 
  IModelService, 
  APIResponse, 
  Model, 
  Provider, 
  ModelQueryParams,
  CompareResult,
  ComparisonCategory,
  PaginatedResponse
} from '@/types/api';

// HttpClient 实现 - 单例模式
class UniHttpClient implements IHttpClient {
  private static instance: UniHttpClient;
  private baseURL: string;

  private constructor(baseURL: string = 'http://localhost:3000/api') {
    this.baseURL = baseURL;
  }

  static getInstance(baseURL?: string): UniHttpClient {
    if (!UniHttpClient.instance) {
      UniHttpClient.instance = new UniHttpClient(baseURL);
    }
    return UniHttpClient.instance;
  }

  async get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
    return new Promise((resolve, reject) => {
      const fullURL = this.buildURL(url, params);
      
      uni.request({
        url: fullURL,
        method: 'GET',
        success: (res) => {
          if (res.statusCode >= 200 && res.statusCode < 300) {
            resolve(res.data as T);
          } else {
            reject(new Error(`HTTP ${res.statusCode}: ${res.errMsg}`));
          }
        },
        fail: (err) => {
          reject(new Error(err.errMsg || 'Request failed'));
        }
      });
    });
  }

  async post<T>(url: string, data?: unknown): Promise<T> {
    return new Promise((resolve, reject) => {
      uni.request({
        url: `${this.baseURL}${url}`,
        method: 'POST',
        data,
        header: {
          'Content-Type': 'application/json'
        },
        success: (res) => {
          if (res.statusCode >= 200 && res.statusCode < 300) {
            resolve(res.data as T);
          } else {
            reject(new Error(`HTTP ${res.statusCode}: ${res.errMsg}`));
          }
        },
        fail: (err) => {
          reject(new Error(err.errMsg || 'Request failed'));
        }
      });
    });
  }

  private buildURL(url: string, params?: Record<string, unknown>): string {
    const fullURL = `${this.baseURL}${url}`;
    if (!params) return fullURL;

    const queryParams = Object.entries(params)
      .filter(([, value]) => value !== undefined && value !== null)
      .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
      .join('&');

    return queryParams ? `${fullURL}?${queryParams}` : fullURL;
  }
}

// ModelService 实现
class ModelService implements IModelService {
  private httpClient: IHttpClient;

  constructor(httpClient: IHttpClient) {
    this.httpClient = httpClient;
  }

  async getModels(params?: ModelQueryParams): Promise<PaginatedResponse<Model[]>> {
    return this.httpClient.get<PaginatedResponse<Model[]>>('/models', params as Record<string, unknown>);
  }

  async getModelById(id: string): Promise<APIResponse<Model>> {
    return this.httpClient.get<APIResponse<Model>>(`/models/detail/${id}`);
  }

  async getFamilies(): Promise<APIResponse<string[]>> {
    return this.httpClient.get<APIResponse<string[]>>('/models/families');
  }

  async getFamilyModels(family: string): Promise<APIResponse<Model[]>> {
    return this.httpClient.get<APIResponse<Model[]>>(`/models/family/${family}`);
  }

  // 兼容旧接口
  async getProviders(): Promise<APIResponse<Provider[]>> {
    return this.httpClient.get<APIResponse<Provider[]>>('/models/providers');
  }

  async compareModels(ids: string[]): Promise<APIResponse<CompareResult>> {
    return this.httpClient.post<APIResponse<CompareResult>>('/models/compare', { ids });
  }

  async getComparisonCategories(): Promise<APIResponse<ComparisonCategory[]>> {
    return this.httpClient.get<APIResponse<ComparisonCategory[]>>('/models/meta/comparison-categories');
  }
}

// 工厂函数
export function createModelService(httpClient?: IHttpClient): IModelService {
  const client = httpClient || UniHttpClient.getInstance();
  return new ModelService(client);
}

// 默认导出
export const httpClient = UniHttpClient.getInstance();
export const modelService = createModelService(httpClient);
