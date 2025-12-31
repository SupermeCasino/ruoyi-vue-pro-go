package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppProductFavoriteHandler struct {
	svc *product.ProductFavoriteService
}

func NewAppProductFavoriteHandler(svc *product.ProductFavoriteService) *AppProductFavoriteHandler {
	return &AppProductFavoriteHandler{svc: svc}
}

// CreateFavorite 添加商品收藏
func (h *AppProductFavoriteHandler) CreateFavorite(c *gin.Context) {
	var r req.AppFavoriteCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	id, err := h.svc.CreateFavorite(c, userId, r.SpuId)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// DeleteFavorite 取消单个商品收藏
func (h *AppProductFavoriteHandler) DeleteFavorite(c *gin.Context) {
	var r req.AppFavoriteReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	if err := h.svc.DeleteFavorite(c, userId, r.SpuId); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetFavoritePage 获得商品收藏分页
func (h *AppProductFavoriteHandler) GetFavoritePage(c *gin.Context) {
	var r req.AppFavoritePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	res, err := h.svc.GetAppFavoritePage(c, userId, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// IsFavoriteExists 检查是否收藏过商品
func (h *AppProductFavoriteHandler) IsFavoriteExists(c *gin.Context) {
	var r req.AppFavoriteReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	favorite, err := h.svc.GetFavorite(c, userId, r.SpuId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.WriteSuccess(c, false)
			return
		}
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, favorite != nil)
}

// GetFavoriteCount 获得商品收藏数量
func (h *AppProductFavoriteHandler) GetFavoriteCount(c *gin.Context) {
	userId := context.GetLoginUserID(c)
	count, err := h.svc.GetFavoriteCount(c, userId)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, count)
}
