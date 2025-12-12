package pay

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PayClientConfig 支付客户端配置 - 联合类型
// 根据 ConfigType 字段判断使用哪种配置
// 对应 Java 中的多态 PayClientConfig 接口
type PayClientConfig struct {
	// ConfigType 配置类型标识 (从 JSON @class 字段解析)
	// 可能的值: "WxPayClientConfig", "AlipayPayClientConfig", "NonePayClientConfig"
	ConfigType string `json:"@class,omitempty"`

	// ========== 微信支付配置 (当 ConfigType 包含 "WxPayClientConfig") ==========
	// 公共字段
	AppID      string `json:"appId,omitempty"`      // 公众号或小程序的 appid
	MchID      string `json:"mchId,omitempty"`      // 商户号
	APIVersion string `json:"apiVersion,omitempty"` // API 版本: "v2" 或 "v3"
	// V2 版本
	MchKey     string `json:"mchKey,omitempty"`     // 商户密钥 (V2)
	KeyContent string `json:"keyContent,omitempty"` // apiclient_cert.p12 证书 base64 (V2)
	// V3 版本
	PrivateKeyContent string `json:"privateKeyContent,omitempty"` // apiclient_key.pem 证书 (V3)
	APIV3Key          string `json:"apiV3Key,omitempty"`          // apiV3 密钥值 (V3)
	CertSerialNo      string `json:"certSerialNo,omitempty"`      // 证书序列号 (V3)
	PublicKeyContent  string `json:"publicKeyContent,omitempty"`  // pub_key.pem 证书 (V3)
	PublicKeyID       string `json:"publicKeyId,omitempty"`       // publicKeyId (V3)

	// ========== 支付宝配置 (当 ConfigType 包含 "AlipayPayClientConfig") ==========
	ServerURL string `json:"serverUrl,omitempty"` // 网关地址
	// AppID 与微信共用
	SignType string `json:"signType,omitempty"` // 签名算法类型，推荐 RSA2
	Mode     *int   `json:"mode,omitempty"`     // 公钥类型: 1=公钥模式, 2=证书模式
	// 公钥模式 (Mode=1)
	PrivateKey      string `json:"privateKey,omitempty"`      // 商户私钥
	AlipayPublicKey string `json:"alipayPublicKey,omitempty"` // 支付宝公钥字符串
	// 证书模式 (Mode=2)
	AppCertContent          string `json:"appCertContent,omitempty"`          // 商户公钥应用证书内容
	AlipayPublicCertContent string `json:"alipayPublicCertContent,omitempty"` // 支付宝公钥证书内容
	RootCertContent         string `json:"rootCertContent,omitempty"`         // 根证书内容
	// 可选加密
	EncryptType string `json:"encryptType,omitempty"` // 接口内容加密方式 "AES"
	EncryptKey  string `json:"encryptKey,omitempty"`  // 接口内容加密私钥
}

// ========== 常量定义 ==========

// ConfigType 常量
const (
	ConfigTypeWxPay  = "cn.iocoder.yudao.module.pay.framework.pay.core.client.impl.weixin.WxPayClientConfig"
	ConfigTypeAlipay = "cn.iocoder.yudao.module.pay.framework.pay.core.client.impl.alipay.AlipayPayClientConfig"
	ConfigTypeNone   = "cn.iocoder.yudao.module.pay.framework.pay.core.client.impl.NonePayClientConfig"
)

// 微信 API 版本
const (
	WxAPIVersionV2 = "v2"
	WxAPIVersionV3 = "v3"
)

// 支付宝公钥类型
const (
	AlipayModePublicKey   = 1 // 公钥模式
	AlipayModeCertificate = 2 // 证书模式
)

// ========== GORM Scanner/Valuer 实现 ==========

// Scan 从数据库读取 JSON 并解析
func (c *PayClientConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("incompatible type for PayClientConfig")
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, c)
}

// Value 序列化为 JSON 存储到数据库
func (c PayClientConfig) Value() (driver.Value, error) {
	if c.ConfigType == "" {
		return nil, nil
	}
	return json.Marshal(c)
}

// ========== 辅助方法 ==========

// IsWxPay 判断是否为微信支付配置
func (c *PayClientConfig) IsWxPay() bool {
	return c.ConfigType == ConfigTypeWxPay
}

// IsAlipay 判断是否为支付宝配置
func (c *PayClientConfig) IsAlipay() bool {
	return c.ConfigType == ConfigTypeAlipay
}

// IsNone 判断是否为空配置
func (c *PayClientConfig) IsNone() bool {
	return c.ConfigType == ConfigTypeNone
}

// ToJSON 序列化为 JSON 字符串 (用于传递给支付客户端)
func (c *PayClientConfig) ToJSON() string {
	if c == nil {
		return ""
	}
	data, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(data)
}
