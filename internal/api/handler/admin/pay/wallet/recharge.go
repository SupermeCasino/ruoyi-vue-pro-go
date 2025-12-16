package wallet

import (
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
