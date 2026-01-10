package handlers

import (
	"encoding/json"
	"testing"
)

func TestSortAPICases(t *testing.T) {
	// 构造测试数据：故意打乱顺序
	testCases := []interface{}{
		map[string]interface{}{
			"screen":   "[用户管理]",
			"url":      "/api/user",
			"method":   "DELETE",
			"response": `{"code": 200}`,
		},
		map[string]interface{}{
			"screen":   "[登录]",
			"url":      "/api/login",
			"method":   "POST",
			"response": `{"code": 401}`,
		},
		map[string]interface{}{
			"screen":   "[用户管理]",
			"url":      "/api/user",
			"method":   "POST",
			"response": `{"code": 201}`,
		},
		map[string]interface{}{
			"screen":   "[登录]",
			"url":      "/api/login",
			"method":   "POST",
			"response": `{"code": 200}`,
		},
		map[string]interface{}{
			"screen":   "[用户管理]",
			"url":      "/api/user",
			"method":   "GET",
			"response": `{"code": 200}`,
		},
		map[string]interface{}{
			"screen":   "[用户管理]",
			"url":      "/api/user",
			"method":   "PUT",
			"response": `{"code": 200}`,
		},
		map[string]interface{}{
			"screen":   "[用户管理]",
			"url":      "/api/user",
			"method":   "GET",
			"response": `{"code": 401}`,
		},
	}

	// 执行排序
	sorted := sortAPICases(testCases)

	// 验证结果
	t.Log("排序后的用例顺序:")
	for i, c := range sorted {
		data := c.(map[string]interface{})
		t.Logf("  %d: screen=%s, url=%s, method=%s, response=%s",
			i+1, data["screen"], data["url"], data["method"], data["response"])
	}

	// 期望顺序（按screen字母序 → url → method → responseCode）:
	// [用户管理] 的UTF-8编码在 [登录] 之前，所以用户管理排在前面
	// 1. [用户管理] /api/user GET 200
	// 2. [用户管理] /api/user GET 401
	// 3. [用户管理] /api/user POST 201
	// 4. [用户管理] /api/user PUT 200
	// 5. [用户管理] /api/user DELETE 200
	// 6. [登录] /api/login POST 200
	// 7. [登录] /api/login POST 401

	expected := []struct {
		screen   string
		method   string
		response string
	}{
		{"[用户管理]", "GET", `{"code": 200}`},
		{"[用户管理]", "GET", `{"code": 401}`},
		{"[用户管理]", "POST", `{"code": 201}`},
		{"[用户管理]", "PUT", `{"code": 200}`},
		{"[用户管理]", "DELETE", `{"code": 200}`},
		{"[登录]", "POST", `{"code": 200}`},
		{"[登录]", "POST", `{"code": 401}`},
	}

	if len(sorted) != len(expected) {
		t.Errorf("排序后用例数量不匹配: got %d, want %d", len(sorted), len(expected))
		return
	}

	for i, exp := range expected {
		data := sorted[i].(map[string]interface{})
		if data["screen"] != exp.screen {
			t.Errorf("用例%d screen不匹配: got %s, want %s", i+1, data["screen"], exp.screen)
		}
		if data["method"] != exp.method {
			t.Errorf("用例%d method不匹配: got %s, want %s", i+1, data["method"], exp.method)
		}
	}

	t.Log("✅ 排序测试通过!")
}

func TestExtractResponseCode(t *testing.T) {
	tests := []struct {
		name     string
		response string
		expected int
	}{
		{"JSON格式 code:200", `{"code": 200, "data": {}}`, 200},
		{"JSON格式 code:401", `{"code": 401, "message": "Unauthorized"}`, 401},
		{"前缀格式 200:", `200: {"data": {}}`, 200},
		{"前缀格式 404:", `404: {"message": "Not Found"}`, 404},
		{"空字符串", "", 200},
		{"无code字段", `{"message": "success"}`, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractResponseCode(tt.response)
			if got != tt.expected {
				t.Errorf("extractResponseCode(%q) = %d, want %d", tt.response, got, tt.expected)
			}
		})
	}
}

func TestGetMethodWeight(t *testing.T) {
	tests := []struct {
		method   string
		expected int
	}{
		{"GET", 1},
		{"get", 1},
		{"POST", 2},
		{"PUT", 3},
		{"PATCH", 4},
		{"DELETE", 5},
		{"OPTIONS", 99},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			got := getMethodWeight(tt.method)
			if got != tt.expected {
				t.Errorf("getMethodWeight(%q) = %d, want %d", tt.method, got, tt.expected)
			}
		})
	}
}

func TestSortAPICasesWithRealData(t *testing.T) {
	// 模拟实际生成的用例数据（类似P1comadmin的情况）
	realCases := []interface{}{
		map[string]interface{}{
			"screen":   "[ログイン履歴]",
			"url":      "/api/logininfo/list",
			"method":   "POST",
			"response": `{"code": 200}`,
		},
		map[string]interface{}{
			"screen":   "[ログイン履歴]",
			"url":      "/api/logininfo/list",
			"method":   "POST",
			"response": `{"code": 401}`,
		},
		map[string]interface{}{
			"screen":   "[ログイン履歴]",
			"url":      "/api/logininfo/export",
			"method":   "POST",
			"response": "Binary CSV file",
		},
		map[string]interface{}{
			"screen":   "[セットアップ履歴]",
			"url":      "/api/installhistory/list",
			"method":   "POST",
			"response": `{"code": 200}`,
		},
		map[string]interface{}{
			"screen":   "[セットアップ履歴]",
			"url":      "/api/installhistory/list",
			"method":   "POST",
			"response": `{"code": 401}`,
		},
	}

	sorted := sortAPICases(realCases)

	t.Log("实际数据排序结果:")
	for i, c := range sorted {
		data := c.(map[string]interface{})
		respJSON, _ := json.Marshal(data["response"])
		t.Logf("  %d: [%s] %s %s → %s",
			i+1, data["screen"], data["method"], data["url"], string(respJSON))
	}

	// 验证：同一screen的用例应该聚合在一起
	// [セットアップ履歴] 应该在 [ログイン履歴] 之前（按字母顺序）
	firstScreen := sorted[0].(map[string]interface{})["screen"].(string)
	if firstScreen != "[セットアップ履歴]" {
		t.Errorf("第一个用例应该是[セットアップ履歴], 实际是 %s", firstScreen)
	}
}
