-- Batch 8: Defects 1185 to 1194
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001185-' || hex(randomblob(8)), '001185', 35, '[Remote Play] Settings panel issue when switching user', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入Remote Play功能', 'D', 'Critical', 'Environment', 'Sometimes', 'v1.1.0', 'Regression', 'TC-PS5-6288', 'dev_suzuki', 'S', 'QA-Shanghai', 'Remote Play > Settings panel', 'v1.2.0', '', 'Network Stack', 'Configuration issue resolved', 'PS5 Slim', 'Tanaka', 'Resolved', 1, 1, '2025-11-11 18:00:00', '2025-11-11 23:00:00'),
('001186-' || hex(randomblob(8)), '001186', 35, '[Audio System] Home screen issue when starting application', 'Audio System', '## 问题描述
在使用 Audio System 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Audio System 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.1.0', '重启主机或重新进入Audio System功能', 'B', 'Minor', 'UIInteraction', 'Once', 'v2.0.0-beta', 'UAT', 'TC-PS5-6103', 'dev_wang', 'S', 'QA-California', 'Audio System > Home screen', 'v1.0.0', '', 'Firmware', 'Code patch applied', 'PS5 Slim', 'Emily Chen', 'Closed', 1, 1, '2026-01-29 05:00:00', '2026-01-29 11:00:00'),
('001187-' || hex(randomblob(8)), '001187', 35, '[Settings] User profile issue when switching user', 'Settings', '## 问题描述
在使用 Settings 功能时，User profile 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Settings 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v2.0.0', '重启主机或重新进入Settings功能', 'C', 'Minor', 'Performance', 'Once', 'v2.0.0', 'IT', 'TC-PS5-4362', 'dev_chen', 'B', 'QA-Shanghai', 'Settings > User profile', '', '', 'Frontend', '', 'PS5 Slim', 'Sato', 'InProgress', 1, 1, '2025-11-09 23:00:00', '2025-11-11 00:00:00'),
('001188-' || hex(randomblob(8)), '001188', 35, '[Game Library] Game menu issue when installing game', 'Game Library', '## 问题描述
在使用 Game Library 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 installing game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入Game Library功能', 'A', 'Minor', 'Functional', 'Often', 'v2.0.0', 'ST', 'TC-PS5-1716', 'dev_tanaka', 'C', 'Dev-Test', 'Game Library > Game menu', 'v1.2.0', '', 'Firmware', 'Fixed in firmware update', 'PS5 Slim', 'Kim', 'Resolved', 1, 1, '2025-12-13 06:00:00', '2025-12-16 06:00:00'),
('001189-' || hex(randomblob(8)), '001189', 35, '[Storage Management] Party interface issue when connecting to server', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 connecting to server 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0', '重启主机或重新进入Storage Management功能', 'B', 'Trivial', 'Compatibility', 'Often', 'v1.2.0', 'IT', 'TC-PS5-9086', 'dev_tanaka', 'S', 'QA-Tokyo', 'Storage Management > Party interface', '', '', 'Frontend', '', 'PS5 Pro', 'Emily Chen', 'Reopened', 1, 1, '2025-11-25 14:00:00', '2025-11-27 08:00:00'),
('001190-' || hex(randomblob(8)), '001190', 35, '[Video Output] User profile issue when capturing screenshot', 'Video Output', '## 问题描述
在使用 Video Output 功能时，User profile 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Video Output 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.2.0', '重启主机或重新进入Video Output功能', 'C', 'Major', 'UI', 'Once', 'v1.2.0', 'IT', 'TC-PS5-2863', 'dev_wang', 'S', 'QA-Tokyo', 'Video Output > User profile', 'v1.2.0', '', 'Network Stack', 'Configuration issue resolved', 'PS5 Standard Edition', 'Li Ming', 'Closed', 1, 1, '2025-11-10 00:00:00', '2025-11-12 13:00:00'),
('001191-' || hex(randomblob(8)), '001191', 35, '[Game Library] Download queue issue when connecting to server', 'Game Library', '## 问题描述
在使用 Game Library 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 connecting to server 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Game Library功能', 'A', 'Major', 'UI', 'Always', 'v1.0.1', 'IT', 'TC-PS5-5594', 'dev_johnson', 'B', 'Dev-Test', 'Game Library > Download queue', '', '', 'Backend', '', 'PS5 Digital Edition', 'Emily Chen', 'InProgress', 1, 1, '2026-01-25 02:00:00', '2026-01-25 17:00:00'),
('001192-' || hex(randomblob(8)), '001192', 35, '[Party Chat] Friends list issue when downloading content', 'Party Chat', '## 问题描述
在使用 Party Chat 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Party Chat 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v2.0.0-beta', '重启主机或重新进入Party Chat功能', 'C', 'Major', 'UI', 'Sometimes', 'v1.0.1', 'UAT', 'TC-PS5-5355', 'dev_tanaka', 'A', 'Beta Tester', 'Party Chat > Friends list', '', '', 'Firmware', '', 'PS5 Slim', 'John Smith', 'Confirmed', 1, 1, '2025-12-30 08:00:00', '2025-12-30 22:00:00'),
('001193-' || hex(randomblob(8)), '001193', 35, '[Settings] Friends list issue when loading game', 'Settings', '## 问题描述
在使用 Settings 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Settings 菜单
3. 执行 loading game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.2.0', '重启主机或重新进入Settings功能', 'D', 'Critical', 'Performance', 'Sometimes', 'v1.0.1', 'ST', 'TC-PS5-3199', 'dev_suzuki', 'C', 'QA-Shanghai', 'Settings > Friends list', '', '', 'Network Stack', '', 'PS5 Slim', 'Emily Chen', 'Confirmed', 1, 1, '2025-11-15 02:00:00', '2025-11-17 05:00:00'),
('001194-' || hex(randomblob(8)), '001194', 35, '[Share/Capture] Settings panel issue when connecting to server', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 connecting to server 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.0.0', '重启主机或重新进入Share/Capture功能', 'C', 'Trivial', 'Security', 'Rarely', 'v2.0.0', 'UT', 'TC-PS5-4495', 'dev_suzuki', 'C', 'QA-California', 'Share/Capture > Settings panel', 'v2.0.0', '', 'Firmware', 'Code patch applied', 'PS5 Slim', 'Kim', 'Resolved', 1, 1, '2025-11-13 00:00:00', '2025-11-13 05:00:00');