package handler

import (
	"fmt"
	"net/url"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type SmsLogHandler struct {
	smsLogSvc *service.SmsLogService
}

func NewSmsLogHandler(smsLogSvc *service.SmsLogService) *SmsLogHandler {
	return &SmsLogHandler{
		smsLogSvc: smsLogSvc,
	}
}

// GetSmsLogPage 获得短信日志分页
func (h *SmsLogHandler) GetSmsLogPage(c *gin.Context) {
	var req req.SmsLogPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.smsLogSvc.GetSmsLogPage(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// ExportSmsLogExcel 导出短信日志 Excel
func (h *SmsLogHandler) ExportSmsLogExcel(c *gin.Context) {
	var req req.SmsLogPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 设置分页大小为 -1（获取所有数据）
	req.PageSize = -1

	// 获取所有数据
	pageResult, err := h.smsLogSvc.GetSmsLogPage(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "短信日志"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	f.SetActiveSheet(index)

	// 设置列头
	headers := []string{"编号", "短信渠道编号", "短信渠道编码", "模板编号", "模板编码", "短信类型", "短信内容", "短信参数",
		"短信 API 的模板编号", "手机号", "用户编号", "用户类型", "发送状态", "发送时间", "短信 API 发送结果的编码",
		"短信 API 发送失败的提示", "短信 API 发送返回的唯一请求 ID", "短信 API 发送返回的序号",
		"接收状态", "接收时间", "API 接收结果的编码", "API 接收结果的说明", "创建时间"}

	for i, header := range headers {
		f.SetCellValue(sheetName, fmt.Sprintf("%c1", 'A'+i), header)
	}

	// 填充数据行
	for i, item := range pageResult.List {
		row := i + 2 // 从第2行开始
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.ChannelId)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.ChannelCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.TemplateId)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), item.TemplateCode)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), item.TemplateType)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), item.TemplateContent)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("%v", item.TemplateParams))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), item.ApiTemplateId)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), item.Mobile)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), item.UserId)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), item.UserType)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), item.SendStatus)
		sendTime := ""
		if item.SendTime != nil {
			sendTime = item.SendTime.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), sendTime)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), item.ApiSendCode)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), item.ApiSendMsg)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), item.ApiRequestId)
		f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), item.ApiSerialNo)
		f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), item.ReceiveStatus)
		receiveTime := ""
		if item.ReceiveTime != nil {
			receiveTime = item.ReceiveTime.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("T%d", row), receiveTime)
		f.SetCellValue(sheetName, fmt.Sprintf("U%d", row), item.ApiReceiveCode)
		f.SetCellValue(sheetName, fmt.Sprintf("V%d", row), item.ApiReceiveMsg)
		f.SetCellValue(sheetName, fmt.Sprintf("W%d", row), item.CreateTime.Format("2006-01-02 15:04:05"))
	}

	// 设置列宽（适应较长的内容）
	for col := 'A'; col <= 'W'; col++ {
		f.SetColWidth(sheetName, string(col), string(col), 18)
	}

	// 设置下载响应头
	filename := fmt.Sprintf("短信日志_%d.xlsx", time.Now().Unix())
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", url.QueryEscape(filename)))

	// 写入 Excel 到响应
	if err := f.Write(c.Writer); err != nil {
		response.WriteBizError(c, err)
		return
	}
}
