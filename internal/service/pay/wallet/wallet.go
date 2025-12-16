package wallet

import (
	"context"
	"errors"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"

	"gorm.io/gorm"
)

type PayWalletService struct {
	q *query.Query
}

func NewPayWalletService(q *query.Query) *PayWalletService {
	return &PayWalletService{q: q}
}

// GetOrCreateWallet 获得会员钱包，不存在则创建
func (s *PayWalletService) GetOrCreateWallet(ctx context.Context, userID int64, userType int) (*pay.PayWallet, error) {
	// 1. 查询
	wallet, err := s.q.PayWallet.WithContext(ctx).
		Where(s.q.PayWallet.UserID.Eq(userID), s.q.PayWallet.UserType.Eq(userType)).
		First()
	if err == nil {
		return wallet, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 2. 创建
	wallet = &pay.PayWallet{
		UserID:        userID,
		UserType:      userType,
		Balance:       0,
		TotalExpense:  0,
		TotalRecharge: 0,
		FreezePrice:   0,
	}
	err = s.q.PayWallet.WithContext(ctx).Create(wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// GetWallet 获得会员钱包
func (s *PayWalletService) GetWallet(ctx context.Context, id int64) (*pay.PayWallet, error) {
	return s.q.PayWallet.WithContext(ctx).Where(s.q.PayWallet.ID.Eq(id)).First()
}

// GetWalletPage 获得会员钱包分页
func (s *PayWalletService) GetWalletPage(ctx context.Context, req *req.PayWalletPageReq) (*core.PageResult[*pay.PayWallet], error) {
	q := s.q.PayWallet.WithContext(ctx)
	if req.UserID > 0 {
		q = q.Where(s.q.PayWallet.UserID.Eq(req.UserID))
	}
	if req.UserType > 0 {
		q = q.Where(s.q.PayWallet.UserType.Eq(req.UserType))
	}
	q = q.Order(s.q.PayWallet.ID.Desc())

	list, total, err := q.FindByPage(req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}
	return core.NewPageResult(list, total), nil
}
