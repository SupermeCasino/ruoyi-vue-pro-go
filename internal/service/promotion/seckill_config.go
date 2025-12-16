package promotion

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type SeckillConfigService struct {
	q *query.Query
}

func NewSeckillConfigService(q *query.Query) *SeckillConfigService {
	return &SeckillConfigService{q: q}
}

// CreateSeckillConfig 创建秒杀时段
func (s *SeckillConfigService) CreateSeckillConfig(ctx context.Context, r *req.SeckillConfigCreateReq) (int64, error) {
	// TODO: Verify params (startTime < endTime) if needed
	do := &promotion.PromotionSeckillConfig{
		Name:          r.Name,
		StartTime:     r.StartTime,
		EndTime:       r.EndTime,
		SliderPicUrls: r.SliderPicUrls,
		Status:        r.Status,
	}
	err := s.q.PromotionSeckillConfig.WithContext(ctx).Create(do)
	if err != nil {
		return 0, err
	}
	return do.ID, nil
}

// UpdateSeckillConfig 更新秒杀时段
func (s *SeckillConfigService) UpdateSeckillConfig(ctx context.Context, r *req.SeckillConfigUpdateReq) error {
	q := s.q.PromotionSeckillConfig
	_, err := q.WithContext(ctx).Where(q.ID.Eq(r.ID)).First()
	if err != nil {
		return errors.NewBizError(1001001000, "秒杀时段不存在")
	}

	_, err = q.WithContext(ctx).Where(q.ID.Eq(r.ID)).Updates(&promotion.PromotionSeckillConfig{
		Name:          r.Name,
		StartTime:     r.StartTime,
		EndTime:       r.EndTime,
		SliderPicUrls: r.SliderPicUrls,
		Status:        r.Status,
	})
	return err
}

// UpdateSeckillConfigStatus 更新秒杀时段状态
func (s *SeckillConfigService) UpdateSeckillConfigStatus(ctx context.Context, id int64, status int) error {
	q := s.q.PromotionSeckillConfig
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1001001000, "秒杀时段不存在")
	}
	_, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Status, status)
	return err
}

// DeleteSeckillConfig 删除秒杀时段
func (s *SeckillConfigService) DeleteSeckillConfig(ctx context.Context, id int64) error {
	q := s.q.PromotionSeckillConfig
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
	if err != nil {
		return errors.NewBizError(1001001000, "秒杀时段不存在")
	}
	_, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Delete()
	return err
}

// GetSeckillConfig 获得秒杀时段
func (s *SeckillConfigService) GetSeckillConfig(ctx context.Context, id int64) (*promotion.PromotionSeckillConfig, error) {
	q := s.q.PromotionSeckillConfig
	return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

// GetSeckillConfigList 获得所有秒杀时段列表
func (s *SeckillConfigService) GetSeckillConfigList(ctx context.Context) ([]*promotion.PromotionSeckillConfig, error) {
	q := s.q.PromotionSeckillConfig
	return q.WithContext(ctx).Order(q.StartTime).Find()
}

// GetSeckillConfigListByStatus 获得指定状态的秒杀时段列表
func (s *SeckillConfigService) GetSeckillConfigListByStatus(ctx context.Context, status int) ([]*promotion.PromotionSeckillConfig, error) {
	q := s.q.PromotionSeckillConfig
	return q.WithContext(ctx).Where(q.Status.Eq(status)).Order(q.StartTime).Find()
}

// GetSeckillConfigPage 分页获得秒杀时段
func (s *SeckillConfigService) GetSeckillConfigPage(ctx context.Context, r *req.SeckillConfigPageReq) (*pagination.PageResult[*promotion.PromotionSeckillConfig], error) {
	q := s.q.PromotionSeckillConfig
	do := q.WithContext(ctx)

	if r.Name != "" {
		do = do.Where(q.Name.Like("%" + r.Name + "%"))
	}
	if r.Status != nil {
		do = do.Where(q.Status.Eq(*r.Status))
	}

	do = do.Order(q.StartTime)

	list, count, err := do.FindByPage(r.PageNo, r.PageSize)
	if err != nil {
		return nil, err
	}
	return &pagination.PageResult[*promotion.PromotionSeckillConfig]{
		List:  list,
		Total: count,
	}, nil
}

// ValidateSeckillConfigExists 校验秒杀时段是否存在
func (s *SeckillConfigService) ValidateSeckillConfigExists(ctx context.Context, configIds []int64) error {
	if len(configIds) == 0 {
		return nil
	}
	q := s.q.PromotionSeckillConfig
	count, err := q.WithContext(ctx).Where(q.ID.In(configIds...)).Count()
	if err != nil {
		return err
	}
	if int(count) != len(configIds) {
		return errors.NewBizError(1001001000, "秒杀时段不存在")
	}
	return nil
}

// GetCurrentSeckillConfig 获得当前秒杀时段
func (s *SeckillConfigService) GetCurrentSeckillConfig(ctx context.Context) (*promotion.PromotionSeckillConfig, error) {
	list, err := s.GetSeckillConfigListByStatus(ctx, 1) // 1=Enabled
	if err != nil {
		return nil, err
	}
	now := time.Now()
	currentTimeStr := now.Format("15:04:05")

	for _, config := range list {
		// Ensure format consistency (add :00 if missing HH:mm -> HH:mm:00)
		start := config.StartTime
		end := config.EndTime
		if len(start) == 5 {
			start += ":00"
		}
		if len(end) == 5 {
			end += ":00"
		}

		if currentTimeStr >= start && currentTimeStr <= end {
			return config, nil
		}
	}
	return nil, nil
}
