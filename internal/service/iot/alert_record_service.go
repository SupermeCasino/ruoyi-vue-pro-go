package iot

import (
	"context"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type AlertRecordService struct {
	alertRecordRepo AlertRecordRepository
}

func NewAlertRecordService(alertRecordRepo AlertRecordRepository) *AlertRecordService {
	return &AlertRecordService{
		alertRecordRepo: alertRecordRepo,
	}
}

func (s *AlertRecordService) GetPage(ctx context.Context, r *iot2.IotAlertRecordPageReqVO) (*pagination.PageResult[*model.IotAlertRecordDO], error) {
	return s.alertRecordRepo.GetPage(ctx, r)
}

func (s *AlertRecordService) Process(ctx context.Context, r *iot2.IotAlertRecordProcessReqVO) error {
	record, err := s.alertRecordRepo.GetByID(ctx, r.ID)
	if err != nil {
		return err
	}
	if record == nil {
		return model.ErrAlertRecordNotExists
	}

	record.ProcessStatus = true
	record.ProcessRemark = r.ProcessRemark
	return s.alertRecordRepo.Update(ctx, record)
}

func (s *AlertRecordService) Get(ctx context.Context, id int64) (*model.IotAlertRecordDO, error) {
	return s.alertRecordRepo.GetByID(ctx, id)
}

func (s *AlertRecordService) GetListBySceneRuleId(ctx context.Context, sceneRuleID int64, deviceID *int64, processStatus *bool) ([]*model.IotAlertRecordDO, error) {
	return s.alertRecordRepo.GetListBySceneRuleId(ctx, sceneRuleID, deviceID, processStatus)
}

func (s *AlertRecordService) CreateAlertRecord(ctx context.Context, config *model.IotAlertConfigDO, sceneRuleID int64, deviceMessage string) (int64, error) {
	record := &model.IotAlertRecordDO{
		ConfigID:      config.ID,
		ConfigName:    config.Name,
		ConfigLevel:   config.Level,
		SceneRuleID:   sceneRuleID,
		DeviceMessage: deviceMessage,
		ProcessStatus: false,
	}
	if err := s.alertRecordRepo.Create(ctx, record); err != nil {
		return 0, err
	}
	return record.ID, nil
}
