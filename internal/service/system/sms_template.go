package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type SmsTemplateService struct {
	q *query.Query
}

func NewSmsTemplateService(q *query.Query) *SmsTemplateService {
	return &SmsTemplateService{
		q: q,
	}
}

// CreateSmsTemplate 创建短信模板
func (s *SmsTemplateService) CreateSmsTemplate(ctx context.Context, req *system.SmsTemplateSaveReq) (int64, error) {
	// 校验 Channel 是否存在
	channel, err := s.q.SystemSmsChannel.WithContext(ctx).Where(s.q.SystemSmsChannel.ID.Eq(req.ChannelId)).First()
	if err != nil {
		return 0, errors.New("短信渠道不存在")
	}

	params := s.parseTemplateContentParams(req.Content)
	status := int32(0) // Default status (e.g., Enable/Normal)
	if req.Status != nil {
		status = *req.Status
	}

	template := &model.SystemSmsTemplate{
		Type:          req.Type,
		Status:        status,
		Code:          req.Code,
		Name:          req.Name,
		Content:       req.Content,
		Params:        params,
		Remark:        req.Remark,
		ApiTemplateId: req.ApiTemplateId,
		ChannelId:     req.ChannelId,
		ChannelCode:   channel.Code,
	}
	err = s.q.SystemSmsTemplate.WithContext(ctx).Create(template)
	return template.ID, err
}

// UpdateSmsTemplate 更新短信模板
func (s *SmsTemplateService) UpdateSmsTemplate(ctx context.Context, req *system.SmsTemplateSaveReq) error {
	t := s.q.SystemSmsTemplate
	_, err := t.WithContext(ctx).Where(t.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("短信模板不存在")
	}

	// 校验 Channel 是否存在
	channel, err := s.q.SystemSmsChannel.WithContext(ctx).Where(s.q.SystemSmsChannel.ID.Eq(req.ChannelId)).First()
	if err != nil {
		return errors.New("短信渠道不存在")
	}

	params := s.parseTemplateContentParams(req.Content)
	paramsBytes, _ := json.Marshal(params)

	updates := map[string]any{
		"type":            req.Type,
		"code":            req.Code,
		"name":            req.Name,
		"content":         req.Content,
		"params":          paramsBytes,
		"remark":          req.Remark,
		"api_template_id": req.ApiTemplateId,
		"channel_id":      req.ChannelId,
		"channel_code":    channel.Code,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	_, err = t.WithContext(ctx).Where(t.ID.Eq(req.ID)).Updates(updates)
	return err
}

// DeleteSmsTemplate 删除短信模板
func (s *SmsTemplateService) DeleteSmsTemplate(ctx context.Context, id int64) error {
	t := s.q.SystemSmsTemplate
	_, err := t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	return err
}

// GetSmsTemplate 获得短信模板
func (s *SmsTemplateService) GetSmsTemplate(ctx context.Context, id int64) (*system.SmsTemplateRespVO, error) {
	t := s.q.SystemSmsTemplate
	item, err := t.WithContext(ctx).Where(t.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetSmsTemplatePage 获得短信模板分页
func (s *SmsTemplateService) GetSmsTemplatePage(ctx context.Context, req *system.SmsTemplatePageReq) (*pagination.PageResult[*system.SmsTemplateRespVO], error) {
	t := s.q.SystemSmsTemplate
	qb := t.WithContext(ctx)

	if req.Type != nil {
		qb = qb.Where(t.Type.Eq(*req.Type))
	}
	if req.Status != nil {
		qb = qb.Where(t.Status.Eq(*req.Status))
	}
	if req.Code != "" {
		qb = qb.Where(t.Code.Like("%" + req.Code + "%"))
	}
	if req.Content != "" {
		qb = qb.Where(t.Content.Like("%" + req.Content + "%"))
	}
	if req.ApiTemplateId != "" {
		qb = qb.Where(t.ApiTemplateId.Like("%" + req.ApiTemplateId + "%"))
	}
	if req.ChannelId != nil {
		qb = qb.Where(t.ChannelId.Eq(*req.ChannelId))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(t.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*system.SmsTemplateRespVO]{
		List:  lo.Map(list, func(item *model.SystemSmsTemplate, _ int) *system.SmsTemplateRespVO { return s.convertResp(item) }),
		Total: total,
	}, nil
}

func (s *SmsTemplateService) convertResp(item *model.SystemSmsTemplate) *system.SmsTemplateRespVO {
	return &system.SmsTemplateRespVO{
		ID:            item.ID,
		Type:          item.Type,
		Status:        item.Status,
		Code:          item.Code,
		Name:          item.Name,
		Content:       item.Content,
		Params:        item.Params,
		Remark:        item.Remark,
		ApiTemplateId: item.ApiTemplateId,
		ChannelId:     item.ChannelId,
		ChannelCode:   item.ChannelCode,
		CreateTime:    item.CreateTime,
	}
}

// FormatSmsTemplateContent 格式化短信模板内容
func (s *SmsTemplateService) FormatSmsTemplateContent(content string, params map[string]interface{}) string {
	for k, v := range params {
		content = regexp.MustCompile(fmt.Sprintf(`\{%s\}`, k)).ReplaceAllString(content, fmt.Sprintf("%v", v))
	}
	return content
}

// parseTemplateContentParams 解析模板内容参数
// 例如：你好，{name}。你长的太{like}啦！ => [name, like]
func (s *SmsTemplateService) parseTemplateContentParams(content string) []string {
	re := regexp.MustCompile(`\{([a-zA-Z0-9]+)\}`)
	matches := re.FindAllStringSubmatch(content, -1)
	var params []string
	for _, match := range matches {
		if len(match) > 1 {
			params = append(params, match[1])
		}
	}
	return params
}
