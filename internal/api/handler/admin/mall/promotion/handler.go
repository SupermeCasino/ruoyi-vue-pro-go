package promotion

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewArticleHandler,
	NewArticleCategoryHandler,
	NewBannerHandler,
	NewBargainActivityHandler,
	NewBargainHelpHandler,
	NewBargainRecordHandler,
	NewCombinationActivityHandler,
	NewCombinationRecordHandler,
	NewCouponHandler,
	NewDiscountActivityHandler,
	NewDiyPageHandler,
	NewDiyTemplateHandler,
	NewKefuHandler,
	NewPointActivityHandler,
	NewRewardActivityHandler,
	NewSeckillActivityHandler,
	NewSeckillConfigHandler,
	NewHandlers,
)

type Handlers struct {
	Article             *ArticleHandler
	ArticleCategory     *ArticleCategoryHandler
	Banner              *BannerHandler
	BargainActivity     *BargainActivityHandler
	BargainHelp         *BargainHelpHandler
	BargainRecord       *BargainRecordHandler
	CombinationActivity *CombinationActivityHandler
	CombinationRecord   *CombinationRecordHandler
	Coupon              *CouponHandler
	DiscountActivity    *DiscountActivityHandler
	DiyPage             *DiyPageHandler
	DiyTemplate         *DiyTemplateHandler
	Kefu                *KefuHandler
	PointActivity       *PointActivityHandler
	RewardActivity      *RewardActivityHandler
	SeckillActivity     *SeckillActivityHandler
	SeckillConfig       *SeckillConfigHandler
}

func NewHandlers(
	article *ArticleHandler,
	articleCategory *ArticleCategoryHandler,
	banner *BannerHandler,
	bargainActivity *BargainActivityHandler,
	bargainHelp *BargainHelpHandler,
	bargainRecord *BargainRecordHandler,
	combinationActivity *CombinationActivityHandler,
	combinationRecord *CombinationRecordHandler,
	coupon *CouponHandler,
	discountActivity *DiscountActivityHandler,
	diyPage *DiyPageHandler,
	diyTemplate *DiyTemplateHandler,
	kefu *KefuHandler,
	pointActivity *PointActivityHandler,
	rewardActivity *RewardActivityHandler,
	seckillActivity *SeckillActivityHandler,
	seckillConfig *SeckillConfigHandler,
) *Handlers {
	return &Handlers{
		Article:             article,
		ArticleCategory:     articleCategory,
		Banner:              banner,
		BargainActivity:     bargainActivity,
		BargainHelp:         bargainHelp,
		BargainRecord:       bargainRecord,
		CombinationActivity: combinationActivity,
		CombinationRecord:   combinationRecord,
		Coupon:              coupon,
		DiscountActivity:    discountActivity,
		DiyPage:             diyPage,
		DiyTemplate:         diyTemplate,
		Kefu:                kefu,
		PointActivity:       pointActivity,
		RewardActivity:      rewardActivity,
		SeckillActivity:     seckillActivity,
		SeckillConfig:       seckillConfig,
	}
}
