package router

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/xxzhwl/gaia/framework/server"

	"gaia-server/app/api"
	"gaia-server/app/middleware"
)

func RegisterRouters(s *server.Server) {
	group := s.Group("/api")
	s.RegisterCommonHandler(group, middleware.AuthMiddleware())

	registerHealth(group)
	registerAuth(group)
	registerMFA(group)
	registerSession(group)
	registerUser(group)
	registerRole(group)
	registerPermission(group)
	registerLog(group)
	registerJobs(group)
	registerAsyncTask(group)
	registerMetrics(group)
	registerConfig(group)
	registerAdmin(group)
}

func registerAuth(s *route.RouterGroup) {
	ctrl := api.NewAuthCtrl()
	g := s.Group("auth")
	g.POST("/login", ctrl.Login())
	g.POST("/register", ctrl.Register())
	g.POST("/mfa/complete", ctrl.CompleteMFA())
	g.POST("/forgot-password/start", ctrl.ForgotPasswordStart())
	g.POST("/forgot-password/complete", ctrl.ForgotPasswordComplete())
	g.POST("/refresh", ctrl.RefreshToken())
	g.POST("/logout", ctrl.Logout())
}

func registerMFA(s *route.RouterGroup) {
	ctrl := api.NewMFACtrl()
	g := s.Group("mfa")
	g.Use(middleware.AuthMiddleware())
	g.POST("/totp/setup", ctrl.SetupTOTP())
	g.POST("/totp/verify", ctrl.VerifyTOTP())
	g.DELETE("/totp", ctrl.DisableTOTP())
	g.POST("/step-up/start", ctrl.StartStepUp())
	g.POST("/step-up/complete", ctrl.CompleteStepUp())
}

func registerSession(s *route.RouterGroup) {
	ctrl := api.NewSessionCtrl()
	g := s.Group("sessions")
	g.Use(middleware.AuthMiddleware())
	g.GET("", ctrl.List())
	g.DELETE("/:id", ctrl.Revoke())
	g.POST("/revoke-others", ctrl.RevokeOthers())
}

func registerUser(s *route.RouterGroup) {
	ctrl := api.NewUserCtrl()
	g := s.Group("user")
	g.Use(middleware.AuthMiddleware())
	g.POST("/info", ctrl.GetUserInfo())
	g.POST("/update", ctrl.UpdateUserInfo())
	g.POST("/change-password", ctrl.ChangePassword())
	g.POST("/list", ctrl.GetUserList())
	g.POST("/create", ctrl.CreateUser())
	g.POST("/update-status", ctrl.UpdateUserStatus())
	g.POST("/reset-mfa", ctrl.ResetUserMFA())
}

func registerRole(s *route.RouterGroup) {
	ctrl := api.NewRoleCtrl()
	g := s.Group("role")
	g.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	g.POST("/create", ctrl.CreateRole())
	g.POST("/update", ctrl.UpdateRole())
	g.POST("/delete", ctrl.DeleteRole())
	g.POST("/list", ctrl.ListRoles())
	g.POST("/permissions", ctrl.GetRolePermissions())
	g.POST("/assign-user", ctrl.AssignUserRole())
	g.POST("/remove-user", ctrl.RemoveUserRole())
}

func registerPermission(s *route.RouterGroup) {
	ctrl := api.NewPermissionCtrl()
	g := s.Group("permission")
	g.Use(middleware.AuthMiddleware())
	g.POST("/create", ctrl.Create())
	g.POST("/update", ctrl.Update())
	g.POST("/delete", ctrl.Delete())
	g.GET("/list", ctrl.List())
	g.POST("/assign-role", ctrl.AssignToRole())
	g.POST("/remove-role", ctrl.RemoveFromRole())
}

func registerLog(s *route.RouterGroup) {
	ctrl := api.NewLogCtrl()
	g := s.Group("log")
	g.Use(middleware.AuthMiddleware())
	g.POST("/query", ctrl.QueryES())
	g.GET("/indices", ctrl.GetIndices())
	g.GET("/mapping", ctrl.GetMapping())
	g.POST("/export", ctrl.ExportLogs())
}

func registerHealth(s *route.RouterGroup) {
	ctrl := api.NewHealthCtrl()
	g := s.Group("health")
	g.Any("", ctrl.CheckHealth())
	g.Any("/detailed", ctrl.CheckHealthDetailed())
	g.Any("/services", ctrl.GetAvailableServices())
}

func registerJobs(s *route.RouterGroup) {
	ctrl := api.NewJobsCtrl()
	g := s.Group("jobs")
	g.Use(middleware.AuthMiddleware())
	g.GET("/list", ctrl.ListJobs())
	g.GET("/detail/:id", ctrl.GetJobDetail())
	g.POST("/create", ctrl.CreateJob())
	g.POST("/update", ctrl.UpdateJob())
	g.DELETE("/delete/:id", ctrl.DeleteJob())
	g.POST("/toggle", ctrl.ToggleJob())
	g.POST("/execute", ctrl.ExecuteJobImmediately())
	g.GET("/records", ctrl.GetJobRecords())
	g.GET("/running", ctrl.GetRunningJobs())
	g.GET("/types", ctrl.GetJobTypes())
	g.GET("/stats", ctrl.GetJobStats())
	g.POST("/stop-executor", ctrl.StopExecutor())
	g.POST("/start-executor", ctrl.StartExecutor())
}

func registerAsyncTask(s *route.RouterGroup) {
	ctrl := api.NewAsyncTaskCtrl()
	g := s.Group("asynctask")
	g.Use(middleware.AuthMiddleware())
	g.GET("/", ctrl.GetSchedulerInfoPoll())
	g.GET("/list", ctrl.ListTasks())
	g.GET("/detail/:id", ctrl.GetTaskDetail())
	g.POST("/retry", ctrl.RetryTask())
	g.POST("/cancel", ctrl.CancelTask())
	g.GET("/records", ctrl.GetTaskExecRecords())
	g.GET("/schedulers", ctrl.GetAllSchedulerStatus())
	g.POST("/stop", ctrl.StopScheduler())
	g.POST("/start", ctrl.StartScheduler())
}

func registerConfig(s *route.RouterGroup) {
	ctrl := api.NewConfigCtrl()
	g := s.Group("config")
	g.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	g.GET("/local", ctrl.GetLocalConfig())
	g.GET("/remote", ctrl.GetRemoteConfig())
	g.GET("/env", ctrl.GetEnvironmentConfig())
	g.GET("/query", ctrl.GetConfigByKey())
}

func registerMetrics(s *route.RouterGroup) {
	ctrl := api.NewMetricsCtrl()
	g := s.Group("metrics")
	g.Use(middleware.AuthMiddleware())
	g.GET("/overview", ctrl.GetMetricsOverview())
	g.GET("/jobs", ctrl.GetJobsMetrics())
	g.GET("/asynctask", ctrl.GetAsyncTaskMetrics())
	g.GET("/jobs-global", ctrl.GetJobsGlobalMetrics())
}

func registerAdmin(s *route.RouterGroup) {
	ctrl := api.NewAdminCtrl()
	g := s.Group("admin")
	g.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	g.GET("/users/:id/sessions", ctrl.ListUserSessions())
	g.GET("/users/:id/permissions", ctrl.GetUserPermissions())
	g.DELETE("/sessions/:id", ctrl.RevokeSession())
	g.GET("/audit", ctrl.QueryAudit())
	g.GET("/audit/archived", ctrl.QueryArchivedAudit())
	g.GET("/audit/events", ctrl.ListAuditEvents())
	g.GET("/audit/:id", ctrl.GetAuditLog())
	g.POST("/audit/restore/:id", ctrl.RestoreAuditLog())
}
