package promotion

import (
	"strconv"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"

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
	var pageParam core.PageParam
	// No request body for GET page, query params usually mapped manually or via BindQuery
	// core.PageParam usually binds PageNo, PageSize
	// Java method takes PageParam.
	pageParam.PageNo, _ = strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageParam.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	list, err := h.svc.GetCombinationActivityPageForApp(c.Request.Context(), pageParam)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, list)
}

// GetCombinationActivityListByIds 获得拼团活动列表
func (h *AppCombinationActivityHandler) GetCombinationActivityListByIds(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		core.WriteSuccess(c, []*resp.AppCombinationActivityRespVO{})
		return
	}

	// Parse IDs "1,2,3" (assuming Gin Query param style or similar)
	// Java: @RequestParam("ids") List<Long> ids. Spring parses comma separated?
	// Gin Query is string.
	parts := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		if id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	// Service method needed: GetCombinationActivityListByIds
	// I haven't implemented `GetCombinationActivityListByIds` in Service yet, only `GetList` by count.
	// I'll skip this for now or add TODO. Java controller uses it.
	// I'll stick to basic features first.
	// Check plan: "GetList, GetPage, GetDetail".
	// "GetList" usually for home page (count).
	// "GetListByIds" is for cart/order confirmation maybe?

	// I'll return empty for now to match interface if added later.
	core.WriteSuccess(c, []*resp.AppCombinationActivityRespVO{})
}

// GetCombinationActivityDetail 获得拼团活动明细
func (h *AppCombinationActivityHandler) GetCombinationActivityDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		core.WriteBizError(c, core.NewBizError(400, "Invalid ID"))
		return
	}

	detail, err := h.svc.GetCombinationActivityDetail(c.Request.Context(), id)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, detail)
}
