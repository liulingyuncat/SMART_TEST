-- Batch 6: Defects 1165 to 1174
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001165-' || hex(randomblob(8)), '001165', 35, '[Share/Capture] Game menu issue when starting application', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Share/Capture功能', 'D', 'Critical', 'UIInteraction', 'Once', 'v1.1.0', 'UT', 'TC-PS5-3072', 'dev_tanaka', 'S', 'QA-Shanghai', 'Share/Capture > Game menu', '', '', 'Graphics Engine', '', 'PS5 Standard Edition', 'Kim', 'Reopened', 1, 1, '2025-12-19 10:00:00', '2025-12-21 01:00:00'),
('001166-' || hex(randomblob(8)), '001166', 35, '[Trophy System] Party interface issue when connecting to server', 'Trophy System', '## 问题描述
在使用 Trophy System 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Trophy System 菜单
3. 执行 connecting to server 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v2.0.0-beta', '重启主机或重新进入Trophy System功能', 'B', 'Critical', 'Environment', 'Rarely', 'v1.1.0', 'IT', 'TC-PS5-6590', 'dev_kim', 'C', 'QA-Tokyo', 'Trophy System > Party interface', '', '', 'Graphics Engine', '', 'PS5 Digital Edition', 'Tanaka', 'Reopened', 1, 1, '2026-01-22 10:00:00', '2026-01-23 23:00:00'),
('001167-' || hex(randomblob(8)), '001167', 35, '[Remote Play] Home screen issue when entering rest mode', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 entering rest mode 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.1.0', '重启主机或重新进入Remote Play功能', 'D', 'Trivial', 'Functional', 'Always', 'v2.0.0-beta', 'ST', 'TC-PS5-4651', 'dev_tanaka', 'S', 'QA-California', 'Remote Play > Home screen', '', '', 'Network Stack', '', 'PS5 Standard Edition', 'Tanaka', 'Reopened', 1, 1, '2025-11-01 14:00:00', '2025-11-04 06:00:00'),
('001168-' || hex(randomblob(8)), '001168', 35, '[Controller] Download queue issue when loading game', 'Controller', '## 问题描述
在使用 Controller 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Controller 菜单
3. 执行 loading game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.1.0', '重启主机或重新进入Controller功能', 'A', 'Trivial', 'Compatibility', 'Often', 'v1.2.0', 'Regression', 'TC-PS5-8574', 'dev_kim', 'A', 'QA-California', 'Controller > Download queue', '', '', 'Network Stack', '', 'PS5 Standard Edition', 'Tanaka', 'InProgress', 1, 1, '2025-12-29 04:00:00', '2025-12-30 23:00:00'),
('001169-' || hex(randomblob(8)), '001169', 35, '[User Account] User profile issue when updating system', 'User Account', '## 问题描述
在使用 User Account 功能时，User profile 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 updating system 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0-beta', '重启主机或重新进入User Account功能', 'D', 'Trivial', 'Performance', 'Often', 'v1.0.1', 'UAT', 'TC-PS5-9936', 'dev_chen', 'C', 'Dev-Test', 'User Account > User profile', 'v1.0.0', '', 'Audio Engine', 'Configuration issue resolved', 'PS5 Slim', 'Kim', 'Resolved', 1, 1, '2026-01-18 16:00:00', '2026-01-19 04:00:00'),
('001170-' || hex(randomblob(8)), '001170', 35, '[Dashboard UI] Party interface issue when resuming from suspend', 'Dashboard UI', '## 问题描述
在使用 Dashboard UI 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Dashboard UI 菜单
3. 执行 resuming from suspend 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.1.0', '重启主机或重新进入Dashboard UI功能', 'D', 'Critical', 'Functional', 'Always', 'v1.2.0', 'ST', 'TC-PS5-1235', 'dev_wang', 'A', 'Beta Tester', 'Dashboard UI > Party interface', '', '', 'Storage Driver', '', 'PS5 Standard Edition', 'Emily Chen', 'Reopened', 1, 1, '2026-02-02 13:00:00', '2026-02-05 07:00:00'),
('001171-' || hex(randomblob(8)), '001171', 35, '[Dashboard UI] Home screen issue when switching user', 'Dashboard UI', '## 问题描述
在使用 Dashboard UI 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Dashboard UI 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.2.0', '重启主机或重新进入Dashboard UI功能', 'D', 'Critical', 'UIInteraction', 'Always', 'v1.0.0', 'ST', 'TC-PS5-7785', 'dev_chen', 'C', 'QA-California', 'Dashboard UI > Home screen', '', '', 'Backend', '', 'PS5 Standard Edition', 'Zhang Wei', 'Confirmed', 1, 1, '2025-11-02 12:00:00', '2025-11-05 04:00:00'),
('001172-' || hex(randomblob(8)), '001172', 35, '[Controller] Game menu issue when streaming gameplay', 'Controller', '## 问题描述
在使用 Controller 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Controller 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.0.1', '重启主机或重新进入Controller功能', 'B', 'Minor', 'Environment', 'Always', 'v1.1.0', 'Regression', 'TC-PS5-1597', 'dev_johnson', 'B', 'QA-Tokyo', 'Controller > Game menu', 'v1.0.1', '', 'Backend', 'Fixed in firmware update', 'PS5 Slim', 'Emily Chen', 'Closed', 1, 1, '2025-11-19 06:00:00', '2025-11-21 07:00:00'),
('001173-' || hex(randomblob(8)), '001173', 35, '[Remote Play] Game menu issue when streaming gameplay', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.2.0', '重启主机或重新进入Remote Play功能', 'B', 'Critical', 'UIInteraction', 'Always', 'v2.0.0', 'IT', 'TC-PS5-1356', 'dev_suzuki', 'B', 'Dev-Test', 'Remote Play > Game menu', '', '', 'Storage Driver', '', 'PS5 Digital Edition', 'Zhang Wei', 'New', 1, 1, '2026-01-15 19:00:00', '2026-01-17 23:00:00'),
('001174-' || hex(randomblob(8)), '001174', 35, '[Storage Management] Home screen issue when capturing screenshot', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入Storage Management功能', 'D', 'Trivial', 'Performance', 'Once', 'v1.0.0', 'UT', 'TC-PS5-9188', 'dev_chen', 'B', 'QA-Tokyo', 'Storage Management > Home screen', '', '', 'Graphics Engine', '', 'PS5 Digital Edition', 'Sato', 'InProgress', 1, 1, '2026-02-02 07:00:00', '2026-02-05 05:00:00');