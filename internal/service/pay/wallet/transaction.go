package wallet

import (
	"backend-go/internal/api/req"
	"backend-go/internal/model/pay"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
	"context"
	"strconv"
	"time"
)

type PayWalletTransactionService struct {
	q *query.Query
}

func NewPayWalletTransactionService(q *query.Query) *PayWalletTransactionService {
	return &PayWalletTransactionService{q: q}
}

// GetWalletTransactionPage 获得会员钱包流水分页
func (s *PayWalletTransactionService) GetWalletTransactionPage(ctx context.Context, req *req.PayWalletTransactionPageReq) (*core.PageResult[*pay.PayWalletTransaction], error) {
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
	return core.NewPageResult(list, total), nil
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
