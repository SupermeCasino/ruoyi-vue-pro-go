package promotion

import (
	"strconv"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppCombinationActivityHandler struct {
	svc promotion.CombinationActivityService
}

func NewAppCombinationActivityHandler(svc promotion.CombinationActivityService) *AppCombinationActivityHandler {
	return &AppCombinationActivityHandler{svc: svc}
}

// GetCombinationActivityPage 获得拼团活动分页
func (h *AppCombinationActivityHandler) GetCombinationActivityPage(c *gin.Context) {
	var pageParam pagination.PageParam
	// No request body for GET page, query params usually mapped manually or via BindQuery
	// pagination.PageParam usually binds PageNo, PageSize
	// Java method takes PageParam.
	pageParam.PageNo, _ = strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageParam.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	list, err := h.svc.GetCombinationActivityPageForApp(c.Request.Context(), pageParam)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}

// GetCombinationActivityListByIds 获得拼团活动列表
func (h *AppCombinationActivityHandler) GetCombinationActivityListByIds(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		response.WriteSuccess(c, []*resp.AppCombinationActivityRespVO{})
		return
	}

	parts := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		if id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		response.WriteSuccess(c, []*resp.AppCombinationActivityRespVO{})
		return
	}

	// 1. 获取活动列表 (已通过 Service 进行 SPU 过滤与详情聚合)
	list, err := h.svc.GetCombinationActivityListByIdsForApp(c.Request.Context(), ids)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, list)
}

// GetCombinationActivityDetail 获得拼团活动明细
func (h *AppCombinationActivityHandler) GetCombinationActivityDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "Invalid ID"))
		return
	}

	detail, err := h.svc.GetCombinationActivityDetail(c.Request.Context(), id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, detail)
}
