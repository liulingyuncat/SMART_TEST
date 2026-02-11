package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"webtest/config"
	"webtest/internal/constants"
	"webtest/internal/handlers"
	"webtest/internal/middleware"
	"webtest/internal/models"
	"webtest/internal/repositories"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// getStorageBasePath 返回存储目录的绝对路径
func getStorageBasePath() string {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: could not get executable path: %v, using relative path", err)
		return "storage"
	}
	exeDir := filepath.Dir(exePath)
	log.Printf("Executable directory: %s", exeDir)

	// 尝试多种路径：
	// 1. Docker 容器路径: /app/server -> /app/storage
	// 2. 本地开发路径: backend/ -> backend/storage
	candidates := []string{
		filepath.Join(exeDir, "storage"),       // Docker/部署: exeDir/storage
		filepath.Join(exeDir, "..", "storage"), // Local: exeDir/../storage
	}

	for _, candidate := range candidates {
		absPath, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		// 检查目录是否存在
		if _, err := os.Stat(absPath); err == nil {
			log.Printf("Found storage at: %s", absPath)
			return absPath
		}
	}

	// 默认返回相对路径
	log.Printf("Warning: could not find storage directory, using relative path")
	return "storage"
}

// getFrontendBuildPath 返回前端 build 目录的绝对路径
func getFrontendBuildPath() string {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: could not get executable path: %v, using relative path", err)
		return "../frontend/build"
	}
	exeDir := filepath.Dir(exePath)
	log.Printf("Executable directory: %s", exeDir)

	// 尝试多种路径：
	// 1. Docker 容器路径: /app/server -> /app/frontend/build
	// 2. 本地开发路径: backend/ -> webtest/frontend/build
	candidates := []string{
		filepath.Join(exeDir, "frontend", "build"),       // Docker: /app/frontend/build
		filepath.Join(exeDir, "..", "frontend", "build"), // Local: backend/ -> webtest/frontend/build
	}

	for _, candidate := range candidates {
		absPath, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		// 检查 index.html 是否存在
		indexPath := filepath.Join(absPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			log.Printf("Found frontend build at: %s", absPath)
			return absPath
		}
	}

	// 默认返回相对路径
	log.Printf("Warning: could not find frontend build, using relative path")
	return "../frontend/build"
}

func main() {
	// 从环境变量读取数据库配置（支持 SQLite 和 PostgreSQL）
	dbConfig := config.GetDatabaseConfigFromEnv()
	log.Printf("Database configuration: type=%s, host=%s, port=%d, dbname=%s",
		dbConfig.Type, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	// 初始化数据库
	db, err := config.InitDatabase(dbConfig)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	log.Printf("Database connected successfully (type: %s)", dbConfig.Type)

	// 自动迁移数据库模型
	if err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.RequirementItem{},
		&models.ViewpointItem{},
		&models.RequirementChunk{}, // T54: 需求Chunk表
		&models.ViewpointChunk{},   // T54: 观点Chunk表
		&models.Version{},
		&models.ManualTestCase{},
		&models.AutoTestCase{},
		&models.ApiTestCase{},        // T46: API测试用例表
		&models.ApiTestCaseVersion{}, // T46: API测试用例版本表
		&models.CaseReview{},
		&models.CaseVersion{},
		&models.AutoTestCaseVersion{},
		&models.ExecutionTask{},
		&models.ExecutionCaseResult{}, // 测试执行结果表
		&models.Defect{},
		&models.DefectAttachment{},
		&models.DefectSubject{},
		&models.DefectPhase{},
		&models.DefectComment{},
		&models.CaseReviewItem{},      // T44: 审阅条目表
		&models.CaseGroup{},           // 用例集表
		&models.WebCaseVersion{},      // T45: Web用例版本表
		&models.AIReport{},            // T47: AI质量报告表
		&models.RawDocument{},         // T48: 原始需求文档表
		&models.Prompt{},              // 提示词表
		&models.UserDefinedVariable{}, // 用户自定义变量表
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("database migration completed")

	// 初始化依赖
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	memberRepo := repositories.NewProjectMemberRepository(db)
	manualCaseRepo := repositories.NewManualTestCaseRepository(db)
	autoCaseRepo := repositories.NewAutoTestCaseRepository(db)
	apiCaseRepo := repositories.NewApiTestCaseRepository(db)
	caseReviewRepo := repositories.NewCaseReviewRepository(db)
	caseVersionRepo := repositories.NewCaseVersionRepository(db)
	caseGroupRepo := repositories.NewCaseGroupRepository(db)

	// 需求管理相关Repository (T42)
	requirementItemRepo := repositories.NewRequirementItemRepository(db)
	viewpointItemRepo := repositories.NewViewpointItemRepository(db)
	versionRepo := repositories.NewVersionRepository(db)

	// 需求Chunk相关Repository (T54)
	requirementChunkRepo := repositories.NewRequirementChunkRepository(db)
	viewpointChunkRepo := repositories.NewViewpointChunkRepository(db)

	// 缺陷管理相关Repository
	defectRepo := repositories.NewDefectRepository(db)
	defectAttachmentRepo := repositories.NewDefectAttachmentRepository(db)
	defectSubjectRepo := repositories.NewDefectSubjectRepository(db)
	defectPhaseRepo := repositories.NewDefectPhaseRepository(db)
	defectCommentRepo := repositories.NewDefectCommentRepository(db)

	// 审阅条目相关Repository (T44)
	reviewItemRepo := repositories.NewReviewItemRepository(db)

	// AI质量报告相关Repository (T47)
	aiReportRepo := repositories.NewAIReportRepository(db)

	// 原始需求文档相关Repository (T48)
	rawDocumentRepo := repositories.NewRawDocumentRepository(db)

	// 用户自定义变量相关Repository
	userDefinedVarRepo := repositories.NewUserDefinedVariableRepository(db)

	authService := services.NewAuthService(userRepo)
	projectService := services.NewProjectService(projectRepo, memberRepo, userRepo, db)
	manualCaseService := services.NewManualTestCaseService(manualCaseRepo, projectService)
	autoCaseService := services.NewAutoTestCaseService(autoCaseRepo, projectService, db, caseGroupRepo)
	apiCaseService := services.NewApiTestCaseService(apiCaseRepo, projectService)
	executionTaskRepo := repositories.NewExecutionTaskRepository(db)
	executionCaseResultRepo := repositories.NewExecutionCaseResultRepository(db)
	excelService := services.NewExcelService(manualCaseRepo, projectRepo, executionCaseResultRepo)
	versionService := services.NewVersionService(db, caseVersionRepo, excelService)
	reviewService := services.NewReviewService(caseReviewRepo)

	// T45: Web用例版本管理Service
	webVersionService := services.NewWebVersionService(db, projectRepo, caseGroupRepo, autoCaseRepo, excelService)

	// 审阅条目相关Service (T44)
	reviewItemService := services.NewReviewItemService(reviewItemRepo)

	// AI质量报告相关Service (T47)
	aiReportService := services.NewAIReportService(aiReportRepo)

	// 需求管理相关Service (T42)
	storageDir := getStorageBasePath()
	requirementItemService := services.NewRequirementItemService(requirementItemRepo, requirementChunkRepo, versionRepo, storageDir)
	viewpointItemService := services.NewViewpointItemService(viewpointItemRepo, viewpointChunkRepo, versionRepo, storageDir)

	// 需求Chunk相关Service (T54)
	requirementChunkService := services.NewRequirementChunkService(requirementChunkRepo)
	viewpointChunkService := services.NewViewpointChunkService(viewpointChunkRepo)

	// 用户自定义变量相关Service (需要在executionTaskService之前初始化)
	userDefinedVarService := services.NewUserDefinedVariableService(userDefinedVarRepo)

	executionTaskService := services.NewExecutionTaskService(executionTaskRepo, projectRepo, executionCaseResultRepo, userRepo, userDefinedVarService)
	executionCaseResultService := services.NewExecutionCaseResultService(
		executionCaseResultRepo,
		executionTaskRepo,
		manualCaseRepo,
		autoCaseRepo,
		apiCaseRepo,
	)

	// 缺陷管理相关Service
	defectService := services.NewDefectService(defectRepo, userRepo)
	defectAttachmentService := services.NewDefectAttachmentService(defectAttachmentRepo, getStorageBasePath())
	defectConfigService := services.NewDefectConfigService(defectSubjectRepo, defectPhaseRepo)
	defectCommentService := services.NewDefectCommentService(defectCommentRepo, defectRepo)

	// 原始需求文档相关Service (T48)
	rawDocumentService := services.NewRawDocumentService(rawDocumentRepo, storageDir)

	// T51: 提示词管理相关Service
	promptService := services.NewPromptService(db)

	// 加载系统提示词到数据库（首次初始化或更新）
	if err := initSystemPrompts(db); err != nil {
		log.Printf("warning: failed to init system prompts: %v", err)
	} else {
		log.Println("system prompts initialized successfully")
	}

	authHandler := handlers.NewAuthHandler(authService)
	projectHandler := handlers.NewProjectHandler(projectService)
	manualCasesHandler := handlers.NewManualCasesHandler(manualCaseService, versionService)
	autoCasesHandler := handlers.NewAutoCasesHandler(autoCaseService)
	apiCasesHandler := handlers.NewApiTestCaseHandler(apiCaseService)
	apiCaseGroupHandler := handlers.NewApiCaseGroupHandler(db)
	exportHandler := handlers.NewExportHandler(excelService)
	importHandler := handlers.NewImportHandler(excelService, caseGroupRepo)
	versionHandler := handlers.NewVersionHandler(versionService, projectService, versionRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	executionTaskHandler := handlers.NewExecutionTaskHandler(executionTaskService)
	executionCaseResultHandler := handlers.NewExecutionCaseResultHandler(executionCaseResultService, executionTaskService)
	caseGroupHandler := handlers.NewCaseGroupHandler(caseGroupRepo, autoCaseRepo, manualCaseRepo, apiCaseRepo)

	// T18: 用户管理 & T22: 个人信息
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	profileHandler := handlers.NewProfileHandlerWithProject(userService, projectService)

	// T45: Web用例版本管理Handler
	webVersionHandler := handlers.NewWebVersionHandler(webVersionService)

	// 需求管理相关Handler (T42)
	requirementItemHandler := handlers.NewRequirementItemHandler(requirementItemService, projectService, storageDir)
	viewpointItemHandler := handlers.NewViewpointItemHandler(viewpointItemService, projectService, storageDir)

	// 需求Chunk相关Handler (T54)
	requirementChunkHandler := handlers.NewRequirementChunkHandler(requirementChunkService)
	viewpointChunkHandler := handlers.NewViewpointChunkHandler(viewpointChunkService)

	// 审阅条目相关Handler (T44)
	reviewItemHandler := handlers.NewReviewItemHandler(reviewItemService, projectService)

	// AI质量报告相关Handler (T47)
	aiReportHandler := handlers.NewAIReportHandler(aiReportService)

	// 缺陷管理相关Handler
	defectHandler := handlers.NewDefectHandler(defectService, projectRepo)
	defectAttachmentHandler := handlers.NewDefectAttachmentHandler(defectAttachmentService)
	defectConfigHandler := handlers.NewDefectConfigHandler(defectConfigService)
	defectCommentHandler := handlers.NewDefectCommentHandler(defectCommentService)

	// 原始需求文档相关Handler (T48)
	rawDocumentHandler := handlers.NewRawDocumentHandler(rawDocumentService)

	// T51: 提示词管理相关Handler
	promptHandler := handlers.NewPromptHandler(promptService)

	// 用户自定义变量相关Handler
	userDefinedVarHandler := handlers.NewUserDefinedVariableHandler(userDefinedVarService)

	// 脚本测试服务和Handler
	scriptTestService := services.NewScriptTestService(userDefinedVarService)
	scriptTestHandler := handlers.NewScriptTestHandler(scriptTestService)

	// 初始化管理员账号
	if err := authService.InitAdminUsers(); err != nil {
		log.Printf("warning: failed to init admin users: %v", err)
	} else {
		log.Println("admin users initialized successfully")
	}

	// 创建 Gin 路由引擎
	r := gin.Default()

	// 配置 CORS - 允许局域网访问
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
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
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.AbortWithStatus(204)
	})

	// 注册路由
	api := r.Group("/api/v1")
	{
		// 健康检查接口（无需认证，供 Docker 健康检查使用）
		api.GET("/health", func(c *gin.Context) {
			// 检查数据库连接
			sqlDB, err := db.DB()
			if err != nil {
				c.JSON(503, gin.H{
					"status": "error",
					"db":     "disconnected",
					"error":  err.Error(),
				})
				return
			}
			if err := sqlDB.Ping(); err != nil {
				c.JSON(503, gin.H{
					"status": "error",
					"db":     "disconnected",
					"error":  err.Error(),
				})
				return
			}
			c.JSON(200, gin.H{
				"status": "ok",
				"db":     "connected",
			})
		})

		// 1. 公开路由(无需认证)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 公开的提示词路由（系统+全员提示词，不需要认证）
		publicPrompts := api.Group("/prompts/public")
		{
			publicPrompts.GET("", promptHandler.ListPrompts)
		}

		// T50: MCP双模式认证路由 - 支持JWT和API Token
		// 这些路由用于MCP服务器访问，需要支持X-API-Token认证
		mcpAuth := api.Group("")
		mcpAuth.Use(middleware.DualAuthMiddleware(authService, userService))
		{
			// T50: MCP Token验证接口 - 返回当前用户信息
			mcpAuth.GET("/auth/me", profileHandler.GetProfile)
		}

		// 2. 认证路由(需要登录,但不限角色) - 使用双模式认证支持MCP
		authenticated := api.Group("")
		authenticated.Use(middleware.DualAuthMiddleware(authService, userService))
		{
			authenticated.POST("/auth/logout", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "logout successful"})
			})
			authenticated.GET("/profile", profileHandler.GetProfile)
			authenticated.PUT("/profile/nickname", profileHandler.UpdateNickname)
			// T23 密码修改与Token功能
			authenticated.PUT("/profile/password", profileHandler.UpdatePassword)
			authenticated.POST("/profile/token", profileHandler.GenerateToken)
			authenticated.GET("/profile/token", profileHandler.GetToken)
			// T50 当前项目管理
			authenticated.PUT("/profile/current-project", profileHandler.SetCurrentProject)
			authenticated.GET("/profile/current-project", profileHandler.GetCurrentProject)
			// MCP工具使用：获取当前项目信息（无权限检查）
			authenticated.GET("/profile/current-project-info", profileHandler.GetCurrentProjectInfo)

			// 手工测试用例模版导出（全局路由，不依赖项目）
			authenticated.GET("/manual-cases/template", versionHandler.ExportTemplate)
		}

		// 3. 项目管理路由(PM + PMemb 可查看,PM 可操作)
		// T50: 使用DualAuthMiddleware支持MCP的API Token访问
		projects := api.Group("/projects")
		projects.Use(middleware.DualAuthMiddleware(authService, userService))
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

			// 需求条目管理路由 (T42 - 新架构)
			projects.GET("/:id/requirement-items",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.ListItems)
			projects.GET("/:id/requirement-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.GetItem)
			projects.POST("/:id/requirement-items",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.CreateItem)
			projects.POST("/:id/requirement-items/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.BulkCreateItems)
			projects.POST("/:id/requirement-items/export",
				func(c *gin.Context) {
					log.Printf("=== [路由层] 收到POST请求: /api/v1/projects/%s/requirement-items/export ===", c.Param("id"))
					c.Next()
				},
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.ExportToZip)
			projects.POST("/:id/requirement-items/import",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.ImportFromZip)
			projects.PUT("/:id/requirement-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.UpdateItem)
			projects.DELETE("/:id/requirement-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.DeleteItem)
			projects.PUT("/:id/requirement-items/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.BulkUpdateItems)
			projects.DELETE("/:id/requirement-items/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.BulkDeleteItems)

			// 需求Chunk管理路由 (T54)
			projects.GET("/:id/requirement-items/:itemId/chunks",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementChunkHandler.ListChunks)
			projects.POST("/:id/requirement-items/:itemId/chunks",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementChunkHandler.CreateChunk)
			projects.PUT("/:id/requirement-items/:itemId/chunks/reorder",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementChunkHandler.ReorderChunks)

			// 观点条目管理路由 (T42 - 新架构)
			projects.GET("/:id/viewpoint-items",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.ListItems)
			projects.GET("/:id/viewpoint-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.GetItem)
			projects.POST("/:id/viewpoint-items",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.CreateItem)
			projects.POST("/:id/viewpoint-items/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.BulkCreateItems)
			projects.POST("/:id/viewpoint-items/export",
				func(c *gin.Context) {
					log.Printf("=== [路由层] 收到POST请求: /api/v1/projects/%s/viewpoint-items/export ===", c.Param("id"))
					c.Next()
				},
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.ExportToZip)
			projects.POST("/:id/viewpoint-items/import",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.ImportFromZip)
			projects.PUT("/:id/viewpoint-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.UpdateItem)
			projects.DELETE("/:id/viewpoint-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.DeleteItem)
			projects.PUT("/:id/viewpoint-items/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.BulkUpdateItems)
			projects.DELETE("/:id/viewpoint-items/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.BulkDeleteItems)

			// 观点Chunk管理路由 (T54)
			projects.GET("/:id/viewpoint-items/:itemId/chunks",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointChunkHandler.ListChunks)
			projects.POST("/:id/viewpoint-items/:itemId/chunks",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointChunkHandler.CreateChunk)
			projects.PUT("/:id/viewpoint-items/:itemId/chunks/reorder",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointChunkHandler.ReorderChunks)

			// 用例集管理路由
			projects.GET("/:id/case-groups",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				caseGroupHandler.GetCaseGroups)
			projects.POST("/:id/case-groups",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				caseGroupHandler.CreateCaseGroup)

			// 用户自定义变量路由
			projects.GET("/:id/case-groups/:groupId/variables",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				userDefinedVarHandler.GetVariables)
			projects.PUT("/:id/case-groups/:groupId/variables",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				userDefinedVarHandler.SaveVariables)
			projects.POST("/:id/case-groups/:groupId/variables",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				userDefinedVarHandler.AddVariable)
			projects.PUT("/:id/case-groups/:groupId/variables/:varId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				userDefinedVarHandler.UpdateVariable)
			projects.DELETE("/:id/case-groups/:groupId/variables/:varId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				userDefinedVarHandler.DeleteVariable)

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
			projects.POST("/:id/manual-cases/save-version",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.SaveMultiLangVersion)

			// 为特定用例集创建手工用例（MCP工具使用）
			projects.POST("/:id/case-groups/:groupId/manual-cases",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.CreateCaseForGroup)
			// 为特定用例集更新手工用例（MCP工具使用）
			projects.PUT("/:id/case-groups/:groupId/manual-cases/:caseId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				manualCasesHandler.UpdateCaseForGroup)

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

			// Web用例模版导入导出路由
			projects.GET("/:id/web-cases/template",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.ExportWebTemplate)
			projects.POST("/:id/web-cases/import",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				autoCasesHandler.ImportWebCases)

			// T45: Web用例版本管理路由
			projects.POST("/:id/web-cases/versions",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				webVersionHandler.SaveWebVersion)
			projects.GET("/:id/web-cases/versions",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				webVersionHandler.GetWebVersionList)
			projects.GET("/:id/web-cases/versions/:versionId/export",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				webVersionHandler.DownloadWebVersion)
			projects.DELETE("/:id/web-cases/versions/:versionId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				webVersionHandler.DeleteWebVersion)
			projects.PUT("/:id/web-cases/versions/:versionId/remark",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				webVersionHandler.UpdateWebVersionRemark)

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

			// 接口测试用例集管理路由
			projects.GET("/:id/api-case-groups",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCaseGroupHandler.GetCaseGroups)
			projects.POST("/:id/api-case-groups",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCaseGroupHandler.CreateCaseGroup)

			// 接口测试用例模版下载和导入路由
			projects.GET("/:id/api-cases/template",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.ExportTemplate)
			projects.POST("/:id/api-cases/import",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCasesHandler.ImportCases)

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
		}

		// 接口测试用例集路由（不依赖项目ID）
		apiCaseGroups := api.Group("/api-case-groups")
		apiCaseGroups.Use(middleware.AuthMiddleware(authService))
		{
			apiCaseGroups.PUT("/:groupId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCaseGroupHandler.UpdateCaseGroup)
			apiCaseGroups.DELETE("/:groupId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				apiCaseGroupHandler.DeleteCaseGroup)
		}

		// 版本管理路由等 - 继续使用DualAuthMiddleware保证MCP兼容性
		{
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

			// 审阅条目路由 (T44)
			projects.GET("/:id/review-items",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				reviewItemHandler.ListReviewItems)
			projects.POST("/:id/review-items",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				reviewItemHandler.CreateReviewItem)
			projects.GET("/:id/review-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				reviewItemHandler.GetReviewItem)
			projects.PUT("/:id/review-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				reviewItemHandler.UpdateReviewItem)
			projects.DELETE("/:id/review-items/:itemId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				reviewItemHandler.DeleteReviewItem)
			projects.GET("/:id/review-items/:itemId/download",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				reviewItemHandler.DownloadReviewItem)

			// AI质量报告路由 (T47)
			projects.GET("/:id/ai-reports",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				aiReportHandler.ListReports)
			projects.POST("/:id/ai-reports",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				aiReportHandler.CreateReport)
			projects.GET("/:id/ai-reports/:reportId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				aiReportHandler.GetReport)
			projects.PUT("/:id/ai-reports/:reportId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				aiReportHandler.UpdateReport)
			projects.DELETE("/:id/ai-reports/:reportId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				aiReportHandler.DeleteReport)

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
			projects.POST("/:id/execution-tasks/:task_uuid/execute",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				executionTaskHandler.ExecuteTask)
			projects.POST("/:id/execution-tasks/:task_uuid/cases/:case_result_id/execute",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				executionTaskHandler.ExecuteSingleCase)

			// 执行任务变量路由
			projects.GET("/:id/execution-tasks/:task_uuid/variables",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				userDefinedVarHandler.GetTaskVariables)
			projects.PUT("/:id/execution-tasks/:task_uuid/variables",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				userDefinedVarHandler.SaveTaskVariables)

			// 脚本测试路由 - 用于用例详情页面的脚本测试功能
			projects.POST("/:id/script-test",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				scriptTestHandler.TestScript)

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
		executionGroup.Use(middleware.DualAuthMiddleware(authService, userService))
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

		// 执行用例结果单独更新路由 (供MCP工具使用)
		executionCaseGroup := api.Group("/execution-task-cases")
		executionCaseGroup.Use(middleware.DualAuthMiddleware(authService, userService))
		{
			executionCaseGroup.PUT("/batch",
				executionCaseResultHandler.BatchUpdateCaseResults)
			executionCaseGroup.PUT("/:id",
				executionCaseResultHandler.UpdateSingleCaseResult)
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

		// 原始需求文档管理路由 (T48)
		// 项目级别路由
		projects.POST("/:id/raw-documents",
			rawDocumentHandler.Upload)
		projects.GET("/:id/raw-documents",
			rawDocumentHandler.List)

		// 文档级别路由 - 不需要认证的操作
		rawDocumentsPublic := api.Group("/raw-documents")
		{
			rawDocumentsPublic.GET("/:id/convert-status",
				rawDocumentHandler.GetConvertStatus)
			rawDocumentsPublic.GET("/:id/download",
				rawDocumentHandler.DownloadOriginal)
			rawDocumentsPublic.GET("/:id/converted/download",
				rawDocumentHandler.DownloadConverted)
			rawDocumentsPublic.GET("/:id/converted/preview",
				rawDocumentHandler.PreviewConverted)
		}

		// 文档级别路由 - 需要认证的操作
		rawDocumentsAuth := api.Group("/raw-documents")
		rawDocumentsAuth.Use(middleware.AuthMiddleware(authService))
		{
			rawDocumentsAuth.POST("/:id/convert",
				rawDocumentHandler.Convert)
			rawDocumentsAuth.DELETE("/:id",
				rawDocumentHandler.DeleteOriginal)
			rawDocumentsAuth.DELETE("/:id/converted",
				rawDocumentHandler.DeleteConverted)
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

		// 需求条目单项CRUD路由 (T42)
		requirementItems := api.Group("/requirement-items")
		requirementItems.Use(middleware.AuthMiddleware(authService))
		{
			requirementItems.GET("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.GetItem)
			requirementItems.PUT("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.UpdateItem)
			requirementItems.DELETE("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.DeleteItem)
			requirementItems.PUT("/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.BulkUpdateItems)
			requirementItems.DELETE("/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementItemHandler.BulkDeleteItems)
		}

		// 观点条目单项CRUD路由 (T42)
		viewpointItems := api.Group("/viewpoint-items")
		viewpointItems.Use(middleware.AuthMiddleware(authService))
		{
			viewpointItems.GET("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.GetItem)
			viewpointItems.PUT("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.UpdateItem)
			viewpointItems.DELETE("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.DeleteItem)
			viewpointItems.PUT("/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.BulkUpdateItems)
			viewpointItems.DELETE("/bulk",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointItemHandler.BulkDeleteItems)
		}

		// 需求Chunk独立路由 (T54)
		requirementChunks := api.Group("/requirement-chunks")
		requirementChunks.Use(middleware.DualAuthMiddleware(authService, userService))
		{
			requirementChunks.GET("/:chunkId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementChunkHandler.GetChunk)
			requirementChunks.PUT("/:chunkId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementChunkHandler.UpdateChunk)
			requirementChunks.DELETE("/:chunkId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				requirementChunkHandler.DeleteChunk)
		}

		// 观点Chunk独立路由 (T54)
		viewpointChunks := api.Group("/viewpoint-chunks")
		viewpointChunks.Use(middleware.DualAuthMiddleware(authService, userService))
		{
			viewpointChunks.GET("/:chunkId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointChunkHandler.GetChunk)
			viewpointChunks.PUT("/:chunkId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointChunkHandler.UpdateChunk)
			viewpointChunks.DELETE("/:chunkId",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				viewpointChunkHandler.DeleteChunk)
		}

		// 用例集通用路由（更新和删除）
		// T50: 使用DualAuthMiddleware支持MCP的API Token访问
		caseGroups := api.Group("/case-groups")
		caseGroups.Use(middleware.DualAuthMiddleware(authService, userService))
		{
			caseGroups.GET("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				caseGroupHandler.GetCaseGroup)
			caseGroups.PUT("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				caseGroupHandler.UpdateCaseGroup)
			caseGroups.DELETE("/:id",
				middleware.RequireRole(constants.RoleProjectManager, constants.RoleProjectMember),
				caseGroupHandler.DeleteCaseGroup)
		}

		// 4. 用户管理路由（系统管理员和项目管理员共用，内部分流）
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
			projectHandler.GetProjectMembers)
		api.PUT("/projects/:id/members",
			middleware.AuthMiddleware(authService),
			middleware.RequireRole(constants.RoleProjectManager),
			projectHandler.UpdateProjectMembers)

		// T51: 提示词管理路由
		// 需要认证的提示词路由（个人提示词）
		prompts := api.Group("/prompts")
		prompts.Use(middleware.DualAuthMiddleware(authService, userService))
		{
			// ListPrompts 在handler内部处理权限检查
			// 有Token时可获取个人提示词
			prompts.GET("", promptHandler.ListPrompts)
			prompts.GET("/by-name", promptHandler.GetPromptByName) // MCP通过名称获取提示词
			prompts.GET("/:id",
				middleware.RequireRole(constants.RoleSystemAdmin, constants.RoleProjectManager, constants.RoleProjectMember),
				promptHandler.GetPromptByID)
			prompts.POST("",
				middleware.RequireRole(constants.RoleSystemAdmin, constants.RoleProjectManager, constants.RoleProjectMember),
				promptHandler.CreatePrompt)
			prompts.PUT("/:id",
				middleware.RequireRole(constants.RoleSystemAdmin, constants.RoleProjectManager, constants.RoleProjectMember),
				promptHandler.UpdatePrompt)
			prompts.DELETE("/:id",
				middleware.RequireRole(constants.RoleSystemAdmin, constants.RoleProjectManager, constants.RoleProjectMember),
				promptHandler.DeletePrompt)
		}
	}

	// 服务前端静态文件
	frontendBuildPath := getFrontendBuildPath()
	log.Printf("Serving frontend from: %s", frontendBuildPath)
	r.Static("/static", filepath.Join(frontendBuildPath, "static"))

	// 对于所有非 API 路由，返回 index.html (支持前端路由)
	indexPath := filepath.Join(frontendBuildPath, "index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File(indexPath)
	})

	// 启动 HTTPS 服务器
	// 证书路径通过环境变量配置，默认在项目根目录的 certs/ 下
	certFile := config.GetCertFilePath()
	keyFile := config.GetKeyFilePath()
	serverAddr := config.GetServerAddr()

	log.Printf("server starting on https://localhost%s (cert: %s)", serverAddr, certFile)
	if err := r.RunTLS(serverAddr, certFile, keyFile); err != nil {
		log.Fatalf("failed to start HTTPS server: %v", err)
	}
}

// initSystemPrompts 从prompts目录动态加载系统提示词到数据库
func initSystemPrompts(db *gorm.DB) error {
	projectID := uint(1)
	promptsDir := config.GetPromptsDir()

	// 记录详细的路径信息用于调试
	log.Printf("[INFO] Initializing system prompts from directory: %s", promptsDir)

	// 检查目录是否存在
	if stat, err := os.Stat(promptsDir); err != nil {
		log.Printf("[ERROR] Prompts directory does not exist or cannot be accessed: %s", promptsDir)
		log.Printf("[ERROR] Error details: %v", err)

		// 尝试打印当前工作目录和环境变量，帮助调试
		if cwd, err := os.Getwd(); err == nil {
			log.Printf("[DEBUG] Current working directory: %s", cwd)
		}
		log.Printf("[DEBUG] PROMPTS_DIR env: %s", os.Getenv("PROMPTS_DIR"))

		return fmt.Errorf("prompts directory not found: %s", promptsDir)
	} else if !stat.IsDir() {
		log.Printf("[ERROR] %s exists but is not a directory", promptsDir)
		return fmt.Errorf("%s is not a directory", promptsDir)
	}

	// 扫描目录中所有 .prompt.md 文件
	files, err := filepath.Glob(filepath.Join(promptsDir, "*.prompt.md"))
	if err != nil {
		return fmt.Errorf("scan prompts directory: %w", err)
	}

	log.Printf("[INFO] Found %d prompt files in %s", len(files), promptsDir)

	// 记录当前文件中的所有提示词名称（用于删除检测）
	filePromptNames := make(map[string]bool)

	for _, filePath := range files {
		// 检查文件是否为空
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("[WARN] Failed to stat file %s: %v", filePath, err)
			continue
		}
		if fileInfo.Size() == 0 {
			log.Printf("[SKIP] Empty file: %s", filePath)
			continue
		}

		// 读取文件内容
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("[WARN] Failed to read prompt file %s: %v", filePath, err)
			continue
		}

		// 解析YAML front matter
		frontMatter, err := parsePromptFrontMatter(string(content), filePath)
		if err != nil {
			log.Printf("[WARN] Failed to parse prompt file %s: %v", filePath, err)
			continue
		}

		filePromptNames[frontMatter.Name] = true

		// 使用 Upsert 方式：先尝试更新，不存在则创建
		var existingPrompt models.Prompt
		result := db.Where("name = ? AND scope = ? AND project_id = ?", frontMatter.Name, "system", projectID).First(&existingPrompt)

		if result.Error == gorm.ErrRecordNotFound {
			// 创建新的提示词记录
			prompt := &models.Prompt{
				ProjectID:   projectID,
				Name:        frontMatter.Name,
				Description: frontMatter.Description,
				Version:     frontMatter.Version,
				Content:     string(content),
				Scope:       "system",
				Arguments:   "[]", // 必须是有效的JSON，PostgreSQL的json类型不接受空字符串
				UserID:      nil,
			}

			if err := db.Create(prompt).Error; err != nil {
				log.Printf("[WARN] Failed to create system prompt %s: %v", frontMatter.Name, err)
				continue
			}

			log.Printf("[INFO] Created system prompt: %s (version %s)", frontMatter.Name, frontMatter.Version)
		} else if result.Error == nil {
			// 更新现有的提示词记录
			if err := db.Model(&existingPrompt).Updates(map[string]interface{}{
				"description": frontMatter.Description,
				"version":     frontMatter.Version,
				"content":     string(content),
			}).Error; err != nil {
				log.Printf("[WARN] Failed to update system prompt %s: %v", frontMatter.Name, err)
				continue
			}

			log.Printf("[INFO] Updated system prompt: %s (version %s)", frontMatter.Name, frontMatter.Version)
		} else {
			log.Printf("[WARN] Error checking prompt %s: %v", frontMatter.Name, result.Error)
			continue
		}
	}

	// 删除数据库中不再需要的系统提示词
	var existingPrompts []models.Prompt
	if err := db.Where("scope = ? AND project_id = ?", "system", projectID).Find(&existingPrompts).Error; err != nil {
		log.Printf("[WARN] Failed to query existing system prompts: %v", err)
	} else {
		for _, p := range existingPrompts {
			if !filePromptNames[p.Name] {
				if err := db.Delete(&p).Error; err != nil {
					log.Printf("[WARN] Failed to delete obsolete system prompt %s: %v", p.Name, err)
				} else {
					log.Printf("[INFO] Deleted obsolete system prompt: %s", p.Name)
				}
			}
		}
	}

	log.Printf("[INFO] System prompts initialized. Active: %d", len(filePromptNames))
	return nil
}

// PromptFrontMatter 提示词文件的YAML头部结构
type PromptFrontMatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// parsePromptFrontMatter 解析提示词文件的YAML front matter
func parsePromptFrontMatter(content string, filePath string) (*PromptFrontMatter, error) {
	// 使用正则匹配 YAML front matter (---开头和结尾)
	frontMatterRegex := regexp.MustCompile(`(?s)^---\s*\n(.+?)\n---`)
	matches := frontMatterRegex.FindStringSubmatch(content)

	baseName := filepath.Base(filePath)
	nameFromFile := strings.TrimSuffix(baseName, ".prompt.md")

	if len(matches) < 2 {
		// 没有front matter，使用文件名作为name
		return &PromptFrontMatter{
			Name:        nameFromFile,
			Description: "",
			Version:     "1.0",
		}, nil
	}

	// 解析YAML
	var frontMatter PromptFrontMatter
	if err := yaml.Unmarshal([]byte(matches[1]), &frontMatter); err != nil {
		return nil, fmt.Errorf("parse yaml front matter: %w", err)
	}

	// 如果name为空，使用文件名
	if frontMatter.Name == "" {
		frontMatter.Name = nameFromFile
	}

	// 如果version为空，默认1.0
	if frontMatter.Version == "" {
		frontMatter.Version = "1.0"
	}

	return &frontMatter, nil
}
