package system

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"gorm.io/gorm"
)

type MailService struct {
	db *gorm.DB
	// Caches
	accountCache  map[int64]*model.SystemMailAccount
	templateCache map[string]*model.SystemMailTemplate
	mu            sync.RWMutex
}

func NewMailService(db *gorm.DB) *MailService {
	s := &MailService{
		db: db,
	}
	s.RefreshCache()
	return s
}

// RefreshCache 刷新缓存
func (s *MailService) RefreshCache() {
	// Accounts
	var accounts []model.SystemMailAccount
	s.db.Find(&accounts)
	accountMap := make(map[int64]*model.SystemMailAccount)
	for i := range accounts {
		accountMap[accounts[i].ID] = &accounts[i]
	}

	// Templates
	var templates []model.SystemMailTemplate
	s.db.Find(&templates)
	templateMap := make(map[string]*model.SystemMailTemplate)
	for i := range templates {
		templateMap[templates[i].Code] = &templates[i]
	}

	s.mu.Lock()
	s.accountCache = accountMap
	s.templateCache = templateMap
	s.mu.Unlock()
}

// ================= Mail Account CRUD =================

func (s *MailService) CreateMailAccount(ctx context.Context, r *system.MailAccountCreateReq) (int64, error) {
	account := &model.SystemMailAccount{
		Mail:           r.Mail,
		Username:       r.Username,
		Password:       r.Password,
		Host:           r.Host,
		Port:           r.Port,
		SslEnable:      model.BitBool(r.SslEnable),
		StarttlsEnable: model.BitBool(r.StarttlsEnable),
	}
	if err := s.db.WithContext(ctx).Create(account).Error; err != nil {
		return 0, err
	}
	s.RefreshCache()
	return account.ID, nil
}

func (s *MailService) UpdateMailAccount(ctx context.Context, r *system.MailAccountUpdateReq) error {
	account := &model.SystemMailAccount{
		ID:             r.ID,
		Mail:           r.Mail,
		Username:       r.Username,
		Password:       r.Password,
		Host:           r.Host,
		Port:           r.Port,
		SslEnable:      model.BitBool(r.SslEnable),
		StarttlsEnable: model.BitBool(r.StarttlsEnable),
	}
	if err := s.db.WithContext(ctx).Updates(account).Error; err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

func (s *MailService) DeleteMailAccount(ctx context.Context, id int64) error {
	if err := s.db.WithContext(ctx).Delete(&model.SystemMailAccount{}, id).Error; err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

func (s *MailService) GetMailAccount(ctx context.Context, id int64) (*model.SystemMailAccount, error) {
	var account model.SystemMailAccount
	if err := s.db.WithContext(ctx).First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *MailService) GetMailAccountPage(ctx context.Context, r *system.MailAccountPageReq) (*pagination.PageResult[*model.SystemMailAccount], error) {
	db := s.db.WithContext(ctx).Model(&model.SystemMailAccount{})
	if r.Mail != "" {
		db = db.Where("mail LIKE ?", "%"+r.Mail+"%")
	}
	if r.Username != "" {
		db = db.Where("username LIKE ?", "%"+r.Username+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []*model.SystemMailAccount
	offset := (r.PageNo - 1) * r.PageSize
	if err := db.Order("id desc").Offset(offset).Limit(r.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return &pagination.PageResult[*model.SystemMailAccount]{List: list, Total: total}, nil
}

func (s *MailService) GetSimpleMailAccountList(ctx context.Context) ([]*model.SystemMailAccount, error) {
	var list []*model.SystemMailAccount
	if err := s.db.WithContext(ctx).Select("id", "mail").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ================= Mail Template CRUD =================

func (s *MailService) CreateMailTemplate(ctx context.Context, r *system.MailTemplateCreateReq) (int64, error) {
	template := &model.SystemMailTemplate{
		Name:      r.Name,
		Code:      r.Code,
		AccountID: r.AccountID,
		Nickname:  r.Nickname,
		Title:     r.Title,
		Content:   r.Content,
		Status:    r.Status,
		Remark:    r.Remark,
	}
	if err := s.db.WithContext(ctx).Create(template).Error; err != nil {
		return 0, err
	}
	s.RefreshCache()
	return template.ID, nil
}

func (s *MailService) UpdateMailTemplate(ctx context.Context, r *system.MailTemplateUpdateReq) error {
	template := &model.SystemMailTemplate{
		ID:        r.ID,
		Name:      r.Name,
		Code:      r.Code,
		AccountID: r.AccountID,
		Nickname:  r.Nickname,
		Title:     r.Title,
		Content:   r.Content,
		Status:    r.Status,
		Remark:    r.Remark,
	}
	if err := s.db.WithContext(ctx).Updates(template).Error; err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

func (s *MailService) DeleteMailTemplate(ctx context.Context, id int64) error {
	if err := s.db.WithContext(ctx).Delete(&model.SystemMailTemplate{}, id).Error; err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

func (s *MailService) GetMailTemplate(ctx context.Context, id int64) (*model.SystemMailTemplate, error) {
	var template model.SystemMailTemplate
	if err := s.db.WithContext(ctx).First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (s *MailService) GetMailTemplatePage(ctx context.Context, r *system.MailTemplatePageReq) (*pagination.PageResult[*model.SystemMailTemplate], error) {
	db := s.db.WithContext(ctx).Model(&model.SystemMailTemplate{})
	if r.Name != "" {
		db = db.Where("name LIKE ?", "%"+r.Name+"%")
	}
	if r.Code != "" {
		db = db.Where("code LIKE ?", "%"+r.Code+"%")
	}
	if r.AccountID != 0 {
		db = db.Where("account_id = ?", r.AccountID)
	}
	if r.Status != nil {
		db = db.Where("status = ?", *r.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []*model.SystemMailTemplate
	offset := (r.PageNo - 1) * r.PageSize
	if err := db.Order("id desc").Offset(offset).Limit(r.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return &pagination.PageResult[*model.SystemMailTemplate]{List: list, Total: total}, nil
}

// ================= Mail Sending Logic =================

func (s *MailService) SendMail(ctx context.Context, userID int64, userType int, toMail string, templateCode string, params map[string]interface{}) (int64, error) {
	// 1. Get Template
	s.mu.RLock()
	template, ok := s.templateCache[templateCode]
	s.mu.RUnlock()
	if !ok || template == nil {
		return 0, errors.NewBizError(1002005001, "邮件模板不存在")
	}

	// 2. Get Account
	s.mu.RLock()
	account, ok := s.accountCache[template.AccountID]
	s.mu.RUnlock()
	if !ok || account == nil {
		return 0, errors.NewBizError(1002005002, "邮箱账号不存在")
	}

	// 3. Render Content
	content := template.Content
	for k, v := range params {
		content = strings.ReplaceAll(content, "{"+k+"}", fmt.Sprintf("%v", v))
	}
	title := template.Title
	for k, v := range params {
		title = strings.ReplaceAll(title, "{"+k+"}", fmt.Sprintf("%v", v))
	}

	// 4. Send
	err := s.doSend(account, toMail, title, content)

	// 5. Log
	sendStatus := 1 // Success
	sendMessage := ""
	if err != nil {
		sendStatus = 0
		sendMessage = err.Error()
	}

	paramsStr, _ := json.Marshal(params)
	now := time.Now()
	log := &model.SystemMailLog{
		UserID:           userID,
		UserType:         userType,
		ToMail:           toMail,
		AccountID:        account.ID,
		FromMail:         account.Mail,
		TemplateID:       template.ID,
		TemplateCode:     template.Code,
		TemplateNickname: template.Nickname,
		TemplateTitle:    title,
		TemplateContent:  content,
		TemplateParams:   string(paramsStr),
		SendStatus:       sendStatus,
		SendTime:         &now,
		SendMessage:      sendMessage,
	}
	s.db.WithContext(ctx).Create(log)

	if err != nil {
		return log.ID, err
	}
	return log.ID, nil
}

func (s *MailService) doSend(account *model.SystemMailAccount, to string, subject string, body string) error {
	addr := fmt.Sprintf("%s:%d", account.Host, account.Port)
	auth := smtp.PlainAuth("", account.Username, account.Password, account.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(addr, auth, account.Mail, []string{to}, msg)
}

// ================= Mail Log =================

func (s *MailService) GetMailLogPage(ctx context.Context, r *system.MailLogPageReq) (*pagination.PageResult[*model.SystemMailLog], error) {
	db := s.db.WithContext(ctx).Model(&model.SystemMailLog{})
	if r.ToMail != "" {
		db = db.Where("to_mail LIKE ?", "%"+r.ToMail+"%")
	}
	if r.AccountID != 0 {
		db = db.Where("account_id = ?", r.AccountID)
	}
	if r.TemplateID != 0 {
		db = db.Where("template_id = ?", r.TemplateID)
	}
	if r.SendStatus != nil {
		db = db.Where("send_status = ?", *r.SendStatus)
	}
	if r.StartDate != "" && r.EndDate != "" {
		db = db.Where("create_time BETWEEN ? AND ?", r.StartDate, r.EndDate)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []*model.SystemMailLog
	offset := (r.PageNo - 1) * r.PageSize
	if err := db.Order("id desc").Offset(offset).Limit(r.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return &pagination.PageResult[*model.SystemMailLog]{List: list, Total: total}, nil
}
