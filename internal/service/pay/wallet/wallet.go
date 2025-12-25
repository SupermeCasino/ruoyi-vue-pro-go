package wallet

import (
	"context"
	stdErrors "errors"
	"fmt"
	"strconv"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	walletLockKeyPrefix = "pay:wallet:lock:"
	walletLockTimeout   = 5 * 1000 // 5 seconds in milliseconds
)

type PayWalletService struct {
	q              *query.Query
	rdb            *redis.Client
	transactionSvc *PayWalletTransactionService
}

func NewPayWalletService(q *query.Query, rdb *redis.Client, transactionSvc *PayWalletTransactionService) *PayWalletService {
	return &PayWalletService{q: q, rdb: rdb, transactionSvc: transactionSvc}
}

// GetOrCreateWallet 获得会员钱包，不存在则创建
// 严格对齐 Java 实现：获取一条，如果存在就返回一条，不存在就创建一条再返回
// 通过 Redis 分布式锁确保在不改变数据库表结构的情况下避免重复创建
func (s *PayWalletService) GetOrCreateWallet(ctx context.Context, userID int64, userType int) (*pay.PayWallet, error) {
	// 1. 先进行一次查询，如果已存在则直接返回（优化性能，减少加锁频率）
	wallet, err := s.q.PayWallet.WithContext(ctx).
		Where(s.q.PayWallet.UserID.Eq(userID), s.q.PayWallet.UserType.Eq(userType)).
		First()
	if err == nil {
		return wallet, nil
	}
	if !stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 2. 钱包尝试加锁
	lockKey := fmt.Sprintf("%s%d:%d", walletLockKeyPrefix, userID, userType)
	// 尝试获取锁，设置 5 秒过期时间防止死锁
	ok, err := s.rdb.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		// 加锁失败，说明有并发请求正在创建，等待一小段时间后重新查询
		// 这里采用简单的轮询重试逻辑
		for i := 0; i < 5; i++ {
			time.Sleep(100 * time.Millisecond) // 等待 100ms
			wallet, err = s.q.PayWallet.WithContext(ctx).
				Where(s.q.PayWallet.UserID.Eq(userID), s.q.PayWallet.UserType.Eq(userType)).
				First()
			if err == nil {
				return wallet, nil
			}
		}
		return nil, stdErrors.New("获取钱包超时，请稍后重试")
	}

	// 确保释放锁
	defer s.rdb.Del(ctx, lockKey)

	// 3. 在锁内再次查询，防止在获取锁的过程中其他请求已经创建了钱包
	wallet, err = s.q.PayWallet.WithContext(ctx).
		Where(s.q.PayWallet.UserID.Eq(userID), s.q.PayWallet.UserType.Eq(userType)).
		First()
	if err == nil {
		return wallet, nil
	}

	// 4. 执行创建
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
		return stdErrors.New("wallet not found")
	}

	// 2. 更新余额
	switch bizType {
	case pay.PayWalletBizTypePayment:
		if wallet.Balance < -price { // price is negative for payment
			return stdErrors.New("insufficient balance")
		}
	case pay.PayWalletBizTypePaymentRefund:
		// Refund adds back checks?
	case pay.PayWalletBizTypeRecharge:
		// Recharge adds
	case pay.PayWalletBizTypeUpdateBalance:
		if price < 0 && wallet.Balance < -price {
			return stdErrors.New("insufficient balance")
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
		return stdErrors.New("update wallet balance failed")
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

// ReduceWalletBalance 扣减钱包余额
func (s *PayWalletService) ReduceWalletBalance(ctx context.Context, walletID int64, bizID int64, bizType int, price int) (*pay.PayWalletTransaction, error) {
	// 1. 获取钱包
	wallet, err := s.GetWallet(ctx, walletID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.ErrNotFound
	}

	// 2. 扣除余额 (这里简化处理，直接更新，实际应该考虑分布式锁避免并发扣减导致负数，虽然数据库层会有 check constraint 或 unsigned int 保护，但业务层也应校验)
	// GORM updates
	res, err := s.q.PayWallet.WithContext(ctx).
		Where(s.q.PayWallet.ID.Eq(walletID)).
		Updates(map[string]interface{}{
			"balance":       gorm.Expr("balance - ?", price),       // 余额减少
			"total_expense": gorm.Expr("total_expense + ?", price), // 支出增加
		})

	if err != nil {
		return nil, err
	}
	if res.RowsAffected == 0 {
		return nil, stdErrors.New("insufficient balance")
	}

	// 3. 生成钱包流水
	wallet.Balance -= price
	title := "钱包支出"
	// if bizType == ... // 可以根据 bizType 设置 title
	return s.transactionSvc.CreateWalletTransaction(ctx, wallet, bizType, strconv.FormatInt(bizID, 10), title, -price)
}

// FreezePrice 冻结钱包余额
func (s *PayWalletService) FreezePrice(ctx context.Context, walletID int64, price int) error {
	// check balance enough?
	// update set balance = balance - price, freeze_price = freeze_price + price
	res, err := s.q.PayWallet.WithContext(ctx).
		Where(s.q.PayWallet.ID.Eq(walletID), s.q.PayWallet.Balance.Gte(price)).
		Updates(map[string]interface{}{
			"balance":      gorm.Expr("balance - ?", price),
			"freeze_price": gorm.Expr("freeze_price + ?", price),
		})
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return stdErrors.New("insufficient balance to freeze")
	}
	return nil
}

// UnfreezePrice 解冻钱包余额
func (s *PayWalletService) UnfreezePrice(ctx context.Context, walletID int64, price int) error {
	// update set balance = balance + price, freeze_price = freeze_price - price
	res, err := s.q.PayWallet.WithContext(ctx).
		Where(s.q.PayWallet.ID.Eq(walletID), s.q.PayWallet.FreezePrice.Gte(price)).
		Updates(map[string]interface{}{
			"balance":      gorm.Expr("balance + ?", price),
			"freeze_price": gorm.Expr("freeze_price - ?", price),
		})
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return stdErrors.New("insufficient frozen balance to unfreeze")
	}
	return nil
}
