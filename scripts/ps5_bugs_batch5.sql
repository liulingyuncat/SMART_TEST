-- Batch 5: Defects 1155 to 1164
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001155-' || hex(randomblob(8)), '001155', 35, '[Share/Capture] Trophy list issue when resuming from suspend', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，Trophy list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 resuming from suspend 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v2.0.0', '重启主机或重新进入Share/Capture功能', 'B', 'Minor', 'UIInteraction', 'Often', 'v2.0.0-beta', 'UT', 'TC-PS5-3489', 'dev_johnson', 'S', 'QA-California', 'Share/Capture > Trophy list', 'v1.0.0', '', 'Audio Engine', 'Fixed in firmware update', 'PS5 Digital Edition', 'Kim', 'Resolved', 1, 1, '2025-11-01 13:00:00', '2025-11-03 21:00:00'),
('001156-' || hex(randomblob(8)), '001156', 35, '[Controller] Settings panel issue when updating system', 'Controller', '## 问题描述
在使用 Controller 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Controller 菜单
3. 执行 updating system 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.1.0', '重启主机或重新进入Controller功能', 'D', 'Minor', 'UIInteraction', 'Rarely', 'v1.0.0', 'IT', 'TC-PS5-2728', 'dev_tanaka', 'S', 'QA-California', 'Controller > Settings panel', '', '', 'Storage Driver', '', 'PS5 Digital Edition', 'Li Ming', 'Reopened', 1, 1, '2025-11-10 07:00:00', '2025-11-11 15:00:00'),
('001157-' || hex(randomblob(8)), '001157', 35, '[Dashboard UI] Game menu issue when voice chat active', 'Dashboard UI', '## 问题描述
在使用 Dashboard UI 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Dashboard UI 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Dashboard UI功能', 'C', 'Critical', 'Compatibility', 'Rarely', 'v1.0.1', 'ST', 'TC-PS5-9073', 'dev_wang', 'B', 'QA-Tokyo', 'Dashboard UI > Game menu', 'v1.0.1', '', 'Audio Engine', 'Configuration issue resolved', 'PS5 Digital Edition', 'Kim', 'Resolved', 1, 1, '2026-02-06 06:00:00', '2026-02-08 00:00:00'),
('001158-' || hex(randomblob(8)), '001158', 35, '[Trophy System] Settings panel issue when voice chat active', 'Trophy System', '## 问题描述
在使用 Trophy System 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Trophy System 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.0.0', '重启主机或重新进入Trophy System功能', 'A', 'Trivial', 'Functional', 'Often', 'v1.1.0', 'IT', 'TC-PS5-6236', 'dev_kim', 'A', 'Dev-Test', 'Trophy System > Settings panel', 'v1.0.1', '', 'Backend', 'Code patch applied', 'PS5 Digital Edition', 'Emily Chen', 'Resolved', 1, 1, '2025-12-22 09:00:00', '2025-12-24 22:00:00'),
('001159-' || hex(randomblob(8)), '001159', 35, '[Share/Capture] Download queue issue when capturing screenshot', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.0.0', '重启主机或重新进入Share/Capture功能', 'B', 'Critical', 'Functional', 'Rarely', 'v1.0.1', 'Regression', 'TC-PS5-5556', 'dev_tanaka', 'A', 'Beta Tester', 'Share/Capture > Download queue', '', '', 'Graphics Engine', '', 'PS5 Slim', 'Li Ming', 'InProgress', 1, 1, '2025-11-28 07:00:00', '2025-11-30 16:00:00'),
('001160-' || hex(randomblob(8)), '001160', 35, '[Party Chat] Download queue issue when voice chat active', 'Party Chat', '## 问题描述
在使用 Party Chat 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Party Chat 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.1.0', '重启主机或重新进入Party Chat功能', 'B', 'Major', 'Security', 'Once', 'v1.2.0', 'UAT', 'TC-PS5-6533', 'dev_johnson', 'A', 'QA-Shanghai', 'Party Chat > Download queue', '', '', 'Storage Driver', '', 'PS5 Pro', 'Tanaka', 'New', 1, 1, '2026-01-25 23:00:00', '2026-01-28 06:00:00'),
('001161-' || hex(randomblob(8)), '001161', 35, '[Party Chat] Trophy list issue when updating system', 'Party Chat', '## 问题描述
在使用 Party Chat 功能时，Trophy list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Party Chat 菜单
3. 执行 updating system 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.0.1', '重启主机或重新进入Party Chat功能', 'B', 'Critical', 'Security', 'Sometimes', 'v1.2.0', 'Regression', 'TC-PS5-7663', 'dev_chen', 'A', 'QA-Shanghai', 'Party Chat > Trophy list', '', '', 'Backend', '', 'PS5 Digital Edition', 'Tanaka', 'Confirmed', 1, 1, '2025-12-09 14:00:00', '2025-12-12 06:00:00'),
('001162-' || hex(randomblob(8)), '001162', 35, '[User Account] Settings panel issue when downloading content', 'User Account', '## 问题描述
在使用 User Account 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.0.1', '重启主机或重新进入User Account功能', 'C', 'Major', 'UI', 'Once', 'v1.2.0', 'UT', 'TC-PS5-7220', 'dev_tanaka', 'C', 'Beta Tester', 'User Account > Settings panel', 'v1.1.0', '', 'Network Stack', 'Configuration issue resolved', 'PS5 Standard Edition', 'Tanaka', 'Closed', 1, 1, '2025-12-01 13:00:00', '2025-12-02 18:00:00'),
('001163-' || hex(randomblob(8)), '001163', 35, '[User Account] Home screen issue when connecting to server', 'User Account', '## 问题描述
在使用 User Account 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 connecting to server 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.0.0', '重启主机或重新进入User Account功能', 'D', 'Major', 'Compatibility', 'Sometimes', 'v1.0.1', 'Regression', 'TC-PS5-1290', 'dev_kim', 'A', 'Dev-Test', 'User Account > Home screen', 'v1.0.1', '', 'Graphics Engine', 'Code patch applied', 'PS5 Pro', 'Li Ming', 'Closed', 1, 1, '2025-12-13 17:00:00', '2025-12-16 07:00:00'),
('001164-' || hex(randomblob(8)), '001164', 35, '[Controller] Settings panel issue when starting application', 'Controller', '## 问题描述
在使用 Controller 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Controller 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入Controller功能', 'C', 'Minor', 'Compatibility', 'Once', 'v2.0.0-beta', 'UAT', 'TC-PS5-9539', 'dev_tanaka', 'C', 'QA-California', 'Controller > Settings panel', '', '', 'Backend', '', 'PS5 Standard Edition', 'Sato', 'Confirmed', 1, 1, '2025-11-11 09:00:00', '2025-11-11 20:00:00');