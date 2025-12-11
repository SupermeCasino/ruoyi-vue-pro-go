package product

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/product"

	"github.com/gin-gonic/gin"
)

type ProductBrowseHistoryHandler struct {
	svc *product.ProductBrowseHistoryService
}

func NewProductBrowseHistoryHandler(svc *product.ProductBrowseHistoryService) *ProductBrowseHistoryHandler {
	return &ProductBrowseHistoryHandler{svc: svc}
}

// GetBrowseHistoryPage 获得商品浏览记录分页 (Admin)
// @Summary 获得商品浏览记录分页
// @Tags 管理后台-商品浏览记录
// @Produce json
// @Param pageNo query int true "页码"
// @Param pageSize query int true "页数"
// @Success 200 {object} core.PageResult[resp.ProductBrowseHistoryResp]
// @Router /admin-api/product/browse-history/page [get]
func (h *ProductBrowseHistoryHandler) GetBrowseHistoryPage(c *gin.Context) {
	var r req.ProductBrowseHistoryPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetBrowseHistoryPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}
