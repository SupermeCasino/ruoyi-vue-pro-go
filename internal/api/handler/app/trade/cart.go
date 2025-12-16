package trade

import (
	"strconv"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppCartHandler struct {
	svc *trade.CartService
}

func NewAppCartHandler(svc *trade.CartService) *AppCartHandler {
	return &AppCartHandler{svc: svc}
}

// AddCart 添加购物车
func (h *AppCartHandler) AddCart(c *gin.Context) {
	var r req.AppCartAddReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	id, err := h.svc.AddCart(c, userId, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(id))
}

// UpdateCartCount 更新购物车数量
func (h *AppCartHandler) UpdateCartCount(c *gin.Context) {
	var r req.AppCartUpdateCountReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	if err := h.svc.UpdateCartCount(c, userId, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// UpdateCartSelected 更新购物车选中状态
func (h *AppCartHandler) UpdateCartSelected(c *gin.Context) {
	var r req.AppCartUpdateSelectedReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	if err := h.svc.UpdateCartSelected(c, userId, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// ResetCart 重置购物车
func (h *AppCartHandler) ResetCart(c *gin.Context) {
	var r req.AppCartResetReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	if err := h.svc.ResetCart(c, userId, &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// DeleteCart 删除购物车
func (h *AppCartHandler) DeleteCart(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		c.JSON(200, errors.ErrParam)
		return
	}
	idStrs := strings.Split(idsStr, ",")
	var ids []int64
	for _, s := range idStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		c.JSON(200, errors.ErrParam)
		return
	}
	userId := context.GetLoginUserID(c)
	if err := h.svc.DeleteCart(c, userId, ids); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// GetCartCount 获取购物车商品数量
func (h *AppCartHandler) GetCartCount(c *gin.Context) {
	userId := context.GetLoginUserID(c)
	count, err := h.svc.GetCartCount(c, userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(count))
}

// GetCartList 获取购物车列表
func (h *AppCartHandler) GetCartList(c *gin.Context) {
	userId := context.GetLoginUserID(c)
	res, err := h.svc.GetCartList(c, userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(res))
}
