package system

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

type DeptService struct {
	q *query.Query
}

func NewDeptService(q *query.Query) *DeptService {
	return &DeptService{
		q: q,
	}
}

func (s *DeptService) CreateDept(ctx context.Context, req *system.DeptSaveReq) (int64, error) {
	d := s.q.SystemDept
	if req.ParentID > 0 {
		_, err := d.WithContext(ctx).Where(d.ID.Eq(req.ParentID)).First()
		if err != nil {
			return 0, errors.New("父部门不存在")
		}
	}

	dept := &model.SystemDept{
		Name:         req.Name,
		ParentID:     req.ParentID,
		Sort:         req.Sort,
		LeaderUserID: req.LeaderUserID,
		Phone:        req.Phone,
		Email:        req.Email,
		Status:       int32(req.Status),
	}
	err := d.WithContext(ctx).Create(dept)
	return dept.ID, err
}

func (s *DeptService) UpdateDept(ctx context.Context, req *system.DeptSaveReq) error {
	d := s.q.SystemDept
	if req.ID == req.ParentID {
		return errors.New("父部门不能是自己")
	}
	_, err := d.WithContext(ctx).Where(d.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("部门不存在")
	}

	if req.ParentID > 0 {
		_, err := d.WithContext(ctx).Where(d.ID.Eq(req.ParentID)).First()
		if err != nil {
			return errors.New("父部门不存在")
		}
	}

	_, err = d.WithContext(ctx).Where(d.ID.Eq(req.ID)).Updates(&model.SystemDept{
		Name:         req.Name,
		ParentID:     req.ParentID,
		Sort:         req.Sort,
		LeaderUserID: req.LeaderUserID,
		Phone:        req.Phone,
		Email:        req.Email,
		Status:       int32(req.Status),
	})
	return err
}

func (s *DeptService) DeleteDept(ctx context.Context, id int64) error {
	d := s.q.SystemDept
	// Check children
	count, err := d.WithContext(ctx).Where(d.ParentID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("存在子部门，无法删除")
	}
	// Check users assigned to this dept
	userCount, err := s.q.SystemUser.WithContext(ctx).Where(s.q.SystemUser.DeptID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if userCount > 0 {
		return errors.New("部门下存在用户，无法删除")
	}

	_, err = d.WithContext(ctx).Where(d.ID.Eq(id)).Delete()
	return err
}

func (s *DeptService) GetDept(ctx context.Context, id int64) (*system.DeptRespVO, error) {
	d := s.q.SystemDept
	item, err := d.WithContext(ctx).Where(d.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return &system.DeptRespVO{
		ID:           item.ID,
		Name:         item.Name,
		ParentID:     item.ParentID,
		Sort:         item.Sort,
		LeaderUserID: item.LeaderUserID,
		Phone:        item.Phone,
		Email:        item.Email,
		Status:       item.Status,
		CreateTime:   item.CreateTime,
	}, nil
}

func (s *DeptService) GetDeptList(ctx context.Context, req *system.DeptListReq) ([]*system.DeptRespVO, error) {
	d := s.q.SystemDept
	qb := d.WithContext(ctx)

	if req.Name != "" {
		qb = qb.Where(d.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(d.Status.Eq(int32(*req.Status)))
	}

	list, err := qb.Order(d.Sort, d.ID).Find()
	if err != nil {
		return nil, err
	}

	var res []*system.DeptRespVO
	for _, item := range list {
		res = append(res, &system.DeptRespVO{
			ID:           item.ID,
			Name:         item.Name,
			ParentID:     item.ParentID,
			Sort:         item.Sort,
			LeaderUserID: item.LeaderUserID,
			Phone:        item.Phone,
			Email:        item.Email,
			Status:       item.Status,
			CreateTime:   item.CreateTime,
		})
	}
	return res, nil
}

func (s *DeptService) GetSimpleDeptList(ctx context.Context) ([]*system.DeptSimpleRespVO, error) {
	d := s.q.SystemDept
	list, err := d.WithContext(ctx).Where(d.Status.Eq(0)).Order(d.Sort, d.ID).Find()
	if err != nil {
		return nil, err
	}

	var res []*system.DeptSimpleRespVO
	for _, item := range list {
		res = append(res, &system.DeptSimpleRespVO{
			ID:       item.ID,
			Name:     item.Name,
			ParentID: item.ParentID,
		})
	}
	return res, nil
}

// GetDeptIdListByParentId 递归获取某个部门的所有子部门ID列表（用于数据权限）
func (s *DeptService) GetDeptIdListByParentId(ctx context.Context, parentId int64) ([]int64, error) {
	d := s.q.SystemDept

	// 获取直接子部门
	children, err := d.WithContext(ctx).Where(d.ParentID.Eq(parentId)).Find()
	if err != nil {
		return nil, err
	}

	var result []int64
	for _, child := range children {
		// 添加当前子部门ID
		result = append(result, child.ID)

		// 递归获取子部门的子部门
		grandchildren, err := s.GetDeptIdListByParentId(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, grandchildren...)
	}

	return result, nil
}
