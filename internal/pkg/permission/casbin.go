package permission

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"gorm.io/gorm"
)

// CasbinRBACModel 定义 RBAC 模型
const CasbinRBACModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

// InitEnforcer 初始化 Casbin Enforcer
func InitEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	// 1. 加载模型
	m, err := model.NewModelFromString(CasbinRBACModel)
	if err != nil {
		return nil, err
	}

	// 2. 创建适配器
	adapter := NewAdapter(db)

	// 3. 创建 Enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 4. 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		log.Printf("Failed to load casbin policy: %v", err)
		return nil, err
	}

	// 5. 开启自动加载 (可选，根据业务需求，这里先不开启，手动 Reload)
	// enforcer.StartAutoLoadPolicy(5 * time.Second)

	return enforcer, nil
}
