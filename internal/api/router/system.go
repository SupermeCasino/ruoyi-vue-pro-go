package router

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes 注册系统管理模块路由
func RegisterSystemRoutes(engine *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tenantHandler *handler.TenantHandler,
	tenantPackageHandler *handler.TenantPackageHandler, // 新增租户套餐
	dictHandler *handler.DictHandler,
	deptHandler *handler.DeptHandler,
	postHandler *handler.PostHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	permissionHandler *handler.PermissionHandler,
	noticeHandler *handler.NoticeHandler,
	loginLogHandler *handler.LoginLogHandler,
	operateLogHandler *handler.OperateLogHandler,
	configHandler *handler.ConfigHandler,
	smsChannelHandler *handler.SmsChannelHandler,
	smsTemplateHandler *handler.SmsTemplateHandler,
	smsLogHandler *handler.SmsLogHandler,
	fileConfigHandler *handler.FileConfigHandler,
	fileHandler *handler.FileHandler,
	jobHandler *handler.JobHandler,
	jobLogHandler *handler.JobLogHandler,
	apiAccessLogHandler *handler.ApiAccessLogHandler,
	apiErrorLogHandler *handler.ApiErrorLogHandler,
	socialClientHandler *handler.SocialClientHandler,
	socialUserHandler *handler.SocialUserHandler,
	mailHandler *handler.MailHandler,
	notifyHandler *handler.NotifyHandler,
	oauth2ClientHandler *handler.OAuth2ClientHandler,
	webSocketHandler *handler.WebSocketHandler,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	api := engine.Group("/admin-api")
	{
		systemGroup := api.Group("/system")
		{
			// ====== Public Routes (No Auth Required) ======
			// Auth Public Routes
			authGroup := systemGroup.Group("/auth")
			{
				authGroup.POST("/login", authHandler.Login)
				authGroup.POST("/logout", authHandler.Logout)
				authGroup.POST("/refresh-token", authHandler.RefreshToken)
				authGroup.POST("/sms-login", authHandler.SmsLogin)
				authGroup.POST("/send-sms-code", authHandler.SendSmsCode)
				authGroup.POST("/register", authHandler.Register)
				authGroup.POST("/reset-password", authHandler.ResetPassword)
				authGroup.GET("/social-auth-redirect", authHandler.SocialAuthRedirect)
				authGroup.POST("/social-login", authHandler.SocialLogin)
			}

			// Tenant Public Routes
			tenantPublicGroup := systemGroup.Group("/tenant")
			{
				tenantPublicGroup.GET("/simple-list", tenantHandler.GetTenantSimpleList)
				tenantPublicGroup.GET("/get-by-website", tenantHandler.GetTenantByWebsite)
				tenantPublicGroup.GET("/get-id-by-name", tenantHandler.GetTenantIdByName)
			}

			// Dict Public Routes
			dictTypePublicGroup := systemGroup.Group("/dict-type")
			{
				dictTypePublicGroup.GET("/simple-list", dictHandler.GetSimpleDictTypeList)
			}

			dictDataPublicGroup := systemGroup.Group("/dict-data")
			{
				dictDataPublicGroup.GET("/simple-list", dictHandler.GetSimpleDictDataList)
				dictDataPublicGroup.GET("/list-all-simple", dictHandler.GetSimpleDictDataList)
			}

			// Dept Public Routes
			deptPublicGroup := systemGroup.Group("/dept")
			{
				deptPublicGroup.GET("/list", deptHandler.GetDeptList)
				deptPublicGroup.GET("/list-all-simple", deptHandler.GetSimpleDeptList)
				deptPublicGroup.GET("/simple-list", deptHandler.GetSimpleDeptList)
			}

			// Post Public Routes
			postPublicGroup := systemGroup.Group("/post")
			{
				postPublicGroup.GET("/simple-list", postHandler.GetSimplePostList)
			}

			// User Public Routes
			userPublicGroup := systemGroup.Group("/user")
			{
				userPublicGroup.GET("/list-all-simple", userHandler.GetSimpleUserList)
				userPublicGroup.GET("/simple-list", userHandler.GetSimpleUserList)
			}

			// Role Public Routes
			rolePublicGroup := systemGroup.Group("/role")
			{
				rolePublicGroup.GET("/list-all-simple", roleHandler.GetSimpleRoleList)
				rolePublicGroup.GET("/simple-list", roleHandler.GetSimpleRoleList)
			}

			// Menu Public Routes
			menuPublicGroup := systemGroup.Group("/menu")
			{
				menuPublicGroup.GET("/simple-list", menuHandler.GetSimpleMenuList)
			}

			// ====== Protected Routes (Auth Required) ======
			// Apply Auth Middleware to all subsequent system routes
			systemGroup.Use(middleware.Auth())

			// Auth Protected Routes
			authProtectedGroup := systemGroup.Group("/auth")
			{
				authProtectedGroup.GET("/get-permission-info", authHandler.GetPermissionInfo)
			}

			// Tenant Protected Routes
			tenantProtectedGroup := systemGroup.Group("/tenant")
			{
				tenantProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:tenant:create"), tenantHandler.CreateTenant)
				tenantProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:tenant:update"), tenantHandler.UpdateTenant)
				tenantProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:tenant:delete"), tenantHandler.DeleteTenant)
				tenantProtectedGroup.DELETE("/delete-list", casbinMiddleware.RequirePermission("system:tenant:delete"), tenantHandler.DeleteTenantList)
				tenantProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:tenant:query"), tenantHandler.GetTenant)
				tenantProtectedGroup.GET("/page", casbinMiddleware.RequirePermission("system:tenant:query"), tenantHandler.GetTenantPage)
				tenantProtectedGroup.GET("/export-excel", casbinMiddleware.RequirePermission("system:tenant:export"), tenantHandler.ExportTenantExcel)
			}

			// Tenant Package Protected Routes
			tenantPackageGroup := systemGroup.Group("/tenant-package")
			{
				tenantPackageGroup.POST("/create", casbinMiddleware.RequirePermission("system:tenant-package:create"), tenantPackageHandler.CreateTenantPackage)
				tenantPackageGroup.PUT("/update", casbinMiddleware.RequirePermission("system:tenant-package:update"), tenantPackageHandler.UpdateTenantPackage)
				tenantPackageGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:tenant-package:delete"), tenantPackageHandler.DeleteTenantPackage)
				tenantPackageGroup.DELETE("/delete-list", casbinMiddleware.RequirePermission("system:tenant-package:delete"), tenantPackageHandler.DeleteTenantPackageList)
				tenantPackageGroup.GET("/get", casbinMiddleware.RequirePermission("system:tenant-package:query"), tenantPackageHandler.GetTenantPackage)
				tenantPackageGroup.GET("/page", casbinMiddleware.RequirePermission("system:tenant-package:query"), tenantPackageHandler.GetTenantPackagePage)
				tenantPackageGroup.GET("/get-simple-list", tenantPackageHandler.GetTenantPackageSimpleList)
				tenantPackageGroup.GET("/simple-list", tenantPackageHandler.GetTenantPackageSimpleList)
			}

			// Dict Type Protected Routes
			dictTypeProtectedGroup := systemGroup.Group("/dict-type")
			{
				dictTypeProtectedGroup.GET("/page", casbinMiddleware.RequirePermission("system:dict:query"), dictHandler.GetDictTypePage)
				dictTypeProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:dict:query"), dictHandler.GetDictType)
				dictTypeProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:dict:create"), dictHandler.CreateDictType)
				dictTypeProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:dict:update"), dictHandler.UpdateDictType)
				dictTypeProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:dict:delete"), dictHandler.DeleteDictType)
				dictTypeProtectedGroup.GET("/export-excel", casbinMiddleware.RequirePermission("system:dict:export"), dictHandler.ExportDictTypeExcel)
			}

			// Dict Data Protected Routes
			dictDataProtectedGroup := systemGroup.Group("/dict-data")
			{
				dictDataProtectedGroup.GET("/page", casbinMiddleware.RequirePermission("system:dict:query"), dictHandler.GetDictDataPage)
				dictDataProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:dict:query"), dictHandler.GetDictData)
				dictDataProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:dict:create"), dictHandler.CreateDictData)
				dictDataProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:dict:update"), dictHandler.UpdateDictData)
				dictDataProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:dict:delete"), dictHandler.DeleteDictData)
			}

			// Dept Protected Routes
			deptProtectedGroup := systemGroup.Group("/dept")
			{
				deptProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:dept:query"), deptHandler.GetDept)
				deptProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:dept:create"), deptHandler.CreateDept)
				deptProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:dept:update"), deptHandler.UpdateDept)
				deptProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:dept:delete"), deptHandler.DeleteDept)
			}

			// Post Protected Routes
			postProtectedGroup := systemGroup.Group("/post")
			{
				postProtectedGroup.GET("/page", casbinMiddleware.RequirePermission("system:post:query"), postHandler.GetPostPage)
				postProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:post:query"), postHandler.GetPost)
				postProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:post:create"), postHandler.CreatePost)
				postProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:post:update"), postHandler.UpdatePost)
				postProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:post:delete"), postHandler.DeletePost)
			}

			// User Protected Routes
			userProtectedGroup := systemGroup.Group("/user")
			{
				userProtectedGroup.GET("/page", casbinMiddleware.RequirePermission("system:user:query"), userHandler.GetUserPage)
				userProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:user:query"), userHandler.GetUser)
				userProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:user:create"), userHandler.CreateUser)
				userProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:user:update"), userHandler.UpdateUser)
				userProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:user:delete"), userHandler.DeleteUser)
				userProtectedGroup.DELETE("/delete-list", casbinMiddleware.RequirePermission("system:user:delete"), userHandler.DeleteUserList)
				userProtectedGroup.PUT("/update-status", casbinMiddleware.RequirePermission("system:user:update"), userHandler.UpdateUserStatus)
				userProtectedGroup.PUT("/update-password", casbinMiddleware.RequirePermission("system:user:update-password"), userHandler.UpdateUserPassword)
				userProtectedGroup.GET("/export", casbinMiddleware.RequirePermission("system:user:export"), userHandler.ExportUser)
				userProtectedGroup.GET("/get-import-template", casbinMiddleware.RequirePermission("system:user:import"), userHandler.GetImportTemplate)
				userProtectedGroup.POST("/import", casbinMiddleware.RequirePermission("system:user:import"), userHandler.ImportUser)
			}

			// Role Protected Routes
			roleProtectedGroup := systemGroup.Group("/role")
			{
				roleProtectedGroup.GET("/page", casbinMiddleware.RequirePermission("system:role:query"), roleHandler.GetRolePage)
				roleProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:role:query"), roleHandler.GetRole)
				roleProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:role:create"), roleHandler.CreateRole)
				roleProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:role:update"), roleHandler.UpdateRole)
				roleProtectedGroup.PUT("/update-status", casbinMiddleware.RequirePermission("system:role:update"), roleHandler.UpdateRoleStatus)
				roleProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:role:delete"), roleHandler.DeleteRole)
			}

			// Permission Protected Routes
			permProtectedGroup := systemGroup.Group("/permission")
			{
				permProtectedGroup.GET("/list-role-menus", casbinMiddleware.RequirePermission("system:permission:assign-role-menu"), permissionHandler.GetRoleMenuList)
				permProtectedGroup.POST("/assign-role-menu", casbinMiddleware.RequirePermission("system:permission:assign-role-menu"), permissionHandler.AssignRoleMenu)
				permProtectedGroup.POST("/assign-role-data-scope", casbinMiddleware.RequirePermission("system:permission:assign-role-data-scope"), permissionHandler.AssignRoleDataScope)
				permProtectedGroup.GET("/list-user-roles", casbinMiddleware.RequirePermission("system:permission:assign-user-role"), permissionHandler.GetUserRoleList)
				permProtectedGroup.POST("/assign-user-role", casbinMiddleware.RequirePermission("system:permission:assign-user-role"), permissionHandler.AssignUserRole)
			}

			// Menu Protected Routes
			menuProtectedGroup := systemGroup.Group("/menu")
			{
				menuProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:menu:create"), menuHandler.CreateMenu)
				menuProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:menu:update"), menuHandler.UpdateMenu)
				menuProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:menu:delete"), menuHandler.DeleteMenu)
				menuProtectedGroup.GET("/list", casbinMiddleware.RequirePermission("system:menu:query"), menuHandler.GetMenuList)
				menuProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:menu:query"), menuHandler.GetMenu)
			}

			// Notice
			noticeGroup := systemGroup.Group("/notice")
			{
				noticeGroup.GET("/page", casbinMiddleware.RequirePermission("system:notice:query"), noticeHandler.GetNoticePage)
				noticeGroup.GET("/get", casbinMiddleware.RequirePermission("system:notice:query"), noticeHandler.GetNotice)
				noticeGroup.POST("/create", casbinMiddleware.RequirePermission("system:notice:create"), noticeHandler.CreateNotice)
				noticeGroup.PUT("/update", casbinMiddleware.RequirePermission("system:notice:update"), noticeHandler.UpdateNotice)
				noticeGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:notice:delete"), noticeHandler.DeleteNotice)
				noticeGroup.POST("/push", casbinMiddleware.RequirePermission("system:notice:create"), noticeHandler.Push)
			}

			// Login Log
			loginLogGroup := systemGroup.Group("/login-log")
			{
				loginLogGroup.GET("/page", casbinMiddleware.RequirePermission("system:login-log:query"), loginLogHandler.GetLoginLogPage)
			}

			// Operate Log
			operateLogGroup := systemGroup.Group("/operate-log")
			{
				operateLogGroup.GET("/page", casbinMiddleware.RequirePermission("system:operate-log:query"), operateLogHandler.GetOperateLogPage)
			}

			// Mail Account
			mailAccountGroup := systemGroup.Group("/mail/account")
			{
				mailAccountGroup.POST("/create", casbinMiddleware.RequirePermission("system:mail-account:create"), mailHandler.CreateMailAccount)
				mailAccountGroup.PUT("/update", casbinMiddleware.RequirePermission("system:mail-account:update"), mailHandler.UpdateMailAccount)
				mailAccountGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:mail-account:delete"), mailHandler.DeleteMailAccount)
				mailAccountGroup.GET("/get", casbinMiddleware.RequirePermission("system:mail-account:query"), mailHandler.GetMailAccount)
				mailAccountGroup.GET("/page", casbinMiddleware.RequirePermission("system:mail-account:query"), mailHandler.GetMailAccountPage)
				mailAccountGroup.GET("/list-all-simple", mailHandler.GetSimpleMailAccountList)
			}

			// Mail Template
			mailTemplateGroup := systemGroup.Group("/mail/template")
			{
				mailTemplateGroup.POST("/create", casbinMiddleware.RequirePermission("system:mail-template:create"), mailHandler.CreateMailTemplate)
				mailTemplateGroup.PUT("/update", casbinMiddleware.RequirePermission("system:mail-template:update"), mailHandler.UpdateMailTemplate)
				mailTemplateGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:mail-template:delete"), mailHandler.DeleteMailTemplate)
				mailTemplateGroup.GET("/get", casbinMiddleware.RequirePermission("system:mail-template:query"), mailHandler.GetMailTemplate)
				mailTemplateGroup.GET("/page", casbinMiddleware.RequirePermission("system:mail-template:query"), mailHandler.GetMailTemplatePage)
				mailTemplateGroup.POST("/send-mail", casbinMiddleware.RequirePermission("system:mail-template:send-mail"), mailHandler.SendMail)
			}

			// Mail Log
			mailLogGroup := systemGroup.Group("/mail/log")
			{
				mailLogGroup.GET("/page", casbinMiddleware.RequirePermission("system:mail-log:query"), mailHandler.GetMailLogPage)
			}

			// Notify Template
			notifyTemplateGroup := systemGroup.Group("/notify-template")
			{
				notifyTemplateGroup.POST("/create", casbinMiddleware.RequirePermission("system:notify-template:create"), notifyHandler.CreateNotifyTemplate)
				notifyTemplateGroup.PUT("/update", casbinMiddleware.RequirePermission("system:notify-template:update"), notifyHandler.UpdateNotifyTemplate)
				notifyTemplateGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:notify-template:delete"), notifyHandler.DeleteNotifyTemplate)
				notifyTemplateGroup.GET("/get", casbinMiddleware.RequirePermission("system:notify-template:query"), notifyHandler.GetNotifyTemplate)
				notifyTemplateGroup.GET("/page", casbinMiddleware.RequirePermission("system:notify-template:query"), notifyHandler.GetNotifyTemplatePage)
				notifyTemplateGroup.POST("/send-notify", casbinMiddleware.RequirePermission("system:notify-template:send-notify"), notifyHandler.SendNotify)
			}

			// Notify Message
			notifyMessageGroup := systemGroup.Group("/notify-message")
			{
				notifyMessageGroup.GET("/get", casbinMiddleware.RequirePermission("system:notify-message:query"), notifyHandler.GetNotifyMessage)
				notifyMessageGroup.GET("/get-unread-count", notifyHandler.GetUnreadNotifyMessageCount)
				notifyMessageGroup.GET("/get-unread-list", notifyHandler.GetUnreadNotifyMessageList)
				notifyMessageGroup.GET("/my-page", notifyHandler.GetMyNotifyMessagePage)
				notifyMessageGroup.GET("/page", casbinMiddleware.RequirePermission("system:notify-message:query"), notifyHandler.GetNotifyMessagePage)
				notifyMessageGroup.PUT("/update-read", notifyHandler.UpdateNotifyMessageRead)
				notifyMessageGroup.PUT("/update-all-read", notifyHandler.UpdateAllNotifyMessageRead)
			}

			// OAuth2 Client
			oauth2ClientGroup := systemGroup.Group("/oauth2-client")
			{
				oauth2ClientGroup.POST("/create", oauth2ClientHandler.CreateOAuth2Client)
				oauth2ClientGroup.PUT("/update", oauth2ClientHandler.UpdateOAuth2Client)
				oauth2ClientGroup.DELETE("/delete", oauth2ClientHandler.DeleteOAuth2Client)
				oauth2ClientGroup.GET("/get", oauth2ClientHandler.GetOAuth2Client)
				oauth2ClientGroup.GET("/page", oauth2ClientHandler.GetOAuth2ClientPage)
			}

			// Social Client
			socialClientProtectedGroup := systemGroup.Group("/social-client")
			{
				socialClientProtectedGroup.POST("/create", casbinMiddleware.RequirePermission("system:social-client:create"), socialClientHandler.CreateSocialClient)
				socialClientProtectedGroup.PUT("/update", casbinMiddleware.RequirePermission("system:social-client:update"), socialClientHandler.UpdateSocialClient)
				socialClientProtectedGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:social-client:delete"), socialClientHandler.DeleteSocialClient)
				socialClientProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:social-client:query"), socialClientHandler.GetSocialClient)
				socialClientProtectedGroup.GET("/page", casbinMiddleware.RequirePermission("system:social-client:query"), socialClientHandler.GetSocialClientPage)
			}

			// Social User
			socialUserProtectedGroup := systemGroup.Group("/social-user")
			{
				socialUserProtectedGroup.POST("/bind", socialUserHandler.BindSocialUser)
				socialUserProtectedGroup.DELETE("/unbind", socialUserHandler.UnbindSocialUser)
				socialUserProtectedGroup.GET("/get-bind-list", socialUserHandler.GetSocialUserList)
				socialUserProtectedGroup.GET("/get", casbinMiddleware.RequirePermission("system:social-user:query"), socialUserHandler.GetSocialUser)
				socialUserProtectedGroup.GET("/page", casbinMiddleware.RequirePermission("system:social-user:query"), socialUserHandler.GetSocialUserPage)
			}
		}

		// ====== SMS Routes (Protected) ======
		// SMS simple-list is public
		smsChannelPublicGroup := api.Group("/system/sms-channel")
		{
			smsChannelPublicGroup.GET("/simple-list", smsChannelHandler.GetSimpleSmsChannelList)
		}

		// SMS Protected Routes
		smsChannelGroup := api.Group("/system/sms-channel", middleware.Auth())
		{
			smsChannelGroup.POST("/create", casbinMiddleware.RequirePermission("system:sms-channel:create"), smsChannelHandler.CreateSmsChannel)
			smsChannelGroup.PUT("/update", casbinMiddleware.RequirePermission("system:sms-channel:update"), smsChannelHandler.UpdateSmsChannel)
			smsChannelGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:sms-channel:delete"), smsChannelHandler.DeleteSmsChannel)
			smsChannelGroup.GET("/get", casbinMiddleware.RequirePermission("system:sms-channel:query"), smsChannelHandler.GetSmsChannel)
			smsChannelGroup.GET("/page", casbinMiddleware.RequirePermission("system:sms-channel:query"), smsChannelHandler.GetSmsChannelPage)
		}

		smsTemplateGroup := api.Group("/system/sms-template", middleware.Auth())
		{
			smsTemplateGroup.POST("/create", casbinMiddleware.RequirePermission("system:sms-template:create"), smsTemplateHandler.CreateSmsTemplate)
			smsTemplateGroup.PUT("/update", casbinMiddleware.RequirePermission("system:sms-template:update"), smsTemplateHandler.UpdateSmsTemplate)
			smsTemplateGroup.DELETE("/delete", casbinMiddleware.RequirePermission("system:sms-template:delete"), smsTemplateHandler.DeleteSmsTemplate)
			smsTemplateGroup.GET("/get", casbinMiddleware.RequirePermission("system:sms-template:query"), smsTemplateHandler.GetSmsTemplate)
			smsTemplateGroup.GET("/page", casbinMiddleware.RequirePermission("system:sms-template:query"), smsTemplateHandler.GetSmsTemplatePage)
			smsTemplateGroup.POST("/send-sms", casbinMiddleware.RequirePermission("system:sms-template:send-sms"), smsTemplateHandler.SendSms)
			smsTemplateGroup.GET("/export-excel", casbinMiddleware.RequirePermission("system:sms-template:export"), smsTemplateHandler.ExportSmsTemplateExcel)
		}

		smsLogGroup := api.Group("/system/sms-log", middleware.Auth())
		{
			smsLogGroup.GET("/page", casbinMiddleware.RequirePermission("system:sms-log:query"), smsLogHandler.GetSmsLogPage)
			smsLogGroup.GET("/export-excel", casbinMiddleware.RequirePermission("system:sms-log:export"), smsLogHandler.ExportSmsLogExcel)
		}

		// ====== Infra Routes (Public) ======
		infraPublicGroup := api.Group("/infra")
		{
			infraPublicGroup.GET("/file/:configId/get/*path", fileHandler.GetFileContent)
		}

		// ====== Infra Routes (Protected) ======
		configGroup := api.Group("/infra/config", middleware.Auth())
		{
			configGroup.GET("/page", casbinMiddleware.RequirePermission("infra:config:query"), configHandler.GetConfigPage)
			configGroup.GET("/get", casbinMiddleware.RequirePermission("infra:config:query"), configHandler.GetConfig)
			configGroup.GET("/get-value-by-key", casbinMiddleware.RequirePermission("infra:config:query"), configHandler.GetConfigKey)
			configGroup.POST("/create", casbinMiddleware.RequirePermission("infra:config:create"), configHandler.CreateConfig)
			configGroup.PUT("/update", casbinMiddleware.RequirePermission("infra:config:update"), configHandler.UpdateConfig)
			configGroup.DELETE("/delete", casbinMiddleware.RequirePermission("infra:config:delete"), configHandler.DeleteConfig)
		}

		infraGroup := api.Group("/infra", middleware.Auth())
		{
			// WebSocket (对齐 Java /infra/ws)
			infraGroup.GET("/ws", webSocketHandler.Handle)

			// File Config
			fileConfigGroup := infraGroup.Group("/file-config")
			{
				fileConfigGroup.POST("/create", casbinMiddleware.RequirePermission("infra:file-config:create"), fileConfigHandler.CreateFileConfig)
				fileConfigGroup.PUT("/update", casbinMiddleware.RequirePermission("infra:file-config:update"), fileConfigHandler.UpdateFileConfig)
				fileConfigGroup.PUT("/update-master", casbinMiddleware.RequirePermission("infra:file-config:update"), fileConfigHandler.UpdateFileConfigMaster)
				fileConfigGroup.DELETE("/delete", casbinMiddleware.RequirePermission("infra:file-config:delete"), fileConfigHandler.DeleteFileConfig)
				fileConfigGroup.GET("/page", casbinMiddleware.RequirePermission("infra:file-config:query"), fileConfigHandler.GetFileConfigPage)
				fileConfigGroup.GET("/get", casbinMiddleware.RequirePermission("infra:file-config:query"), fileConfigHandler.GetFileConfig)
				fileConfigGroup.GET("/test", casbinMiddleware.RequirePermission("infra:file-config:query"), fileConfigHandler.TestFileConfig)
			}

			// File
			fileGroup := infraGroup.Group("/file")
			{
				fileGroup.POST("/upload", fileHandler.UploadFile)
				fileGroup.DELETE("/delete", casbinMiddleware.RequirePermission("infra:file:delete"), fileHandler.DeleteFile)
				fileGroup.GET("/page", casbinMiddleware.RequirePermission("infra:file:query"), fileHandler.GetFilePage)
				fileGroup.GET("/presigned-url", casbinMiddleware.RequirePermission("infra:file:query"), fileHandler.GetFilePresignedUrl)
				fileGroup.POST("/create", casbinMiddleware.RequirePermission("infra:file:create"), fileHandler.CreateFile)
			}

			// Job
			jobGroup := infraGroup.Group("/job")
			{
				jobGroup.POST("/create", casbinMiddleware.RequirePermission("infra:job:create"), jobHandler.CreateJob)
				jobGroup.PUT("/update", casbinMiddleware.RequirePermission("infra:job:update"), jobHandler.UpdateJob)
				jobGroup.PUT("/update-status", casbinMiddleware.RequirePermission("infra:job:update"), jobHandler.UpdateJobStatus)
				jobGroup.DELETE("/delete", casbinMiddleware.RequirePermission("infra:job:delete"), jobHandler.DeleteJob)
				jobGroup.GET("/get", casbinMiddleware.RequirePermission("infra:job:query"), jobHandler.GetJob)
				jobGroup.GET("/page", casbinMiddleware.RequirePermission("infra:job:query"), jobHandler.GetJobPage)
				jobGroup.PUT("/trigger", casbinMiddleware.RequirePermission("infra:job:trigger"), jobHandler.TriggerJob)
				jobGroup.POST("/sync", casbinMiddleware.RequirePermission("infra:job:create"), jobHandler.SyncJob)
				jobGroup.GET("/export-excel", casbinMiddleware.RequirePermission("infra:job:export"), jobHandler.ExportJobExcel)
				jobGroup.GET("/get_next_times", casbinMiddleware.RequirePermission("infra:job:query"), jobHandler.GetJobNextTimes)
			}

			// Job Log
			jobLogGroup := infraGroup.Group("/job-log")
			{
				jobLogGroup.GET("/get", jobLogHandler.GetJobLog)
				jobLogGroup.GET("/page", jobLogHandler.GetJobLogPage)
				jobLogGroup.GET("/export-excel", casbinMiddleware.RequirePermission("infra:job:export"), jobLogHandler.ExportJobLogExcel)
			}

			// API Access Log
			apiAccessLogGroup := infraGroup.Group("/api-access-log")
			{
				apiAccessLogGroup.GET("/page", casbinMiddleware.RequirePermission("infra:api-access-log:query"), apiAccessLogHandler.GetApiAccessLogPage)
			}

			// API Error Log
			apiErrorLogGroup := infraGroup.Group("/api-error-log")
			{
				apiErrorLogGroup.GET("/page", casbinMiddleware.RequirePermission("infra:api-error-log:query"), apiErrorLogHandler.GetApiErrorLogPage)
				apiErrorLogGroup.PUT("/update-status", casbinMiddleware.RequirePermission("infra:api-error-log:update-status"), apiErrorLogHandler.UpdateApiErrorLogProcess)
			}
		}
	}
}

// RegisterAreaRoutes 注册地区路由 (Public - 不需要认证)
func RegisterAreaRoutes(engine *gin.Engine, areaHandler *handler.AreaHandler) {
	api := engine.Group("/admin-api")
	{
		// Area 地区 (Public Routes)
		areaGroup := api.Group("/system/area")
		{
			areaGroup.GET("/tree", areaHandler.GetAreaTree)
			areaGroup.GET("/get-by-ip", areaHandler.GetAreaByIP)
		}
	}
}
