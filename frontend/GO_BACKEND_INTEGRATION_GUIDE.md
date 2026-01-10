# Go后端与React+Recharts前端集成指南

## 概述

本文档说明如何在Go后端与React+Recharts前端协作中，确保数据格式正确且图表能正确显示。

---

## 📊 数据格式规范

### 1. 基础数据格式

所有用于Recharts图表的数据必须是**对象数组**格式：

```json
[
  { "field1": "value1", "field2": 100, "field3": 50 },
  { "field1": "value2", "field2": 150, "field3": 60 },
  { "field1": "value3", "field2": 120, "field3": 70 }
]
```

### 2. Go中的数据结构定义

#### 缺陷趋势数据
```go
// 缺陷趋势数据点
type DefectTrendData struct {
    Date   string `json:"date"`
    Total  int    `json:"total"`
    Closed int    `json:"closed"`
    Open   int    `json:"open"`
}

// 缺陷趋势响应
type DefectTrendResponse struct {
    Success bool                `json:"success"`
    Data    []DefectTrendData   `json:"data"`
    Message string              `json:"message"`
}
```

#### 其他常见图表数据
```go
// 柱状图数据
type BarChartData struct {
    Name     string  `json:"name"`
    Value    int     `json:"value"`
    Value2   int     `json:"value2"`
}

// 饼图数据
type PieChartData struct {
    Name  string `json:"name"`
    Value int    `json:"value"`
}

// 多维度数据（雷达图、其他）
type MultiDimensionData struct {
    Subject string `json:"subject"`
    MetricA int    `json:"metricA"`
    MetricB int    `json:"metricB"`
    FullMark int   `json:"fullMark"`
}
```

---

## 🔄 Go API实现示例

### 1. 缺陷趋势数据API

```go
package api

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
)

// GetDefectTrend 获取缺陷趋势数据
// @Summary 获取缺陷趋势
// @Description 返回指定项目的缺陷趋势数据（用于Recharts展示）
// @Tags defects
// @Param projectId path string true "项目ID"
// @Success 200 {object} DefectTrendResponse
// @Router /api/v1/projects/{projectId}/defects/trend [get]
func (c *Controller) GetDefectTrend(ctx *gin.Context) {
    projectId := ctx.Param("projectId")
    
    // 获取所有缺陷
    defects, err := c.defectService.GetAllDefects(projectId)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }
    
    // 生成趋势数据
    trendData := generateDefectTrend(defects)
    
    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": trendData,
    })
}

// generateDefectTrend 生成趋势数据
func generateDefectTrend(defects []Defect) []DefectTrendData {
    // 按日期分组统计
    trendMap := make(map[string]*DefectTrendData)
    
    for _, d := range defects {
        dateKey := d.CreatedAt.Format("2006-01-02")
        
        if _, exists := trendMap[dateKey]; !exists {
            trendMap[dateKey] = &DefectTrendData{
                Date:   dateKey,
                Total:  0,
                Closed: 0,
            }
        }
        
        trendMap[dateKey].Total++
        
        if d.Status == "Closed" {
            trendMap[dateKey].Closed++
        }
    }
    
    // 转换为排序的数组
    var result []DefectTrendData
    for date := minDate; date.Before(time.Now()); date = date.AddDate(0, 0, 1) {
        dateStr := date.Format("2006-01-02")
        if trend, exists := trendMap[dateStr]; exists {
            trend.Open = trend.Total - trend.Closed
            result = append(result, *trend)
        } else {
            // 填充空日期
            if len(result) > 0 {
                result = append(result, DefectTrendData{
                    Date:   dateStr,
                    Total:  result[len(result)-1].Total,
                    Closed: result[len(result)-1].Closed,
                    Open:   result[len(result)-1].Open,
                })
            }
        }
    }
    
    return result
}
```

### 2. 柱状图数据API

```go
// GetDefectStats 获取缺陷统计（柱状图）
func (c *Controller) GetDefectStats(ctx *gin.Context) {
    projectId := ctx.Param("projectId")
    
    // 按优先级统计
    stats, err := c.defectService.GetDefectsByPriority(projectId)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }
    
    // 转换为图表数据格式
    var chartData []BarChartData
    for priority, count := range stats {
        chartData = append(chartData, BarChartData{
            Name:  priority,
            Value: count,
        })
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": chartData,
    })
}
```

### 3. 饼图数据API

```go
// GetDefectDistribution 获取缺陷分布（饼图）
func (c *Controller) GetDefectDistribution(ctx *gin.Context) {
    projectId := ctx.Param("projectId")
    
    // 按状态统计
    distribution, err := c.defectService.GetDefectsByStatus(projectId)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}, http.StatusInternalServerError)
        return
    }
    
    // 转换为图表数据
    var chartData []PieChartData
    for status, count := range distribution {
        chartData = append(chartData, PieChartData{
            Name:  status,
            Value: count,
        })
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": chartData,
    })
}
```

---

## 📋 数据格式检查清单

Go后端返回的数据应满足以下要求：

### ✅ 必须条件

- [ ] 数据格式为JSON数组：`[{ }, { }]`
- [ ] 每个对象包含必要的字段
- [ ] 日期格式统一：`YYYY-MM-DD`
- [ ] 数值字段为数字类型（不是字符串）
- [ ] 不包含null值（使用默认值或0）
- [ ] 响应包含success和data字段

### ⚠️ 常见错误

```go
// ❌ 错误：字段名与前端dataKey不匹配
type DefectTrendData struct {
    DateField   string  `json:"dateField"`    // 前端期望：date
    DefectCount int     `json:"defectCount"`  // 前端期望：total
}

// ✅ 正确：字段名与前端一致
type DefectTrendData struct {
    Date   string `json:"date"`
    Total  int    `json:"total"`
    Closed int    `json:"closed"`
}
```

```go
// ❌ 错误：返回格式不是数组
ctx.JSON(http.StatusOK, DefectTrendData{})

// ✅ 正确：返回数组格式
ctx.JSON(http.StatusOK, gin.H{
    "success": true,
    "data": []DefectTrendData{},
})
```

```go
// ❌ 错误：数值作为字符串
type BarChartData struct {
    Value string `json:"value"` // "100"
}

// ✅ 正确：数值作为数字
type BarChartData struct {
    Value int `json:"value"` // 100
}
```

---

## 🔗 API端点规范

### 推荐的API端点设计

```
GET /api/v1/projects/{projectId}/defects/trend
    返回：[]DefectTrendData - 缺陷趋势数据

GET /api/v1/projects/{projectId}/defects/stats
    返回：[]BarChartData - 缺陷统计数据

GET /api/v1/projects/{projectId}/defects/distribution
    返回：[]PieChartData - 缺陷分布数据

GET /api/v1/projects/{projectId}/defects/priority
    返回：[]BarChartData - 按优先级统计

GET /api/v1/projects/{projectId}/defects/severity
    返回：[]BarChartData - 按严重级别统计
```

---

## 🧪 测试数据生成

### Go测试数据示例

```go
// 生成测试数据
func generateTestDefectTrendData() []DefectTrendData {
    now := time.Now()
    var data []DefectTrendData
    
    for i := 0; i < 30; i++ {
        date := now.AddDate(0, 0, -i)
        data = append(data, DefectTrendData{
            Date:   date.Format("2006-01-02"),
            Total:  10 + i,
            Closed: 3 + i/2,
        })
    }
    
    return data
}

// 在测试中使用
func TestGetDefectTrend(t *testing.T) {
    data := generateTestDefectTrendData()
    
    // 序列化为JSON验证格式
    jsonData, _ := json.Marshal(data)
    t.Logf("Generated JSON: %s", string(jsonData))
}
```

---

## 🔄 前后端调用流程

```
前端React组件
    ↓
[调用API] useEffect(() => { fetchDefects(projectId) })
    ↓
Go后端API
    ↓
数据库查询
    ↓
[数据处理] 按日期分组、统计、排序
    ↓
[JSON响应] {success: true, data: [...]}
    ↓
前端接收数据
    ↓
[数据过滤] 根据日期范围筛选
    ↓
[缓存数据] useMemo、useSta
    ↓
Recharts渲染
    ↓
用户看到图表
```

---

## 💾 数据库查询示例

### SQL查询

```sql
-- 缺陷趋势：按创建日期统计
SELECT 
    DATE(created_at) as date,
    COUNT(*) as total,
    SUM(CASE WHEN status = 'Closed' THEN 1 ELSE 0 END) as closed
FROM defects
WHERE project_id = ?
GROUP BY DATE(created_at)
ORDER BY date;

-- 按优先级统计
SELECT 
    priority as name,
    COUNT(*) as value
FROM defects
WHERE project_id = ?
GROUP BY priority;

-- 按状态分布
SELECT 
    status as name,
    COUNT(*) as value
FROM defects
WHERE project_id = ?
GROUP BY status;
```

### Go ORM示例（GORM）

```go
// 获取趋势数据
func (s *DefectService) GetTrendData(projectId string) ([]DefectTrendData, error) {
    var results []DefectTrendData
    
    err := s.db.
        Model(&Defect{}).
        Where("project_id = ?", projectId).
        Select(
            "DATE(created_at) as date",
            "COUNT(*) as total",
            "SUM(CASE WHEN status = 'Closed' THEN 1 ELSE 0 END) as closed",
        ).
        Group("DATE(created_at)").
        Order("date").
        Scan(&results).Error
    
    return results, err
}
```

---

## ✅ 验证清单

部署前请检查：

- [ ] API返回的字段名与前端dataKey完全匹配
- [ ] 日期格式为 `YYYY-MM-DD`
- [ ] 数值字段为数字类型（int/float）
- [ ] 返回结果为数组格式
- [ ] 包含success和data字段
- [ ] 处理了空数据的情况
- [ ] 添加了适当的错误处理
- [ ] 数据已按时间排序
- [ ] 包含必要的日期记录（无空白）
- [ ] 通过Postman或curl测试了API

---

## 🐛 常见问题排查

### Q1: 图表为什么不显示？
**A:** 检查以下几点：
1. API是否返回正确的JSON数组
2. 字段名是否与前端dataKey一致
3. 数据是否为空
4. 浏览器控制台是否有错误

### Q2: 数据点很少或不连续？
**A:** 
1. 确保查询包含所有日期（可能需要生成完整的日期序列）
2. 使用LEFT JOIN或FULL OUTER JOIN处理无数据的日期
3. 填充缺失日期的累积值

### Q3: 性能太慢？
**A:**
1. 添加数据库索引：`CREATE INDEX idx_project_date ON defects(project_id, created_at)`
2. 使用分页或时间范围限制
3. 考虑缓存查询结果

### Q4: 前端收到数据但不显示？
**A:**
1. 检查ResponsiveContainer的高度是否设置
2. 确认dataKey与JSON字段名完全匹配（大小写敏感）
3. 验证数据类型（数值不能是字符串）

---

## 📖 参考资源

- [Gin框架文档](https://gin-gonic.com/)
- [GORM文档](https://gorm.io/)
- [JSON Tag说明](https://golang.org/pkg/encoding/json/)
- [Recharts前端集成指南](./RECHARTS_INTEGRATION_GUIDE.md)

---

## 🚀 部署检查

### 上线前验证

```bash
# 1. 检查API响应格式
curl http://localhost:8080/api/v1/projects/123/defects/trend | jq

# 2. 验证JSON格式（使用在线JSON验证器）
# 确保：
# - 是有效的JSON数组
# - 所有字段都有值
# - 数值类型正确

# 3. 在浏览器开发工具中检查
# - Network标签查看API响应
# - Console标签检查错误
# - 确保请求状态码为200
```

---

## 📞 前后端沟通

### 集成前确认清单

与前端开发人员确认以下内容：

- [ ] 确认需要的图表类型
- [ ] 确认数据字段名和格式
- [ ] 确认API端点路径
- [ ] 确认数据刷新频率
- [ ] 确认是否需要分页
- [ ] 确认时间范围限制
- [ ] 确认错误处理方式

---

**文档版本**：1.0  
**更新时间**：2025-12-26  
**状态**：生产就绪

