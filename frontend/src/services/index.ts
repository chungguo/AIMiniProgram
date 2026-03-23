// 服务统一导出
export { httpClient, modelService, createModelService } from './modelService';
export { paperService, createPaperService } from './paperService';
export { artificialAnalysisService, createArtificialAnalysisService } from './artificialAnalysisService';

// 类型导出
export type { 
  IModelService, 
  IPaperService, 
  IArtificialAnalysisService,
  IHttpClient 
} from '@/types/api';
