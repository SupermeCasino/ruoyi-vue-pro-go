package iot

import (
	"context"
	"encoding/json"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/datatypes"
)

type SceneRuleService struct {
	sceneRuleRepo SceneRuleRepository
}

func NewSceneRuleService(sceneRuleRepo SceneRuleRepository) *SceneRuleService {
	return &SceneRuleService{
		sceneRuleRepo: sceneRuleRepo,
	}
}

func (s *SceneRuleService) Create(ctx context.Context, r *iot2.IotSceneRuleSaveReqVO) (int64, error) {
	triggers, _ := json.Marshal(r.Triggers)
	actions, _ := json.Marshal(r.Actions)
	rule := &model.IotSceneRuleDO{
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		Triggers:    datatypes.JSON(triggers),
		Actions:     datatypes.JSON(actions),
	}
	if err := s.sceneRuleRepo.Create(ctx, rule); err != nil {
		return 0, err
	}
	return rule.ID, nil
}

func (s *SceneRuleService) Update(ctx context.Context, r *iot2.IotSceneRuleSaveReqVO) error {
	rule, err := s.sceneRuleRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if rule == nil {
		return model.ErrSceneRuleNotExists
	}

	triggers, _ := json.Marshal(r.Triggers)
	actions, _ := json.Marshal(r.Actions)

	rule.Name = r.Name
	rule.Description = r.Description
	rule.Status = r.Status
	rule.Triggers = datatypes.JSON(triggers)
	rule.Actions = datatypes.JSON(actions)

	return s.sceneRuleRepo.Update(ctx, rule)
}

func (s *SceneRuleService) UpdateStatus(ctx context.Context, id int64, status int8) error {
	rule, err := s.sceneRuleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rule == nil {
		return model.ErrSceneRuleNotExists
	}
	rule.Status = status
	return s.sceneRuleRepo.Update(ctx, rule)
}

func (s *SceneRuleService) Delete(ctx context.Context, id int64) error {
	return s.sceneRuleRepo.Delete(ctx, id)
}

func (s *SceneRuleService) Get(ctx context.Context, id int64) (*model.IotSceneRuleDO, error) {
	return s.sceneRuleRepo.GetByID(ctx, id)
}

func (s *SceneRuleService) GetPage(ctx context.Context, r *iot2.IotSceneRulePageReqVO) (*pagination.PageResult[*model.IotSceneRuleDO], error) {
	return s.sceneRuleRepo.GetPage(ctx, r)
}

func (s *SceneRuleService) GetListByStatus(ctx context.Context, status int8) ([]*model.IotSceneRuleDO, error) {
	return s.sceneRuleRepo.GetListByStatus(ctx, status)
}
