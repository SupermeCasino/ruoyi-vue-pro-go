package wallet

import (
	"context"
	"strconv"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type PayWalletTransactionService struct {
	q *query.Query
}

func NewPayWalletTransactionService(q *query.Query) *PayWalletTransactionService {
	return &PayWalletTransactionService{q: q}
}

// GetWalletTransactionPage 获得会员钱包流水分页
func (s *PayWalletTransactionService) GetWalletTransactionPage(ctx context.Context, req *req.PayWalletTransactionPageReq) (*pagination.PageResult[*pay.PayWalletTransaction], error) {
	q := s.q.PayWalletTransaction.WithContext(ctx)
	if req.WalletID > 0 {
		q = q.Where(s.q.PayWalletTransaction.WalletID.Eq(req.WalletID))
	}
	if req.BizType > 0 {
		q = q.Where(s.q.PayWalletTransaction.BizType.Eq(req.BizType))
	}
	if req.BizID != "" {
		q = q.Where(s.q.PayWalletTransaction.BizID.Eq(req.BizID))
	}
	if req.No != "" {
		q = q.Where(s.q.PayWalletTransaction.No.Like("%" + req.No + "%"))
	}
	if req.Title != "" {
		q = q.Where(s.q.PayWalletTransaction.Title.Like("%" + req.Title + "%"))
	}
	q = q.Order(s.q.PayWalletTransaction.ID.Desc())

	list, total, err := q.FindByPage(req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}
	return pagination.NewPageResult(list, total), nil
}

// CreateWalletTransaction 创建/记录钱包流水
func (s *PayWalletTransactionService) CreateWalletTransaction(ctx context.Context, wallet *pay.PayWallet, bizType int, bizID string, title string, price int) (*pay.PayWalletTransaction, error) {
	// 生成流水号
	no := "WT" + strconv.FormatInt(time.Now().UnixNano(), 10)

	trx := &pay.PayWalletTransaction{
		WalletID: wallet.ID,
		BizType:  bizType,
		BizID:    bizID,
		No:       no,
		Title:    title,
		Price:    price,
		Balance:  wallet.Balance, // 交易后的余额? 需要确认为 updated balance
	}
	// Note: wallet.Balance should be the balance AFTER the transaction if we follow the ledger logic.

	err := s.q.PayWalletTransaction.WithContext(ctx).Create(trx)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

// GetWalletTransactionSummary 获得钱包流水统计
func (s *PayWalletTransactionService) GetWalletTransactionSummary(ctx context.Context, userId int64, userType int, createTime []time.Time) (totalIncome int, totalExpense int, err error) {
	// 1. 先查询钱包 ID
	wallet, err := s.q.PayWallet.WithContext(ctx).Where(s.q.PayWallet.UserID.Eq(userId), s.q.PayWallet.UserType.Eq(userType)).First()
	if err != nil {
		return 0, 0, err
	}

	q := s.q.PayWalletTransaction.WithContext(ctx).Where(s.q.PayWalletTransaction.WalletID.Eq(wallet.ID))

	if len(createTime) == 2 {
		q = q.Where(s.q.PayWalletTransaction.CreateTime.Between(createTime[0], createTime[1]))
	}

	// 统计支出 (Price < 0)
	err = q.Where(s.q.PayWalletTransaction.Price.Lt(0)).Select(s.q.PayWalletTransaction.Price.Sum()).Scan(&totalExpense)
	if err != nil {
		return 0, 0, err
	}
	// 统计收入 (Price > 0)
	err = q.Where(s.q.PayWalletTransaction.Price.Gt(0)).Select(s.q.PayWalletTransaction.Price.Sum()).Scan(&totalIncome)
	if err != nil {
		return 0, 0, err
	}

	// 支出取绝对值（对齐 Java，Java 返回的是正数支出）
	if totalExpense < 0 {
		totalExpense = -totalExpense
	}

	return totalIncome, totalExpense, nil
}
