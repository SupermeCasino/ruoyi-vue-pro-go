package router

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPromotionRoutes 注册营销活动模块路由
func RegisterPromotionRoutes(engine *gin.Engine,
	handlers *promotion.Handlers,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	promotionGroup := engine.Group("/admin-api/promotion")
	promotionGroup.Use(middleware.Auth())
	{
		// Coupon Template
		templateGroup := promotionGroup.Group("/coupon-template")
		{
			templateGroup.POST("/create", handlers.Coupon.CreateCouponTemplate)
			templateGroup.PUT("/update", handlers.Coupon.UpdateCouponTemplate)
			templateGroup.PUT("/update-status", handlers.Coupon.UpdateCouponTemplateStatus)
			templateGroup.DELETE("/delete", handlers.Coupon.DeleteCouponTemplate)
			templateGroup.GET("/get", handlers.Coupon.GetCouponTemplate)
			templateGroup.GET("/page", handlers.Coupon.GetCouponTemplatePage)
			templateGroup.GET("/list", handlers.Coupon.GetCouponTemplateList)
		}

		// Coupon
		couponGroup := promotionGroup.Group("/coupon")
		{
			couponGroup.DELETE("/delete", handlers.Coupon.DeleteCoupon)
			couponGroup.GET("/page", handlers.Coupon.GetCouponPage)
			couponGroup.POST("/send", handlers.Coupon.SendCoupon)
		}

		// Banner
		bannerGroup := promotionGroup.Group("/banner")
		{
			bannerGroup.POST("/create", handlers.Banner.CreateBanner)
			bannerGroup.PUT("/update", handlers.Banner.UpdateBanner)
			bannerGroup.DELETE("/delete", handlers.Banner.DeleteBanner)
			bannerGroup.GET("/get", handlers.Banner.GetBanner)
			bannerGroup.GET("/page", handlers.Banner.GetBannerPage)
		}

		// Reward Activity
		rewardGroup := promotionGroup.Group("/reward-activity")
		{
			rewardGroup.POST("/create", handlers.RewardActivity.CreateRewardActivity)
			rewardGroup.PUT("/update", handlers.RewardActivity.UpdateRewardActivity)
			rewardGroup.PUT("/close", handlers.RewardActivity.CloseRewardActivity)
			rewardGroup.DELETE("/delete", handlers.RewardActivity.DeleteRewardActivity)
			rewardGroup.GET("/get", handlers.RewardActivity.GetRewardActivity)
			rewardGroup.GET("/page", handlers.RewardActivity.GetRewardActivityPage)
		}

		// Seckill Config
		seckillConfigGroup := promotionGroup.Group("/seckill-config")
		{
			seckillConfigGroup.POST("/create", handlers.SeckillConfig.CreateSeckillConfig)
			seckillConfigGroup.PUT("/update", handlers.SeckillConfig.UpdateSeckillConfig)
			seckillConfigGroup.PUT("/update-status", handlers.SeckillConfig.UpdateSeckillConfigStatus)
			seckillConfigGroup.DELETE("/delete", handlers.SeckillConfig.DeleteSeckillConfig)
			seckillConfigGroup.GET("/get", handlers.SeckillConfig.GetSeckillConfig)
			seckillConfigGroup.GET("/page", handlers.SeckillConfig.GetSeckillConfigPage)
			seckillConfigGroup.GET("/list", handlers.SeckillConfig.GetSeckillConfigList)
			seckillConfigGroup.GET("/simple-list", handlers.SeckillConfig.GetSeckillConfigSimpleList)
		}

		// Seckill Activity
		seckillActivityGroup := promotionGroup.Group("/seckill-activity")
		{
			seckillActivityGroup.POST("/create", handlers.SeckillActivity.CreateSeckillActivity)
			seckillActivityGroup.PUT("/update", handlers.SeckillActivity.UpdateSeckillActivity)
			seckillActivityGroup.DELETE("/delete", handlers.SeckillActivity.DeleteSeckillActivity)
			seckillActivityGroup.PUT("/close", handlers.SeckillActivity.CloseSeckillActivity)
			seckillActivityGroup.GET("/get", handlers.SeckillActivity.GetSeckillActivity)
			seckillActivityGroup.GET("/page", handlers.SeckillActivity.GetSeckillActivityPage)
			seckillActivityGroup.GET("/list-by-ids", handlers.SeckillActivity.GetSeckillActivityListByIds)
		}

		// Bargain Activity
		bargainActivityGroup := promotionGroup.Group("/bargain-activity")
		{
			bargainActivityGroup.POST("/create", handlers.BargainActivity.CreateBargainActivity)
			bargainActivityGroup.PUT("/update", handlers.BargainActivity.UpdateBargainActivity)
			bargainActivityGroup.DELETE("/delete", handlers.BargainActivity.DeleteBargainActivity)
			bargainActivityGroup.PUT("/close", handlers.BargainActivity.CloseBargainActivity)
			bargainActivityGroup.GET("/get", handlers.BargainActivity.GetBargainActivity)
			bargainActivityGroup.GET("/page", handlers.BargainActivity.GetBargainActivityPage)
		}

		// Combination Activity
		combinationActivityGroup := promotionGroup.Group("/combination-activity")
		{
			combinationActivityGroup.POST("/create", handlers.CombinationActivity.CreateCombinationActivity)
			combinationActivityGroup.PUT("/update", handlers.CombinationActivity.UpdateCombinationActivity)
			combinationActivityGroup.PUT("/close", handlers.CombinationActivity.CloseCombinationActivity)
			combinationActivityGroup.DELETE("/delete", handlers.CombinationActivity.DeleteCombinationActivity)
			combinationActivityGroup.GET("/get", handlers.CombinationActivity.GetCombinationActivity)
			combinationActivityGroup.GET("/list-by-ids", handlers.CombinationActivity.GetCombinationActivityListByIds)
			combinationActivityGroup.GET("/page", handlers.CombinationActivity.GetCombinationActivityPage)
		}

		// Discount Activity
		discountActivityGroup := promotionGroup.Group("/discount-activity")
		{
			discountActivityGroup.POST("/create", handlers.DiscountActivity.CreateDiscountActivity)
			discountActivityGroup.PUT("/update", handlers.DiscountActivity.UpdateDiscountActivity)
			discountActivityGroup.POST("/close", handlers.DiscountActivity.CloseDiscountActivity)
			discountActivityGroup.DELETE("/delete", handlers.DiscountActivity.DeleteDiscountActivity)
			discountActivityGroup.GET("/get", handlers.DiscountActivity.GetDiscountActivity)
			discountActivityGroup.GET("/page", handlers.DiscountActivity.GetDiscountActivityPage)
		}

		// Article Category
		articleCategoryGroup := promotionGroup.Group("/article-category")
		{
			articleCategoryGroup.POST("/create", handlers.ArticleCategory.CreateArticleCategory)
			articleCategoryGroup.PUT("/update", handlers.ArticleCategory.UpdateArticleCategory)
			articleCategoryGroup.DELETE("/delete", handlers.ArticleCategory.DeleteArticleCategory)
			articleCategoryGroup.GET("/get", handlers.ArticleCategory.GetArticleCategory)
			articleCategoryGroup.GET("/list", handlers.ArticleCategory.GetArticleCategoryList)
			articleCategoryGroup.GET("/list-all-simple", handlers.ArticleCategory.GetSimpleList)
			articleCategoryGroup.GET("/page", handlers.ArticleCategory.GetArticleCategoryPage)
		}

		// Article
		articleGroup := promotionGroup.Group("/article")
		{
			articleGroup.POST("/create", handlers.Article.CreateArticle)
			articleGroup.PUT("/update", handlers.Article.UpdateArticle)
			articleGroup.DELETE("/delete", handlers.Article.DeleteArticle)
			articleGroup.GET("/get", handlers.Article.GetArticle)
			articleGroup.GET("/page", handlers.Article.GetArticlePage)
		}

		// DIY Template
		diyTemplateGroup := promotionGroup.Group("/diy-template")
		{
			diyTemplateGroup.POST("/create", handlers.DiyTemplate.CreateDiyTemplate)
			diyTemplateGroup.PUT("/update", handlers.DiyTemplate.UpdateDiyTemplate)
			diyTemplateGroup.PUT("/use", handlers.DiyTemplate.UseDiyTemplate)
			diyTemplateGroup.DELETE("/delete", handlers.DiyTemplate.DeleteDiyTemplate)
			diyTemplateGroup.GET("/get", handlers.DiyTemplate.GetDiyTemplate)
			diyTemplateGroup.GET("/page", handlers.DiyTemplate.GetDiyTemplatePage)
			diyTemplateGroup.GET("/get-property", handlers.DiyTemplate.GetDiyTemplateProperty)
			diyTemplateGroup.PUT("/update-property", handlers.DiyTemplate.UpdateDiyTemplateProperty)
		}

		// DIY Page
		diyPageGroup := promotionGroup.Group("/diy-page")
		{
			diyPageGroup.POST("/create", handlers.DiyPage.CreateDiyPage)
			diyPageGroup.PUT("/update", handlers.DiyPage.UpdateDiyPage)
			diyPageGroup.DELETE("/delete", handlers.DiyPage.DeleteDiyPage)
			diyPageGroup.GET("/get", handlers.DiyPage.GetDiyPage)
			diyPageGroup.GET("/list", handlers.DiyPage.GetDiyPageList)
			diyPageGroup.GET("/page", handlers.DiyPage.GetDiyPagePage)
			diyPageGroup.GET("/get-property", handlers.DiyPage.GetDiyPageProperty)
			diyPageGroup.PUT("/update-property", handlers.DiyPage.UpdateDiyPageProperty)
		}

		// Kefu Conversation (Admin)
		kefuConversationGroup := promotionGroup.Group("/kefu-conversation")
		{
			kefuConversationGroup.GET("/get", handlers.Kefu.GetConversation)
			kefuConversationGroup.GET("/list", handlers.Kefu.GetConversationList)
			kefuConversationGroup.PUT("/update-conversation-pinned", handlers.Kefu.UpdateConversationPinned)
			kefuConversationGroup.DELETE("/delete", handlers.Kefu.DeleteConversation)
		}

		// Kefu Message (Admin)
		kefuMessageGroup := promotionGroup.Group("/kefu-message")
		{
			kefuMessageGroup.POST("/send", handlers.Kefu.SendMessage)
			kefuMessageGroup.PUT("/update-read-status", handlers.Kefu.UpdateMessageReadStatus)
			kefuMessageGroup.GET("/list", handlers.Kefu.GetMessageList)
		}

		// Point Activity
		pointActivityGroup := promotionGroup.Group("/point-activity")
		{
			pointActivityGroup.POST("/create", handlers.PointActivity.CreatePointActivity)
			pointActivityGroup.PUT("/update", handlers.PointActivity.UpdatePointActivity)
			pointActivityGroup.PUT("/close", handlers.PointActivity.ClosePointActivity)
			pointActivityGroup.DELETE("/delete", handlers.PointActivity.DeletePointActivity)
			pointActivityGroup.GET("/get", handlers.PointActivity.GetPointActivity)
			pointActivityGroup.GET("/page", handlers.PointActivity.GetPointActivityPage)
			pointActivityGroup.GET("/list-by-ids", handlers.PointActivity.GetPointActivityListByIds)
		}

		// Bargain Record
		bargainRecordGroup := promotionGroup.Group("/bargain-record")
		{
			bargainRecordGroup.GET("/page", handlers.BargainRecord.GetBargainRecordPage)
		}

		// Combination Record
		combinationRecordGroup := promotionGroup.Group("/combination-record")
		{
			combinationRecordGroup.GET("/page", handlers.CombinationRecord.GetCombinationRecordPage)
			combinationRecordGroup.GET("/get-summary", handlers.CombinationRecord.GetCombinationRecordSummary)
		}

		// Bargain Help
		bargainHelpGroup := promotionGroup.Group("/bargain-help")
		{
			bargainHelpGroup.GET("/page", handlers.BargainHelp.GetBargainHelpPage)
		}
	}
}
