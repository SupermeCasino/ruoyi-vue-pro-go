package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppCouponHandler struct {
	svc *promotion.CouponUserService
}

func NewAppCouponHandler(svc *promotion.CouponUserService) *AppCouponHandler {
	return &AppCouponHandler{svc: svc}
}

// TakeCoupon 领取优惠券
func (h *AppCouponHandler) TakeCoupon(c *gin.Context) {
	var r req.AppCouponTakeReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	uid := c.GetInt64(context.CtxUserIDKey)
	id, err := h.svc.TakeCoupon(c, uid, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, id)
}

// GetCouponPage 我的优惠券
func (h *AppCouponHandler) GetCouponPage(c *gin.Context) {
	var r req.AppCouponPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	uid := c.GetInt64(context.CtxUserIDKey)
	list, err := h.svc.GetCouponPage(c, uid, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, list)
}

// GetCouponMatchList 获得匹配的优惠券列表
func (h *AppCouponHandler) GetCouponMatchList(c *gin.Context) {
	var r req.AppCouponMatchReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	uid := c.GetInt64(context.CtxUserIDKey)

	// SpuIDs and CategoryIDs need to be resolved via Product Service usually.
	// But AppCouponMatchReq takes them as IDs.
	// If the frontend sends them, we use them.
	// NOTE: Secure implementation would verify these with Product Service, but we trust frontend for query/match context for now.

	list, err := h.svc.GetCouponMatchList(c, uid, r.Price, r.SpuIDs, r.CategoryIDs)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, list)
}
