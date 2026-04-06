package handler

import (
	"context"
	"net/http"

	"aiminiprogram/backend-trpc/internal/wechat"

	"github.com/gin-gonic/gin"
)

// WeChatHandler 微信处理器
type WeChatHandler struct {
	client *wechat.Client
}

// NewWeChatHandler 创建微信处理器
func NewWeChatHandler(client *wechat.Client) *WeChatHandler {
	return &WeChatHandler{client: client}
}

// ==================== 1. 小程序登录（高优先级） ====================

// MiniProgramLoginRequest 小程序登录请求
type MiniProgramLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// MiniProgramLogin 小程序登录
// POST /api/wechat/login
func (h *WeChatHandler) MiniProgramLogin(c *gin.Context) {
	var req MiniProgramLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "code is required",
		})
		return
	}

	// 调用微信接口获取 session
	session, err := h.client.MiniProgram.Auth.Session(context.Background(), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "wechat login failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login success",
		"data": gin.H{
			"openid":      session.OpenID,
			"session_key": session.SessionKey,
			"unionid":     session.UnionID,
		},
	})
}

// ==================== 2. 解密手机号（高优先级） ====================

// DecryptPhoneRequest 解密手机号请求
type DecryptPhoneRequest struct {
	SessionKey    string `json:"session_key" binding:"required"`
	EncryptedData string `json:"encrypted_data" binding:"required"`
	IV            string `json:"iv" binding:"required"`
}

// DecryptPhoneNumber 解密手机号
// POST /api/wechat/phone
func (h *WeChatHandler) DecryptPhoneNumber(c *gin.Context) {
	var req DecryptPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid params",
		})
		return
	}

	// 解密手机号
	_, err := h.client.MiniProgram.Encryptor.DecryptData(
		req.EncryptedData,
		req.SessionKey,
		req.IV,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "decrypt failed",
		})
		return
	}

	// 简化返回
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "decrypt success (actual implementation needs JSON parsing)",
	})
}

// ==================== 3. 订阅消息（中优先级） ====================

// SendSubscribeMessageRequest 发送订阅消息请求
type SendSubscribeMessageRequest struct {
	ToUser     string                 `json:"touser" binding:"required"`
	TemplateID string                 `json:"template_id" binding:"required"`
	Page       string                 `json:"page,omitempty"`
	Data       map[string]interface{} `json:"data" binding:"required"`
}

// SendSubscribeMessage 发送订阅消息
// POST /api/wechat/subscribe-message
func (h *WeChatHandler) SendSubscribeMessage(c *gin.Context) {
	var req SendSubscribeMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid params",
		})
		return
	}

	// 简化实现，实际需要构建正确的订阅消息请求
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "send success (mock)",
		"data": gin.H{
			"touser":      req.ToUser,
			"template_id": req.TemplateID,
		},
	})
}

// ==================== 4. 微信支付（中优先级） ====================

// CreateOrderRequest 创建支付订单请求
type CreateOrderRequest struct {
	OpenID      string `json:"openid" binding:"required"`
	Body        string `json:"body" binding:"required"`
	OutTradeNo  string `json:"out_trade_no" binding:"required"`
	TotalFee    int64  `json:"total_fee" binding:"required"`
	Description string `json:"description,omitempty"`
}

// CreateOrder 创建微信支付订单
// POST /api/wechat/pay/order
func (h *WeChatHandler) CreateOrder(c *gin.Context) {
	if h.client.Payment == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"message": "payment not configured",
		})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid params",
		})
		return
	}

	// 简化实现
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "order created (mock)",
		"data": gin.H{
			"appId":     "wx_appid",
			"timeStamp": "1234567890",
			"nonceStr":  "nonce_str",
			"package":   "prepay_id=xxx",
			"signType":  "RSA",
			"paySign":   "sign",
			"out_trade_no": req.OutTradeNo,
		},
	})
}

// PayNotify 微信支付回调通知
// POST /api/wechat/pay/notify
func (h *WeChatHandler) PayNotify(c *gin.Context) {
	if h.client.Payment == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"message": "payment not configured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "notify received",
	})
}

// ==================== 5. 内容安全（低优先级） ====================

// CheckContentRequest 内容安全检查请求
type CheckContentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CheckContent 检查文字内容是否违规
// POST /api/wechat/security/check-content
func (h *WeChatHandler) CheckContent(c *gin.Context) {
	var req CheckContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "content is required",
		})
		return
	}

	// 简化实现
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "check success (mock)",
		"data": gin.H{
			"is_risky": false,
		},
	})
}

// CheckImage 检查图片是否违规
// POST /api/wechat/security/check-image
func (h *WeChatHandler) CheckImage(c *gin.Context) {
	_, err := c.FormFile("media")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "image is required",
		})
		return
	}

	// 简化实现
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "check success (mock)",
		"data": gin.H{
			"is_risky": false,
		},
	})
}

// ==================== 6. 二维码生成（低优先级） ====================

// CreateQRCodeRequest 创建二维码请求
type CreateQRCodeRequest struct {
	Scene string `json:"scene" binding:"required"`
	Page  string `json:"page,omitempty"`
	Width int    `json:"width,omitempty"`
}

// CreateQRCode 创建小程序二维码
// POST /api/wechat/qrcode
func (h *WeChatHandler) CreateQRCode(c *gin.Context) {
	var req CreateQRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "scene is required",
		})
		return
	}

	if req.Width == 0 {
		req.Width = 430
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "qrcode created (mock)",
		"data": gin.H{
			"scene": req.Scene,
			"page":  req.Page,
			"width": req.Width,
		},
	})
}
