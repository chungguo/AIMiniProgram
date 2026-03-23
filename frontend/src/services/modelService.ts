import type { 
  IHttpClient, 
  IModelService, 
  Model, 
  ModelQueryParams,
  CompareResult,
  ComparisonCategory,
  PaginatedData
} from '@/types/api';

// 从响应头中提取错误信息，支持 base64 编码
function extractErrorMessage(res: { statusCode: number; header?: Record<string, unknown> }): string {
  let errorMessage = (res.header?.['X-Error-Message'] as string) || `HTTP ${res.statusCode}`;
  const encoding = res.header?.['X-Error-Message-Encoding'] as string;
  
  // 如果是 base64 编码，进行解码
  if (encoding === 'base64' && errorMessage) {
    try {
      errorMessage = atob(errorMessage);
    } catch {
      // 解码失败，使用原始消息
    }
  }
  
  return errorMessage;
}

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
            reject(new Error(extractErrorMessage(res)));
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
            reject(new Error(extractErrorMessage(res)));
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

  async getModels(params?: ModelQueryParams): Promise<PaginatedData<Model[]>> {
    return this.httpClient.get<PaginatedData<Model[]>>('/models', params as Record<string, unknown>);
  }

  async getModelById(id: string): Promise<Model> {
    return this.httpClient.get<Model>(`/models/detail/${id}`);
  }

  async getFamilies(): Promise<string[]> {
    return this.httpClient.get<string[]>('/models/families');
  }

  async getFamilyModels(family: string): Promise<Model[]> {
    return this.httpClient.get<Model[]>(`/models/family/${family}`);
  }

  async compareModels(ids: string[]): Promise<CompareResult> {
    return this.httpClient.post<CompareResult>('/models/compare', { ids });
  }

  async getComparisonCategories(): Promise<ComparisonCategory[]> {
    return this.httpClient.get<ComparisonCategory[]>('/models/meta/comparison-categories');
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
