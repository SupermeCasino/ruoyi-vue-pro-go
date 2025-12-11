package member

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	memberModel "backend-go/internal/model/member"
	"backend-go/internal/pkg/core"
	memberSvc "backend-go/internal/service/member"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AppMemberSignInRecordHandler struct {
	svc *memberSvc.MemberSignInRecordService
}

func NewAppMemberSignInRecordHandler(svc *memberSvc.MemberSignInRecordService) *AppMemberSignInRecordHandler {
	return &AppMemberSignInRecordHandler{svc: svc}
}

// GetSignInRecordSummary 获得个人签到统计
func (h *AppMemberSignInRecordHandler) GetSignInRecordSummary(c *gin.Context) {
	userId := c.GetInt64(core.CtxUserIDKey)
	summary, err := h.svc.GetSignInRecordSummary(c, userId)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}
	core.WriteSuccess(c, summary)
}

// CreateSignInRecord 创建签到记录
func (h *AppMemberSignInRecordHandler) CreateSignInRecord(c *gin.Context) {
	userId := c.GetInt64(core.CtxUserIDKey)
	record, err := h.svc.CreateSignInRecord(c, userId)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	// Convert record to resp (simplified, or just return record)
	// Java doesn't return much? It returns the full record.
	// We'll return simplified App resp.
	core.WriteSuccess(c, resp.AppMemberSignInRecordResp{
		ID:         record.ID,
		Day:        record.Day,
		Point:      record.Point,
		Experience: record.Experience,
		CreatedAt:  record.CreatedAt,
	})
}

// GetSignInRecordPage 获得个人签到分页
func (h *AppMemberSignInRecordHandler) GetSignInRecordPage(c *gin.Context) {
	var r req.AppMemberSignInRecordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteBizError(c, core.ErrParam)
		return
	}
	userId := c.GetInt64(core.CtxUserIDKey)

	// Reuse service method but fill UserID
	pageReq := req.MemberSignInRecordPageReq{
		PageParam: core.PageParam{
			PageNo:   r.PageNo,
			PageSize: r.PageSize,
		},
		UserID: userId,
	}

	pageResult, err := h.svc.GetSignInRecordPage(c, &pageReq)
	if err != nil {
		core.WriteBizError(c, err)
		return
	}

	respList := lo.Map(pageResult.List, func(item *memberModel.MemberSignInRecord, _ int) resp.AppMemberSignInRecordResp {
		return resp.AppMemberSignInRecordResp{
			ID:         item.ID,
			Day:        item.Day,
			Point:      item.Point,
			Experience: item.Experience,
			CreatedAt:  item.CreatedAt,
		}
	})
	core.WriteSuccess(c, core.NewPageResult(respList, pageResult.Total))
}
