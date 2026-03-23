import type { IHttpClient, IHttpInterceptor } from '@/types/api';

// 拦截器管理器
class HttpInterceptorManager {
  private requestInterceptors: IHttpInterceptor['request'][] = [];
  private responseInterceptors: IHttpInterceptor['response'][] = [];
  private errorInterceptors: IHttpInterceptor['error'][] = [];

  // 注册请求拦截器
  useRequest(interceptor: IHttpInterceptor['request']): () => void {
    this.requestInterceptors.push(interceptor);
    // 返回取消注册的函数
    return () => {
      const index = this.requestInterceptors.indexOf(interceptor);
      if (index > -1) {
        this.requestInterceptors.splice(index, 1);
      }
    };
  }

  // 注册响应拦截器
  useResponse(interceptor: IHttpInterceptor['response']): () => void {
    this.responseInterceptors.push(interceptor);
    return () => {
      const index = this.responseInterceptors.indexOf(interceptor);
      if (index > -1) {
        this.responseInterceptors.splice(index, 1);
      }
    };
  }

  // 注册错误拦截器
  useError(interceptor: IHttpInterceptor['error']): () => void {
    this.errorInterceptors.push(interceptor);
    return () => {
      const index = this.errorInterceptors.indexOf(interceptor);
      if (index > -1) {
        this.errorInterceptors.splice(index, 1);
      }
    };
  }

  // 执行请求拦截器
  async executeRequestInterceptors(
    url: string,
    options: { method: string; data?: unknown; header?: Record<string, string> }
  ): Promise<{ url: string; options: typeof options }> {
    let result = { url, options };
    for (const interceptor of this.requestInterceptors) {
      result = await interceptor(result.url, result.options);
    }
    return result;
  }

  // 执行响应拦截器
  async executeResponseInterceptors<T>(response: T): Promise<T> {
    let result = response;
    for (const interceptor of this.responseInterceptors) {
      result = await interceptor(result);
    }
    return result;
  }

  // 执行错误拦截器
  async executeErrorInterceptors(error: Error): Promise<Error> {
    let result = error;
    for (const interceptor of this.errorInterceptors) {
      result = await interceptor(result);
    }
    return result;
  }
}

// 创建全局拦截器管理器实例
export const interceptorManager = new HttpInterceptorManager();

// 带拦截器的 HttpClient 包装器
export class InterceptableHttpClient implements IHttpClient {
  private baseClient: IHttpClient;

  constructor(baseClient: IHttpClient) {
    this.baseClient = baseClient;
  }

  async get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
    try {
      // 执行请求拦截器
      const { url: processedUrl, options } = await interceptorManager.executeRequestInterceptors(
        url,
        { method: 'GET' }
      );

      // 处理 URL 参数
      const fullURL = this.buildURL(processedUrl, params);

      // 发送请求
      const response = await this.baseClient.get<T>(fullURL);

      // 执行响应拦截器
      return await interceptorManager.executeResponseInterceptors(response);
    } catch (error) {
      // 执行错误拦截器
      const processedError = await interceptorManager.executeErrorInterceptors(error as Error);
      throw processedError;
    }
  }

  async post<T>(url: string, data?: unknown): Promise<T> {
    try {
      // 执行请求拦截器
      const { url: processedUrl, options } = await interceptorManager.executeRequestInterceptors(
        url,
        { method: 'POST', data, header: { 'Content-Type': 'application/json' } }
      );

      // 发送请求
      const response = await this.baseClient.post<T>(processedUrl, options.data);

      // 执行响应拦截器
      return await interceptorManager.executeResponseInterceptors(response);
    } catch (error) {
      // 执行错误拦截器
      const processedError = await interceptorManager.executeErrorInterceptors(error as Error);
      throw processedError;
    }
  }

  private buildURL(url: string, params?: Record<string, unknown>): string {
    if (!params || Object.keys(params).length === 0) return url;

    const queryParams = Object.entries(params)
      .filter(([, value]) => value !== undefined && value !== null)
      .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
      .join('&');

    return queryParams ? `${url}?${queryParams}` : url;
  }
}

// 默认的错误处理拦截器（示例）
export function setupDefaultInterceptors(): void {
  // 请求拦截器 - 添加通用请求头
  interceptorManager.useRequest(async (url, options) => {
    // 可以在这里添加 token 等
    const token = uni.getStorageSync('token');
    if (token) {
      options.header = {
        ...options.header,
        'Authorization': `Bearer ${token}`
      };
    }
    return { url, options };
  });

  // 响应拦截器 - 统一处理响应
  interceptorManager.useResponse(async (response) => {
    // 可以在这里统一处理响应数据
    return response;
  });

  // 错误拦截器 - 统一处理错误
  interceptorManager.useError(async (error) => {
    // 统一错误提示
    if (error.message) {
      uni.showToast({
        title: error.message,
        icon: 'none',
        duration: 3000
      });
    }
    return error;
  });
}

// 登录态自动续期拦截器（示例）
export function setupAuthInterceptor(refreshTokenFn: () => Promise<string>): () => void {
  const removeRequestInterceptor = interceptorManager.useRequest(async (url, options) => {
    // 检查 token 是否即将过期，如果是则先刷新
    const tokenExpireTime = uni.getStorageSync('tokenExpireTime');
    if (tokenExpireTime && Date.now() > tokenExpireTime - 5 * 60 * 1000) {
      // token 将在 5 分钟内过期，尝试刷新
      try {
        const newToken = await refreshTokenFn();
        uni.setStorageSync('token', newToken);
        options.header = {
          ...options.header,
          'Authorization': `Bearer ${newToken}`
        };
      } catch {
        // 刷新失败，跳转到登录页
        uni.navigateTo({ url: '/pages/login/login' });
      }
    }
    return { url, options };
  });

  const removeErrorInterceptor = interceptorManager.useError(async (error) => {
    // 处理 401 未授权错误
    if (error.message.includes('401')) {
      uni.showToast({
        title: '登录已过期，请重新登录',
        icon: 'none'
      });
      uni.navigateTo({ url: '/pages/login/login' });
    }
    return error;
  });

  // 返回取消注册的函数
  return () => {
    removeRequestInterceptor();
    removeErrorInterceptor();
  };
}
