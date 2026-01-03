package wallet

import (
	"strconv"

	pay2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	payData "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/wallet"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
	var r pay2.PayWalletPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetWalletPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// Convert list
	newRes := pagination.NewPageResult(make([]*pay2.PayWalletResp, 0, len(res.List)), res.Total)
	for _, item := range res.List {
		newRes.List = append(newRes.List, convertWalletResp(item))
	}
	response.WriteSuccess(c, newRes)
}

// GetWallet 获得会员钱包
func (h *PayWalletHandler) GetWallet(c *gin.Context) {
	userIdStr := c.Query("userId")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// 对齐 Java: payWalletService.getOrCreateWallet(reqVO.getUserId(), MEMBER.getValue())
	wallet, err := h.svc.GetOrCreateWallet(c, userId, 1) // 1: Member
	if err != nil {
		c.Error(err)
		return
	}
	response.WriteSuccess(c, convertWalletResp(wallet))
}

// UpdateWalletBalance 更新会员用户余额
func (h *PayWalletHandler) UpdateWalletBalance(c *gin.Context) {
	var r pay2.PayWalletUpdateBalanceReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 获得用户钱包
	wallet, err := h.svc.GetOrCreateWallet(c, r.UserID, 1) // 1: Member
	if err != nil {
		c.Error(err)
		return
	}
	if wallet == nil {
		response.WriteBizError(c, errors.ErrNotFound)
		return
	}

	// 更新钱包余额
	// walletID, bizID, bizType, price
	err = h.svc.AddWalletBalance(c, wallet.ID, strconv.FormatInt(r.UserID, 10), consts.PayWalletBizTypeUpdateBalance, r.Balance)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func convertWalletResp(item *pay.PayWallet) *pay2.PayWalletResp {
	if item == nil {
		return nil
	}
	return &pay2.PayWalletResp{
		ID:            item.ID,
		UserID:        item.UserID,
		UserType:      item.UserType,
		Balance:       item.Balance,
		TotalExpense:  item.TotalExpense,
		TotalRecharge: item.TotalRecharge,
		FreezePrice:   item.FreezePrice,
		CreateTime:    item.CreateTime,
	}
}
