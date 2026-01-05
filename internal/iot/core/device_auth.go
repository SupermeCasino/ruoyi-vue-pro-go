package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// DeviceAuthUtils 设备认证工具
type DeviceAuthUtils struct{}

// NewDeviceAuthUtils 创建设备认证工具
func NewDeviceAuthUtils() *DeviceAuthUtils {
	return &DeviceAuthUtils{}
}

// DeviceInfo 设备信息（从用户名解析）
type DeviceInfo struct {
	ProductKey string
	DeviceName string
}

// ParseUsername 解析 MQTT 用户名获取设备信息
// 用户名格式: {productKey}&{deviceName}
func (u *DeviceAuthUtils) ParseUsername(username string) (*DeviceInfo, error) {
	parts := strings.Split(username, "&")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid username format: %s", username)
	}
	return &DeviceInfo{
		ProductKey: parts[0],
		DeviceName: parts[1],
	}, nil
}

// BuildPassword 构建 MQTT 认证密码
// password = HMAC-SHA256(deviceSecret, content)
func (u *DeviceAuthUtils) BuildPassword(deviceSecret, content string) string {
	h := hmac.New(sha256.New, []byte(deviceSecret))
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}

// BuildAuthContent 构建认证内容
// content = deviceName{deviceName}productKey{productKey}
func (u *DeviceAuthUtils) BuildAuthContent(deviceName, productKey string) string {
	return fmt.Sprintf("deviceName%sproductKey%s", deviceName, productKey)
}

// ValidatePassword 验证密码
func (u *DeviceAuthUtils) ValidatePassword(deviceSecret, deviceName, productKey, password string) bool {
	content := u.BuildAuthContent(deviceName, productKey)
	expectedPassword := u.BuildPassword(deviceSecret, content)
	return expectedPassword == password
}
