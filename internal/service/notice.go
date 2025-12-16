package service

import (
	"context"
	"errors"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"

	"github.com/samber/lo"
)

type NoticeService struct {
	q *query.Query
}

func NewNoticeService(q *query.Query) *NoticeService {
	return &NoticeService{
		q: q,
	}
}

// CreateNotice 创建通知公告
func (s *NoticeService) CreateNotice(ctx context.Context, req *req.NoticeSaveReq) (int64, error) {
	notice := &model.SystemNotice{
		Title:   req.Title,
		Type:    *req.Type,
		Content: req.Content,
		Status:  *req.Status,
	}
	err := s.q.SystemNotice.WithContext(ctx).Create(notice)
	return notice.ID, err
}

// UpdateNotice 修改通知公告
func (s *NoticeService) UpdateNotice(ctx context.Context, req *req.NoticeSaveReq) error {
	n := s.q.SystemNotice
	_, err := n.WithContext(ctx).Where(n.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("公告不存在")
	}

	_, err = n.WithContext(ctx).Where(n.ID.Eq(req.ID)).Updates(&model.SystemNotice{
		Title:   req.Title,
		Type:    *req.Type,
		Content: req.Content,
		Status:  *req.Status,
	})
	return err
}

// DeleteNotice 删除通知公告
func (s *NoticeService) DeleteNotice(ctx context.Context, id int64) error {
	n := s.q.SystemNotice
	_, err := n.WithContext(ctx).Where(n.ID.Eq(id)).Delete()
	return err
}

// DeleteNoticeList 批量删除通知公告
func (s *NoticeService) DeleteNoticeList(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	n := s.q.SystemNotice
	_, err := n.WithContext(ctx).Where(n.ID.In(ids...)).Delete()
	return err
}

// GetNotice 获得通知公告
func (s *NoticeService) GetNotice(ctx context.Context, id int64) (*resp.NoticeRespVO, error) {
	n := s.q.SystemNotice
	item, err := n.WithContext(ctx).Where(n.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetNoticePage 获得通知公告分页
func (s *NoticeService) GetNoticePage(ctx context.Context, req *req.NoticePageReq) (*core.PageResult[*resp.NoticeRespVO], error) {
	n := s.q.SystemNotice
	qb := n.WithContext(ctx)

	if req.Title != "" {
		qb = qb.Where(n.Title.Like("%" + req.Title + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(n.Status.Eq(*req.Status))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(n.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*resp.NoticeRespVO]{
		List:  lo.Map(list, func(item *model.SystemNotice, _ int) *resp.NoticeRespVO { return s.convertResp(item) }),
		Total: total,
	}, nil
}

func (s *NoticeService) convertResp(item *model.SystemNotice) *resp.NoticeRespVO {
	return &resp.NoticeRespVO{
		ID:         item.ID,
		Title:      item.Title,
		Type:       item.Type,
		Content:    item.Content,
		Status:     item.Status,
		CreateTime: item.CreatedAt,
	}
}
