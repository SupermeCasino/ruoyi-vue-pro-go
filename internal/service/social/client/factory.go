package client

import (
	"context"
	"fmt"

	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

type SocialPlatformFactoryImpl struct {
	q *query.Query
}

func NewSocialPlatformFactory(q *query.Query) *SocialPlatformFactoryImpl {
	return &SocialPlatformFactoryImpl{q: q}
}

func (f *SocialPlatformFactoryImpl) GetPlatform(ctx context.Context, socialType int, userType int) (SocialPlatform, error) {
	// 1. 查询 SocialClient 配置
	client, err := f.q.SocialClient.WithContext(ctx).
		Where(f.q.SocialClient.SocialType.Eq(socialType), f.q.SocialClient.UserType.Eq(userType), f.q.SocialClient.Status.Eq(0)). // 0 = Enable
		First()
	if err != nil {
		return nil, core.NewBizError(1002004001, fmt.Sprintf("社交客户端不存在: type=%d", socialType))
	}

	// 2. 根据类型创建客户端
	switch socialType {
	case 31, 30: // WeChat Mini / MP
		return NewWeChatClient(client), nil
	// case 20: DingTalk
	// ...
	default:
		return nil, core.NewBizError(1002004002, fmt.Sprintf("不支持的社交类型: %d", socialType))
	}
}
