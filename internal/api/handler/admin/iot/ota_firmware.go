package iot

import (
	"strconv"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

// Create 创建固件
func (h *OtaFirmwareHandler) Create(c *gin.Context) {
	var r iot2.IotOtaFirmwareSaveReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	id, err := h.svc.Create(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// Update 更新固件
func (h *OtaFirmwareHandler) Update(c *gin.Context) {
	var r iot2.IotOtaFirmwareSaveReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.BindingErr(err))
		return
	}
	if err := h.svc.Update(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Delete 删除固件
func (h *OtaFirmwareHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.Delete(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Get 获取固件
func (h *OtaFirmwareHandler) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	firmware, err := h.svc.Get(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if firmware == nil {
		response.WriteSuccess(c, nil)
		return
	}
	resp := &iot2.IotOtaFirmwareRespVO{
		ID:                  firmware.ID,
		Name:                firmware.Name,
		Description:         firmware.Description,
		Version:             firmware.Version,
		ProductID:           firmware.ProductID,
		FileURL:             firmware.FileURL,
		FileSize:            firmware.FileSize,
		FileDigestAlgorithm: firmware.FileDigestAlgorithm,
		FileDigestValue:     firmware.FileDigestValue,
		CreateTime:          firmware.CreateTime,
	}
	// 补全产品名称
	if p, _ := h.productSvc.Get(c, firmware.ProductID); p != nil {
		resp.ProductName = p.Name
	}
	response.WriteSuccess(c, resp)
}

// Page 获取固件分页
func (h *OtaFirmwareHandler) Page(c *gin.Context) {
	var r iot2.IotOtaFirmwarePageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]*iot2.IotOtaFirmwareRespVO, 0, len(page.List))
	for _, item := range page.List {
		resp := &iot2.IotOtaFirmwareRespVO{
			ID:                  item.ID,
			Name:                item.Name,
			Description:         item.Description,
			Version:             item.Version,
			ProductID:           item.ProductID,
			FileURL:             item.FileURL,
			FileSize:            item.FileSize,
			FileDigestAlgorithm: item.FileDigestAlgorithm,
			FileDigestValue:     item.FileDigestValue,
			CreateTime:          item.CreateTime,
		}
		// TOOD: 补全产品名称 (实际开发中建议批量获取以优化性能)
		if p, _ := h.productSvc.Get(c, item.ProductID); p != nil {
			resp.ProductName = p.Name
		}
		list = append(list, resp)
	}
	response.WritePage(c, page.Total, list)
}
