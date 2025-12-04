# 测试清空AI用例API
# 使用方法: .\test_clear_ai_cases.ps1 <project_id> <auth_token>

param(
    [Parameter(Mandatory=$true)]
    [int]$ProjectId,
    
    [Parameter(Mandatory=$true)]
    [string]$Token
)

$baseUrl = "http://localhost:8080/api/v1"
$headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type" = "application/json"
}

Write-Host "测试清空AI用例API..." -ForegroundColor Cyan
Write-Host "项目ID: $ProjectId" -ForegroundColor Yellow

try {
    # 调用清空AI用例API
    $response = Invoke-RestMethod -Uri "$baseUrl/projects/$ProjectId/manual-cases/clear-ai" `
        -Method Delete `
        -Headers $headers `
        -ErrorAction Stop
    
    Write-Host "`n✅ 成功!" -ForegroundColor Green
    Write-Host "响应数据:" -ForegroundColor Cyan
    $response | ConvertTo-Json -Depth 3
    
} catch {
    Write-Host "`n❌ 失败!" -ForegroundColor Red
    Write-Host "状态码: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
    Write-Host "错误信息: $($_.Exception.Message)" -ForegroundColor Red
    
    if ($_.ErrorDetails) {
        Write-Host "详细错误:" -ForegroundColor Red
        $_.ErrorDetails.Message
    }
}
