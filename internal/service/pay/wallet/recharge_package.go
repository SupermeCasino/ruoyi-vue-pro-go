package wallet

import (
	"context"
	"errors"

	pay2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	pkgErrors "github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"gorm.io/gorm"
)

type PayWalletRechargePackageService struct {
	q *query.Query
}

func NewPayWalletRechargePackageService(q *query.Query) *PayWalletRechargePackageService {
	return &PayWalletRechargePackageService{q: q}
}

// CreateWalletRechargePackage 创建充值套餐
func (s *PayWalletRechargePackageService) CreateWalletRechargePackage(ctx context.Context, req *pay2.PayWalletRechargePackageCreateReq) (int64, error) {
	// 校验名是否重复
	exists, err := s.q.PayWalletRechargePackage.WithContext(ctx).Where(s.q.PayWalletRechargePackage.Name.Eq(req.Name)).First()
	if err == nil && exists != nil {
		return 0, pkgErrors.NewBizError(1006004000, "充值套餐名已存在") // PAY_WALLET_RECHARGE_PACKAGE_NAME_EXISTS
	}

	pkg := &pay.PayWalletRechargePackage{
		Name:       req.Name,
		PayPrice:   req.PayPrice,
		BonusPrice: req.BonusPrice,
		Status:     req.Status,
	}
	err = s.q.PayWalletRechargePackage.WithContext(ctx).Create(pkg)
	if err != nil {
		return 0, err
	}
	return pkg.ID, nil
}

// UpdateWalletRechargePackage 更新充值套餐
func (s *PayWalletRechargePackageService) UpdateWalletRechargePackage(ctx context.Context, req *pay2.PayWalletRechargePackageUpdateReq) error {
	// 校验存在
	oldPkg, err := s.validatePackageExists(ctx, req.ID)
	if err != nil {
		return err
	}

	// 校验名是否重复
	if req.Name != oldPkg.Name {
		exists, err := s.q.PayWalletRechargePackage.WithContext(ctx).Where(s.q.PayWalletRechargePackage.Name.Eq(req.Name)).First()
		if err == nil && exists != nil {
			return pkgErrors.NewBizError(1006004000, "充值套餐名已存在")
		}
	}

	_, err = s.q.PayWalletRechargePackage.WithContext(ctx).Where(s.q.PayWalletRechargePackage.ID.Eq(req.ID)).Updates(pay.PayWalletRechargePackage{
		Name:       req.Name,
		PayPrice:   req.PayPrice,
		BonusPrice: req.BonusPrice,
		Status:     req.Status,
	})
	return err
}

// DeleteWalletRechargePackage 删除充值套餐
func (s *PayWalletRechargePackageService) DeleteWalletRechargePackage(ctx context.Context, id int64) error {
	// 校验存在
	if _, err := s.validatePackageExists(ctx, id); err != nil {
		return err
	}
	// 删除
	_, err := s.q.PayWalletRechargePackage.WithContext(ctx).Where(s.q.PayWalletRechargePackage.ID.Eq(id)).Delete()
	return err
}

// GetWalletRechargePackage 获得充值套餐
func (s *PayWalletRechargePackageService) GetWalletRechargePackage(ctx context.Context, id int64) (*pay.PayWalletRechargePackage, error) {
	return s.q.PayWalletRechargePackage.WithContext(ctx).Where(s.q.PayWalletRechargePackage.ID.Eq(id)).First()
}

// GetWalletRechargePackagePage 获得充值套餐分页
func (s *PayWalletRechargePackageService) GetWalletRechargePackagePage(ctx context.Context, req *pay2.PayWalletRechargePackagePageReq) (*pagination.PageResult[*pay.PayWalletRechargePackage], error) {
	q := s.q.PayWalletRechargePackage.WithContext(ctx)
	if req.Name != "" {
		q = q.Where(s.q.PayWalletRechargePackage.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		q = q.Where(s.q.PayWalletRechargePackage.Status.Eq(*req.Status))
	}
	q = q.Order(s.q.PayWalletRechargePackage.PayPrice.Desc()) // 价格倒序

	list, total, err := q.FindByPage(req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}
	return pagination.NewPageResult(list, total), nil
}

// GetWalletRechargePackageList 获得充值套餐列表
func (s *PayWalletRechargePackageService) GetWalletRechargePackageList(ctx context.Context) ([]*pay.PayWalletRechargePackage, error) {
	// 只返回开启的，且删除标识未删除
	return s.q.PayWalletRechargePackage.WithContext(ctx).
		Where(s.q.PayWalletRechargePackage.Status.Eq(0)). // 0: 开启
		Order(s.q.PayWalletRechargePackage.PayPrice.Desc()).
		Find()
}

func (s *PayWalletRechargePackageService) validatePackageExists(ctx context.Context, id int64) (*pay.PayWalletRechargePackage, error) {
	pkg, err := s.q.PayWalletRechargePackage.WithContext(ctx).Where(s.q.PayWalletRechargePackage.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewBizError(1006004001, "充值套餐不存在") // PAY_WALLET_RECHARGE_PACKAGE_NOT_FOUND
		}
		return nil, err
	}
	return pkg, nil
}
