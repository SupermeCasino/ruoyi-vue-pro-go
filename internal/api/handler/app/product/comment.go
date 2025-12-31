package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppProductCommentHandler struct {
	svc *product.ProductCommentService
}

func NewAppProductCommentHandler(svc *product.ProductCommentService) *AppProductCommentHandler {
	return &AppProductCommentHandler{svc: svc}
}

// GetCommentPage 获得商品评价分页 (App)
func (h *AppProductCommentHandler) GetCommentPage(c *gin.Context) {
	var r req.AppProductCommentPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetAppCommentPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// CreateComment 创建商品评价
func (h *AppProductCommentHandler) CreateComment(c *gin.Context) {
	var r req.AppProductCommentCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.CreateAppComment(c, context.GetUserId(c), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res.ID)
}
