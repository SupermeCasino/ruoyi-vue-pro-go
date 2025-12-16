package product

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/gorm"
)

// ProductStatisticsRepositoryImpl 商品统计 Repository 实现 - 使用 gorm gen Query
type ProductStatisticsRepositoryImpl struct {
	q  *query.Query
	db *gorm.DB
}

// NewProductStatisticsRepository 创建商品统计 Repository
func NewProductStatisticsRepository(q *query.Query, db *gorm.DB) service.ProductStatisticsRepository {
	return &ProductStatisticsRepositoryImpl{q: q, db: db}
}

// GetByDateRange 查询指定日期范围的商品统计数据
func (r *ProductStatisticsRepositoryImpl) GetByDateRange(ctx context.Context, beginTime, endTime time.Time) ([]*resp.ProductStatisticsRespVO, error) {
	ps := r.q.ProductStatistics

	// 查询指定时间范围内的统计数据，按 SPU 聚合
	var results []struct {
		SpuID           int64 `gorm:"column:spu_id"`
		BrowseCount     int64 `gorm:"column:browse_count"`
		FavoriteCount   int64 `gorm:"column:favorite_count"`
		CartCount       int64 `gorm:"column:cart_count"`
		OrderCount      int64 `gorm:"column:order_count"`
		BuyCount        int64 `gorm:"column:buy_count"`
		BuyPrice        int64 `gorm:"column:buy_price"`
		AfterSaleCount  int64 `gorm:"column:after_sale_count"`
		AfterSaleRefund int64 `gorm:"column:after_sale_refund"`
	}

	err := ps.WithContext(ctx).
		Select(
			ps.SpuID,
			ps.BrowseCount.Sum().As("browse_count"),
			ps.FavoriteCount.Sum().As("favorite_count"),
			ps.CartCount.Sum().As("cart_count"),
			ps.OrderCount.Sum().As("order_count"),
			ps.OrderPayCount.Sum().As("buy_count"),
			ps.OrderPayPrice.Sum().As("buy_price"),
			ps.AfterSaleCount.Sum().As("after_sale_count"),
			ps.AfterSaleRefundPrice.Sum().As("after_sale_refund"),
		).
		Where(ps.Time.Between(beginTime, endTime)).
		// Where(ps.Deleted.Eq(false)). // Soft delete handled by GORM
		Group(ps.SpuID).
		Scan(&results)
	if err != nil {
		return nil, err
	}

	// 转换为 VO
	voList := make([]*resp.ProductStatisticsRespVO, 0, len(results))
	for _, r := range results {
		voList = append(voList, &resp.ProductStatisticsRespVO{
			SpuID:         r.SpuID,
			BrowseCount:   r.BrowseCount,
			FavoriteCount: r.FavoriteCount,
			BuyCount:      r.BuyCount,
			BuyPrice:      r.BuyPrice,
		})
	}

	return voList, nil
}

// GetSummaryByDateRange 获取指定日期范围的汇总统计数据
func (r *ProductStatisticsRepositoryImpl) GetSummaryByDateRange(ctx context.Context, beginTime, endTime time.Time) (*resp.ProductStatisticsRespVO, error) {
	ps := r.q.ProductStatistics

	var result struct {
		BrowseCount   int64 `gorm:"column:browse_count"`
		FavoriteCount int64 `gorm:"column:favorite_count"`
		CartCount     int64 `gorm:"column:cart_count"`
		BuyCount      int64 `gorm:"column:buy_count"`
		BuyPrice      int64 `gorm:"column:buy_price"`
		CommentCount  int64 `gorm:"column:comment_count"`
	}

	err := ps.WithContext(ctx).
		Select(
			ps.BrowseCount.Sum().As("browse_count"),
			ps.FavoriteCount.Sum().As("favorite_count"),
			ps.CartCount.Sum().As("cart_count"),
			ps.OrderPayCount.Sum().As("buy_count"),
			ps.OrderPayPrice.Sum().As("buy_price"),
		).
		Where(ps.Time.Between(beginTime, endTime)).
		Scan(&result)
	if err != nil {
		return nil, err
	}

	return &resp.ProductStatisticsRespVO{
		BrowseCount:   result.BrowseCount,
		FavoriteCount: result.FavoriteCount,
		BuyCount:      result.BuyCount,
		BuyPrice:      result.BuyPrice,
	}, nil
}

// GetPageGroupBySpuId 分页获取按 SPU 分组的统计数据
func (r *ProductStatisticsRepositoryImpl) GetPageGroupBySpuId(ctx context.Context, reqVO *req.ProductStatisticsReqVO, pageParam *pagination.PageParam) (*pagination.PageResult[*resp.ProductStatisticsRespVO], error) {
	// 先获取所有数据，然后内存分页
	list, err := r.GetByDateRange(ctx, reqVO.Times[0], reqVO.Times[1])
	if err != nil {
		return nil, err
	}

	total := int64(len(list))
	start := (pageParam.PageNo - 1) * pageParam.PageSize
	end := start + pageParam.PageSize
	if start > int(total) {
		start = int(total)
	}
	if end > int(total) {
		end = int(total)
	}

	return &pagination.PageResult[*resp.ProductStatisticsRespVO]{
		List:  list[start:end],
		Total: total,
	}, nil
}

// CountByDateRange 统计指定日期范围内的记录数
func (r *ProductStatisticsRepositoryImpl) CountByDateRange(ctx context.Context, beginTime, endTime time.Time) (int64, error) {
	ps := r.q.ProductStatistics
	return ps.WithContext(ctx).Where(ps.Time.Between(beginTime, endTime)).Count()
}

// StatisticsProductByDateRange 统计指定日期范围内的商品数据并入库
// 对应 Java: ProductStatisticsMapper.selectStatisticsResultPageByTimeBetween
func (r *ProductStatisticsRepositoryImpl) StatisticsProductByDateRange(ctx context.Context, date time.Time, beginTime, endTime time.Time) error {
	db := r.db

	// 分页统计，避免商品表数据较多时出现超时问题
	const pageSize = 100
	offset := 0

	for {
		// 使用原生 SQL 查询统计数据，对应 Java XML 映射
		var records []struct {
			SpuID                int64 `gorm:"column:spu_id"`
			BrowseCount          int   `gorm:"column:browse_count"`
			BrowseUserCount      int   `gorm:"column:browse_user_count"`
			FavoriteCount        int   `gorm:"column:favorite_count"`
			CartCount            int   `gorm:"column:cart_count"`
			OrderCount           int   `gorm:"column:order_count"`
			OrderPayCount        int   `gorm:"column:order_pay_count"`
			OrderPayPrice        int64 `gorm:"column:order_pay_price"`
			AfterSaleCount       int   `gorm:"column:after_sale_count"`
			AfterSaleRefundPrice int64 `gorm:"column:after_sale_refund_price"`
		}

		sql := `
			SELECT spu.id AS spu_id
				-- 浏览量：一个用户可以有多次
				, (SELECT COUNT(1) FROM product_browse_history 
				   WHERE spu_id = spu.id AND create_time BETWEEN ? AND ?) AS browse_count
				-- 访客量：按用户去重计数
				, (SELECT COUNT(DISTINCT user_id) FROM product_browse_history 
				   WHERE spu_id = spu.id AND create_time BETWEEN ? AND ?) AS browse_user_count
				-- 收藏数量：按用户去重计数
				, (SELECT COUNT(DISTINCT user_id) FROM product_favorite 
				   WHERE spu_id = spu.id AND create_time BETWEEN ? AND ?) AS favorite_count
				-- 加购数量：按用户去重计数
				, (SELECT COUNT(DISTINCT user_id) FROM trade_cart 
				   WHERE spu_id = spu.id AND create_time BETWEEN ? AND ?) AS cart_count
				-- 下单件数
				, (SELECT IFNULL(SUM(count), 0) FROM trade_order_item 
				   WHERE spu_id = spu.id AND create_time BETWEEN ? AND ?) AS order_count
				-- 支付件数
				, (SELECT IFNULL(SUM(item.count), 0) FROM trade_order_item item 
				   JOIN trade_order o ON item.order_id = o.id 
				   WHERE spu_id = spu.id AND o.pay_status = TRUE 
				   AND item.create_time BETWEEN ? AND ?) AS order_pay_count
				-- 支付金额
				, (SELECT IFNULL(SUM(item.pay_price), 0) FROM trade_order_item item 
				   JOIN trade_order o ON item.order_id = o.id 
				   WHERE spu_id = spu.id AND o.pay_status = TRUE 
				   AND item.create_time BETWEEN ? AND ?) AS order_pay_price
				-- 退款件数
				, (SELECT IFNULL(SUM(count), 0) FROM trade_after_sale 
				   WHERE spu_id = spu.id AND refund_time IS NOT NULL 
				   AND create_time BETWEEN ? AND ?) AS after_sale_count
				-- 退款金额
				, (SELECT IFNULL(SUM(refund_price), 0) FROM trade_after_sale 
				   WHERE spu_id = spu.id AND refund_time IS NOT NULL 
				   AND create_time BETWEEN ? AND ?) AS after_sale_refund_price
			FROM product_spu spu
			WHERE spu.deleted = 0
			ORDER BY spu.id
			LIMIT ? OFFSET ?
		`

		err := db.WithContext(ctx).Raw(sql,
			beginTime, endTime, // browse_count
			beginTime, endTime, // browse_user_count
			beginTime, endTime, // favorite_count
			beginTime, endTime, // cart_count
			beginTime, endTime, // order_count
			beginTime, endTime, // order_pay_count
			beginTime, endTime, // order_pay_price
			beginTime, endTime, // after_sale_count
			beginTime, endTime, // after_sale_refund_price
			pageSize, offset,
		).Scan(&records).Error
		if err != nil {
			return err
		}

		// 如果没有数据，退出循环
		if len(records) == 0 {
			break
		}

		// 计算访客支付转化率并批量插入
		dateOnly := date.Format("2006-01-02")
		for _, record := range records {
			// 计算访客支付转化率（百分比）
			browseConvertPercent := 0
			if record.BrowseUserCount > 0 {
				browseConvertPercent = 100 * record.OrderPayCount / record.BrowseUserCount
			}

			// 插入统计记录
			insertSQL := `
				INSERT INTO product_statistics 
					(spu_id, time, browse_count, browse_user_count, favorite_count, cart_count, 
					 order_count, order_pay_count, order_pay_price, after_sale_count, 
					 after_sale_refund_price, browse_convert_percent, create_time, update_time, deleted, tenant_id)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), 0, 0)
			`
			err := db.WithContext(ctx).Exec(insertSQL,
				record.SpuID, dateOnly,
				record.BrowseCount, record.BrowseUserCount, record.FavoriteCount, record.CartCount,
				record.OrderCount, record.OrderPayCount, record.OrderPayPrice,
				record.AfterSaleCount, record.AfterSaleRefundPrice, browseConvertPercent,
			).Error
			if err != nil {
				return err
			}
		}

		offset += pageSize
	}

	return nil
}
