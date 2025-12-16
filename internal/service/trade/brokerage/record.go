package brokerage

import (
	"context"
	"errors"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp/app/trade"
	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	tradeSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"go.uber.org/zap"
)

type BrokerageRecordService struct {
	q              *query.Query
	logger         *zap.Logger
	tradeConfigSvc *tradeSvc.TradeConfigService
	spuSvc         *product.ProductSpuService
	skuSvc         *product.ProductSkuService
}

func NewBrokerageRecordService(q *query.Query, logger *zap.Logger, tradeConfigSvc *tradeSvc.TradeConfigService, spuSvc *product.ProductSpuService, skuSvc *product.ProductSkuService) *BrokerageRecordService {
	return &BrokerageRecordService{
		q:              q,
		logger:         logger,
		tradeConfigSvc: tradeConfigSvc,
		spuSvc:         spuSvc,
		skuSvc:         skuSvc,
	}
}

// GetSummaryPriceByUserId 获得分销佣金统计
func (s *BrokerageRecordService) GetSummaryPriceByUserId(ctx context.Context, userId int64, bizType int, status int, beginTime, endTime time.Time) (int, error) {
	q := s.q.BrokerageRecord.WithContext(ctx).
		Where(s.q.BrokerageRecord.UserID.Eq(userId))

	if bizType > 0 {
		q = q.Where(s.q.BrokerageRecord.BizType.Eq(bizType))
	}
	if status > 0 {
		q = q.Where(s.q.BrokerageRecord.Status.Eq(status))
	}
	if !beginTime.IsZero() && !endTime.IsZero() {
		q = q.Where(s.q.BrokerageRecord.CreatedAt.Between(beginTime, endTime))
	}

	var sum int
	err := q.Select(s.q.BrokerageRecord.Price.Sum()).Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// GetBrokerageRecord 获得分销记录
func (s *BrokerageRecordService) GetBrokerageRecord(ctx context.Context, id int64) (*brokerage.BrokerageRecord, error) {
	return s.q.BrokerageRecord.WithContext(ctx).Where(s.q.BrokerageRecord.ID.Eq(id)).First()
}

// GetBrokerageRecordPage 获得分销记录分页
func (s *BrokerageRecordService) GetBrokerageRecordPage(ctx context.Context, r *req.BrokerageRecordPageReq) (*pagination.PageResult[*brokerage.BrokerageRecord], error) {
	q := s.q.BrokerageRecord.WithContext(ctx)

	if r.UserID > 0 {
		q = q.Where(s.q.BrokerageRecord.UserID.Eq(r.UserID))
	}
	// BizType 过滤: 请求传入字符串，需要映射为整数
	if r.BizType != "" {
		bizTypeInt := s.mapBizTypeToInt(r.BizType)
		if bizTypeInt > 0 {
			q = q.Where(s.q.BrokerageRecord.BizType.Eq(bizTypeInt))
		}
	}
	if r.Status > 0 {
		q = q.Where(s.q.BrokerageRecord.Status.Eq(r.Status))
	}
	if r.BizID != "" {
		q = q.Where(s.q.BrokerageRecord.BizID.Eq(r.BizID))
	}
	if len(r.CreateTime) == 2 {
		q = q.Where(s.q.BrokerageRecord.CreatedAt.Between(parseTime(r.CreateTime[0]), parseTime(r.CreateTime[1])))
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	offset := (r.PageNo - 1) * r.PageSize
	list, err := q.Limit(r.PageSize).Offset(offset).Order(s.q.BrokerageRecord.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*brokerage.BrokerageRecord]{
		List:  list,
		Total: total,
	}, nil
}

// mapBizTypeToInt 将业务类型字符串映射为整数
func (s *BrokerageRecordService) mapBizTypeToInt(bizType string) int {
	switch bizType {
	case "order":
		return tradeModel.BrokerageRecordBizTypeOrder
	case "withdraw":
		return tradeModel.BrokerageRecordBizTypeWithdraw
	case "withdraw_reject":
		return tradeModel.BrokerageRecordBizTypeWithdrawReject
	default:
		return 0
	}
}

// AddBrokerage 添加分销记录（增加佣金）
func (s *BrokerageRecordService) AddBrokerage(ctx context.Context, userID int64, bizType string, bizID string, price int, title string) error {
	// 1. 更新用户佣金余额
	u := s.q.BrokerageUser
	if _, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).UpdateSimple(u.BrokeragePrice.Add(price)); err != nil {
		return err
	}

	// 2. 创建佣金记录
	bizTypeInt := s.mapBizTypeToInt(bizType)
	record := &brokerage.BrokerageRecord{
		UserID:      userID,
		BizType:     bizTypeInt,
		BizID:       bizID,
		Price:       price,
		Title:       title,
		Description: title,
		Status:      tradeModel.BrokerageRecordStatusSettlement,
		CreatedAt:   time.Now(),
		TotalPrice:  price,
	}
	return s.q.BrokerageRecord.WithContext(ctx).Create(record)
}

// ReduceBrokerageForWithdraw 提现扣减佣金
func (s *BrokerageRecordService) ReduceBrokerageForWithdraw(ctx context.Context, userID int64, bizID string, price int) error {
	// 1. 校验佣金余额并更新（原子操作）
	u := s.q.BrokerageUser
	info, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).UpdateSimple(
		u.BrokeragePrice.Sub(price),
		u.FrozenPrice.Add(price),
	)
	if err != nil {
		return err
	}
	if info.RowsAffected == 0 {
		return errors.New("佣金不足")
	}

	// 2. 创建佣金记录
	record := &brokerage.BrokerageRecord{
		UserID:      userID,
		BizType:     tradeModel.BrokerageRecordBizTypeWithdraw,
		BizID:       bizID,
		Price:       -price, // 负数表示扣减
		Title:       "佣金提现",
		Description: "佣金提现",
		Status:      tradeModel.BrokerageRecordStatusSettlement,
		CreatedAt:   time.Now(),
		TotalPrice:  price,
	}
	return s.q.BrokerageRecord.WithContext(ctx).Create(record)
}

// CalculateProductBrokeragePrice 计算商品佣金
func (s *BrokerageRecordService) CalculateProductBrokeragePrice(ctx context.Context, userId int64, spuId int64) (*trade.AppBrokerageProductPriceRespVO, error) {
	resp := &trade.AppBrokerageProductPriceRespVO{
		BrokerageEnabled: false,
		BrokeragePrice:   0,
	}

	// 1. 校验分销功能是否开启
	config, err := s.tradeConfigSvc.GetTradeConfig(ctx)
	if err != nil || config == nil || !config.BrokerageEnabled {
		return resp, nil
	}

	// 2. 校验用户是否有分销资格
	user, err := s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(userId)).First()
	if err != nil || user == nil || !user.BrokerageEnabled {
		return resp, nil
	}
	resp.BrokerageEnabled = true

	// 3. 校验商品是否存在
	spu, err := s.spuSvc.GetSpu(ctx, spuId)
	if err != nil || spu == nil {
		return resp, nil
	}

	// 4. 计算佣金
	skus, err := s.skuSvc.GetSkuListBySpuId(ctx, spuId)
	if err != nil {
		return resp, nil
	}

	minPrice := 0
	maxPrice := 0
	percent := config.BrokerageFirstPercent

	for _, sku := range skus {
		var brokeragePrice int
		if spu.SubCommissionType {
			// 商品单独分佣模式：使用 SKU 的固定佣金
			brokeragePrice = sku.FirstBrokeragePrice
		} else {
			// 全局分佣模式：根据商品价格比例计算
			brokeragePrice = sku.Price * percent / 100
		}

		if minPrice == 0 || brokeragePrice < minPrice {
			minPrice = brokeragePrice
		}
		if brokeragePrice > maxPrice {
			maxPrice = brokeragePrice
		}
	}

	// 使用最大佣金作为展示值
	resp.BrokeragePrice = maxPrice

	return resp, nil
}

// GetBrokerageUserRankPageByPrice 获得分销用户排行分页（基于佣金）
func (s *BrokerageRecordService) GetBrokerageUserRankPageByPrice(ctx context.Context, r *req.AppBrokerageUserRankPageReq) (*pagination.PageResult[*trade.AppBrokerageUserRankByPriceRespVO], error) {
	// 解析时间范围
	var beginTime, endTime time.Time
	if len(r.Times) >= 2 {
		beginTime = parseTime(r.Times[0])
		endTime = parseTime(r.Times[1])
	}

	// 使用 Gen 生成的字段和表名
	br := s.q.BrokerageRecord
	tableName := br.TableName()
	userIDCol := br.UserID.ColumnName().String()
	priceCol := br.Price.ColumnName().String()
	bizTypeCol := br.BizType.ColumnName().String()
	statusCol := br.Status.ColumnName().String()
	createTimeCol := br.CreatedAt.ColumnName().String()

	// 构建基础查询条件统计总数
	q := br.WithContext(ctx).
		Where(br.BizType.Eq(tradeModel.BrokerageRecordBizTypeOrder)).
		Where(br.Status.Eq(tradeModel.BrokerageRecordStatusSettlement))

	if !beginTime.IsZero() && !endTime.IsZero() {
		q = q.Where(br.CreatedAt.Between(beginTime, endTime))
	}

	// 获取总数 (不同用户数)
	total, err := q.Distinct(br.UserID).Count()
	if err != nil {
		return nil, err
	}

	// 分组查询需要使用原生 GORM，因为 Gen 不支持 GROUP BY 聚合
	// 使用 Gen 生成的字段名确保类型安全
	db := br.WithContext(ctx).UnderlyingDB()
	offset := (r.PageNo - 1) * r.PageSize

	type RankResult struct {
		UserID int64 `gorm:"column:user_id"`
		Price  int   `gorm:"column:price"`
	}

	selectClause := userIDCol + ", SUM(" + priceCol + ") as price"
	query := db.Table(tableName).
		Select(selectClause).
		Where(bizTypeCol+" = ? AND "+statusCol+" = ?", tradeModel.BrokerageRecordBizTypeOrder, tradeModel.BrokerageRecordStatusSettlement).
		Where("deleted = 0")

	if !beginTime.IsZero() && !endTime.IsZero() {
		query = query.Where(createTimeCol+" BETWEEN ? AND ?", beginTime, endTime)
	}

	var results []RankResult
	query.Group(userIDCol).
		Order("price DESC").
		Limit(r.PageSize).
		Offset(offset).
		Scan(&results)

	// 转换为 VO
	list := make([]*trade.AppBrokerageUserRankByPriceRespVO, len(results))
	for i, item := range results {
		list[i] = &trade.AppBrokerageUserRankByPriceRespVO{
			ID:             item.UserID,
			BrokeragePrice: item.Price,
		}
	}

	return &pagination.PageResult[*trade.AppBrokerageUserRankByPriceRespVO]{
		List:  list,
		Total: total,
	}, nil
}

// GetUserRankByPrice 获得分销用户排行（基于佣金）
func (s *BrokerageRecordService) GetUserRankByPrice(ctx context.Context, userId int64, times []time.Time) (int, error) {
	var beginTime, endTime time.Time
	if len(times) >= 2 {
		beginTime = times[0]
		endTime = times[1]
	}

	// 1. 获取用户的分销佣金总额
	userPrice, err := s.GetSummaryPriceByUserId(ctx, userId, tradeModel.BrokerageRecordBizTypeOrder, tradeModel.BrokerageRecordStatusSettlement, beginTime, endTime)
	if err != nil {
		return 0, err
	}

	// 使用 Gen 生成的字段和表名
	br := s.q.BrokerageRecord
	tableName := br.TableName()
	userIDCol := br.UserID.ColumnName().String()
	priceCol := br.Price.ColumnName().String()
	bizTypeCol := br.BizType.ColumnName().String()
	statusCol := br.Status.ColumnName().String()
	createTimeCol := br.CreatedAt.ColumnName().String()

	// 2. 获取比用户佣金高的用户数量
	db := br.WithContext(ctx).UnderlyingDB()

	// 子查询：获取每个用户的佣金总额大于当前用户的数量
	subQuery := db.Table(tableName).
		Select(userIDCol+", SUM("+priceCol+") as total_price").
		Where(bizTypeCol+" = ? AND "+statusCol+" = ?", tradeModel.BrokerageRecordBizTypeOrder, tradeModel.BrokerageRecordStatusSettlement).
		Where("deleted = 0")
	if !beginTime.IsZero() && !endTime.IsZero() {
		subQuery = subQuery.Where(createTimeCol+" BETWEEN ? AND ?", beginTime, endTime)
	}
	subQuery = subQuery.Group(userIDCol).Having("SUM("+priceCol+") > ?", userPrice)

	var greaterCount int64
	db.Table("(?) as ranked", subQuery).Count(&greaterCount)

	// 3. 返回排名 (比自己高的人数 + 1)
	return int(greaterCount) + 1, nil
}
