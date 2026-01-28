package services

import (
	"fmt"
	"webtest/internal/models"
	"webtest/internal/repositories"

	"github.com/google/uuid"
)

// AIReportService AI报告服务接口
type AIReportService interface {
	ListReports(projectID uint, reportType string) ([]*models.AIReport, error)
	CreateReport(projectID uint, reportType, name string) (*models.AIReport, error)
	GetReport(reportID string) (*models.AIReport, error)
	UpdateReport(reportID string, name, content *string) (*models.AIReport, error)
	DeleteReport(reportID string) error
}

// aiReportService AI报告服务实现
type aiReportService struct {
	repo repositories.AIReportRepository
}

// NewAIReportService 创建AI报告服务实例
func NewAIReportService(repo repositories.AIReportRepository) AIReportService {
	return &aiReportService{repo: repo}
}

// ListReports 获取项目报告列表(按类型过滤)
func (s *aiReportService) ListReports(projectID uint, reportType string) ([]*models.AIReport, error) {
	var reports []*models.AIReport
	var err error

	if reportType != "" {
		reports, err = s.repo.FindByProjectIDAndType(projectID, reportType)
	} else {
		reports, err = s.repo.FindByProjectID(projectID)
	}

	if err != nil {
		return nil, fmt.Errorf("获取报告列表失败: %w", err)
	}
	if reports == nil {
		reports = make([]*models.AIReport, 0)
	}
	return reports, nil
}

// CreateReport 创建报告(包含名称检重,支持类型)
func (s *aiReportService) CreateReport(projectID uint, reportType, name string) (*models.AIReport, error) {
	// 验证类型,默认为O(其他)
	if reportType == "" {
		reportType = "O"
	}
	validTypes := map[string]bool{"R": true, "A": true, "T": true, "O": true}
	if !validTypes[reportType] {
		return nil, fmt.Errorf("无效的报告类型: %s", reportType)
	}

	// 检查同类型下名称是否重复
	exists, err := s.repo.ExistsByProjectTypeAndName(projectID, reportType, name, "")
	if err != nil {
		return nil, fmt.Errorf("检查报告名称重复失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("报告名称已存在")
	}

	// 生成ID: report_前缀加UUID(简化版雪花ID)
	reportID := "report_" + uuid.New().String()

	report := &models.AIReport{
		ID:        reportID,
		ProjectID: projectID,
		Type:      reportType,
		Name:      name,
		Content:   "", // 初始化为空
	}

	if err := s.repo.Create(report); err != nil {
		return nil, fmt.Errorf("创建报告失败: %w", err)
	}

	return report, nil
}

// GetReport 获取报告详情
func (s *aiReportService) GetReport(reportID string) (*models.AIReport, error) {
	report, err := s.repo.FindByID(reportID)
	if err != nil {
		return nil, fmt.Errorf("报告不存在: %w", err)
	}
	return report, nil
}

// UpdateReport 更新报告(支持可选更新name和content)
func (s *aiReportService) UpdateReport(reportID string, name, content *string) (*models.AIReport, error) {
	// 获取现有报告
	report, err := s.repo.FindByID(reportID)
	if err != nil {
		return nil, fmt.Errorf("报告不存在: %w", err)
	}

	// 如果更新名称,检查重名(排除当前报告ID)
	if name != nil && *name != report.Name {
		exists, err := s.repo.ExistsByProjectAndName(report.ProjectID, *name, reportID)
		if err != nil {
			return nil, fmt.Errorf("检查报告名称重复失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("报告名称已存在")
		}
		report.Name = *name
	}

	// 更新内容
	if content != nil {
		report.Content = *content
	}

	if err := s.repo.Update(report); err != nil {
		return nil, fmt.Errorf("更新报告失败: %w", err)
	}

	return report, nil
}

// DeleteReport 删除报告(软删除)
func (s *aiReportService) DeleteReport(reportID string) error {
	if err := s.repo.Delete(reportID); err != nil {
		return fmt.Errorf("删除报告失败: %w", err)
	}
	return nil
}
