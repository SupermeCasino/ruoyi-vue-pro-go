package wallet

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewPayWalletRechargeHandler,
	NewPayWalletRechargePackageHandler,
	NewPayWalletTransactionHandler,
	NewPayWalletHandler,
	NewHandlers,
)

type Handlers struct {
	Recharge        *PayWalletRechargeHandler
	RechargePackage *PayWalletRechargePackageHandler
	Transaction     *PayWalletTransactionHandler
	Wallet          *PayWalletHandler
}

func NewHandlers(
	recharge *PayWalletRechargeHandler,
	rechargePackage *PayWalletRechargePackageHandler,
	transaction *PayWalletTransactionHandler,
	wallet *PayWalletHandler,
) *Handlers {
	return &Handlers{
		Recharge:        recharge,
		RechargePackage: rechargePackage,
		Transaction:     transaction,
		Wallet:          wallet,
	}
}
