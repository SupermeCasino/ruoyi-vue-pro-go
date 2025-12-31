package trade

import (
	"github.com/google/wire"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall/trade/brokerage"
)

var ProviderSet = wire.NewSet(
	NewTradeAfterSaleHandler,
	NewTradeConfigHandler,
	NewDeliveryExpressHandler,
	NewDeliveryPickUpStoreHandler,
	NewDeliveryExpressTemplateHandler,
	NewTradeOrderHandler,
	NewHandlers,
	brokerage.ProviderSet,
)

type Handlers struct {
	AfterSale               *TradeAfterSaleHandler
	Config                  *TradeConfigHandler
	DeliveryExpress         *DeliveryExpressHandler
	DeliveryPickUpStore     *DeliveryPickUpStoreHandler
	DeliveryExpressTemplate *DeliveryExpressTemplateHandler
	Order                   *TradeOrderHandler
	Brokerage               *brokerage.Handlers
}

func NewHandlers(
	afterSale *TradeAfterSaleHandler,
	config *TradeConfigHandler,
	deliveryExpress *DeliveryExpressHandler,
	deliveryPickUpStore *DeliveryPickUpStoreHandler,
	deliveryExpressTemplate *DeliveryExpressTemplateHandler,
	order *TradeOrderHandler,
	brokerageHandlers *brokerage.Handlers,
) *Handlers {
	return &Handlers{
		AfterSale:               afterSale,
		Config:                  config,
		DeliveryExpress:         deliveryExpress,
		DeliveryPickUpStore:     deliveryPickUpStore,
		DeliveryExpressTemplate: deliveryExpressTemplate,
		Order:                   order,
		Brokerage:               brokerageHandlers,
	}
}
