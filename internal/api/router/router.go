package router

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler"
	adminHandler "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin" // Statistics Handlers
	memberAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/member"
	payAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay"
	payWallet "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay/wallet"
	productHandler "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/product"
	promotionAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/promotion"
	tradeAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/trade"
	tradeBrokerageAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/trade/brokerage"
	memberHandler "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/member"
	payApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/pay"
	productApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/product"
	promotionApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/promotion"
	tradeApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/trade"
	appBrokerage "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/trade/brokerage"

	"fmt"

	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/datascope"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, rdb *redis.Client,
	_ *datascope.PluginRegistered, // 确保Plugin在Router初始化前已注册
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tenantHandler *handler.TenantHandler,
	dictHandler *handler.DictHandler,
	deptHandler *handler.DeptHandler,
	postHandler *handler.PostHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	permissionHandler *handler.PermissionHandler,
	noticeHandler *handler.NoticeHandler,
	configHandler *handler.ConfigHandler,
	smsChannelHandler *handler.SmsChannelHandler,
	smsTemplateHandler *handler.SmsTemplateHandler,
	smsLogHandler *handler.SmsLogHandler,
	fileConfigHandler *handler.FileConfigHandler,
	fileHandler *handler.FileHandler,
	appAuthHandler *memberHandler.AppAuthHandler,
	appMemberUserHandler *memberHandler.AppMemberUserHandler,
	appMemberAddressHandler *memberHandler.AppMemberAddressHandler,
	productCategoryHandler *productHandler.ProductCategoryHandler,
	productPropertyHandler *productHandler.ProductPropertyHandler,
	productBrandHandler *productHandler.ProductBrandHandler,
	productSpuHandler *productHandler.ProductSpuHandler,
	productCommentHandler *productHandler.ProductCommentHandler,
	productFavoriteHandler *productHandler.ProductFavoriteHandler,
	productBrowseHistoryHandler *productHandler.ProductBrowseHistoryHandler,
	appProductCategoryHandler *productApp.AppCategoryHandler,
	appProductFavoriteHandler *productApp.AppProductFavoriteHandler,
	appProductBrowseHistoryHandler *productApp.AppProductBrowseHistoryHandler,
	appProductSpuHandler *productApp.AppProductSpuHandler,
	appProductCommentHandler *productApp.AppProductCommentHandler,
	appCartHandler *tradeApp.AppCartHandler,
	tradeOrderHandler *tradeAdmin.TradeOrderHandler,
	appTradeOrderHandler *tradeApp.AppTradeOrderHandler,
	tradeAfterSaleHandler *tradeAdmin.TradeAfterSaleHandler,
	appTradeAfterSaleHandler *tradeApp.AppTradeAfterSaleHandler,
// Promotion
	couponHandler *promotionAdmin.CouponHandler,
	combinationActivityHandler *promotionAdmin.CombinationActivityHandler,
	discountActivityHandler *promotionAdmin.DiscountActivityHandler,
	appCombinationActivityHandler *promotionApp.AppCombinationActivityHandler,
	appCombinationRecordHandler *promotionApp.AppCombinationRecordHandler,
	appCouponHandler *promotionApp.AppCouponHandler,
	appCouponTemplateHandler *promotionApp.AppCouponTemplateHandler, // 新增
	deliveryExpressHandler *tradeAdmin.DeliveryExpressHandler,
	deliveryPickUpStoreHandler *tradeAdmin.DeliveryPickUpStoreHandler,
	deliveryExpressTemplateHandler *tradeAdmin.DeliveryExpressTemplateHandler,
	bannerHandler *promotionAdmin.BannerHandler,
	rewardActivityHandler *promotionAdmin.RewardActivityHandler,
	seckillConfigHandler *promotionAdmin.SeckillConfigHandler,
	seckillActivityHandler *promotionAdmin.SeckillActivityHandler,
	bargainActivityHandler *promotionAdmin.BargainActivityHandler,
	appBannerHandler *promotionApp.AppBannerHandler,
	memberLevelHandler *memberAdmin.MemberLevelHandler,
	memberGroupHandler *memberAdmin.MemberGroupHandler,
	memberTagHandler *memberAdmin.MemberTagHandler,
	memberConfigHandler *memberAdmin.MemberConfigHandler,
	memberPointRecordHandler *memberAdmin.MemberPointRecordHandler,
	appMemberPointRecordHandler *memberHandler.AppMemberPointRecordHandler,
	memberSignInConfigHandler *memberAdmin.MemberSignInConfigHandler,
	memberSignInRecordHandler *memberAdmin.MemberSignInRecordHandler,
	appMemberSignInRecordHandler *memberHandler.AppMemberSignInRecordHandler,
	appSocialUserHandler *memberHandler.AppSocialUserHandler,
	memberUserHandler *memberAdmin.MemberUserHandler,
	payAppHandler *payAdmin.PayAppHandler,
	payChannelHandler *payAdmin.PayChannelHandler,
	payOrderHandler *payAdmin.PayOrderHandler,
	payRefundHandler *payAdmin.PayRefundHandler,
	payNotifyHandler *payAdmin.PayNotifyHandler,
	payTransferHandler *payAdmin.PayTransferHandler,
// Wallet
	payWalletHandler *payWallet.PayWalletHandler,
	payWalletRechargeHandler *payWallet.PayWalletRechargeHandler,
	payWalletRechargePackageHandler *payWallet.PayWalletRechargePackageHandler,
	payWalletTransactionHandler *payWallet.PayWalletTransactionHandler,
	loginLogHandler *handler.LoginLogHandler,
	operateLogHandler *handler.OperateLogHandler,
	jobHandler *handler.JobHandler,
	jobLogHandler *handler.JobLogHandler,
	apiAccessLogHandler *handler.ApiAccessLogHandler,
	apiErrorLogHandler *handler.ApiErrorLogHandler,
	socialClientHandler *handler.SocialClientHandler,
	socialUserHandler *handler.SocialUserHandler,
	mailHandler *handler.MailHandler,
	notifyHandler *handler.NotifyHandler,
	oauth2ClientHandler *handler.OAuth2ClientHandler, // Added OAuth2ClientHandler
	appBargainActivityHandler *promotionApp.AppBargainActivityHandler,
	appBargainRecordHandler *promotionApp.AppBargainRecordHandler,
	appBargainHelpHandler *promotionApp.AppBargainHelpHandler,
	appSeckillActivityHandler *promotionApp.AppSeckillActivityHandler, // 新增
	appSeckillConfigHandler *promotionApp.AppSeckillConfigHandler,     // 新增
// Article
	articleCategoryHandler *promotionAdmin.ArticleCategoryHandler,
	articleHandler *promotionAdmin.ArticleHandler,
	appArticleHandler *promotionApp.AppArticleHandler,
// DIY
	diyTemplateHandler *promotionAdmin.DiyTemplateHandler,
	diyPageHandler *promotionAdmin.DiyPageHandler,
	appDiyPageHandler *promotionApp.AppDiyPageHandler,
	appDiyTemplateHandler *promotionApp.AppDiyTemplateHandler,
// Kefu
	kefuHandler *promotionAdmin.KefuHandler,
	appKefuHandler *promotionApp.AppKefuHandler,
// Point Activity
	pointActivityHandler *promotionAdmin.PointActivityHandler,
// Record Handlers (Added Phase 3)
	bargainRecordHandler *promotionAdmin.BargainRecordHandler,
	combinationRecordHandler *promotionAdmin.CombinationRecordHandler,
	bargainHelpHandler *promotionAdmin.BargainHelpHandler,
// Trade Config
	tradeConfigHandler *tradeAdmin.TradeConfigHandler,
	appTradeConfigHandler *tradeApp.AppTradeConfigHandler,
	brokerageUserHandler *tradeBrokerageAdmin.BrokerageUserHandler,
	brokerageRecordHandler *tradeBrokerageAdmin.BrokerageRecordHandler,
	brokerageWithdrawHandler *tradeBrokerageAdmin.BrokerageWithdrawHandler,
// Statistics
	tradeStatisticsHandler *adminHandler.TradeStatisticsHandler,
	productStatisticsHandler *adminHandler.ProductStatisticsHandler,
	memberStatisticsHandler *adminHandler.MemberStatisticsHandler,
	payStatisticsHandler *adminHandler.PayStatisticsHandler,
	appBrokerageUserHandler *appBrokerage.AppBrokerageUserHandler,
	appBrokerageRecordHandler *appBrokerage.AppBrokerageRecordHandler,
	appBrokerageWithdrawHandler *appBrokerage.AppBrokerageWithdrawHandler,
	appPayOrderHandler *payApp.AppPayOrderHandler,
	appPayWalletHandler *payApp.AppPayWalletHandler,
	appPayChannelHandler *payApp.AppPayChannelHandler,
	appPayTransferHandler *payApp.AppPayTransferHandler,
	appPayWalletTransactionHandler *payApp.AppPayWalletTransactionHandler,
	appPayWalletRechargePackageHandler *payApp.AppPayWalletRechargePackageHandler,
	appActivityHandler *promotionApp.AppActivityHandler,
	appPointActivityHandler *promotionApp.AppPointActivityHandler,
	webSocketHandler *handler.WebSocketHandler,
	casbinMiddleware *middleware.CasbinMiddleware,
) *gin.Engine {
	// Debug log to confirm router init
	fmt.Println("Initializing Router...")
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}))
	r.Use(gin.Logger())
	// 注入 gin.Context 到 request context，供 GORM Hook 使用
	r.Use(middleware.InjectContext())

	// 基础路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// ========== 模块化路由注册 ==========

	// System 模块 (Auth, Tenant, Dict, Dept, Post, User, Role, Permission, Logs, SMS, File, Infra)
	// System 模块 (Auth, Tenant, Dict, Dept, Post, User, Role, Permission, Logs, SMS, File, Infra)
	// WebSocket (Register at root /infra/ws to match Java path)
	r.GET("/infra/ws", webSocketHandler.Handle)

	RegisterSystemRoutes(r,
		authHandler, userHandler, tenantHandler, dictHandler, deptHandler,
		postHandler, roleHandler, menuHandler, permissionHandler, noticeHandler,
		loginLogHandler, operateLogHandler, configHandler,
		smsChannelHandler, smsTemplateHandler, smsLogHandler,
		fileConfigHandler, fileHandler,
		jobHandler, jobLogHandler, apiAccessLogHandler, apiErrorLogHandler,
		socialClientHandler, socialUserHandler, mailHandler, notifyHandler, oauth2ClientHandler, webSocketHandler,
		casbinMiddleware,
	)

	// Product 模块
	RegisterProductRoutes(r,
		productCategoryHandler, productBrandHandler, productPropertyHandler,
		productSpuHandler, productCommentHandler, productFavoriteHandler,
		productBrowseHistoryHandler,
		casbinMiddleware,
	)

	// Promotion 模块
	RegisterPromotionRoutes(r,
		couponHandler, bannerHandler, rewardActivityHandler,
		seckillConfigHandler, seckillActivityHandler, bargainActivityHandler,
		combinationActivityHandler, discountActivityHandler,
		articleCategoryHandler, articleHandler,
		diyTemplateHandler, diyPageHandler, kefuHandler,
		pointActivityHandler,
		bargainRecordHandler, combinationRecordHandler, bargainHelpHandler,
		casbinMiddleware,
	)

	// Trade 模块
	RegisterTradeRoutes(r,
		tradeOrderHandler, tradeAfterSaleHandler,
		deliveryExpressHandler, deliveryPickUpStoreHandler, deliveryExpressTemplateHandler,
		tradeConfigHandler,
		brokerageUserHandler,
		brokerageRecordHandler,
		brokerageWithdrawHandler,
		casbinMiddleware,
	)

	// Member 模块 (Admin)
	RegisterMemberRoutes(r,
		memberSignInConfigHandler, memberSignInRecordHandler,
		memberPointRecordHandler,
		memberConfigHandler, memberGroupHandler, memberLevelHandler, memberTagHandler,
		memberUserHandler,
		casbinMiddleware,
	)

	// Pay 模块
	RegisterPayRoutes(r,
		payAppHandler, payChannelHandler, payOrderHandler, payRefundHandler, payNotifyHandler, payTransferHandler,
		payWalletHandler, payWalletRechargeHandler, payWalletRechargePackageHandler, payWalletTransactionHandler,
		casbinMiddleware,
	)

	// App 模块 (移动端)
	areaHandler := handler.NewAreaHandler() // 创建 AreaHandler
	RegisterAppRoutes(r,
		// System
		tenantHandler,
		areaHandler,
		// Member
		appAuthHandler, appMemberUserHandler, appMemberAddressHandler,
		appMemberPointRecordHandler, appMemberSignInRecordHandler,
		appSocialUserHandler,
		// Product
		appProductCategoryHandler,
		appProductFavoriteHandler, appProductBrowseHistoryHandler,
		appProductSpuHandler, appProductCommentHandler,
		// Trade
		appCartHandler, appTradeOrderHandler, appTradeAfterSaleHandler, appTradeConfigHandler,
		// Promotion
		appCouponHandler, appCouponTemplateHandler, appBannerHandler, appArticleHandler, // DIY
		appDiyPageHandler,
		appDiyTemplateHandler,
		// Kefu
		appKefuHandler,
		appCombinationActivityHandler, appCombinationRecordHandler,
		appBargainActivityHandler, appBargainRecordHandler, appBargainHelpHandler,
		appSeckillActivityHandler, appSeckillConfigHandler,
		appBrokerageUserHandler,
		appBrokerageRecordHandler,
		appBrokerageWithdrawHandler,
		appPayOrderHandler,
		appPayWalletHandler,
		appPayChannelHandler,
		appPayTransferHandler,
		appPayWalletTransactionHandler,
		appPayWalletRechargePackageHandler,
		appActivityHandler,
		appPointActivityHandler,
	)

	// Statistics 模块
	RegisterStatisticsRoutes(r,
		tradeStatisticsHandler, productStatisticsHandler,
		memberStatisticsHandler, payStatisticsHandler,
		casbinMiddleware,
	)

	return r
}
