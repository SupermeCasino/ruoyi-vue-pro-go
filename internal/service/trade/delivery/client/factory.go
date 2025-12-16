package client

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/config"
)

type ExpressClientFactory interface {
	GetDefaultExpressClient() ExpressClient
	GetOrCreateExpressClient(client string) ExpressClient
}

type ExpressClientFactoryImpl struct {
	conf config.TradeConfig
}

func NewExpressClientFactory() *ExpressClientFactoryImpl {
	return &ExpressClientFactoryImpl{conf: config.C.Trade}
}

func (f *ExpressClientFactoryImpl) GetDefaultExpressClient() ExpressClient {
	return f.GetOrCreateExpressClient(f.conf.Express.Client)
}

func (f *ExpressClientFactoryImpl) GetOrCreateExpressClient(client string) ExpressClient {
	switch client {
	case "kd100":
		return NewKd100ExpressClient(f.conf.Express.Kd100)
	default:
		return nil
	}
}
