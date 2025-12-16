package product

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
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
// @Summary 添加商品收藏
// @Tags 用户 APP - 商品收藏
// @Produce json
// @Param req body req.AppFavoriteCreateReq true "请求参数"
// @Router /app-api/product/favorite/create [post]
func (h *AppProductFavoriteHandler) CreateFavorite(c *gin.Context) {
	var r req.AppFavoriteCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	id, err := h.svc.CreateFavorite(c, userId, r.SpuId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(id))
}

// DeleteFavorite 取消单个商品收藏
// @Summary 取消单个商品收藏
// @Tags 用户 APP - 商品收藏
// @Produce json
// @Param req body req.AppFavoriteReq true "请求参数"
// @Router /app-api/product/favorite/delete [delete]
func (h *AppProductFavoriteHandler) DeleteFavorite(c *gin.Context) {
	var r req.AppFavoriteReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	if err := h.svc.DeleteFavorite(c, userId, r.SpuId); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// GetFavoritePage 获得商品收藏分页
// @Summary 获得商品收藏分页
// @Tags 用户 APP - 商品收藏
// @Produce json
// @Param pageNo query int true "页码"
// @Param pageSize query int true "页数"
// @Success 200 {object} pagination.PageResult[resp.AppFavoriteResp]
// @Router /app-api/product/favorite/page [get]
func (h *AppProductFavoriteHandler) GetFavoritePage(c *gin.Context) {
	var r req.AppFavoritePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	res, err := h.svc.GetAppFavoritePage(c, userId, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(res))
}

// IsFavoriteExists 检查是否收藏过商品
// @Summary 检查是否收藏过商品
// @Tags 用户 APP - 商品收藏
// @Produce json
// @Param spuId query int true "商品SPU编号"
// @Success 200 {object} bool
// @Router /app-api/product/favorite/exits [get]
func (h *AppProductFavoriteHandler) IsFavoriteExists(c *gin.Context) {
	var r req.AppFavoriteReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	favorite, err := h.svc.GetFavorite(c, userId, r.SpuId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(200, response.Success(false))
			return
		}
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(favorite != nil))
}

// GetFavoriteCount 获得商品收藏数量
// @Summary 获得商品收藏数量
// @Tags 用户 APP - 商品收藏
// @Produce json
// @Success 200 {object} int64
// @Router /app-api/product/favorite/get-count [get]
func (h *AppProductFavoriteHandler) GetFavoriteCount(c *gin.Context) {
	userId := context.GetLoginUserID(c)
	count, err := h.svc.GetFavoriteCount(c, userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(count))
}
