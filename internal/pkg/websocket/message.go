package websocket

import (
	"encoding/json"
)

// Message WebSocket JSON 消息格式 (与 Java JsonWebSocketMessage 对齐)
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// NewMessage 创建新消息
func NewMessage(msgType string, content interface{}) (*Message, error) {
	var contentStr string
	switch v := content.(type) {
	case string:
		contentStr = v
	default:
		data, err := json.Marshal(content)
		if err != nil {
			return nil, err
		}
		contentStr = string(data)
	}
	return &Message{
		Type:    msgType,
		Content: contentStr,
	}, nil
}

// ToJSON 序列化为 JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// ParseMessage 解析 JSON 消息
func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
