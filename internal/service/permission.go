package service

import (
	"context"

	"github.com/samber/lo"

	"backend-go/internal/model"
	"backend-go/internal/repo/query"
)

type PermissionService struct {
	q       *query.Query
	roleSvc *RoleService
}

func NewPermissionService(q *query.Query, roleSvc *RoleService) *PermissionService {
	return &PermissionService{
		q:       q,
		roleSvc: roleSvc,
	}
}

// GetUserRoleIdListByUserId 获取用户的角色ID列表
func (s *PermissionService) GetUserRoleIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	ur := s.q.SystemUserRole
	list, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *model.SystemUserRole, _ int) int64 {
		return item.RoleID
	}), nil
}

// GetRoleMenuListByRoleId 获取角色的菜单ID列表
func (s *PermissionService) GetRoleMenuListByRoleId(ctx context.Context, roleIds []int64) ([]int64, error) {
	if len(roleIds) == 0 {
		return []int64{}, nil
	}
	rm := s.q.SystemRoleMenu
	list, err := rm.WithContext(ctx).Where(rm.RoleID.In(roleIds...)).Find()
	if err != nil {
		return nil, err
	}
	// Extract MenuIDs, distinct
	return lo.Uniq(lo.Map(list, func(item *model.SystemRoleMenu, _ int) int64 {
		return item.MenuID
	})), nil
}

// AssignRoleMenu 赋予角色菜单
func (s *PermissionService) AssignRoleMenu(ctx context.Context, roleId int64, menuIds []int64) error {
	// Transaction
	return s.q.Transaction(func(tx *query.Query) error {
		// 1. Delete old
		rm := tx.SystemRoleMenu
		if _, err := rm.WithContext(ctx).Where(rm.RoleID.Eq(roleId)).Delete(); err != nil {
			return err
		}

		// 2. Insert new
		if len(menuIds) > 0 {
			var bat []*model.SystemRoleMenu
			for _, mid := range menuIds {
				bat = append(bat, &model.SystemRoleMenu{
					RoleID: roleId,
					MenuID: mid,
				})
			}
			// Use batch create
			if err := rm.WithContext(ctx).Create(bat...); err != nil {
				return err
			}
		}
		return nil
	})
}

// AssignRoleDataScope 赋予角色数据权限
func (s *PermissionService) AssignRoleDataScope(ctx context.Context, roleId int64, dataScope int, deptIds []int64) error {
	return s.roleSvc.UpdateRoleDataScope(ctx, roleId, dataScope, deptIds)
}

// AssignUserRole 赋予用户角色
func (s *PermissionService) AssignUserRole(ctx context.Context, userId int64, roleIds []int64) error {
	return s.q.Transaction(func(tx *query.Query) error {
		ur := tx.SystemUserRole
		// 1. Delete old
		if _, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userId)).Delete(); err != nil {
			return err
		}

		// 2. Insert new
		if len(roleIds) > 0 {
			var bat []*model.SystemUserRole
			for _, rid := range roleIds {
				bat = append(bat, &model.SystemUserRole{
					UserID: userId,
					RoleID: rid,
				})
			}
			if err := ur.WithContext(ctx).Create(bat...); err != nil {
				return err
			}
		}
		return nil
	})
}
