package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/promotion"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SeckillActivityHandler struct {
	svc *promotion.SeckillActivityService
}

func NewSeckillActivityHandler(svc *promotion.SeckillActivityService) *SeckillActivityHandler {
	return &SeckillActivityHandler{svc: svc}
}

// CreateSeckillActivity 创建
func (h *SeckillActivityHandler) CreateSeckillActivity(c *gin.Context) {
	var r req.SeckillActivityCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreateSeckillActivity(c.Request.Context(), &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateSeckillActivity 更新
func (h *SeckillActivityHandler) UpdateSeckillActivity(c *gin.Context) {
	var r req.SeckillActivityUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateSeckillActivity(c.Request.Context(), &r); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteSeckillActivity 删除
func (h *SeckillActivityHandler) DeleteSeckillActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteSeckillActivity(c.Request.Context(), id); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// CloseSeckillActivity 关闭
func (h *SeckillActivityHandler) CloseSeckillActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.CloseSeckillActivity(c.Request.Context(), id); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetSeckillActivity 获得详情
func (h *SeckillActivityHandler) GetSeckillActivity(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	act, err := h.svc.GetSeckillActivity(c.Request.Context(), id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	if act == nil {
		core.WriteSuccess(c, nil)
		return
	}
	products, err := h.svc.GetSeckillProductListByActivityID(c.Request.Context(), id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	detail := resp.SeckillActivityDetailResp{}
	detail.SeckillActivityResp = resp.SeckillActivityResp{
		ID:               act.ID,
		SpuID:            act.SpuID,
		Name:             act.Name,
		Status:           act.Status,
		Remark:           act.Remark,
		StartTime:        act.StartTime,
		EndTime:          act.EndTime,
		Sort:             act.Sort,
		ConfigIds:        act.ConfigIds,
		TotalLimitCount:  act.TotalLimitCount,
		SingleLimitCount: act.SingleLimitCount,
		Stock:            act.Stock,
		TotalStock:       act.TotalStock,
		CreateTime:       act.CreatedAt,
	}
	for _, p := range products {
		detail.Products = append(detail.Products, resp.SeckillProductResp{
			ID:           p.ID,
			ActivityID:   p.ActivityID,
			SpuID:        p.SpuID,
			SkuID:        p.SkuID,
			SeckillPrice: p.SeckillPrice,
			Stock:        p.Stock,
		})
	}
	core.WriteSuccess(c, detail)
}

// GetSeckillActivityPage 分页
func (h *SeckillActivityHandler) GetSeckillActivityPage(c *gin.Context) {
	var r req.SeckillActivityPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	res, err := h.svc.GetSeckillActivityPage(c.Request.Context(), &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	// Convert to Resp
	list := make([]resp.SeckillActivityResp, len(res.List))
	for i, v := range res.List {
		list[i] = resp.SeckillActivityResp{
			ID:               v.ID,
			SpuID:            v.SpuID,
			Name:             v.Name,
			Status:           v.Status,
			Remark:           v.Remark,
			StartTime:        v.StartTime,
			EndTime:          v.EndTime,
			Sort:             v.Sort,
			ConfigIds:        v.ConfigIds,
			TotalLimitCount:  v.TotalLimitCount,
			SingleLimitCount: v.SingleLimitCount,
			Stock:            v.Stock,
			TotalStock:       v.TotalStock,
			CreateTime:       v.CreatedAt,
		}
	}

	core.WriteSuccess(c, core.PageResult[resp.SeckillActivityResp]{
		List:  list,
		Total: res.Total,
	})
}
