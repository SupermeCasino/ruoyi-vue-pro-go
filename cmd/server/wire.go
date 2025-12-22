//go:build wireinject
// +build wireinject

package main

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler"
	adminHandler "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin" // Statistics Handlers
	memberAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/member"
	payAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay"
	payWallet "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay/wallet"
	productHandler "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/product"
	promotionAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/promotion"
	tradeAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/trade"
	brokerage "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/trade/brokerage"
	memberHandler "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/member"
	payApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/pay"
	productApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/product"
	promotionApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/promotion"
	tradeApp "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/trade"
	appBrokerage "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/router"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/datascope"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/permission"
	ws "github.com/wxlbd/ruoyi-mall-go/internal/pkg/websocket"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo" // Pay Repo
	payRepo "github.com/wxlbd/ruoyi-mall-go/internal/repo/pay"
	tradeRepo "github.com/wxlbd/ruoyi-mall-go/internal/repo/trade"

	productRepo "github.com/wxlbd/ruoyi-mall-go/internal/repo/product" // Product Statistics Repo
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client"
	_ "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client/alipay"
	_ "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client/weixin"
	payWalletSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/wallet"
	product "github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	promotionSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	tradeSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	tradeBrokerageSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/trade/brokerage"
	deliveryClient "github.com/wxlbd/ruoyi-mall-go/internal/service/trade/delivery/client"

	"github.com/wxlbd/ruoyi-mall-go/pkg/cache"
	"github.com/wxlbd/ruoyi-mall-go/pkg/database"
	"github.com/wxlbd/ruoyi-mall-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitApp() (*gin.Engine, error) {
	wire.Build(
		database.InitDB,
		cache.InitRedis,
		logger.NewLogger,
		// Repo (GORM Gen)
		repo.NewQuery,
		// Data Scope Plugin (需要在Services之后注册)
		datascope.RegisterPlugin,
		// Service
		service.NewOAuth2TokenService,
		service.NewAuthService,
		service.NewMenuService,
		service.NewRoleService,
		service.NewPermissionService,
		service.NewTenantService,
		service.NewUserService,
		service.NewDictService,
		service.NewDeptService,
		service.NewPostService,
		service.NewNoticeService,
		service.NewConfigService,
		service.NewSmsClientFactory,            // Added SmsClientFactory
		service.NewSmsChannelService,           // Added SmsChannelService
		service.NewSmsTemplateService,          // Added SmsTemplateService
		service.NewSmsLogService,               // Added SmsLogService
		service.NewSmsSendService,              // Added SmsSendService
		service.NewFileConfigService,           // Added FileConfigService
		service.NewFileService,                 // Added FileService
		service.NewSmsCodeService,              // Added SmsCodeService
		service.NewLoginLogService,             // Added LoginLogService
		service.NewOperateLogService,           // Added OperateLogService
		service.NewScheduler,                   // Added Scheduler
		service.NewJobService,                  // Added JobService
		service.NewJobLogService,               // Added JobLogService
		service.NewApiAccessLogService,         // Added ApiAccessLogService
		service.NewApiErrorLogService,          // Added ApiErrorLogService
		service.NewSocialClientService,         // Added SocialClientService
		service.NewSocialUserService,           // Added SocialUserService
		service.NewSensitiveWordService,        // Added SensitiveWordService
		service.NewMailService,                 // Added MailService
		service.NewNotifyService,               // Added NotifyService
		service.NewOAuth2ClientService,         // Added OAuth2ClientService
		memberSvc.NewMemberAuthService,         // Added MemberAuthService
		memberSvc.NewMemberUserService,         // Added MemberUserService
		memberSvc.NewMemberAddressService,      // Added MemberAddressService
		memberSvc.NewMemberLevelService,        // Added MemberLevelService
		memberSvc.NewMemberGroupService,        // Added MemberGroupService
		memberSvc.NewMemberTagService,          // Added MemberTagService
		memberSvc.NewMemberConfigService,       // Added MemberConfigService
		memberSvc.NewMemberPointRecordService,  // Added MemberPointRecordService
		product.NewProductCategoryService,      // Added ProductCategoryService
		product.NewProductPropertyService,      // Added ProductPropertyService
		product.NewProductPropertyValueService, // Added ProductPropertyValueService
		product.NewProductBrandService,         // Added ProductBrandService
		product.NewProductSkuService,           // Added ProductSkuService
		product.NewProductSpuService,           // Added ProductSpuService
		product.NewProductCommentService,       // Added ProductCommentService
		product.NewProductFavoriteService,      // Added ProductFavoriteService
		product.NewProductBrowseHistoryService, // Added ProductBrowseHistoryService
		// Member Sign-in
		memberSvc.NewMemberSignInConfigService,
		memberSvc.NewMemberSignInRecordService,

		// Handler
		handler.NewAuthHandler,
		handler.NewUserHandler,
		handler.NewTenantHandler,
		handler.NewDictHandler,
		handler.NewDeptHandler,
		handler.NewPostHandler,
		handler.NewRoleHandler,
		handler.NewMenuHandler, // Added MenuHandler
		handler.NewPermissionHandler,
		handler.NewNoticeHandler,
		handler.NewConfigHandler,
		handler.NewSmsChannelHandler,    // Added SmsChannelHandler
		handler.NewSmsTemplateHandler,   // Added SmsTemplateHandler
		handler.NewSmsLogHandler,        // Added SmsLogHandler
		handler.NewFileConfigHandler,    // Added FileConfigHandler
		handler.NewFileHandler,          // Added FileHandler
		memberHandler.NewAppAuthHandler, // Added AppAuthHandler
		handler.NewLoginLogHandler,      // Added LoginLogHandler
		handler.NewOperateLogHandler,    // Added OperateLogHandler
		handler.NewJobHandler,           // Added JobHandler
		handler.NewJobLogHandler,        // Added JobLogHandler
		handler.NewApiAccessLogHandler,  // Added ApiAccessLogHandler
		handler.NewApiErrorLogHandler,   // Added ApiErrorLogHandler
		handler.NewSocialClientHandler,  // Added SocialClientHandler
		handler.NewSocialUserHandler,    // Added SocialUserHandler
		handler.NewSensitiveWordHandler, // Added SensitiveWordHandler
		handler.NewMailHandler,          // Added MailHandler
		handler.NewNotifyHandler,        // Added NotifyHandler
		handler.NewOAuth2ClientHandler,  // Added OAuth2ClientHandler
		// Member
		memberAdmin.NewMemberLevelHandler,             // Added MemberLevelHandler for admin
		memberAdmin.NewMemberGroupHandler,             // Added MemberGroupHandler for admin
		memberAdmin.NewMemberTagHandler,               // Added MemberTagHandler for admin
		memberAdmin.NewMemberConfigHandler,            // Added MemberConfigHandler for admin
		memberAdmin.NewMemberPointRecordHandler,       // Added MemberPointRecordHandler for admin
		memberAdmin.NewMemberSignInConfigHandler,      // Added MemberSignInConfigHandler
		memberAdmin.NewMemberSignInRecordHandler,      // Added MemberSignInRecordHandler
		memberAdmin.NewMemberUserHandler,              // Added MemberUserHandler
		memberHandler.NewAppMemberUserHandler,         // Added AppMemberUserHandler
		memberHandler.NewAppMemberAddressHandler,      // Added AppMemberAddressHandler
		memberHandler.NewAppMemberPointRecordHandler,  // Added AppMemberPointRecordHandler
		memberHandler.NewAppMemberSignInRecordHandler, // Added AppMemberSignInRecordHandler
		memberHandler.NewAppSocialUserHandler,         // Added AppSocialUserHandler
		productHandler.NewProductCategoryHandler,      // Added ProductCategoryHandler
		productHandler.NewProductPropertyHandler,      // Added ProductPropertyHandler
		productHandler.NewProductBrandHandler,         // Added ProductBrandHandler
		productHandler.NewProductSpuHandler,           // Added ProductSpuHandler
		productHandler.NewProductCommentHandler,
		productHandler.NewProductFavoriteHandler,
		productHandler.NewProductBrowseHistoryHandler,

		// App handlers
		productApp.NewAppCategoryHandler,
		productApp.NewAppProductFavoriteHandler,
		productApp.NewAppProductBrowseHistoryHandler,
		productApp.NewAppProductSpuHandler,
		productApp.NewAppProductCommentHandler,
		// Trade
		tradeSvc.NewCartService,
		tradeSvc.NewTradeOrderQueryService,
		tradeSvc.NewTradePriceService,
		tradeSvc.NewTradeOrderUpdateService,
		tradeSvc.NewTradeAfterSaleService,
		tradeSvc.NewAfterSaleLogService,  // Added
		tradeSvc.NewTradeConfigService,   // Added Config
		tradeSvc.NewTradeOrderLogService, // Added Log
		tradeApp.NewAppCartHandler,
		tradeApp.NewAppTradeOrderHandler,
		tradeApp.NewAppTradeAfterSaleHandler,
		tradeApp.NewAppTradeConfigHandler, // Added Config
		tradeAdmin.NewTradeOrderHandler,
		tradeAdmin.NewTradeAfterSaleHandler,
		tradeAdmin.NewTradeConfigHandler, // Added Config
		// Delivery
		tradeSvc.NewDeliveryExpressService,
		tradeSvc.NewDeliveryPickUpStoreService,
		tradeSvc.NewDeliveryExpressTemplateService,
		tradeAdmin.NewDeliveryExpressHandler,
		tradeAdmin.NewDeliveryPickUpStoreHandler,
		tradeAdmin.NewDeliveryExpressTemplateHandler,
		tradeBrokerageSvc.NewBrokerageUserService,
		// Brokerage
		tradeBrokerageSvc.NewBrokerageRecordService,
		tradeBrokerageSvc.NewBrokerageWithdrawService, // Added

		// Wallet Services
		payWalletSvc.NewPayWalletService,
		payWalletSvc.NewPayWalletRechargeService,
		payWalletSvc.NewPayWalletRechargePackageService,
		payWalletSvc.NewPayWalletTransactionService,
		brokerage.NewBrokerageUserHandler,
		brokerage.NewBrokerageRecordHandler,
		brokerage.NewBrokerageWithdrawHandler,
		appBrokerage.NewAppBrokerageUserHandler, // Added
		appBrokerage.NewAppBrokerageRecordHandler,
		appBrokerage.NewAppBrokerageWithdrawHandler,

		// Pay Repositories
		payRepo.NewPayTransferRepository,
		payRepo.NewPayNoRedisDAO,
		// Trade Repositories
		tradeRepo.NewTradeNoRedisDAO,
		// Statistics
		repo.NewTradeStatisticsRepository,
		repo.NewTradeOrderStatisticsRepository,
		repo.NewTradeOrderLogRepository, // Added Log Repo
		repo.NewAfterSaleLogRepository,  // Added AfterSale Log Repo
		repo.NewAfterSaleStatisticsRepository,
		repo.NewBrokerageStatisticsRepository,
		repo.NewMemberStatisticsRepository,
		repo.NewApiAccessLogStatisticsRepository,
		repo.NewPayWalletStatisticsRepository,
		productRepo.NewProductStatisticsRepository, // Product
		service.NewProductStatisticsService,        // Product
		service.NewTradeStatisticsService,          // Trade
		service.NewTradeOrderStatisticsServiceV2,
		service.NewAfterSaleStatisticsService,
		service.NewBrokerageStatisticsService,
		service.NewMemberStatisticsService,
		service.NewApiAccessLogStatisticsService,
		service.NewPayWalletStatisticsService,
		adminHandler.NewTradeStatisticsHandler,
		adminHandler.NewProductStatisticsHandler,
		adminHandler.NewMemberStatisticsHandler,
		service.NewPayTransferSyncJob, // Added PayTransferSyncJob
		service.NewPayNotifyJob,       // Added PayNotifyJob
		service.NewPayOrderSyncJob,    // Added PayOrderSyncJob
		service.NewPayOrderExpireJob,  // Added PayOrderExpireJob
		service.NewPayRefundSyncJob,   // Added PayRefundSyncJob
		adminHandler.NewPayStatisticsHandler,

		// Promotion
		promotionSvc.NewCouponService,
		promotionSvc.NewCouponUserService,
		promotionSvc.NewPromotionBannerService,     // Added Banner
		promotionSvc.NewRewardActivityService,      // Added Activity
		promotionSvc.NewSeckillConfigService,       // Added Seckill Config
		promotionSvc.NewSeckillActivityService,     // Added Seckill Activity
		promotionSvc.NewBargainActivityService,     // Added Bargain Activity
		promotionSvc.NewBargainRecordService,       // Added Bargain Record
		promotionSvc.NewBargainHelpService,         // Added Bargain Help
		promotionSvc.NewCombinationActivityService, // Added Combination Activity
		promotionSvc.NewCombinationRecordService,   // Added Combination Record
		promotionSvc.NewDiscountActivityService,    // Added Discount Activity
		promotionSvc.NewArticleCategoryService,     // Added Article Category
		promotionSvc.NewArticleService,             // Added Article
		promotionSvc.NewDiyTemplateService,         // Added Diy Template
		promotionSvc.NewDiyPageService,             // Added Diy Page
		promotionSvc.NewPointActivityService, promotionSvc.NewKefuService,
		promotionAdmin.NewCouponHandler,
		promotionAdmin.NewBannerHandler,              // Added Banner
		promotionAdmin.NewRewardActivityHandler,      // Added Activity
		promotionAdmin.NewSeckillConfigHandler,       // Added Seckill Config
		promotionAdmin.NewSeckillActivityHandler,     // Added Seckill Activity
		promotionAdmin.NewBargainActivityHandler,     // Added Bargain Activity
		promotionAdmin.NewCombinationActivityHandler, // Added Combination Activity
		promotionAdmin.NewDiscountActivityHandler,    // Added Discount Activity
		promotionAdmin.NewArticleCategoryHandler,     // Added Article Category
		promotionAdmin.NewArticleHandler,             // Added Article
		promotionAdmin.NewDiyTemplateHandler,         // Added Diy Template
		promotionAdmin.NewDiyPageHandler,             // Added Diy Page
		promotionAdmin.NewPointActivityHandler, promotionAdmin.NewKefuHandler,
		promotionAdmin.NewBargainRecordHandler,
		promotionAdmin.NewCombinationRecordHandler,
		promotionAdmin.NewBargainHelpHandler, promotionApp.NewAppKefuHandler,
		promotionApp.NewAppCouponHandler,
		promotionApp.NewAppCouponTemplateHandler, // 新增
		promotionApp.NewAppBannerHandler,         // Added Banner
		promotionApp.NewAppBargainActivityHandler,
		promotionApp.NewAppBargainRecordHandler,
		promotionApp.NewAppBargainHelpHandler,
		promotionApp.NewAppSeckillActivityHandler,     // 新增
		promotionApp.NewAppSeckillConfigHandler,       // 新增
		promotionApp.NewAppCombinationActivityHandler, // Added Combination Activity
		promotionApp.NewAppCombinationRecordHandler,   // Added Combination Record
		promotionApp.NewAppArticleHandler,             // Added Article
		promotionApp.NewAppDiyPageHandler,             // Added Diy Page
		promotionApp.NewAppDiyTemplateHandler,
		promotionApp.NewAppActivityHandler,
		// WebSocket
		ws.ProviderSet,
		handler.NewWebSocketHandler,

		// Casbin
		permission.InitEnforcer,
		middleware.NewCasbinMiddleware,

		// Pay
		paySvc.NewPayAppService,
		paySvc.NewPayChannelService,
		paySvc.NewPayOrderService,
		paySvc.NewPayRefundService,
		paySvc.NewPayNotifyService,
		paySvc.NewPayTransferService,
		client.NewPayClientFactory,

		deliveryClient.NewExpressClientFactory, // Added ExpressClientFactory
		wire.Bind(new(deliveryClient.ExpressClientFactory), new(*deliveryClient.ExpressClientFactoryImpl)),

		payAdmin.NewPayAppHandler,
		payAdmin.NewPayChannelHandler,
		payAdmin.NewPayOrderHandler,
		payAdmin.NewPayRefundHandler,
		payAdmin.NewPayNotifyHandler,
		payAdmin.NewPayTransferHandler,
		// Wallet Handlers
		payWallet.NewPayWalletHandler,
		payWallet.NewPayWalletRechargeHandler,
		payWallet.NewPayWalletRechargePackageHandler,
		payWallet.NewPayWalletTransactionHandler,
		payApp.NewAppPayOrderHandler,
		payApp.NewAppPayWalletHandler,
		payApp.NewAppPayChannelHandler,
		payApp.NewAppPayTransferHandler,
		payApp.NewAppPayWalletTransactionHandler,
		payApp.NewAppPayWalletRechargePackageHandler,

		// Router
		router.InitRouter,

		// Job Handlers (Slice Injection for Scheduler)
		service.ProvideJobHandlers,
	)
	return &gin.Engine{}, nil
}
