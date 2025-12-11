package weixin

// WxPayClientConfig 微信支付客户端配置
// 与 Java WxPayClientConfig 保持一致
type WxPayClientConfig struct {
	// ========== 通用参数 ==========
	// 公众号或者小程序的 appid
	AppID string `json:"appId"`
	// 商户号
	MchID string `json:"mchId"`
	// API 版本 (v2 / v3)
	APIVersion string `json:"apiVersion"`

	// ========== V2 版本的参数 ==========
	// 商户密钥
	MchKey string `json:"mchKey,omitempty"`
	// apiclient_cert.p12 证书文件的对应字符串【base64 格式】
	KeyContent string `json:"keyContent,omitempty"`

	// ========== V3 版本的参数 ==========
	// apiclient_key.pem 证书文件的对应字符串
	PrivateKeyContent string `json:"privateKeyContent,omitempty"`
	// apiV3 密钥值
	APIV3Key string `json:"apiV3Key,omitempty"`
	// 证书序列号（merchantSerialNumber）
	CertSerialNo string `json:"certSerialNo,omitempty"`
	// pub_key.pem 证书文件的对应字符串
	PublicKeyContent string `json:"publicKeyContent,omitempty"`
	// 公钥 ID
	PublicKeyID string `json:"publicKeyId,omitempty"`
}

const (
	APIVersionV2 = "v2"
	APIVersionV3 = "v3"
)
