package datascope

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PluginRegistered 是一个空标记类型,表示数据权限插件已注册
type PluginRegistered struct{}

// RegisterPlugin 注册数据权限插件到GORM (用于Wire依赖注入)
// 返回PluginRegistered标记以便Wire知道插件已注册
func RegisterPlugin(
	db *gorm.DB,
	logger *zap.Logger,
	permissionSvc *system.PermissionService,
	deptSvc *system.DeptService,
) (*PluginRegistered, error) {
	plugin := NewPlugin(logger, permissionSvc, deptSvc)
	if err := db.Use(plugin); err != nil {
		return nil, err
	}
	logger.Info("Data scope plugin registered successfully")
	return &PluginRegistered{}, nil
}
