-- Batch 7: Defects 1175 to 1184
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001175-' || hex(randomblob(8)), '001175', 35, '[Trophy System] Party interface issue when downloading content', 'Trophy System', '## 问题描述
在使用 Trophy System 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Trophy System 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.0.1', '重启主机或重新进入Trophy System功能', 'D', 'Critical', 'UIInteraction', 'Always', 'v2.0.0-beta', 'UT', 'TC-PS5-7516', 'dev_chen', 'B', 'Beta Tester', 'Trophy System > Party interface', '', '', 'Storage Driver', '', 'PS5 Standard Edition', 'Tanaka', 'New', 1, 1, '2025-11-29 07:00:00', '2025-12-01 21:00:00'),
('001176-' || hex(randomblob(8)), '001176', 35, '[Video Output] Settings panel issue when downloading content', 'Video Output', '## 问题描述
在使用 Video Output 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Video Output 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入Video Output功能', 'B', 'Minor', 'Security', 'Once', 'v1.2.0', 'Regression', 'TC-PS5-2779', 'dev_wang', 'C', 'QA-Tokyo', 'Video Output > Settings panel', 'v2.0.0', '', 'Network Stack', 'Configuration issue resolved', 'PS5 Slim', 'John Smith', 'Closed', 1, 1, '2026-01-28 23:00:00', '2026-01-31 13:00:00'),
('001177-' || hex(randomblob(8)), '001177', 35, '[PlayStation Store] Store page issue when voice chat active', 'PlayStation Store', '## 问题描述
在使用 PlayStation Store 功能时，Store page 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 PlayStation Store 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.2.0', '重启主机或重新进入PlayStation Store功能', 'B', 'Minor', 'UIInteraction', 'Once', 'v1.0.0', 'UT', 'TC-PS5-4171', 'dev_kim', 'C', 'Beta Tester', 'PlayStation Store > Store page', '', '', 'Backend', '', 'PS5 Pro', 'Yamamoto', 'Confirmed', 1, 1, '2025-11-22 19:00:00', '2025-11-25 18:00:00'),
('001178-' || hex(randomblob(8)), '001178', 35, '[Dashboard UI] User profile issue when resuming from suspend', 'Dashboard UI', '## 问题描述
在使用 Dashboard UI 功能时，User profile 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Dashboard UI 菜单
3. 执行 resuming from suspend 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0-beta', '重启主机或重新进入Dashboard UI功能', 'A', 'Critical', 'Security', 'Often', 'v1.0.0', 'Regression', 'TC-PS5-2206', 'dev_wang', 'C', 'Beta Tester', 'Dashboard UI > User profile', 'v2.0.0', '', 'Graphics Engine', 'Configuration issue resolved', 'PS5 Pro', 'Sato', 'Resolved', 1, 1, '2025-11-06 22:00:00', '2025-11-07 02:00:00'),
('001179-' || hex(randomblob(8)), '001179', 35, '[Game Library] Party interface issue when switching user', 'Game Library', '## 问题描述
在使用 Game Library 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.1.0', '重启主机或重新进入Game Library功能', 'A', 'Critical', 'Performance', 'Always', 'v1.2.0', 'UAT', 'TC-PS5-4207', 'dev_suzuki', 'B', 'QA-California', 'Game Library > Party interface', 'v1.1.0', '', 'Audio Engine', 'Configuration issue resolved', 'PS5 Slim', 'Emily Chen', 'Closed', 1, 1, '2026-01-06 17:00:00', '2026-01-08 13:00:00'),
('001180-' || hex(randomblob(8)), '001180', 35, '[Game Library] Friends list issue when capturing screenshot', 'Game Library', '## 问题描述
在使用 Game Library 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Game Library功能', 'C', 'Trivial', 'UIInteraction', 'Rarely', 'v1.0.0', 'ST', 'TC-PS5-6640', 'dev_tanaka', 'S', 'Beta Tester', 'Game Library > Friends list', '', '', 'Frontend', '', 'PS5 Digital Edition', 'John Smith', 'New', 1, 1, '2025-11-25 08:00:00', '2025-11-27 21:00:00'),
('001181-' || hex(randomblob(8)), '001181', 35, '[PlayStation Store] Trophy list issue when installing game', 'PlayStation Store', '## 问题描述
在使用 PlayStation Store 功能时，Trophy list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 PlayStation Store 菜单
3. 执行 installing game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.0.1', '重启主机或重新进入PlayStation Store功能', 'C', 'Critical', 'Environment', 'Rarely', 'v2.0.0', 'ST', 'TC-PS5-9335', 'dev_kim', 'B', 'Beta Tester', 'PlayStation Store > Trophy list', 'v1.0.1', '', 'Audio Engine', 'Code patch applied', 'PS5 Slim', 'John Smith', 'Closed', 1, 1, '2026-01-05 13:00:00', '2026-01-07 08:00:00'),
('001182-' || hex(randomblob(8)), '001182', 35, '[Share/Capture] Store page issue when streaming gameplay', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，Store page 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.0.0', '重启主机或重新进入Share/Capture功能', 'A', 'Major', 'UIInteraction', 'Once', 'v2.0.0', 'UT', 'TC-PS5-7535', 'dev_tanaka', 'A', 'QA-Shanghai', 'Share/Capture > Store page', 'v1.0.0', '', 'Frontend', 'Fixed in firmware update', 'PS5 Slim', 'Zhang Wei', 'Closed', 1, 1, '2025-11-23 19:00:00', '2025-11-25 17:00:00'),
('001183-' || hex(randomblob(8)), '001183', 35, '[Storage Management] Download queue issue when starting application', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Storage Management功能', 'A', 'Minor', 'UIInteraction', 'Once', 'v1.1.0', 'ST', 'TC-PS5-4585', 'dev_tanaka', 'S', 'QA-Tokyo', 'Storage Management > Download queue', '', '', 'Graphics Engine', '', 'PS5 Standard Edition', 'Sato', 'New', 1, 1, '2025-12-18 20:00:00', '2025-12-20 23:00:00'),
('001184-' || hex(randomblob(8)), '001184', 35, '[Audio System] Home screen issue when installing game', 'Audio System', '## 问题描述
在使用 Audio System 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Audio System 菜单
3. 执行 installing game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.1.0', '重启主机或重新进入Audio System功能', 'D', 'Minor', 'Performance', 'Sometimes', 'v1.0.0', 'ST', 'TC-PS5-7214', 'dev_tanaka', 'S', 'Beta Tester', 'Audio System > Home screen', 'v1.0.1', '', 'Firmware', 'Configuration issue resolved', 'PS5 Pro', 'Emily Chen', 'Closed', 1, 1, '2026-01-15 00:00:00', '2026-01-15 09:00:00');