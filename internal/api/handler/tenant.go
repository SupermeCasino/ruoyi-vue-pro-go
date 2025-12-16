package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"

	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type TenantHandler struct {
	svc *service.TenantService
}

func NewTenantHandler(svc *service.TenantService) *TenantHandler {
	return &TenantHandler{svc: svc}
}

// GetTenantSimpleList 获取租户精简列表
// @Router /system/tenant/simple-list [get]
func (h *TenantHandler) GetTenantSimpleList(c *gin.Context) {
	list, err := h.svc.GetTenantSimpleList(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(list))
}

// CreateTenant 创建租户
// @Router /system/tenant/create [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var r req.TenantCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateTenant(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(id))
}

// UpdateTenant 更新租户
// @Router /system/tenant/update [put]
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	var r req.TenantUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateTenant(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// DeleteTenant 删除租户
// @Router /system/tenant/delete [delete]
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteTenant(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// GetTenant 获得租户
// @Router /system/tenant/get [get]
func (h *TenantHandler) GetTenant(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	item, err := h.svc.GetTenant(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(item))
}

// GetTenantPage 获得租户分页
// @Router /system/tenant/page [get]
func (h *TenantHandler) GetTenantPage(c *gin.Context) {
	var r req.TenantPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	page, err := h.svc.GetTenantPage(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(page))
}

// ExportTenantExcel 导出租户 Excel
// @Router /system/tenant/export-excel [get]
func (h *TenantHandler) ExportTenantExcel(c *gin.Context) {
	var r req.TenantExportReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	list, err := h.svc.GetTenantList(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}

	// Create Excel File
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			c.Error(err)
			return
		}
	}()

	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		c.Error(err)
		return
	}
	f.SetActiveSheet(index)

	// Headers
	headers := []string{"租户编号", "租户名", "联系人", "联系手机", "状态", "绑定域名", "过期时间", "账号数量", "创建时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Data
	for i, item := range list {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.ContactName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.ContactMobile)
		statusStr := "开启"
		if item.Status != 0 {
			statusStr = "关闭"
		}
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), statusStr)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), item.Domain)

		expireStr := ""
		if item.ExpireDate > 0 {
			expireStr = time.UnixMilli(item.ExpireDate).Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), expireStr)

		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), item.AccountCount)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), item.CreateTime.Format("2006-01-02 15:04:05"))
	}

	// Response
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=tenant_list.xlsx")
	if err := f.Write(c.Writer); err != nil {
		c.Error(err)
		return
	}
}

// GetTenantByWebsite 根据域名获取租户
// @Router /system/tenant/get-by-website [get]
func (h *TenantHandler) GetTenantByWebsite(c *gin.Context) {
	website := c.Query("website")
	tenant, err := h.svc.GetTenantByWebsite(c.Request.Context(), website)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(tenant))
}

// GetTenantIdByName 根据租户名获取租户ID
// @Router /system/tenant/get-id-by-name [get]
func (h *TenantHandler) GetTenantIdByName(c *gin.Context) {
	name := c.Query("name")
	tenantId, err := h.svc.GetTenantIdByName(c.Request.Context(), name)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(tenantId))
}
