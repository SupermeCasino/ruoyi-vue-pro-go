package service

import (
	"sync"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/sms/client"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/sms/client/aliyun"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/sms/client/debug"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/sms/client/tencent"
)

type SmsClientFactory struct {
	// channelId -> SmsClient
	clients sync.Map
}

func NewSmsClientFactory() *SmsClientFactory {
	return &SmsClientFactory{}
}

func (f *SmsClientFactory) GetClient(channelId int64) client.SmsClient {
	if v, ok := f.clients.Load(channelId); ok {
		return v.(client.SmsClient)
	}
	return nil
}

func (f *SmsClientFactory) CreateOrUpdateClient(channel *model.SystemSmsChannel) {
	c := f.createClient(channel)
	f.clients.Store(channel.ID, c)
}

func (f *SmsClientFactory) createClient(channel *model.SystemSmsChannel) client.SmsClient {
	switch channel.Code {
	case "aliyun":
		return aliyun.NewSmsClient(channel)
	case "tencent":
		return tencent.NewSmsClient(channel)
	case "debug":
		return debug.NewSmsClient(channel)
	default:
		// Fallback to debug
		return debug.NewSmsClient(channel)
	}
}

// InitClients 初始化所有客户端
func (f *SmsClientFactory) InitClients(channels []*model.SystemSmsChannel) {
	for _, channel := range channels {
		f.CreateOrUpdateClient(channel)
	}
}
