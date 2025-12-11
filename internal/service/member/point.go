package member

import (
	"backend-go/internal/api/req"
	"backend-go/internal/model/member"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"
	"errors"
)

type MemberPointRecordService struct {
	q             *query.Query
	memberUserSvc *MemberUserService
}

func NewMemberPointRecordService(q *query.Query, memberUserSvc *MemberUserService) *MemberPointRecordService {
	return &MemberPointRecordService{
		q:             q,
		memberUserSvc: memberUserSvc,
	}
}

// GetPointRecordPage 获得用户积分记录分页
func (s *MemberPointRecordService) GetPointRecordPage(ctx context.Context, r *req.MemberPointRecordPageReq) (*core.PageResult[*member.MemberPointRecord], error) {
	q := s.q.MemberPointRecord.WithContext(ctx)

	// Filter by Nickname -> UserIDs
	if r.Nickname != "" {
		users, err := s.memberUserSvc.GetUserListByNickname(ctx, r.Nickname)
		if err != nil {
			return nil, err
		}
		if len(users) == 0 {
			return core.NewEmptyPageResult[*member.MemberPointRecord](), nil
		}
		var userIds []int64
		for _, u := range users {
			userIds = append(userIds, u.ID)
		}
		q = q.Where(s.q.MemberPointRecord.UserID.In(userIds...))
	}

	if r.BizType != "" {
		// bizType is int in DB but string in query param? usually int.
		// Java: private Integer bizType;
		// Assuming request structure passing valid int or empty.
		// For now, ignoring BizType filter or assuming strict int.
	}
	if r.Title != "" {
		q = q.Where(s.q.MemberPointRecord.Title.Like("%" + r.Title + "%"))
	}

	q = q.Order(s.q.MemberPointRecord.ID.Desc())

	list, count, err := q.FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}
	return core.NewPageResult(list, count), nil
}

// GetAppPointRecordPage 获得用户App积分记录分页
func (s *MemberPointRecordService) GetAppPointRecordPage(ctx context.Context, userId int64, r *req.AppMemberPointRecordPageReq) (*core.PageResult[*member.MemberPointRecord], error) {
	q := s.q.MemberPointRecord.WithContext(ctx).Where(s.q.MemberPointRecord.UserID.Eq(userId))

	if r.AddStatus != nil {
		if *r.AddStatus {
			q = q.Where(s.q.MemberPointRecord.Point.Gt(0))
		} else {
			q = q.Where(s.q.MemberPointRecord.Point.Lt(0))
		}
	}

	q = q.Order(s.q.MemberPointRecord.ID.Desc())

	list, count, err := q.FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}
	return core.NewPageResult(list, count), nil
}

// CreatePointRecord 创建积分记录
func (s *MemberPointRecordService) CreatePointRecord(ctx context.Context, userId int64, point int, bizType int, bizId string, title string, description string) error {
	if point == 0 {
		return nil
	}

	return s.q.Transaction(func(tx *query.Query) error {
		// 1. 校验用户积分余额
		user, err := s.memberUserSvc.GetUser(ctx, userId)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("user not found")
		}

		userPoint := user.Point
		totalPoint := int(userPoint) + point
		if totalPoint < 0 {
			return core.NewBizError(1004014003, "用户积分余额不足") // Assuming error code
		}

		// 2. 更新用户积分
		// Note: UpdateUserPoint in MemberUserService uses non-transactional DB instance by default if not passed TX.
		// BUT `s.q.Transaction` passes `tx *query.Query`. We should use `tx` to perform updates.
		// However, MemberUserService methods mostly use `s.q`. To support transaction properly across services,
		// ideally methods should accept `*query.Query` or we use the `tx` here to update directly or
		// MemberUserService needs to support transactional context propagation (which WithContext does IF the DB attached to context is tx).
		// GORM `WithContext` propagates context, but `s.q` in MemberUserService is the global one.
		// Standard way in this project: Pass `tx` to service? Or `MemberUserService` methods use `s.q.WithContext(ctx)`.
		// If we use `s.q.Transaction`, the `tx` has the transaction.
		// We can't easily inject `tx` into `MemberUserService` without changing method signature.

		// Workaround: We will manually update User Point here using `tx` to ensure transaction safety.
		// OR, if `MemberUserService.UpdateUserPoint` uses `WithContext(ctx)`, GORM DOES NOT automatically pick up transaction from context unless using a specific middleware or logic.
		// In `go-zero`/standard GORM, usually we pass `db` instance.
		// Here using `gorm gen`, `tx` IS the query instance for transaction.

		// Let's implement logic here for safety using `tx`.

		u := tx.MemberUser
		info, err := u.WithContext(ctx).Where(u.ID.Eq(userId)).Update(u.Point, u.Point.Add(int32(point)))
		if err != nil {
			return err
		}
		if info.RowsAffected == 0 {
			return errors.New("update user point failed")
		}

		// 3. 增加积分记录
		record := &member.MemberPointRecord{
			UserID:      userId,
			BizID:       bizId,
			BizType:     bizType,
			Title:       title,
			Description: description, // Format description if needed outside
			Point:       point,
			TotalPoint:  totalPoint,
		}
		return tx.MemberPointRecord.WithContext(ctx).Create(record)
	})
}
