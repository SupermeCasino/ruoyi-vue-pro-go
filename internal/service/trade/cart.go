package trade

import (
	"context"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	productSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	pkgErrors "github.com/wxlbd/ruoyi-mall-go/pkg/errors"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CartService struct {
	q      *query.Query
	skuSvc *productSvc.ProductSkuService
	spuSvc *productSvc.ProductSpuService
}

func NewCartService(q *query.Query, skuSvc *productSvc.ProductSkuService, spuSvc *productSvc.ProductSpuService) *CartService {
	return &CartService{
		q:      q,
		skuSvc: skuSvc,
		spuSvc: spuSvc,
	}
}

// AddCart 添加购物车
func (s *CartService) AddCart(ctx context.Context, userId int64, r *req.AppCartAddReq) (int64, error) {
	c := s.q.Cart
	// 查询是否已存在
	cart, err := c.WithContext(ctx).Where(c.UserID.Eq(userId), c.SkuID.Eq(r.SkuID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	// 校验 SKU
	sku, err := s.skuSvc.GetSku(ctx, r.SkuID)
	if err != nil {
		return 0, err
	}
	if sku == nil {
		return 0, pkgErrors.NewBizError(1007001001, "商品 SKU 不存在")
	}
	// 校验库存
	newCount := r.Count
	if cart != nil {
		newCount = cart.Count + r.Count
	}
	if sku.Stock < newCount {
		return 0, pkgErrors.NewBizError(1007001002, "库存不足")
	}

	if cart != nil {
		// 更新数量
		_, err = c.WithContext(ctx).Where(c.ID.Eq(cart.ID)).Update(c.Count, newCount)
		return cart.ID, err
	}

	// 新增
	newCart := &trade.Cart{
		UserID:   userId,
		SpuID:    sku.SpuID,
		SkuID:    r.SkuID,
		Count:    r.Count,
		Selected: true,
	}
	err = c.WithContext(ctx).Create(newCart)
	return newCart.ID, err
}

// UpdateCartCount 更新购物车数量
func (s *CartService) UpdateCartCount(ctx context.Context, userId int64, r *req.AppCartUpdateCountReq) error {
	c := s.q.Cart
	cart, err := c.WithContext(ctx).Where(c.ID.Eq(r.ID), c.UserID.Eq(userId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.NewBizError(1007001003, "购物车项不存在")
		}
		return err
	}
	// 校验库存
	sku, err := s.skuSvc.GetSku(ctx, cart.SkuID)
	if err != nil {
		return err
	}
	if sku != nil && sku.Stock < r.Count {
		return pkgErrors.NewBizError(1007001002, "库存不足")
	}
	_, err = c.WithContext(ctx).Where(c.ID.Eq(r.ID)).Update(c.Count, r.Count)
	return err
}

// UpdateCartSelected 更新购物车选中状态
func (s *CartService) UpdateCartSelected(ctx context.Context, userId int64, r *req.AppCartUpdateSelectedReq) error {
	c := s.q.Cart
	_, err := c.WithContext(ctx).Where(c.UserID.Eq(userId), c.ID.In(r.IDs...)).Update(c.Selected, *r.Selected)
	return err
}

// ResetCart 重置购物车（换 SKU）
func (s *CartService) ResetCart(ctx context.Context, userId int64, r *req.AppCartResetReq) error {
	c := s.q.Cart
	// 校验原购物车项
	cart, err := c.WithContext(ctx).Where(c.ID.Eq(r.ID), c.UserID.Eq(userId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.NewBizError(1007001003, "购物车项不存在")
		}
		return err
	}

	// 校验新 SKU
	sku, err := s.skuSvc.GetSku(ctx, r.SkuID)
	if err != nil {
		return err
	}
	if sku == nil {
		return pkgErrors.NewBizError(1007001001, "商品 SKU 不存在")
	}
	if sku.Stock < r.Count {
		return pkgErrors.NewBizError(1007001002, "库存不足")
	}

	// 如果新 SKU 和原 SKU 相同，只更新数量
	if cart.SkuID == r.SkuID {
		_, err = c.WithContext(ctx).Where(c.ID.Eq(r.ID)).Update(c.Count, r.Count)
		return err
	}

	// 检查新 SKU 是否已在购物车
	existCart, _ := c.WithContext(ctx).Where(c.UserID.Eq(userId), c.SkuID.Eq(r.SkuID)).First()
	if existCart != nil {
		// 合并到已有项，删除原项
		_, _ = c.WithContext(ctx).Where(c.ID.Eq(existCart.ID)).Update(c.Count, existCart.Count+r.Count)
		_, _ = c.WithContext(ctx).Where(c.ID.Eq(r.ID)).Delete()
		return nil
	}

	// 更新为新 SKU
	_, err = c.WithContext(ctx).Where(c.ID.Eq(r.ID)).Updates(map[string]interface{}{
		"sku_id": r.SkuID,
		"spu_id": sku.SpuID,
		"count":  r.Count,
	})
	return err
}

// DeleteCart 删除购物车
func (s *CartService) DeleteCart(ctx context.Context, userId int64, ids []int64) error {
	c := s.q.Cart
	_, err := c.WithContext(ctx).Where(c.UserID.Eq(userId), c.ID.In(ids...)).Delete()
	return err
}

// GetCartCount 获取购物车商品数量
func (s *CartService) GetCartCount(ctx context.Context, userId int64) (int, error) {
	c := s.q.Cart
	var total int
	err := c.WithContext(ctx).Where(c.UserID.Eq(userId)).Select(c.Count.Sum()).Scan(&total)
	return total, err
}

// GetCartList 获取购物车列表
func (s *CartService) GetCartList(ctx context.Context, userId int64) (*resp.AppCartListResp, error) {
	c := s.q.Cart
	carts, err := c.WithContext(ctx).Where(c.UserID.Eq(userId)).Order(c.UpdateTime.Desc()).Find()
	if err != nil {
		return nil, err
	}

	if len(carts) == 0 {
		return &resp.AppCartListResp{ValidList: []resp.AppCartItem{}, InvalidList: []resp.AppCartItem{}}, nil
	}

	// 获取 SKU 信息
	skuIds := lo.Map(carts, func(item *trade.Cart, _ int) int64 { return item.SkuID })
	skuList, _ := s.skuSvc.GetSkuList(ctx, skuIds)
	skuMap := lo.KeyBy(skuList, func(item *resp.ProductSkuResp) int64 { return item.ID })

	// 获取 SPU 信息
	spuIds := lo.Uniq(lo.Map(carts, func(item *trade.Cart, _ int) int64 { return item.SpuID }))
	spuList, _ := s.spuSvc.GetSpuList(ctx, spuIds)
	spuMap := lo.KeyBy(spuList, func(item *resp.ProductSpuResp) int64 { return item.ID })

	var validList, invalidList []resp.AppCartItem
	for _, cart := range carts {
		sku := skuMap[cart.SkuID]
		spu := spuMap[cart.SpuID]

		item := resp.AppCartItem{
			ID:       cart.ID,
			Count:    cart.Count,
			Selected: bool(cart.Selected),
		}

		// 判断有效性
		isValid := sku != nil && spu != nil && spu.Status == 1 && sku.Stock > 0

		if spu != nil {
			item.Spu = &resp.AppCartSpuInfo{
				ID:     spu.ID,
				Name:   spu.Name,
				PicURL: spu.PicURL,
			}
		}
		if sku != nil {
			item.Sku = &resp.AppCartSkuInfo{
				ID:     sku.ID,
				PicURL: sku.PicURL,
				Price:  sku.Price,
				Stock:  sku.Stock,
			}
		}

		if isValid {
			validList = append(validList, item)
		} else {
			invalidList = append(invalidList, item)
		}
	}

	if validList == nil {
		validList = []resp.AppCartItem{}
	}
	if invalidList == nil {
		invalidList = []resp.AppCartItem{}
	}

	return &resp.AppCartListResp{
		ValidList:   validList,
		InvalidList: invalidList,
	}, nil
}
