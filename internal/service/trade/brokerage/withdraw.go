package brokerage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	tradeReq "github.com/wxlbd/ruoyi-mall-go/internal/api/req/app/trade"
	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
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

	// 1.1 校验存在
	withdraw, err := w.WithContext(ctx).Where(w.ID.Eq(id)).First()
	if err != nil {
		return errors.New("提现记录不存在")
	}

	// 1.2 特殊：【重新转账】如果是提现失败，允许重新审核
	if withdraw.Status == tradeModel.BrokerageWithdrawStatusWithdrawFail {
		// 重置为审核中状态
		if _, err := w.WithContext(ctx).Where(w.ID.Eq(id)).
			Updates(map[string]interface{}{
				"status":             tradeModel.BrokerageWithdrawStatusAuditing,
				"transfer_error_msg": "",
			}); err != nil {
			return err
		}
		withdraw.Status = tradeModel.BrokerageWithdrawStatusAuditing
		withdraw.TransferErrorMsg = ""
	}

	// 1.3 校验状态为审核中
	if withdraw.Status != tradeModel.BrokerageWithdrawStatusAuditing {
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

	// 3. 处理审批结果
	switch status {
	case tradeModel.BrokerageWithdrawStatusAuditSuccess:
		// 3.1 审批通过的后续处理
		s.auditBrokerageWithdrawSuccess(ctx, withdraw)
	case tradeModel.BrokerageWithdrawStatusAuditFail:
		// 3.2 审批不通过：退还佣金
		s.recordSvc.AddBrokerage(ctx, withdraw.UserID, "withdraw_reject",
			strconv.FormatInt(withdraw.ID, 10), withdraw.Price, "提现驳回")
	default:
		return fmt.Errorf("不支持的提现状态：%d", status)
	}

	return nil
}

// auditBrokerageWithdrawSuccess 审批通过的后续处理
func (s *BrokerageWithdrawService) auditBrokerageWithdrawSuccess(ctx context.Context, withdraw *brokerage.BrokerageWithdraw) {
	// 情况一：通过 API 转账（钱包/微信/支付宝）
	if s.isApiWithdrawType(withdraw.Type) {
		s.createPayTransfer(ctx, withdraw)
		return
	}

	// 情况二：非 API 转账（银行卡，手动打款）
	s.q.BrokerageWithdraw.WithContext(ctx).Where(s.q.BrokerageWithdraw.ID.Eq(withdraw.ID)).
		Update(s.q.BrokerageWithdraw.Status, tradeModel.BrokerageWithdrawStatusWithdrawSuccess)
}

// isApiWithdrawType 判断是否为 API 提现类型
func (s *BrokerageWithdrawService) isApiWithdrawType(withdrawType int) bool {
	return withdrawType == tradeModel.BrokerageWithdrawTypeWallet ||
		withdrawType == tradeModel.BrokerageWithdrawTypeWechat ||
		withdrawType == tradeModel.BrokerageWithdrawTypeAlipay
}

// createPayTransfer 创建支付转账
func (s *BrokerageWithdrawService) createPayTransfer(ctx context.Context, withdraw *brokerage.BrokerageWithdraw) error {
	// 1.1 获取基础信息
	userAccount := withdraw.UserAccount
	userName := withdraw.UserName
	channelCode := ""
	var channelExtras map[string]string

	switch withdraw.Type {
	case tradeModel.BrokerageWithdrawTypeAlipay:
		channelCode = "alipay_pc"
	case tradeModel.BrokerageWithdrawTypeWechat:
		channelCode = withdraw.TransferChannelCode
		userAccount = withdraw.UserAccount
		// 特殊：微信需要有报备信息
		channelExtras = map[string]string{
			"scene":    "佣金提现",
			"sceneId":  "1000",
			"userName": userName,
		}
	case tradeModel.BrokerageWithdrawTypeWallet:
		// 钱包转账
		channelCode = "wallet"
	}

	// 1.2 获取交易配置
	tradeConfig, err := s.tradeConfigSvc.GetTradeConfig(ctx)
	if err != nil {
		return err
	}

	// 1.3 构建请求
	createReq := &pay.PayTransferCreateReqDTO{
		AppID:              tradeConfig.AppID,
		ChannelCode:        channelCode,
		MerchantTransferID: strconv.FormatInt(withdraw.ID, 10),
		Subject:            "佣金提现",
		Price:              withdraw.Price,
		UserAccount:        userAccount,
		UserName:           userName,
		UserIP:             "127.0.0.1", // TODO: 暂用默认值，后续可通过 context 传递
		ChannelExtras:      channelExtras,
	}

	// 1.4 发起请求
	resp, err := s.payTransferSvc.CreateTransfer(ctx, createReq)
	if err != nil {
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
	// 1.1 校验提现金额
	config, err := s.tradeConfigSvc.GetTradeConfig(ctx)
	if err != nil {
		return 0, err
	}
	if config.BrokerageWithdrawMinPrice > 0 && reqVO.Price < config.BrokerageWithdrawMinPrice {
		return 0, errors.New("提现金额低于最低提现金额")
	}

	// 2.1 计算手续费
	feePrice := 0
	if config.BrokerageWithdrawFeePercent > 0 {
		feePrice = reqVO.Price * config.BrokerageWithdrawFeePercent / 100
	}

	// 2.2 创建佣金提现记录
	withdraw := &brokerage.BrokerageWithdraw{
		UserID:      userId,
		Price:       reqVO.Price,
		FeePrice:    feePrice,
		TotalPrice:  reqVO.Price, // Java: setTotalPrice(price)
		Type:        reqVO.Type,
		UserName:    reqVO.Name,
		UserAccount: reqVO.Account,
		BankName:    reqVO.BankName,
		BankAddress: reqVO.BankAddress,
		QrCodeURL:   reqVO.QrCodeUrl,
		Status:      tradeModel.BrokerageWithdrawStatusAuditing,
	}

	err = s.q.Transaction(func(tx *query.Query) error {
		// 创建提现记录
		err := tx.BrokerageWithdraw.WithContext(ctx).Create(withdraw)
		if err != nil {
			return err
		}

		// 3. 创建用户佣金记录（扣减佣金）
		// 注意：佣金是否充足，ReduceBrokerageForWithdraw 已经进行校验
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

	// 1.1 校验提现单是否存在
	withdraw, err := w.WithContext(ctx).Where(w.ID.Eq(id)).First()
	if err != nil {
		s.logger.Error("提现单不存在", zap.Int64("id", id), zap.Int64("payTransferId", payTransferId))
		return errors.New("提现记录不存在")
	}

	// 1.2 校验提现单已经结束（成功或失败）
	if withdraw.Status == tradeModel.BrokerageWithdrawStatusWithdrawSuccess ||
		withdraw.Status == tradeModel.BrokerageWithdrawStatusWithdrawFail {
		// 特殊：转账单编号相同，直接返回，说明重复回调
		if withdraw.PayTransferID == payTransferId {
			s.logger.Warn("提现单已结束，且转账单编号相同，直接返回",
				zap.Int64("id", id), zap.Int64("payTransferId", payTransferId))
			return nil
		}
		// 异常：转账单编号不同，说明转账单编号错误
		s.logger.Error("转账单不匹配", zap.Int64("id", id), zap.Int64("payTransferId", payTransferId))
		return errors.New("转账单不匹配")
	}

	// 2. 校验转账单的合法性
	payTransfer, err := s.payTransferSvc.GetTransfer(ctx, payTransferId)
	if err != nil || payTransfer == nil {
		s.logger.Error("转账单不存在", zap.Int64("id", id), zap.Int64("payTransferId", payTransferId))
		return errors.New("转账单不存在")
	}

	// 2.1 校验转账单已成功或关闭
	if payTransfer.Status != tradeModel.PayTransferStatusSuccess && payTransfer.Status != tradeModel.PayTransferStatusClosed {
		s.logger.Error("转账单未结束", zap.Int64("id", id), zap.Int64("payTransferId", payTransferId))
		return errors.New("转账单未结束")
	}

	// 2.2 校验转账金额一致
	if payTransfer.Price != withdraw.Price {
		s.logger.Error("转账金额不匹配", zap.Int64("id", id), zap.Int("withdrawPrice", withdraw.Price), zap.Int("transferPrice", payTransfer.Price))
		return errors.New("转账金额不匹配")
	}

	// 3. 更新提现单状态
	var newStatus int
	if payTransfer.Status == tradeModel.PayTransferStatusSuccess {
		newStatus = tradeModel.BrokerageWithdrawStatusWithdrawSuccess
	} else {
		newStatus = tradeModel.BrokerageWithdrawStatusWithdrawFail
	}

	_, err = w.WithContext(ctx).Where(w.ID.Eq(id)).Updates(map[string]interface{}{
		"status":             newStatus,
		"transfer_time":      payTransfer.SuccessTime,
		"transfer_error_msg": payTransfer.ChannelErrorMsg,
	})
	return err
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
