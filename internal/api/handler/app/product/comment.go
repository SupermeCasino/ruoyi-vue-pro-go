package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"

	"github.com/gin-gonic/gin"
)

type AppProductCommentHandler struct {
	svc *product.ProductCommentService
}

func NewAppProductCommentHandler(svc *product.ProductCommentService) *AppProductCommentHandler {
	return &AppProductCommentHandler{svc: svc}
}

// GetCommentPage 获得商品评价分页 (App)
// @Summary 获得商品评价分页
// @Tags 用户 APP - 商品评价
// @Produce json
// @Param spuId query int false "商品SPU编号"
// @Param pageNo query int true "页码"
// @Param pageSize query int true "页数"
// @Success 200 {object} core.PageResult[resp.AppProductCommentResp]
// @Router /app-api/product/comment/page [get]
func (h *AppProductCommentHandler) GetCommentPage(c *gin.Context) {
	var r req.AppProductCommentPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetAppCommentPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// CreateComment 创建商品评价
func (h *AppProductCommentHandler) CreateComment(c *gin.Context) {
	var r req.AppProductCommentCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.CreateAppComment(c, core.GetUserId(c), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res.ID))
}
