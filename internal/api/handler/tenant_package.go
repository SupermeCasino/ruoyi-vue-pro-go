package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
)

type TenantPackageHandler struct {
	svc *service.TenantPackageService
}

func NewTenantPackageHandler(svc *service.TenantPackageService) *TenantPackageHandler {
	return &TenantPackageHandler{svc: svc}
}

// CreateTenantPackage 创建租户套餐
func (h *TenantPackageHandler) CreateTenantPackage(c *gin.Context) {
	var r req.TenantPackageSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateTenantPackage(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateTenantPackage 更新租户套餐
func (h *TenantPackageHandler) UpdateTenantPackage(c *gin.Context) {
	var r req.TenantPackageSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateTenantPackage(c.Request.Context(), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteTenantPackage 删除租户套餐
func (h *TenantPackageHandler) DeleteTenantPackage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteTenantPackage(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteTenantPackageList 批量删除租户套餐
func (h *TenantPackageHandler) DeleteTenantPackageList(c *gin.Context) {
	ids := utils.ParseIDs(c.QueryArray("ids"))
	if len(ids) == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteTenantPackageList(c.Request.Context(), ids); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetTenantPackage 获得租户套餐
func (h *TenantPackageHandler) GetTenantPackage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	result, err := h.svc.GetTenantPackage(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, result)
}

// GetTenantPackageSimpleList 获得租户套餐精简列表
func (h *TenantPackageHandler) GetTenantPackageSimpleList(c *gin.Context) {
	result, err := h.svc.GetTenantPackageSimpleList(c.Request.Context())
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, result)
}

// GetTenantPackagePage 获得租户套餐分页
func (h *TenantPackageHandler) GetTenantPackagePage(c *gin.Context) {
	var r req.TenantPackagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	result, err := h.svc.GetTenantPackagePage(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, result)
}
