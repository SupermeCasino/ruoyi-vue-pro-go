package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	bzErr "github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"go.uber.org/zap"
)

// ========== SMS 验证码场景常量 ==========

// SmsSceneEnum 短信验证码场景枚举
type SmsSceneEnum struct {
	Scene        int32
	TemplateCode string
	Description  string
}

// SMS 验证码场景定义（对齐 Java 版本）
var (
	// 会员场景
	SmsSceneMemberLogin       = SmsSceneEnum{1, "user-sms-login", "会员用户 - 手机号登陆"}
	SmsSceneMemberUpdateMob   = SmsSceneEnum{2, "user-update-mobile", "会员用户 - 修改手机"}
	SmsSceneMemberUpdatePwd   = SmsSceneEnum{3, "user-update-password", "会员用户 - 修改密码"}
	SmsSceneMemberResetPwd    = SmsSceneEnum{4, "user-reset-password", "会员用户 - 忘记密码"}

	// 后台用户场景
	SmsSceneAdminLogin        = SmsSceneEnum{21, "admin-sms-login", "后台用户 - 手机号登录"}
	SmsSceneAdminRegister     = SmsSceneEnum{22, "admin-sms-register", "后台用户 - 手机号注册"}
	SmsSceneAdminResetPwd     = SmsSceneEnum{23, "admin-reset-password", "后台用户 - 忘记密码"}
)

// SceneMap 场景到枚举的映射
var SceneMap = map[int32]SmsSceneEnum{
	1:  SmsSceneMemberLogin,
	2:  SmsSceneMemberUpdateMob,
	3:  SmsSceneMemberUpdatePwd,
	4:  SmsSceneMemberResetPwd,
	21: SmsSceneAdminLogin,
	22: SmsSceneAdminRegister,
	23: SmsSceneAdminResetPwd,
}

// GetSceneEnum 根据 scene 获取枚举（对齐 Java 的 getCodeByScene）
func GetSceneEnum(scene int32) *SmsSceneEnum {
	if se, ok := SceneMap[scene]; ok {
		return &se
	}
	return nil
}

// ========== SMS 验证码配置常量 ==========

const (
	// SmsCodeExpire 验证码过期时间（10 分钟，对齐 Java）
	SmsCodeExpire = 10 * time.Minute

	// SmsCodeSendFrequency 发送频率限制（1 分钟，对齐 Java）
	SmsCodeSendFrequency = 1 * time.Minute

	// SmsCodeMaxPerDay 每日最大发送数量（对齐 Java）
	SmsCodeMaxPerDay = 10

	// SmsCodeCacheKeyPrefix Redis 缓存 key 前缀
	SmsCodeCacheKeyPrefix = "sms:code:"

	// SmsCodeRateLimitPrefix 发送频率限制 key 前缀
	SmsCodeRateLimitPrefix = "sms:rate:"
)

// ========== 错误码定义（对齐 Java 版本）==========

var (
	// ErrSmsCodeNotFound 验证码不存在
	ErrSmsCodeNotFound = bzErr.NewBizError(1_002_014_000, "验证码不存在")

	// ErrSmsCodeExpired 验证码已过期
	ErrSmsCodeExpired = bzErr.NewBizError(1_002_014_001, "验证码已过期")

	// ErrSmsCodeUsed 验证码已使用
	ErrSmsCodeUsed = bzErr.NewBizError(1_002_014_002, "验证码已使用")

	// ErrSmsCodeSendTooFast 短信发送过于频繁
	ErrSmsCodeSendTooFast = bzErr.NewBizError(1_002_014_005, "短信发送过于频繁，请稍后再试")

	// ErrSmsCodeExceedMaxPerDay 超过每日发送数量限制
	ErrSmsCodeExceedMaxPerDay = bzErr.NewBizError(1_002_014_004, "今日短信发送数量已达上限")

	// ErrSmsSceneInvalid 短信场景无效
	ErrSmsSceneInvalid = bzErr.NewBizError(400, "短信场景无效")
)

// ========== SMS 验证码服务 ==========

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

// SendSmsCode 发送短信验证码（完整版本，对齐 Java）
func (s *SmsCodeService) SendSmsCode(ctx context.Context, mobile string, scene int32, createIp string) error {
	// 1. 验证 scene 有效性
	sceneEnum := GetSceneEnum(scene)
	if sceneEnum == nil {
		return ErrSmsSceneInvalid
	}

	// 2. 检查发送频率（1 分钟内最多发送一次）
	rateLimitKey := fmt.Sprintf("%s%s:%d", SmsCodeRateLimitPrefix, mobile, scene)
	exists, err := s.rdb.Exists(ctx, rateLimitKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return ErrSmsCodeSendTooFast
	}

	// 3. 检查每日发送数量限制
	// 查询最后一条记录
	lastCode, err := s.getLastSmsCode(ctx, mobile, scene)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return err
	}

	var todayIndex int32 = 1
	if lastCode != nil && isToday(lastCode.CreateTime) {
		// 如果今天已发过，检查是否超过限制
		if lastCode.TodayIndex >= int32(SmsCodeMaxPerDay) {
			return ErrSmsCodeExceedMaxPerDay
		}
		todayIndex = lastCode.TodayIndex + 1
	}

	// 4. 生成验证码（对齐 Java：4-6 位数字）
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	// 5. 保存到 Redis（用于快速查询和过期管理）
	key := s.getCacheKey(mobile, scene)
	if err := s.rdb.Set(ctx, key, code, SmsCodeExpire).Err(); err != nil {
		return err
	}

	// 6. 设置频率限制 key
	if err := s.rdb.Set(ctx, rateLimitKey, "1", SmsCodeSendFrequency).Err(); err != nil {
		return err
	}

	// 7. 保存到数据库（完整记录生命周期，对齐 Java）
	smsCode := &model.SystemSmsCode{
		Mobile:     mobile,
		Code:       code,
		Scene:      scene,
		Used:       false,
		TodayIndex: todayIndex,
		CreateIp:   createIp,
	}
	if err := s.q.SystemSmsCode.WithContext(ctx).Create(smsCode); err != nil {
		zap.L().Error("Failed to save SMS code to DB", zap.Error(err))
		// 即使数据库保存失败，也继续发送短信（但记录日志）
	}

	// 8. 发送短信（通过模板编码）
	params := map[string]interface{}{
		"code": code,
	}
	_, err = s.smsSendService.SendSingleSmsToMember(ctx, mobile, 0, sceneEnum.TemplateCode, params)
	if err != nil {
		zap.L().Error("Failed to send SMS code",
			zap.String("mobile", mobile),
			zap.Int32("scene", scene),
			zap.Error(err))
		return err
	}

	return nil
}

// ValidateSmsCode 仅验证验证码（不标记为已使用，对齐 Java）
func (s *SmsCodeService) ValidateSmsCode(ctx context.Context, mobile string, scene int32, code string) error {
	// 1. 查询 Redis 中的验证码
	key := s.getCacheKey(mobile, scene)
	val, err := s.rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return ErrSmsCodeNotFound
	}
	if err != nil {
		return err
	}

	// 2. 验证码比对
	if val != code {
		return ErrSmsCodeNotFound
	}

	// 3. 仅验证，不执行任何修改操作
	return nil
}

// UseSmsCode 验证并标记为已使用（对齐 Java）
func (s *SmsCodeService) UseSmsCode(ctx context.Context, mobile string, scene int32, code string, usedIp string) error {
	// 1. 先验证有效性
	if err := s.ValidateSmsCode(ctx, mobile, scene, code); err != nil {
		return err
	}

	// 2. 从 Redis 删除（一次性使用）
	key := s.getCacheKey(mobile, scene)
	s.rdb.Del(ctx, key)

	// 3. 更新数据库中的最后一条记录为已使用（对齐 Java）
	lastCode, err := s.getLastSmsCode(ctx, mobile, scene)
	if err != nil {
		return err
	}

	now := time.Now()
	_, updateErr := s.q.SystemSmsCode.WithContext(ctx).
		Where(s.q.SystemSmsCode.ID.Eq(lastCode.ID)).
		Updates(map[string]interface{}{
			"used":     true,
			"used_time": now,
			"used_ip":   usedIp,
		})
	if updateErr != nil {
		zap.L().Error("Failed to mark SMS code as used", zap.Error(updateErr))
		// 即使数据库更新失败，也不返回错误（Redis 已删除，一次性使用已生效）
	}

	return nil
}

// ========== 私有辅助方法 ==========

// getCacheKey 获取 Redis 缓存 key
func (s *SmsCodeService) getCacheKey(mobile string, scene int32) string {
	return fmt.Sprintf("%s%s:%d", SmsCodeCacheKeyPrefix, mobile, scene)
}

// getLastSmsCode 获取最后一条短信验证码记录（数据库查询）
func (s *SmsCodeService) getLastSmsCode(ctx context.Context, mobile string, scene int32) (*model.SystemSmsCode, error) {
	l := s.q.SystemSmsCode
	return l.WithContext(ctx).
		Where(l.Mobile.Eq(mobile), l.Scene.Eq(scene)).
		Order(l.ID.Desc()).
		First()
}

// isToday 判断时间是否是今天
func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() &&
		t.Month() == now.Month() &&
		t.Day() == now.Day()
}
