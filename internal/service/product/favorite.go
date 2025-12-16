package product

import (
	"context"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"

	"github.com/samber/lo"
)

type ProductFavoriteService struct {
	q      *query.Query
	spuSvc *ProductSpuService
}

func NewProductFavoriteService(q *query.Query, spuSvc *ProductSpuService) *ProductFavoriteService {
	return &ProductFavoriteService{
		q:      q,
		spuSvc: spuSvc,
	}
}

// CreateFavorite 创建商品收藏
func (s *ProductFavoriteService) CreateFavorite(ctx context.Context, userId, spuId int64) (int64, error) {
	// 1. 校验 SPU 是否存在 (Optional, but good practice. Java checks Favorite exists first)
	// Java impl: check if favorite exists.
	f := s.q.ProductFavorite
	count, err := f.WithContext(ctx).Where(f.UserID.Eq(userId), f.SpuID.Eq(spuId)).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, core.NewBizError(1006003000, "商品已收藏") // FAVORITE_EXISTS
	}

	// 2. 插入
	entity := &product.ProductFavorite{
		UserID: userId,
		SpuID:  spuId,
	}
	if err := f.WithContext(ctx).Create(entity); err != nil {
		return 0, err
	}
	return entity.ID, nil
}

// DeleteFavorite 取消商品收藏
func (s *ProductFavoriteService) DeleteFavorite(ctx context.Context, userId, spuId int64) error {
	f := s.q.ProductFavorite
	info, err := f.WithContext(ctx).Where(f.UserID.Eq(userId), f.SpuID.Eq(spuId)).First()
	if err != nil {
		return core.NewBizError(1006003001, "商品未收藏") // FAVORITE_NOT_EXISTS
	}

	if _, err := f.WithContext(ctx).Where(f.ID.Eq(info.ID)).Delete(); err != nil {
		return err
	}
	return nil
}

// GetFavoritePage (Admin)
func (s *ProductFavoriteService) GetFavoritePage(ctx context.Context, r *req.ProductFavoritePageReq) (*core.PageResult[resp.ProductFavoriteResp], error) {
	f := s.q.ProductFavorite
	q := f.WithContext(ctx)

	if r.UserId > 0 {
		q = q.Where(f.UserID.Eq(r.UserId))
	}
	if r.SpuId > 0 {
		q = q.Where(f.SpuID.Eq(r.SpuId))
	}

	list, total, err := q.FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}

	// Fill SPU Info
	spuIds := lo.Map(list, func(item *product.ProductFavorite, _ int) int64 {
		return item.SpuID
	})
	spuList, err := s.spuSvc.GetSpuList(ctx, spuIds)
	if err != nil {
		return nil, err
	}
	spuMap := lo.KeyBy(spuList, func(item *resp.ProductSpuResp) int64 { return item.ID })

	result := lo.Map(list, func(item *product.ProductFavorite, _ int) resp.ProductFavoriteResp {
		r := resp.ProductFavoriteResp{
			ID:        item.ID,
			UserID:    item.UserID,
			SpuID:     item.SpuID,
			CreatedAt: item.CreatedAt,
		}
		if spu, ok := spuMap[item.SpuID]; ok {
			r.SpuName = spu.Name
			r.PicURL = spu.PicURL
			r.Price = int64(spu.Price)
		}
		return r
	})

	return &core.PageResult[resp.ProductFavoriteResp]{
		List:  result,
		Total: total,
	}, nil
}

// GetAppFavoritePage (App)
func (s *ProductFavoriteService) GetAppFavoritePage(ctx context.Context, userId int64, r *req.AppFavoritePageReq) (*core.PageResult[resp.AppFavoriteResp], error) {
	f := s.q.ProductFavorite
	q := f.WithContext(ctx).Where(f.UserID.Eq(userId))

	list, total, err := q.FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}

	spuIds := lo.Map(list, func(item *product.ProductFavorite, _ int) int64 {
		return item.SpuID
	})
	spuList, err := s.spuSvc.GetSpuList(ctx, spuIds)
	if err != nil {
		return nil, err
	}
	spuMap := lo.KeyBy(spuList, func(item *resp.ProductSpuResp) int64 { return item.ID })

	result := lo.Map(list, func(item *product.ProductFavorite, _ int) resp.AppFavoriteResp {
		r := resp.AppFavoriteResp{
			ID:        item.ID,
			SpuID:     item.SpuID,
			CreatedAt: item.CreatedAt,
		}
		if spu, ok := spuMap[item.SpuID]; ok {
			r.SpuName = spu.Name
			r.PicURL = spu.PicURL
			r.Price = int64(spu.Price)
		}
		return r
	})

	return &core.PageResult[resp.AppFavoriteResp]{
		List:  result,
		Total: total,
	}, nil
}

// GetFavorite 检查是否收藏
func (s *ProductFavoriteService) GetFavorite(ctx context.Context, userId, spuId int64) (*product.ProductFavorite, error) {
	f := s.q.ProductFavorite
	return f.WithContext(ctx).Where(f.UserID.Eq(userId), f.SpuID.Eq(spuId)).First()
}

// GetFavoriteCount 获得收藏数量
func (s *ProductFavoriteService) GetFavoriteCount(ctx context.Context, userId int64) (int64, error) {
	f := s.q.ProductFavorite
	return f.WithContext(ctx).Where(f.UserID.Eq(userId)).Count()
}
