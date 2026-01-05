package iot

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	iot2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type DevicePropertyHandler struct {
	svc           *iotsvc.DevicePropertyService
	deviceSvc     *iotsvc.DeviceService
	thingModelSvc *iotsvc.ThingModelService
}

func NewDevicePropertyHandler(
	svc *iotsvc.DevicePropertyService,
	deviceSvc *iotsvc.DeviceService,
	thingModelSvc *iotsvc.ThingModelService,
) *DevicePropertyHandler {
	return &DevicePropertyHandler{
		svc:           svc,
		deviceSvc:     deviceSvc,
		thingModelSvc: thingModelSvc,
	}
}

func (h *DevicePropertyHandler) GetLatest(c *gin.Context) {
	deviceIDStr := c.Query("deviceId")
	if deviceIDStr == "" {
		response.WriteBizError(c, fmt.Errorf("设备编号不能为空"))
		return
	}
	var deviceID int64
	if _, err := fmt.Sscanf(deviceIDStr, "%d", &deviceID); err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 1.1 获取设备信息
	device, err := h.deviceSvc.Get(c.Request.Context(), deviceID)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	if device == nil {
		response.WriteBizError(c, fmt.Errorf("设备不存在"))
		return
	}

	// 1.2 获取设备最新属性
	properties, err := h.svc.GetLatestDeviceProperties(c.Request.Context(), deviceID)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 1.3 根据 productId + type 查询属性类型的物模型
	// type=1 为属性 (PROPERTY)
	thingModels, err := h.thingModelSvc.GetThingModelListByProductIdAndType(c.Request.Context(), device.ProductID, 1)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 2. 基于 thingModels 遍历，拼接 properties
	results := make([]*iot2.IotDevicePropertyDetailRespVO, 0, len(thingModels))
	for _, thingModel := range thingModels {
		detail := &iot2.IotDevicePropertyDetailRespVO{
			Name: thingModel.Name,
			IotDevicePropertyRespVO: iot2.IotDevicePropertyRespVO{
				Identifier: thingModel.Identifier,
			},
		}

		var property iot2.ThingModelProperty
		if len(thingModel.Property) > 0 {
			if err := json.Unmarshal(thingModel.Property, &property); err == nil {
				detail.DataType = property.DataType
				detail.DataSpecs = property.DataSpecs
				detail.DataSpecsList = property.DataSpecsList
			}
		}

		if p, ok := properties[thingModel.Identifier]; ok {
			var val interface{}
			if err := json.Unmarshal([]byte(p.Value), &val); err == nil {
				detail.Value = val
			} else {
				detail.Value = p.Value
			}
			detail.UpdateTime = p.UpdateTime.UnixMilli()
		}
		results = append(results, detail)
	}

	response.WriteSuccess(c, results)
}

func (h *DevicePropertyHandler) GetHistoryList(c *gin.Context) {
	var req iot2.IotDevicePropertyHistoryListReqVO
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, err)
		return
	}

	list, err := h.svc.GetHistoryDevicePropertyList(c.Request.Context(), &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	results := make([]*iot2.IotDevicePropertyRespVO, 0, len(list))
	for _, p := range list {
		vo := &iot2.IotDevicePropertyRespVO{
			Identifier: p.Identifier,
			UpdateTime: p.UpdateTime.UnixMilli(),
		}
		var val interface{}
		if err := json.Unmarshal([]byte(p.Value), &val); err == nil {
			vo.Value = val
		} else {
			vo.Value = p.Value
		}
		results = append(results, vo)
	}

	response.WriteSuccess(c, results)
}
