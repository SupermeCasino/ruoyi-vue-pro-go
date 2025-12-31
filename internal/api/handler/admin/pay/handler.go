package pay

import (
	"github.com/google/wire"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay/wallet"
)

var ProviderSet = wire.NewSet(
	NewPayAppHandler,
	NewPayChannelHandler,
	NewPayNotifyHandler,
	NewPayOrderHandler,
	NewPayRefundHandler,
	NewPayTransferHandler,
	NewHandlers,
	wallet.ProviderSet,
)

type Handlers struct {
	App      *PayAppHandler
	Channel  *PayChannelHandler
	Notify   *PayNotifyHandler
	Order    *PayOrderHandler
	Refund   *PayRefundHandler
	Transfer *PayTransferHandler
	Wallet   *wallet.Handlers
}

func NewHandlers(
	app *PayAppHandler,
	channel *PayChannelHandler,
	notify *PayNotifyHandler,
	order *PayOrderHandler,
	refund *PayRefundHandler,
	transfer *PayTransferHandler,
	wallet *wallet.Handlers,
) *Handlers {
	return &Handlers{
		App:      app,
		Channel:  channel,
		Notify:   notify,
		Order:    order,
		Refund:   refund,
		Transfer: transfer,
		Wallet:   wallet,
	}
}
