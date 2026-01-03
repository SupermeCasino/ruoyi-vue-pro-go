package system

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type NotifyService struct {
	q *query.Query
	// Cache
	templateCache map[string]*model.SystemNotifyTemplate
	mu            sync.RWMutex
}

func NewNotifyService(q *query.Query) *NotifyService {
	s := &NotifyService{
		q: q,
	}
	s.RefreshCache()
	return s
}

func (s *NotifyService) RefreshCache() {
	t := s.q.SystemNotifyTemplate
	list, err := t.WithContext(context.Background()).Find()
	if err != nil {
		return
	}
	m := make(map[string]*model.SystemNotifyTemplate)
	for _, item := range list {
		m[item.Code] = item
	}
	s.mu.Lock()
	s.templateCache = m
	s.mu.Unlock()
}

// ================= Template CRUD =================

func (s *NotifyService) CreateNotifyTemplate(ctx context.Context, r *pay.NotifyTemplateCreateReq) (int64, error) {
	t := s.q.SystemNotifyTemplate
	template := &model.SystemNotifyTemplate{
		Name:     r.Name,
		Code:     r.Code,
		Nickname: r.Nickname,
		Content:  r.Content,
		Type:     r.Type,
		Status:   r.Status,
		Remark:   r.Remark,
	}
	if err := t.WithContext(ctx).Create(template); err != nil {
		return 0, err
	}
	s.RefreshCache()
	return template.ID, nil
}

func (s *NotifyService) UpdateNotifyTemplate(ctx context.Context, r *pay.NotifyTemplateUpdateReq) error {
	t := s.q.SystemNotifyTemplate
	_, err := t.WithContext(ctx).Where(t.ID.Eq(r.ID)).Updates(&model.SystemNotifyTemplate{
		Name:     r.Name,
		Code:     r.Code,
		Nickname: r.Nickname,
		Content:  r.Content,
		Type:     r.Type,
		Status:   r.Status,
		Remark:   r.Remark,
	})
	if err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

func (s *NotifyService) DeleteNotifyTemplate(ctx context.Context, id int64) error {
	t := s.q.SystemNotifyTemplate
	_, err := t.WithContext(ctx).Where(t.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

func (s *NotifyService) GetNotifyTemplate(ctx context.Context, id int64) (*model.SystemNotifyTemplate, error) {
	t := s.q.SystemNotifyTemplate
	return t.WithContext(ctx).Where(t.ID.Eq(id)).First()
}

func (s *NotifyService) GetNotifyTemplatePage(ctx context.Context, r *pay.NotifyTemplatePageReq) (*pagination.PageResult[*model.SystemNotifyTemplate], error) {
	t := s.q.SystemNotifyTemplate
	qb := t.WithContext(ctx)

	if r.Name != "" {
		qb = qb.Where(t.Name.Like("%" + r.Name + "%"))
	}
	if r.Code != "" {
		qb = qb.Where(t.Code.Like("%" + r.Code + "%"))
	}
	if r.Status != nil {
		qb = qb.Where(t.Status.Eq(*r.Status))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	offset := (r.PageNo - 1) * r.PageSize
	list, err := qb.Order(t.ID.Desc()).Offset(offset).Limit(r.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.SystemNotifyTemplate]{List: list, Total: total}, nil
}

// ================= Message Logic =================

func (s *NotifyService) SendNotify(ctx context.Context, userID int64, userType int, templateCode string, params map[string]interface{}) (int64, error) {
	s.mu.RLock()
	template, ok := s.templateCache[templateCode]
	s.mu.RUnlock()
	if !ok || template == nil {
		return 0, errors.NewBizError(1002006001, "站内信模板不存在")
	}

	// 对齐 Java: NotifySendServiceImpl.sendSingleNotify - 校验模板状态
	// Status: 0=开启, 1=禁用 (CommonStatusEnum: ENABLE=0, DISABLE=1)
	if template.Status == 1 {
		// 模板已禁用，静默返回（对齐 Java 的 log.info 并 return null）
		return 0, nil
	}

	// 对齐 Java: NotifySendServiceImpl.validateTemplateParams - 校验模板参数完整性
	if template.Params != "" {
		var requiredParams []string
		if err := json.Unmarshal([]byte(template.Params), &requiredParams); err == nil {
			for _, key := range requiredParams {
				if _, exists := params[key]; !exists {
					return 0, errors.NewBizError(1002006002, fmt.Sprintf("站内信模板参数 [%s] 缺失", key))
				}
			}
		}
	}

	content := template.Content
	for k, v := range params {
		content = strings.ReplaceAll(content, "{"+k+"}", fmt.Sprintf("%v", v))
	}

	paramsStr, _ := json.Marshal(params)
	msg := &model.SystemNotifyMessage{
		UserID:           userID,
		UserType:         userType,
		TemplateID:       template.ID,
		TemplateCode:     template.Code,
		TemplateNickname: template.Nickname,
		TemplateContent:  content,
		TemplateType:     template.Type,
		TemplateParams:   string(paramsStr),
		ReadStatus:       false,
	}

	m := s.q.SystemNotifyMessage
	if err := m.WithContext(ctx).Create(msg); err != nil {
		return 0, err
	}
	return msg.ID, nil
}

func (s *NotifyService) GetNotifyMessagePage(ctx context.Context, r *pay.NotifyMessagePageReq) (*pagination.PageResult[*model.SystemNotifyMessage], error) {
	m := s.q.SystemNotifyMessage
	qb := m.WithContext(ctx)

	if r.UserID != 0 {
		qb = qb.Where(m.UserID.Eq(r.UserID))
	}
	if r.UserType != 0 {
		qb = qb.Where(m.UserType.Eq(r.UserType))
	}
	if r.TemplateCode != "" {
		qb = qb.Where(m.TemplateCode.Like("%" + r.TemplateCode + "%"))
	}
	if r.TemplateType != nil {
		qb = qb.Where(m.TemplateType.Eq(*r.TemplateType))
	}
	if r.ReadStatus != nil {
		qb = qb.Where(m.ReadStatus.Is(*r.ReadStatus))
	}
	if r.StartDate != "" && r.EndDate != "" {
		// 使用 Between 查询时间范围
		qb = qb.Where(m.CreateTime.Between(parseTime(r.StartDate), parseTime(r.EndDate)))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	offset := (r.PageNo - 1) * r.PageSize
	list, err := qb.Order(m.ID.Desc()).Offset(offset).Limit(r.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.SystemNotifyMessage]{List: list, Total: total}, nil
}

func (s *NotifyService) GetMyNotifyMessagePage(ctx context.Context, userID int64, userType int, r *pay.MyNotifyMessagePageReq) (*pagination.PageResult[*model.SystemNotifyMessage], error) {
	m := s.q.SystemNotifyMessage
	qb := m.WithContext(ctx).Where(m.UserID.Eq(userID), m.UserType.Eq(userType))

	if r.ReadStatus != nil {
		qb = qb.Where(m.ReadStatus.Is(*r.ReadStatus))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	offset := (r.PageNo - 1) * r.PageSize
	list, err := qb.Order(m.ID.Desc()).Offset(offset).Limit(r.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*model.SystemNotifyMessage]{List: list, Total: total}, nil
}

func (s *NotifyService) UpdateNotifyMessageRead(ctx context.Context, userID int64, userType int, ids []int64) error {
	m := s.q.SystemNotifyMessage
	now := time.Now()
	_, err := m.WithContext(ctx).
		Where(m.ID.In(ids...), m.UserID.Eq(userID), m.UserType.Eq(userType)).
		UpdateSimple(m.ReadStatus.Value(true), m.ReadTime.Value(now))
	return err
}

func (s *NotifyService) UpdateAllNotifyMessageRead(ctx context.Context, userID int64, userType int) error {
	m := s.q.SystemNotifyMessage
	now := time.Now()
	_, err := m.WithContext(ctx).
		Where(m.UserID.Eq(userID), m.UserType.Eq(userType), m.ReadStatus.Is(false)).
		UpdateSimple(m.ReadStatus.Value(true), m.ReadTime.Value(now))
	return err
}

func (s *NotifyService) GetUnreadNotifyMessageCount(ctx context.Context, userID int64, userType int) (int64, error) {
	m := s.q.SystemNotifyMessage
	return m.WithContext(ctx).
		Where(m.UserID.Eq(userID), m.UserType.Eq(userType), m.ReadStatus.Is(false)).
		Count()
}

// GetNotifyMessage 获取单条站内信 (对齐 Java: NotifyMessageService.getNotifyMessage)
func (s *NotifyService) GetNotifyMessage(ctx context.Context, id int64) (*model.SystemNotifyMessage, error) {
	m := s.q.SystemNotifyMessage
	return m.WithContext(ctx).Where(m.ID.Eq(id)).First()
}

// GetUnreadNotifyMessageList 获取未读站内信列表 (对齐 Java: NotifyMessageService.getUnreadNotifyMessageList)
func (s *NotifyService) GetUnreadNotifyMessageList(ctx context.Context, userID int64, userType int, size int) ([]*model.SystemNotifyMessage, error) {
	m := s.q.SystemNotifyMessage
	if size <= 0 {
		size = 10 // Default size
	}
	return m.WithContext(ctx).
		Where(m.UserID.Eq(userID), m.UserType.Eq(userType), m.ReadStatus.Is(false)).
		Order(m.ID.Desc()).
		Limit(size).
		Find()
}

// parseTime 解析日期字符串
func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}
