package pay

import (
	"context"
	"time"

	reqPay "github.com/wxlbd/ruoyi-mall-go/internal/api/req/pay"
	modelPay "github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PayTransferRepository interface {
	SelectPage(ctx context.Context, req *reqPay.PayTransferPageReq) (*pagination.PageResult[*modelPay.PayTransfer], error)
	SelectById(ctx context.Context, id int64) (*modelPay.PayTransfer, error)
	Create(ctx context.Context, transfer *modelPay.PayTransfer) error
	Update(ctx context.Context, transfer *modelPay.PayTransfer) error
	UpdateByIdAndStatus(ctx context.Context, id int64, whereStatuses []int, updateObj *modelPay.PayTransfer) (int64, error)
	SelectByAppIdAndMerchantTransferId(ctx context.Context, appId int64, merchantTransferId string) (*modelPay.PayTransfer, error)
	SelectByAppIdAndNo(ctx context.Context, appId int64, no string) (*modelPay.PayTransfer, error)
	SelectListByStatus(ctx context.Context, statuses []int) ([]*modelPay.PayTransfer, error)
	SelectByNo(ctx context.Context, no string) (*modelPay.PayTransfer, error)
}

type PayTransferRepositoryImpl struct {
	q *query.Query
}

func NewPayTransferRepository(q *query.Query) PayTransferRepository {
	return &PayTransferRepositoryImpl{q: q}
}

func (r *PayTransferRepositoryImpl) SelectPage(ctx context.Context, req *reqPay.PayTransferPageReq) (*pagination.PageResult[*modelPay.PayTransfer], error) {
	q := r.q.PayTransfer.WithContext(ctx)

	if req.No != "" {
		q = q.Where(r.q.PayTransfer.No.Eq(req.No))
	}
	if req.AppID != 0 {
		q = q.Where(r.q.PayTransfer.AppID.Eq(req.AppID))
	}
	if req.ChannelCode != "" {
		q = q.Where(r.q.PayTransfer.ChannelCode.Eq(req.ChannelCode))
	}
	if req.MerchantOrderId != "" {
		q = q.Where(r.q.PayTransfer.MerchantTransferID.Eq(req.MerchantOrderId))
	}
	if req.Status != nil {
		q = q.Where(r.q.PayTransfer.Status.Eq(*req.Status))
	}
	if req.UserName != "" {
		q = q.Where(r.q.PayTransfer.UserName.Like("%" + req.UserName + "%"))
	}
	if req.UserAccount != "" {
		q = q.Where(r.q.PayTransfer.UserAccount.Like("%" + req.UserAccount + "%"))
	}
	if req.ChannelTransferNo != "" {
		q = q.Where(r.q.PayTransfer.ChannelTransferNo.Eq(req.ChannelTransferNo))
	}
	// 对齐 Java: .betweenIfPresent(PayTransferDO::getCreateTime, reqVO.getCreateTime())
	if len(req.CreateTime) == 2 {
		// 解析时间字符串 "2006-01-02 15:04:05"
		const layout = "2006-01-02 15:04:05"
		if startTime, err := time.Parse(layout, req.CreateTime[0]); err == nil {
			if endTime, err := time.Parse(layout, req.CreateTime[1]); err == nil {
				q = q.Where(r.q.PayTransfer.CreatedAt.Between(startTime, endTime))
			}
		}
	}
	q = q.Order(r.q.PayTransfer.ID.Desc())

	list, count, err := q.FindByPage(req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}
	return pagination.NewPageResult(list, count), nil
}

func (r *PayTransferRepositoryImpl) SelectById(ctx context.Context, id int64) (*modelPay.PayTransfer, error) {
	return r.q.PayTransfer.WithContext(ctx).Where(r.q.PayTransfer.ID.Eq(id)).First()
}

func (r *PayTransferRepositoryImpl) Create(ctx context.Context, transfer *modelPay.PayTransfer) error {
	return r.q.PayTransfer.WithContext(ctx).Create(transfer)
}

func (r *PayTransferRepositoryImpl) Update(ctx context.Context, transfer *modelPay.PayTransfer) error {
	return r.q.PayTransfer.WithContext(ctx).Save(transfer)
}

func (r *PayTransferRepositoryImpl) UpdateByIdAndStatus(ctx context.Context, id int64, whereStatuses []int, updateObj *modelPay.PayTransfer) (int64, error) {
	q := r.q.PayTransfer.WithContext(ctx).Where(r.q.PayTransfer.ID.Eq(id)).Where(r.q.PayTransfer.Status.In(whereStatuses...))
	result, err := q.Updates(updateObj)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

func (r *PayTransferRepositoryImpl) SelectByAppIdAndMerchantTransferId(ctx context.Context, appId int64, merchantTransferId string) (*modelPay.PayTransfer, error) {
	return r.q.PayTransfer.WithContext(ctx).
		Where(r.q.PayTransfer.AppID.Eq(appId)).
		Where(r.q.PayTransfer.MerchantTransferID.Eq(merchantTransferId)).
		First()
}

func (r *PayTransferRepositoryImpl) SelectByAppIdAndNo(ctx context.Context, appId int64, no string) (*modelPay.PayTransfer, error) {
	return r.q.PayTransfer.WithContext(ctx).
		Where(r.q.PayTransfer.AppID.Eq(appId)).
		Where(r.q.PayTransfer.No.Eq(no)).
		First()
}

func (r *PayTransferRepositoryImpl) SelectListByStatus(ctx context.Context, statuses []int) ([]*modelPay.PayTransfer, error) {
	return r.q.PayTransfer.WithContext(ctx).Where(r.q.PayTransfer.Status.In(statuses...)).Find()
}

func (r *PayTransferRepositoryImpl) SelectByNo(ctx context.Context, no string) (*modelPay.PayTransfer, error) {
	return r.q.PayTransfer.WithContext(ctx).Where(r.q.PayTransfer.No.Eq(no)).First()
}
