package system

// CaptchaGetReq 获取验证码请求
// 兼容 Java aj-captcha 的请求格式
type CaptchaGetReq struct {
	CaptchaType string `json:"captchaType"` // 验证码类型：slide / blockPuzzle
	BrowserInfo string `json:"browserInfo"` // 浏览器指纹（IP + User-Agent）
}

// CaptchaRepData 验证码响应数据（aj-captcha 格式）
type CaptchaRepData struct {
	OriginalImageBase64 string `json:"originalImageBase64"` // 主图 Base64
	JigsawImageBase64   string `json:"jigsawImageBase64"`   // 拼图块 Base64
	Token               string `json:"token"`               // 验证码唯一标识
	SecretKey           string `json:"secretKey"`           // 加密密钥（可选，用于前端加密坐标）
}

// CaptchaGetResp 获取验证码响应
// 兼容 Java aj-captcha 的响应格式，包装在项目标准 Result 中
type CaptchaGetResp struct {
	RepCode string          `json:"repCode"` // 状态码：0000 成功
	RepMsg  string          `json:"repMsg"`  // 状态信息
	RepData *CaptchaRepData `json:"repData"` // 验证码数据
}

// CaptchaCheckReq 校验验证码请求
type CaptchaCheckReq struct {
	Token       string `json:"token"`       // 验证码唯一标识
	PointJson   string `json:"pointJson"`   // 坐标 JSON，格式：{"x":100,"y":50}
	CaptchaType string `json:"captchaType"` // 验证码类型
	BrowserInfo string `json:"browserInfo"` // 浏览器指纹
}

// CaptchaCheckResp 校验验证码响应
type CaptchaCheckResp struct {
	RepCode string `json:"repCode"` // 状态码：0000 成功
	RepMsg  string `json:"repMsg"`  // 状态信息
}

// CaptchaPoint 验证码坐标点
// 注意：前端 aj-captcha 发送的坐标是浮点数，需要用 float64 接收
type CaptchaPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
