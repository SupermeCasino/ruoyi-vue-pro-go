package member

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppMemberAddressHandler struct {
	svc *member.MemberAddressService
}

func NewAppMemberAddressHandler(svc *member.MemberAddressService) *AppMemberAddressHandler {
	return &AppMemberAddressHandler{svc: svc}
}

// CreateAddress 创建收件地址
// @Router /member/address/create [post]
func (h *AppMemberAddressHandler) CreateAddress(c *gin.Context) {
	var r req.AppAddressCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateAddress(c, c.GetInt64(core.CtxUserIDKey), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

// UpdateAddress 更新收件地址
// @Router /member/address/update [put]
func (h *AppMemberAddressHandler) UpdateAddress(c *gin.Context) {
	var r req.AppAddressUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateAddress(c, c.GetInt64(core.CtxUserIDKey), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// DeleteAddress 删除收件地址
// @Router /member/address/delete [delete]
func (h *AppMemberAddressHandler) DeleteAddress(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.DeleteAddress(c, c.GetInt64(core.CtxUserIDKey), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// GetAddress 获得收件地址
// @Router /member/address/get [get]
func (h *AppMemberAddressHandler) GetAddress(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetAddress(c, c.GetInt64(core.CtxUserIDKey), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetDefaultUserAddress 获得默认收件地址
// @Router /member/address/get-default [get]
func (h *AppMemberAddressHandler) GetDefaultUserAddress(c *gin.Context) {
	res, err := h.svc.GetDefaultUserAddress(c, c.GetInt64(core.CtxUserIDKey))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetAddressList 获得收件地址列表
// @Router /member/address/list [get]
func (h *AppMemberAddressHandler) GetAddressList(c *gin.Context) {
	res, err := h.svc.GetAddressList(c, c.GetInt64(core.CtxUserIDKey))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}
