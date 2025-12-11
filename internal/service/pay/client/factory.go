package client

import (
	"errors"
	"sync"
)

// PayClientFactory 支付客户端工厂
type PayClientFactory struct {
	clients map[int64]PayClient
	mutex   sync.RWMutex
}

func NewPayClientFactory() *PayClientFactory {
	return &PayClientFactory{
		clients: make(map[int64]PayClient),
	}
}

// GetPayClient 获得支付客户端
func (f *PayClientFactory) GetPayClient(channelID int64) PayClient {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.clients[channelID]
}

// RegisterClient 注册客户端 (用于扩展，避免硬编码 switch)
// Currently simplifying: we will hardcode the switch in CreateOrUpdate for simplicity unless we want a registry pattern.
// Given Go's static nature, a registry or simple switch in a "provider" package implies circular deps if not careful.
// Best approach: Define 'Creator' function type.

type ClientCreator func(channelID int64, config string) (PayClient, error)

var creators = make(map[string]ClientCreator)

func RegisterCreator(channelCode string, creator ClientCreator) {
	creators[channelCode] = creator
}

// CreateOrUpdatePayClient 创建或更新支付客户端
func (f *PayClientFactory) CreateOrUpdatePayClient(channelID int64, channelCode string, config string) (PayClient, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	client := f.clients[channelID]
	if client != nil {
		// Verify if change is needed, implementation dependent. For now, we recreate.
		// Real logic: client.Refresh(config)
	}

	creator, ok := creators[channelCode]
	if !ok {
		// Use "mock" as fallback or specific codes
		if channelCode == "mock" {
			// return new MockClient
		}
		// Try generic creators if needed, or specific prefixes?
		// For now if not found, error
		// return nil, fmt.Errorf("channel code %s not supported", channelCode)

		// Temporary Logic: If no creator found, return logic for Mock?
		// We will implement Alipay/WxPay creators and register them at init.
		return nil, errors.New("channel not supported")
	}

	newClient, err := creator(channelID, config)
	if err != nil {
		return nil, err
	}
	if err := newClient.Init(); err != nil {
		return nil, err
	}
	f.clients[channelID] = newClient
	return newClient, nil
}
