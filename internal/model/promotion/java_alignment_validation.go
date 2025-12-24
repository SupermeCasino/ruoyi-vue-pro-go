package promotion

import (
	"fmt"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// JavaAlignmentValidation 验证Go常量与Java枚举值的对齐情况
type JavaAlignmentValidation struct {
	ConstantName string
	GoValue      int
	JavaValue    int
	JavaEnum     string
	IsAligned    bool
}

// ValidateJavaAlignment 验证所有常量与Java版本的对齐情况
func ValidateJavaAlignment() []JavaAlignmentValidation {
	var validations []JavaAlignmentValidation

	// 验证优惠券状态常量 (对齐 Java CouponStatusEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"CouponStatusUnused", CouponStatusUnused, 1, "CouponStatusEnum.UNUSED", CouponStatusUnused == 1},
		{"CouponStatusUsed", CouponStatusUsed, 2, "CouponStatusEnum.USED", CouponStatusUsed == 2},
		{"CouponStatusExpired", CouponStatusExpired, 3, "CouponStatusEnum.EXPIRE", CouponStatusExpired == 3},
	}...)

	// 验证优惠券领取类型常量 (对齐 Java CouponTakeTypeEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"CouponTakeTypeUser", CouponTakeTypeUser, 1, "CouponTakeTypeEnum.USER", CouponTakeTypeUser == 1},
		{"CouponTakeTypeAdmin", CouponTakeTypeAdmin, 2, "CouponTakeTypeEnum.ADMIN", CouponTakeTypeAdmin == 2},
		{"CouponTakeTypeRegister", CouponTakeTypeRegister, 3, "CouponTakeTypeEnum.REGISTER", CouponTakeTypeRegister == 3},
	}...)

	// 验证优惠券有效期类型常量 (对齐 Java CouponTemplateValidityTypeEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"CouponValidityTypeDate", CouponValidityTypeDate, 1, "CouponTemplateValidityTypeEnum.DATE", CouponValidityTypeDate == 1},
		{"CouponValidityTypeTerm", CouponValidityTypeTerm, 2, "CouponTemplateValidityTypeEnum.TERM", CouponValidityTypeTerm == 2},
	}...)

	// 验证商品范围常量 (对齐 Java PromotionProductScopeEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"ProductScopeAll", ProductScopeAll, 1, "PromotionProductScopeEnum.ALL", ProductScopeAll == 1},
		{"ProductScopeSpu", ProductScopeSpu, 2, "PromotionProductScopeEnum.SPU", ProductScopeSpu == 2},
		{"ProductScopeCategory", ProductScopeCategory, 3, "PromotionProductScopeEnum.CATEGORY", ProductScopeCategory == 3},
	}...)

	// 验证折扣类型常量 (对齐 Java PromotionDiscountTypeEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"DiscountTypePrice", DiscountTypePrice, 1, "PromotionDiscountTypeEnum.PRICE", DiscountTypePrice == 1},
		{"DiscountTypePercent", DiscountTypePercent, 2, "PromotionDiscountTypeEnum.PERCENT", DiscountTypePercent == 2},
	}...)

	// 验证活动状态常量 (对齐 Java PromotionActivityStatusEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"ActivityStatusWait", ActivityStatusWait, 10, "PromotionActivityStatusEnum.WAIT", ActivityStatusWait == 10},
		{"ActivityStatusRun", ActivityStatusRun, 20, "PromotionActivityStatusEnum.RUN", ActivityStatusRun == 20},
		{"ActivityStatusEnd", ActivityStatusEnd, 30, "PromotionActivityStatusEnum.END", ActivityStatusEnd == 30},
		{"ActivityStatusClose", ActivityStatusClose, 40, "PromotionActivityStatusEnum.CLOSE", ActivityStatusClose == 40},
	}...)

	// 验证通用状态常量 (对齐 Java CommonStatusEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"CommonStatusEnable", model.CommonStatusEnable, 0, "CommonStatusEnum.ENABLE", model.CommonStatusEnable == 0},
		{"CommonStatusDisable", model.CommonStatusDisable, 1, "CommonStatusEnum.DISABLE", model.CommonStatusDisable == 1},
	}...)

	// 验证营销类型常量 (对齐 Java PromotionTypeEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"PromotionTypeSeckillActivity", PromotionTypeSeckillActivity, 1, "PromotionTypeEnum.SECKILL_ACTIVITY", PromotionTypeSeckillActivity == 1},
		{"PromotionTypeBargainActivity", PromotionTypeBargainActivity, 2, "PromotionTypeEnum.BARGAIN_ACTIVITY", PromotionTypeBargainActivity == 2},
		{"PromotionTypeCombinationActivity", PromotionTypeCombinationActivity, 3, "PromotionTypeEnum.COMBINATION_ACTIVITY", PromotionTypeCombinationActivity == 3},
		{"PromotionTypeDiscountActivity", PromotionTypeDiscountActivity, 4, "PromotionTypeEnum.DISCOUNT_ACTIVITY", PromotionTypeDiscountActivity == 4},
		{"PromotionTypeRewardActivity", PromotionTypeRewardActivity, 5, "PromotionTypeEnum.REWARD_ACTIVITY", PromotionTypeRewardActivity == 5},
		{"PromotionTypeMemberLevel", PromotionTypeMemberLevel, 6, "PromotionTypeEnum.MEMBER_LEVEL", PromotionTypeMemberLevel == 6},
		{"PromotionTypeCoupon", PromotionTypeCoupon, 7, "PromotionTypeEnum.COUPON", PromotionTypeCoupon == 7},
		{"PromotionTypePoint", PromotionTypePoint, 8, "PromotionTypeEnum.POINT", PromotionTypePoint == 8},
	}...)

	// 验证拼团记录状态常量 (对齐 Java PromotionCombinationRecordStatusEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"CombinationRecordStatusInProgress", model.PromotionCombinationRecordStatusInProgress, 0, "PromotionCombinationRecordStatusEnum.IN_PROGRESS", model.PromotionCombinationRecordStatusInProgress == 0},
		{"CombinationRecordStatusSuccess", model.PromotionCombinationRecordStatusSuccess, 1, "PromotionCombinationRecordStatusEnum.SUCCESS", model.PromotionCombinationRecordStatusSuccess == 1},
		{"CombinationRecordStatusFailed", model.PromotionCombinationRecordStatusFailed, 2, "PromotionCombinationRecordStatusEnum.FAILED", model.PromotionCombinationRecordStatusFailed == 2},
	}...)

	// 验证砍价记录状态常量 (对齐 Java BargainRecordStatusEnum)
	validations = append(validations, []JavaAlignmentValidation{
		{"BargainRecordStatusInProgress", model.BargainRecordStatusInProgress, 0, "BargainRecordStatusEnum.IN_PROGRESS", model.BargainRecordStatusInProgress == 0},
		{"BargainRecordStatusSuccess", model.BargainRecordStatusSuccess, 1, "BargainRecordStatusEnum.SUCCESS", model.BargainRecordStatusSuccess == 1},
		{"BargainRecordStatusFailed", model.BargainRecordStatusFailed, 2, "BargainRecordStatusEnum.FAILED", model.BargainRecordStatusFailed == 2},
	}...)

	return validations
}

// PrintJavaAlignmentReport 打印Java对齐验证报告
func PrintJavaAlignmentReport() {
	validations := ValidateJavaAlignment()

	fmt.Println("=== Java 常量对齐验证报告 ===")
	fmt.Println()

	alignedCount := 0
	misalignedCount := 0

	for _, validation := range validations {
		if validation.IsAligned {
			alignedCount++
			fmt.Printf("✅ %s = %d (Java: %s = %d) - 对齐\n",
				validation.ConstantName, validation.GoValue, validation.JavaEnum, validation.JavaValue)
		} else {
			misalignedCount++
			fmt.Printf("❌ %s = %d (Java: %s = %d) - 不对齐\n",
				validation.ConstantName, validation.GoValue, validation.JavaEnum, validation.JavaValue)
		}
	}

	fmt.Println()
	fmt.Printf("总计: %d 个常量\n", len(validations))
	fmt.Printf("对齐: %d 个\n", alignedCount)
	fmt.Printf("不对齐: %d 个\n", misalignedCount)

	if misalignedCount == 0 {
		fmt.Println("🎉 所有常量都与Java版本完全对齐！")
	} else {
		fmt.Printf("⚠️  发现 %d 个常量与Java版本不对齐，需要修复\n", misalignedCount)
	}
}

// GetMisalignedConstants 获取所有不对齐的常量
func GetMisalignedConstants() []JavaAlignmentValidation {
	validations := ValidateJavaAlignment()
	var misaligned []JavaAlignmentValidation

	for _, validation := range validations {
		if !validation.IsAligned {
			misaligned = append(misaligned, validation)
		}
	}

	return misaligned
}

// IsAllConstantsAligned 检查是否所有常量都与Java版本对齐
func IsAllConstantsAligned() bool {
	return len(GetMisalignedConstants()) == 0
}
