package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var r req.UserSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateUser(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var r req.UserSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateUser(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteUser(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	item, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(item))
}

func (h *UserHandler) GetUserPage(c *gin.Context) {
	var r req.UserPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	page, err := h.svc.GetUserPage(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(page))
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	var r req.UserUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateUserStatus(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *UserHandler) ResetUserPassword(c *gin.Context) {
	var r req.UserResetPasswordReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.ResetUserPassword(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
	var r req.UserUpdatePasswordReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	// Note: This API typically checks old password, but Admin reset usually doesn't.
	// Admin changing other's password vs User changing own password.
	// This handler seems to be for Admin (UpdateUserPassword) or User Profile?
	// Based on Java Controller, there is usually /system/user/update-password (Profile) and /system/user/profile/update-password.
	// Checked Java Controller:
	// @PutMapping("update-password") Admin updates users password.
	if err := h.svc.UpdateUserPassword(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}
