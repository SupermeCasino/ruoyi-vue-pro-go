package wallet

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"context"
	"strconv"
	"time"
)

type PayWalletRechargeService struct {
	q         *query.Query
	walletSvc *PayWalletService
	trxSvc    *PayWalletTransactionService
	pkgSvc    *PayWalletRechargePackageService
}

func NewPayWalletRechargeService(q *query.Query, walletSvc *PayWalletService, trxSvc *PayWalletTransactionService, pkgSvc *PayWalletRechargePackageService) *PayWalletRechargeService {
	return &PayWalletRechargeService{
		q:         q,
		walletSvc: walletSvc,
		trxSvc:    trxSvc,
		pkgSvc:    pkgSvc,
	}
}

// CreateWalletRecharge 创建充值记录 (发起充值)
func (s *PayWalletRechargeService) CreateWalletRecharge(ctx context.Context, req *req.PayWalletRechargeCreateReq, userIP string) (*pay.PayWalletRecharge, error) {
	// 1. 校验钱包是否存在
	wallet, err := s.walletSvc.GetOrCreateWallet(ctx, req.UserID, req.UserType)
	if err != nil {
		return nil, err
	}

	// 2. 计算金额
	payPrice := req.PayPrice
	bonusPrice := req.BonusPrice
	if req.PackageID > 0 {
		pkg, err := s.pkgSvc.GetWalletRechargePackage(ctx, req.PackageID)
		if err != nil {
			return nil, err
		}
		if pkg == nil {
			return nil, errors.NewBizError(1006004001, "充值套餐不存在")
		}
		payPrice = pkg.PayPrice
		bonusPrice = pkg.BonusPrice
	}

	// 3. 创建充值记录
	recharge := &pay.PayWalletRecharge{
		WalletID:   wallet.ID,
		TotalPrice: payPrice + bonusPrice,
		PayPrice:   payPrice,
		BonusPrice: bonusPrice,
		PackageID:  req.PackageID,
		PayStatus:  false, // Waiting
	}
	err = s.q.PayWalletRecharge.WithContext(ctx).Create(recharge)
	if err != nil {
		return nil, err
	}

	return recharge, nil
}

// UpdateWalletRechargePaid 更新充值支付成功
func (s *PayWalletRechargeService) UpdateWalletRechargePaid(ctx context.Context, id int64, payOrderID int64) error {
	// 1. 获取充值记录
	recharge, err := s.q.PayWalletRecharge.WithContext(ctx).Where(s.q.PayWalletRecharge.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	if recharge.PayStatus {
		return nil // 已经支付，重复回调
	}

	// 2. 更新状态
	now := time.Now()
	_, err = s.q.PayWalletRecharge.WithContext(ctx).Where(s.q.PayWalletRecharge.ID.Eq(id)).Updates(pay.PayWalletRecharge{
		PayStatus:  true,
		PayOrderID: payOrderID,
		PayTime:    &now,
	})
	if err != nil {
		return err
	}

	// 3. 更新钱包余额
	return s.updateWalletBalance(ctx, recharge)
}

// updateWalletBalance 更新钱包余额
func (s *PayWalletRechargeService) updateWalletBalance(ctx context.Context, recharge *pay.PayWalletRecharge) error {
	wallet, err := s.walletSvc.GetWallet(ctx, recharge.WalletID)
	if err != nil {
		return err
	}

	// 更新余额
	newBalance := wallet.Balance + recharge.TotalPrice
	newTotalRecharge := wallet.TotalRecharge + recharge.PayPrice

	_, err = s.q.PayWallet.WithContext(ctx).Where(s.q.PayWallet.ID.Eq(wallet.ID)).Updates(pay.PayWallet{
		Balance:       newBalance,
		TotalRecharge: newTotalRecharge,
	})
	if err != nil {
		return err
	}

	// 记录流水
	_, err = s.trxSvc.CreateWalletTransaction(ctx, wallet, 1, strconv.FormatInt(recharge.ID, 10), "钱包充值", recharge.TotalPrice) // 1=充值
	return err
}

func (s *PayWalletRechargeService) GetWalletRechargePage(ctx context.Context, req *req.PayWalletRechargePageReq) (*pagination.PageResult[*pay.PayWalletRecharge], error) {
	q := s.q.PayWalletRecharge.WithContext(ctx)
	if req.PayStatus != nil {
		q = q.Where(s.q.PayWalletRecharge.PayStatus.Is(*req.PayStatus))
	}
	q = q.Order(s.q.PayWalletRecharge.ID.Desc())

	list, total, err := q.FindByPage(req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}
	return pagination.NewPageResult(list, total), nil
}
