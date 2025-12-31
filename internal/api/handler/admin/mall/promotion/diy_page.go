package promotion

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type DiyPageHandler struct {
	svc promotion.DiyPageService
}

func NewDiyPageHandler(svc promotion.DiyPageService) *DiyPageHandler {
	return &DiyPageHandler{svc: svc}
}

func (h *DiyPageHandler) CreateDiyPage(c *gin.Context) {
	var r req.DiyPageCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateDiyPage(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *DiyPageHandler) UpdateDiyPage(c *gin.Context) {
	var r req.DiyPageUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.UpdateDiyPage(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DiyPageHandler) DeleteDiyPage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.DeleteDiyPage(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DiyPageHandler) GetDiyPage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetDiyPage(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

func (h *DiyPageHandler) GetDiyPagePage(c *gin.Context) {
	var r req.DiyPagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetDiyPagePage(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetDiyPageProperty 获得装修页面属性
func (h *DiyPageHandler) GetDiyPageProperty(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetDiyPageProperty(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetDiyPageList 获得装修页面列表
// Java: DiyPageController#getDiyPageList
func (h *DiyPageHandler) GetDiyPageList(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		response.WriteSuccess(c, []interface{}{})
		return
	}
	// 使用 model.IntListFromCSV 解析 ID 列表
	var ids model.IntListFromCSV
	if err := ids.Scan(idsStr); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 转换为 []int64
	ids64 := make([]int64, len(ids))
	for i, id := range ids {
		ids64[i] = int64(id)
	}
	res, err := h.svc.GetDiyPageList(c, ids64)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// UpdateDiyPageProperty 更新装修页面属性
// Java: DiyPageController#updateDiyPageProperty
func (h *DiyPageHandler) UpdateDiyPageProperty(c *gin.Context) {
	var r req.DiyPagePropertyUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.UpdateDiyPageProperty(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}
