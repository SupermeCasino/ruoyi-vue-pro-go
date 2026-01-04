package system

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

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
func (s *NoticeService) CreateNotice(ctx context.Context, req *system.NoticeSaveReq) (int64, error) {
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
func (s *NoticeService) UpdateNotice(ctx context.Context, req *system.NoticeSaveReq) error {
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
// 对应Java: NoticeServiceImpl.deleteNotice() - 删除前需要校验公告是否存在
func (s *NoticeService) DeleteNotice(ctx context.Context, id int64) error {
	n := s.q.SystemNotice

	// 1. 校验是否存在 (对应 Java: validateNoticeExists(id))
	notice, err := n.WithContext(ctx).Where(n.ID.Eq(id)).First()
	if err != nil {
		return errors.New("公告不存在")
	}
	if notice == nil {
		return errors.New("公告不存在")
	}

	// 2. 删除通知公告
	_, err = n.WithContext(ctx).Where(n.ID.Eq(id)).Delete()
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
func (s *NoticeService) GetNotice(ctx context.Context, id int64) (*system.NoticeRespVO, error) {
	n := s.q.SystemNotice
	item, err := n.WithContext(ctx).Where(n.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetNoticePage 获得通知公告分页
func (s *NoticeService) GetNoticePage(ctx context.Context, req *system.NoticePageReq) (*pagination.PageResult[*system.NoticeRespVO], error) {
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

	return &pagination.PageResult[*system.NoticeRespVO]{
		List:  lo.Map(list, func(item *model.SystemNotice, _ int) *system.NoticeRespVO { return s.convertResp(item) }),
		Total: total,
	}, nil
}

func (s *NoticeService) convertResp(item *model.SystemNotice) *system.NoticeRespVO {
	return &system.NoticeRespVO{
		ID:         item.ID,
		Title:      item.Title,
		Type:       item.Type,
		Content:    item.Content,
		Status:     item.Status,
		CreateTime: item.CreateTime,
		// ← 补充新增字段的映射 (对应Java: NoticeRespVO 继承自 BaseDO)
		UpdateTime: item.UpdateTime,
		Creator:    item.Creator,
		Updater:    item.Updater,
	}
}
