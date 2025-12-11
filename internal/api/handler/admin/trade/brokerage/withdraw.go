package brokerage

import (
	"github.com/gin-gonic/gin"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	brokerageModel "backend-go/internal/model/trade/brokerage"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service/member"
	"backend-go/internal/service/trade/brokerage"
)

type BrokerageWithdrawHandler struct {
	withdrawSvc   *brokerage.BrokerageWithdrawService
	memberUserSvc *member.MemberUserService
}

func NewBrokerageWithdrawHandler(withdrawSvc *brokerage.BrokerageWithdrawService, memberUserSvc *member.MemberUserService) *BrokerageWithdrawHandler {
	return &BrokerageWithdrawHandler{
		withdrawSvc:   withdrawSvc,
		memberUserSvc: memberUserSvc,
	}
}

// ApproveBrokerageWithdraw 通过申请
func (h *BrokerageWithdrawHandler) ApproveBrokerageWithdraw(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "参数错误")
		return
	}

	// 10: AUDIT_SUCCESS (See Enum in Java)
	if err := h.withdrawSvc.AuditBrokerageWithdraw(c, id, 10, ""); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	core.WriteSuccess(c, true)
}

// RejectBrokerageWithdraw 驳回申请
func (h *BrokerageWithdrawHandler) RejectBrokerageWithdraw(c *gin.Context) {
	var r req.BrokerageWithdrawRejectReq
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}

	// 20: AUDIT_FAIL
	if err := h.withdrawSvc.AuditBrokerageWithdraw(c, r.ID, 20, r.AuditReason); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	core.WriteSuccess(c, true)
}

// UpdateBrokerageWithdrawTransferred 更新佣金提现的转账结果
func (h *BrokerageWithdrawHandler) UpdateBrokerageWithdrawTransferred(c *gin.Context) {
	// Placeholder: Req Struct
	var r struct {
		MerchantTransferID string `json:"merchantTransferId"`
		PayTransferID      int64  `json:"payTransferId"`
	}
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}

	id := core.ParseInt64(r.MerchantTransferID)
	if err := h.withdrawSvc.UpdateBrokerageWithdrawTransferred(c, id, r.PayTransferID); err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	core.WriteSuccess(c, true)
}

// GetBrokerageWithdraw 获得佣金提现
func (h *BrokerageWithdrawHandler) GetBrokerageWithdraw(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	if id == 0 {
		core.WriteError(c, 400, "参数错误")
		return
	}

	withdraw, err := h.withdrawSvc.GetBrokerageWithdraw(c, id)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}
	if withdraw == nil {
		core.WriteError(c, 404, "提现记录不存在")
		return
	}

	res := h.convert(withdraw)
	core.WriteSuccess(c, res)
}

// GetBrokerageWithdrawPage 获得佣金提现分页
func (h *BrokerageWithdrawHandler) GetBrokerageWithdrawPage(c *gin.Context) {
	var r req.BrokerageWithdrawPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteError(c, 400, "参数错误")
		return
	}

	pageResult, err := h.withdrawSvc.GetBrokerageWithdrawPage(c, &r)
	if err != nil {
		core.WriteError(c, 500, err.Error())
		return
	}

	// Aggregate User Info
	userIds := make([]int64, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		userIds = append(userIds, item.UserID)
	}
	userMap, _ := h.memberUserSvc.GetUserMap(c, userIds)

	list := make([]resp.BrokerageWithdrawResp, len(pageResult.List))
	for i, item := range pageResult.List {
		res := h.convert(item)
		if u, ok := userMap[item.UserID]; ok {
			res.BrokerageUserResp.Nickname = u.Nickname
			res.BrokerageUserResp.Avatar = u.Avatar
		}
		list[i] = res
	}

	core.WriteSuccess(c, core.PageResult[resp.BrokerageWithdrawResp]{
		List:  list,
		Total: pageResult.Total,
	})
}

func (h *BrokerageWithdrawHandler) convert(do *brokerageModel.BrokerageWithdraw) resp.BrokerageWithdrawResp {
	payTransferId := int64(0)
	if do.PayTransferID > 0 {
		payTransferId = do.PayTransferID
	}
	return resp.BrokerageWithdrawResp{
		ID:                  do.ID,
		UserID:              do.UserID,
		Price:               do.Price,
		FeePrice:            do.FeePrice,
		TotalPrice:          do.TotalPrice,
		Type:                do.Type,
		UserName:            do.UserName,
		UserAccount:         do.UserAccount,
		QRCodeUrl:           do.QrCodeURL,
		BankName:            do.BankName,
		BankAddress:         do.BankAddress,
		Status:              do.Status,
		AuditReason:         do.AuditReason,
		AuditTime:           do.AuditTime,
		Remark:              do.Remark,
		PayTransferID:       payTransferId,
		TransferChannelCode: do.TransferChannelCode,
		TransferTime:        do.TransferTime,
		TransferErrorMsg:    do.TransferErrorMsg,
		CreateTime:          do.CreatedAt,
		BrokerageUserResp: resp.BrokerageUserResp{
			ID: do.UserID,
			// Note: BrokerageUserResp expects BrokerageUserId (the user IS the brokerage user here)
			// But BrokerageUserResp fields (Price, Frozen...) might need to be fetched?
			// Java controller: `BrokerageWithdrawConvert.INSTANCE.convertPage(pageResult, userMap)`
			// Standard UserInfo (Nickname/Avatar) is aggregated.
			// We only fill Nickname/Avatar in GetPage loop.
		},
	}
}
