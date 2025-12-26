package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	memberModel "github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.SaveConfig(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetConfig 获得会员配置
func (h *MemberConfigHandler) GetConfig(c *gin.Context) {
	config, err := h.svc.GetConfig(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, h.convertResp(config))
}

func (h *MemberConfigHandler) convertResp(item *memberModel.MemberConfig) *resp.MemberConfigResp {
	if item == nil {
		return nil
	}
	pointTradeDeductEnable := 0
	if item.PointTradeDeductEnable {
		pointTradeDeductEnable = 1
	}
	return &resp.MemberConfigResp{
		ID:                        item.ID,
		PointTradeDeductEnable:    pointTradeDeductEnable,
		PointTradeDeductUnitPrice: item.PointTradeDeductUnitPrice,
		PointTradeDeductMaxPrice:  item.PointTradeDeductMaxPrice,
		PointTradeGivePoint:       item.PointTradeGivePoint,
	}
}
