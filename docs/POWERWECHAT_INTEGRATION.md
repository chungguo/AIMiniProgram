# PowerWeChat 集成文档

ModelLens 已集成 PowerWeChat SDK，支持以下微信功能：

## 已实现功能

| 优先级 | 功能 | 端点 | 状态 |
|--------|------|------|------|
| 🔴 高 | 小程序登录 | POST /api/wechat/login | ✅ |
| 🔴 高 | 解密手机号 | POST /api/wechat/phone | ✅ |
| 🟡 中 | 订阅消息 | POST /api/wechat/subscribe-message | ✅ |
| 🟡 中 | 微信支付 | POST /api/wechat/pay/order | ✅ |
| 🟢 低 | 内容安全 | POST /api/wechat/security/check-content | ✅ |
| 🟢 低 | 二维码生成 | POST /api/wechat/qrcode | ✅ |

---

## 配置方法

### 1. 设置环境变量

```bash
# 复制示例文件
cp backend-trpc/.env.example backend-trpc/.env

# 编辑 .env 文件，填写你的微信小程序配置
WECHAT_APPID=wx1234567890abcdef
WECHAT_SECRET=your_app_secret_here
WECHAT_MCH_ID=1234567890          # 可选，支付用
WECHAT_MCH_KEY=your_mch_key_here   # 可选，支付用
```

### 2. 获取微信小程序配置

登录 [微信公众平台](https://mp.weixin.qq.com/)：
- **AppID**: 开发 > 开发设置 > AppID
- **AppSecret**: 开发 > 开发设置 > AppSecret（需重置获取）
- **支付商户号**: 微信支付 > 商户号

---

## 前端调用示例

### 1. 小程序登录

```typescript
// services/wechat.ts
export async function wechatLogin(): Promise<{ openid: string; session_key: string }> {
  // 获取微信登录 code
  const [loginErr, loginRes] = await uni.login({ provider: 'weixin' });
  
  if (loginErr) {
    throw new Error('微信登录失败');
  }

  // 发送到后端换取 session
  const response = await fetch('/api/wechat/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ code: loginRes.code })
  });

  const result = await response.json();
  if (!result.success) {
    throw new Error(result.message);
  }

  // 保存登录态
  uni.setStorageSync('openid', result.data.openid);
  uni.setStorageSync('session_key', result.data.session_key);

  return result.data;
}
```

### 2. 获取手机号

```vue
<template>
  <button 
    open-type="getPhoneNumber" 
    @getphonenumber="onGetPhoneNumber"
  >
    获取手机号
  </button>
</template>

<script setup>
async function onGetPhoneNumber(e: any) {
  if (e.detail.errMsg !== 'getPhoneNumber:ok') {
    uni.showToast({ title: '用户拒绝授权', icon: 'none' });
    return;
  }

  const sessionKey = uni.getStorageSync('session_key');
  
  const response = await fetch('/api/wechat/phone', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      session_key: sessionKey,
      encrypted_data: e.detail.encryptedData,
      iv: e.detail.iv
    })
  });

  const result = await response.json();
  if (result.success) {
    console.log('手机号:', result.data.phoneNumber);
  }
}
</script>
```

### 3. 请求订阅消息授权

```typescript
// 请求订阅消息授权
async function requestSubscribeMessage() {
  const tmplIds = ['xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx']; // 你的模板ID
  
  const [err, res] = await uni.requestSubscribeMessage({ tmplIds });
  
  if (err) {
    console.error('订阅失败:', err);
    return;
  }

  console.log('订阅结果:', res);
  // res[templateId] === 'accept' 表示用户同意
}

// 发送订阅消息
async function sendSubscribeMessage(data: any) {
  const response = await fetch('/api/wechat/subscribe-message', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      touser: uni.getStorageSync('openid'),
      template_id: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
      page: 'pages/index/index',
      data: {
        thing1: { value: '预约提醒' },
        time2: { value: '2024-01-01 10:00' },
        thing3: { value: '请准时就诊' }
      }
    })
  });

  return response.json();
}
```

### 4. 微信支付

```typescript
// 创建订单并调起支付
async function createPayment(orderInfo: any) {
  // 1. 创建订单
  const orderRes = await fetch('/api/wechat/pay/order', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      openid: uni.getStorageSync('openid'),
      body: orderInfo.title,
      out_trade_no: orderInfo.orderNo,
      total_fee: orderInfo.amount * 100, // 单位：分
      description: orderInfo.description
    })
  });

  const orderResult = await orderRes.json();
  if (!orderResult.success) {
    throw new Error(orderResult.message);
  }

  // 2. 调起微信支付
  const [payErr, payRes] = await uni.requestPayment({
    provider: 'wxpay',
    ...orderResult.data
  });

  if (payErr) {
    throw new Error('支付失败: ' + payErr.errMsg);
  }

  return payRes;
}
```

### 5. 内容安全检测

```typescript
// 检测文字内容
async function checkContent(content: string): Promise<boolean> {
  const response = await fetch('/api/wechat/security/check-content', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ content })
  });

  const result = await response.json();
  return !result.data.is_risky;
}

// 检测图片
async function checkImage(filePath: string): Promise<boolean> {
  const response = await fetch('/api/wechat/security/check-image', {
    method: 'POST',
    body: (() => {
      const formData = new FormData();
      formData.append('media', filePath);
      return formData;
    })()
  });

  const result = await response.json();
  return !result.data.is_risky;
}
```

---

## 后端 API 文档

### POST /api/wechat/login
小程序登录

**请求参数：**
```json
{
  "code": "wx_login_code_from_frontend"
}
```

**响应：**
```json
{
  "success": true,
  "data": {
    "openid": "o1234567890abcdef",
    "session_key": "xxxxx",
    "unionid": "xxxxx"
  }
}
```

### POST /api/wechat/phone
解密手机号

**请求参数：**
```json
{
  "session_key": "from_login",
  "encrypted_data": "from_getPhoneNumber",
  "iv": "from_getPhoneNumber"
}
```

### POST /api/wechat/subscribe-message
发送订阅消息

**请求参数：**
```json
{
  "touser": "openid",
  "template_id": "template_id",
  "page": "pages/index/index",
  "data": {
    "thing1": { "value": "内容" },
    "time2": { "value": "2024-01-01" }
  }
}
```

### POST /api/wechat/pay/order
创建支付订单

**请求参数：**
```json
{
  "openid": "user_openid",
  "body": "商品描述",
  "out_trade_no": "ORDER_123456",
  "total_fee": 100,
  "description": "订单描述"
}
```

---

## 注意事项

1. **HTTPS**: 生产环境必须使用 HTTPS
2. **IP 白名单**: 在微信公众平台配置服务器 IP 地址
3. **域名配置**: 在小程序后台配置 request 合法域名
4. **模板消息**: 订阅消息模板需要在小程序后台申请
5. **支付授权**: 微信支付需要单独申请开通

---

## 参考链接

- [PowerWeChat 文档](https://github.com/ArtisanCloud/PowerWeChat)
- [微信小程序登录](https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/login.html)
- [小程序订阅消息](https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/subscribe-message.html)
- [微信支付](https://pay.weixin.qq.com/wiki/doc/apiv3/index.shtml)
