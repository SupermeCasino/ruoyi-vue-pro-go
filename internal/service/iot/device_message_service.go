package iot

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	iotcore "github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type DeviceMessageService struct {
	deviceMessageRepo DeviceMessageRepository
	deviceRepo        DeviceRepository
	devicePropertySvc *DevicePropertyService
	otaTaskSvc        *OtaTaskService
	messageBus        iotcore.MessageBus
}

func NewDeviceMessageService(
	deviceMessageRepo DeviceMessageRepository,
	deviceRepo DeviceRepository,
	devicePropertySvc *DevicePropertyService,
	otaTaskSvc *OtaTaskService,
	messageBus iotcore.MessageBus,
) *DeviceMessageService {
	return &DeviceMessageService{
		deviceMessageRepo: deviceMessageRepo,
		deviceRepo:        deviceRepo,
		devicePropertySvc: devicePropertySvc,
		otaTaskSvc:        otaTaskSvc,
		messageBus:        messageBus,
	}
}

// GetDeviceMessagePage 分页查询设备消息
func (s *DeviceMessageService) GetDeviceMessagePage(ctx context.Context, req *iot.IotDeviceMessagePageReqVO) (*pagination.PageResult[*model.IotDeviceMessageDO], error) {
	return s.deviceMessageRepo.GetPage(ctx, req)
}

// GetDeviceMessageListByRequestIdsAndReply 根据请求ID和回复标志查询消息列表
func (s *DeviceMessageService) GetDeviceMessageListByRequestIdsAndReply(ctx context.Context, deviceID int64, requestIDs []string, reply bool) ([]*model.IotDeviceMessageDO, error) {
	return s.deviceMessageRepo.GetListByRequestIdsAndReply(ctx, deviceID, requestIDs, reply)
}

// SendDeviceMessage 发送设备消息
func (s *DeviceMessageService) SendDeviceMessage(ctx context.Context, req *iot.IotDeviceMessageSendReqVO) error {
	// 1. 验证设备存在
	device, err := s.deviceRepo.GetByID(ctx, req.DeviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return model.ErrDeviceNotExists
	}

	// 2. 构建设备消息
	var params map[string]any
	if req.Params != nil {
		if p, ok := req.Params.(map[string]any); ok {
			params = p
		}
	}
	message := &iotcore.IotDeviceMessage{
		ID:         generateMessageID(),
		RequestID:  generateMessageID(),
		Method:     req.Method,
		Params:     params,
		DeviceID:   device.ID,
		TenantID:   device.TenantID,
		ReportTime: time.Now(),
	}

	// 3. 发送消息到消息总线
	return s.sendDeviceMessageInternal(ctx, message, device, "")
}

// SendDeviceMessageCore 发送核心设备消息（供内部调用）
func (s *DeviceMessageService) SendDeviceMessageCore(ctx context.Context, message *iotcore.IotDeviceMessage) error {
	device, err := s.deviceRepo.GetByID(ctx, message.DeviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return model.ErrDeviceNotExists
	}
	return s.sendDeviceMessageInternal(ctx, message, device, "")
}

// sendDeviceMessageInternal 内部发送消息逻辑
func (s *DeviceMessageService) sendDeviceMessageInternal(ctx context.Context, message *iotcore.IotDeviceMessage, device *model.IotDeviceDO, serverID string) error {
	// 1. 补充消息信息
	s.appendDeviceMessage(message, device)

	// 2.1 上行消息：发送到消息总线
	if message.IsUpstreamMessage() {
		s.messageBus.Post(iotcore.DeviceMessageTopic, message)
		return nil
	}

	// 2.2 下行消息：需要发送到指定网关
	if serverID == "" {
		// 从属性服务获取设备连接的 serverID
		if s.devicePropertySvc != nil {
			serverID = s.devicePropertySvc.GetDeviceServerId(ctx, device.ID)
		}
		if serverID == "" {
			serverID = "iot-gateway-01" // 默认网关
		}
	}

	// 发送到指定网关的主题
	gatewayTopic := "iot_gateway_downstream_" + serverID
	s.messageBus.Post(gatewayTopic, message)

	// 记录下行消息日志（异步）
	go s.createDeviceLogAsync(ctx, message)

	return nil
}

// HandleUpstreamDeviceMessage 处理上行设备消息
func (s *DeviceMessageService) HandleUpstreamDeviceMessage(ctx context.Context, message *iotcore.IotDeviceMessage, device *model.IotDeviceDO) {
	log.Printf("[DeviceMessageService] Processing upstream message: deviceId=%d, method=%s",
		device.ID, message.Method)

	// 1. 处理消息
	var replyData any
	var err error
	switch message.Method {
	case consts.IotDeviceMessageMethodStateUpdate:
		// 设备上下线
		err = s.handleStateUpdate(ctx, message, device)
	case consts.IotDeviceMessageMethodPropertyPost:
		// 属性上报
		err = s.handlePropertyPost(ctx, message, device)
	case consts.IotDeviceMessageMethodOtaProgress:
		// OTA 进度上报
		err = s.handleOtaProgress(ctx, message, device)
	default:
		log.Printf("[DeviceMessageService] Unknown method: %s", message.Method)
	}

	// 2. 记录消息日志
	go s.createDeviceLogAsync(ctx, message)

	// 3. 发送回复消息（如果需要）
	if !isReplyMessage(message.Method) && !isReplyDisabled(message.Method) && message.ServerID != "" {
		s.sendReplyMessage(ctx, message, device, replyData, err)
	}
}

// handleStateUpdate 处理状态更新
func (s *DeviceMessageService) handleStateUpdate(ctx context.Context, message *iotcore.IotDeviceMessage, device *model.IotDeviceDO) error {
	state, ok := message.Params["state"].(string)
	if !ok {
		return nil
	}

	var stateValue int8
	switch state {
	case "online":
		stateValue = consts.IotDeviceStateOnline
	case "offline":
		stateValue = consts.IotDeviceStateOffline
	default:
		return nil
	}

	// 更新设备状态
	device.State = stateValue
	return s.deviceRepo.Update(ctx, device)
}

// handlePropertyPost 处理属性上报
func (s *DeviceMessageService) handlePropertyPost(ctx context.Context, message *iotcore.IotDeviceMessage, device *model.IotDeviceDO) error {
	if message.Params == nil {
		return nil
	}

	// 调用属性服务保存属性
	if s.devicePropertySvc != nil {
		return s.devicePropertySvc.SaveDeviceProperty(ctx, device.ID, message.Params)
	}

	log.Printf("[DeviceMessageService] Property post: deviceId=%d, params=%v",
		device.ID, message.Params)
	return nil
}

// handleOtaProgress 处理 OTA 进度上报
func (s *DeviceMessageService) handleOtaProgress(ctx context.Context, message *iotcore.IotDeviceMessage, device *model.IotDeviceDO) error {
	log.Printf("[DeviceMessageService] OTA progress: deviceId=%d, params=%v",
		device.ID, message.Params)

	// 1. 解析 OTA 进度参数
	if message.Params == nil {
		return nil
	}

	// OTA 进度报文格式: {"version": "1.0.1", "status": 10, "progress": 50, "description": "downloading"}
	version, _ := message.Params["version"].(string)
	statusFloat, _ := message.Params["status"].(float64)
	status := int(statusFloat)
	progressFloat, _ := message.Params["progress"].(float64)
	progress := int(progressFloat)
	description, _ := message.Params["description"].(string)

	if version == "" {
		log.Printf("[DeviceMessageService] OTA progress missing version, skipping")
		return nil
	}

	log.Printf("[DeviceMessageService] OTA progress parsed: deviceId=%d, version=%s, status=%d, progress=%d, desc=%s",
		device.ID, version, status, progress, description)

	// 2. 调用 OtaTaskService 更新 OTA 记录进度
	if s.otaTaskSvc != nil {
		return s.otaTaskSvc.UpdateOtaRecordProgress(ctx, device, version, status, progress, description)
	}

	return nil
}

// sendReplyMessage 发送回复消息
func (s *DeviceMessageService) sendReplyMessage(ctx context.Context, message *iotcore.IotDeviceMessage, device *model.IotDeviceDO, data any, err error) {
	replyMethod := message.Method + "_reply"
	var code int
	var msg string
	if err != nil {
		code = 500
		msg = err.Error()
	}

	replyMessage := &iotcore.IotDeviceMessage{
		ID:         generateMessageID(),
		RequestID:  message.RequestID,
		Method:     replyMethod,
		Data:       data,
		Code:       code,
		Msg:        msg,
		DeviceID:   device.ID,
		TenantID:   device.TenantID,
		ServerID:   message.ServerID,
		ReportTime: time.Now(),
	}

	s.sendDeviceMessageInternal(ctx, replyMessage, device, message.ServerID)
}

// appendDeviceMessage 补充消息信息
func (s *DeviceMessageService) appendDeviceMessage(message *iotcore.IotDeviceMessage, device *model.IotDeviceDO) {
	if message.ID == "" {
		message.ID = generateMessageID()
	}
	if message.RequestID == "" {
		message.RequestID = message.ID
	}
	if message.ReportTime.IsZero() {
		message.ReportTime = time.Now()
	}
	message.DeviceID = device.ID
	message.TenantID = device.TenantID
}

// createDeviceLogAsync 异步创建设备日志
func (s *DeviceMessageService) createDeviceLogAsync(ctx context.Context, message *iotcore.IotDeviceMessage) {
	// 转换为数据库对象
	paramsJSON, _ := json.Marshal(message.Params)
	dataJSON, _ := json.Marshal(message.Data)

	messageDO := &model.IotDeviceMessageDO{
		ID:         message.ID,
		DeviceID:   message.DeviceID,
		ServerID:   message.ServerID,
		Method:     message.Method,
		RequestID:  message.RequestID,
		Params:     string(paramsJSON),
		Data:       string(dataJSON),
		Code:       message.Code,
		Msg:        message.Msg,
		Upstream:   message.IsUpstreamMessage(),
		Reply:      isReplyMessage(message.Method),
		Identifier: getIdentifier(message),
		ReportTime: message.ReportTime.UnixMilli(),
		TS:         message.ReportTime.UnixMilli(),
	}

	if err := s.deviceMessageRepo.Create(ctx, messageDO); err != nil {
		log.Printf("[DeviceMessageService] Create device log failed: %v", err)
	}
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return uuid.New().String()
}

// isReplyMessage 判断是否为回复消息
func isReplyMessage(method string) bool {
	return len(method) > 6 && method[len(method)-6:] == "_reply"
}

// isReplyDisabled 判断是否禁用回复
func isReplyDisabled(method string) bool {
	// 某些消息类型不需要回复
	disabledMethods := map[string]bool{
		"thing.lifecycle.state.update": true,
	}
	return disabledMethods[method]
}

// getIdentifier 获取消息标识符
func getIdentifier(message *iotcore.IotDeviceMessage) string {
	// 从 Params 中提取标识符
	if message.Params == nil {
		return ""
	}
	if id, ok := message.Params["identifier"].(string); ok {
		return id
	}
	// 对于状态更新，使用 state 作为标识符
	if state, ok := message.Params["state"].(string); ok {
		return state
	}
	return ""
}
