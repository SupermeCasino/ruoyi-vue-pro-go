package iot

import (
	"context"
	"encoding/json"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/datatypes"
)

type AlertConfigService struct {
	alertConfigRepo AlertConfigRepository
}

func NewAlertConfigService(alertConfigRepo AlertConfigRepository) *AlertConfigService {
	return &AlertConfigService{
		alertConfigRepo: alertConfigRepo,
	}
}

func (s *AlertConfigService) Create(ctx context.Context, r *iot2.IotAlertConfigSaveReqVO) (int64, error) {
	sceneRuleIDs, _ := json.Marshal(r.SceneRuleIDs)
	receiveUserIDs, _ := json.Marshal(r.ReceiveUserIDs)
	receiveTypes, _ := json.Marshal(r.ReceiveTypes)

	config := &model.IotAlertConfigDO{
		Name:           r.Name,
		Description:    r.Description,
		Level:          r.Level,
		Status:         r.Status,
		SceneRuleIDs:   datatypes.JSON(sceneRuleIDs),
		ReceiveUserIDs: datatypes.JSON(receiveUserIDs),
		ReceiveTypes:   datatypes.JSON(receiveTypes),
	}
	if err := s.alertConfigRepo.Create(ctx, config); err != nil {
		return 0, err
	}
	return config.ID, nil
}

func (s *AlertConfigService) Update(ctx context.Context, r *iot2.IotAlertConfigSaveReqVO) error {
	c, err := s.alertConfigRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if c == nil {
		return model.ErrAlertConfigNotExists
	}

	sceneRuleIDs, _ := json.Marshal(r.SceneRuleIDs)
	receiveUserIDs, _ := json.Marshal(r.ReceiveUserIDs)
	receiveTypes, _ := json.Marshal(r.ReceiveTypes)

	c.Name = r.Name
	c.Description = r.Description
	c.Level = r.Level
	c.Status = r.Status
	c.SceneRuleIDs = datatypes.JSON(sceneRuleIDs)
	c.ReceiveUserIDs = datatypes.JSON(receiveUserIDs)
	c.ReceiveTypes = datatypes.JSON(receiveTypes)

	return s.alertConfigRepo.Update(ctx, c)
}

func (s *AlertConfigService) Delete(ctx context.Context, id int64) error {
	return s.alertConfigRepo.Delete(ctx, id)
}

func (s *AlertConfigService) Get(ctx context.Context, id int64) (*model.IotAlertConfigDO, error) {
	return s.alertConfigRepo.GetByID(ctx, id)
}

func (s *AlertConfigService) GetPage(ctx context.Context, r *iot2.IotAlertConfigPageReqVO) (*pagination.PageResult[*model.IotAlertConfigDO], error) {
	return s.alertConfigRepo.GetPage(ctx, r)
}

func (s *AlertConfigService) GetListByStatus(ctx context.Context, status int8) ([]*model.IotAlertConfigDO, error) {
	return s.alertConfigRepo.GetListByStatus(ctx, status)
}

func (s *AlertConfigService) GetListBySceneRuleIdAndStatus(ctx context.Context, sceneRuleID int64, status int8) ([]*model.IotAlertConfigDO, error) {
	return s.alertConfigRepo.GetListBySceneRuleIdAndStatus(ctx, sceneRuleID, status)
}
