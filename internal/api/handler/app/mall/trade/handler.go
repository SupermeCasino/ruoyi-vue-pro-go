package trade

import (
	"github.com/google/wire"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/mall/trade/brokerage"
)

var ProviderSet = wire.NewSet(
	NewAppTradeAfterSaleHandler,
	NewAppCartHandler,
	NewAppTradeConfigHandler,
	NewAppTradeOrderHandler,
	NewHandlers,
	brokerage.ProviderSet,
)

type Handlers struct {
	AfterSale *AppTradeAfterSaleHandler
	Cart      *AppCartHandler
	Config    *AppTradeConfigHandler
	Order     *AppTradeOrderHandler
	Brokerage *brokerage.Handlers
}

func NewHandlers(
	afterSale *AppTradeAfterSaleHandler,
	cart *AppCartHandler,
	config *AppTradeConfigHandler,
	order *AppTradeOrderHandler,
	brokerageHandlers *brokerage.Handlers,
) *Handlers {
	return &Handlers{
		AfterSale: afterSale,
		Cart:      cart,
		Config:    config,
		Order:     order,
		Brokerage: brokerageHandlers,
	}
}
