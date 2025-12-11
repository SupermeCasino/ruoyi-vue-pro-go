package pay

import (
	"backend-go/internal/api/req"
	"backend-go/internal/model/pay"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"
)

type PayRefundService struct {
	q *query.Query
}

func NewPayRefundService(q *query.Query) *PayRefundService {
	return &PayRefundService{q: q}
}

// GetRefund 获得退款订单
func (s *PayRefundService) GetRefund(ctx context.Context, id int64) (*pay.PayRefund, error) {
	return s.q.PayRefund.WithContext(ctx).Where(s.q.PayRefund.ID.Eq(id)).First()
}

// GetRefundPage 获得退款订单分页
func (s *PayRefundService) GetRefundPage(ctx context.Context, req *req.PayRefundPageReq) (*core.PageResult[*pay.PayRefund], error) {
	q := s.q.PayRefund.WithContext(ctx)
	if req.AppID > 0 {
		q = q.Where(s.q.PayRefund.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(s.q.PayRefund.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(s.q.PayRefund.MerchantOrderId.Eq(req.MerchantOrderId))
	}
	if req.MerchantRefundId != "" {
		q = q.Where(s.q.PayRefund.MerchantRefundId.Eq(req.MerchantRefundId))
	}
	if req.ChannelOrderNo != "" {
		q = q.Where(s.q.PayRefund.ChannelOrderNo.Eq(req.ChannelOrderNo))
	}
	if req.ChannelRefundNo != "" {
		q = q.Where(s.q.PayRefund.ChannelRefundNo.Eq(req.ChannelRefundNo))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayRefund.Status.Eq(*req.Status))
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}
	list, err := q.Limit(req.GetLimit()).Offset(req.GetOffset()).Order(s.q.PayRefund.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*pay.PayRefund]{
		List:  list,
		Total: total,
	}, nil
}

// GetRefundList 获得退款订单列表 (Export)
func (s *PayRefundService) GetRefundList(ctx context.Context, req *req.PayRefundExportReq) ([]*pay.PayRefund, error) {
	q := s.q.PayRefund.WithContext(ctx)
	if req.AppID > 0 {
		q = q.Where(s.q.PayRefund.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(s.q.PayRefund.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(s.q.PayRefund.MerchantOrderId.Eq(req.MerchantOrderId))
	}
	if req.MerchantRefundId != "" {
		q = q.Where(s.q.PayRefund.MerchantRefundId.Eq(req.MerchantRefundId))
	}
	if req.ChannelOrderNo != "" {
		q = q.Where(s.q.PayRefund.ChannelOrderNo.Eq(req.ChannelOrderNo))
	}
	if req.ChannelRefundNo != "" {
		q = q.Where(s.q.PayRefund.ChannelRefundNo.Eq(req.ChannelRefundNo))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayRefund.Status.Eq(*req.Status))
	}
	return q.Order(s.q.PayRefund.ID.Desc()).Find()
}
