package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductFavoriteHandler struct {
	svc *product.ProductFavoriteService
}

func NewProductFavoriteHandler(svc *product.ProductFavoriteService) *ProductFavoriteHandler {
	return &ProductFavoriteHandler{svc: svc}
}

// GetFavoritePage 获得商品收藏分页 (Admin)
// @Summary 获得商品收藏分页
// @Tags 管理后台-商品收藏
// @Produce json
// @Param pageNo query int true "页码"
// @Param pageSize query int true "页数"
// @Success 200 {object} pagination.PageResult[resp.ProductFavoriteResp]
// @Router /admin-api/product/favorite/page [get]
func (h *ProductFavoriteHandler) GetFavoritePage(c *gin.Context) {
	var r req.ProductFavoritePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetFavoritePage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
