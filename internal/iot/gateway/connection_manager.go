package gateway

import (
	"log"
	"sync"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
)

// ConnectionInfo 连接信息
type ConnectionInfo struct {
	DeviceID      int64
	ProductKey    string
	DeviceName    string
	ClientID      string
	Authenticated bool
	RemoteAddress string
	ConnectedAt   time.Time
	LastHeartbeat time.Time
}

// ConnectionManager 设备连接管理器
// 管理设备连接状态、心跳维持及在线状态同步
type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*ConnectionInfo // key: clientID
	deviceIndex map[int64]string           // key: deviceID, value: clientID
	messageBus  core.MessageBus
	serverID    string
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(messageBus core.MessageBus, serverID string) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*ConnectionInfo),
		deviceIndex: make(map[int64]string),
		messageBus:  messageBus,
		serverID:    serverID,
	}
}

// RegisterConnection 注册连接
func (m *ConnectionManager) RegisterConnection(clientID string, info *ConnectionInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	info.ConnectedAt = time.Now()
	info.LastHeartbeat = time.Now()
	m.connections[clientID] = info
	m.deviceIndex[info.DeviceID] = clientID

	log.Printf("[ConnectionManager] Device connected: clientID=%s, deviceId=%d", clientID, info.DeviceID)

	// 发送设备上线消息
	m.sendStateMessage(info, "online")
}

// UnregisterConnection 注销连接
func (m *ConnectionManager) UnregisterConnection(clientID string) {
	m.mu.Lock()
	info := m.connections[clientID]
	if info != nil {
		delete(m.connections, clientID)
		delete(m.deviceIndex, info.DeviceID)
	}
	m.mu.Unlock()

	if info != nil {
		log.Printf("[ConnectionManager] Device disconnected: clientID=%s, deviceId=%d", clientID, info.DeviceID)
		// 发送设备离线消息
		m.sendStateMessage(info, "offline")
	}
}

// UpdateHeartbeat 更新心跳时间
func (m *ConnectionManager) UpdateHeartbeat(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if info, ok := m.connections[clientID]; ok {
		info.LastHeartbeat = time.Now()
	}
}

// GetConnectionInfo 获取连接信息
func (m *ConnectionManager) GetConnectionInfo(clientID string) *ConnectionInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connections[clientID]
}

// GetConnectionByDeviceID 根据设备 ID 获取连接信息
func (m *ConnectionManager) GetConnectionByDeviceID(deviceID int64) *ConnectionInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if clientID, ok := m.deviceIndex[deviceID]; ok {
		return m.connections[clientID]
	}
	return nil
}

// IsDeviceOnline 检查设备是否在线
func (m *ConnectionManager) IsDeviceOnline(deviceID int64) bool {
	return m.GetConnectionByDeviceID(deviceID) != nil
}

// GetOnlineDeviceCount 获取在线设备数量
func (m *ConnectionManager) GetOnlineDeviceCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

// sendStateMessage 发送状态消息到消息总线
func (m *ConnectionManager) sendStateMessage(info *ConnectionInfo, state string) {
	message := &core.IotDeviceMessage{
		Method:     "thing.lifecycle.state.update",
		Params:     map[string]any{"state": state},
		DeviceID:   info.DeviceID,
		ServerID:   m.serverID,
		ReportTime: time.Now(),
	}
	m.messageBus.Post(core.DeviceMessageTopic, message)
}

// StartHeartbeatChecker 启动心跳检查器
// 定期检查连接状态，超时则断开连接
func (m *ConnectionManager) StartHeartbeatChecker(timeout time.Duration, checkInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.checkHeartbeats(timeout)
		}
	}()
	log.Printf("[ConnectionManager] Heartbeat checker started: timeout=%v, interval=%v", timeout, checkInterval)
}

// checkHeartbeats 检查心跳超时
func (m *ConnectionManager) checkHeartbeats(timeout time.Duration) {
	m.mu.RLock()
	var expiredClients []string
	now := time.Now()
	for clientID, info := range m.connections {
		if now.Sub(info.LastHeartbeat) > timeout {
			expiredClients = append(expiredClients, clientID)
		}
	}
	m.mu.RUnlock()

	// 断开超时连接
	for _, clientID := range expiredClients {
		log.Printf("[ConnectionManager] Heartbeat timeout: clientID=%s", clientID)
		m.UnregisterConnection(clientID)
	}
}
