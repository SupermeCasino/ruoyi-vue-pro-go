package iot

import (
	"context"
	"encoding/json"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"gorm.io/datatypes"
)

type DataRuleService struct {
	dataRuleRepo DataRuleRepository
}

func NewDataRuleService(dataRuleRepo DataRuleRepository) *DataRuleService {
	return &DataRuleService{
		dataRuleRepo: dataRuleRepo,
	}
}

func (s *DataRuleService) Create(ctx context.Context, r *iot2.IotDataRuleSaveReqVO) (int64, error) {
	sourceConfigs, _ := json.Marshal(r.SourceConfigs)
	sinkIDs, _ := json.Marshal(r.SinkIDs)
	rule := &model.IotDataRuleDO{
		Name:          r.Name,
		Description:   r.Description,
		Status:        r.Status,
		SourceConfigs: datatypes.JSON(sourceConfigs),
		SinkIDs:       datatypes.JSON(sinkIDs),
	}
	if err := s.dataRuleRepo.Create(ctx, rule); err != nil {
		return 0, err
	}
	return rule.ID, nil
}

func (s *DataRuleService) Update(ctx context.Context, r *iot2.IotDataRuleSaveReqVO) error {
	rule, err := s.dataRuleRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if rule == nil {
		return model.ErrDataRuleNotExists
	}

	sourceConfigs, _ := json.Marshal(r.SourceConfigs)
	sinkIDs, _ := json.Marshal(r.SinkIDs)

	rule.Name = r.Name
	rule.Description = r.Description
	rule.Status = r.Status
	rule.SourceConfigs = datatypes.JSON(sourceConfigs)
	rule.SinkIDs = datatypes.JSON(sinkIDs)

	return s.dataRuleRepo.Update(ctx, rule)
}

func (s *DataRuleService) Delete(ctx context.Context, id int64) error {
	return s.dataRuleRepo.Delete(ctx, id)
}

func (s *DataRuleService) Get(ctx context.Context, id int64) (*model.IotDataRuleDO, error) {
	return s.dataRuleRepo.GetByID(ctx, id)
}

func (s *DataRuleService) GetPage(ctx context.Context, r *iot2.IotDataRulePageReqVO) (*pagination.PageResult[*model.IotDataRuleDO], error) {
	return s.dataRuleRepo.GetPage(ctx, r)
}
