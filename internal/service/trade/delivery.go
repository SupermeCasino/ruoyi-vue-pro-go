package trade

import (
	"context"

	"backend-go/internal/api/req"
	"backend-go/internal/model/trade"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
)

type DeliveryExpressService struct {
	q *query.Query
}

func NewDeliveryExpressService(q *query.Query) *DeliveryExpressService {
	return &DeliveryExpressService{q: q}
}

// CreateDeliveryExpress 创建物流公司
func (s *DeliveryExpressService) CreateDeliveryExpress(ctx context.Context, r *req.DeliveryExpressSaveReq) (int64, error) {
	express := &trade.TradeDeliveryExpress{
		Code:   r.Code,
		Name:   r.Name,
		Logo:   r.Logo,
		Sort:   r.Sort,
		Status: r.Status,
	}
	if err := s.q.TradeDeliveryExpress.WithContext(ctx).Create(express); err != nil {
		return 0, err
	}
	return express.ID, nil
}

// UpdateDeliveryExpress 更新物流公司
func (s *DeliveryExpressService) UpdateDeliveryExpress(ctx context.Context, r *req.DeliveryExpressSaveReq) error {
	_, err := s.q.TradeDeliveryExpress.WithContext(ctx).Where(s.q.TradeDeliveryExpress.ID.Eq(*r.ID)).Updates(map[string]interface{}{
		"code":   r.Code,
		"name":   r.Name,
		"logo":   r.Logo,
		"sort":   r.Sort,
		"status": r.Status,
	})
	return err
}

// DeleteDeliveryExpress 删除物流公司
func (s *DeliveryExpressService) DeleteDeliveryExpress(ctx context.Context, id int64) error {
	_, err := s.q.TradeDeliveryExpress.WithContext(ctx).Where(s.q.TradeDeliveryExpress.ID.Eq(id)).Delete()
	return err
}

// GetDeliveryExpress 获取物流公司
func (s *DeliveryExpressService) GetDeliveryExpress(ctx context.Context, id int64) (*trade.TradeDeliveryExpress, error) {
	return s.q.TradeDeliveryExpress.WithContext(ctx).Where(s.q.TradeDeliveryExpress.ID.Eq(id)).First()
}

// GetDeliveryExpressPage 获取物流公司分页
func (s *DeliveryExpressService) GetDeliveryExpressPage(ctx context.Context, r *req.DeliveryExpressPageReq) (*core.PageResult[*trade.TradeDeliveryExpress], error) {
	q := s.q.TradeDeliveryExpress.WithContext(ctx)
	if r.Code != "" {
		q = q.Where(s.q.TradeDeliveryExpress.Code.Like("%" + r.Code + "%"))
	}
	if r.Name != "" {
		q = q.Where(s.q.TradeDeliveryExpress.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.TradeDeliveryExpress.Status.Eq(*r.Status))
	}

	pageNo := r.PageNo
	pageSize := r.PageSize
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Order(s.q.TradeDeliveryExpress.Sort.Asc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*trade.TradeDeliveryExpress]{
		List:  list,
		Total: total,
	}, nil
}

// GetSimpleDeliveryExpressList 获取物流公司精简列表
func (s *DeliveryExpressService) GetSimpleDeliveryExpressList(ctx context.Context) ([]*trade.TradeDeliveryExpress, error) {
	return s.q.TradeDeliveryExpress.WithContext(ctx).Where(s.q.TradeDeliveryExpress.Status.Eq(1)).Order(s.q.TradeDeliveryExpress.Sort.Asc()).Find()
}

type DeliveryPickUpStoreService struct {
	q *query.Query
}

func NewDeliveryPickUpStoreService(q *query.Query) *DeliveryPickUpStoreService {
	return &DeliveryPickUpStoreService{q: q}
}

// CreateDeliveryPickUpStore 创建自提门店
func (s *DeliveryPickUpStoreService) CreateDeliveryPickUpStore(ctx context.Context, r *req.DeliveryPickUpStoreSaveReq) (int64, error) {
	store := &trade.TradeDeliveryPickUpStore{
		Name:          r.Name,
		Introduction:  r.Introduction,
		Phone:         r.Phone,
		AreaID:        r.AreaID,
		DetailAddress: r.DetailAddress,
		Logo:          r.Logo,
		Latitude:      r.Latitude,
		Longitude:     r.Longitude,
		Status:        r.Status,
		Sort:          r.Sort,
	}
	if err := s.q.TradeDeliveryPickUpStore.WithContext(ctx).Create(store); err != nil {
		return 0, err
	}
	return store.ID, nil
}

// UpdateDeliveryPickUpStore 更新自提门店
func (s *DeliveryPickUpStoreService) UpdateDeliveryPickUpStore(ctx context.Context, r *req.DeliveryPickUpStoreSaveReq) error {
	_, err := s.q.TradeDeliveryPickUpStore.WithContext(ctx).Where(s.q.TradeDeliveryPickUpStore.ID.Eq(*r.ID)).Updates(map[string]interface{}{
		"name":           r.Name,
		"introduction":   r.Introduction,
		"phone":          r.Phone,
		"area_id":        r.AreaID,
		"detail_address": r.DetailAddress,
		"logo":           r.Logo,
		"latitude":       r.Latitude,
		"longitude":      r.Longitude,
		"status":         r.Status,
		"sort":           r.Sort,
	})
	return err
}

// DeleteDeliveryPickUpStore 删除自提门店
func (s *DeliveryPickUpStoreService) DeleteDeliveryPickUpStore(ctx context.Context, id int64) error {
	_, err := s.q.TradeDeliveryPickUpStore.WithContext(ctx).Where(s.q.TradeDeliveryPickUpStore.ID.Eq(id)).Delete()
	return err
}

// GetDeliveryPickUpStore 获取自提门店
func (s *DeliveryPickUpStoreService) GetDeliveryPickUpStore(ctx context.Context, id int64) (*trade.TradeDeliveryPickUpStore, error) {
	return s.q.TradeDeliveryPickUpStore.WithContext(ctx).Where(s.q.TradeDeliveryPickUpStore.ID.Eq(id)).First()
}

// GetDeliveryPickUpStorePage 获取自提门店分页
func (s *DeliveryPickUpStoreService) GetDeliveryPickUpStorePage(ctx context.Context, r *req.DeliveryPickUpStorePageReq) (*core.PageResult[*trade.TradeDeliveryPickUpStore], error) {
	q := s.q.TradeDeliveryPickUpStore.WithContext(ctx)
	if r.Name != "" {
		q = q.Where(s.q.TradeDeliveryPickUpStore.Name.Like("%" + r.Name + "%"))
	}
	if r.Phone != "" {
		q = q.Where(s.q.TradeDeliveryPickUpStore.Phone.Like("%" + r.Phone + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.TradeDeliveryPickUpStore.Status.Eq(*r.Status))
	}

	pageNo := r.PageNo
	pageSize := r.PageSize
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Order(s.q.TradeDeliveryPickUpStore.Sort.Asc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*trade.TradeDeliveryPickUpStore]{
		List:  list,
		Total: total,
	}, nil
}
