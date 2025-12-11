package client

// BaseClient 支付客户端基类
type BaseClient struct {
	ChannelID   int64
	ChannelCode string
	Config      string // JSON Configuration
}

func NewBaseClient(channelID int64, channelCode string, config string) *BaseClient {
	return &BaseClient{
		ChannelID:   channelID,
		ChannelCode: channelCode,
		Config:      config,
	}
}

func (c *BaseClient) GetID() int64 {
	return c.ChannelID
}

func (c *BaseClient) GetConfig() string {
	return c.Config
}

// Common validation logic or helpers can go here
