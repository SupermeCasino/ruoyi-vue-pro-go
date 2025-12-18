package repo

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

type AfterSaleLogRepository struct {
	q *query.Query
}

func NewAfterSaleLogRepository(q *query.Query) *AfterSaleLogRepository {
	return &AfterSaleLogRepository{
		q: q,
	}
}

func (r *AfterSaleLogRepository) Create(ctx context.Context, log *trade.AfterSaleLog) error {
	return r.q.AfterSaleLog.WithContext(ctx).Create(log)
}

func (r *AfterSaleLogRepository) GetListByAfterSaleId(ctx context.Context, afterSaleId int64) ([]*trade.AfterSaleLog, error) {
	return r.q.AfterSaleLog.WithContext(ctx).Where(r.q.AfterSaleLog.AfterSaleID.Eq(afterSaleId)).Find()
}
