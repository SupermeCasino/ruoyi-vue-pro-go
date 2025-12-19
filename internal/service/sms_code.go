package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"go.uber.org/zap"
)

const (
	SmsCodeCacheKeyPrefix = "sms:code:"
	SmsCodeExpire         = 5 * time.Minute
)

type SmsCodeService struct {
	q              *query.Query
	rdb            *redis.Client
	smsSendService *SmsSendService
}

func NewSmsCodeService(q *query.Query, rdb *redis.Client, smsSendService *SmsSendService) *SmsCodeService {
	return &SmsCodeService{
		q:              q,
		rdb:            rdb,
		smsSendService: smsSendService,
	}
}

// SendSmsCode 发送短信验证码
func (s *SmsCodeService) SendSmsCode(ctx context.Context, mobile string, scene int) error {
	// 1. 校验频率 (1分钟内只能发一次)
	rateLimitKey := fmt.Sprintf("sms:rate:%s:%d", mobile, scene)
	exists, err := s.rdb.Exists(ctx, rateLimitKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.NewBizError(1004003001, "发送过于频繁，请稍后再试")
	}
	// Set rate limit key with 60s expiry
	if err := s.rdb.Set(ctx, rateLimitKey, "1", 60*time.Second).Err(); err != nil {
		return err
	}

	// 2. 生成验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 3. 保存到 Redis
	key := s.getCacheKey(mobile, scene)
	if err := s.rdb.Set(ctx, key, code, SmsCodeExpire).Err(); err != nil {
		return err
	}

	// 4. 发送短信
	params := map[string]interface{}{
		"code": code,
	}
	// TODO: 根据 scene 获取 templateCode
	// 暂时 hardcode 一个测试用 code, 实际应查表或配置
	templateCode := "USER_SMS_LOGIN"

	// 默认发送给 Member
	_, err = s.smsSendService.SendSingleSmsToMember(ctx, mobile, 0, templateCode, params)
	if err != nil {
		zap.L().Error("Send SMS code failed", zap.Error(err))
		return err
	}

	return nil
}

// ValidateSmsCode 校验验证码
func (s *SmsCodeService) ValidateSmsCode(ctx context.Context, mobile string, scene int, code string) error {
	key := s.getCacheKey(mobile, scene)
	val, err := s.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return errors.NewBizError(1004003003, "验证码已过期或不存在")
	}
	if err != nil {
		return err
	}
	if val != code {
		return errors.NewBizError(1004003004, "验证码错误")
	}

	// 验证成功后删除，避免重复使用
	s.rdb.Del(ctx, key)
	return nil
}

func (s *SmsCodeService) getCacheKey(mobile string, scene int) string {
	return fmt.Sprintf("%s%s:%d", SmsCodeCacheKeyPrefix, mobile, scene)
}
