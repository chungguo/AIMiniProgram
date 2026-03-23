# Frontend

uni-app 微信小程序前端。

## 技术栈

- Vue 3 + TypeScript
- uni-app (微信小程序/H5/App)
- UnoCSS (原子化 CSS)

## 结构

```
src/
├── pages/      # 页面
├── services/   # API 服务
├── types/      # TypeScript 类型
└── utils/      # 工具函数
```

## 运行

```bash
npm run dev:h5        # H5 开发
npm run dev:mp-weixin # 微信小程序
```

## 服务层

- `modelService` - 模型 API
- `paperService` - 论文 API
- `artificialAnalysisService` - 评测数据 API
- `interceptorManager` - HTTP 拦截器
