package product

import (
	product2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
// @Success 200 {object} pagination.PageResult[resp.ProductBrowseHistoryResp]
// @Router /admin-api/product/browse-history/page [get]
func (h *ProductBrowseHistoryHandler) GetBrowseHistoryPage(c *gin.Context) {
	var r product2.ProductBrowseHistoryPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetBrowseHistoryPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
