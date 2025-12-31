package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type DictHandler struct {
	svc *system.DictService
}

func NewDictHandler(svc *system.DictService) *DictHandler {
	return &DictHandler{
		svc: svc,
	}
}

// --- DictType ---

func (h *DictHandler) CreateDictType(c *gin.Context) {
	var r req.DictTypeSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateDictType(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *DictHandler) UpdateDictType(c *gin.Context) {
	var r req.DictTypeSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateDictType(c.Request.Context(), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DictHandler) DeleteDictType(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteDictType(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DictHandler) GetDictType(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	item, err := h.svc.GetDictType(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, item)
}

func (h *DictHandler) GetDictTypePage(c *gin.Context) {
	var r req.DictTypePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetDictTypePage(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}

func (h *DictHandler) GetSimpleDictTypeList(c *gin.Context) {
	list, err := h.svc.GetSimpleDictTypeList(c.Request.Context())
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}

// ExportDictTypeExcel 导出字典类型 Excel
func (h *DictHandler) ExportDictTypeExcel(c *gin.Context) {
	var r req.DictTypePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// TODO: 实现 Excel 导出逻辑
	// 这里简化实现，返回空数据
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=dict-type.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", []byte{})
}

// --- DictData ---

func (h *DictHandler) CreateDictData(c *gin.Context) {
	var r req.DictDataSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateDictData(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *DictHandler) UpdateDictData(c *gin.Context) {
	var r req.DictDataSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateDictData(c.Request.Context(), &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DictHandler) DeleteDictData(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteDictData(c.Request.Context(), id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DictHandler) GetDictData(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	item, err := h.svc.GetDictData(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, item)
}

func (h *DictHandler) GetDictDataPage(c *gin.Context) {
	var r req.DictDataPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetDictDataPage(c.Request.Context(), &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}

func (h *DictHandler) GetSimpleDictDataList(c *gin.Context) {
	list, err := h.svc.GetSimpleDictDataList(c.Request.Context())
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}
