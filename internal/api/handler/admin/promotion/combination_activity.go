package promotion

import (
	"strconv"

	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/promotion"

	"github.com/gin-gonic/gin"
)

type CombinationActivityHandler struct {
	svc promotion.CombinationActivityService
}

func NewCombinationActivityHandler(svc promotion.CombinationActivityService) *CombinationActivityHandler {
	return &CombinationActivityHandler{svc: svc}
}

// CreateCombinationActivity 创建拼团活动
func (h *CombinationActivityHandler) CreateCombinationActivity(c *gin.Context) {
	var r req.CombinationActivityCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.NewBizError(400, "Invalid Request"))
		return
	}

	id, err := h.svc.CreateCombinationActivity(c.Request.Context(), r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateCombinationActivity 更新拼团活动
func (h *CombinationActivityHandler) UpdateCombinationActivity(c *gin.Context) {
	var r req.CombinationActivityUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteBizError(c, core.NewBizError(400, "Invalid Request"))
		return
	}

	err := h.svc.UpdateCombinationActivity(c.Request.Context(), r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteCombinationActivity 删除拼团活动
func (h *CombinationActivityHandler) DeleteCombinationActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		core.WriteBizError(c, core.NewBizError(400, "Invalid ID"))
		return
	}

	err = h.svc.DeleteCombinationActivity(c.Request.Context(), id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetCombinationActivity 获得拼团活动
func (h *CombinationActivityHandler) GetCombinationActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		core.WriteBizError(c, core.NewBizError(400, "Invalid ID"))
		return
	}

	activity, err := h.svc.GetCombinationActivity(c.Request.Context(), id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, activity)
}

// GetCombinationActivityPage 获得拼团活动分页
func (h *CombinationActivityHandler) GetCombinationActivityPage(c *gin.Context) {
	var r req.CombinationActivityPageReq
	// Bind Query
	if err := c.ShouldBindQuery(&r); err != nil {
		// default bindings might fail for int, manually bind if needed
		// But ShouldBindQuery handles it usually.
	}
	// Manual defaults if zero
	if r.PageNo == 0 {
		r.PageNo = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 10
	}

	list, err := h.svc.GetCombinationActivityPage(c.Request.Context(), r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, list)
}
