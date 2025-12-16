package trade

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type TradeAfterSaleService struct {
	q        *query.Query
	orderSvc *TradeOrderUpdateService
}

func NewTradeAfterSaleService(q *query.Query, orderSvc *TradeOrderUpdateService) *TradeAfterSaleService {
	return &TradeAfterSaleService{
		q:        q,
		orderSvc: orderSvc,
	}
}

// CreateAfterSale 创建售后
func (s *TradeAfterSaleService) CreateAfterSale(ctx context.Context, userId int64, r *req.AppAfterSaleCreateReq) (int64, error) {
	// 1. Get Order Item
	item, err := s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.ID.Eq(r.OrderItemID), s.q.TradeOrderItem.UserID.Eq(userId)).First()
	if err != nil {
		return 0, err
	}

	// Validate Status
	if item.AfterSaleStatus != 0 { // 0: None
		return 0, fmt.Errorf("该订单项已申请售后")
	}

	// Fetch Order for OrderNo
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(item.OrderID)).First()
	if err != nil {
		return 0, err
	}

	// 3. Create AfterSale & Update Item Status
	pics, _ := json.Marshal(r.ApplyPicURLs)
	props, _ := json.Marshal(item.Properties)

	afterSale := &trade.AfterSale{
		No:               fmt.Sprintf("%d", time.Now().UnixNano()), // TODO: Better No Gen
		Status:           10,                                       // Applying
		Way:              r.Way,
		Type:             r.Type,
		UserID:           userId,
		ApplyReason:      r.ApplyReason,
		ApplyDescription: r.ApplyDescription,
		ApplyPicURLs:     string(pics),
		OrderID:          item.OrderID,
		OrderNo:          order.No,
		OrderItemID:      item.ID,
		SpuID:            item.SpuID,
		SpuName:          item.SpuName,
		SkuID:            item.SkuID,
		Properties:       string(props),
		PicURL:           item.PicURL,
		Count:            r.Count,
		RefundPrice:      r.RefundPrice,
	}

	err = s.q.Transaction(func(tx *query.Query) error {
		if err := tx.AfterSale.WithContext(ctx).Create(afterSale); err != nil {
			return err
		}
		// Update OrderItem
		if _, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.ID.Eq(item.ID)).Updates(map[string]interface{}{
			"after_sale_id":     afterSale.ID,
			"after_sale_status": 10,
		}); err != nil {
			return err
		}
		return nil
	})
	return afterSale.ID, err
}

// GetAfterSale 获得售后详情
func (s *TradeAfterSaleService) GetAfterSale(ctx context.Context, userId int64, id int64) (*resp.AppAfterSaleResp, error) {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(id), s.q.AfterSale.UserID.Eq(userId)).First()
	if err != nil {
		return nil, err
	}

	var pics []string
	_ = json.Unmarshal([]byte(as.ApplyPicURLs), &pics)

	return &resp.AppAfterSaleResp{
		ID:               as.ID,
		No:               as.No,
		Status:           as.Status,
		Way:              as.Way,
		Type:             as.Type,
		ApplyReason:      as.ApplyReason,
		ApplyDescription: as.ApplyDescription,
		ApplyPicURLs:     pics,
		OrderNo:          as.OrderNo,
		SpuName:          as.SpuName,
		PicURL:           as.PicURL,
		Count:            as.Count,
		RefundPrice:      as.RefundPrice,
		AuditTime:        as.AuditTime,
		AuditReason:      as.AuditReason,
		CreateTime:       as.CreatedAt,
	}, nil
}

// GetAfterSalePage 获得售后分页 (Admin)
func (s *TradeAfterSaleService) GetAfterSalePage(ctx context.Context, r *req.TradeAfterSalePageReq) (*pagination.PageResult[*trade.AfterSale], error) {
	q := s.q.AfterSale.WithContext(ctx)
	if r.No != "" {
		q = q.Where(s.q.AfterSale.No.Like("%" + r.No + "%"))
	}
	if r.UserID != nil {
		q = q.Where(s.q.AfterSale.UserID.Eq(*r.UserID))
	}
	if r.Status != nil {
		q = q.Where(s.q.AfterSale.Status.Eq(*r.Status))
	}

	list, total, err := q.Order(s.q.AfterSale.ID.Desc()).FindByPage(r.GetOffset(), r.PageSize)
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*trade.AfterSale]{
		List:  list,
		Total: total,
	}, nil
}

// CancelAfterSale 取消售后
func (s *TradeAfterSaleService) CancelAfterSale(ctx context.Context, userId int64, id int64) error {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(id), s.q.AfterSale.UserID.Eq(userId)).First()
	if err != nil {
		return err
	}
	if as.Status != 10 {
		return fmt.Errorf("状态不允许取消")
	}

	return s.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.AfterSale.WithContext(ctx).Where(tx.AfterSale.ID.Eq(id)).Update(tx.AfterSale.Status, 60); err != nil { // 60: Cancelled
			return err
		}
		// Update OrderItem Status to 0 (None)? Or keep history?
		// Usually set back to 0 so user can apply again? Or maybe 60.
		// If 0, user can apply again.
		if _, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.ID.Eq(as.OrderItemID)).Updates(map[string]interface{}{
			"after_sale_status": 0, // Reset
		}); err != nil {
			return err
		}
		return nil
	})
}

// AgreeAfterSale 同意售后
func (s *TradeAfterSaleService) AgreeAfterSale(ctx context.Context, id int64) error {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(id)).First()
	if err != nil {
		return err
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// Update AfterSale
		if _, err := tx.AfterSale.WithContext(ctx).Where(tx.AfterSale.ID.Eq(id)).Updates(trade.AfterSale{
			Status:    20, // Approved
			AuditTime: time.Now(),
		}); err != nil {
			return err
		}
		// Update OrderItem Status
		if _, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.ID.Eq(as.OrderItemID)).Update(tx.TradeOrderItem.AfterSaleStatus, 20); err != nil {
			return err
		}
		return nil
	})
}

// DisagreeAfterSale 拒绝售后 (审核不通过)
func (s *TradeAfterSaleService) DisagreeAfterSale(ctx context.Context, req *req.TradeAfterSaleDisagreeReq) error {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(req.ID)).First()
	if err != nil {
		return err
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// Update AfterSale
		if _, err := tx.AfterSale.WithContext(ctx).Where(tx.AfterSale.ID.Eq(req.ID)).Updates(trade.AfterSale{
			Status:      50, // Disagree (Audit Failed)
			AuditReason: req.AuditReason,
			AuditTime:   time.Now(),
		}); err != nil {
			return err
		}
		// Update OrderItem Status (Reset to 0? Or 50?)
		// If refused, usually user can apply again or see Refused status. Default: 50.
		if _, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.ID.Eq(as.OrderItemID)).Update(tx.TradeOrderItem.AfterSaleStatus, 50); err != nil {
			return err
		}
		return nil
	})
}

// RefundAfterSale 退款
func (s *TradeAfterSaleService) RefundAfterSale(ctx context.Context, id int64) error {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(id)).First()
	if err != nil {
		return err
	}

	// TODO: Call Pay Service to Refund

	return s.q.Transaction(func(tx *query.Query) error {
		// Update AfterSale
		if _, err := tx.AfterSale.WithContext(ctx).Where(tx.AfterSale.ID.Eq(id)).Updates(trade.AfterSale{
			Status:     30, // Completed/Refunded
			RefundTime: time.Now(),
		}); err != nil {
			return err
		}
		// Update OrderItem Status
		if _, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.ID.Eq(as.OrderItemID)).Update(tx.TradeOrderItem.AfterSaleStatus, 30); err != nil {
			return err
		}
		return nil
	})
}

// GetAfterSaleDetail 获得售后订单详情 (Admin)
func (s *TradeAfterSaleService) GetAfterSaleDetail(ctx context.Context, id int64) (*resp.TradeAfterSaleDetailResp, error) {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	// 获取订单信息
	order, _ := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(as.OrderID)).First()
	// 获取订单项信息
	orderItem, _ := s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.ID.Eq(as.OrderItemID)).First()

	var pics []string
	_ = json.Unmarshal([]byte(as.ApplyPicURLs), &pics)

	result := &resp.TradeAfterSaleDetailResp{
		ID:               as.ID,
		No:               as.No,
		Status:           as.Status,
		Way:              as.Way,
		Type:             as.Type,
		UserID:           as.UserID,
		ApplyReason:      as.ApplyReason,
		ApplyDescription: as.ApplyDescription,
		ApplyPicURLs:     pics,
		OrderID:          as.OrderID,
		OrderNo:          as.OrderNo,
		OrderItemID:      as.OrderItemID,
		SpuID:            as.SpuID,
		SpuName:          as.SpuName,
		SkuID:            as.SkuID,
		PicURL:           as.PicURL,
		Count:            as.Count,
		RefundPrice:      as.RefundPrice,
		AuditTime:        as.AuditTime,
		AuditReason:      as.AuditReason,
		RefundTime:       as.RefundTime,
		CreateTime:       as.CreatedAt,
	}

	if order != nil {
		result.OrderPayPrice = order.PayPrice
	}
	if orderItem != nil {
		result.OrderItemPayPrice = orderItem.PayPrice
	}

	return result, nil
}

// ReceiveAfterSale 确认收货 (Admin)
func (s *TradeAfterSaleService) ReceiveAfterSale(ctx context.Context, id int64) error {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(id)).First()
	if err != nil {
		return err
	}

	// 校验状态：只有待收货状态才能确认收货
	if as.Status != 20 { // 20: 待收货/已同意
		return fmt.Errorf("售后状态不允许确认收货")
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// Update AfterSale Status to 待退款
		if _, err := tx.AfterSale.WithContext(ctx).Where(tx.AfterSale.ID.Eq(id)).Updates(trade.AfterSale{
			Status:      25, // 待退款
			ReceiveTime: time.Now(),
		}); err != nil {
			return err
		}
		// Update OrderItem Status
		if _, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.ID.Eq(as.OrderItemID)).Update(tx.TradeOrderItem.AfterSaleStatus, 25); err != nil {
			return err
		}
		return nil
	})
}

// DeliveryAfterSale 用户退回货物 (App)
func (s *TradeAfterSaleService) DeliveryAfterSale(ctx context.Context, userId int64, req *req.AppAfterSaleDeliveryReq) error {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(req.ID), s.q.AfterSale.UserID.Eq(userId)).First()
	if err != nil {
		return fmt.Errorf("售后单不存在")
	}

	// 校验状态：只有已同意状态才能填写物流信息
	if as.Status != 20 { // 20: 已同意/待退货
		return fmt.Errorf("售后状态不允许填写物流信息")
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// Update AfterSale with logistics info
		if _, err := tx.AfterSale.WithContext(ctx).Where(tx.AfterSale.ID.Eq(req.ID)).Updates(trade.AfterSale{
			LogisticsID:  req.LogisticsId,
			LogisticsNo:  req.LogisticsNo,
			DeliveryTime: time.Now(),
		}); err != nil {
			return err
		}
		return nil
	})
}

// UpdateAfterSaleRefunded 更新售后单为已退款 (Callback from Pay)
func (s *TradeAfterSaleService) UpdateAfterSaleRefunded(ctx context.Context, afterSaleId int64, payRefundId int64) error {
	as, err := s.q.AfterSale.WithContext(ctx).Where(s.q.AfterSale.ID.Eq(afterSaleId)).First()
	if err != nil {
		return err
	}
	if as.Status == 30 { // Already refunded
		return nil
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// Update AfterSale
		if _, err := tx.AfterSale.WithContext(ctx).Where(tx.AfterSale.ID.Eq(afterSaleId)).Updates(trade.AfterSale{
			Status:      30, // Completed/Refunded
			RefundTime:  time.Now(),
			PayRefundID: payRefundId,
		}); err != nil {
			return err
		}
		// Update OrderItem Status
		if _, err := tx.TradeOrderItem.WithContext(ctx).Where(tx.TradeOrderItem.ID.Eq(as.OrderItemID)).Update(tx.TradeOrderItem.AfterSaleStatus, 30); err != nil {
			return err
		}
		return nil
	})
}

// UpdateRefunded 更新退款状态 (Unified Facade)
func (s *TradeAfterSaleService) UpdateRefunded(ctx context.Context, req *req.PayRefundNotifyReqDTO) error {
	if strings.HasPrefix(req.MerchantRefundId, "order-") {
		orderId, err := strconv.ParseInt(strings.TrimPrefix(req.MerchantRefundId, "order-"), 10, 64)
		if err != nil {
			return err
		}
		return s.orderSvc.UpdatePaidOrderRefunded(ctx, orderId, req.PayRefundId)
	} else {
		afterSaleId, err := strconv.ParseInt(req.MerchantRefundId, 10, 64)
		if err != nil {
			return err
		}
		return s.UpdateAfterSaleRefunded(ctx, afterSaleId, req.PayRefundId)
	}
}

// GetUserAfterSaleCount 获得用户的售后数量 (进行中)
func (s *TradeAfterSaleService) GetUserAfterSaleCount(ctx context.Context, userId int64) (int64, error) {
	// Status: 10 (Applying), 20 (Approved/Processing)
	return s.q.AfterSale.WithContext(ctx).
		Where(s.q.AfterSale.UserID.Eq(userId)).
		Where(s.q.AfterSale.Status.In(10, 20)).
		Count()
}
