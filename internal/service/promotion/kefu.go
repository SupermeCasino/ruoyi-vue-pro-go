package promotion

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

// KefuService 客服 Service
type KefuService interface {
	// CreateMessage 发送消息
	CreateMessage(ctx context.Context, r req.KefuMessageCreateReq, senderID int64, senderType int) (int64, error)
	// GetMessagePage 获得消息分页
	GetMessagePage(ctx context.Context, r req.KefuMessagePageReq) (*core.PageResult[resp.KefuMessageResp], error)
	// GetConversationPage 获得会话分页
	GetConversationPage(ctx context.Context, r req.KefuConversationPageReq) (*core.PageResult[resp.KefuConversationResp], error)
	// DeleteConversation 删除会话
	DeleteConversation(ctx context.Context, id int64) error
	// GetConversation 获得会话
	GetConversation(ctx context.Context, id int64) (*promotion.PromotionKefuConversation, error)
}

type kefuService struct {
	q *query.Query
}

func NewKefuService(q *query.Query) KefuService {
	return &kefuService{q: q}
}

func (s *kefuService) CreateMessage(ctx context.Context, r req.KefuMessageCreateReq, senderID int64, senderType int) (int64, error) {
	// 1. 自动查找或创建会话
	var convo *promotion.PromotionKefuConversation
	var err error
	convoRepo := s.q.PromotionKefuConversation
	msgRepo := s.q.PromotionKefuMessage

	// 简单逻辑：如果是用户发消息，且没有 conversationID，尝试查找该用户与任意客服的活跃会话，或者创建一个新的
	if r.ConversationID == 0 {
		if senderType == 1 { // User
			// 查找是否存在未结束的会话 (简化：同一用户只允许一个活跃会话? 或者这里先简单创建一个新的)
			// 实际业务可能需要分配客服逻辑。这里简化为：如果没有指定会话，就创建一个新的会话，AdminID 暂空或指定默认
			convo = &promotion.PromotionKefuConversation{
				UserID:          senderID,
				AdminID:         0, // 待分配
				Status:          0, // 接待中
				CreateTime:      time.Now(),
				UpdateTime:      time.Now(),
				LastMessageTime: time.Now(),
			}
			if err = convoRepo.WithContext(ctx).Create(convo); err != nil {
				return 0, err
			}
			r.ConversationID = convo.ID
		} else {
			return 0, core.NewBizError(400, "客服发送消息必须指定会话ID")
		}
	} else {
		convo, err = convoRepo.WithContext(ctx).Where(convoRepo.ID.Eq(r.ConversationID)).First()
		if err != nil {
			return 0, core.NewBizError(404, "会话不存在")
		}
	}

	// 2. 创建消息
	msg := &promotion.PromotionKefuMessage{
		ConversationID: r.ConversationID,
		SenderID:       senderID,
		SenderType:     senderType,
		ReceiverID:     0, // 需根据 convo 确定 receiver
		ReceiverType:   0,
		ContentType:    r.ContentType,
		Content:        r.Content,
		ReadStatus:     false,
		CreateTime:     time.Now(),
	}

	// 确定接收者
	if senderType == 1 { // User -> Admin
		msg.ReceiverID = convo.AdminID
		msg.ReceiverType = 2
		convo.AdminUnreadCount++
	} else { // Admin -> User
		msg.ReceiverID = convo.UserID
		msg.ReceiverType = 1
		convo.UserUnreadCount++
	}

	if err = msgRepo.WithContext(ctx).Create(msg); err != nil {
		return 0, err
	}

	// 3. 更新会话状态
	convo.LastMessageTime = time.Now()
	convo.LastMessageContent = r.Content
	convo.LastMessageType = r.ContentType
	if _, err = convoRepo.WithContext(ctx).Where(convoRepo.ID.Eq(convo.ID)).Updates(convo); err != nil {
		return 0, err
	}

	return msg.ID, nil
}

func (s *kefuService) GetMessagePage(ctx context.Context, r req.KefuMessagePageReq) (*core.PageResult[resp.KefuMessageResp], error) {
	msgRepo := s.q.PromotionKefuMessage
	q := msgRepo.WithContext(ctx).Where(msgRepo.ConversationID.Eq(r.ConversationID))

	// 分页查询
	list, count, err := q.Order(msgRepo.CreateTime.Desc()).FindByPage(r.PageNo-1, r.PageSize) // 注意：通常消息是倒序查还是正序？前端习惯倒序加载历史
	if err != nil {
		return nil, err
	}

	resList := make([]resp.KefuMessageResp, len(list))
	for i, v := range list {
		resList[i] = resp.KefuMessageResp{
			ID:             v.ID,
			ConversationID: v.ConversationID,
			SenderID:       v.SenderID,
			SenderType:     v.SenderType,
			ReceiverID:     v.ReceiverID,
			ReceiverType:   v.ReceiverType,
			ContentType:    v.ContentType,
			Content:        v.Content,
			ReadStatus:     v.ReadStatus,
			CreateTime:     v.CreateTime,
		}
	}

	return &core.PageResult[resp.KefuMessageResp]{
		List:  resList,
		Total: count,
	}, nil
}

func (s *kefuService) GetConversationPage(ctx context.Context, r req.KefuConversationPageReq) (*core.PageResult[resp.KefuConversationResp], error) {
	convoRepo := s.q.PromotionKefuConversation
	q := convoRepo.WithContext(ctx)
	// TODO: Add filters based on r

	list, count, err := q.Order(convoRepo.LastMessageTime.Desc()).FindByPage(r.PageNo-1, r.PageSize)
	if err != nil {
		return nil, err
	}

	resList := make([]resp.KefuConversationResp, len(list))
	for i, v := range list {
		// 这里简化了 UserInfo 的获取，实际应调用 memberService 获取用户昵称头像
		resList[i] = resp.KefuConversationResp{
			ID:                 v.ID,
			UserID:             v.UserID,
			AdminID:            v.AdminID,
			AdminUnreadCount:   v.AdminUnreadCount,
			UserUnreadCount:    v.UserUnreadCount,
			LastMessageTime:    v.LastMessageTime,
			LastMessageContent: v.LastMessageContent,
			LastMessageType:    v.LastMessageType,
			Status:             v.Status,
			CreateTime:         v.CreateTime,
		}
	}

	return &core.PageResult[resp.KefuConversationResp]{
		List:  resList,
		Total: count,
	}, nil
}

func (s *kefuService) DeleteConversation(ctx context.Context, id int64) error {
	convoRepo := s.q.PromotionKefuConversation
	// 软删除？这里示例为物理删除或标记删除
	// TODO: Implement Logic, e.g., check if exists
	_, err := convoRepo.WithContext(ctx).Where(convoRepo.ID.Eq(id)).Delete()
	return err
}

func (s *kefuService) GetConversation(ctx context.Context, id int64) (*promotion.PromotionKefuConversation, error) {
	convoRepo := s.q.PromotionKefuConversation
	return convoRepo.WithContext(ctx).Where(convoRepo.ID.Eq(id)).First()
}
