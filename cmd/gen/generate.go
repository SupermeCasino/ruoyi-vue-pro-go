package main

import (
	"backend-go/internal/model"
	"backend-go/internal/model/member"
	"backend-go/internal/model/pay"
	"backend-go/internal/model/product"
	"backend-go/internal/model/promotion" // Added system import
	"backend-go/internal/model/trade"
	"backend-go/internal/model/trade/brokerage" // Added brokerage import

	"gorm.io/gen"
)

func main() {
	// 1. 不需要连接数据库，直接基于 Struct 生成

	// 2. 配置生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./internal/repo/query",
		ModelPkgPath:  "./internal/model",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// 3. 基于 Struct 生成
	g.ApplyBasic(
		// System
		model.SystemUser{},
		model.SystemRole{},
		model.SystemMenu{},
		model.SystemTenant{},
		model.SystemDictData{},
		model.SystemDictType{},
		model.SystemDept{},
		model.SystemPost{},
		model.SystemUserPost{},
		model.SystemNotice{},
		model.SystemUserRole{},
		model.SystemRoleMenu{},
		model.SystemConfig{},
		model.SystemSmsChannel{},
		model.SystemSmsTemplate{},
		model.SystemSmsLog{},
		model.InfraFileConfig{},
		model.InfraFile{},
		model.InfraFileConfig{},
		model.InfraFile{},
		model.SocialUser{},
		model.SocialUserBind{},
		model.SocialClient{},
		model.SystemLoginLog{},
		model.SystemOperateLog{},
		model.InfraJob{},
		model.InfraJobLog{},
		model.InfraApiAccessLog{},
		model.InfraApiErrorLog{},
		// Member
		member.MemberUser{},
		member.MemberAddress{},
		// Product
		product.ProductCategory{},
		product.ProductProperty{},
		product.ProductPropertyValue{},
		product.ProductBrand{},
		product.ProductSpu{},
		product.ProductSku{},
		&product.ProductComment{},
		&product.ProductFavorite{},
		&product.ProductBrowseHistory{},
	)
	// Trade
	g.ApplyBasic(
		trade.Cart{},
		// Trade
		trade.TradeOrder{},
		trade.TradeOrderItem{},
		trade.AfterSale{},
		trade.TradeConfig{},
		trade.TradeOrderLog{},
		trade.TradeStatistics{}, // 统计

		// DeliverySale{},
		trade.TradeDeliveryExpress{},
		trade.TradeDeliveryPickUpStore{},
		trade.TradeDeliveryExpressTemplate{},
		trade.TradeDeliveryExpressTemplateCharge{},
		trade.TradeDeliveryExpressTemplateFree{},
		// Brokerage
		brokerage.BrokerageUser{},
		brokerage.BrokerageRecord{},
		brokerage.BrokerageWithdraw{},
	)
	// Statistics
	g.ApplyBasic(
		product.ProductStatistics{}, // 商品统计
	)
	// Promotion DIY
	g.ApplyBasic(
		promotion.PromotionDiyTemplate{},
		promotion.PromotionDiyPage{},
	)
	// Promotion Kefu
	g.ApplyBasic(promotion.PromotionKefuConversation{}, promotion.PromotionKefuMessage{})
	// Promotion
	g.ApplyBasic(
		promotion.PromotionCouponTemplate{},
		promotion.PromotionCouponTemplate{},
		promotion.PromotionCoupon{},
		promotion.PromotionSeckillActivity{},     // Added
		promotion.PromotionSeckillProduct{},      // Added
		promotion.PromotionSeckillConfig{},       // Added
		promotion.PromotionBargainActivity{},     // Added Bargain
		promotion.PromotionBargainRecord{},       // Added Bargain
		promotion.PromotionBargainHelp{},         // Added Bargain
		promotion.PromotionCombinationActivity{}, // Added Combination
		promotion.PromotionCombinationProduct{},  // Added Combination
		promotion.PromotionCombinationRecord{},   // Added Combination
		promotion.PromotionDiscountActivity{},    // Added Discount
		promotion.PromotionDiscountProduct{},     // Added Discount
		promotion.PromotionArticle{},             // Added Article
		promotion.PromotionArticleCategory{},     // Added Article
		promotion.PromotionBanner{},
		promotion.PromotionRewardActivity{},
		promotion.PromotionPointActivity{},
		promotion.PromotionPointProduct{},
		member.MemberLevel{},
		member.MemberGroup{},
		member.MemberTag{},
		member.MemberConfig{},
		member.MemberPointRecord{},
		member.MemberSignInConfig{},
		member.MemberSignInRecord{},
	)

	// Pay Models
	g.ApplyBasic(
		pay.PayApp{},
		pay.PayChannel{},
		pay.PayOrder{},
		pay.PayOrderExtension{},
		pay.PayRefund{},
		pay.PayNotifyTask{},
		pay.PayNotifyLog{},
		pay.PayWallet{},
		pay.PayWalletRecharge{},
		pay.PayWalletTransaction{},
		pay.PayWalletRechargePackage{},
	)

	// 4. 执行生成
	g.Execute()
}
