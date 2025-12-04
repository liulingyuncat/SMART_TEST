# AI用例重排序功能测试脚本

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "AI用例重排序功能测试" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "问题描述:" -ForegroundColor Yellow
Write-Host "- 当前是10条/页" -ForegroundColor White
Write-Host "- 进入第二页，插入行，编号为No.88" -ForegroundColor White
Write-Host "- 第二页显示：No.88, No.10, No.11" -ForegroundColor White
Write-Host "- 点击重新排序后，No.88变成了No.4（错误）" -ForegroundColor Red
Write-Host "- 期望：No.88应该变成No.11（正确）" -ForegroundColor Green
Write-Host ""

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "修复方案" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "1. 问题根因:" -ForegroundColor Yellow
Write-Host "   前端ReorderModal使用getCasesList获取所有用例，然后按ID排序" -ForegroundColor White
Write-Host "   导致插入在第二页的No.88被排序到第88位" -ForegroundColor White
Write-Host ""

Write-Host "2. 解决方案:" -ForegroundColor Yellow
Write-Host "   ✓ 后端新增 /api/v1/projects/:id/manual-cases/reorder-all API" -ForegroundColor Green
Write-Host "   ✓ 该API按现有ID顺序重新编号所有用例（1,2,3...）" -ForegroundColor Green
Write-Host "   ✓ 前端调用新API，不再自己获取和排序用例" -ForegroundColor Green
Write-Host ""

Write-Host "3. 修改文件:" -ForegroundColor Yellow
Write-Host "   后端:" -ForegroundColor White
Write-Host "   - backend/cmd/server/main.go (添加路由)" -ForegroundColor Gray
Write-Host "   - backend/internal/services/manual_test_case_service.go (修复方法调用)" -ForegroundColor Gray
Write-Host ""
Write-Host "   前端:" -ForegroundColor White
Write-Host "   - frontend/src/api/manualCase.js (添加reorderAllCasesByID函数)" -ForegroundColor Gray
Write-Host "   - frontend/src/pages/ProjectDetail/ManualTestTabs/components/ReorderModal.jsx (使用新API)" -ForegroundColor Gray
Write-Host ""

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "测试步骤" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "手动测试步骤:" -ForegroundColor Yellow
Write-Host "1. 启动后端服务" -ForegroundColor White
Write-Host "   cd backend" -ForegroundColor Gray
Write-Host "   go run cmd/server/main.go" -ForegroundColor Gray
Write-Host ""
Write-Host "2. 启动前端服务" -ForegroundColor White
Write-Host "   cd frontend" -ForegroundColor Gray
Write-Host "   npm start" -ForegroundColor Gray
Write-Host ""
Write-Host "3. 登录系统，进入项目的AI用例页" -ForegroundColor White
Write-Host ""
Write-Host "4. 设置每页显示10条" -ForegroundColor White
Write-Host ""
Write-Host "5. 翻到第二页（显示No.11-20）" -ForegroundColor White
Write-Host ""
Write-Host "6. 在第二页点击插入行按钮" -ForegroundColor White
Write-Host "   - 新行的No应该是一个较大的数字（如No.88）" -ForegroundColor Gray
Write-Host ""
Write-Host "7. 点击重新排序按钮" -ForegroundColor White
Write-Host ""
Write-Host "8. 验证结果:" -ForegroundColor White
Write-Host "   预期: No.88应该变成No.11（第二页第一条）" -ForegroundColor Green
Write-Host "   错误: No.88变成No.4（修复前的错误）" -ForegroundColor Red
Write-Host ""

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "API测试（可选）" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "使用curl测试新API:" -ForegroundColor Yellow
Write-Host "curl -X POST http://localhost:8080/api/v1/projects/1/manual-cases/reorder-all" -ForegroundColor Gray
Write-Host "  -H Content-Type: application/json" -ForegroundColor Gray
Write-Host "  -H Authorization: Bearer YOUR_TOKEN" -ForegroundColor Gray
Write-Host "  -d JSON_DATA" -ForegroundColor Gray
Write-Host ""

Write-Host "期望返回JSON:" -ForegroundColor Yellow
Write-Host "包含message和count字段" -ForegroundColor Gray
Write-Host ""

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "完成" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
