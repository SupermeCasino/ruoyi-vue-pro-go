package service

import (
	"context"
	"errors"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"

	"github.com/samber/lo"
)

type ConfigService struct {
	q *query.Query
}

func NewConfigService(q *query.Query) *ConfigService {
	return &ConfigService{
		q: q,
	}
}

// CreateConfig 创建参数配置
func (s *ConfigService) CreateConfig(ctx context.Context, req *req.ConfigSaveReq) (int64, error) {
	// TODO: Check key uniqueness if necessary
	config := &model.SystemConfig{
		Category:  req.Category,
		Name:      req.Name,
		ConfigKey: req.Key,
		Value:     req.Value,
		Type:      1, // 默认值 1，Java Service 中不从请求设置
		Visible:   model.BitBool(*req.Visible),
		Remark:    req.Remark,
	}
	err := s.q.SystemConfig.WithContext(ctx).Create(config)
	return config.ID, err
}

// UpdateConfig 修改参数配置
func (s *ConfigService) UpdateConfig(ctx context.Context, req *req.ConfigSaveReq) error {
	c := s.q.SystemConfig
	_, err := c.WithContext(ctx).Where(c.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("参数配置不存在")
	}

	_, err = c.WithContext(ctx).Where(c.ID.Eq(req.ID)).Updates(&model.SystemConfig{
		Category:  req.Category,
		Name:      req.Name,
		ConfigKey: req.Key,
		Value:     req.Value,
		// Type 不从请求更新，保持原值
		Visible: model.BitBool(*req.Visible),
		Remark:  req.Remark,
	})
	return err
}

// DeleteConfig 删除参数配置
func (s *ConfigService) DeleteConfig(ctx context.Context, id int64) error {
	c := s.q.SystemConfig
	_, err := c.WithContext(ctx).Where(c.ID.Eq(id)).Delete()
	return err
}

// DeleteConfigList 批量删除参数配置
func (s *ConfigService) DeleteConfigList(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	c := s.q.SystemConfig
	_, err := c.WithContext(ctx).Where(c.ID.In(ids...)).Delete()
	return err
}

// GetConfig 获得参数配置
func (s *ConfigService) GetConfig(ctx context.Context, id int64) (*resp.ConfigRespVO, error) {
	c := s.q.SystemConfig
	item, err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetConfigByKey 根据Key获得参数配置 (Internal use)
func (s *ConfigService) GetConfigByKey(ctx context.Context, key string) (*model.SystemConfig, error) {
	c := s.q.SystemConfig
	return c.WithContext(ctx).Where(c.ConfigKey.Eq(key)).First()
}

// GetConfigPage 获得参数配置分页
func (s *ConfigService) GetConfigPage(ctx context.Context, req *req.ConfigPageReq) (*core.PageResult[*resp.ConfigRespVO], error) {
	c := s.q.SystemConfig
	qb := c.WithContext(ctx)

	if req.Name != "" {
		qb = qb.Where(c.Name.Like("%" + req.Name + "%"))
	}
	if req.Key != "" {
		qb = qb.Where(c.ConfigKey.Like("%" + req.Key + "%"))
	}
	if req.Type != nil {
		qb = qb.Where(c.Type.Eq(*req.Type))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(c.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*resp.ConfigRespVO]{
		List:  lo.Map(list, func(item *model.SystemConfig, _ int) *resp.ConfigRespVO { return s.convertResp(item) }),
		Total: total,
	}, nil
}

func (s *ConfigService) convertResp(item *model.SystemConfig) *resp.ConfigRespVO {
	return &resp.ConfigRespVO{
		ID:         item.ID,
		Category:   item.Category,
		Name:       item.Name,
		Key:        item.ConfigKey,
		Value:      item.Value,
		Type:       item.Type,
		Visible:    bool(item.Visible),
		Remark:     item.Remark,
		CreateTime: item.CreatedAt,
	}
}
