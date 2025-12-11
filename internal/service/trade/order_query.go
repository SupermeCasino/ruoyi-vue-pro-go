package trade

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/trade"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"backend-go/internal/service/trade/delivery/client"
	"context"
)

type TradeOrderQueryService struct {
	q                    *query.Query
	expressClientFactory client.ExpressClientFactory
	deliveryExpressSvc   *DeliveryExpressService
}

func NewTradeOrderQueryService(q *query.Query, expressClientFactory client.ExpressClientFactory, deliveryExpressSvc *DeliveryExpressService) *TradeOrderQueryService {
	return &TradeOrderQueryService{
		q:                    q,
		expressClientFactory: expressClientFactory,
		deliveryExpressSvc:   deliveryExpressSvc,
	}
}

// ... GetOrder, GetOrderPage methods ...
// (Skipping middle methods for brevity in edit, but replace_file_content needs exact match or range.
// I will just replace the struct/constructor first, then the method separately to avoid massive context issues)

// GetOrder 获得交易订单
func (s *TradeOrderQueryService) GetOrder(ctx context.Context, id int64) (*trade.TradeOrder, error) {
	return s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id)).First()
}

// GetOrderPage 获得交易订单分页
func (s *TradeOrderQueryService) GetOrderPage(ctx context.Context, uId int64, r *req.AppTradeOrderPageReq) (*core.PageResult[*trade.TradeOrder], error) {
	q := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.UserID.Eq(uId))

	if r.Status != nil {
		q = q.Where(s.q.TradeOrder.Status.Eq(*r.Status))
	}
	if r.CommentStatus != nil {
		q = q.Where(s.q.TradeOrder.CommentStatus.Is(*r.CommentStatus))
	}

	list, total, err := q.Order(s.q.TradeOrder.ID.Desc()).FindByPage(r.GetOffset(), r.PageSize)
	if err != nil {
		return nil, err
	}

	result := &core.PageResult[*trade.TradeOrder]{
		List:  list,
		Total: total,
	}
	return result, nil
}

// GetOrderItemListByOrderId 获得交易订单项列表
func (s *TradeOrderQueryService) GetOrderItemListByOrderId(ctx context.Context, orderId int64) ([]*trade.TradeOrderItem, error) {
	return s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.OrderID.Eq(orderId)).Find()
}

// GetOrderItemListByOrderIds 获得交易订单项列表
func (s *TradeOrderQueryService) GetOrderItemListByOrderIds(ctx context.Context, orderIds []int64) ([]*trade.TradeOrderItem, error) {
	return s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.OrderID.In(orderIds...)).Find()
}

// GetOrderCount 获得交易订单数量
func (s *TradeOrderQueryService) GetOrderCount(ctx context.Context, userId int64, status *int, commentStatus *bool) (int64, error) {
	q := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.UserID.Eq(userId))
	if status != nil {
		q = q.Where(s.q.TradeOrder.Status.Eq(*status))
	}
	if commentStatus != nil {
		q = q.Where(s.q.TradeOrder.CommentStatus.Is(*commentStatus))
	}
	return q.Count()
}

// GetOrderItem 获得交易订单项
func (s *TradeOrderQueryService) GetOrderItem(ctx context.Context, userId int64, id int64) (*trade.TradeOrderItem, error) {
	return s.q.TradeOrderItem.WithContext(ctx).Where(s.q.TradeOrderItem.ID.Eq(id), s.q.TradeOrderItem.UserID.Eq(userId)).First()
}

// GetOrderPageForAdmin 获得交易订单分页 (Admin)
func (s *TradeOrderQueryService) GetOrderPageForAdmin(ctx context.Context, r *req.TradeOrderPageReq) (*core.PageResult[*trade.TradeOrder], error) {
	q := s.q.TradeOrder.WithContext(ctx)

	if r.No != "" {
		q = q.Where(s.q.TradeOrder.No.Like("%" + r.No + "%"))
	}
	if r.UserID != nil {
		q = q.Where(s.q.TradeOrder.UserID.Eq(*r.UserID))
	}
	if r.Status != nil {
		q = q.Where(s.q.TradeOrder.Status.Eq(*r.Status))
	}
	// Add more filters as needed (e.g. create_time, type, etc.)

	list, total, err := q.Order(s.q.TradeOrder.ID.Desc()).FindByPage(r.GetOffset(), r.PageSize)
	if err != nil {
		return nil, err
	}

	result := &core.PageResult[*trade.TradeOrder]{
		List:  list,
		Total: total,
	}
	return result, nil
}

// GetExpressTrackList 获得物流轨迹 (App - requires UserId)
func (s *TradeOrderQueryService) GetExpressTrackList(ctx context.Context, id int64, userId int64) ([]*resp.ExpressTrackRespVO, error) {
	// 查询订单
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id), s.q.TradeOrder.UserID.Eq(userId)).First()
	if err != nil {
		return nil, core.NewBizError(2002001, "订单不存在") // ORDER_NOT_FOUND code
	}
	return s.getExpressTrackList(ctx, order)
}

// GetExpressTrackListById 获得物流轨迹 (Admin - no UserId check)
func (s *TradeOrderQueryService) GetExpressTrackListById(ctx context.Context, id int64) ([]*resp.ExpressTrackRespVO, error) {
	// 查询订单
	order, err := s.q.TradeOrder.WithContext(ctx).Where(s.q.TradeOrder.ID.Eq(id)).First()
	if err != nil {
		return nil, core.NewBizError(2002001, "订单不存在") // ORDER_NOT_FOUND code
	}
	return s.getExpressTrackList(ctx, order)
}

func (s *TradeOrderQueryService) getExpressTrackList(ctx context.Context, order *trade.TradeOrder) ([]*resp.ExpressTrackRespVO, error) {
	if order.LogisticsID == 0 {
		return []*resp.ExpressTrackRespVO{}, nil
	}
	// 查询物流公司
	express, err := s.deliveryExpressSvc.GetDeliveryExpress(ctx, order.LogisticsID)
	if err != nil || express == nil {
		return nil, core.NewBizError(2002015, "物流公司不存在") // EXPRESS_NOT_EXISTS
	}

	// 获得客户端
	expressClient := s.expressClientFactory.GetDefaultExpressClient()
	if expressClient == nil {
		// Mock Data if no client configured ? Or return error?
		// User wants Kd100 implemented. If configured, it should work.
		return nil, core.NewBizError(500, "物流客户端未配置")
	}

	// 查询物流轨迹
	tracks, err := expressClient.GetExpressTrackList(&client.ExpressTrackQueryReqDTO{
		ExpressCode: express.Code,
		LogisticsNo: order.LogisticsNo,
		Phone:       order.ReceiverMobile,
	})
	if err != nil {
		return nil, err
	}

	// Convert to []resp.ExpressTrackRespVO
	var res []*resp.ExpressTrackRespVO
	for _, t := range tracks {
		res = append(res, &resp.ExpressTrackRespVO{
			Time:    t.Time,
			Content: t.Context,
		})
	}
	return res, nil
}

// GetOrderSummary 获得交易订单统计
func (s *TradeOrderQueryService) GetOrderSummary(ctx context.Context, r *req.TradeOrderPageReq) (*resp.TradeOrderSummaryResp, error) {
	// 1. Construct Query with filters
	q := s.q.TradeOrder.WithContext(ctx)
	if r.UserID != nil {
		q = q.Where(s.q.TradeOrder.UserID.Eq(*r.UserID))
	}
	if r.No != "" {
		q = q.Where(s.q.TradeOrder.No.Like("%" + r.No + "%"))
	}
	// Add other filters as needed to match Page Request behavior if consistent with Java

	// 2. Define result struct for aggregation
	type AggResult struct {
		RefundStatus int   `gorm:"column:refund_status"` // Group by RefundStatus
		Count        int64 `gorm:"column:count"`         // Count(*)
		Price        int64 `gorm:"column:price"`         // Sum(pay_price)
	}
	var results []AggResult

	// 3. Execute Group By Query
	// Select: refund_status, count(*) as count, sum(pay_price) as price
	err := q.Select(
		s.q.TradeOrder.RefundStatus,
		s.q.TradeOrder.ID.Count().As("count"),
		s.q.TradeOrder.PayPrice.Sum().As("price"),
	).Group(s.q.TradeOrder.RefundStatus).Scan(&results)

	if err != nil {
		return nil, err
	}

	// 4. Aggregate into Summary Response
	summary := &resp.TradeOrderSummaryResp{}
	for _, res := range results {
		if res.RefundStatus == 0 { // None (Unrefunded)
			summary.OrderCount += res.Count
			summary.OrderPayPrice += res.Price
		} else { // Any Refund Status (Applied, Successful, etc.)
			summary.AfterSaleCount += res.Count
			summary.AfterSalePrice += res.Price
		}
	}

	return summary, nil
}

// GetOrderLogListByOrderId 获得交易订单日志列表
func (s *TradeOrderQueryService) GetOrderLogListByOrderId(ctx context.Context, orderId int64) ([]*trade.TradeOrderLog, error) {
	return s.q.TradeOrderLog.WithContext(ctx).Where(s.q.TradeOrderLog.OrderID.Eq(orderId)).Order(s.q.TradeOrderLog.CreatedAt.Desc()).Find()
}
