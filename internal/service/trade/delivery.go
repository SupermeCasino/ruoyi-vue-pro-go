package trade

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
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
func (s *DeliveryExpressService) GetDeliveryExpressPage(ctx context.Context, r *req.DeliveryExpressPageReq) (*pagination.PageResult[*trade.TradeDeliveryExpress], error) {
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

	return &pagination.PageResult[*trade.TradeDeliveryExpress]{
		List:  list,
		Total: total,
	}, nil
}

// GetSimpleDeliveryExpressList 获取物流公司精简列表
func (s *DeliveryExpressService) GetSimpleDeliveryExpressList(ctx context.Context) ([]*trade.TradeDeliveryExpress, error) {
	return s.q.TradeDeliveryExpress.WithContext(ctx).Where(s.q.TradeDeliveryExpress.Status.Eq(trade.DeliveryStatusEnabled)).Order(s.q.TradeDeliveryExpress.Sort.Asc()).Find()
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
		OpeningTime:   r.OpeningTime,
		ClosingTime:   r.ClosingTime,
		Latitude:      r.Latitude,
		Longitude:     r.Longitude,
		Status:        r.Status,
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
		"opening_time":   r.OpeningTime,
		"closing_time":   r.ClosingTime,
		"latitude":       r.Latitude,
		"longitude":      r.Longitude,
		"status":         r.Status,
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
func (s *DeliveryPickUpStoreService) GetDeliveryPickUpStorePage(ctx context.Context, r *req.DeliveryPickUpStorePageReq) (*pagination.PageResult[*trade.TradeDeliveryPickUpStore], error) {
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

	list, err := q.Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*trade.TradeDeliveryPickUpStore]{
		List:  list,
		Total: total,
	}, nil
}

// GetSimpleDeliveryPickUpStoreList 获取自提门店精简列表
func (s *DeliveryPickUpStoreService) GetSimpleDeliveryPickUpStoreList(ctx context.Context) ([]*trade.TradeDeliveryPickUpStore, error) {
	return s.q.TradeDeliveryPickUpStore.WithContext(ctx).Where(s.q.TradeDeliveryPickUpStore.Status.Eq(trade.DeliveryStatusEnabled)).Find()
}

// BindDeliveryPickUpStore 绑定自提门店核销员工
// 对齐 Java: DeliveryPickUpStoreServiceImpl.bindDeliveryPickUpStore
func (s *DeliveryPickUpStoreService) BindDeliveryPickUpStore(ctx context.Context, bindReq *req.DeliveryPickUpBindReq) error {
	// 1. 校验门店存在
	store, err := s.q.TradeDeliveryPickUpStore.WithContext(ctx).Where(s.q.TradeDeliveryPickUpStore.ID.Eq(bindReq.ID)).First()
	if err != nil {
		return err
	}
	if store == nil {
		return errors.New("自提门店不存在")
	}

	// 2. 更新核销员工ID列表
	_, err = s.q.TradeDeliveryPickUpStore.WithContext(ctx).Where(s.q.TradeDeliveryPickUpStore.ID.Eq(bindReq.ID)).Updates(map[string]any{
		"verify_user_ids": bindReq.VerifyUserIds,
	})
	return err
}

type DeliveryExpressTemplateService struct {
	q *query.Query
}

func NewDeliveryExpressTemplateService(q *query.Query) *DeliveryExpressTemplateService {
	return &DeliveryExpressTemplateService{q: q}
}

// CreateDeliveryExpressTemplate 创建运费模板
func (s *DeliveryExpressTemplateService) CreateDeliveryExpressTemplate(ctx context.Context, r *req.DeliveryFreightTemplateSaveReq) (int64, error) {
	template := &trade.TradeDeliveryExpressTemplate{
		Name:       r.Name,
		ChargeMode: r.ChargeMode,
		Sort:       r.Sort,
	}

	err := s.q.Transaction(func(tx *query.Query) error {
		// 1. 保存模板
		if err := tx.TradeDeliveryExpressTemplate.WithContext(ctx).Create(template); err != nil {
			return err
		}

		// 2. 保存计费规则
		if len(r.Charges) > 0 {
			var charges []*trade.TradeDeliveryExpressTemplateCharge
			for _, chargeReq := range r.Charges {
				areaIDs := s.convertAreaIDsToString(chargeReq.AreaIDs)
				charges = append(charges, &trade.TradeDeliveryExpressTemplateCharge{
					TemplateID: template.ID,
					AreaIDs:    areaIDs,
					ChargeMode: r.ChargeMode, // Inherit from template
					StartCount: chargeReq.StartCount,
					StartPrice: chargeReq.StartPrice,
					ExtraCount: chargeReq.ExtraCount,
					ExtraPrice: chargeReq.ExtraPrice,
				})
			}
			if err := tx.TradeDeliveryExpressTemplateCharge.WithContext(ctx).Create(charges...); err != nil {
				return err
			}
		}

		// 3. 保存包邮规则
		if len(r.Frees) > 0 {
			var frees []*trade.TradeDeliveryExpressTemplateFree
			for _, freeReq := range r.Frees {
				areaIDs := s.convertAreaIDsToString(freeReq.AreaIDs)
				frees = append(frees, &trade.TradeDeliveryExpressTemplateFree{
					TemplateID: template.ID,
					AreaIDs:    areaIDs,
					FreePrice:  freeReq.FreePrice,
					FreeCount:  freeReq.FreeCount,
				})
			}
			if err := tx.TradeDeliveryExpressTemplateFree.WithContext(ctx).Create(frees...); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return template.ID, nil
}

// UpdateDeliveryExpressTemplate 更新运费模板
func (s *DeliveryExpressTemplateService) UpdateDeliveryExpressTemplate(ctx context.Context, r *req.DeliveryFreightTemplateSaveReq) error {
	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 更新模板
		if _, err := tx.TradeDeliveryExpressTemplate.WithContext(ctx).Where(tx.TradeDeliveryExpressTemplate.ID.Eq(r.ID)).Updates(map[string]interface{}{
			"name":        r.Name,
			"charge_mode": r.ChargeMode,
			"sort":        r.Sort,
		}); err != nil {
			return err
		}

		// 2. 删除旧的计费规则
		if _, err := tx.TradeDeliveryExpressTemplateCharge.WithContext(ctx).Where(tx.TradeDeliveryExpressTemplateCharge.TemplateID.Eq(r.ID)).Delete(); err != nil {
			return err
		}

		// 3. 保存新的计费规则
		if len(r.Charges) > 0 {
			var charges []*trade.TradeDeliveryExpressTemplateCharge
			for _, chargeReq := range r.Charges {
				areaIDs := s.convertAreaIDsToString(chargeReq.AreaIDs)
				charges = append(charges, &trade.TradeDeliveryExpressTemplateCharge{
					TemplateID: r.ID,
					AreaIDs:    areaIDs,
					ChargeMode: r.ChargeMode,
					StartCount: chargeReq.StartCount,
					StartPrice: chargeReq.StartPrice,
					ExtraCount: chargeReq.ExtraCount,
					ExtraPrice: chargeReq.ExtraPrice,
				})
			}
			if err := tx.TradeDeliveryExpressTemplateCharge.WithContext(ctx).Create(charges...); err != nil {
				return err
			}
		}

		// 4. 删除旧的包邮规则
		if _, err := tx.TradeDeliveryExpressTemplateFree.WithContext(ctx).Where(tx.TradeDeliveryExpressTemplateFree.TemplateID.Eq(r.ID)).Delete(); err != nil {
			return err
		}

		// 5. 保存新的包邮规则
		if len(r.Frees) > 0 {
			var frees []*trade.TradeDeliveryExpressTemplateFree
			for _, freeReq := range r.Frees {
				areaIDs := s.convertAreaIDsToString(freeReq.AreaIDs)
				frees = append(frees, &trade.TradeDeliveryExpressTemplateFree{
					TemplateID: r.ID,
					AreaIDs:    areaIDs,
					FreePrice:  freeReq.FreePrice,
					FreeCount:  freeReq.FreeCount,
				})
			}
			if err := tx.TradeDeliveryExpressTemplateFree.WithContext(ctx).Create(frees...); err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteDeliveryExpressTemplate 删除运费模板
func (s *DeliveryExpressTemplateService) DeleteDeliveryExpressTemplate(ctx context.Context, id int64) error {
	return s.q.Transaction(func(tx *query.Query) error {
		// 删除模板
		if _, err := tx.TradeDeliveryExpressTemplate.WithContext(ctx).Where(tx.TradeDeliveryExpressTemplate.ID.Eq(id)).Delete(); err != nil {
			return err
		}
		// 删除计费规则
		if _, err := tx.TradeDeliveryExpressTemplateCharge.WithContext(ctx).Where(tx.TradeDeliveryExpressTemplateCharge.TemplateID.Eq(id)).Delete(); err != nil {
			return err
		}
		// 删除包邮规则
		if _, err := tx.TradeDeliveryExpressTemplateFree.WithContext(ctx).Where(tx.TradeDeliveryExpressTemplateFree.TemplateID.Eq(id)).Delete(); err != nil {
			return err
		}
		return nil
	})
}

// GetDeliveryExpressTemplate 获取运费模板详情
func (s *DeliveryExpressTemplateService) GetDeliveryExpressTemplate(ctx context.Context, id int64) (*resp.DeliveryFreightTemplateResp, error) {
	// 1. 获取模板
	template, err := s.q.TradeDeliveryExpressTemplate.WithContext(ctx).Where(s.q.TradeDeliveryExpressTemplate.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	result := &resp.DeliveryFreightTemplateResp{
		ID:         template.ID,
		Name:       template.Name,
		ChargeMode: template.ChargeMode,
		Sort:       template.Sort,
		CreateTime: template.CreateTime,
	}

	// 2. 获取计费规则
	charges, err := s.q.TradeDeliveryExpressTemplateCharge.WithContext(ctx).Where(s.q.TradeDeliveryExpressTemplateCharge.TemplateID.Eq(id)).Find()
	if err != nil {
		return nil, err
	}
	for _, charge := range charges {
		areaIDs := s.convertAreaIDsToIntSlice(charge.AreaIDs)
		result.Charges = append(result.Charges, resp.DeliveryFreightTemplateChargeResp{
			AreaIDs:    areaIDs,
			StartCount: charge.StartCount,
			StartPrice: charge.StartPrice,
			ExtraCount: charge.ExtraCount,
			ExtraPrice: charge.ExtraPrice,
		})
	}

	// 3. 获取包邮规则
	frees, err := s.q.TradeDeliveryExpressTemplateFree.WithContext(ctx).Where(s.q.TradeDeliveryExpressTemplateFree.TemplateID.Eq(id)).Find()
	if err != nil {
		return nil, err
	}
	for _, free := range frees {
		areaIDs := s.convertAreaIDsToIntSlice(free.AreaIDs)
		result.Frees = append(result.Frees, resp.DeliveryFreightTemplateFreeResp{
			AreaIDs:   areaIDs,
			FreePrice: free.FreePrice,
			FreeCount: free.FreeCount,
		})
	}

	return result, nil
}

// GetDeliveryExpressTemplatePage 获取运费模板分页
func (s *DeliveryExpressTemplateService) GetDeliveryExpressTemplatePage(ctx context.Context, r *req.DeliveryFreightTemplatePageReq) (*pagination.PageResult[*trade.TradeDeliveryExpressTemplate], error) {
	q := s.q.TradeDeliveryExpressTemplate.WithContext(ctx)
	if r.Name != "" {
		q = q.Where(s.q.TradeDeliveryExpressTemplate.Name.Like("%" + r.Name + "%"))
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

	list, err := q.Order(s.q.TradeDeliveryExpressTemplate.Sort.Asc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*trade.TradeDeliveryExpressTemplate]{
		List:  list,
		Total: total,
	}, nil
}

// GetSimpleDeliveryExpressTemplateList 获取所有运费模板精简列表
func (s *DeliveryExpressTemplateService) GetSimpleDeliveryExpressTemplateList(ctx context.Context) ([]*resp.SimpleDeliveryFreightTemplateResp, error) {
	list, err := s.q.TradeDeliveryExpressTemplate.WithContext(ctx).Order(s.q.TradeDeliveryExpressTemplate.Sort.Asc()).Find()
	if err != nil {
		return nil, err
	}
	var res []*resp.SimpleDeliveryFreightTemplateResp
	for _, item := range list {
		res = append(res, &resp.SimpleDeliveryFreightTemplateResp{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	return res, nil
}

// 辅助方法: 转换 AreaIDs 数组为逗号分隔字符串
func (s *DeliveryExpressTemplateService) convertAreaIDsToString(ids []int) string {
	if len(ids) == 0 {
		return ""
	}
	var strIDs []string
	for _, id := range ids {
		strIDs = append(strIDs, strconv.Itoa(id))
	}
	return strings.Join(strIDs, ",")
}

// 辅助方法: 转换逗号分隔字符串为 AreaIDs 数组
func (s *DeliveryExpressTemplateService) convertAreaIDsToIntSlice(str string) []int {
	if str == "" {
		return []int{}
	}
	// Check if it starts with [
	if strings.HasPrefix(str, "[") {
		var ids []int
		json.Unmarshal([]byte(str), &ids)
		return ids
	}

	parts := strings.Split(str, ",")
	var ids []int
	for _, p := range parts {
		if id, err := strconv.Atoi(p); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// CalculateFreight 计算运费
func (s *DeliveryExpressTemplateService) CalculateFreight(ctx context.Context, templateID int64, areaID int, count float64, price int) (int, error) {
	// 1. Get Template
	template, err := s.GetDeliveryExpressTemplate(ctx, templateID)
	if err != nil {
		return 0, err
	}
	if template == nil {
		return 0, errors.New("运费模板不存在")
	}

	// 2. Check Free Rules
	frees, err := s.q.TradeDeliveryExpressTemplateFree.WithContext(ctx).
		Where(s.q.TradeDeliveryExpressTemplateFree.TemplateID.Eq(templateID)).
		Find()
	if err != nil {
		return 0, err
	}
	for _, free := range frees {
		areaIds := s.convertAreaIDsToIntSlice(free.AreaIDs)
		found := false
		for _, aid := range areaIds {
			if aid == areaID {
				found = true
				break
			}
		}
		if found {
			if free.FreeCount > 0 && count >= float64(free.FreeCount) {
				return 0, nil
			}
			if free.FreePrice > 0 && price >= free.FreePrice {
				return 0, nil
			}
		}
	}

	// 3. Calculate Charge
	charges, err := s.q.TradeDeliveryExpressTemplateCharge.WithContext(ctx).
		Where(s.q.TradeDeliveryExpressTemplateCharge.TemplateID.Eq(templateID)).
		Find()
	if err != nil {
		return 0, err
	}

	var matchCharge *trade.TradeDeliveryExpressTemplateCharge
	for _, charge := range charges {
		areaIds := s.convertAreaIDsToIntSlice(charge.AreaIDs)
		for _, aid := range areaIds {
			if aid == areaID {
				matchCharge = charge
				break
			}
		}
		if matchCharge != nil {
			break
		}
	}

	if matchCharge == nil {
		return 0, errors.New("该区域不支持配送")
	}

	if count <= matchCharge.StartCount {
		return matchCharge.StartPrice, nil
	}

	extraNum := math.Ceil((count - matchCharge.StartCount) / matchCharge.ExtraCount)
	totalPrice := matchCharge.StartPrice + int(extraNum)*matchCharge.ExtraPrice
	return totalPrice, nil
}
