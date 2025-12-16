package pay

import (
	"time"
)

// TODO: Pay module
// Placeholder APIs for Pay module (to be implemented later)

type PayTransferCreateReqDTO struct {
	AppID              int64             `json:"appId"`
	ChannelCode        string            `json:"channelCode"`
	MerchantTransferID string            `json:"merchantTransferId"`
	Subject            string            `json:"subject"`
	Price              int               `json:"price"`
	UserAccount        string            `json:"userAccount"`
	UserName           string            `json:"userName"`
	UserIP             string            `json:"userIp"`
	OpenID             string            `json:"openid"` // Added OpenID for WeChat
	ChannelExtras      map[string]string `json:"channelExtras"`
}

type PayTransferCreateRespDTO struct {
	ID     int64 `json:"id"`
	Status int   `json:"status"`
}

type PayTransferRespDTO struct {
	ID                 int64             `json:"id"`
	Status             int               `json:"status"`
	Price              int               `json:"price"`
	MerchantTransferID string            `json:"merchantTransferId"`
	ChannelCode        string            `json:"channelCode"`
	SuccessTime        *time.Time        `json:"successTime"`
	ChannelErrorMsg    string            `json:"channelErrorMsg"`
	ChannelExtras      map[string]string `json:"channelExtras"`
}

type PayWalletRespDTO struct {
	ID int64 `json:"id"`
}

// PayWalletService (Placeholder)
type PayWalletService struct{}

func NewPayWalletService() *PayWalletService {
	return &PayWalletService{}
}

func (s *PayWalletService) GetOrCreateWallet(userId int64, userType int) (*PayWalletRespDTO, error) {
	return &PayWalletRespDTO{ID: 0}, nil // Mock
}
