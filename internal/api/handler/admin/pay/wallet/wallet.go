package wallet

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/pay"
	"backend-go/internal/pkg/core"
	payData "backend-go/internal/service/pay/wallet"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PayWalletHandler struct {
	svc *payData.PayWalletService
}

func NewPayWalletHandler(svc *payData.PayWalletService) *PayWalletHandler {
	return &PayWalletHandler{svc: svc}
}

// GetWalletPage 获得会员钱包分页
func (h *PayWalletHandler) GetWalletPage(c *gin.Context) {
	var r req.PayWalletPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetWalletPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	// Convert list
	newRes := core.NewPageResult(make([]*resp.PayWalletResp, 0, len(res.List)), res.Total)
	for _, item := range res.List {
		newRes.List = append(newRes.List, convertWalletResp(item))
	}
	c.JSON(200, core.Success(newRes))
}

// GetWallet 获得会员钱包
func (h *PayWalletHandler) GetWallet(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	wallet, err := h.svc.GetWallet(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(convertWalletResp(wallet)))
}

func convertWalletResp(item *pay.PayWallet) *resp.PayWalletResp {
	if item == nil {
		return nil
	}
	return &resp.PayWalletResp{
		ID:            item.ID,
		UserID:        item.UserID,
		UserType:      item.UserType,
		Balance:       item.Balance,
		TotalExpense:  item.TotalExpense,
		TotalRecharge: item.TotalRecharge,
		FreezePrice:   item.FreezePrice,
		CreateTime:    item.CreatedAt,
	}
}
