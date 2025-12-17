package router

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler"
	memberApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/member"
	productApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/product"
	promotionApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/promotion"
	tradeApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/trade"
	appBrokerage "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAppRoutes 注册 App 端路由
func RegisterAppRoutes(engine *gin.Engine,
	// System
	tenantHandler *handler.TenantHandler,
	// Member
	appAuthHandler *memberApp.AppAuthHandler,
	appMemberUserHandler *memberApp.AppMemberUserHandler,
	appMemberAddressHandler *memberApp.AppMemberAddressHandler,
	appMemberPointRecordHandler *memberApp.AppMemberPointRecordHandler,
	appMemberSignInRecordHandler *memberApp.AppMemberSignInRecordHandler,
	// Product
	appProductCategoryHandler *productApp.AppCategoryHandler,
	appProductFavoriteHandler *productApp.AppProductFavoriteHandler,
	appProductBrowseHistoryHandler *productApp.AppProductBrowseHistoryHandler,
	appProductSpuHandler *productApp.AppProductSpuHandler,
	appProductCommentHandler *productApp.AppProductCommentHandler,
	// Trade
	appCartHandler *tradeApp.AppCartHandler,
	appTradeOrderHandler *tradeApp.AppTradeOrderHandler,
	appTradeAfterSaleHandler *tradeApp.AppTradeAfterSaleHandler,
	appTradeConfigHandler *tradeApp.AppTradeConfigHandler,
	// Promotion
	appCouponHandler *promotionApp.AppCouponHandler,
	appBannerHandler *promotionApp.AppBannerHandler,
	appArticleHandler *promotionApp.AppArticleHandler,
	// DIY
	appDiyPageHandler *promotionApp.AppDiyPageHandler,
	appDiyTemplateHandler *promotionApp.AppDiyTemplateHandler,
	// Kefu
	appKefuHandler *promotionApp.AppKefuHandler,
	appCombinationActivityHandler *promotionApp.AppCombinationActivityHandler,
	appCombinationRecordHandler *promotionApp.AppCombinationRecordHandler,
	appBargainActivityHandler *promotionApp.AppBargainActivityHandler,
	appBargainRecordHandler *promotionApp.AppBargainRecordHandler,
	appBargainHelpHandler *promotionApp.AppBargainHelpHandler,
	// Brokerage
	appBrokerageUserHandler *appBrokerage.AppBrokerageUserHandler,
	appBrokerageRecordHandler *appBrokerage.AppBrokerageRecordHandler,
	appBrokerageWithdrawHandler *appBrokerage.AppBrokerageWithdrawHandler,
) {
	appGroup := engine.Group("/app-api")
	{
		// ========== System ==========
		systemGroup := appGroup.Group("/system")
		{
			// Tenant
			tenantGroup := systemGroup.Group("/tenant")
			{
				tenantGroup.GET("/get-by-website", tenantHandler.GetTenantByWebsite)
			}
		}

		// ========== Member ==========
		memberGroup := appGroup.Group("/member")
		{
			// Auth
			authGroup := memberGroup.Group("/auth")
			{
				authGroup.POST("/login", appAuthHandler.Login)
				authGroup.POST("/sms-login", appAuthHandler.SmsLogin)
				authGroup.POST("/social-login", appAuthHandler.SocialLogin)
				authGroup.POST("/send-sms-code", appAuthHandler.SendSmsCode)
				authGroup.POST("/validate-sms-code", appAuthHandler.ValidateSmsCode)
				authGroup.POST("/logout", appAuthHandler.Logout)
				authGroup.POST("/refresh-token", appAuthHandler.RefreshToken)
			}

			// User (Auth Required)
			userGroup := memberGroup.Group("/user")
			userGroup.Use(middleware.Auth())
			{
				userGroup.GET("/get", appMemberUserHandler.GetUserInfo)
				userGroup.PUT("/update", appMemberUserHandler.UpdateUser)
				userGroup.PUT("/update-mobile", appMemberUserHandler.UpdateUserMobile)
				userGroup.PUT("/update-password", appMemberUserHandler.UpdateUserPassword)
			}

			// User (Public)
			userPublicGroup := memberGroup.Group("/user")
			{
				userPublicGroup.PUT("/reset-password", appMemberUserHandler.ResetUserPassword)
			}

			// Address (Auth Required)
			addressGroup := memberGroup.Group("/address")
			addressGroup.Use(middleware.Auth())
			{
				addressGroup.POST("/create", appMemberAddressHandler.CreateAddress)
				addressGroup.PUT("/update", appMemberAddressHandler.UpdateAddress)
				addressGroup.DELETE("/delete", appMemberAddressHandler.DeleteAddress)
				addressGroup.GET("/get", appMemberAddressHandler.GetAddress)
				addressGroup.GET("/get-default", appMemberAddressHandler.GetDefaultUserAddress)
				addressGroup.GET("/list", appMemberAddressHandler.GetAddressList)
			}

			// Point Record (Auth Required)
			pointRecordGroup := memberGroup.Group("/point/record")
			pointRecordGroup.Use(middleware.Auth())
			{
				pointRecordGroup.GET("/page", appMemberPointRecordHandler.GetPointRecordPage)
			}

			// Sign-in Record (App)
			signInGroup := memberGroup.Group("/sign-in/record")
			{
				signInGroup.GET("/get-summary", appMemberSignInRecordHandler.GetSignInRecordSummary)
				signInGroup.POST("/create", appMemberSignInRecordHandler.CreateSignInRecord)
				signInGroup.GET("/page", appMemberSignInRecordHandler.GetSignInRecordPage)
			}
		}

		// ========== Product ==========
		productGroup := appGroup.Group("/product")
		{
			// Category
			categoryGroup := productGroup.Group("/category")
			{
				categoryGroup.GET("/list", appProductCategoryHandler.GetCategoryList)
				categoryGroup.GET("/list-by-ids", appProductCategoryHandler.GetCategoryListByIds)
			}

			// Favorite (Auth Required)
			favoriteGroup := productGroup.Group("/favorite")
			favoriteGroup.Use(middleware.Auth())
			{
				favoriteGroup.POST("/create", appProductFavoriteHandler.CreateFavorite)
				favoriteGroup.DELETE("/delete", appProductFavoriteHandler.DeleteFavorite)
				favoriteGroup.GET("/page", appProductFavoriteHandler.GetFavoritePage)
				favoriteGroup.GET("/exits", appProductFavoriteHandler.IsFavoriteExists)
				favoriteGroup.GET("/get-count", appProductFavoriteHandler.GetFavoriteCount)
			}

			// Browse History (Auth Required)
			browseHistoryGroup := productGroup.Group("/browse-history")
			browseHistoryGroup.Use(middleware.Auth())
			{
				browseHistoryGroup.DELETE("/delete", appProductBrowseHistoryHandler.DeleteBrowseHistory)
				browseHistoryGroup.DELETE("/clean", appProductBrowseHistoryHandler.CleanBrowseHistory)
				browseHistoryGroup.GET("/page", appProductBrowseHistoryHandler.GetBrowseHistoryPage)
			}

			// SPU
			spuGroup := productGroup.Group("/spu")
			{
				spuGroup.GET("/get-detail", appProductSpuHandler.GetSpuDetail)
				spuGroup.GET("/page", appProductSpuHandler.GetSpuPage)
				spuGroup.GET("/list-by-ids", appProductSpuHandler.GetSpuList)
			}

			// Comment
			commentGroup := productGroup.Group("/comment")
			{
				commentGroup.GET("/page", appProductCommentHandler.GetCommentPage)
				commentGroup.POST("/create", middleware.Auth(), appProductCommentHandler.CreateComment)
			}
		}

		// ========== Trade ==========
		tradeGroup := appGroup.Group("/trade")
		tradeGroup.Use(middleware.Auth())
		{
			// Cart
			cartGroup := tradeGroup.Group("/cart")
			{
				cartGroup.POST("/add", appCartHandler.AddCart)
				cartGroup.PUT("/update-count", appCartHandler.UpdateCartCount)
				cartGroup.PUT("/update-selected", appCartHandler.UpdateCartSelected)
				cartGroup.PUT("/reset", appCartHandler.ResetCart)
				cartGroup.DELETE("/delete", appCartHandler.DeleteCart)
				cartGroup.GET("/get-count", appCartHandler.GetCartCount)
				cartGroup.GET("/list", appCartHandler.GetCartList)
			}

			// Order
			orderGroup := tradeGroup.Group("/order")
			{
				orderGroup.GET("/settlement", appTradeOrderHandler.SettlementOrder)
				orderGroup.GET("/settlement-product", appTradeOrderHandler.SettlementProduct)
				orderGroup.POST("/create", appTradeOrderHandler.CreateOrder)
				orderGroup.GET("/get-detail", appTradeOrderHandler.GetOrderDetail)
				orderGroup.GET("/item/get", appTradeOrderHandler.GetOrderItem)
				orderGroup.GET("/page", appTradeOrderHandler.GetOrderPage)
				orderGroup.GET("/get-count", appTradeOrderHandler.GetOrderCount)
				orderGroup.PUT("/receive", appTradeOrderHandler.ReceiveOrder)
				orderGroup.DELETE("/cancel", appTradeOrderHandler.CancelOrder)
				orderGroup.GET("/get-express-track-list", appTradeOrderHandler.GetOrderExpressTrackList)
			}

			// AfterSale
			afterSaleGroup := tradeGroup.Group("/after-sale")
			{
				afterSaleGroup.POST("/create", appTradeAfterSaleHandler.CreateAfterSale)
				afterSaleGroup.GET("/page", appTradeAfterSaleHandler.GetAfterSalePage)
				afterSaleGroup.GET("/get", appTradeAfterSaleHandler.GetAfterSale)
				afterSaleGroup.DELETE("/cancel", appTradeAfterSaleHandler.CancelAfterSale)
				afterSaleGroup.POST("/delivery", appTradeAfterSaleHandler.DeliveryAfterSale)
			}

			// Brokerage User
			brokerageUserGroup := tradeGroup.Group("/brokerage-user")
			{
				brokerageUserGroup.GET("/get", appBrokerageUserHandler.GetBrokerageUser)
				brokerageUserGroup.GET("/get-summary", appBrokerageUserHandler.GetBrokerageUserSummary)
				brokerageUserGroup.GET("/child-summary-page", appBrokerageUserHandler.GetBrokerageUserChildSummaryPage)
				brokerageUserGroup.PUT("/bind", appBrokerageUserHandler.BindBrokerageUser)
			}
			brokerageRecordGroup := tradeGroup.Group("/brokerage-record")
			{
				brokerageRecordGroup.GET("/page", appBrokerageRecordHandler.GetBrokerageRecordPage)
				brokerageRecordGroup.GET("/get-product-brokerage-price", appBrokerageRecordHandler.GetProductBrokeragePrice)
			}
			brokerageWithdrawGroup := tradeGroup.Group("/brokerage-withdraw")
			{
				brokerageWithdrawGroup.GET("/page", appBrokerageWithdrawHandler.GetBrokerageWithdrawPage)
				brokerageWithdrawGroup.GET("/get", appBrokerageWithdrawHandler.GetBrokerageWithdraw) // GET /get?id=...
				brokerageWithdrawGroup.POST("/create", appBrokerageWithdrawHandler.CreateBrokerageWithdraw)
			}
		}

		// ========== Trade (Public) ==========
		tradePublicGroup := appGroup.Group("/trade")
		{
			orderPublicGroup := tradePublicGroup.Group("/order")
			{
				orderPublicGroup.POST("/update-paid", appTradeOrderHandler.UpdateOrderPaid)
			}
		}

		// Trade Config (Public)
		tradeConfigGroup := appGroup.Group("/trade/config")
		{
			tradeConfigGroup.GET("/get", appTradeConfigHandler.GetTradeConfig)
		}

		// ========== Promotion ==========
		promotionGroup := appGroup.Group("/promotion")
		{
			// Coupon (Auth Required)
			couponGroup := promotionGroup.Group("/coupon")
			couponGroup.Use(middleware.Auth())
			{
				couponGroup.POST("/take", appCouponHandler.TakeCoupon)
				couponGroup.GET("/page", appCouponHandler.GetCouponPage)
				couponGroup.POST("/match-list", appCouponHandler.GetCouponMatchList)
			}

			// Banner (Public)
			engine.GET("/app-api/promotion/banner/list", appBannerHandler.GetBannerList)

			// Article (Public)
			articleGroup := promotionGroup.Group("/article")
			{
				articleGroup.GET("/list-category", appArticleHandler.GetArticleCategoryList)
				articleGroup.GET("/page", appArticleHandler.GetArticlePage)
				articleGroup.GET("/get", appArticleHandler.GetArticle)
			}

			// DIY Page (Public)			// DIY
			diyTemplateGroup := promotionGroup.Group("/diy-template")
			{
				diyTemplateGroup.GET("/used", appDiyTemplateHandler.GetUsedDiyTemplate)
				diyTemplateGroup.GET("/get", appDiyTemplateHandler.GetDiyTemplate)
			}
			diyPageGroup := promotionGroup.Group("/diy-page")
			{
				diyPageGroup.GET("/get", appDiyPageHandler.GetDiyPage)
			}

			// Kefu Message
			kefuMessageGroup := promotionGroup.Group("/kefu-message")
			{
				kefuMessageGroup.POST("/send", appKefuHandler.SendMessage)
				kefuMessageGroup.PUT("/update-read-status", appKefuHandler.UpdateMessageReadStatus)
				kefuMessageGroup.GET("/list", appKefuHandler.GetMessageList)
			}

			// Combination Activity
			combinationActivityGroup := promotionGroup.Group("/combination-activity")
			{
				combinationActivityGroup.GET("/list-by-ids", appCombinationActivityHandler.GetCombinationActivityListByIds)
				combinationActivityGroup.GET("/get-detail", appCombinationActivityHandler.GetCombinationActivityDetail)
				combinationActivityGroup.GET("/page", appCombinationActivityHandler.GetCombinationActivityPage)
			}

			// Combination Record
			combinationRecordGroup := promotionGroup.Group("/combination-record")
			{
				combinationRecordGroup.GET("/get-summary", appCombinationRecordHandler.GetCombinationRecordSummary)
				combinationRecordGroup.GET("/get-head-list", appCombinationRecordHandler.GetHeadCombinationRecordList)
				combinationRecordGroup.GET("/get-detail", appCombinationRecordHandler.GetCombinationRecordDetail)
				combinationRecordGroup.Use(middleware.Auth())
				combinationRecordGroup.GET("/page", appCombinationRecordHandler.GetCombinationRecordPage)
			}

			// Bargain Activity (Public)
			bargainActivityGroup := promotionGroup.Group("/bargain-activity")
			{
				bargainActivityGroup.GET("/list", appBargainActivityHandler.GetBargainActivityList)
				bargainActivityGroup.GET("/page", appBargainActivityHandler.GetBargainActivityPage)
				bargainActivityGroup.GET("/get-detail", appBargainActivityHandler.GetBargainActivityDetail)
			}

			// Bargain Record
			bargainRecordGroup := promotionGroup.Group("/bargain-record")
			{
				bargainRecordGroup.GET("/get-summary", appBargainRecordHandler.GetBargainRecordSummary)
				bargainRecordGroup.GET("/get-detail", appBargainRecordHandler.GetBargainRecordDetail)
				// Auth Required
				bargainRecordGroup.POST("/create", middleware.Auth(), appBargainRecordHandler.CreateBargainRecord)
			}

			// Bargain Help
			bargainHelpGroup := promotionGroup.Group("/bargain-help")
			{
				bargainHelpGroup.GET("/list", appBargainHelpHandler.GetBargainHelpList)
				bargainHelpGroup.POST("/create", middleware.Auth(), appBargainHelpHandler.CreateBargainHelp)
			}
		}
	}
}
