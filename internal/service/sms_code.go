package service

import (
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	SmsCodeCacheKeyPrefix = "sms:code:"
	SmsCodeExpire         = 5 * time.Minute
)

type SmsCodeService struct {
	q   *query.Query
	rdb *redis.Client
	// channelSvc *SmsChannelService // If we need to read config for real sending
}

func NewSmsCodeService(q *query.Query, rdb *redis.Client) *SmsCodeService {
	return &SmsCodeService{
		q:   q,
		rdb: rdb,
	}
}

// SendSmsCode 发送短信验证码
func (s *SmsCodeService) SendSmsCode(ctx context.Context, mobile string, scene int) error {
	// 1. 校验频率 (例如 1分钟内只能发一次)
	// TODO: Rate Limit

	// 2. 生成验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 3. 保存到 Redis
	key := s.getCacheKey(mobile, scene)
	if err := s.rdb.Set(ctx, key, code, SmsCodeExpire).Err(); err != nil {
		return err
	}

	// 4. 发送短信 (Mock for now, log it)
	// In real world, query enabled channel from s.q.SystemSmsChannel and use SDK
	zap.L().Info("Send SMS Code", zap.String("mobile", mobile), zap.String("code", code), zap.Int("scene", scene))

	return nil
}

// ValidateSmsCode 校验验证码
func (s *SmsCodeService) ValidateSmsCode(ctx context.Context, mobile string, scene int, code string) error {
	key := s.getCacheKey(mobile, scene)
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return core.NewBizError(1004003003, "验证码已过期或不存在")
	}
	if err != nil {
		return err
	}
	if val != code {
		return core.NewBizError(1004003004, "验证码错误")
	}

	// 验证成功后删除，避免重复使用
	s.rdb.Del(ctx, key)
	return nil
}

func (s *SmsCodeService) getCacheKey(mobile string, scene int) string {
	return fmt.Sprintf("%s%s:%d", SmsCodeCacheKeyPrefix, mobile, scene)
}
