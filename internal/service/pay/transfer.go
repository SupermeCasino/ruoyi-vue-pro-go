package pay

import (
	"backend-go/internal/model/pay"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"backend-go/internal/service/pay/client"
	"context"
	"fmt"
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
