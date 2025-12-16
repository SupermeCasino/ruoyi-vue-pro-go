package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
)

type DeptHandler struct {
	svc *service.DeptService
}

func NewDeptHandler(svc *service.DeptService) *DeptHandler {
	return &DeptHandler{
		svc: svc,
	}
}

func (h *DeptHandler) CreateDept(c *gin.Context) {
	var r req.DeptSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateDept(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

func (h *DeptHandler) UpdateDept(c *gin.Context) {
	var r req.DeptSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdateDept(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *DeptHandler) DeleteDept(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteDept(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *DeptHandler) GetDept(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	item, err := h.svc.GetDept(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(item))
}

func (h *DeptHandler) GetDeptList(c *gin.Context) {
	var r req.DeptListReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	list, err := h.svc.GetDeptList(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}

func (h *DeptHandler) GetSimpleDeptList(c *gin.Context) {
	list, err := h.svc.GetSimpleDeptList(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}
