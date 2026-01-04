package websocket

import (
	"encoding/json"
)

// Message WebSocket JSON 消息格式 (与 Java JsonWebSocketMessage 对齐)
// Content 支持任何 JSON 可序列化的对象（对象、数组、字符串等）
// 序列化时，content 将作为完整的 JSON 对象包含在消息中，而不是字符串
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// NewMessage 创建新消息
// content 可以是任何类型（对象、数组、字符串等），将直接作为 JSON 对象包含在消息中
func NewMessage(msgType string, content interface{}) (*Message, error) {
	return &Message{
		Type:    msgType,
		Content: content,
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
