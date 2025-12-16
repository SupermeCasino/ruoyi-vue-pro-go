package promotion

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type PromotionBannerService struct {
	q *query.Query
}

func NewPromotionBannerService(q *query.Query) *PromotionBannerService {
	return &PromotionBannerService{q: q}
}

// CreateBanner 创建 Banner
func (s *PromotionBannerService) CreateBanner(ctx context.Context, r *req.PromotionBannerCreateReq) (int64, error) {
	banner := &promotion.PromotionBanner{
		Title:    r.Title,
		PicURL:   r.PicURL,
		Url:      r.Url,
		Status:   r.Status,
		Sort:     r.Sort,
		Position: r.Position,
		Memo:     r.Memo,
	}
	err := s.q.PromotionBanner.WithContext(ctx).Create(banner)
	return banner.ID, err
}

// UpdateBanner 更新 Banner
func (s *PromotionBannerService) UpdateBanner(ctx context.Context, r *req.PromotionBannerUpdateReq) error {
	_, err := s.q.PromotionBanner.WithContext(ctx).Where(s.q.PromotionBanner.ID.Eq(r.ID)).First()
	if err != nil {
		return errors.NewBizError(1004001000, "Banner不存在") // TODO: Error Code
	}
	_, err = s.q.PromotionBanner.WithContext(ctx).Where(s.q.PromotionBanner.ID.Eq(r.ID)).Updates(promotion.PromotionBanner{
		Title:    r.Title,
		PicURL:   r.PicURL,
		Url:      r.Url,
		Status:   r.Status,
		Sort:     r.Sort,
		Position: r.Position,
		Memo:     r.Memo,
	})
	return err
}

// DeleteBanner 删除 Banner
func (s *PromotionBannerService) DeleteBanner(ctx context.Context, id int64) error {
	_, err := s.q.PromotionBanner.WithContext(ctx).Where(s.q.PromotionBanner.ID.Eq(id)).Delete()
	return err
}

// GetBanner 获得 Banner
func (s *PromotionBannerService) GetBanner(ctx context.Context, id int64) (*resp.PromotionBannerResp, error) {
	item, err := s.q.PromotionBanner.WithContext(ctx).Where(s.q.PromotionBanner.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetBannerPage 获得 Banner 分页 (Admin)
func (s *PromotionBannerService) GetBannerPage(ctx context.Context, r *req.PromotionBannerPageReq) (*pagination.PageResult[*resp.PromotionBannerResp], error) {
	q := s.q.PromotionBanner.WithContext(ctx)
	if r.Title != "" {
		q = q.Where(s.q.PromotionBanner.Title.Like("%" + r.Title + "%"))
	}
	if r.Status != nil {
		q = q.Where(s.q.PromotionBanner.Status.Eq(*r.Status))
	}

	list, total, err := q.Order(s.q.PromotionBanner.Sort.Desc(), s.q.PromotionBanner.ID.Desc()).FindByPage(r.GetOffset(), r.PageSize)
	if err != nil {
		return nil, err
	}

	resList := lo.Map(list, func(item *promotion.PromotionBanner, _ int) *resp.PromotionBannerResp {
		return s.convertResp(item)
	})

	return &pagination.PageResult[*resp.PromotionBannerResp]{
		List:  resList,
		Total: total,
	}, nil
}

// GetInfoList 获得 App Banner 列表
// 对应 Java /app-api/promotion/banner/list
func (s *PromotionBannerService) GetAppBannerList(ctx context.Context, position int) ([]*resp.PromotionBannerResp, error) {
	q := s.q.PromotionBanner.WithContext(ctx).Where(s.q.PromotionBanner.Status.Eq(0)) // Enable
	if position > 0 {
		q = q.Where(s.q.PromotionBanner.Position.Eq(position))
	}
	list, err := q.Order(s.q.PromotionBanner.Sort.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *promotion.PromotionBanner, _ int) *resp.PromotionBannerResp {
		return s.convertResp(item)
	}), nil
}

func (s *PromotionBannerService) convertResp(item *promotion.PromotionBanner) *resp.PromotionBannerResp {
	return &resp.PromotionBannerResp{
		ID:        item.ID,
		Title:     item.Title,
		PicURL:    item.PicURL,
		Url:       item.Url,
		Status:    item.Status,
		Sort:      item.Sort,
		Position:  item.Position,
		Memo:      item.Memo,
		CreatedAt: item.CreatedAt,
	}
}
