package iot

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type OtaFirmwareRepositoryImpl struct {
	q *query.Query
}

func NewOtaFirmwareRepository(q *query.Query) iotsvc.OtaFirmwareRepository {
	return &OtaFirmwareRepositoryImpl{q: q}
}

func (r *OtaFirmwareRepositoryImpl) Create(ctx context.Context, firmware *model.IotOtaFirmwareDO) error {
	return r.q.IotOtaFirmwareDO.WithContext(ctx).Create(firmware)
}

func (r *OtaFirmwareRepositoryImpl) Update(ctx context.Context, firmware *model.IotOtaFirmwareDO) error {
	_, err := r.q.IotOtaFirmwareDO.WithContext(ctx).Where(r.q.IotOtaFirmwareDO.ID.Eq(firmware.ID)).Updates(firmware)
	return err
}

func (r *OtaFirmwareRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.q.IotOtaFirmwareDO.WithContext(ctx).Where(r.q.IotOtaFirmwareDO.ID.Eq(id)).Delete()
	return err
}

func (r *OtaFirmwareRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.IotOtaFirmwareDO, error) {
	return r.q.IotOtaFirmwareDO.WithContext(ctx).Where(r.q.IotOtaFirmwareDO.ID.Eq(id)).First()
}

func (r *OtaFirmwareRepositoryImpl) GetPage(ctx context.Context, req *iot.IotOtaFirmwarePageReqVO) (*pagination.PageResult[*model.IotOtaFirmwareDO], error) {
	f := r.q.IotOtaFirmwareDO
	db := f.WithContext(ctx)
	if req.Name != "" {
		db = db.Where(f.Name.Like("%" + req.Name + "%"))
	}
	if req.ProductID != 0 {
		db = db.Where(f.ProductID.Eq(req.ProductID))
	}
	list, total, err := db.Order(f.ID.Desc()).FindByPage((req.PageNo-1)*req.PageSize, req.PageSize)
	return &pagination.PageResult[*model.IotOtaFirmwareDO]{List: list, Total: total}, err
}
