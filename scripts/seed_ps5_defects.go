//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// PS5测试相关的模块
var ps5Modules = []string{
	"Dashboard UI",
	"Game Library",
	"PlayStation Store",
	"Settings",
	"Network",
	"Controller",
	"Audio System",
	"Video Output",
	"Storage Management",
	"User Account",
	"Trophy System",
	"Share/Capture",
	"Party Chat",
	"Remote Play",
	"VR Support",
}

// 缺陷类型
var defectTypes = []string{
	"Functional",
	"UI",
	"UIInteraction",
	"Compatibility",
	"Performance",
	"Security",
	"Environment",
}

// 严重程度
var severities = []string{
	"Critical",
	"Major",
	"Minor",
	"Trivial",
}

// 优先级
var priorities = []string{
	"A",
	"B",
	"C",
	"D",
}

// 状态
var statuses = []string{
	"New",
	"InProgress",
	"Confirmed",
	"Resolved",
	"Closed",
	"Reopened",
}

// 测试阶段
var phases = []string{
	"UT",
	"IT",
	"ST",
	"UAT",
	"Regression",
}

// 版本
var versions = []string{
	"v1.0.0",
	"v1.0.1",
	"v1.1.0",
	"v1.2.0",
	"v2.0.0-beta",
	"v2.0.0",
}

// 组件
var components = []string{
	"Frontend",
	"Backend",
	"Firmware",
	"Network Stack",
	"Graphics Engine",
	"Audio Engine",
	"Storage Driver",
}

// 检测团队
var detectionTeams = []string{
	"QA-Tokyo",
	"QA-Shanghai",
	"QA-California",
	"Dev-Test",
	"Beta Tester",
}

// 检测人
var detectedByUsers = []string{
	"Tanaka",
	"Yamamoto",
	"Zhang Wei",
	"Li Ming",
	"John Smith",
	"Emily Chen",
	"Sato",
	"Kim",
}

// 指派人
var assignees = []string{
	"dev_tanaka",
	"dev_suzuki",
	"dev_wang",
	"dev_chen",
	"dev_johnson",
	"dev_kim",
}

// 机型
var models = []string{
	"PS5 Standard Edition",
	"PS5 Digital Edition",
	"PS5 Slim",
	"PS5 Pro",
}

// 复现频率
var frequencies = []string{
	"Always",
	"Often",
	"Sometimes",
	"Rarely",
	"Once",
}

// Bug标题模板
var titleTemplates = []string{
	"[%s] %s crashes when %s",
	"[%s] Cannot %s in %s mode",
	"[%s] %s display issue on %s",
	"[%s] Performance drop when %s",
	"[%s] %s button not responding",
	"[%s] Audio glitch during %s",
	"[%s] Network error when %s",
	"[%s] %s freezes after %s",
	"[%s] Incorrect %s in %s screen",
	"[%s] Memory leak in %s",
}

var actions = []string{
	"loading game",
	"starting application",
	"connecting to server",
	"downloading content",
	"updating system",
	"switching user",
	"entering rest mode",
	"resuming from suspend",
	"capturing screenshot",
	"streaming gameplay",
	"voice chat active",
	"installing game",
	"browsing store",
	"syncing trophies",
	"checking notifications",
}

var objects = []string{
	"Home screen",
	"Store page",
	"Game menu",
	"Settings panel",
	"Trophy list",
	"Friends list",
	"Download queue",
	"Capture gallery",
	"User profile",
	"Party interface",
}

func main() {
	db, err := sql.Open("sqlite3", "./webtest.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	projectID := 35 // CZ项目
	userID := 1     // 默认用户

	rand.Seed(time.Now().UnixNano())

	// 获取当前最大的defect_id序号
	var maxSeq int
	err = db.QueryRow("SELECT COALESCE(MAX(CAST(defect_id AS INTEGER)), 0) FROM defects WHERE project_id = ?", projectID).Scan(&maxSeq)
	if err != nil {
		log.Printf("Warning: Could not get max defect_id: %v, starting from 1", err)
		maxSeq = 0
	}

	log.Printf("Starting from defect_id sequence: %d", maxSeq+1)

	// 生成100个Bug
	startDate := time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2026, 2, 6, 0, 0, 0, 0, time.Local)
	dateRange := endDate.Sub(startDate)

	stmt, err := db.Prepare(`
		INSERT INTO defects (
			id, defect_id, project_id, title, subject, description, 
			recovery_method, priority, severity, type, frequency,
			detected_version, phase, case_id, assignee, recovery_rank,
			detection_team, location, fix_version, sqa_memo, component,
			resolution, models, detected_by, status, created_by, updated_by,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	successCount := 0
	for i := 1; i <= 100; i++ {
		defectSeq := maxSeq + i
		defectID := fmt.Sprintf("%06d", defectSeq)
		id := uuid.New().String()

		// 随机选择各字段值
		module := ps5Modules[rand.Intn(len(ps5Modules))]
		defectType := defectTypes[rand.Intn(len(defectTypes))]
		severity := severities[rand.Intn(len(severities))]
		priority := priorities[rand.Intn(len(priorities))]
		status := statuses[rand.Intn(len(statuses))]
		phase := phases[rand.Intn(len(phases))]
		version := versions[rand.Intn(len(versions))]
		component := components[rand.Intn(len(components))]
		team := detectionTeams[rand.Intn(len(detectionTeams))]
		detectedBy := detectedByUsers[rand.Intn(len(detectedByUsers))]
		assignee := assignees[rand.Intn(len(assignees))]
		model := models[rand.Intn(len(models))]
		frequency := frequencies[rand.Intn(len(frequencies))]

		// 生成标题
		titleTemplate := titleTemplates[rand.Intn(len(titleTemplates))]
		action := actions[rand.Intn(len(actions))]
		object := objects[rand.Intn(len(objects))]
		title := fmt.Sprintf(titleTemplate, module, object, action)

		// 生成描述
		description := fmt.Sprintf(`## 问题描述
在使用 %s 功能时，%s 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 %s 菜单
3. 执行 %s 操作
4. 观察到问题发生

## 预期结果
功能应正常工作，无错误提示。

## 实际结果
出现异常行为，影响用户体验。

## 测试环境
- 机型: %s
- 固件版本: %s
- 网络: 有线连接`, module, object, module, action, model, version)

		// 生成恢复方法
		recoveryMethod := fmt.Sprintf("重启主机或重新进入%s功能", module)

		// 随机生成创建时间
		randomDuration := time.Duration(rand.Int63n(int64(dateRange)))
		createdAt := startDate.Add(randomDuration)
		updatedAt := createdAt.Add(time.Duration(rand.Intn(72)) * time.Hour)

		// 修复版本（只有Resolved或Closed状态才有）
		fixVersion := ""
		if status == "Resolved" || status == "Closed" {
			versionIdx := rand.Intn(len(versions))
			fixVersion = versions[versionIdx]
		}

		// 解决方案（只有Resolved或Closed状态才有）
		resolution := ""
		if status == "Resolved" || status == "Closed" {
			resolutions := []string{
				"Fixed in firmware update",
				"Code patch applied",
				"Configuration issue resolved",
				"Hardware compatibility issue addressed",
				"Network protocol updated",
			}
			resolution = resolutions[rand.Intn(len(resolutions))]
		}

		// SQA备注
		sqaMemo := ""
		if rand.Float32() > 0.5 {
			memos := []string{
				"Verified fix in latest build",
				"Needs further testing",
				"Related to issue #" + fmt.Sprintf("%d", rand.Intn(1000)),
				"Priority escalated by PM",
				"Customer reported issue",
			}
			sqaMemo = memos[rand.Intn(len(memos))]
		}

		// 位置信息
		location := fmt.Sprintf("%s > %s", module, object)

		// Case ID
		caseID := fmt.Sprintf("TC-PS5-%04d", rand.Intn(9999))

		// 恢复等级
		recoveryRanks := []string{"S", "A", "B", "C"}
		recoveryRank := recoveryRanks[rand.Intn(len(recoveryRanks))]

		_, err := stmt.Exec(
			id, defectID, projectID, title, module, description,
			recoveryMethod, priority, severity, defectType, frequency,
			version, phase, caseID, assignee, recoveryRank,
			team, location, fixVersion, sqaMemo, component,
			resolution, model, detectedBy, status, userID, userID,
			createdAt, updatedAt,
		)
		if err != nil {
			log.Printf("Error inserting defect %s: %v", defectID, err)
			continue
		}
		successCount++

		if i%10 == 0 {
			log.Printf("Inserted %d defects...", i)
		}
	}

	log.Printf("Successfully inserted %d PS5 test defects into CZ project (project_id=%d)", successCount, projectID)
}
