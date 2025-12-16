package product

import (
	"context"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"time"
)

// ProductStatisticsRepositoryImpl 商品统计 Repository 实现 - 使用 gorm gen Query
type ProductStatisticsRepositoryImpl struct {
	q *query.Query
}

// NewProductStatisticsRepository 创建商品统计 Repository
func NewProductStatisticsRepository(q *query.Query) service.ProductStatisticsRepository {
	return &ProductStatisticsRepositoryImpl{q: q}
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
