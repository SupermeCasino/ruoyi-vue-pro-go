package service

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type TenantPackageService struct {
	q *query.Query
}

func NewTenantPackageService(q *query.Query) *TenantPackageService {
	return &TenantPackageService{q: q}
}

// CreateTenantPackage 创建租户套餐
func (s *TenantPackageService) CreateTenantPackage(ctx context.Context, r *req.TenantPackageSaveReq) (int64, error) {
	pkg := &model.SystemTenantPackage{
		Name:    r.Name,
		Status:  int32(r.Status),
		Remark:  r.Remark,
		MenuIDs: r.MenuIds,
	}
	if err := s.q.SystemTenantPackage.WithContext(ctx).Create(pkg); err != nil {
		return 0, err
	}
	return pkg.ID, nil
}

// UpdateTenantPackage 更新租户套餐
func (s *TenantPackageService) UpdateTenantPackage(ctx context.Context, r *req.TenantPackageSaveReq) error {
	t := s.q.SystemTenantPackage
	// 1. 校验存在
	if _, err := t.WithContext(ctx).Where(t.ID.Eq(r.ID)).First(); err != nil {
		return err
	}
	// TODO: 校验状态，如果从开启变为关闭，需要做一些处理？Java 版通常有校验
	_, err := t.WithContext(ctx).Where(t.ID.Eq(r.ID)).Updates(&model.SystemTenantPackage{
		Name:    r.Name,
		Status:  int32(r.Status),
		Remark:  r.Remark,
		MenuIDs: r.MenuIds,
	})
	return err
}

// DeleteTenantPackage 删除租户套餐
func (s *TenantPackageService) DeleteTenantPackage(ctx context.Context, id int64) error {
	// TODO: 校验是否有租户正在使用该套餐
	t := s.q.SystemTenantPackage
	_, err := t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

// DeleteTenantPackageList 批量删除租户套餐
func (s *TenantPackageService) DeleteTenantPackageList(ctx context.Context, ids []int64) error {
	// TODO: 校验是否有租户正在使用这些套餐
	t := s.q.SystemTenantPackage
	_, err := t.WithContext(ctx).Where(t.ID.In(ids...)).Delete()
	return err
}

// GetTenantPackage 获得租户套餐
func (s *TenantPackageService) GetTenantPackage(ctx context.Context, id int64) (*resp.TenantPackageResp, error) {
	t := s.q.SystemTenantPackage
	pkg, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return &resp.TenantPackageResp{
		ID:         pkg.ID,
		Name:       pkg.Name,
		Status:     int(pkg.Status),
		Remark:     pkg.Remark,
		MenuIds:    pkg.MenuIDs,
		CreateTime: pkg.CreateTime,
	}, nil
}

// GetTenantPackageSimpleList 获取租户套餐精简列表
func (s *TenantPackageService) GetTenantPackageSimpleList(ctx context.Context) ([]*resp.TenantPackageResp, error) {
	t := s.q.SystemTenantPackage
	list, err := t.WithContext(ctx).Where(t.Status.Eq(0)).Find() // 0 = 开启
	if err != nil {
		return nil, err
	}
	respList := make([]*resp.TenantPackageResp, len(list))
	for i, pkg := range list {
		respList[i] = &resp.TenantPackageResp{
			ID:   pkg.ID,
			Name: pkg.Name,
		}
	}
	return respList, nil
}

// GetTenantPackagePage 获得租户套餐分页
func (s *TenantPackageService) GetTenantPackagePage(ctx context.Context, r *req.TenantPackagePageReq) (*pagination.PageResult[*resp.TenantPackageResp], error) {
	t := s.q.SystemTenantPackage
	q := t.WithContext(ctx)

	// 条件过滤
	if r.Name != "" {
		q = q.Where(t.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		q = q.Where(t.Status.Eq(int32(*r.Status)))
	}
	if r.Remark != "" {
		q = q.Where(t.Remark.Like("%" + r.Remark + "%"))
	}
	if len(r.CreateTime) == 2 {
		startTime, _ := time.Parse("2006-01-02 15:04:05", r.CreateTime[0])
		endTime, _ := time.Parse("2006-01-02 15:04:05", r.CreateTime[1])
		q = q.Where(t.CreateTime.Between(startTime, endTime))
	}

	// 分页查询
	offset := (r.PageNo - 1) * r.PageSize
	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	list, err := q.Order(t.ID.Desc()).Offset(offset).Limit(r.PageSize).Find()
	if err != nil {
		return nil, err
	}

	// Entity → DTO 转换
	respList := make([]*resp.TenantPackageResp, len(list))
	for i, pkg := range list {
		respList[i] = &resp.TenantPackageResp{
			ID:         pkg.ID,
			Name:       pkg.Name,
			Status:     int(pkg.Status),
			Remark:     pkg.Remark,
			MenuIds:    pkg.MenuIDs,
			CreateTime: pkg.CreateTime,
		}
	}

	return &pagination.PageResult[*resp.TenantPackageResp]{
		List:  respList,
		Total: count,
	}, nil
}
