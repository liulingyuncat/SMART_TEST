# AI用例重排序修复 - 快速测试脚本

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  AI用例重排序功能修复 - 测试指南" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "修复状态:" -ForegroundColor Green
Write-Host "  [✓] 后端编译成功" -ForegroundColor Green
Write-Host "  [✓] 添加API路由: POST /api/v1/projects/:id/manual-cases/reorder-all" -ForegroundColor Green
Write-Host "  [✓] 修复Service接口定义" -ForegroundColor Green
Write-Host "  [✓] 前端API接口已添加" -ForegroundColor Green
Write-Host "  [✓] ReorderModal组件已更新" -ForegroundColor Green
Write-Host ""

Write-Host "----------------------------------------" -ForegroundColor Yellow
Write-Host "测试步骤:" -ForegroundColor Yellow
Write-Host "----------------------------------------" -ForegroundColor White
Write-Host ""
Write-Host "1. 启动后端服务" -ForegroundColor Cyan
Write-Host "   cd backend" -ForegroundColor Gray
Write-Host "   .\server.exe" -ForegroundColor Gray
Write-Host "   (或使用: go run cmd/server/main.go)" -ForegroundColor DarkGray
Write-Host ""

Write-Host "2. 启动前端服务（新窗口）" -ForegroundColor Cyan
Write-Host "   cd frontend" -ForegroundColor Gray
Write-Host "   npm start" -ForegroundColor Gray
Write-Host ""

Write-Host "3. 测试重排序功能" -ForegroundColor Cyan
Write-Host "   a) 登录系统" -ForegroundColor White
Write-Host "   b) 进入项目的AI用例页面" -ForegroundColor White
Write-Host "   c) 设置每页显示10条" -ForegroundColor White
Write-Host "   d) 翻到第二页（显示No.11-20）" -ForegroundColor White
Write-Host "   e) 点击'插入行'按钮" -ForegroundColor White
Write-Host "      -> 新行编号应该是较大的数字（如No.88）" -ForegroundColor DarkGray
Write-Host "   f) 点击'重新排序'按钮" -ForegroundColor White
Write-Host ""

Write-Host "4. 验证结果" -ForegroundColor Cyan
Write-Host "   期望: No.88 -> No.11" -ForegroundColor Green
Write-Host "   错误: No.88 -> No.4 (修复前的bug)" -ForegroundColor Red
Write-Host ""

Write-Host "----------------------------------------" -ForegroundColor Yellow
Write-Host "问题说明:" -ForegroundColor Yellow
Write-Host "----------------------------------------" -ForegroundColor White
Write-Host "修复前: 前端获取所有数据并按ID排序，导致新插入的" -ForegroundColor White
Write-Host "        No.88被排到第88位" -ForegroundColor White
Write-Host ""
Write-Host "修复后: 后端直接按现有ID顺序重新编号，保持插入" -ForegroundColor White
Write-Host "        位置的正确性" -ForegroundColor White
Write-Host ""

Write-Host "详细文档: docs\ai-case-reorder-fix.md" -ForegroundColor Cyan
Write-Host ""
