package system

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	contract "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	svc *system.CaptchaService
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler(svc *system.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{svc: svc}
}

// Get 获取验证码
// @Router /system/captcha/get [post]
// 注意：验证码接口不使用项目标准响应格式，直接返回 aj-captcha 原始格式
// 前端 aj-captcha 组件期望响应结构为 {repCode, repData}，而非 {code, msg, data}
func (h *CaptchaHandler) Get(c *gin.Context) {
	var req contract.CaptchaGetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 即使绑定失败也继续，因为请求体可能为空
	}

	result, err := h.svc.Generate(c.Request.Context())
	if err != nil {
		// 直接返回 aj-captcha 格式
		c.JSON(200, contract.CaptchaGetResp{
			RepCode: "6110",
			RepMsg:  "获取验证码失败: " + err.Error(),
			RepData: nil,
		})
		return
	}

	// 直接返回 aj-captcha 格式，前端期望 res.data.repCode 能直接访问
	c.JSON(200, contract.CaptchaGetResp{
		RepCode: "0000",
		RepMsg:  "success",
		RepData: &contract.CaptchaRepData{
			OriginalImageBase64: result.OriginalImageBase64,
			JigsawImageBase64:   result.JigsawImageBase64,
			Token:               result.Token,
			SecretKey:           "", // 暂不使用加密
		},
	})
}

// Check 校验验证码
// @Router /system/captcha/check [post]
// 注意：验证码接口不使用项目标准响应格式，直接返回 aj-captcha 原始格式
func (h *CaptchaHandler) Check(c *gin.Context) {
	var req contract.CaptchaCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, contract.CaptchaCheckResp{
			RepCode: "6111",
			RepMsg:  "参数错误",
		})
		return
	}

	// 解析坐标
	var point contract.CaptchaPoint
	if err := json.Unmarshal([]byte(req.PointJson), &point); err != nil {
		c.JSON(200, contract.CaptchaCheckResp{
			RepCode: "6111",
			RepMsg:  "坐标格式错误",
		})
		return
	}

	// 校验（将浮点坐标转换为整数）
	valid, err := h.svc.Verify(c.Request.Context(), req.Token, int(point.X), int(point.Y))
	if err != nil {
		c.JSON(200, contract.CaptchaCheckResp{
			RepCode: "6112",
			RepMsg:  "校验失败: " + err.Error(),
		})
		return
	}

	if !valid {
		c.JSON(200, contract.CaptchaCheckResp{
			RepCode: "6111",
			RepMsg:  "验证失败",
		})
		return
	}

	c.JSON(200, contract.CaptchaCheckResp{
		RepCode: "0000",
		RepMsg:  "success",
	})
}
