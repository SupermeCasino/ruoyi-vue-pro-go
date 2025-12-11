package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"
)

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{
		svc: svc,
	}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var r req.RoleSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateRole(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var r req.RoleSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateRole(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *RoleHandler) UpdateRoleStatus(c *gin.Context) {
	var r req.RoleUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateRoleStatus(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteRole(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	item, err := h.svc.GetRole(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(item))
}

func (h *RoleHandler) GetRolePage(c *gin.Context) {
	var r req.RolePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	page, err := h.svc.GetRolePage(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(page))
}

func (h *RoleHandler) GetSimpleRoleList(c *gin.Context) {
	// Status 0: Enable
	list, err := h.svc.GetRoleListByStatus(c.Request.Context(), 0)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}
