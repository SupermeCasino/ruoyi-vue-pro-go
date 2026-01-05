package iot

import (
	"context"
	"log"

	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type OtaTaskService struct {
	otaFirmwareRepo   OtaFirmwareRepository
	otaTaskRepo       OtaTaskRepository
	otaTaskRecordRepo OtaTaskRecordRepository
	deviceSvc         *DeviceService
}

func NewOtaTaskService(
	otaFirmwareRepo OtaFirmwareRepository,
	otaTaskRepo OtaTaskRepository,
	otaTaskRecordRepo OtaTaskRecordRepository,
	deviceSvc *DeviceService,
) *OtaTaskService {
	return &OtaTaskService{
		otaFirmwareRepo:   otaFirmwareRepo,
		otaTaskRepo:       otaTaskRepo,
		otaTaskRecordRepo: otaTaskRecordRepo,
		deviceSvc:         deviceSvc,
	}
}

func (s *OtaTaskService) Create(ctx context.Context, r *iot2.IotOtaTaskCreateReqVO) (int64, error) {
	firmware, err := s.otaFirmwareRepo.GetByID(ctx, r.FirmwareID)
	if err != nil {
		return 0, err
	}
	if firmware == nil {
		return 0, model.ErrOtaFirmwareNotExists
	}

	task := &model.IotOtaTaskDO{
		Name:        r.Name,
		Description: r.Description,
		FirmwareID:  r.FirmwareID,
		Status:      consts.IotOtaTaskStatusWait,
		DeviceScope: r.DeviceScope,
	}

	if err := s.otaTaskRepo.Create(ctx, task); err != nil {
		return 0, err
	}

	// Upgrade records logic (simplified, ideally in transaction)
	var deviceIDs []int64
	if r.DeviceScope == consts.IotOtaDeviceScopeAll {
		// All devices for product logic...
	} else {
		deviceIDs = r.DeviceIDs
	}

	if len(deviceIDs) > 0 {
		records := make([]*model.IotOtaTaskRecordDO, 0, len(deviceIDs))
		for _, dID := range deviceIDs {
			records = append(records, &model.IotOtaTaskRecordDO{
				FirmwareID: r.FirmwareID,
				TaskID:     task.ID,
				DeviceID:   dID,
				Status:     consts.IotOtaRecordStatusWait,
			})
		}
		s.otaTaskRecordRepo.Create(ctx, records)
		task.DeviceTotalCount = int32(len(deviceIDs))
		s.otaTaskRepo.Update(ctx, task)
	}

	return task.ID, nil
}

func (s *OtaTaskService) GetPage(ctx context.Context, r *iot2.IotOtaTaskPageReqVO) (*pagination.PageResult[*model.IotOtaTaskDO], error) {
	return s.otaTaskRepo.GetPage(ctx, r)
}

func (s *OtaTaskService) Get(ctx context.Context, id int64) (*model.IotOtaTaskDO, error) {
	return s.otaTaskRepo.GetByID(ctx, id)
}

func (s *OtaTaskService) Cancel(ctx context.Context, id int64) error {
	task, err := s.otaTaskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return model.ErrOtaTaskNotExists
	}
	if task.Status != consts.IotOtaTaskStatusWait {
		return model.ErrOtaTaskStatusNotAllowCancel
	}
	task.Status = consts.IotOtaTaskStatusCancel
	return s.otaTaskRepo.Update(ctx, task)
}

func (s *OtaTaskService) GetRecordPage(ctx context.Context, r *iot2.IotOtaTaskRecordPageReqVO) (*pagination.PageResult[*model.IotOtaTaskRecordDO], error) {
	return s.otaTaskRecordRepo.GetPage(ctx, r)
}

// UpdateOtaRecordProgress 更新 OTA 升级记录进度
// version: 固件版本
// status: 升级状态
// progress: 升级进度 (0-100)
// description: 状态描述
func (s *OtaTaskService) UpdateOtaRecordProgress(ctx context.Context, device *model.IotDeviceDO, version string, status int, progress int, description string) error {
	// 1. 查询进行中的 OTA 升级记录
	records, err := s.otaTaskRecordRepo.GetListByDeviceIdAndStatus(ctx, device.ID, []int{
		consts.IotOtaRecordStatusWait,
		consts.IotOtaRecordStatusPushed,
		consts.IotOtaRecordStatusUpgrading,
	})
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return model.ErrOtaTaskRecordNotExists
	}

	// 取第一条记录进行更新
	record := records[0]

	// 2. 查询固件信息验证版本
	firmware, err := s.otaFirmwareRepo.GetByID(ctx, record.FirmwareID)
	if err != nil {
		return err
	}
	if firmware == nil {
		return model.ErrOtaFirmwareNotExists
	}
	if firmware.Version != version {
		// 版本不匹配，记录日志但不阻断
		log.Printf("[OtaTaskService] Version mismatch: expected=%s, got=%s", firmware.Version, version)
	}

	// 3. 更新 OTA 升级记录状态与进度
	record.Status = int8(status)
	record.Progress = int32(progress)
	record.Description = description
	if err := s.otaTaskRecordRepo.Update(ctx, record); err != nil {
		return err
	}

	// 4. 如果升级成功，更新设备固件版本
	if status == consts.IotOtaRecordStatusSuccess {
		if err := s.deviceSvc.UpdateDeviceFirmware(ctx, device.ID, firmware.ID); err != nil {
			log.Printf("[OtaTaskService] Failed to update device %d firmware: %v", device.ID, err)
		} else {
			log.Printf("[OtaTaskService] Device %d OTA success, firmware updated to %d", device.ID, firmware.ID)
		}
	}

	// 5. 检查是否所有记录都已完成，更新任务状态
	s.checkAndUpdateTaskStatus(ctx, record.TaskID)

	return nil
}

// checkAndUpdateTaskStatus 检查并更新任务状态
func (s *OtaTaskService) checkAndUpdateTaskStatus(ctx context.Context, taskID int64) {
	// 查询是否还有进行中的记录
	inProgressRecords, err := s.otaTaskRecordRepo.GetListByTaskIdAndStatus(ctx, taskID, []int{
		consts.IotOtaRecordStatusWait,
		consts.IotOtaRecordStatusPushed,
		consts.IotOtaRecordStatusUpgrading,
	})
	if err != nil {
		log.Printf("[OtaTaskService] Check task status error: %v", err)
		return
	}

	// 如果还有进行中的记录，不更新任务状态
	if len(inProgressRecords) > 0 {
		return
	}

	// 所有记录都已完成，更新任务状态为已结束
	task, err := s.otaTaskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return
	}
	task.Status = consts.IotOtaTaskStatusDone
	s.otaTaskRepo.Update(ctx, task)
	log.Printf("[OtaTaskService] Task %d completed", taskID)
}
