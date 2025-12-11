package promotion

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/promotion"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SeckillConfigHandler struct {
	svc *promotion.SeckillConfigService
}

func NewSeckillConfigHandler(svc *promotion.SeckillConfigService) *SeckillConfigHandler {
	return &SeckillConfigHandler{svc: svc}
}

// CreateSeckillConfig 创建
func (h *SeckillConfigHandler) CreateSeckillConfig(c *gin.Context) {
	var r req.SeckillConfigCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	id, err := h.svc.CreateSeckillConfig(c.Request.Context(), &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, id)
}

// UpdateSeckillConfig 更新
func (h *SeckillConfigHandler) UpdateSeckillConfig(c *gin.Context) {
	var r req.SeckillConfigUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateSeckillConfig(c.Request.Context(), &r); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// UpdateSeckillConfigStatus 更新状态
func (h *SeckillConfigHandler) UpdateSeckillConfigStatus(c *gin.Context) {
	var r req.SeckillConfigUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateSeckillConfigStatus(c.Request.Context(), r.ID, r.Status); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// DeleteSeckillConfig 删除
func (h *SeckillConfigHandler) DeleteSeckillConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeleteSeckillConfig(c.Request.Context(), id); err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, true)
}

// GetSeckillConfig 获得
func (h *SeckillConfigHandler) GetSeckillConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	res, err := h.svc.GetSeckillConfig(c.Request.Context(), id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// GetSeckillConfigPage 分页
func (h *SeckillConfigHandler) GetSeckillConfigPage(c *gin.Context) {
	var r req.SeckillConfigPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, err.Error())
		return
	}
	res, err := h.svc.GetSeckillConfigPage(c.Request.Context(), &r)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, res)
}

// GetSeckillConfigList 获得列表
func (h *SeckillConfigHandler) GetSeckillConfigList(c *gin.Context) {
	res, err := h.svc.GetSeckillConfigList(c.Request.Context())
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	var respList []resp.SeckillConfigResp
	for _, v := range res {
		respList = append(respList, resp.SeckillConfigResp{
			ID:            v.ID,
			Name:          v.Name,
			StartTime:     v.StartTime,
			EndTime:       v.EndTime,
			SliderPicUrls: v.SliderPicUrls,
			Status:        v.Status,
			CreateTime:    v.CreatedAt,
		})
	}
	core.WriteSuccess(c, respList)
}

// GetSeckillConfigSimpleList 精简列表
func (h *SeckillConfigHandler) GetSeckillConfigSimpleList(c *gin.Context) {
	res, err := h.svc.GetSeckillConfigListByStatus(c.Request.Context(), 1) // Enable
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	var respList []resp.SeckillConfigSimpleResp
	for _, v := range res {
		respList = append(respList, resp.SeckillConfigSimpleResp{
			ID:        v.ID,
			Name:      v.Name,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		})
	}
	core.WriteSuccess(c, respList)
}
