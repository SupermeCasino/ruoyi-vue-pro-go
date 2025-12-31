package system

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type OperateLogService struct {
	q *query.Query
}

func NewOperateLogService(q *query.Query) *OperateLogService {
	return &OperateLogService{q: q}
}

// GetOperateLogPage 获取操作日志分页
func (s *OperateLogService) GetOperateLogPage(ctx context.Context, r *req.OperateLogPageReq) (*pagination.PageResult[*model.SystemOperateLog], error) {
	q := s.q.SystemOperateLog.WithContext(ctx)

	// 过滤条件
	if r.UserID != nil {
		q = q.Where(s.q.SystemOperateLog.UserID.Eq(*r.UserID))
	}
	if r.BizID != nil {
		q = q.Where(s.q.SystemOperateLog.BizID.Eq(*r.BizID))
	}
	if r.Type != "" {
		q = q.Where(s.q.SystemOperateLog.Type.Like("%" + r.Type + "%"))
	}
	if r.SubType != "" {
		q = q.Where(s.q.SystemOperateLog.SubType.Like("%" + r.SubType + "%"))
	}
	if r.Action != "" {
		q = q.Where(s.q.SystemOperateLog.Action.Like("%" + r.Action + "%"))
	}
	if len(r.CreateTime) == 2 {
		q = q.Where(s.q.SystemOperateLog.CreateTime.Between(r.CreateTime[0], r.CreateTime[1]))
	}

	// 分页
	pageNo := r.PageNo
	pageSize := r.PageSize
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Order(s.q.SystemOperateLog.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.SystemOperateLog]{
		List:  list,
		Total: total,
	}, nil
}
