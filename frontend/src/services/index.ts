// 服务统一导出
export { httpClient, modelService, createModelService } from './modelService';
export { paperService, createPaperService } from './paperService';

// 类型导出
export type { 
  IModelService, 
  IPaperService, 
  IHttpClient 
} from '@/types/api';
