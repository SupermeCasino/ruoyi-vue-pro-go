package promotion

import (
	"strconv"

	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type RewardActivityHandler struct {
	svc *promotion.RewardActivityService
}

func NewRewardActivityHandler(svc *promotion.RewardActivityService) *RewardActivityHandler {
	return &RewardActivityHandler{svc: svc}
}

// CreateRewardActivity 创建活动
// @Summary 创建活动
// @Router /admin-api/promotion/reward-activity/create [post]
func (h *RewardActivityHandler) CreateRewardActivity(c *gin.Context) {
	var r promotion2.PromotionRewardActivityCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateRewardActivity(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateRewardActivity 更新活动
// @Summary 更新活动
// @Router /admin-api/promotion/reward-activity/update [put]
func (h *RewardActivityHandler) UpdateRewardActivity(c *gin.Context) {
	var r promotion2.PromotionRewardActivityUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateRewardActivity(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteRewardActivity 删除活动
// @Summary 删除活动
// @Router /admin-api/promotion/reward-activity/delete [delete]
func (h *RewardActivityHandler) DeleteRewardActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteRewardActivity(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// CloseRewardActivity 关闭活动
// @Summary 关闭活动
// @Router /admin-api/promotion/reward-activity/close [put]
// Java: RewardActivityController#closeRewardActivity
func (h *RewardActivityHandler) CloseRewardActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.CloseRewardActivity(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetRewardActivity 获得活动
// @Summary 获得活动
// @Router /admin-api/promotion/reward-activity/get [get]
func (h *RewardActivityHandler) GetRewardActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetRewardActivity(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetRewardActivityPage 获得活动分页
// @Summary 获得活动分页
// @Router /admin-api/promotion/reward-activity/page [get]
func (h *RewardActivityHandler) GetRewardActivityPage(c *gin.Context) {
	var r promotion2.PromotionRewardActivityPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetRewardActivityPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
