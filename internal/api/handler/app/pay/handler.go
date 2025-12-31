package pay

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAppPayChannelHandler,
	NewAppPayOrderHandler,
	NewAppPayTransferHandler,
	NewAppPayWalletHandler,
	NewAppPayWalletRechargePackageHandler,
	NewAppPayWalletTransactionHandler,
	NewHandlers,
)

type Handlers struct {
	Channel               *AppPayChannelHandler
	Order                 *AppPayOrderHandler
	Transfer              *AppPayTransferHandler
	Wallet                *AppPayWalletHandler
	WalletRechargePackage *AppPayWalletRechargePackageHandler
	WalletTransaction     *AppPayWalletTransactionHandler
}

func NewHandlers(
	channel *AppPayChannelHandler,
	order *AppPayOrderHandler,
	transfer *AppPayTransferHandler,
	wallet *AppPayWalletHandler,
	walletRechargePackage *AppPayWalletRechargePackageHandler,
	walletTransaction *AppPayWalletTransactionHandler,
) *Handlers {
	return &Handlers{
		Channel:               channel,
		Order:                 order,
		Transfer:              transfer,
		Wallet:                wallet,
		WalletRechargePackage: walletRechargePackage,
		WalletTransaction:     walletTransaction,
	}
}
