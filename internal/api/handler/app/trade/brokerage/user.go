package brokerage

import (
	tradeReq "backend-go/internal/api/req/app/trade"
	tradeResp "backend-go/internal/api/resp/app/trade"
	model "backend-go/internal/model/trade/brokerage"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/trade/brokerage"
	"time"

	"github.com/gin-gonic/gin"
)

type AppBrokerageUserHandler struct {
	userSvc     *brokerage.BrokerageUserService
	recordSvc   *brokerage.BrokerageRecordService
	withdrawSvc *brokerage.BrokerageWithdrawService
}

func NewAppBrokerageUserHandler(userSvc *brokerage.BrokerageUserService, recordSvc *brokerage.BrokerageRecordService, withdrawSvc *brokerage.BrokerageWithdrawService) *AppBrokerageUserHandler {
	return &AppBrokerageUserHandler{
		userSvc:     userSvc,
		recordSvc:   recordSvc,
		withdrawSvc: withdrawSvc,
	}
}

// GetBrokerageUser 获得个人分销信息
func (h *AppBrokerageUserHandler) GetBrokerageUser(c *gin.Context) {
	userId := c.GetInt64("userId")
	user, err := h.userSvc.GetOrCreateBrokerageUser(c, userId)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	respVO := &tradeResp.AppBrokerageUserRespVO{
		BrokerageEnabled: user.BrokerageEnabled,
		BrokeragePrice:   user.BrokeragePrice,
		FrozenPrice:      user.FrozenPrice,
	}
	core.WriteSuccess(c, respVO)
}

// BindBrokerageUser 绑定推广员
func (h *AppBrokerageUserHandler) BindBrokerageUser(c *gin.Context) {
	var r tradeReq.AppBrokerageUserBindReqVO
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}
	userId := c.GetInt64("userId")
	success, err := h.userSvc.BindBrokerageUser(c, userId, r.BindUserID)
	if err != nil {
		core.WriteError(c, 500, err.Error()) // Or 400
		return
	}
	core.WriteSuccess(c, success)
}

// GetBrokerageUserSummary 获得个人分销统计
func (h *AppBrokerageUserHandler) GetBrokerageUserSummary(c *gin.Context) {
	userId := c.GetInt64("userId")
	user, err := h.userSvc.GetBrokerageUser(c, userId)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	if user == nil {
		// Should create? Java uses getBrokerageUser and assuming it exists?
		// Actually Java calls getOrCreate in /get, but here it calls get.
		// If nil, use default 0
		user = &model.BrokerageUser{} // Empty
	}

	// 1. Yesterday Price (Call Record Service)
	yesterday := time.Now().AddDate(0, 0, -1)
	beginOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, time.Local)

	// BizType: ORDER=1 (Assume). Status: SETTLEMENT=1 (Assume).
	yesterdayPrice, err := h.recordSvc.GetSummaryPriceByUserId(c, userId, 1, 1, beginOfDay, endOfDay)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	// 2. Withdraw Price
	// Status: AUDIT_SUCCESS(10), WITHDRAW_SUCCESS(11)
	summaries, err := h.withdrawSvc.GetWithdrawSummaryListByUserId(c, []int64{userId}, []int{10, 11})
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	withdrawPrice := 0
	if len(summaries) > 0 {
		withdrawPrice = summaries[0].Price
	}

	// 3. Count
	firstCount, _ := h.userSvc.GetBrokerageUserCountByBindUserId(c, userId, 1)
	secondCount, _ := h.userSvc.GetBrokerageUserCountByBindUserId(c, userId, 2)

	respVO := &tradeResp.AppBrokerageUserMySummaryRespVO{
		YesterdayPrice:           yesterdayPrice,
		WithdrawPrice:            withdrawPrice,
		FirstBrokerageUserCount:  int(firstCount),
		SecondBrokerageUserCount: int(secondCount),
		BrokeragePrice:           user.BrokeragePrice,
		FrozenPrice:              user.FrozenPrice,
	}
	core.WriteSuccess(c, respVO)
}

// GetBrokerageUserChildSummaryPage 获得下级分销统计分页
func (h *AppBrokerageUserHandler) GetBrokerageUserChildSummaryPage(c *gin.Context) {
	var r tradeReq.AppBrokerageUserChildSummaryPageReqVO
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}
	userId := c.GetInt64("userId")
	pageResult, err := h.userSvc.GetBrokerageUserChildSummaryPage(c, &r, userId)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	// Convert to RESP VO (Fetch User Info)
	// Placeholder conversion
	list := make([]tradeResp.AppBrokerageUserChildSummaryRespVO, len(pageResult.List))
	for i, u := range pageResult.List {
		list[i] = tradeResp.AppBrokerageUserChildSummaryRespVO{
			ID:             u.ID,
			BrokeragePrice: u.BrokeragePrice,
			// Nickname/Avatar need member service
		}
	}

	core.WriteSuccess(c, &core.PageResult[tradeResp.AppBrokerageUserChildSummaryRespVO]{
		List:  list,
		Total: pageResult.Total,
	})
}
