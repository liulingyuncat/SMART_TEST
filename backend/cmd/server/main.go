package main

import (
	"fmt"
	"log"
	"time"
	"webtest/config"
	"webtest/internal/constants"
	"webtest/internal/handlers"
	"webtest/internal/middleware"
	"webtest/internal/models"
	"webtest/internal/repositories"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 数据库配置
	dbConfig := &config.DatabaseConfig{
		Type:            "sqlite",        // 开发环境使用 SQLite
		DBName:          "webtest.db",    // 数据库文件路径
		MaxOpenConns:    25,              // 最大打开连接数
		MaxIdleConns:    10,              // 最大空闲连接数
		ConnMaxLifetime: 5 * time.Minute, // 连接最大生命周期
		ConnMaxIdleTime: 1 * time.Minute, // 连接最大空闲时间
	}

	// 初始化数据库
	db, err := config.InitDatabase(dbConfig)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// 自动迁移数据库模型
	if err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Requirement{},
		&models.ManualTestCase{},
		&models.AutoTestCase{},
		&models.CaseReview{},
		&models.CaseVersion{},
		&models.AutoTestCaseVersion{},
		&models.ExecutionTask{},
		&models.Defect{},
		&models.DefectAttachment{},
		&models.DefectSubject{},
		&models.DefectPhase{},
		&models.DefectComment{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("database migration completed")

	// 初始化依赖
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	memberRepo := repositories.NewProjectMemberRepository(db)
	requirementRepo := repositories.NewRequirementRepository(db)
	manualCaseRepo := repositories.NewManualTestCaseRepository(db)
	autoCaseRepo := repositories.NewAutoTestCaseRepository(db)
	apiCaseRepo := repositories.NewApiTestCaseRepository(db)
	caseReviewRepo := repositories.NewCaseReviewRepository(db)
	caseVersionRepo := repositories.NewCaseVersionRepository(db)

	// 缺陷管理相关Repository
	defectRepo := repositories.NewDefectRepository(db)
	defectAttachmentRepo := repositories.NewDefectAttachmentRepository(db)
	defectSubjectRepo := repositories.NewDefectSubjectRepository(db)
	defectPhaseRepo := repositories.NewDefectPhaseRepository(db)
	defectCommentRepo := repositories.NewDefectCommentRepository(db)

	authService := services.NewAuthService(userRepo)
	projectService := services.NewProjectService(projectRepo, memberRepo, db)
	requirementService := services.NewRequirementService(requirementRepo)
	manualCaseService := services.NewManualTestCaseService(manualCaseRepo, projectService)
	autoCaseService := services.NewAutoTestCaseService(autoCaseRepo, projectService, db)
	apiCaseService := services.NewApiTestCaseService(apiCaseRepo, projectService)
	executionTaskRepo := repositories.NewExecutionTaskRepository(db)
	executionCaseResultRepo := repositories.NewExecutionCaseResultRepository(db)
	excelService := services.NewExcelService(manualCaseRepo, projectRepo, executionCaseResultRepo)
	versionService := services.NewVersionService(db, caseVersionRepo, excelService)
	reviewService := services.NewReviewService(caseReviewRepo)
	executionTaskService := services.NewExecutionTaskService(executionTaskRepo, projectRepo)
	executionCaseResultService := services.NewExecutionCaseResultService(
		executionCaseResultRepo,
		executionTaskRepo,
		manualCaseRepo,
		autoCaseRepo,
		apiCaseRepo,
	)

	// 缺陷管理相关Service
	defectService := services.NewDefectService(defectRepo)
	defectAttachmentService := services.NewDefectAttachmentService(defectAttachmentRepo, "storage")
	defectConfigService := services.NewDefectConfigService(defectSubjectRepo, defectPhaseRepo)
	defectCommentService := services.NewDefectCommentService(defectCommentRepo, defectRepo)

	authHandler := handlers.NewAuthHandler(authService)
	projectHandler := handlers.NewProjectHandler(projectService)
	requirementHandler := handlers.NewRequirementHandler(requirementService, projectService)
	manualCasesHandler := handlers.NewManualCasesHandler(manualCaseService)
	autoCasesHandler := handlers.NewAutoCasesHandler(autoCaseService)
	apiCasesHandler := handlers.NewApiTestCaseHandler(apiCaseService)
	exportHandler := handlers.NewExportHandler(excelService)
	importHandler := handlers.NewImportHandler(excelService)
	versionHandler := handlers.NewVersionHandler(versionService, requirementService, projectService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	executionTaskHandler := handlers.NewExecutionTaskHandler(executionTaskService)
	executionCaseResultHandler := handlers.NewExecutionCaseResultHandler(executionCaseResultService, executionTaskService)

	// 缺陷管理相关Handler
	defectHandler := handlers.NewDefectHandler(defectService)
	defectAttachmentHandler := handlers.NewDefectAttachmentHandler(defectAttachmentService)
	defectConfigHandler := handlers.NewDefectConfigHandler(defectConfigService)
	defectCommentHandler := handlers.NewDefectCommentHandler(defectCommentService)

	// 初始化管理员账号
	if err := authService.InitAdminUsers(); err != nil {
		log.Printf("warning: failed to init admin users: %v", err)
	} else {
		log.Println("admin users initialized successfully")
	}

	// 创建 Gin 路由引擎
	r := gin.Default()

	// 配置 CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 为所有路由添加OPTIONS处理
	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.AbortWithStatus(204)
	})

	// 注册路由
	api := r.Group("/api/v1")
	{
		// 1. 公开路由(无需认证)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 2. 认证路由(需要登录,但不限角色)
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware(authService))
		{
			authenticated.POST("/auth/logout", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "logout successful"})
			})
			authenticated.GET("/profile", func(c *gin.Context) {
				username, _ := c.Get("username")
				userID, _ := c.Get("userID")
				role, _ := c.Get("role")
				c.JSON(200, gin.H{
					"user_id":  userID,
					"username": username,
					"role":     role,
					"message":  "authenticated user profile",
				})
			})
			// TODO: 添加 PUT /profile/nickname 和 PUT /profile/password
		}

		// 3. 项目管理路由(PM + PMemb 可查看,PM 可操作)
		projects := api.Group("/projects")
		projects.Use(middleware.AuthMiddleware(authService))
		{
			// GET - 允许PM和PM Member查看
			projects.GET("",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				projectHandler.GetProjects)

			// POST - 仅PM可创建
			projects.POST("",
				middleware.RequireRole(constants.RoleProjectManager),
				projectHandler.CreateProject)

			// PUT - 仅PM可更新
			projects.PUT("/:id",
				middleware.RequireRole(constants.RoleProjectManager),
				projectHandler.UpdateProject)

			// DELETE - 仅PM可删除
			projects.DELETE("/:id",
				middleware.RequireRole(constants.RoleProjectManager),
				projectHandler.DeleteProject)

			// GET BY ID - 允许PM和PM Member查看项目详情
			projects.GET("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				projectHandler.GetProjectByID)

			// 需求文档路由 - 允许PM和PM Member访问
			projects.GET("/:id/requirements/:type",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementHandler.GetRequirement)
			projects.PUT("/:id/requirements/:type",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementHandler.UpdateRequirement)

			// 手工测试用例路由 - 允许PM和PM Member访问
			projects.GET("/:id/manual-cases/metadata",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.GetMetadata)
			projects.PUT("/:id/manual-cases/metadata",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.UpdateMetadata)
			projects.GET("/:id/manual-cases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.GetCases)
			projects.POST("/:id/manual-cases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.CreateCase)
			projects.PATCH("/:id/manual-cases/:caseId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.UpdateCase)
			projects.DELETE("/:id/manual-cases/:caseId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.DeleteCase)
			projects.POST("/:id/manual-cases/reorder",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.ReorderCases)
			projects.POST("/:id/manual-cases/reorder-drag",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.ReorderCasesByDrag)
			projects.POST("/:id/manual-cases/reorder-all",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.ReorderAllCasesByID)
			projects.DELETE("/:id/manual-cases/clear-ai",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.ClearAICases)
			projects.POST("/:id/manual-cases/insert",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.InsertCase)
			projects.POST("/:id/manual-cases/batch-delete",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.BatchDeleteCases)
			projects.POST("/:id/manual-cases/reassign-ids",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.ReassignIDs)

			// 自动化测试用例路由 - 允许PM和PM Member访问
			projects.GET("/:id/auto-cases/metadata",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.GetMetadata)
			projects.PUT("/:id/auto-cases/metadata",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.UpdateMetadata)
			projects.GET("/:id/auto-cases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.GetCases)
			projects.POST("/:id/auto-cases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.CreateCase)
			projects.PATCH("/:id/auto-cases/:caseId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.UpdateCase)
			projects.DELETE("/:id/auto-cases/:caseId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.DeleteCase)
			projects.POST("/:id/auto-cases/reorder",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.ReorderAllCases)
			projects.POST("/:id/auto-cases/insert",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.InsertCase)
			projects.POST("/:id/auto-cases/batch-delete",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.BatchDeleteCases)
			projects.POST("/:id/auto-cases/reassign-ids",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.ReassignIDs)

			// 自动化测试用例版本管理路由
			projects.POST("/:id/auto-cases/versions",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.BatchSaveVersion)
			projects.GET("/:id/auto-cases/versions",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.GetAutoVersions)
			projects.GET("/:id/auto-cases/versions/:versionId/export",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.DownloadAutoVersion)
			projects.DELETE("/:id/auto-cases/versions/:versionId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.DeleteAutoVersion)
			projects.PUT("/:id/auto-cases/versions/:versionId/remark",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.UpdateAutoVersionRemark)

			// 接口测试用例路由 - 允许PM和PM Member访问
			projects.GET("/:id/api-cases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.GetCases)
			projects.POST("/:id/api-cases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.CreateCase)
			projects.POST("/:id/api-cases/insert",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.InsertCase)
			projects.DELETE("/:id/api-cases/:caseId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.DeleteCase)
			projects.POST("/:id/api-cases/batch-delete",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.BatchDeleteCases)
			projects.PATCH("/:id/api-cases/:caseId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.UpdateCase)

			// 接口测试用例版本管理路由
			projects.POST("/:id/api-cases/versions",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.SaveVersion)
			projects.GET("/:id/api-cases/versions",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.GetVersions)
			projects.GET("/:id/api-cases/versions/:versionId/export",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.DownloadVersion)
			projects.DELETE("/:id/api-cases/versions/:versionId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.DeleteVersion)
			projects.PUT("/:id/api-cases/versions/:versionId/remark",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.UpdateVersionRemark)

			// 导出导入路由
			projects.GET("/:id/manual-cases/export/ai",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				exportHandler.ExportAICases)
			projects.GET("/:id/manual-cases/export/template",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				exportHandler.ExportTemplate)
			projects.GET("/:id/manual-cases/export/cases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				exportHandler.ExportCases)
			projects.POST("/:id/manual-cases/import",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				func(c *gin.Context) {
					fmt.Println("[ROUTE] Import route matched!")
					importHandler.ImportCases(c)
				})

			// 版本管理路由
			projects.POST("/:id/versions/save",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				versionHandler.SaveVersion)
			projects.GET("/:id/versions",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				versionHandler.GetVersionList)
			projects.GET("/:id/versions/:versionID/download",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				versionHandler.DownloadVersion)
			projects.DELETE("/:id/versions/:versionID",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				versionHandler.DeleteVersion)

			// 评审内容路由
			projects.GET("/:id/review",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				reviewHandler.GetCaseReview)
			projects.POST("/:id/review",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				reviewHandler.SaveCaseReview)

			// 测试执行任务路由
			projects.GET("/:id/execution-tasks",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				executionTaskHandler.GetTasks)
			projects.POST("/:id/execution-tasks",
				middleware.RequireRole(constants.RoleProjectManager),
				executionTaskHandler.CreateTask)
			projects.PUT("/:id/execution-tasks/:task_uuid",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				executionTaskHandler.UpdateTask)
			projects.DELETE("/:id/execution-tasks/:task_uuid",
				middleware.RequireRole(constants.RoleProjectManager),
				executionTaskHandler.DeleteTask)

			// 缺陷管理路由 - 允许PM和PM Member访问
			projects.GET("/:id/defects",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectHandler.GetDefects)
			projects.POST("/:id/defects",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectHandler.CreateDefect)
			projects.GET("/:id/defects/template",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectHandler.ExportTemplate)
			projects.POST("/:id/defects/import",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectHandler.ImportDefects)
			projects.GET("/:id/defects/export",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectHandler.ExportDefects)
			projects.GET("/:id/defects/:defectId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectHandler.GetDefect)
			projects.PUT("/:id/defects/:defectId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectHandler.UpdateDefect)
			projects.DELETE("/:id/defects/:defectId",
				middleware.RequireRole(constants.RoleProjectManager),
				defectHandler.DeleteDefect)

			// 缺陷说明管理路由
			projects.GET("/:id/defects/:defectId/comments",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectCommentHandler.GetComments)
			projects.POST("/:id/defects/:defectId/comments",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectCommentHandler.CreateComment)
			projects.PUT("/:id/defects/:defectId/comments/:commentId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectCommentHandler.UpdateComment)
			projects.DELETE("/:id/defects/:defectId/comments/:commentId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectCommentHandler.DeleteComment)

			// 缺陷配置管理路由 - Subject
			projects.GET("/:id/defect-subjects",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectConfigHandler.GetSubjects)
			projects.POST("/:id/defect-subjects",
				middleware.RequireRole(constants.RoleProjectManager),
				defectConfigHandler.CreateSubject)
			projects.PUT("/:id/defect-subjects/:subjectId",
				middleware.RequireRole(constants.RoleProjectManager),
				defectConfigHandler.UpdateSubject)
			projects.DELETE("/:id/defect-subjects/:subjectId",
				middleware.RequireRole(constants.RoleProjectManager),
				defectConfigHandler.DeleteSubject)

			// 缺陷配置管理路由 - Phase
			projects.GET("/:id/defect-phases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				defectConfigHandler.GetPhases)
			projects.POST("/:id/defect-phases",
				middleware.RequireRole(constants.RoleProjectManager),
				defectConfigHandler.CreatePhase)
			projects.PUT("/:id/defect-phases/:phaseId",
				middleware.RequireRole(constants.RoleProjectManager),
				defectConfigHandler.UpdatePhase)
			projects.DELETE("/:id/defect-phases/:phaseId",
				middleware.RequireRole(constants.RoleProjectManager),
				defectConfigHandler.DeletePhase)
		}

		// 测试执行用例结果路由
		executionGroup := api.Group("/execution-tasks")
		executionGroup.Use(middleware.AuthMiddleware(authService))
		{
			executionGroup.GET("/:taskUuid/case-results",
				executionCaseResultHandler.GetExecutionCaseResults)
			executionGroup.PATCH("/:taskUuid/case-results",
				executionCaseResultHandler.SaveExecutionCaseResults)
			executionGroup.GET("/:taskUuid/statistics",
				executionCaseResultHandler.GetExecutionStatistics)
			executionGroup.POST("/:taskUuid/case-results/init",
				executionCaseResultHandler.InitExecutionResults)
		}

		// 缺陷附件管理路由(在projects路由组下)
		projectDefects := projects.Group("/:id/defects")
		{
			projectDefects.GET("/:defectId/attachments",
				defectAttachmentHandler.List)
			projectDefects.POST("/:defectId/attachments",
				defectAttachmentHandler.Upload)
			projectDefects.GET("/:defectId/attachments/:attId",
				defectAttachmentHandler.Download)
			projectDefects.DELETE("/:defectId/attachments/:attId",
				defectAttachmentHandler.Delete)
		}

		// 缺陷附件管理路由(兼容旧路由)
		defectAttachments := api.Group("/defects")
		defectAttachments.Use(middleware.AuthMiddleware(authService))
		{
			defectAttachments.POST("/:defectId/attachments",
				defectAttachmentHandler.Upload)
			defectAttachments.GET("/:defectId/attachments/:attId",
				defectAttachmentHandler.Download)
			defectAttachments.DELETE("/:defectId/attachments/:attId",
				defectAttachmentHandler.Delete)
		}

		// 通用版本管理路由(支持需求管理)
		apiVersions := api.Group("/versions")
		apiVersions.Use(middleware.AuthMiddleware(authService))
		{
			apiVersions.POST("", versionHandler.SaveVersionGeneric)
			apiVersions.GET("", versionHandler.GetVersionListGeneric)
			apiVersions.GET("/:id/download", versionHandler.DownloadVersionGeneric)
			apiVersions.DELETE("/:id", versionHandler.DeleteVersionGeneric)
			apiVersions.PUT("/:id/remark", versionHandler.UpdateVersionRemarkGeneric)
		}

		// 4. 用户管理路由（系统管理员和项目管理员共用，内部分流）
		userService := services.NewUserService(userRepo)
		userHandler := handlers.NewUserHandler(userService)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(authService))
		users.GET("", userHandler.GetUsersAuto)
		users.GET("/check", userHandler.CheckUnique)
		users.POST("", userHandler.CreateUserAuto)
		users.PUT("/:id", userHandler.UpdateNickname)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.POST("/:id/reset-password", userHandler.ResetPassword)

		// 5. 人员分配路由(仅 PM)
		api.GET("/projects/:id/members",
			middleware.AuthMiddleware(authService),
			middleware.RequireRole(constants.RoleProjectManager),
			func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "project members"})
			})
		api.PUT("/projects/:id/members",
			middleware.AuthMiddleware(authService),
			middleware.RequireRole(constants.RoleProjectManager),
			func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "update project members"})
			})
	}

	// 服务前端静态文件
	r.Static("/static", "../frontend/build/static")

	// 对于所有非 API 路由，返回 index.html (支持前端路由)
	r.NoRoute(func(c *gin.Context) {
		c.File("../frontend/build/index.html")
	})
	log.Println("server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
