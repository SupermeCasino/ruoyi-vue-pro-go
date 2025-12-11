package member

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	memberModel "backend-go/internal/model/member"
	"backend-go/internal/pkg/core"
	memberSvc "backend-go/internal/service/member"

	"github.com/gin-gonic/gin"
)

type MemberConfigHandler struct {
	svc *memberSvc.MemberConfigService
}

func NewMemberConfigHandler(svc *memberSvc.MemberConfigService) *MemberConfigHandler {
	return &MemberConfigHandler{svc: svc}
}

// SaveConfig 保存会员配置
func (h *MemberConfigHandler) SaveConfig(c *gin.Context) {
	var r req.MemberConfigSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	err := h.svc.SaveConfig(c, &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetConfig 获得会员配置
func (h *MemberConfigHandler) GetConfig(c *gin.Context) {
	config, err := h.svc.GetConfig(c)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, h.convertResp(config))
}

func (h *MemberConfigHandler) convertResp(item *memberModel.MemberConfig) *resp.MemberConfigResp {
	if item == nil {
		return nil
	}
	return &resp.MemberConfigResp{
		ID:                        item.ID,
		PointTradeDeductEnable:    item.PointTradeDeductEnable,
		PointTradeDeductUnitPrice: item.PointTradeDeductUnitPrice,
		PointTradeDeductMaxPrice:  item.PointTradeDeductMaxPrice,
		PointTradeGivePoint:       item.PointTradeGivePoint,
	}
}
