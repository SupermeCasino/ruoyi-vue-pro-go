package product

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type ProductBrowseHistoryService struct {
	q      *query.Query
	spuSvc *ProductSpuService
}

func NewProductBrowseHistoryService(q *query.Query, spuSvc *ProductSpuService) *ProductBrowseHistoryService {
	return &ProductBrowseHistoryService{
		q:      q,
		spuSvc: spuSvc,
	}
}

// CreateBrowseHistory 创建浏览记录 (Async logic handled by caller or here? Java uses @Async. Here just standard sync, or go routine if truly needed. We'll do sync for simplicity first, or go routine inside)
func (s *ProductBrowseHistoryService) CreateBrowseHistory(ctx context.Context, userId, spuId int64) error {
	// Java: Check if exists? Java impl doesn't check existence usually, it just inserts logs.
	// But Wait, DO key is ID.
	// Typically browse history is "Latest view".
	// If allow duplicate?
	// Java ServiceImpl logic:
	// ProductBrowseHistoryDO history = browseHistoryMapper.selectByUserIdAndSpuId(userId, spuId);
	// if (history != null) {
	//     browseHistoryMapper.updateById(new ProductBrowseHistoryDO().setId(history.getId()).setUpdateTime(LocalDateTime.now()));
	//     return;
	// }
	// browseHistoryMapper.insert(new ProductBrowseHistoryDO().setUserId(userId).setSpuId(spuId));

	h := s.q.ProductBrowseHistory
	history, err := h.WithContext(ctx).Where(h.UserID.Eq(userId), h.SpuID.Eq(spuId)).First()
	if err == nil && history != nil {
		// Update time
		_, err := h.WithContext(ctx).Where(h.ID.Eq(history.ID)).Update(h.UpdatedAt, time.Now())
		return err
	}

	// Insert new
	newHistory := &product.ProductBrowseHistory{
		UserID:      userId,
		SpuID:       spuId,
		UserDeleted: false,
	}
	return h.WithContext(ctx).Create(newHistory)
}

// HideUserBrowseHistory 隐藏(删除)用户浏览记录
func (s *ProductBrowseHistoryService) HideUserBrowseHistory(ctx context.Context, userId int64, spuIds []int64) error {
	h := s.q.ProductBrowseHistory
	q := h.WithContext(ctx).Where(h.UserID.Eq(userId))

	if len(spuIds) > 0 {
		q = q.Where(h.SpuID.In(spuIds...))
	}

	// Java uses logical delete by updating user_deleted = true
	// "void hideUserBrowseHistory(Long userId, Collection<Long> spuId);"
	// productBrowseHistoryMapper.updateUserDeleted(userId, spuIds, true);

	_, err := q.Update(h.UserDeleted, true)
	return err
}

// GetBrowseHistoryPage (Admin & App share similar logic but differing reqs?)
// Java has separate ReqVOs but logic is similar. Admin uses BrowseHistoryPageReqVO, App uses AppBrowseHistoryPageReqVO
func (s *ProductBrowseHistoryService) GetBrowseHistoryPage(ctx context.Context, r *req.ProductBrowseHistoryPageReq) (*pagination.PageResult[resp.ProductBrowseHistoryResp], error) {
	h := s.q.ProductBrowseHistory
	q := h.WithContext(ctx).Where(h.UserDeleted.Is(false)) // Only show non-deleted to user? Java: "Boolean userDeleted" in params.

	if r.UserId > 0 {
		q = q.Where(h.UserID.Eq(r.UserId))
	}
	if r.SpuId > 0 {
		q = q.Where(h.SpuID.Eq(r.SpuId))
	}

	list, total, err := q.Order(h.UpdatedAt.Desc()).FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}

	// Fill SPU Info
	spuIds := lo.Map(list, func(item *product.ProductBrowseHistory, _ int) int64 {
		return item.SpuID
	})
	spuList, err := s.spuSvc.GetSpuList(ctx, spuIds)
	if err != nil {
		return nil, err
	}
	spuMap := lo.KeyBy(spuList, func(item *resp.ProductSpuResp) int64 { return item.ID })

	result := lo.Map(list, func(item *product.ProductBrowseHistory, _ int) resp.ProductBrowseHistoryResp {
		r := resp.ProductBrowseHistoryResp{
			ID:        item.ID,
			UserID:    item.UserID,
			SpuID:     item.SpuID,
			CreatedAt: item.CreatedAt,
		}
		if spu, ok := spuMap[item.SpuID]; ok {
			r.SpuName = spu.Name
			r.PicURL = spu.PicURL
			r.Price = int64(spu.Price)
			r.SalesCount = spu.SalesCount
			r.Stock = spu.Stock
		}
		return r
	})

	return &pagination.PageResult[resp.ProductBrowseHistoryResp]{
		List:  result,
		Total: total,
	}, nil
}

// GetAppBrowseHistoryPage
func (s *ProductBrowseHistoryService) GetAppBrowseHistoryPage(ctx context.Context, userId int64, r *req.AppProductBrowseHistoryPageReq) (*pagination.PageResult[resp.AppProductBrowseHistoryResp], error) {
	h := s.q.ProductBrowseHistory
	q := h.WithContext(ctx).Where(h.UserID.Eq(userId)).Where(h.UserDeleted.Is(false))

	list, total, err := q.Order(h.UpdatedAt.Desc()).FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}

	spuIds := lo.Map(list, func(item *product.ProductBrowseHistory, _ int) int64 {
		return item.SpuID
	})
	spuList, err := s.spuSvc.GetSpuList(ctx, spuIds)
	if err != nil {
		return nil, err
	}
	spuMap := lo.KeyBy(spuList, func(item *resp.ProductSpuResp) int64 { return item.ID })

	result := lo.Map(list, func(item *product.ProductBrowseHistory, _ int) resp.AppProductBrowseHistoryResp {
		r := resp.AppProductBrowseHistoryResp{
			ID:        item.ID,
			SpuID:     item.SpuID,
			CreatedAt: item.CreatedAt,
		}
		if spu, ok := spuMap[item.SpuID]; ok {
			r.SpuName = spu.Name
			r.PicURL = spu.PicURL
			r.Price = int64(spu.Price)
			r.SalesCount = spu.SalesCount
			r.Stock = spu.Stock
		}
		return r
	})

	return &pagination.PageResult[resp.AppProductBrowseHistoryResp]{
		List:  result,
		Total: total,
	}, nil
}
