package wallet

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"gorm.io/gorm"
)

type PayWalletService struct {
	q              *query.Query
	transactionSvc *PayWalletTransactionService
}

func NewPayWalletService(q *query.Query, transactionSvc *PayWalletTransactionService) *PayWalletService {
	return &PayWalletService{q: q, transactionSvc: transactionSvc}
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
func (s *PayWalletService) GetWalletPage(ctx context.Context, req *req.PayWalletPageReq) (*pagination.PageResult[*pay.PayWallet], error) {
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
	return pagination.NewPageResult(list, total), nil
}

// AddWalletBalance 增加钱包余额
// price: 变动金额 (正数增加，负数减少)
func (s *PayWalletService) AddWalletBalance(ctx context.Context, walletID int64, bizID string, bizType int, price int) error {
	// 1. 获取钱包
	wallet, err := s.GetWallet(ctx, walletID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// 2. 更新余额
	switch bizType {
	case pay.PayWalletBizTypePayment:
		if wallet.Balance < -price { // price is negative for payment
			return errors.New("insufficient balance")
		}
	case pay.PayWalletBizTypePaymentRefund:
		// Refund adds back checks?
	case pay.PayWalletBizTypeRecharge:
		// Recharge adds
	case pay.PayWalletBizTypeUpdateBalance:
		if price < 0 && wallet.Balance < -price {
			return errors.New("insufficient balance")
		}
	}

	// Optimistic Update
	result, err := s.q.PayWallet.WithContext(ctx).
		Where(s.q.PayWallet.ID.Eq(wallet.ID)).
		Updates(map[string]interface{}{
			"balance": gorm.Expr("balance + ?", price),
			"total_expense": gorm.Expr("total_expense + ?", func() int {
				if price < 0 {
					return -price
				}
				return 0
			}()),
			"total_recharge": gorm.Expr("total_recharge + ?", func() int {
				if price > 0 {
					return price
				}
				return 0
			}()),
		})
	if err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return errors.New("update wallet balance failed")
	}

	// 3. 记录流水
	wallet.Balance += price // Approximate new balance for log
	title := "钱包余额更新"
	if bizType == pay.PayWalletBizTypeUpdateBalance {
		title = "管理员修改"
	}
	_, err = s.transactionSvc.CreateWalletTransaction(ctx, wallet, bizType, bizID, title, price)
	return err
}
