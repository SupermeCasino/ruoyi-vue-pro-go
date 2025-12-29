package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	pkgContext "github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"
)

type TenantService struct {
	q             *query.Query
	roleSvc       *RoleService
	permissionSvc *PermissionService
}

func NewTenantService(q *query.Query, roleSvc *RoleService, permissionSvc *PermissionService) *TenantService {
	return &TenantService{
		q:             q,
		roleSvc:       roleSvc,
		permissionSvc: permissionSvc,
	}
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, req *req.TenantCreateReq) (int64, error) {
	// 1. 校验租户名是否重复
	if err := s.checkNameUnique(ctx, req.Name, 0); err != nil {
		return 0, err
	}

	// 2. 事务执行
	var tenantId int64
	err := s.q.Transaction(func(tx *query.Query) error {
		// 2.1 创建租户
		tenant := &model.SystemTenant{
			Name:          req.Name,
			ContactName:   req.ContactName,
			ContactMobile: req.ContactMobile,
			Status:        int32(req.Status),
			PackageID:     req.PackageID,
			AccountCount:  int32(req.AccountCount),
			ExpireDate:    time.UnixMilli(req.ExpireTime),
			Websites:      req.Website,
		}
		if err := tx.SystemTenant.WithContext(ctx).Create(tenant); err != nil {
			return err
		}
		tenantId = tenant.ID

		// 角色
		role := &model.SystemRole{
			Name:   "租户管理员",
			Code:   "tenant_admin",
			Sort:   0,
			Status: 0, // 启用
			Type:   2, // 自定义角色
			Remark: "系统自动生成",
		}
		role.TenantID = tenantId
		if err := tx.SystemRole.WithContext(ctx).Create(role); err != nil {
			return err
		}

		// 用户
		hashedPwd, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		user := &model.SystemUser{
			Username: req.Username,
			Password: hashedPwd,
			Nickname: req.ContactName,
			Mobile:   req.ContactMobile,
			Status:   0, // 启用
		}
		user.TenantID = tenantId
		if err := tx.SystemUser.WithContext(ctx).Create(user); err != nil {
			return err
		}

		// 赋予套餐中的菜单权限
		// 获取套餐信息
		var pkg model.SystemTenantPackage
		if err := tx.SystemTenant.WithContext(ctx).UnderlyingDB().Model(&model.SystemTenantPackage{}).Where("id = ?", req.PackageID).First(&pkg).Error; err != nil {
			return err
		}
		menuIds := pkg.MenuIDs

		if len(menuIds) > 0 {
			roleMenus := make([]*model.SystemRoleMenu, len(menuIds))
			for i, mid := range menuIds {
				roleMenus[i] = &model.SystemRoleMenu{
					RoleID: role.ID,
					MenuID: mid,
				}
				roleMenus[i].TenantID = tenantId
			}
			if err := tx.SystemRoleMenu.WithContext(ctx).Create(roleMenus...); err != nil {
				return err
			}
		}

		// 2.3 更新租户的联系人用户ID (ContactUserID)
		if _, err := tx.SystemTenant.WithContext(ctx).Where(tx.SystemTenant.ID.Eq(tenantId)).Update(tx.SystemTenant.ContactUserID, user.ID); err != nil {
			return err
		}

		return nil
	})

	return tenantId, err
}

// UpdateTenant 更新租户
func (s *TenantService) UpdateTenant(ctx context.Context, req *req.TenantUpdateReq) error {
	// 1. 校验存在
	t := s.q.SystemTenant
	tenant, err := t.WithContext(ctx).Where(t.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("租户不存在")
	}

	// 2. 校验名字唯一
	if err := s.checkNameUnique(ctx, req.Name, req.ID); err != nil {
		return err
	}

	// 3. 更新
	tenantObj := &model.SystemTenant{
		Name:          req.Name,
		ContactName:   req.ContactName,
		ContactMobile: req.ContactMobile,
		Status:        int32(req.Status),
		PackageID:     req.PackageID,
		AccountCount:  int32(req.AccountCount),
		ExpireDate:    time.UnixMilli(req.ExpireTime),
		Websites:      req.Website,
	}
	_, err = t.WithContext(ctx).Where(t.ID.Eq(req.ID)).Updates(tenantObj)
	if err != nil {
		return err
	}

	// 4. 如果套餐发生变化，则修改其角色的权限
	if tenant.PackageID != req.PackageID {
		// 获得套餐菜单
		var pkg model.SystemTenantPackage
		if err := s.q.SystemTenant.WithContext(ctx).UnderlyingDB().Model(&model.SystemTenantPackage{}).Where("id = ?", req.PackageID).First(&pkg).Error; err != nil {
			return err
		}
		menuIds := pkg.MenuIDs
		if err := s.updateTenantRoleMenu(ctx, req.ID, menuIds); err != nil {
			return err
		}
	}

	return nil
}

// DeleteTenant 删除租户
func (s *TenantService) DeleteTenant(ctx context.Context, id int64) error {
	t := s.q.SystemTenant
	_, err := t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

// GetTenant 获得租户
func (s *TenantService) GetTenant(ctx context.Context, id int64) (*resp.TenantRespVO, error) {
	t := s.q.SystemTenant
	tenant, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return &resp.TenantRespVO{
		ID:            tenant.ID,
		Name:          tenant.Name,
		ContactUserID: tenant.ContactUserID,
		ContactName:   tenant.ContactName,
		ContactMobile: tenant.ContactMobile,
		Status:        int(tenant.Status),
		Website:       tenant.Websites,
		PackageID:     tenant.PackageID,
		AccountCount:  int(tenant.AccountCount),
		ExpireTime:    tenant.ExpireDate.UnixMilli(),
		CreateTime:    tenant.CreateTime,
	}, nil
}

// GetTenantPage 获得租户分页
func (s *TenantService) GetTenantPage(ctx context.Context, req *req.TenantPageReq) (*pagination.PageResult[*resp.TenantRespVO], error) {
	t := s.q.SystemTenant
	qb := t.WithContext(ctx)

	if req.Name != "" {
		qb = qb.Where(t.Name.Like("%" + req.Name + "%"))
	}
	if req.ContactName != "" {
		qb = qb.Where(t.ContactName.Like("%" + req.ContactName + "%"))
	}
	if req.ContactMobile != "" {
		qb = qb.Where(t.ContactMobile.Like("%" + req.ContactMobile + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(t.Status.Eq(int32(*req.Status)))
	}
	if req.CreateTimeGe != nil {
		qb = qb.Where(t.CreateTime.Gte(*req.CreateTimeGe))
	}
	if req.CreateTimeLe != nil {
		qb = qb.Where(t.CreateTime.Lte(*req.CreateTimeLe))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(t.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	var data []*resp.TenantRespVO
	for _, item := range list {
		data = append(data, &resp.TenantRespVO{
			ID:            item.ID,
			Name:          item.Name,
			ContactUserID: item.ContactUserID,
			ContactName:   item.ContactName,
			ContactMobile: item.ContactMobile,
			Status:        int(item.Status),
			Website:       item.Websites,
			PackageID:     item.PackageID,
			AccountCount:  int(item.AccountCount),
			ExpireTime:    item.ExpireDate.UnixMilli(),
			CreateTime:    item.CreateTime,
		})
	}

	return &pagination.PageResult[*resp.TenantRespVO]{
		List:  data,
		Total: total,
	}, nil
}

// GetTenantList 获得租户列表 (用于导出)
func (s *TenantService) GetTenantList(ctx context.Context, req *req.TenantExportReq) ([]*resp.TenantRespVO, error) {
	t := s.q.SystemTenant
	qb := t.WithContext(ctx)

	if req.Name != "" {
		qb = qb.Where(t.Name.Like("%" + req.Name + "%"))
	}
	if req.ContactName != "" {
		qb = qb.Where(t.ContactName.Like("%" + req.ContactName + "%"))
	}
	if req.ContactMobile != "" {
		qb = qb.Where(t.ContactMobile.Like("%" + req.ContactMobile + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(t.Status.Eq(int32(*req.Status)))
	}
	if req.CreateTimeGe != nil {
		qb = qb.Where(t.CreateTime.Gte(*req.CreateTimeGe))
	}
	if req.CreateTimeLe != nil {
		qb = qb.Where(t.CreateTime.Lte(*req.CreateTimeLe))
	}

	list, err := qb.Order(t.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	var data []*resp.TenantRespVO
	for _, item := range list {
		data = append(data, &resp.TenantRespVO{
			ID:            item.ID,
			Name:          item.Name,
			ContactUserID: item.ContactUserID,
			ContactName:   item.ContactName,
			ContactMobile: item.ContactMobile,
			Status:        int(item.Status),
			Website:       item.Websites,
			PackageID:     item.PackageID,
			ExpireTime:    item.ExpireDate.UnixMilli(), // 毫秒级时间戳
			AccountCount:  int(item.AccountCount),
			CreateTime:    item.CreateTime,
		})
	}
	return data, nil
}

func (s *TenantService) checkNameUnique(ctx context.Context, name string, excludeId int64) error {
	t := s.q.SystemTenant
	qb := t.WithContext(ctx).Where(t.Name.Eq(name))
	if excludeId > 0 {
		qb = qb.Where(t.ID.Neq(excludeId))
	}
	count, err := qb.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("租户名已存在")
	}
	return nil
}

// GetTenantSimpleList 获取启用状态的租户精简列表
func (s *TenantService) GetTenantSimpleList(ctx context.Context) ([]resp.TenantSimpleResp, error) {
	tenantRepo := s.q.SystemTenant
	list, err := tenantRepo.WithContext(ctx).Where(tenantRepo.Status.Eq(0)).Find() // 0 = 启用
	if err != nil {
		return nil, err
	}

	result := make([]resp.TenantSimpleResp, 0, len(list))
	for _, t := range list {
		result = append(result, resp.TenantSimpleResp{
			ID:   t.ID,
			Name: t.Name,
		})
	}
	return result, nil
}

// GetTenantByWebsite 根据域名查询租户
func (s *TenantService) GetTenantByWebsite(ctx context.Context, website string) (*resp.TenantSimpleResp, error) {
	tenantRepo := s.q.SystemTenant
	list, err := tenantRepo.WithContext(ctx).Find()
	if err != nil {
		return nil, nil
	}

	// 应用层过滤
	for _, t := range list {
		if t.Status == 0 && strings.Contains(t.Websites, website) {
			return &resp.TenantSimpleResp{
				ID:   t.ID,
				Name: t.Name,
			}, nil
		}
	}
	return nil, nil // 未找到返回空
}

// GetTenantIdByName 根据租户名获取租户ID
func (s *TenantService) GetTenantIdByName(ctx context.Context, name string) (*int64, error) {
	tenantRepo := s.q.SystemTenant
	tenant, err := tenantRepo.WithContext(ctx).Where(tenantRepo.Name.Eq(name)).First()
	if err != nil {
		return nil, nil // 未找到返回 nil
	}
	return &tenant.ID, nil
}

// GetTenantByName 根据租户名获取租户（供 AuthService 使用）
func (s *TenantService) GetTenantByName(ctx context.Context, name string) (*resp.TenantSimpleResp, error) {
	tenantRepo := s.q.SystemTenant
	tenant, err := tenantRepo.WithContext(ctx).Where(tenantRepo.Name.Eq(name)).First()
	if err != nil {
		return nil, err
	}
	return &resp.TenantSimpleResp{
		ID:   tenant.ID,
		Name: tenant.Name,
	}, nil
}

// HandleTenantMenu 处理租户菜单过滤
// handler 接收租户允许的菜单ID列表，并在回调中移除不在列表中的菜单
// 如果是系统租户，传入 nil 表示允许所有菜单
func (s *TenantService) HandleTenantMenu(ctx context.Context, handler func(allowedMenuIds []int64)) error {
	// 1. 获得租户ID
	var tenantId int64
	if c, ok := ctx.(*gin.Context); ok {
		tenantId = pkgContext.GetTenantId(c)
	}
	if tenantId == 0 {
		return nil // 如果没有租户上下文（如admin），不做过滤
	}

	// 2. 获得租户
	tenantRepo := s.q.SystemTenant
	tenant, err := tenantRepo.WithContext(ctx).Where(tenantRepo.ID.Eq(tenantId)).First()
	if err != nil {
		return err
	}

	// 3. 如果是系统租户 (PackageID=0), 允许所有菜单
	if tenant.PackageID == 0 {
		handler(nil) // nil 表示允许所有
		return nil
	}

	// 4. 读取租户套餐
	var pkg model.SystemTenantPackage
	if err := s.q.SystemTenant.WithContext(ctx).UnderlyingDB().Model(&model.SystemTenantPackage{}).Where("id = ?", tenant.PackageID).First(&pkg).Error; err != nil {
		// 如果套餐不存在，当作无权限处理
		handler([]int64{})
		return nil
	}

	allowedMenuIds := pkg.MenuIDs

	// 6. 执行处理
	handler(allowedMenuIds)
	return nil
}

// updateTenantRoleMenu 更新租户下所有角色的菜单权限
func (s *TenantService) updateTenantRoleMenu(ctx context.Context, tenantId int64, menuIds []int64) error {
	// 通过 Q 的 Where 条件保证租户隔离
	r := s.q.SystemRole
	roles, err := r.WithContext(ctx).Where(r.TenantID.Eq(tenantId)).Find()
	if err != nil {
		return err
	}

	for _, role := range roles {
		if role.Code == "tenant_admin" {
			// 超级管理员：直接分配套餐的所有功能
			if err := s.permissionSvc.AssignRoleMenu(ctx, role.ID, menuIds); err != nil {
				return err
			}
		} else {
			// 普通用户：裁剪掉超出套餐功能的权限 (交集)
			// 1. 获取角色当前拥有的菜单
			roleMenuIds, err := s.permissionSvc.GetRoleMenuListByRoleId(ctx, []int64{role.ID})
			if err != nil {
				return err
			}
			// 2. 取交集
			newRoleMenuIds := utils.Intersect(roleMenuIds, menuIds)
			// 3. 重新分配
			if err := s.permissionSvc.AssignRoleMenu(ctx, role.ID, newRoleMenuIds); err != nil {
				return err
			}
		}
	}
	return nil
}
