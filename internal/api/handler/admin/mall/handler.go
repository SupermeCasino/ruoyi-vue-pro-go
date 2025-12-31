package mall

import (
	"github.com/google/wire"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall/trade"
)

var ProviderSet = wire.NewSet(
	NewHandlers,
	product.ProviderSet,
	promotion.ProviderSet,
	trade.ProviderSet,
)

type Handlers struct {
	Product   *product.Handlers
	Promotion *promotion.Handlers
	Trade     *trade.Handlers
}

func NewHandlers(
	product *product.Handlers,
	promotion *promotion.Handlers,
	trade *trade.Handlers,
) *Handlers {
	return &Handlers{
		Product:   product,
		Promotion: promotion,
		Trade:     trade,
	}
}
