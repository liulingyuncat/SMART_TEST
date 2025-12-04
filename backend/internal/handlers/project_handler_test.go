package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"webtest/internal/models"
	"webtest/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectService Mock项目服务
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetUserProjects(userID uint, role string) ([]models.Project, error) {
	args := m.Called(userID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Project), args.Error(1)
}

func (m *MockProjectService) IsProjectMember(projectID uint, userID uint) (bool, error) {
	args := m.Called(projectID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectService) CreateProject(name string, description string, creatorID uint) (*models.Project, error) {
	args := m.Called(name, description, creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

// setupTestRouter 创建测试路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TestProjectHandler_GetProjects_Success 测试成功获取项目列表
func TestProjectHandler_GetProjects_Success(t *testing.T) {
	mockService := new(MockProjectService)
	handler := NewProjectHandler(mockService)

	// Mock返回项目列表
	projects := []models.Project{
		{ID: 1, Name: "项目1", Description: "描述1"},
		{ID: 2, Name: "项目2", Description: "描述2"},
	}
	mockService.On("GetUserProjects", uint(1), "project_manager").Return(projects, nil)

	// 创建测试路由
	router := setupTestRouter()
	router.GET("/projects", func(c *gin.Context) {
		// 模拟中间件注入的userID和role
		c.Set("userID", uint(1))
		c.Set("role", "project_manager")
		handler.GetProjects(c)
	})

	// 发送请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects", nil)
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(0), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestProjectHandler_CreateProject_Success 测试成功创建项目
func TestProjectHandler_CreateProject_Success(t *testing.T) {
	mockService := new(MockProjectService)
	handler := NewProjectHandler(mockService)

	// Mock返回新创建的项目
	newProject := &models.Project{
		ID:          1,
		Name:        "新项目",
		Description: "项目描述",
	}
	mockService.On("CreateProject", "新项目", "项目描述", uint(1)).Return(newProject, nil)

	// 创建测试路由
	router := setupTestRouter()
	router.POST("/projects", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.CreateProject(c)
	})

	// 准备请求体
	requestBody := map[string]string{
		"name":        "新项目",
		"description": "项目描述",
	}
	body, _ := json.Marshal(requestBody)

	// 发送请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(0), response["code"])

	mockService.AssertExpectations(t)
}

// TestProjectHandler_CreateProject_DuplicateName 测试创建重名项目
func TestProjectHandler_CreateProject_DuplicateName(t *testing.T) {
	mockService := new(MockProjectService)
	handler := NewProjectHandler(mockService)

	// Mock返回项目名已存在错误
	mockService.On("CreateProject", "已存在项目", "描述", uint(1)).
		Return(nil, services.ErrProjectNameExists)

	// 创建测试路由
	router := setupTestRouter()
	router.POST("/projects", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.CreateProject(c)
	})

	// 准备请求体
	requestBody := map[string]string{
		"name":        "已存在项目",
		"description": "描述",
	}
	body, _ := json.Marshal(requestBody)

	// 发送请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(400), response["code"])

	mockService.AssertExpectations(t)
}

// TestProjectHandler_CreateProject_InvalidRequest 测试无效请求
func TestProjectHandler_CreateProject_InvalidRequest(t *testing.T) {
	mockService := new(MockProjectService)
	handler := NewProjectHandler(mockService)

	// 创建测试路由
	router := setupTestRouter()
	router.POST("/projects", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.CreateProject(c)
	})

	// 准备无效请求体(缺少name字段)
	requestBody := map[string]string{
		"description": "描述",
	}
	body, _ := json.Marshal(requestBody)

	// 发送请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// 集成测试说明
// ==========================================
// 完整的Handler集成测试应包括:
//
// 1. 认证测试:
//    - TestGetProjects_NoAuth: 缺少Authorization头返回401
//    - TestGetProjects_InvalidToken: 无效Token返回401
//
// 2. 授权测试:
//    - TestGetProjects_SystemAdmin: 系统管理员返回空列表或403
//    - TestCreateProject_NotProjectManager: 非PM角色返回403
//
// 3. 参数验证测试:
//    - TestCreateProject_EmptyName: 空项目名返回400
//    - TestCreateProject_TooLongName: 项目名超长返回400
//
// 4. 业务逻辑测试:
//    - TestGetProjects_EmptyList: 用户无项目返回空数组
//    - TestCreateProject_ServerError: 服务层错误返回500
//
// 验收标准:
// - HTTP状态码正确
// - 响应格式符合统一规范
// - 错误信息明确
// - 覆盖率 >= 60%
//
// 实际项目中建议使用:
// - httpexpect 库简化HTTP测试
// - 真实AuthMiddleware集成测试
// - 端到端测试覆盖完整请求链路
