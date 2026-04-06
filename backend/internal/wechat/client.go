package wechat

import (
	"os"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment"
)

// Client 微信客户端
type Client struct {
	MiniProgram *miniProgram.MiniProgram
	Payment     *payment.Payment
}

// Config 微信配置
type Config struct {
	AppID     string
	Secret    string
	MchID     string
	MchKey    string
	NotifyURL string
}

// NewClient 创建微信客户端
func NewClient(cfg *Config) (*Client, error) {
	client := &Client{}

	// 初始化小程序客户端
	miniProgramApp, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:  cfg.AppID,
		Secret: cfg.Secret,
	})
	if err != nil {
		return nil, err
	}
	client.MiniProgram = miniProgramApp

	// 初始化支付（如果配置了商户信息）
	if cfg.MchID != "" {
		paymentApp, err := payment.NewPayment(&payment.UserConfig{
			AppID:       cfg.AppID,
			MchID:       cfg.MchID,
			MchApiV3Key: cfg.MchKey,
			Key:         cfg.MchKey,
			NotifyURL:   cfg.NotifyURL,
		})
		if err != nil {
			return nil, err
		}
		client.Payment = paymentApp
	}

	return client, nil
}

// NewClientFromEnv 从环境变量创建客户端
func NewClientFromEnv() (*Client, error) {
	cfg := &Config{
		AppID:     os.Getenv("WECHAT_APPID"),
		Secret:    os.Getenv("WECHAT_SECRET"),
		MchID:     os.Getenv("WECHAT_MCH_ID"),
		MchKey:    os.Getenv("WECHAT_MCH_KEY"),
		NotifyURL: os.Getenv("WECHAT_NOTIFY_URL"),
	}
	return NewClient(cfg)
}
