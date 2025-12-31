package router

import (
	memberAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterMemberRoutes 注册会员管理模块路由
func RegisterMemberRoutes(engine *gin.Engine,
	handlers *memberAdmin.Handlers,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	api := engine.Group("/admin-api")
	api.Use(middleware.Auth())

	// Member Point Record
	pointRecordGroup := api.Group("/member/point/record")
	{
		pointRecordGroup.GET("/page", handlers.PointRecord.GetPointRecordPage)
	}

	// Member Sign-in Config
	signInConfigGroup := api.Group("/member/sign-in/config")
	{
		signInConfigGroup.POST("/create", handlers.SignInConfig.CreateSignInConfig)
		signInConfigGroup.PUT("/update", handlers.SignInConfig.UpdateSignInConfig)
		signInConfigGroup.DELETE("/delete", handlers.SignInConfig.DeleteSignInConfig)
		signInConfigGroup.GET("/get", handlers.SignInConfig.GetSignInConfig)
		signInConfigGroup.GET("/list", handlers.SignInConfig.GetSignInConfigList)
	}

	// Member Sign-in Record (Admin)
	signInRecordGroup := api.Group("/member/sign-in/record")
	{
		signInRecordGroup.GET("/page", handlers.SignInRecord.GetSignInRecordPage)
	}

	// Member Config 会员配置
	configGroup := api.Group("/member/config")
	{
		configGroup.PUT("/save", casbinMiddleware.RequirePermission("member:config:save"), handlers.Config.SaveConfig)
		configGroup.GET("/get", casbinMiddleware.RequirePermission("member:config:query"), handlers.Config.GetConfig)
	}

	// Member Group 用户分组
	groupGroup := api.Group("/member/group")
	{
		groupGroup.POST("/create", casbinMiddleware.RequirePermission("member:group:create"), handlers.Group.CreateGroup)
		groupGroup.PUT("/update", casbinMiddleware.RequirePermission("member:group:update"), handlers.Group.UpdateGroup)
		groupGroup.DELETE("/delete", casbinMiddleware.RequirePermission("member:group:delete"), handlers.Group.DeleteGroup)
		groupGroup.GET("/get", casbinMiddleware.RequirePermission("member:group:query"), handlers.Group.GetGroup)
		groupGroup.GET("/list-all-simple", handlers.Group.GetSimpleGroupList)
		groupGroup.GET("/page", casbinMiddleware.RequirePermission("member:group:query"), handlers.Group.GetGroupPage)
	}

	// Member Level 会员等级
	levelGroup := api.Group("/member/level")
	{
		levelGroup.POST("/create", casbinMiddleware.RequirePermission("member:level:create"), handlers.Level.CreateLevel)
		levelGroup.PUT("/update", casbinMiddleware.RequirePermission("member:level:update"), handlers.Level.UpdateLevel)
		levelGroup.DELETE("/delete", casbinMiddleware.RequirePermission("member:level:delete"), handlers.Level.DeleteLevel)
		levelGroup.GET("/get", casbinMiddleware.RequirePermission("member:level:query"), handlers.Level.GetLevel)
		levelGroup.GET("/list-all-simple", handlers.Level.GetLevelListSimple)
		levelGroup.GET("/list", handlers.Level.GetLevelListSimple)
	}

	// Member Tag 会员标签
	tagGroup := api.Group("/member/tag")
	{
		tagGroup.POST("/create", casbinMiddleware.RequirePermission("member:tag:create"), handlers.Tag.CreateTag)
		tagGroup.PUT("/update", casbinMiddleware.RequirePermission("member:tag:update"), handlers.Tag.UpdateTag)
		tagGroup.DELETE("/delete", casbinMiddleware.RequirePermission("member:tag:delete"), handlers.Tag.DeleteTag)
		tagGroup.GET("/get", casbinMiddleware.RequirePermission("member:tag:query"), handlers.Tag.GetTag)
		tagGroup.GET("/list-all-simple", handlers.Tag.GetSimpleTagList)
		tagGroup.GET("/page", casbinMiddleware.RequirePermission("member:tag:query"), handlers.Tag.GetTagPage)
	}

	// Member User 会员用户
	userGroup := api.Group("/member/user")
	{
		userGroup.PUT("/update", casbinMiddleware.RequirePermission("member:user:update"), handlers.User.UpdateUser)
		userGroup.PUT("/update-level", casbinMiddleware.RequirePermission("member:user:update-level"), handlers.User.UpdateUserLevel)
		userGroup.PUT("/update-point", casbinMiddleware.RequirePermission("member:user:update-point"), handlers.User.UpdateUserPoint)
		userGroup.GET("/get", casbinMiddleware.RequirePermission("member:user:query"), handlers.User.GetUser)
		userGroup.GET("/page", casbinMiddleware.RequirePermission("member:user:query"), handlers.User.GetUserPage)
	}
}
