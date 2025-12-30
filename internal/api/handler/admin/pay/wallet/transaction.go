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
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetWalletTransactionPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// Convert list
	newRes := pagination.NewPageResult(make([]*resp.PayWalletTransactionResp, 0, len(res.List)), res.Total)
	for _, item := range res.List {
		newRes.List = append(newRes.List, convertTransactionResp(item))
	}
	response.WriteSuccess(c, newRes)
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
		CreateTime: item.CreateTime,
	}
}
