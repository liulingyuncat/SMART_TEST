# 测试自动化用例版本保存API
# 使用前先登录获取 token

$baseUrl = "http://localhost:8080/api/v1"

# 1. 登录获取token
Write-Host "=== 步骤1: 登录 ===" -ForegroundColor Cyan
$loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body (@{
    username = "admin"
    password = "admin123"
} | ConvertTo-Json) -ContentType "application/json"

$token = $loginResponse.data.token
Write-Host "登录成功，token: $($token.Substring(0,20))..." -ForegroundColor Green

# 2. 获取项目列表
Write-Host "`n=== 步骤2: 获取项目列表 ===" -ForegroundColor Cyan
$projectsResponse = Invoke-RestMethod -Uri "$baseUrl/projects" -Method Get -Headers @{
    "Authorization" = "Bearer $token"
}

if ($projectsResponse.data.projects.Count -eq 0) {
    Write-Host "没有项目，测试终止" -ForegroundColor Red
    exit 1
}

$projectId = $projectsResponse.data.projects[0].id
$projectName = $projectsResponse.data.projects[0].name
Write-Host "使用项目: ID=$projectId, Name=$projectName" -ForegroundColor Green

# 3. 检查项目中是否有自动化用例
Write-Host "`n=== 步骤3: 检查ROLE1用例 ===" -ForegroundColor Cyan
try {
    $casesResponse = Invoke-RestMethod -Uri "$baseUrl/projects/$projectId/auto-cases?case_type=role1&page=1&size=10" -Method Get -Headers @{
        "Authorization" = "Bearer $token"
    }
    Write-Host "ROLE1用例数量: $($casesResponse.data.total)" -ForegroundColor Green
} catch {
    Write-Host "获取用例失败: $_" -ForegroundColor Red
}

# 4. 保存版本
Write-Host "`n=== 步骤4: 批量保存版本 ===" -ForegroundColor Cyan
try {
    $saveResponse = Invoke-RestMethod -Uri "$baseUrl/projects/$projectId/auto-cases/versions" -Method Post -Headers @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    }
    
    Write-Host "版本保存成功！" -ForegroundColor Green
    Write-Host "Version ID: $($saveResponse.data.version_id)" -ForegroundColor Yellow
    Write-Host "保存时间: $($saveResponse.data.saved_at)" -ForegroundColor Yellow
    Write-Host "总用例数: $($saveResponse.data.total_cases)" -ForegroundColor Yellow
    Write-Host "文件列表:" -ForegroundColor Yellow
    foreach ($file in $saveResponse.data.files) {
        Write-Host "  - $($file.role): $($file.filename) ($($file.count) cases)" -ForegroundColor Gray
    }
    
    $versionId = $saveResponse.data.version_id
} catch {
    Write-Host "版本保存失败: $_" -ForegroundColor Red
    Write-Host "错误详情: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "服务器响应: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
    exit 1
}

# 5. 获取版本列表
Write-Host "`n=== 步骤5: 获取版本列表 ===" -ForegroundColor Cyan
try {
    $versionsResponse = Invoke-RestMethod -Uri "$baseUrl/projects/$projectId/auto-cases/versions?page=1&size=10" -Method Get -Headers @{
        "Authorization" = "Bearer $token"
    }
    
    Write-Host "版本列表总数: $($versionsResponse.data.total)" -ForegroundColor Green
    if ($versionsResponse.data.versions.Count -gt 0) {
        Write-Host "`n最新版本:" -ForegroundColor Yellow
        $latestVersion = $versionsResponse.data.versions[0]
        Write-Host "  Version ID: $($latestVersion.version_id)" -ForegroundColor Gray
        Write-Host "  保存时间: $($latestVersion.saved_at)" -ForegroundColor Gray
        Write-Host "  总用例数: $($latestVersion.total_cases)" -ForegroundColor Gray
        Write-Host "  文件数量: $($latestVersion.files.Count)" -ForegroundColor Gray
    } else {
        Write-Host "版本列表为空！" -ForegroundColor Red
    }
} catch {
    Write-Host "获取版本列表失败: $_" -ForegroundColor Red
    Write-Host "错误详情: $($_.Exception.Message)" -ForegroundColor Red
}

# 6. 检查文件是否生成
Write-Host "`n=== 步骤6: 检查存储文件 ===" -ForegroundColor Cyan
$storageDir = "D:\VSCode\webtest\backend\storage\versions\auto-cases\$projectId\$versionId"
if (Test-Path $storageDir) {
    Write-Host "存储目录存在: $storageDir" -ForegroundColor Green
    $files = Get-ChildItem $storageDir -Filter "*.xlsx"
    Write-Host "Excel文件数量: $($files.Count)" -ForegroundColor Yellow
    foreach ($file in $files) {
        Write-Host "  - $($file.Name) ($([math]::Round($file.Length/1KB, 2)) KB)" -ForegroundColor Gray
    }
} else {
    Write-Host "存储目录不存在: $storageDir" -ForegroundColor Red
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Cyan
