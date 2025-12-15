package handler

import (
	"net/http"
	"time"

	"backend-go/internal/pkg/utils"
	ws "backend-go/internal/pkg/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

type WebSocketHandler struct {
	manager *ws.Manager
	logger  *zap.Logger
}

func NewWebSocketHandler(manager *ws.Manager, logger *zap.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		manager: manager,
		logger:  logger,
	}
}

// Handle 处理 WebSocket 连接
// Java: LoginUserHandshakeInterceptor + JsonWebSocketMessageHandler + WebSocketSessionHandlerDecorator
func (h *WebSocketHandler) Handle(c *gin.Context) {
	// 1. 从 query 参数获取 token 进行认证
	// Java: TokenAuthenticationFilter 在握手前通过 ?token=xxx 认证
	token := c.Query("token")
	if token == "" {
		h.logger.Warn("WebSocket 连接缺少 token")
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "缺少 token"})
		return
	}

	// 2. 验证 token 并获取用户信息
	claims, err := utils.ParseToken(token)
	if err != nil {
		h.logger.Warn("WebSocket token 验证失败", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token 无效"})
		return
	}

	// 3. 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("WebSocket 升级失败", zap.Error(err))
		return
	}

	// 4. 创建 Session 并添加到 Manager
	// Java: WebSocketSessionHandlerDecorator.afterConnectionEstablished
	session := &ws.Session{
		ID:       uuid.New().String(),
		Conn:     conn,
		UserID:   claims.UserID,
		UserType: claims.UserType,
		TenantID: claims.TenantID,
	}
	h.manager.Add(session)
	h.logger.Info("WebSocket 连接建立",
		zap.String("sessionId", session.ID),
		zap.Int64("userId", session.UserID),
		zap.Int("userType", session.UserType),
	)

	// 5. 启动消息读取循环
	go h.readLoop(session)
}

// readLoop 消息读取循环
// Java: JsonWebSocketMessageHandler.handleTextMessage
func (h *WebSocketHandler) readLoop(session *ws.Session) {
	defer func() {
		// 连接关闭时从 Manager 移除
		// Java: WebSocketSessionHandlerDecorator.afterConnectionClosed
		h.manager.Remove(session.ID)
		_ = session.Conn.Close()
		h.logger.Info("WebSocket 连接关闭", zap.String("sessionId", session.ID))
	}()

	// 设置读超时和 Pong 处理
	session.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	session.Conn.SetPongHandler(func(string) error {
		session.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		messageType, message, err := session.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("WebSocket 读取错误", zap.Error(err))
			}
			break
		}

		// 只处理文本消息
		if messageType != websocket.TextMessage {
			continue
		}

		// 空消息跳过
		// Java: if (message.getPayloadLength() == 0) return;
		if len(message) == 0 {
			continue
		}

		// ping 心跳消息，直接返回 pong
		// Java: if (message.getPayloadLength() == 4 && Objects.equals(message.getPayload(), "ping"))
		if string(message) == "ping" {
			_ = session.SendText("pong")
			session.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			continue
		}

		// 解析 JSON 消息
		// Java: JsonWebSocketMessage jsonMessage = JsonUtils.parseObject(message.getPayload(), JsonWebSocketMessage.class);
		wsMsg, err := ws.ParseMessage(message)
		if err != nil {
			h.logger.Error("WebSocket 消息解析失败", zap.Error(err), zap.String("message", string(message)))
			continue
		}

		if wsMsg.Type == "" {
			h.logger.Warn("WebSocket 消息类型为空", zap.String("message", string(message)))
			continue
		}

		// TODO: 根据 Type 路由到不同的 MessageListener
		// Java: WebSocketMessageListener<Object> messageListener = listeners.get(jsonMessage.getType());
		// 目前仅记录日志，后续可扩展消息监听器框架
		h.logger.Info("收到 WebSocket 消息",
			zap.String("sessionId", session.ID),
			zap.String("type", wsMsg.Type),
			zap.String("content", wsMsg.Content),
		)
	}
}

// GetManager 获取 Manager (供其他服务调用推送消息)
func (h *WebSocketHandler) GetManager() *ws.Manager {
	return h.manager
}

// SendToUser 发送消息给指定用户
// Java: WebSocketMessageSender.send(userType, userId, messageType, messageContent)
func (h *WebSocketHandler) SendToUser(userID int64, messageType string, content interface{}) error {
	msg, err := ws.NewMessage(messageType, content)
	if err != nil {
		return err
	}
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}
	h.manager.Send(userID, data)
	return nil
}

// Broadcast 广播消息给所有用户
// Java: WebSocketMessageSender.send(userType, messageType, messageContent)
func (h *WebSocketHandler) Broadcast(messageType string, content interface{}) error {
	msg, err := ws.NewMessage(messageType, content)
	if err != nil {
		return err
	}
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}
	h.manager.Broadcast(data)
	return nil
}

// BroadcastByUserType 广播消息给指定用户类型
func (h *WebSocketHandler) BroadcastByUserType(userType int, messageType string, content interface{}) error {
	msg, err := ws.NewMessage(messageType, content)
	if err != nil {
		return err
	}
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}
	h.manager.BroadcastByUserType(userType, data)
	return nil
}
