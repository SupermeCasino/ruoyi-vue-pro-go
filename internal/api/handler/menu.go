package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type MenuHandler struct {
	svc *system.MenuService
}

func NewMenuHandler(svc *system.MenuService) *MenuHandler {
	return &MenuHandler{
		svc: svc,
	}
}

// CreateMenu 创建菜单
// @Router /system/menu/create [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var r req.MenuCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateMenu(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateMenu 更新菜单
// @Router /system/menu/update [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	var r req.MenuUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateMenu(c.Request.Context(), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteMenu 删除菜单
// @Router /system/menu/delete [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := h.svc.DeleteMenu(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetMenuList 获取菜单列表
// @Router /system/menu/list [get]
func (h *MenuHandler) GetMenuList(c *gin.Context) {
	var r req.MenuListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	list, err := h.svc.GetMenuList(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}

// GetMenu 获取菜单详情
// @Router /system/menu/get [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	item, err := h.svc.GetMenu(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, item)
}

// GetSimpleMenuList 获取精简菜单列表
// @Router /system/menu/simple-list [get]
func (h *MenuHandler) GetSimpleMenuList(c *gin.Context) {
	list, err := h.svc.GetSimpleMenuList(c.Request.Context())
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}
