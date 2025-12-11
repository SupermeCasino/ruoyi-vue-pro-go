package service

import (
	"context"

	"backend-go/internal/api/req"
	"backend-go/internal/model"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
)

type LoginLogService struct {
	q *query.Query
}

func NewLoginLogService(q *query.Query) *LoginLogService {
	return &LoginLogService{q: q}
}

// GetLoginLogPage 获取登录日志分页
func (s *LoginLogService) GetLoginLogPage(ctx context.Context, r *req.LoginLogPageReq) (*core.PageResult[*model.SystemLoginLog], error) {
	q := s.q.SystemLoginLog.WithContext(ctx)

	// 过滤条件
	if r.UserIP != "" {
		q = q.Where(s.q.SystemLoginLog.UserIP.Like("%" + r.UserIP + "%"))
	}
	if r.Username != "" {
		q = q.Where(s.q.SystemLoginLog.Username.Like("%" + r.Username + "%"))
	}
	if r.Status != nil {
		// status = true means result = 0 (success), status = false means result != 0
		if *r.Status {
			q = q.Where(s.q.SystemLoginLog.Result.Eq(0))
		} else {
			q = q.Where(s.q.SystemLoginLog.Result.Neq(0))
		}
	}
	if len(r.CreateTime) == 2 {
		q = q.Where(s.q.SystemLoginLog.CreatedAt.Between(r.CreateTime[0], r.CreateTime[1]))
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

	list, err := q.Order(s.q.SystemLoginLog.ID.Desc()).Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*model.SystemLoginLog]{
		List:  list,
		Total: total,
	}, nil
}
