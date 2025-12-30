package member

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateAddress(c, c.GetInt64(context.CtxUserIDKey), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateAddress 更新收件地址
// @Router /member/address/update [put]
func (h *AppMemberAddressHandler) UpdateAddress(c *gin.Context) {
	var r req.AppAddressUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateAddress(c, c.GetInt64(context.CtxUserIDKey), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteAddress 删除收件地址
// @Router /member/address/delete [delete]
func (h *AppMemberAddressHandler) DeleteAddress(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteAddress(c, c.GetInt64(context.CtxUserIDKey), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetAddress 获得收件地址
// @Router /member/address/get [get]
func (h *AppMemberAddressHandler) GetAddress(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetAddress(c, c.GetInt64(context.CtxUserIDKey), id)
	if err != nil {
			response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetDefaultUserAddress 获得默认收件地址
// @Router /member/address/get-default [get]
func (h *AppMemberAddressHandler) GetDefaultUserAddress(c *gin.Context) {
	res, err := h.svc.GetDefaultUserAddress(c, c.GetInt64(context.CtxUserIDKey))
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetAddressList 获得收件地址列表
// @Router /member/address/list [get]
func (h *AppMemberAddressHandler) GetAddressList(c *gin.Context) {
	res, err := h.svc.GetAddressList(c, c.GetInt64(context.CtxUserIDKey))
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
