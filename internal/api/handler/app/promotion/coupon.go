package promotion

import (
	"strconv"

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

// TakeCoupon 领取优惠券 (对齐 Java: AppCouponController.takeCoupon)
// 返回 Boolean: 是否可继续领取
func (h *AppCouponHandler) TakeCoupon(c *gin.Context) {
	var r req.AppCouponTakeReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	uid := c.GetInt64(context.CtxUserIDKey)
	canTakeAgain, err := h.svc.TakeCoupon(c, uid, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, canTakeAgain)
}

// GetCouponPage 我的优惠券 (对齐 Java: AppCouponController.getCouponPage)
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

// GetCoupon 获得优惠劵 (对齐 Java: AppCouponController.getCoupon)
func (h *AppCouponHandler) GetCoupon(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}
	uid := c.GetInt64(context.CtxUserIDKey)
	coupon, err := h.svc.GetCoupon(c, uid, id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, coupon)
}

// GetUnusedCouponCount 获得未使用的优惠劵数量 (对齐 Java: AppCouponController.getUnusedCouponCount)
func (h *AppCouponHandler) GetUnusedCouponCount(c *gin.Context) {
	uid := c.GetInt64(context.CtxUserIDKey)
	count, err := h.svc.GetUnusedCouponCount(c, uid)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, count)
}

// GetCouponMatchList 获得匹配的优惠券列表
func (h *AppCouponHandler) GetCouponMatchList(c *gin.Context) {
	var r req.AppCouponMatchReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}
	uid := c.GetInt64(context.CtxUserIDKey)

	list, err := h.svc.GetCouponMatchList(c, uid, r.Price, r.SpuIDs, r.CategoryIDs)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, list)
}
