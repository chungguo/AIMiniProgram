// 服务统一导出
export { httpClient, modelService, createModelService } from './modelService';
export { paperService, createPaperService } from './paperService';
export { artificialAnalysisService, createArtificialAnalysisService } from './artificialAnalysisService';

// 拦截器导出
export {
  interceptorManager,
  InterceptableHttpClient,
  setupDefaultInterceptors,
  setupAuthInterceptor
} from './interceptor';

// Auth 导出
export {
  getCurrentToken,
  setToken,
  clearToken,
  getValidToken,
  setupEnhancedInterceptors,
  retryRequest
} from './auth';

// 类型导出
export type { 
  IModelService, 
  IPaperService, 
  IArtificialAnalysisService,
  IHttpClient,
  IHttpInterceptor,
  IInterceptorManager
} from '@/types/api';
