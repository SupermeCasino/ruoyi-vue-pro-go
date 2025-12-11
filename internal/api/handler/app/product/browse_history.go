package product

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/product"

	"github.com/gin-gonic/gin"
)

type AppProductBrowseHistoryHandler struct {
	svc *product.ProductBrowseHistoryService
}

func NewAppProductBrowseHistoryHandler(svc *product.ProductBrowseHistoryService) *AppProductBrowseHistoryHandler {
	return &AppProductBrowseHistoryHandler{svc: svc}
}

// DeleteBrowseHistory 删除商品浏览记录
// @Summary 删除商品浏览记录
// @Tags 用户 APP - 商品浏览记录
// @Produce json
// @Param req body req.AppProductBrowseHistoryDeleteReq true "请求参数"
// @Router /app-api/product/browse-history/delete [delete]
func (h *AppProductBrowseHistoryHandler) DeleteBrowseHistory(c *gin.Context) {
	var r req.AppProductBrowseHistoryDeleteReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	userId := core.GetLoginUserID(c)
	if err := h.svc.HideUserBrowseHistory(c, userId, r.SpuIds); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// CleanBrowseHistory 清空商品浏览记录
// @Summary 清空商品浏览记录
// @Tags 用户 APP - 商品浏览记录
// @Produce json
// @Router /app-api/product/browse-history/clean [delete]
func (h *AppProductBrowseHistoryHandler) CleanBrowseHistory(c *gin.Context) {
	userId := core.GetLoginUserID(c)
	if err := h.svc.HideUserBrowseHistory(c, userId, nil); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// GetBrowseHistoryPage 获得商品浏览记录分页
// @Summary 获得商品浏览记录分页
// @Tags 用户 APP - 商品浏览记录
// @Produce json
// @Param pageNo query int true "页码"
// @Param pageSize query int true "页数"
// @Success 200 {object} core.PageResult[resp.AppProductBrowseHistoryResp]
// @Router /app-api/product/browse-history/page [get]
func (h *AppProductBrowseHistoryHandler) GetBrowseHistoryPage(c *gin.Context) {
	var r req.AppProductBrowseHistoryPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	userId := core.GetLoginUserID(c)
	res, err := h.svc.GetAppBrowseHistoryPage(c, userId, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}
