package brokerage

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	tradeReq "github.com/wxlbd/ruoyi-mall-go/internal/api/req/app/trade"
	tradeResp "github.com/wxlbd/ruoyi-mall-go/internal/api/resp/app/trade"
	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	brokerageSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AppBrokerageWithdrawHandler struct {
	withdrawSvc    *brokerageSvc.BrokerageWithdrawService
	payTransferSvc *pay.PayTransferService // Needed? Java Controller uses PayTransferApi.
}

func NewAppBrokerageWithdrawHandler(withdrawSvc *brokerageSvc.BrokerageWithdrawService, payTransferSvc *pay.PayTransferService) *AppBrokerageWithdrawHandler {
	return &AppBrokerageWithdrawHandler{
		withdrawSvc:    withdrawSvc,
		payTransferSvc: payTransferSvc,
	}
}

// GetBrokerageWithdrawPage 获得分销提现分页
func (h *AppBrokerageWithdrawHandler) GetBrokerageWithdrawPage(c *gin.Context) {
	var reqVO tradeReq.AppBrokerageWithdrawPageReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	// userId := context.GetLoginUserId(c)
	userId := context.GetLoginUserID(c)
	pageReq := &req.BrokerageWithdrawPageReq{
		PageParam: reqVO.PageParam,
		UserID:    userId,
		Type:      reqVO.Type,
		Status:    reqVO.Status,
	}

	pageResult, err := h.withdrawSvc.GetBrokerageWithdrawPage(c, pageReq)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	writeResp := pagination.PageResult[*tradeResp.AppBrokerageWithdrawRespVO]{
		Total: pageResult.Total,
		List: lo.Map(pageResult.List, func(item *brokerage.BrokerageWithdraw, _ int) *tradeResp.AppBrokerageWithdrawRespVO {
			return &tradeResp.AppBrokerageWithdrawRespVO{
				ID:          item.ID,
				UserID:      item.UserID,
				Price:       item.Price,
				FeePrice:    item.FeePrice,
				TotalPrice:  item.TotalPrice,
				Type:        item.Type,
				Name:        item.UserName,    // Map UserName -> Name
				Account:     item.UserAccount, // Map UserAccount -> Account
				BankName:    item.BankName,
				Status:      item.Status,
				AuditReason: item.AuditReason,
				AuditTime:   item.AuditTime,
				Remark:      item.Remark,
				CreatedAt:   item.CreatedAt,
				// TypeName, StatusName -> Dict lookup (Frontend can handle or backend add logic)
			}
		}),
	}
	response.WriteSuccess(c, writeResp)
}

// GetBrokerageWithdraw 获得佣金提现详情
func (h *AppBrokerageWithdrawHandler) GetBrokerageWithdraw(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	userId := context.GetLoginUserID(c)
	withdraw, err := h.withdrawSvc.GetBrokerageWithdraw(c, id)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	if withdraw == nil || withdraw.UserID != userId {
		response.WriteSuccess(c, nil)
		return
	}

	// VO Conversion
	respVO := &tradeResp.AppBrokerageWithdrawRespVO{
		Status:      withdraw.Status,
		AuditReason: withdraw.AuditReason,
		AuditTime:   withdraw.AuditTime,
		Remark:      withdraw.Remark,
		CreatedAt:   withdraw.CreatedAt,
	}

	// Wechat Transfer Info Logic
	// Status: AUDIT_SUCCESS(10), Type: WECHAT(3)
	// We check against constants.
	if withdraw.Status == tradeModel.BrokerageWithdrawStatusAuditSuccess && withdraw.Type == tradeModel.BrokerageWithdrawTypeWechat && withdraw.PayTransferID > 0 {
		transfer, err := h.payTransferSvc.GetTransfer(c.Request.Context(), int64(withdraw.PayTransferID))
		if err != nil {
			response.WriteError(c, 500, err.Error())
			return
		}
		if transfer != nil {
			if transfer.ChannelExtras != nil {
				if val, ok := transfer.ChannelExtras["package_info"]; ok {
					respVO.TransferChannelPackageInfo = val
				}
				if val, ok := transfer.ChannelExtras["mch_id"]; ok {
					respVO.TransferChannelMchId = val
				}
			}
		}
	}

	response.WriteSuccess(c, respVO)
}

// CreateBrokerageWithdraw 创建分销提现
func (h *AppBrokerageWithdrawHandler) CreateBrokerageWithdraw(c *gin.Context) {
	var reqVO tradeReq.AppBrokerageWithdrawCreateReqVO
	if err := c.ShouldBindJSON(&reqVO); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	userId := context.GetLoginUserID(c)
	id, err := h.withdrawSvc.CreateBrokerageWithdraw(c, userId, &reqVO)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, id)
}
