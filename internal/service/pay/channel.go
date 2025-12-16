package pay

import (
	"context"
	"errors"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client"

	"gorm.io/gorm"
)

type PayChannelService struct {
	q             *query.Query
	clientFactory *client.PayClientFactory
}

func NewPayChannelService(q *query.Query, clientFactory *client.PayClientFactory) *PayChannelService {
	return &PayChannelService{
		q:             q,
		clientFactory: clientFactory,
	}
}

// CreateChannel 创建支付渠道
func (s *PayChannelService) CreateChannel(ctx context.Context, req *req.PayChannelCreateReq) (int64, error) {
	// 1. 校验是否重复 (AppID + Code)
	exists, err := s.GetChannelByAppIdAndCode(ctx, req.AppID, req.Code)
	if err != nil {
		return 0, err
	}
	if exists != nil {
		return 0, core.NewBizError(1006002000, "支付渠道已存在") // PAY_CHANNEL_EXIST_SAME_CHANNEL_ERROR
	}

	// 2. 插入
	channel := &pay.PayChannel{
		Code:    req.Code,
		Status:  req.Status,
		FeeRate: req.FeeRate,
		Remark:  req.Remark,
		AppID:   req.AppID,
		Config:  req.Config,
	}
	err = s.q.PayChannel.WithContext(ctx).Create(channel)
	if err != nil {
		return 0, err
	}
	return channel.ID, nil
}

// UpdateChannel 更新支付渠道
func (s *PayChannelService) UpdateChannel(ctx context.Context, req *req.PayChannelUpdateReq) error {
	// 1. 校验存在
	_, err := s.validateChannelExists(ctx, req.ID)
	if err != nil {
		return err
	}

	// 2. 更新
	_, err = s.q.PayChannel.WithContext(ctx).Where(s.q.PayChannel.ID.Eq(req.ID)).Updates(pay.PayChannel{
		FeeRate: req.FeeRate,
		Remark:  req.Remark,
		Config:  req.Config,
	})
	return err
}

// DeleteChannel 删除支付渠道
func (s *PayChannelService) DeleteChannel(ctx context.Context, id int64) error {
	// 1. 校验存在
	if _, err := s.validateChannelExists(ctx, id); err != nil {
		return err
	}
	// 2. 删除
	_, err := s.q.PayChannel.WithContext(ctx).Where(s.q.PayChannel.ID.Eq(id)).Delete()
	return err
}

// GetChannel 获得支付渠道
func (s *PayChannelService) GetChannel(ctx context.Context, id int64) (*pay.PayChannel, error) {
	return s.q.PayChannel.WithContext(ctx).Where(s.q.PayChannel.ID.Eq(id)).First()
}

// GetChannelListByAppIds 根据 AppID 集合获支付渠道列表
func (s *PayChannelService) GetChannelListByAppIds(ctx context.Context, appIds []int64) ([]*pay.PayChannel, error) {
	if len(appIds) == 0 {
		return nil, nil
	}
	return s.q.PayChannel.WithContext(ctx).Where(s.q.PayChannel.AppID.In(appIds...)).Find()
}

// Private Methods

// GetChannelByAppIdAndCode 根据 AppID 和 Code 获得支付渠道
func (s *PayChannelService) GetChannelByAppIdAndCode(ctx context.Context, appId int64, code string) (*pay.PayChannel, error) {
	return s.q.PayChannel.WithContext(ctx).Where(s.q.PayChannel.AppID.Eq(appId), s.q.PayChannel.Code.Eq(code)).First()
}

// GetEnableChannelList 获得指定应用的开启的支付渠道列表
func (s *PayChannelService) GetEnableChannelList(ctx context.Context, appId int64) ([]*pay.PayChannel, error) {
	return s.q.PayChannel.WithContext(ctx).
		Where(s.q.PayChannel.AppID.Eq(appId), s.q.PayChannel.Status.Eq(0)). // 0 = Enabled
		Find()
}

func (s *PayChannelService) validateChannelExists(ctx context.Context, id int64) (*pay.PayChannel, error) {
	channel, err := s.q.PayChannel.WithContext(ctx).Where(s.q.PayChannel.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.NewBizError(1006002002, "支付渠道不存在") // PAY_CHANNEL_NOT_FOUND
		}
		return nil, err
	}
	return channel, nil
}

// ValidPayChannel 校验支付渠道是否有效
func (s *PayChannelService) ValidPayChannel(ctx context.Context, id int64) (*pay.PayChannel, error) {
	channel, err := s.validateChannelExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if channel.Status != 0 { // 0 = Enabled
		return nil, core.NewBizError(1006002001, "支付渠道处于关闭状态") // PAY_CHANNEL_IS_DISABLE
	}
	return channel, nil
}

// GetPayClient 获得支付客户端
// 对齐 Java: PayChannelService.getPayClient(Long id)
func (s *PayChannelService) GetPayClient(channelID int64) client.PayClient {
	return s.clientFactory.GetPayClient(channelID)
}
