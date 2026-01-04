package system

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PostService struct {
	q *query.Query
}

func NewPostService(q *query.Query) *PostService {
	return &PostService{
		q: q,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *system.PostSaveReq) (int64, error) {
	p := s.q.SystemPost
	// Check Code/Name unique?
	count, err := p.WithContext(ctx).Where(p.Name.Eq(req.Name)).Or(p.Code.Eq(req.Code)).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("岗位名称或编码已存在")
	}

	post := &model.SystemPost{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: int32(*req.Status),
		Remark: req.Remark,
	}
	err = p.WithContext(ctx).Create(post)
	return post.ID, err
}

func (s *PostService) UpdatePost(ctx context.Context, req *system.PostSaveReq) error {
	p := s.q.SystemPost
	_, err := p.WithContext(ctx).Where(p.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("岗位不存在")
	}
	// Check Code/Name unique excluding self
	count, err := p.WithContext(ctx).Where(p.ID.Neq(req.ID)).Where(p.Name.Eq(req.Name)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("岗位名称已存在")
	}

	count, err = p.WithContext(ctx).Where(p.ID.Neq(req.ID)).Where(p.Code.Eq(req.Code)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("岗位编码已存在")
	}

	_, err = p.WithContext(ctx).Where(p.ID.Eq(req.ID)).Updates(&model.SystemPost{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: int32(*req.Status),
		Remark: req.Remark,
	})
	return err
}

func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	p := s.q.SystemPost
	_, err := p.WithContext(ctx).Where(p.ID.Eq(id)).Delete()
	return err
}

func (s *PostService) GetPost(ctx context.Context, id int64) (*system.PostRespVO, error) {
	p := s.q.SystemPost
	item, err := p.WithContext(ctx).Where(p.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return &system.PostRespVO{
		ID:         item.ID,
		Name:       item.Name,
		Code:       item.Code,
		Sort:       item.Sort,
		Status:     item.Status,
		Remark:     item.Remark,
		CreateTime: item.CreateTime,
	}, nil
}

func (s *PostService) GetPostPage(ctx context.Context, req *system.PostPageReq) (*pagination.PageResult[*system.PostRespVO], error) {
	p := s.q.SystemPost
	qb := p.WithContext(ctx)

	if req.Code != "" {
		qb = qb.Where(p.Code.Like("%" + req.Code + "%"))
	}
	if req.Name != "" {
		qb = qb.Where(p.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(p.Status.Eq(int32(*req.Status)))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(p.Sort, p.ID).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	var data []*system.PostRespVO
	for _, item := range list {
		data = append(data, &system.PostRespVO{
			ID:         item.ID,
			Name:       item.Name,
			Code:       item.Code,
			Sort:       item.Sort,
			Status:     item.Status,
			Remark:     item.Remark,
			CreateTime: item.CreateTime,
		})
	}

	return &pagination.PageResult[*system.PostRespVO]{
		List:  data,
		Total: total,
	}, nil
}

func (s *PostService) GetSimplePostList(ctx context.Context) ([]*system.PostSimpleRespVO, error) {
	p := s.q.SystemPost
	list, err := p.WithContext(ctx).Where(p.Status.Eq(0)).Order(p.Sort, p.ID).Find()
	if err != nil {
		return nil, err
	}

	var res []*system.PostSimpleRespVO
	for _, item := range list {
		res = append(res, &system.PostSimpleRespVO{
			ID:   item.ID,
			Name: item.Name,
		})
	}
	return res, nil
}
