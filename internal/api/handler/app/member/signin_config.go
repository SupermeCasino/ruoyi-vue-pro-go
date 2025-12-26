package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// AppMemberSignInConfigHandler App端签到配置 Handler
type AppMemberSignInConfigHandler struct {
	svc *memberSvc.MemberSignInConfigService
}

// NewAppMemberSignInConfigHandler 创建 App端签到配置 Handler
func NewAppMemberSignInConfigHandler(svc *memberSvc.MemberSignInConfigService) *AppMemberSignInConfigHandler {
	return &AppMemberSignInConfigHandler{svc: svc}
}

// GetSignInConfigList 获得签到规则列表
// @Summary 获得签到规则列表
// @Description 获取启用状态的签到规则列表，按天数升序排列
// @Tags App - 签到规则
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]resp.AppMemberSignInConfigResp}
// @Router /app-api/member/sign-in/config/list [get]
func (h *AppMemberSignInConfigHandler) GetSignInConfigList(c *gin.Context) {
	// 对齐 Java: 只获取启用状态的配置 (CommonStatusEnum.ENABLE.getStatus() = 0)
	status := 0
	list, err := h.svc.GetSignInConfigList(c, &status)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 转换为 App 端响应结构 (只返回 day 和 point)
	respList := lo.Map(list, func(item *memberModel.MemberSignInConfig, _ int) resp.AppMemberSignInConfigResp {
		return resp.AppMemberSignInConfigResp{
			Day:   item.Day,
			Point: item.Point,
		}
	})
	response.WriteSuccess(c, respList)
}
