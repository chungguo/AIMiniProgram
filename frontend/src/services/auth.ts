import { interceptorManager } from './interceptor';

// ==================== Token 管理 ====================

interface TokenState {
  value: string;
  isRefreshing: boolean;
  refreshPromise: Promise<string> | null;
}

const tokenState: TokenState = {
  value: uni.getStorageSync('token') || '',
  isRefreshing: false,
  refreshPromise: null
};

// ==================== 工具函数 ====================

/**
 * 获取当前 token
 */
export const getCurrentToken = (): string => tokenState.value;

/**
 * 设置 token
 */
export const setToken = (token: string): void => {
  tokenState.value = token;
  uni.setStorageSync('token', token);
};

/**
 * 清空 token
 */
export const clearToken = (): void => {
  tokenState.value = '';
  tokenState.isRefreshing = false;
  tokenState.refreshPromise = null;
  uni.removeStorageSync('token');
  uni.removeStorageSync('tokenExpireTime');
};

/**
 * 刷新 token（需要项目实现具体的刷新逻辑）
 */
export const refreshToken = async (): Promise<string> => {
  // 这里需要根据实际项目实现 token 刷新逻辑
  // 例如：调用刷新接口，使用 refresh_token 换取新的 access_token
  throw new Error('refreshToken not implemented');
};

// ==================== 增强的 Token 管理 ====================

/**
 * 获取有效的 token（带自动刷新）
 */
export const getValidToken = async (): Promise<string> => {
  // 如果已有有效 token，直接返回
  if (tokenState.value) {
    return tokenState.value;
  }

  // 如果正在刷新中，排队等待
  if (tokenState.isRefreshing && tokenState.refreshPromise) {
    return tokenState.refreshPromise;
  }

  // 发起新的刷新请求
  tokenState.isRefreshing = true;
  tokenState.refreshPromise = refreshToken().finally(() => {
    tokenState.isRefreshing = false;
    tokenState.refreshPromise = null;
  });

  return tokenState.refreshPromise;
};

// ==================== 增强拦截器配置 ====================

/**
 * 设置增强的请求拦截器（带 token 自动刷新）
 */
export function setupEnhancedInterceptors(): void {
  // 请求拦截器 - 自动添加 token
  interceptorManager.useRequest(async (url, options) => {
    const token = getCurrentToken();
    if (token) {
      options.header = {
        ...options.header,
        'Authorization': `Bearer ${token}`
      };
    }
    return { url, options };
  });

  // 响应拦截器 - 统一处理响应
  interceptorManager.useResponse(async (response: unknown) => {
    const res = response as { success?: boolean; message?: string; data?: unknown };
    
    // 统一错误处理
    if (res.success === false) {
      uni.showToast({
        title: res.message || '请求失败',
        icon: 'none',
        duration: 3000
      });
    }
    
    return response;
  });

  // 错误拦截器 - 处理 401 未授权
  interceptorManager.useError(async (error: Error) => {
    if (error.message.includes('401')) {
      // Token 过期，尝试刷新
      try {
        await getValidToken();
        // 刷新成功，可以在这里触发重试逻辑
      } catch {
        // 刷新失败，跳转到登录
        uni.showToast({
          title: '登录已过期，请重新登录',
          icon: 'none'
        });
        // 清除登录态
        clearToken();
        // 延迟跳转登录页
        setTimeout(() => {
          uni.navigateTo({ url: '/pages/login/login' });
        }, 1500);
      }
    }
    return error;
  });
}

/**
 * 重试失败的请求
 */
export async function retryRequest<T>(
  requestFn: () => Promise<T>,
  maxRetries: number = 3
): Promise<T> {
  let lastError: Error;
  
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await requestFn();
    } catch (error) {
      lastError = error as Error;
      // 如果不是网络错误，直接抛出
      if (!isNetworkError(lastError)) {
        throw lastError;
      }
      // 网络错误，等待后重试
      await delay(1000 * Math.pow(2, i)); // 指数退避
    }
  }
  
  throw lastError!;
}

/**
 * 判断是否为网络错误
 */
function isNetworkError(error: Error): boolean {
  const networkErrorMessages = [
    'request:fail',
    'timeout',
    'network',
    '连接超时',
    '网络错误'
  ];
  return networkErrorMessages.some(msg => 
    error.message.toLowerCase().includes(msg.toLowerCase())
  );
}

/**
 * 延迟函数
 */
function delay(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}
