/**
 * 拦截器使用示例
 * 
 * 在应用的入口文件（如 main.ts）中配置拦截器
 */

import { 
  interceptorManager, 
  setupDefaultInterceptors,
  setupAuthInterceptor 
} from '@/services';

// ============ 1. 基础配置 ============

// 设置默认拦截器（错误提示、请求头等）
setupDefaultInterceptors();

// ============ 2. 请求拦截器示例 ============

// 添加 Token 到请求头
const removeAuthHeaderInterceptor = interceptorManager.useRequest(async (url, options) => {
  const token = uni.getStorageSync('access_token');
  if (token) {
    options.header = {
      ...options.header,
      'Authorization': `Bearer ${token}`
    };
  }
  return { url, options };
});

// 添加时间戳防止缓存
interceptorManager.useRequest(async (url, options) => {
  if (options.method === 'GET') {
    url += (url.includes('?') ? '&' : '?') + `_t=${Date.now()}`;
  }
  return { url, options };
});

// ============ 3. 响应拦截器示例 ============

// 统一处理响应数据格式
interceptorManager.useResponse(async (response) => {
  // 可以在这里对响应数据进行转换
  console.log('Response:', response);
  return response;
});

// ============ 4. 错误拦截器示例 ============

// 统一错误提示
interceptorManager.useError(async (error) => {
  // 根据错误类型显示不同的提示
  if (error.message.includes('401')) {
    uni.showToast({
      title: '登录已过期，请重新登录',
      icon: 'none',
      duration: 2000
    });
  } else if (error.message.includes('500')) {
    uni.showToast({
      title: '服务器错误，请稍后重试',
      icon: 'none',
      duration: 2000
    });
  } else {
    uni.showToast({
      title: error.message || '网络错误',
      icon: 'none',
      duration: 2000
    });
  }
  return error;
});

// ============ 5. 登录态自动续期 ============

// 配置 token 刷新逻辑
const removeAuthInterceptor = setupAuthInterceptor(async () => {
  // 调用刷新 token 接口
  const refreshToken = uni.getStorageSync('refresh_token');
  const res = await uni.request({
    url: '/api/auth/refresh',
    method: 'POST',
    data: { refreshToken }
  });
  
  if (res.statusCode === 200) {
    const newToken = (res.data as { token: string }).token;
    uni.setStorageSync('access_token', newToken);
    uni.setStorageSync('tokenExpireTime', Date.now() + 2 * 60 * 60 * 1000); // 2小时
    return newToken;
  }
  throw new Error('Refresh token failed');
});

// ============ 6. 取消拦截器 ============

// 如果需要取消某个拦截器
// removeAuthHeaderInterceptor();
// removeAuthInterceptor();

// ============ 7. 自定义拦截器示例 ============

// 请求日志拦截器
interceptorManager.useRequest(async (url, options) => {
  console.log(`[Request] ${options.method} ${url}`, options.data);
  return { url, options };
});

// 响应日志拦截器
interceptorManager.useResponse(async (response) => {
  console.log('[Response]', response);
  return response;
});

// 网络错误重试拦截器（示例）
interceptorManager.useError(async (error) => {
  // 可以实现自动重试逻辑
  // 注意：需要处理重试次数避免无限循环
  console.log('[Error]', error.message);
  return error;
});

/**
 * 使用方式：
 * 
 * 1. 在 main.ts 中导入并执行此文件：
 *    import '@/utils/interceptor-example';
 * 
 * 2. 然后正常调用服务：
 *    import { modelService } from '@/services';
 *    const models = await modelService.getModels();
 * 
 * 3. 所有请求/响应/错误都会经过拦截器处理
 */
