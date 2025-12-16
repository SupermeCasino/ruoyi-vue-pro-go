package member

import (
	"context"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"

	"gorm.io/gorm"
)

type MemberTagService struct {
	q                 *query.Query
	memberUserService *MemberUserService
}

func NewMemberTagService(q *query.Query, memberUserService *MemberUserService) *MemberTagService {
	return &MemberTagService{
		q:                 q,
		memberUserService: memberUserService,
	}
}

// CreateTag 创建用户标签
func (s *MemberTagService) CreateTag(ctx context.Context, r *req.MemberTagCreateReq) (int64, error) {
	// 校验名称唯一
	if err := s.validateNameUnique(ctx, 0, r.Name); err != nil {
		return 0, err
	}
	tag := &member.MemberTag{
		Name:   r.Name,
		Remark: r.Remark,
	}
	err := s.q.MemberTag.WithContext(ctx).Create(tag)
	if err != nil {
		return 0, err
	}
	return tag.ID, nil
}

// UpdateTag 更新用户标签
func (s *MemberTagService) UpdateTag(ctx context.Context, r *req.MemberTagUpdateReq) error {
	// 校验存在
	_, err := s.q.MemberTag.WithContext(ctx).Where(s.q.MemberTag.ID.Eq(r.ID)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewBizError(1004014002, "标签不存在")
		}
		return err
	}
	// 校验名称唯一
	if err := s.validateNameUnique(ctx, r.ID, r.Name); err != nil {
		return err
	}

	// Update
	_, err = s.q.MemberTag.WithContext(ctx).Where(s.q.MemberTag.ID.Eq(r.ID)).Updates(&member.MemberTag{
		Name:   r.Name,
		Remark: r.Remark,
	})
	return err
}

// DeleteTag 删除用户标签
func (s *MemberTagService) DeleteTag(ctx context.Context, id int64) error {
	// 校验存在
	_, err := s.q.MemberTag.WithContext(ctx).Where(s.q.MemberTag.ID.Eq(id)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return core.NewBizError(1004014002, "标签不存在")
		}
		return err
	}

	// 校验标签下是否有用户
	count, err := s.memberUserService.GetUserCountByTagId(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return core.NewBizError(1004014003, "标签下存在用户，无法删除")
	}

	_, err = s.q.MemberTag.WithContext(ctx).Where(s.q.MemberTag.ID.Eq(id)).Delete()
	return err
}

// GetTag 获得用户标签
func (s *MemberTagService) GetTag(ctx context.Context, id int64) (*member.MemberTag, error) {
	return s.q.MemberTag.WithContext(ctx).Where(s.q.MemberTag.ID.Eq(id)).First()
}

// GetTagPage 获得用户标签分页
func (s *MemberTagService) GetTagPage(ctx context.Context, r *req.MemberTagPageReq) (*core.PageResult[*member.MemberTag], error) {
	q := s.q.MemberTag.WithContext(ctx)
	if r.Name != nil && *r.Name != "" {
		q = q.Where(s.q.MemberTag.Name.Like("%" + *r.Name + "%"))
	}
	q = q.Order(s.q.MemberTag.ID.Desc())

	list, count, err := q.FindByPage(r.GetOffset(), r.PageSize)
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*member.MemberTag]{
		List:  list,
		Total: count,
	}, nil
}

// GetTagList 获得开启的用户标签列表
func (s *MemberTagService) GetTagList(ctx context.Context) ([]*member.MemberTag, error) {
	return s.q.MemberTag.WithContext(ctx).Find()
}

func (s *MemberTagService) validateNameUnique(ctx context.Context, id int64, name string) error {
	q := s.q.MemberTag.WithContext(ctx).Where(s.q.MemberTag.Name.Eq(name))
	if id > 0 {
		q = q.Where(s.q.MemberTag.ID.Neq(id))
	}
	count, err := q.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return core.NewBizError(1004014000, "标签名称已存在")
	}
	return nil
}

// GetTagListByIds 根据 ID 列表获得标签列表
func (s *MemberTagService) GetTagListByIds(ctx context.Context, ids []int64) ([]*member.MemberTag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.q.MemberTag.WithContext(ctx).Where(s.q.MemberTag.ID.In(ids...)).Find()
}
