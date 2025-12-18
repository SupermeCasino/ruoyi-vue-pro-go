package pay

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
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

	list := make([]resp.AppPayWalletPackageResp, len(pkgs))
	for i, item := range pkgs {
		list[i] = resp.AppPayWalletPackageResp{
			ID:         item.ID,
			Name:       item.Name,
			PayPrice:   item.PayPrice,
			BonusPrice: item.BonusPrice,
		}
	}

	response.WriteSuccess(c, list)
}
