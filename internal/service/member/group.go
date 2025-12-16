package member

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type MemberGroupService struct {
	q *query.Query
}

func NewMemberGroupService(q *query.Query) *MemberGroupService {
	return &MemberGroupService{q: q}
}

// CreateGroup 创建用户分组
func (s *MemberGroupService) CreateGroup(ctx context.Context, r *req.MemberGroupCreateReq) (int64, error) {
	// TODO: 校验名字是否重复? (Optional)
	g := &member.MemberGroup{
		Name:   r.Name,
		Remark: r.Remark,
		Status: r.Status,
	}
	err := s.q.WithContext(ctx).MemberGroup.Create(g)
	if err != nil {
		return 0, err
	}
	return g.ID, nil
}

// UpdateGroup 更新用户分组
func (s *MemberGroupService) UpdateGroup(ctx context.Context, r *req.MemberGroupUpdateReq) error {
	_, err := s.q.WithContext(ctx).MemberGroup.Where(s.q.MemberGroup.ID.Eq(r.ID)).Updates(&member.MemberGroup{
		Name:   r.Name,
		Remark: r.Remark,
		Status: r.Status,
	})
	return err
}

// DeleteGroup 删除用户分组
func (s *MemberGroupService) DeleteGroup(ctx context.Context, id int64) error {
	_, err := s.q.WithContext(ctx).MemberGroup.Where(s.q.MemberGroup.ID.Eq(id)).Delete()
	return err
}

// GetGroup 获得用户分组
func (s *MemberGroupService) GetGroup(ctx context.Context, id int64) (*member.MemberGroup, error) {
	return s.q.WithContext(ctx).MemberGroup.Where(s.q.MemberGroup.ID.Eq(id)).First()
}

// GetGroupPage 获得用户分组分页
func (s *MemberGroupService) GetGroupPage(ctx context.Context, r *req.MemberGroupPageReq) (*pagination.PageResult[*member.MemberGroup], error) {
	q := s.q.MemberGroup.WithContext(ctx)
	if r.Name != "" {
		q = q.Where(s.q.MemberGroup.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.MemberGroup.Status.Eq(*r.Status))
	}
	q = q.Order(s.q.MemberGroup.ID.Desc())

	list, count, err := q.FindByPage(r.GetOffset(), r.PageSize)
	if err != nil {
		return nil, err
	}
	return &pagination.PageResult[*member.MemberGroup]{
		List:  list,
		Total: count,
	}, nil
}

// GetEnableGroupList 获得开启的用户分组列表
func (s *MemberGroupService) GetEnableGroupList(ctx context.Context) ([]*member.MemberGroup, error) {
	return s.q.WithContext(ctx).MemberGroup.Where(s.q.MemberGroup.Status.Eq(0)).Find() // 0 = Enable
}

// GetGroupListByIds 根据 ID 列表获得分组列表
func (s *MemberGroupService) GetGroupListByIds(ctx context.Context, ids []int64) ([]*member.MemberGroup, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.q.MemberGroup.WithContext(ctx).Where(s.q.MemberGroup.ID.In(ids...)).Find()
}
