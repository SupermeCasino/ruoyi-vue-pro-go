package product

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/product"

	"github.com/gin-gonic/gin"
)

type ProductCommentHandler struct {
	svc *product.ProductCommentService
}

func NewProductCommentHandler(svc *product.ProductCommentService) *ProductCommentHandler {
	return &ProductCommentHandler{svc: svc}
}

// GetCommentPage 获得商品评价分页 (Admin)
// @Summary 获得商品评价分页
// @Tags 管理后台-商品评价
// @Produce json
// @Param pageNo query int true "页码"
// @Param pageSize query int true "页数"
// @Success 200 {object} core.PageResult[resp.ProductCommentResp]
// @Router /admin-api/product/comment/page [get]
func (h *ProductCommentHandler) GetCommentPage(c *gin.Context) {
	var req req.ProductCommentPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetCommentPage(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// UpdateCommentVisible 更新商品评价可见性
// @Summary 更新商品评价可见性
// @Tags 管理后台-商品评价
// @Produce json
// @Param req body req.ProductCommentUpdateVisibleReq true "请求参数"
// @Router /admin-api/product/comment/update-visible [put]
func (h *ProductCommentHandler) UpdateCommentVisible(c *gin.Context) {
	var req req.ProductCommentUpdateVisibleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateCommentVisible(c, &req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// ReplyComment 回复商品评价
// @Summary 回复商品评价
// @Tags 管理后台-商品评价
// @Produce json
// @Param req body req.ProductCommentReplyReq true "请求参数"
// @Router /admin-api/product/comment/reply [put]
func (h *ProductCommentHandler) ReplyComment(c *gin.Context) {
	var req req.ProductCommentReplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	userId := core.GetLoginUserID(c)
	if err := h.svc.ReplyComment(c, &req, userId); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// CreateComment 添加自评
// @Summary 添加自评
// @Tags 管理后台-商品评价
// @Produce json
// @Param req body req.ProductCommentCreateReq true "请求参数"
// @Router /admin-api/product/comment/create [post]
func (h *ProductCommentHandler) CreateComment(c *gin.Context) {
	var req req.ProductCommentCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.CreateComment(c, &req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}
