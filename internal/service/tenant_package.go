package service

import (
	"context"
	"errors"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
)

type TenantPackageService struct {
	q         *query.Query
	tenantSvc *TenantService
}

func NewTenantPackageService(q *query.Query, tenantSvc *TenantService) *TenantPackageService {
	return &TenantPackageService{
		q:         q,
		tenantSvc: tenantSvc,
	}
}

// CreateTenantPackage 创建租户套餐
func (s *TenantPackageService) CreateTenantPackage(ctx context.Context, r *req.TenantPackageSaveReq) (int64, error) {
	// 校验名称唯一
	if err := s.validateTenantPackageNameUnique(ctx, r.Name, 0); err != nil {
		return 0, err
	}

	pkg := &model.SystemTenantPackage{
		Name:    r.Name,
		Status:  int32(*r.Status),
		Remark:  r.Remark,
		MenuIDs: r.MenuIDs,
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
	oldPkg, err := t.WithContext(ctx).Where(t.ID.Eq(r.ID)).First()
	if err != nil {
		return errors.New("租户套餐不存在")
	}
	// 2. 校验名称唯一
	if err := s.validateTenantPackageNameUnique(ctx, r.Name, r.ID); err != nil {
		return err
	}

	// 3. 更新
	_, err = t.WithContext(ctx).Where(t.ID.Eq(r.ID)).Updates(&model.SystemTenantPackage{
		Name:    r.Name,
		Status:  int32(*r.Status),
		Remark:  r.Remark,
		MenuIDs: r.MenuIDs,
	})
	if err != nil {
		return err
	}

	// 4. 如果菜单发生变化，则修改每个租户的权限 - 对齐 Java: updateTenantPackage 数据同步逻辑
	if !utils.IsEqualList(oldPkg.MenuIDs, r.MenuIDs) {
		// 查找所有使用该套餐的租户
		tenants, err := s.q.SystemTenant.WithContext(ctx).Where(s.q.SystemTenant.PackageID.Eq(r.ID)).Find()
		if err == nil {
			for _, tenant := range tenants {
				_ = s.tenantSvc.updateTenantRoleMenu(ctx, tenant.ID, r.MenuIDs)
			}
		}
	}
	return nil
}

// DeleteTenantPackage 删除租户套餐
func (s *TenantPackageService) DeleteTenantPackage(ctx context.Context, id int64) error {
	// 校验是否有租户正在使用该套餐
	if err := s.validateTenantUsed(ctx, id); err != nil {
		return err
	}

	t := s.q.SystemTenantPackage
	_, err := t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

// DeleteTenantPackageList 批量删除租户套餐
func (s *TenantPackageService) DeleteTenantPackageList(ctx context.Context, ids []int64) error {
	// 校验是否有租户正在使用这些套餐
	for _, id := range ids {
		if err := s.validateTenantUsed(ctx, id); err != nil {
			return err
		}
	}
	t := s.q.SystemTenantPackage
	_, err := t.WithContext(ctx).Where(t.ID.In(ids...)).Delete()
	return err
}

// validateTenantPackageNameUnique 校验套餐名是否唯一 (对齐 Java: validateTenantPackageNameUnique)
func (s *TenantPackageService) validateTenantPackageNameUnique(ctx context.Context, name string, id int64) error {
	t := s.q.SystemTenantPackage
	pkg, err := t.WithContext(ctx).Where(t.Name.Eq(name)).First()
	if err == nil && pkg.ID != id {
		return errors.New("套餐名已存在")
	}
	return nil
}

// validateTenantUsed 校验套餐是否正在被使用 (对齐 Java: validateTenantUsed)
func (s *TenantPackageService) validateTenantUsed(ctx context.Context, id int64) error {
	count, err := s.q.SystemTenant.WithContext(ctx).Where(s.q.SystemTenant.PackageID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该套餐正在被使用，无法删除")
	}
	return nil
}

// ValidTenantPackage 校验套餐有效性 (存在且启用) - 对齐 Java: validTenantPackage
func (s *TenantPackageService) ValidTenantPackage(ctx context.Context, id int64) (*model.SystemTenantPackage, error) {
	t := s.q.SystemTenantPackage
	pkg, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.New("租户套餐不存在")
	}
	if pkg.Status != 0 {
		return nil, errors.New("租户套餐已禁用: " + pkg.Name)
	}
	return pkg, nil
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
		MenuIDs:    pkg.MenuIDs,
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
			MenuIDs:    pkg.MenuIDs,
			CreateTime: pkg.CreateTime,
		}
	}

	return &pagination.PageResult[*resp.TenantPackageResp]{
		List:  respList,
		Total: count,
	}, nil
}
