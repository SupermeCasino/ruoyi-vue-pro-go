package promotion

import (
	"strconv"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// AppCouponTemplateHandler App 端优惠券模板 Handler (对齐 Java: AppCouponTemplateController)
type AppCouponTemplateHandler struct {
	svc *promotion.CouponService
}

func NewAppCouponTemplateHandler(svc *promotion.CouponService) *AppCouponTemplateHandler {
	return &AppCouponTemplateHandler{svc: svc}
}

// GetCouponTemplate 获得优惠劵模版 (对齐 Java: getCouponTemplate)
// @PermitAll - 公开接口
func (h *AppCouponTemplateHandler) GetCouponTemplate(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	// 获取当前用户 ID (可能未登录)
	userId := context.GetUserId(c)

	template, err := h.svc.GetCouponTemplateForApp(c, id, userId)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, template)
}

// GetCouponTemplateList 获得优惠劵模版列表 (对齐 Java: getCouponTemplateList - 带查询条件)
// @PermitAll - 公开接口
func (h *AppCouponTemplateHandler) GetCouponTemplateList(c *gin.Context) {
	var r req.AppCouponTemplateListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	userId := context.GetUserId(c)

	list, err := h.svc.GetCouponTemplateListForApp(c, r.SpuID, r.ProductScope, r.Count, userId)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, list)
}

// GetCouponTemplateListByIds 获得优惠劵模版列表 (按ID) (对齐 Java: getCouponTemplateList - 按ids)
// @PermitAll - 公开接口
func (h *AppCouponTemplateHandler) GetCouponTemplateListByIds(c *gin.Context) {
	idsStr := c.Query("ids")
	var ids []int64
	if idsStr != "" {
		for _, s := range strings.Split(idsStr, ",") {
			if id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
				ids = append(ids, id)
			}
		}
	}

	userId := context.GetUserId(c)

	list, err := h.svc.GetCouponTemplateListByIdsForApp(c, ids, userId)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, list)
}

// GetCouponTemplatePage 获得优惠劵模版分页 (对齐 Java: getCouponTemplatePage)
// @PermitAll - 公开接口
func (h *AppCouponTemplateHandler) GetCouponTemplatePage(c *gin.Context) {
	var r req.AppCouponTemplatePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	userId := context.GetUserId(c)

	page, err := h.svc.GetCouponTemplatePageForApp(c, &r, userId)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, page)
}
