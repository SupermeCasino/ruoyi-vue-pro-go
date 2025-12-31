package promotion

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAppActivityHandler,
	NewAppArticleHandler,
	NewAppBannerHandler,
	NewAppBargainActivityHandler,
	NewAppBargainHelpHandler,
	NewAppBargainRecordHandler,
	NewAppCombinationActivityHandler,
	NewAppCombinationRecordHandler,
	NewAppCouponHandler,
	NewAppCouponTemplateHandler,
	NewAppDiyPageHandler,
	NewAppDiyTemplateHandler,
	NewAppKefuHandler,
	NewAppPointActivityHandler,
	NewAppRewardActivityHandler,
	NewAppSeckillActivityHandler,
	NewAppSeckillConfigHandler,
	NewHandlers,
)

type Handlers struct {
	Activity            *AppActivityHandler
	Article             *AppArticleHandler
	Banner              *AppBannerHandler
	BargainActivity     *AppBargainActivityHandler
	BargainHelp         *AppBargainHelpHandler
	BargainRecord       *AppBargainRecordHandler
	CombinationActivity *AppCombinationActivityHandler
	CombinationRecord   *AppCombinationRecordHandler
	Coupon              *AppCouponHandler
	CouponTemplate      *AppCouponTemplateHandler
	DiyPage             *AppDiyPageHandler
	DiyTemplate         *AppDiyTemplateHandler
	Kefu                *AppKefuHandler
	PointActivity       *AppPointActivityHandler
	RewardActivity      *AppRewardActivityHandler
	SeckillActivity     *AppSeckillActivityHandler
	SeckillConfig       *AppSeckillConfigHandler
}

func NewHandlers(
	activity *AppActivityHandler,
	article *AppArticleHandler,
	banner *AppBannerHandler,
	bargainActivity *AppBargainActivityHandler,
	bargainHelp *AppBargainHelpHandler,
	bargainRecord *AppBargainRecordHandler,
	combinationActivity *AppCombinationActivityHandler,
	combinationRecord *AppCombinationRecordHandler,
	coupon *AppCouponHandler,
	couponTemplate *AppCouponTemplateHandler,
	diyPage *AppDiyPageHandler,
	diyTemplate *AppDiyTemplateHandler,
	kefu *AppKefuHandler,
	pointActivity *AppPointActivityHandler,
	rewardActivity *AppRewardActivityHandler,
	seckillActivity *AppSeckillActivityHandler,
	seckillConfig *AppSeckillConfigHandler,
) *Handlers {
	return &Handlers{
		Activity:            activity,
		Article:             article,
		Banner:              banner,
		BargainActivity:     bargainActivity,
		BargainHelp:         bargainHelp,
		BargainRecord:       bargainRecord,
		CombinationActivity: combinationActivity,
		CombinationRecord:   combinationRecord,
		Coupon:              coupon,
		CouponTemplate:      couponTemplate,
		DiyPage:             diyPage,
		DiyTemplate:         diyTemplate,
		Kefu:                kefu,
		PointActivity:       pointActivity,
		RewardActivity:      rewardActivity,
		SeckillActivity:     seckillActivity,
		SeckillConfig:       seckillConfig,
	}
}
