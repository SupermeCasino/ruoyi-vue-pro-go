package wallet

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	payData "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/wallet"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type PayWalletRechargeHandler struct {
	svc *payData.PayWalletRechargeService
}

func NewPayWalletRechargeHandler(svc *payData.PayWalletRechargeService) *PayWalletRechargeHandler {
	return &PayWalletRechargeHandler{svc: svc}
}

// GetWalletRechargePage 获得会员钱包充值分页
func (h *PayWalletRechargeHandler) GetWalletRechargePage(c *gin.Context) {
	var r req.PayWalletRechargePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	res, err := h.svc.GetWalletRechargePage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	// Convert list
	newRes := pagination.NewPageResult(make([]*resp.PayWalletRechargeResp, 0, len(res.List)), res.Total)
	for _, item := range res.List {
		newRes.List = append(newRes.List, convertRechargeResp(item))
	}
	c.JSON(200, response.Success(newRes))
}

// UpdateWalletRechargePaid 更新钱包充值为已支付
func (h *PayWalletRechargeHandler) UpdateWalletRechargePaid(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	payOrderIdStr := c.Query("payOrderId")
	payOrderId, err := strconv.ParseInt(payOrderIdStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	if err := h.svc.UpdateWalletRechargerPaid(c, id, payOrderId); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// RefundWalletRecharge 发起钱包充值退款
func (h *PayWalletRechargeHandler) RefundWalletRecharge(c *gin.Context) {
	/*
		// Java: @RequestParam("id") Long id. So it is form/query param.
		// Check gin bind: ShouldBindQuery or ShouldBindJSON.
	*/
	idStr := c.Query("id")
	if idStr == "" {
		idStr = c.PostForm("id")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	if err := h.svc.RefundWalletRecharge(c, id, c.ClientIP()); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// UpdateWalletRechargeRefunded 更新钱包充值为已退款
func (h *PayWalletRechargeHandler) UpdateWalletRechargeRefunded(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	payRefundIdStr := c.Query("payRefundId")
	payRefundId, err := strconv.ParseInt(payRefundIdStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	if err := h.svc.UpdateWalletRechargeRefunded(c, id, payRefundId); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

func convertRechargeResp(item *pay.PayWalletRecharge) *resp.PayWalletRechargeResp {
	if item == nil {
		return nil
	}
	return &resp.PayWalletRechargeResp{
		ID:               item.ID,
		WalletID:         item.WalletID,
		TotalPrice:       item.TotalPrice,
		PayPrice:         item.PayPrice,
		BonusPrice:       item.BonusPrice,
		PackageID:        item.PackageID,
		PayStatus:        item.PayStatus,
		PayOrderID:       item.PayOrderID,
		PayChannelCode:   item.PayChannelCode,
		PayTime:          item.PayTime,
		RefundStatus:     item.RefundStatus,
		PayRefundID:      item.PayRefundID,
		RefundTotalPrice: item.RefundTotalPrice,
		RefundPayPrice:   item.RefundPayPrice,
		RefundBonusPrice: item.RefundBonusPrice,
		RefundTime:       item.RefundTime,
		CreateTime:       item.CreatedAt,
	}
}
