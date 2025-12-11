package service

import (
	"context"
	"strings"

	"backend-go/internal/api/resp"
	"backend-go/internal/repo/query"
)

type TenantService struct {
	q *query.Query
}

func NewTenantService(q *query.Query) *TenantService {
	return &TenantService{q: q}
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
	// 注意：数据库中可能没有 websites 列，需要优雅降级
	// 如果数据库不支持此查询，直接返回 nil（表示未找到）
	tenantRepo := s.q.SystemTenant
	list, err := tenantRepo.WithContext(ctx).Find()
	if err != nil {
		return nil, nil // 查询失败，返回 nil，不报错
	}

	// 在应用层进行过滤（兼容没有 websites 列的情况）
	for _, t := range list {
		if t.Status == 0 && strings.Contains(t.Websites, website) {
			return &resp.TenantSimpleResp{
				ID:   t.ID,
				Name: t.Name,
			}, nil
		}
	}
	return nil, nil // 未找到返回 nil，不报错
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
