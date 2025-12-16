package wallet

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	payData "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/wallet"

	"github.com/gin-gonic/gin"
)

type PayWalletTransactionHandler struct {
	svc *payData.PayWalletTransactionService
}

func NewPayWalletTransactionHandler(svc *payData.PayWalletTransactionService) *PayWalletTransactionHandler {
	return &PayWalletTransactionHandler{svc: svc}
}

// GetWalletTransactionPage 获得会员钱包流水分页
func (h *PayWalletTransactionHandler) GetWalletTransactionPage(c *gin.Context) {
	var r req.PayWalletTransactionPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetWalletTransactionPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	// Convert list
	newRes := core.NewPageResult(make([]*resp.PayWalletTransactionResp, 0, len(res.List)), res.Total)
	for _, item := range res.List {
		newRes.List = append(newRes.List, convertTransactionResp(item))
	}
	c.JSON(200, core.Success(newRes))
}

func convertTransactionResp(item *pay.PayWalletTransaction) *resp.PayWalletTransactionResp {
	if item == nil {
		return nil
	}
	return &resp.PayWalletTransactionResp{
		ID:         item.ID,
		WalletID:   item.WalletID,
		BizType:    item.BizType,
		BizID:      item.BizID,
		No:         item.No,
		Title:      item.Title,
		Price:      item.Price,
		Balance:    item.Balance,
		CreateTime: item.CreatedAt,
	}
}
