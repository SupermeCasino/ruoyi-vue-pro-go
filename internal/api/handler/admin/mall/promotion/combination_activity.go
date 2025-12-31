package promotion

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type CombinationActivityHandler struct {
	svc       promotion.CombinationActivityService
	recordSvc promotion.CombinationRecordService
	spuSvc    *product.ProductSpuService
}

func NewCombinationActivityHandler(
	svc promotion.CombinationActivityService,
	recordSvc promotion.CombinationRecordService,
	spuSvc *product.ProductSpuService,
) *CombinationActivityHandler {
	return &CombinationActivityHandler{svc: svc, recordSvc: recordSvc, spuSvc: spuSvc}
}

// CreateCombinationActivity 创建拼团活动
func (h *CombinationActivityHandler) CreateCombinationActivity(c *gin.Context) {
	var r req.CombinationActivityCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "Invalid Request"))
		return
	}

	id, err := h.svc.CreateCombinationActivity(c.Request.Context(), r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateCombinationActivity 更新拼团活动
func (h *CombinationActivityHandler) UpdateCombinationActivity(c *gin.Context) {
	var r req.CombinationActivityUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "Invalid Request"))
		return
	}

	err := h.svc.UpdateCombinationActivity(c.Request.Context(), r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteCombinationActivity 删除拼团活动
func (h *CombinationActivityHandler) DeleteCombinationActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "Invalid ID"))
		return
	}

	err = h.svc.DeleteCombinationActivity(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// CloseCombinationActivity 关闭拼团活动
func (h *CombinationActivityHandler) CloseCombinationActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "Invalid ID"))
		return
	}

	err = h.svc.CloseCombinationActivity(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetCombinationActivity 获得拼团活动
func (h *CombinationActivityHandler) GetCombinationActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "Invalid ID"))
		return
	}

	activity, err := h.svc.GetCombinationActivity(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, activity)
}

// GetCombinationActivityListByIds 获得拼团活动列表，基于活动编号数组
// Java: CombinationActivityController#getCombinationActivityListByIds
func (h *CombinationActivityHandler) GetCombinationActivityListByIds(c *gin.Context) {
	var req struct {
		Ids model.IntListFromCSV `form:"ids" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "Invalid IDs"))
		return
	}

	// 将 []int 转换为 []int64
	ids := make([]int64, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = int64(id)
	}

	list, err := h.svc.GetCombinationActivityListByIds(c.Request.Context(), ids)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
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
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}
