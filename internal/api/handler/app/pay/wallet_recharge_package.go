package pay

import (
	"github.com/gin-gonic/gin"
	pay2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/pay"
	payWalletSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/wallet"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type AppPayWalletRechargePackageHandler struct {
	svc *payWalletSvc.PayWalletRechargePackageService
}

func NewAppPayWalletRechargePackageHandler(svc *payWalletSvc.PayWalletRechargePackageService) *AppPayWalletRechargePackageHandler {
	return &AppPayWalletRechargePackageHandler{svc: svc}
}

// GetWalletRechargePackageList 获得钱包充值套餐列表
func (h *AppPayWalletRechargePackageHandler) GetWalletRechargePackageList(c *gin.Context) {
	pkgs, err := h.svc.GetWalletRechargePackageList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	list := make([]pay2.AppPayWalletPackageResp, len(pkgs))
	for i, item := range pkgs {
		list[i] = pay2.AppPayWalletPackageResp{
			ID:         item.ID,
			Name:       item.Name,
			PayPrice:   item.PayPrice,
			BonusPrice: item.BonusPrice,
		}
	}

	response.WriteSuccess(c, list)
}
