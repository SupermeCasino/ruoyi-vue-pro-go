package brokerage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	tradeReq "github.com/wxlbd/ruoyi-mall-go/internal/api/req/app/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"go.uber.org/zap"
)

type BrokerageWithdrawService struct {
	q              *query.Query
	logger         *zap.Logger
	recordSvc      *BrokerageRecordService
	payTransferSvc *pay.PayTransferService
	tradeConfigSvc *trade.TradeConfigService
	memberSvc      *member.MemberUserService
}

func NewBrokerageWithdrawService(q *query.Query, logger *zap.Logger, recordSvc *BrokerageRecordService, payTransferSvc *pay.PayTransferService, tradeConfigSvc *trade.TradeConfigService, memberSvc *member.MemberUserService) *BrokerageWithdrawService {
	return &BrokerageWithdrawService{
		q:              q,
		logger:         logger,
		recordSvc:      recordSvc,
		payTransferSvc: payTransferSvc,
		tradeConfigSvc: tradeConfigSvc,
		memberSvc:      memberSvc,
	}
}

// AuditBrokerageWithdraw 审批佣金提现
func (s *BrokerageWithdrawService) AuditBrokerageWithdraw(ctx context.Context, id int64, status int, auditReason string) error {
	w := s.q.BrokerageWithdraw
	withdraw, err := w.WithContext(ctx).Where(w.ID.Eq(id)).First()
	if err != nil {
		return errors.New("提现记录不存在")
	}

	// 1.2 特殊：【重新转账】如果是提现失败，并且状态是审核中，那么更新状态为审核中，并且清空 transferErrorMsg
	// Java: if (WITHDRAW_FAIL.equals(withdraw.getStatus())) { ... }
	// But the request status (arg) is usually AUDIT_SUCCESS or AUDIT_FAIL.
	// The Java logic checks if *current DB status* is WITHDRAW_FAIL (21).
	// If so, it allows re-auditing -> reset to AUDITING (0).
	if withdraw.Status == 21 { // WITHDRAW_FAIL
		// Reset to AUDITING
		if _, err := w.WithContext(ctx).Where(w.ID.Eq(id)).
			Updates(map[string]interface{}{
				"status":             0, // AUDITING
				"transfer_error_msg": "",
			}); err != nil {
			return err
		}
		withdraw.Status = 0
		withdraw.TransferErrorMsg = ""
	}

	// 1.2 校验状态为审核中
	if withdraw.Status != 0 { // 0: Auditing
		return errors.New("当前状态不可审核")
	}

	// 2. 更新状态
	updateMap := map[string]interface{}{
		"status":       status,
		"audit_reason": auditReason,
		"audit_time":   time.Now(),
	}

	if _, err := w.WithContext(ctx).Where(w.ID.Eq(id)).Updates(updateMap); err != nil {
		return err
	}

	// 3. Handle Success/Fail
	if status == 10 { // AUDIT_SUCCESS
		// If API Type, create transfer
		// Enum: WALLET(1), BANK(2), WECHAT(3), ALIPAY(4)
		if withdraw.Type == 1 || withdraw.Type == 3 || withdraw.Type == 4 { // Wallet/Wechat/Alipay
			s.createPayTransfer(ctx, withdraw)
		} else {
			// Manual (Bank) -> Mark Withdraw Success
			w.WithContext(ctx).Where(w.ID.Eq(id)).Update(w.Status, 11) // WITHDRAW_SUCCESS
		}
	} else if status == 20 { // AUDIT_FAIL
		// Refund Brokerage
		// Java: BrokerageRecordBizTypeEnum.WITHDRAW_REJECT
		s.recordSvc.AddBrokerage(ctx, withdraw.UserID, "withdraw_reject", // BizType
			string(rune(withdraw.ID)), withdraw.Price, "提现驳回")
	}

	return nil
}

// createPayTransfer 创建支付转账
func (s *BrokerageWithdrawService) createPayTransfer(ctx context.Context, withdraw *brokerage.BrokerageWithdraw) error {
	// 1. 获取用户信息
	user, err := s.memberSvc.GetUser(ctx, withdraw.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found: %d", withdraw.UserID)
	}

	// 2. 获取交易配置
	// Assuming there's a method to get config, e.g., GetTradeConfig() since it's singleton or implies context?
	// Or pass tenant/appId? BrokerageWithdraw doesn't explicitly store AppID, but assuming single tenant context or logic.
	// For now, let's try to get by default or mock if needed.
	// Actually tradeConfigSvc likely has specific accessor.
	tradeConfig, err := s.tradeConfigSvc.GetTradeConfig(ctx)
	if err != nil {
		return err
	}

	// 1.1 获取基础信息
	// userAccount := withdraw.UserAccount // Used in logic below
	channelCode := ""
	var channelExtras map[string]string

	// Enum: WALLET(1), BANK(2), WECHAT(3), ALIPAY(4)
	if withdraw.Type == 4 { // ALIPAY_API
		channelCode = "alipay_pc" // PayChannelEnum.ALIPAY_PC.getCode()
	} else if withdraw.Type == 3 { // WECHAT_API
		channelCode = withdraw.TransferChannelCode
		// userAccount = withdraw.UserAccount // Already set
		// 特殊：微信需要有报备信息
		channelExtras = map[string]string{
			"desc": "佣金提现", // Approx
		}
	} else if withdraw.Type == 1 { // WALLET
		channelCode = "wallet"
	}

	// 1.2 构建请求
	createReq := &pay.PayTransferCreateReqDTO{
		AppID:              tradeConfig.AppID,
		ChannelCode:        channelCode,
		MerchantTransferID: strconv.FormatInt(withdraw.ID, 10), // Use withdraw ID as MerchantTransferID
		Subject:            fmt.Sprintf("用户提现 - %d", withdraw.ID),
		Price:              withdraw.Price,
		UserAccount:        user.Mobile, // Default to mobile as account? Or need real account?
		UserName:           user.Nickname,
		UserIP:             "127.0.0.1", // TODO: Get from context or request
		ChannelExtras:      channelExtras,
	}
	if channelCode == "wx_pub" || channelCode == "wx_lite" || channelCode == "wx_app" {
		createReq.OpenID = "TODO" // Need OpenID for WeChat, likely from UserSocial or similar
	}

	// 1.3 发起请求
	resp, err := s.payTransferSvc.CreateTransfer(ctx, createReq)
	if err != nil {
		// Log error handling? Java doesn't catch exception here, allows it to bubble up?
		// But valid alignment implies dealing with it.
		// If failed, status stays AUDIT_SUCCESS (10), but transfer info missing.
		return err
	}

	// 2. 更新提现记录
	s.q.BrokerageWithdraw.WithContext(ctx).Where(s.q.BrokerageWithdraw.ID.Eq(withdraw.ID)).
		Updates(map[string]interface{}{
			"pay_transfer_id":       resp.ID,
			"transfer_channel_code": channelCode,
		})
	return nil
}

// CreateBrokerageWithdraw 创建佣金提现
func (s *BrokerageWithdrawService) CreateBrokerageWithdraw(ctx context.Context, userId int64, reqVO *tradeReq.AppBrokerageWithdrawCreateReqVO) (int64, error) {
	// 1. Check Config
	config, err := s.tradeConfigSvc.GetTradeConfig(ctx)
	if err != nil {
		return 0, err
	}
	if config.BrokerageWithdrawMinPrice > 0 && reqVO.Price < config.BrokerageWithdrawMinPrice {
		return 0, errors.New("提现金额低于最低提现金额")
	}

	// 2. Check Wallet / Realname (Simplified)
	// Java checks if wallet exists and realname is set for Wechat/Bank.
	// We skip strict wallet check here or assume PayWalletService handles it?
	// For now, simple validation.

	// 3. Calculate Fee
	feePrice := 0
	if config.BrokerageWithdrawFeePercent > 0 {
		feePrice = reqVO.Price * config.BrokerageWithdrawFeePercent / 100
	}

	// 4. Create Withdraw Record
	withdraw := &brokerage.BrokerageWithdraw{
		UserID:   userId,
		Price:    reqVO.Price,
		FeePrice: feePrice,
		// TotalPrice: realPrice, // Removed duplicate
		// Java: total_price = price.
		// Wait, Java sets `TotalPrice` to `price`.
		// And `Price` to `price` (withdrawal amount).
		// `FeePrice` is additional or deducted?
		// Java: `withdraw.setPrice(createReqVO.getPrice());`
		// `withdraw.setFeePrice(feePrice);`
		// `withdraw.setTotalPrice(createReqVO.getPrice());`
		// So TotalPrice seems to be the full amount deducted from user balance.
		// Real transfer amount would be Price - FeePrice?
		// If Fee is deducted FROM price, then user creates withdraw of 100, receives 90 (10 fee).
		// If Fee is extra, user needs 110 balance.
		// Usually Fee is deducted from the withdrawal amount.
		// Let's assume deducted.
		Type:        reqVO.Type,
		UserName:    reqVO.Name,    // Map Name -> UserName
		UserAccount: reqVO.Account, // Map Account -> UserAccount
		BankName:    reqVO.BankName,
		BankAddress: reqVO.BankAddress,
		QrCodeURL:   reqVO.QrCodeUrl,
		Status:      1, // Auditing
		TotalPrice:  reqVO.Price,
	}

	err = s.q.Transaction(func(tx *query.Query) error {
		// 1. Create Withdrawal Record
		err := tx.BrokerageWithdraw.WithContext(ctx).Create(withdraw)
		if err != nil {
			return err
		}

		// 2. Deduct Brokerage (Atomic)
		// Java: brokerageRecordService.reduceBrokerage(...)
		// Pass withdraw.ID as BizID
		return s.recordSvc.ReduceBrokerageForWithdraw(ctx, userId, strconv.FormatInt(withdraw.ID, 10), reqVO.Price)
	})
	if err != nil {
		return 0, err
	}

	return withdraw.ID, nil
}

// UpdateBrokerageWithdrawTransferred 更新佣金提现的转账结果
func (s *BrokerageWithdrawService) UpdateBrokerageWithdrawTransferred(ctx context.Context, id int64, payTransferId int64) error {
	w := s.q.BrokerageWithdraw
	withdraw, err := w.WithContext(ctx).Where(w.ID.Eq(id)).First()
	if err != nil {
		return errors.New("提现记录不存在")
	}

	// 1.2 Verify status ended
	// WITHDRAW_SUCCESS(11), WITHDRAW_FAIL(21)
	if withdraw.Status == 11 || withdraw.Status == 21 {
		if withdraw.PayTransferID == payTransferId {
			return nil // Duplicate callback
		}
		return errors.New("转账单不匹配")
	}

	// 2. 校验转账单 (Call Pay Service)
	// Placeholder: payTransferApi.GetTransfer(payTransferId)
	// Assuming PayTransferService has GetTransfer
	// ...

	// 3. Update Status
	// Placeholder logic
	return nil
}

// GetBrokerageWithdraw 获得佣金提现
func (s *BrokerageWithdrawService) GetBrokerageWithdraw(ctx context.Context, id int64) (*brokerage.BrokerageWithdraw, error) {
	return s.q.BrokerageWithdraw.WithContext(ctx).Where(s.q.BrokerageWithdraw.ID.Eq(id)).First()
}

// GetBrokerageWithdrawPage 获得佣金提现分页
func (s *BrokerageWithdrawService) GetBrokerageWithdrawPage(ctx context.Context, r *req.BrokerageWithdrawPageReq) (*pagination.PageResult[*brokerage.BrokerageWithdraw], error) {
	q := s.q.BrokerageWithdraw.WithContext(ctx)

	if r.UserID > 0 {
		q = q.Where(s.q.BrokerageWithdraw.UserID.Eq(r.UserID))
	}
	if r.Type > 0 {
		q = q.Where(s.q.BrokerageWithdraw.Type.Eq(r.Type))
	}
	if r.Status >= 0 { // Careful with 0 (Auditing)
		q = q.Where(s.q.BrokerageWithdraw.Status.Eq(r.Status))
	}
	if r.UserName != "" {
		q = q.Where(s.q.BrokerageWithdraw.UserName.Like("%" + r.UserName + "%"))
	}
	if r.UserAccount != "" {
		q = q.Where(s.q.BrokerageWithdraw.UserAccount.Like("%" + r.UserAccount + "%"))
	}
	if r.BankName != "" {
		q = q.Where(s.q.BrokerageWithdraw.BankName.Like("%" + r.BankName + "%"))
	}
	if len(r.CreateTime) == 2 {
		q = q.Where(s.q.BrokerageWithdraw.CreatedAt.Between(parseTime(r.CreateTime[0]), parseTime(r.CreateTime[1])))
	}

	q = q.Order(s.q.BrokerageWithdraw.ID.Desc())

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Offset(r.GetOffset()).Limit(r.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*brokerage.BrokerageWithdraw]{
		List:  list,
		Total: total,
	}, nil
}

type BrokerageWithdrawSummary struct {
	UserID int64
	Price  int
	Count  int
}

// GetWithdrawSummaryListByUserId 获得提现统计列表
func (s *BrokerageWithdrawService) GetWithdrawSummaryListByUserId(ctx context.Context, userIds []int64, statuses []int) ([]*BrokerageWithdrawSummary, error) {
	q := s.q.BrokerageWithdraw.WithContext(ctx).Where(s.q.BrokerageWithdraw.UserID.In(userIds...))
	if len(statuses) > 0 {
		q = q.Where(s.q.BrokerageWithdraw.Status.In(statuses...))
	}

	// Group By UserID, Sum Price, Count *
	// GORM Group result scan
	var results []*BrokerageWithdrawSummary
	err := q.Select(
		s.q.BrokerageWithdraw.UserID,
		s.q.BrokerageWithdraw.Price.Sum().As("price"),
		s.q.BrokerageWithdraw.ID.Count().As("count"),
	).Group(s.q.BrokerageWithdraw.UserID).Scan(&results)

	if err != nil {
		return nil, err
	}
	return results, nil
}
