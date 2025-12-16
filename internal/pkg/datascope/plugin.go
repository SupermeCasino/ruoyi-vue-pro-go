package datascope

import (
	"context"
	"fmt"

	pkgcontext "github.com/wxlbd/ruoyi-mall-go/pkg/context"

	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Plugin GORM数据权限插件
type Plugin struct {
	logger        *zap.Logger
	permissionSvc *service.PermissionService
	deptSvc       *service.DeptService
}

// NewPlugin 创建数据权限插件
func NewPlugin(
	logger *zap.Logger,
	permissionSvc *service.PermissionService,
	deptSvc *service.DeptService,
) *Plugin {
	return &Plugin{
		logger:        logger,
		permissionSvc: permissionSvc,
		deptSvc:       deptSvc,
	}
}

// Name 返回插件名称
func (p *Plugin) Name() string {
	return "datascope"
}

// Initialize 初始化插件，注册GORM回调
func (p *Plugin) Initialize(db *gorm.DB) error {
	// 在查询前检查并应用数据权限
	return db.Callback().Query().Before("gorm:query").
		Register("datascope:before_query", p.beforeQuery)
}

// beforeQuery 在查询前应用数据权限过滤
func (p *Plugin) beforeQuery(db *gorm.DB) {
	ctx := db.Statement.Context
	if ctx == nil {
		return
	}

	// 1. 检查是否跳过数据权限
	if ShouldSkipDataScope(ctx) {
		p.logger.Debug("Skipping data scope check (context flag)")
		return
	}

	// 2. 获取当前登录用户
	loginUser := pkgcontext.GetLoginUserFromContext(ctx)
	if loginUser == nil {
		// 没有登录用户,不应用数据权限(可能是公开API)
		return
	}

	// 3. 为权限查询创建跳过数据权限的context,避免无限递归
	// beforeQuery -> IsSuperAdmin -> GetUserRoleIdListByUserId -> beforeQuery (递归)
	skipCtx := SkipDataScope(ctx)

	// 4. 超级管理员跳过数据权限
	isSuperAdmin, err := p.permissionSvc.IsSuperAdmin(skipCtx, loginUser.UserID)
	if err != nil {
		p.logger.Error("Failed to check super admin", zap.Error(err))
		return
	}
	if isSuperAdmin {
		p.logger.Debug("Skipping data scope check (super admin)",
			zap.Int64("user_id", loginUser.UserID))
		return
	}

	// 5. 获取用户的所有角色
	roleIDs, err := p.permissionSvc.GetUserRoleIdListByUserId(skipCtx, loginUser.UserID)
	if err != nil {
		p.logger.Error("Failed to get user roles", zap.Error(err), zap.Int64("user_id", loginUser.UserID))
		return
	}

	if len(roleIDs) == 0 {
		// 用户没有任何角色,默认只能看到自己创建的数据
		p.applySelfScope(db, loginUser.UserID)
		return
	}

	// 6. 计算最宽松的数据范围(数字越小权限越大)
	// 如果任一角色有全部数据权限,则不限制
	dataScope, customDeptIDs, err := p.calculateDataScope(skipCtx, roleIDs)
	if err != nil {
		p.logger.Error("Failed to calculate data scope", zap.Error(err))
		return
	}

	// 7. 根据数据范围应用SQL条件
	p.applyDataScope(db, dataScope, customDeptIDs, loginUser)
}

// calculateDataScope 计算用户的最宽松数据范围
func (p *Plugin) calculateDataScope(ctx context.Context, roleIDs []int64) (DataScope, []int64, error) {
	maxScope := DataScopeSelf // 默认最严格
	var customDeptIDs []int64

	for _, roleID := range roleIDs {
		role, err := p.permissionSvc.GetRoleById(ctx, roleID)
		if err != nil {
			return maxScope, nil, err
		}

		roleScope := DataScope(role.DataScope)

		// 取最宽松的权限
		if roleScope < maxScope {
			maxScope = roleScope
			// 如果是自定义部门，收集部门ID
			if roleScope == DataScopeDeptCustom {
				deptIDs, err := p.permissionSvc.GetRoleDeptIdListByRoleId(ctx, roleID)
				if err != nil {
					return maxScope, nil, err
				}
				customDeptIDs = append(customDeptIDs, deptIDs...)
			}
		} else if roleScope == DataScopeDeptCustom && maxScope == DataScopeDeptCustom {
			// 多个角色都是自定义部门，合并部门ID
			deptIDs, err := p.permissionSvc.GetRoleDeptIdListByRoleId(ctx, roleID)
			if err != nil {
				return maxScope, nil, err
			}
			customDeptIDs = append(customDeptIDs, deptIDs...)
		}
	}

	return maxScope, customDeptIDs, nil
}

// applyDataScope 应用数据范围SQL条件
func (p *Plugin) applyDataScope(db *gorm.DB, scope DataScope, customDeptIDs []int64, user *pkgcontext.LoginUser) {
	switch scope {
	case DataScopeAll:
		// 全部数据权限，不添加任何条件
		p.logger.Debug("Data scope: ALL (no restriction)",
			zap.Int64("user_id", user.UserID))
		return

	case DataScopeDeptCustom:
		// 自定义部门数据权限
		if len(customDeptIDs) > 0 {
			p.applyDeptScope(db, customDeptIDs, "DEPT_CUSTOM", user.UserID)
		} else {
			// 没有自定义部门，降级为仅本人
			p.applySelfScope(db, user.UserID)
		}

	case DataScopeDeptOnly:
		// 本部门数据权限
		if user.DeptID != nil && *user.DeptID > 0 {
			p.applyDeptScope(db, []int64{*user.DeptID}, "DEPT_ONLY", user.UserID)
		} else {
			// 用户没有部门，降级为仅本人
			p.applySelfScope(db, user.UserID)
		}

	case DataScopeDeptAndChild:
		// 本部门及子部门数据权限
		if user.DeptID != nil && *user.DeptID > 0 {
			deptIDs, err := p.deptSvc.GetDeptIdListByParentId(context.Background(), *user.DeptID)
			if err != nil {
				p.logger.Error("Failed to get child dept IDs", zap.Error(err))
				p.applySelfScope(db, user.UserID)
				return
			}
			// 包含当前部门
			allDeptIDs := append([]int64{*user.DeptID}, deptIDs...)
			p.applyDeptScope(db, allDeptIDs, "DEPT_AND_CHILD", user.UserID)
		} else {
			p.applySelfScope(db, user.UserID)
		}

	case DataScopeSelf:
		// 仅本人数据权限
		p.applySelfScope(db, user.UserID)

	default:
		p.logger.Warn("Unknown data scope, applying SELF", zap.Int("scope", int(scope)))
		p.applySelfScope(db, user.UserID)
	}
}

// applyDeptScope 应用部门范围条件
func (p *Plugin) applyDeptScope(db *gorm.DB, deptIDs []int64, scopeType string, userID int64) {
	if p.hasColumn(db, "dept_id") {
		db.Where("dept_id IN ?", deptIDs)
		p.logger.Debug(fmt.Sprintf("Applied data scope: %s", scopeType),
			zap.Int64("user_id", userID),
			zap.Int64s("dept_ids", deptIDs))
	} else {
		// 表没有dept_id字段，降级为creator过滤
		p.applySelfScope(db, userID)
	}
}

// applySelfScope 应用仅本人范围条件
func (p *Plugin) applySelfScope(db *gorm.DB, userID int64) {
	if p.hasColumn(db, "creator") {
		db.Where("creator = ?", userID)
		p.logger.Debug("Applied data scope: SELF",
			zap.Int64("user_id", userID))
	} else {
		p.logger.Debug("No creator column, skipping data scope",
			zap.String("table", db.Statement.Table))
	}
}

// hasColumn 检查表是否有指定列
func (p *Plugin) hasColumn(db *gorm.DB, columnName string) bool {
	if db.Statement.Schema == nil {
		return false
	}
	_, ok := db.Statement.Schema.FieldsByDBName[columnName]
	return ok
}
