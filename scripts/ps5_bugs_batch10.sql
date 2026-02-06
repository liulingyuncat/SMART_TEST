-- Batch 10: Defects 1205 to 1214
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001205-' || hex(randomblob(8)), '001205', 35, '[Game Library] User profile issue when switching user', 'Game Library', '## 问题描述
在使用 Game Library 功能时，User profile 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v2.0.0', '重启主机或重新进入Game Library功能', 'C', 'Major', 'Compatibility', 'Often', 'v1.2.0', 'IT', 'TC-PS5-2949', 'dev_johnson', 'S', 'QA-California', 'Game Library > User profile', '', '', 'Frontend', '', 'PS5 Slim', 'Kim', 'InProgress', 1, 1, '2025-12-08 12:00:00', '2025-12-10 10:00:00'),
('001206-' || hex(randomblob(8)), '001206', 35, '[VR Support] Home screen issue when resuming from suspend', 'VR Support', '## 问题描述
在使用 VR Support 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 VR Support 菜单
3. 执行 resuming from suspend 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0-beta', '重启主机或重新进入VR Support功能', 'A', 'Major', 'Environment', 'Often', 'v1.2.0', 'UT', 'TC-PS5-5035', 'dev_wang', 'C', 'QA-California', 'VR Support > Home screen', '', '', 'Backend', '', 'PS5 Slim', 'Tanaka', 'Reopened', 1, 1, '2026-01-21 11:00:00', '2026-01-22 07:00:00'),
('001207-' || hex(randomblob(8)), '001207', 35, '[User Account] Home screen issue when switching user', 'User Account', '## 问题描述
在使用 User Account 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.2.0', '重启主机或重新进入User Account功能', 'A', 'Minor', 'Compatibility', 'Once', 'v2.0.0-beta', 'IT', 'TC-PS5-8153', 'dev_suzuki', 'B', 'QA-California', 'User Account > Home screen', '', '', 'Backend', '', 'PS5 Slim', 'Li Ming', 'Reopened', 1, 1, '2025-11-18 08:00:00', '2025-11-19 00:00:00'),
('001208-' || hex(randomblob(8)), '001208', 35, '[Network] Home screen issue when voice chat active', 'Network', '## 问题描述
在使用 Network 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Network 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0-beta', '重启主机或重新进入Network功能', 'C', 'Trivial', 'Security', 'Often', 'v1.0.1', 'UAT', 'TC-PS5-8724', 'dev_kim', 'C', 'QA-California', 'Network > Home screen', 'v2.0.0-beta', '', 'Firmware', 'Configuration issue resolved', 'PS5 Standard Edition', 'John Smith', 'Resolved', 1, 1, '2025-11-13 06:00:00', '2025-11-13 08:00:00'),
('001209-' || hex(randomblob(8)), '001209', 35, '[Controller] Store page issue when resuming from suspend', 'Controller', '## 问题描述
在使用 Controller 功能时，Store page 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Controller 菜单
3. 执行 resuming from suspend 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.2.0', '重启主机或重新进入Controller功能', 'D', 'Major', 'UI', 'Rarely', 'v1.0.1', 'ST', 'TC-PS5-6258', 'dev_suzuki', 'S', 'QA-California', 'Controller > Store page', '', '', 'Backend', '', 'PS5 Digital Edition', 'Yamamoto', 'New', 1, 1, '2026-01-07 10:00:00', '2026-01-10 07:00:00'),
('001210-' || hex(randomblob(8)), '001210', 35, '[Storage Management] Trophy list issue when loading game', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Trophy list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 loading game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.0.1', '重启主机或重新进入Storage Management功能', 'A', 'Trivial', 'UIInteraction', 'Always', 'v1.1.0', 'UAT', 'TC-PS5-1404', 'dev_kim', 'S', 'QA-Shanghai', 'Storage Management > Trophy list', '', '', 'Frontend', '', 'PS5 Pro', 'Tanaka', 'Reopened', 1, 1, '2026-01-02 14:00:00', '2026-01-02 19:00:00'),
('001211-' || hex(randomblob(8)), '001211', 35, '[PlayStation Store] Friends list issue when streaming gameplay', 'PlayStation Store', '## 问题描述
在使用 PlayStation Store 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 PlayStation Store 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.2.0', '重启主机或重新进入PlayStation Store功能', 'A', 'Critical', 'Compatibility', 'Rarely', 'v2.0.0', 'Regression', 'TC-PS5-8491', 'dev_tanaka', 'S', 'Beta Tester', 'PlayStation Store > Friends list', 'v1.2.0', '', 'Frontend', 'Configuration issue resolved', 'PS5 Pro', 'Zhang Wei', 'Resolved', 1, 1, '2025-12-13 11:00:00', '2025-12-16 03:00:00'),
('001212-' || hex(randomblob(8)), '001212', 35, '[Remote Play] Capture gallery issue when installing game', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Capture gallery 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 installing game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Remote Play功能', 'B', 'Trivial', 'Compatibility', 'Sometimes', 'v1.0.0', 'IT', 'TC-PS5-8270', 'dev_suzuki', 'B', 'QA-California', 'Remote Play > Capture gallery', '', '', 'Firmware', '', 'PS5 Standard Edition', 'Zhang Wei', 'New', 1, 1, '2025-12-08 22:00:00', '2025-12-11 01:00:00'),
('001213-' || hex(randomblob(8)), '001213', 35, '[User Account] Download queue issue when starting application', 'User Account', '## 问题描述
在使用 User Account 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入User Account功能', 'A', 'Critical', 'UIInteraction', 'Once', 'v1.0.0', 'UT', 'TC-PS5-2975', 'dev_johnson', 'B', 'QA-California', 'User Account > Download queue', '', '', 'Backend', '', 'PS5 Slim', 'Zhang Wei', 'InProgress', 1, 1, '2025-12-03 06:00:00', '2025-12-04 03:00:00'),
('001214-' || hex(randomblob(8)), '001214', 35, '[Audio System] Game menu issue when capturing screenshot', 'Audio System', '## 问题描述
在使用 Audio System 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Audio System 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.2.0', '重启主机或重新进入Audio System功能', 'D', 'Major', 'Compatibility', 'Often', 'v1.1.0', 'Regression', 'TC-PS5-3865', 'dev_kim', 'A', 'Dev-Test', 'Audio System > Game menu', '', '', 'Graphics Engine', '', 'PS5 Standard Edition', 'Tanaka', 'New', 1, 1, '2025-12-22 18:00:00', '2025-12-25 11:00:00');