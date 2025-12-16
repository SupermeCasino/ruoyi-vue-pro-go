package pay

import (
	"context"
	"fmt"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay/client"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PayTransferService struct {
	db            *gorm.DB
	logger        *zap.Logger
	appSvc        *PayAppService
	channelSvc    *PayChannelService
	notifySvc     *PayNotifyService
	clientFactory *client.PayClientFactory
	q             *query.Query
}

func NewPayTransferService(db *gorm.DB, logger *zap.Logger, appSvc *PayAppService, channelSvc *PayChannelService, notifySvc *PayNotifyService, clientFactory *client.PayClientFactory, q *query.Query) *PayTransferService {
	return &PayTransferService{
		db:            db,
		logger:        logger,
		appSvc:        appSvc,
		channelSvc:    channelSvc,
		notifySvc:     notifySvc,
		clientFactory: clientFactory,
		q:             q,
	}
}

// CreateTransfer 创建转账
func (s *PayTransferService) CreateTransfer(ctx context.Context, req *PayTransferCreateReqDTO) (*PayTransferCreateRespDTO, error) {
	// 1. 校验应用
	app, err := s.appSvc.ValidPayApp(ctx, req.AppID)
	if err != nil {
		return nil, err
	}

	// 2. 校验渠道
	channel, err := s.channelSvc.GetChannelByAppIdAndCode(ctx, req.AppID, req.ChannelCode)
	if err != nil {
		return nil, err
	}
	if channel == nil {
		return nil, core.NewBizError(1006002002, "支付渠道不存在") // PAY_CHANNEL_NOT_FOUND
	}
	if channel.Status != 0 {
		return nil, core.NewBizError(1006002001, "支付渠道处于关闭状态") // PAY_CHANNEL_IS_DISABLE
	}

	// 3. 获得支付客户端
	payClient := s.clientFactory.GetPayClient(channel.ID)
	if payClient == nil {
		return nil, core.NewBizError(1006000003, "支付渠道客户端不存在") // PAY_CHANNEL_CLIENT_NOT_FOUND
	}

	// 3. 创建转账单
	// subject and no generation logic
	no := s.generateNo()
	transfer := &pay.PayTransfer{
		AppID:              req.AppID,
		ChannelID:          channel.ID,
		ChannelCode:        channel.Code,
		MerchantTransferID: req.MerchantTransferID,
		Subject:            req.Subject,
		Price:              req.Price,
		UserAccount:        req.UserAccount,
		UserName:           req.UserName,
		Status:             0, // WAITING
		No:                 no,
		NotifyURL:          app.TransferNotifyURL,
		UserIP:             req.UserIP,
		ChannelExtras:      req.ChannelExtras,
	}
	if err := s.db.WithContext(ctx).Create(transfer).Error; err != nil {
		return nil, err
	}

	// 4. 调用支付渠道
	unifiedReq := &client.UnifiedTransferReq{
		OutTradeNo:    transfer.No,
		Subject:       transfer.Subject,
		Price:         transfer.Price,
		ChannelExtras: req.ChannelExtras,
		UserIP:        req.UserIP,
		ChannelUserID: req.OpenID, // Assumed passed in extras or somehow?
		UserName:      req.UserName,
		UserAccount:   req.UserAccount,
	}
	// Note: OpenID might be needed for WeChat. Often passed in ChannelExtras or a separate field in CreateReq?
	// In Java reqDTO has `channelExtras`.

	resp, err := payClient.UnifiedTransfer(ctx, unifiedReq)
	if err != nil {
		// Log error and update status to CLOSED or FAIL?
		s.logger.Error("UnifiedTransfer failed", zap.Error(err))
		// Optional: Update transfer status to FAIL
		return nil, err
	}

	// 5. 更新转账单状态
	// Map client status to system status
	// Assuming resp.Status is aligned or needs mapping.
	// 0: Waiting, 10: Success, 20: Closed/Fail?
	transfer.Status = resp.Status
	transfer.ChannelTransferNo = resp.ChannelTransferNo
	if !resp.SuccessTime.IsZero() {
		transfer.SuccessTime = &resp.SuccessTime
	}
	// channel extras, error code etc.

	s.db.WithContext(ctx).Save(transfer)

	return &PayTransferCreateRespDTO{
		ID:     transfer.ID,
		Status: transfer.Status,
	}, nil
}

func (s *PayTransferService) GetTransfer(ctx context.Context, id int64) (*PayTransferRespDTO, error) {
	var transfer pay.PayTransfer
	if err := s.db.WithContext(ctx).First(&transfer, id).Error; err != nil {
		return nil, err
	}
	return s.convertRespDTO(&transfer), nil
}

func (s *PayTransferService) convertRespDTO(t *pay.PayTransfer) *PayTransferRespDTO {
	return &PayTransferRespDTO{
		ID:                 t.ID,
		Status:             t.Status,
		Price:              t.Price,
		MerchantTransferId: t.MerchantTransferID,
		ChannelCode:        t.ChannelCode,
		SuccessTime:        t.SuccessTime,
		ChannelErrorMsg:    t.ChannelErrorMsg,
		ChannelExtras:      t.ChannelExtras,
	}
}

func (s *PayTransferService) generateNo() string {
	return fmt.Sprintf("T%d", time.Now().UnixNano())
}

// NotifyTransfer 通知并更新转账结果
// 对齐 Java: PayTransferService.notifyTransfer(Long channelId, PayTransferRespDTO notify)
func (s *PayTransferService) NotifyTransfer(ctx context.Context, channelID int64, notify *client.TransferResp) error {
	// 校验渠道是否有效
	channel, err := s.channelSvc.ValidPayChannel(ctx, channelID)
	if err != nil {
		return err
	}

	// 转账成功的回调
	if notify.Status == PayTransferStatusSuccess {
		return s.notifyTransferSuccess(ctx, channel, notify)
	}

	// 转账关闭的回调
	if notify.Status == PayTransferStatusClosed {
		return s.notifyTransferClosed(ctx, channel, notify)
	}

	// 转账处理中的回调
	if notify.Status == PayTransferStatusProcessing {
		return s.notifyTransferProcessing(ctx, channel, notify)
	}

	// WAITING 状态无需处理
	return nil
}

// notifyTransferProcessing 处理转账进行中
func (s *PayTransferService) notifyTransferProcessing(ctx context.Context, channel *pay.PayChannel, notify *client.TransferResp) error {
	// 1. 查询转账单
	var transfer pay.PayTransfer
	if err := s.db.WithContext(ctx).
		Where("app_id = ? AND no = ?", channel.AppID, notify.OutTradeNo).
		First(&transfer).Error; err != nil {
		return fmt.Errorf("转账单不存在")
	}

	// 如果已经是转账中，直接返回
	if transfer.Status == PayTransferStatusProcessing {
		return nil
	}

	// 校验状态，必须是等待状态
	if transfer.Status != PayTransferStatusWaiting {
		return fmt.Errorf("转账单状态不是待转账")
	}

	// 2. 更新状态 (使用乐观锁)
	result := s.db.WithContext(ctx).
		Model(&pay.PayTransfer{}).
		Where("id = ? AND status = ?", transfer.ID, PayTransferStatusWaiting).
		Updates(map[string]interface{}{
			"status": PayTransferStatusProcessing,
		})

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("转账单状态不是待转账")
	}

	return nil
}

// notifyTransferSuccess 处理转账成功
func (s *PayTransferService) notifyTransferSuccess(ctx context.Context, channel *pay.PayChannel, notify *client.TransferResp) error {
	// 1. 查询转账单
	var transfer pay.PayTransfer
	if err := s.db.WithContext(ctx).
		Where("app_id = ? AND no = ?", channel.AppID, notify.OutTradeNo).
		First(&transfer).Error; err != nil {
		return fmt.Errorf("转账单不存在")
	}

	// 如果已经是成功，直接返回
	if transfer.Status == PayTransferStatusSuccess {
		return nil
	}

	// 校验状态，必须是等待或进行中状态
	if transfer.Status != PayTransferStatusWaiting && transfer.Status != PayTransferStatusProcessing {
		return fmt.Errorf("转账单状态不是待转账或转账中")
	}

	// 2. 更新状态 (使用乐观锁)
	result := s.db.WithContext(ctx).
		Model(&pay.PayTransfer{}).
		Where("id = ? AND status IN (?, ?)", transfer.ID, PayTransferStatusWaiting, PayTransferStatusProcessing).
		Updates(map[string]interface{}{
			"status":              PayTransferStatusSuccess,
			"success_time":        notify.SuccessTime,
			"channel_transfer_no": notify.ChannelTransferNo,
		})

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("转账单状态不是待转账或转账中")
	}

	// 3. 插入转账通知记录
	s.notifySvc.CreatePayNotifyTask(ctx, PayNotifyTypeTransfer, transfer.ID)

	return nil
}

// notifyTransferClosed 处理转账关闭
func (s *PayTransferService) notifyTransferClosed(ctx context.Context, channel *pay.PayChannel, notify *client.TransferResp) error {
	// 1. 查询转账单
	var transfer pay.PayTransfer
	if err := s.db.WithContext(ctx).
		Where("app_id = ? AND no = ?", channel.AppID, notify.OutTradeNo).
		First(&transfer).Error; err != nil {
		return fmt.Errorf("转账单不存在")
	}

	// 如果已经是关闭，直接返回
	if transfer.Status == PayTransferStatusClosed {
		return nil
	}

	// 校验状态，必须是等待或进行中状态
	if transfer.Status != PayTransferStatusWaiting && transfer.Status != PayTransferStatusProcessing {
		return fmt.Errorf("转账单状态不是待转账或转账中")
	}

	// 2. 更新状态 (使用乐观锁)
	result := s.db.WithContext(ctx).
		Model(&pay.PayTransfer{}).
		Where("id = ? AND status IN (?, ?)", transfer.ID, PayTransferStatusWaiting, PayTransferStatusProcessing).
		Updates(map[string]interface{}{
			"status":              PayTransferStatusClosed,
			"channel_transfer_no": notify.ChannelTransferNo,
			"channel_error_code":  notify.ChannelErrorCode,
			"channel_error_msg":   notify.ChannelErrorMsg,
		})

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("转账单状态不是待转账或转账中")
	}

	// 3. 插入转账通知记录
	s.notifySvc.CreatePayNotifyTask(ctx, PayNotifyTypeTransfer, transfer.ID)

	return nil
}
