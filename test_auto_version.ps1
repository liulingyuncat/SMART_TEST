# 自动化测试用例版本管理集成测试脚本
# 测试流程：保存版本 -> 查询列表 -> 下载压缩包 -> 更新备注 -> 删除版本

$ErrorActionPreference = "Stop"
$baseUrl = "http://localhost:8080/api/v1"
$projectId = 1  # 替换为实际项目ID
$token = ""     # 替换为实际token

Write-Host "========== 自动化测试用例版本管理集成测试 ==========" -ForegroundColor Cyan

# 1. 批量保存版本
Write-Host "`n[测试1] 批量保存版本..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/projects/$projectId/auto-cases/versions" `
        -Method POST `
        -Headers @{ "Authorization" = "Bearer $token" } `
        -ContentType "application/json"
    
    Write-Host "✓ 版本保存成功" -ForegroundColor Green
    Write-Host "  版本ID: $($response.data.version_id)" -ForegroundColor Gray
    Write-Host "  项目名称: $($response.data.project_name)" -ForegroundColor Gray
    Write-Host "  文件数量: $($response.data.files.Count)" -ForegroundColor Gray
    
    $versionId = $response.data.version_id
    
    foreach ($file in $response.data.files) {
        Write-Host "  - $($file.filename): $($file.case_count) 条用例, $([math]::Round($file.file_size/1024, 2)) KB" -ForegroundColor Gray
    }
} catch {
    Write-Host "✗ 版本保存失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Start-Sleep -Seconds 1

# 2. 获取版本列表
Write-Host "`n[测试2] 获取版本列表..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/projects/$projectId/auto-cases/versions?page=1&size=10" `
        -Method GET `
        -Headers @{ "Authorization" = "Bearer $token" }
    
    Write-Host "✓ 版本列表获取成功" -ForegroundColor Green
    Write-Host "  总数: $($response.data.total), 当前页: $($response.data.page)/$([math]::Ceiling($response.data.total/$response.data.size))" -ForegroundColor Gray
    
    foreach ($version in $response.data.versions) {
        Write-Host "  版本: $($version.version_id)" -ForegroundColor Gray
        Write-Host "    备注: $($version.remark)" -ForegroundColor Gray
        Write-Host "    创建时间: $($version.created_at)" -ForegroundColor Gray
    }
} catch {
    Write-Host "✗ 版本列表获取失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Start-Sleep -Seconds 1

# 3. 下载版本压缩包
Write-Host "`n[测试3] 下载版本压缩包..." -ForegroundColor Yellow
try {
    $outputFile = "test_version_$versionId.zip"
    Invoke-WebRequest -Uri "$baseUrl/projects/$projectId/auto-cases/versions/$versionId/export" `
        -Method GET `
        -Headers @{ "Authorization" = "Bearer $token" } `
        -OutFile $outputFile
    
    $fileSize = (Get-Item $outputFile).Length
    Write-Host "✓ 压缩包下载成功" -ForegroundColor Green
    Write-Host "  文件: $outputFile ($([math]::Round($fileSize/1024, 2)) KB)" -ForegroundColor Gray
    
    # 清理测试文件
    Remove-Item $outputFile -Force
    Write-Host "  已清理测试文件" -ForegroundColor Gray
} catch {
    Write-Host "✗ 压缩包下载失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Start-Sleep -Seconds 1

# 4. 更新版本备注
Write-Host "`n[测试4] 更新版本备注..." -ForegroundColor Yellow
try {
    $newRemark = "集成测试备注 - $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
    $body = @{ remark = $newRemark } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "$baseUrl/projects/$projectId/auto-cases/versions/$versionId/remark" `
        -Method PUT `
        -Headers @{ "Authorization" = "Bearer $token"; "Content-Type" = "application/json" } `
        -Body $body
    
    Write-Host "✓ 备注更新成功" -ForegroundColor Green
    Write-Host "  新备注: $newRemark" -ForegroundColor Gray
} catch {
    Write-Host "✗ 备注更新失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Start-Sleep -Seconds 1

# 5. 删除版本（可选，取消注释以测试）
# Write-Host "`n[测试5] 删除版本..." -ForegroundColor Yellow
# try {
#     $response = Invoke-RestMethod -Uri "$baseUrl/projects/$projectId/auto-cases/versions/$versionId" `
#         -Method DELETE `
#         -Headers @{ "Authorization" = "Bearer $token" }
#     
#     Write-Host "✓ 版本删除成功" -ForegroundColor Green
# } catch {
#     Write-Host "✗ 版本删除失败: $($_.Exception.Message)" -ForegroundColor Red
#     exit 1
# }

Write-Host "`n========== 所有测试通过 ==========" -ForegroundColor Green
Write-Host "提示：删除版本测试已注释，如需测试请取消注释第5部分" -ForegroundColor Yellow
