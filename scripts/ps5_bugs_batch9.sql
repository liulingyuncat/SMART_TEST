-- Batch 9: Defects 1195 to 1204
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001195-' || hex(randomblob(8)), '001195', 35, '[Controller] Settings panel issue when entering rest mode', 'Controller', '## 问题描述
在使用 Controller 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Controller 菜单
3. 执行 entering rest mode 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入Controller功能', 'A', 'Trivial', 'UI', 'Rarely', 'v1.2.0', 'Regression', 'TC-PS5-4395', 'dev_wang', 'C', 'Dev-Test', 'Controller > Settings panel', '', '', 'Storage Driver', '', 'PS5 Standard Edition', 'John Smith', 'Confirmed', 1, 1, '2025-11-08 16:00:00', '2025-11-09 09:00:00'),
('001196-' || hex(randomblob(8)), '001196', 35, '[User Account] Home screen issue when starting application', 'User Account', '## 问题描述
在使用 User Account 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.0.1', '重启主机或重新进入User Account功能', 'A', 'Trivial', 'Performance', 'Once', 'v2.0.0-beta', 'IT', 'TC-PS5-1456', 'dev_suzuki', 'B', 'QA-Shanghai', 'User Account > Home screen', 'v1.0.1', '', 'Graphics Engine', 'Configuration issue resolved', 'PS5 Pro', 'Yamamoto', 'Resolved', 1, 1, '2025-12-16 14:00:00', '2025-12-16 20:00:00'),
('001197-' || hex(randomblob(8)), '001197', 35, '[Settings] Game menu issue when installing game', 'Settings', '## 问题描述
在使用 Settings 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Settings 菜单
3. 执行 installing game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.2.0', '重启主机或重新进入Settings功能', 'D', 'Trivial', 'Environment', 'Once', 'v1.0.1', 'UAT', 'TC-PS5-3388', 'dev_johnson', 'S', 'Dev-Test', 'Settings > Game menu', '', '', 'Network Stack', '', 'PS5 Standard Edition', 'John Smith', 'Confirmed', 1, 1, '2026-02-04 16:00:00', '2026-02-04 21:00:00'),
('001198-' || hex(randomblob(8)), '001198', 35, '[Remote Play] Store page issue when loading game', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Store page 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 loading game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.0.0', '重启主机或重新进入Remote Play功能', 'D', 'Major', 'UI', 'Once', 'v2.0.0-beta', 'ST', 'TC-PS5-6317', 'dev_suzuki', 'C', 'QA-California', 'Remote Play > Store page', '', '', 'Frontend', '', 'PS5 Standard Edition', 'Emily Chen', 'InProgress', 1, 1, '2025-11-05 18:00:00', '2025-11-06 01:00:00'),
('001199-' || hex(randomblob(8)), '001199', 35, '[User Account] Party interface issue when voice chat active', 'User Account', '## 问题描述
在使用 User Account 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.0.0', '重启主机或重新进入User Account功能', 'B', 'Critical', 'UIInteraction', 'Rarely', 'v1.0.1', 'UAT', 'TC-PS5-2653', 'dev_chen', 'S', 'Beta Tester', 'User Account > Party interface', 'v2.0.0-beta', '', 'Backend', 'Code patch applied', 'PS5 Standard Edition', 'Emily Chen', 'Resolved', 1, 1, '2025-12-25 14:00:00', '2025-12-27 22:00:00'),
('001200-' || hex(randomblob(8)), '001200', 35, '[PlayStation Store] Home screen issue when entering rest mode', 'PlayStation Store', '## 问题描述
在使用 PlayStation Store 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 PlayStation Store 菜单
3. 执行 entering rest mode 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入PlayStation Store功能', 'C', 'Minor', 'UIInteraction', 'Sometimes', 'v2.0.0-beta', 'IT', 'TC-PS5-2407', 'dev_suzuki', 'A', 'Beta Tester', 'PlayStation Store > Home screen', '', '', 'Audio Engine', '', 'PS5 Pro', 'Sato', 'InProgress', 1, 1, '2026-01-18 09:00:00', '2026-01-19 21:00:00'),
('001201-' || hex(randomblob(8)), '001201', 35, '[Game Library] User profile issue when capturing screenshot', 'Game Library', '## 问题描述
在使用 Game Library 功能时，User profile 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入Game Library功能', 'B', 'Trivial', 'Compatibility', 'Often', 'v1.0.0', 'ST', 'TC-PS5-2883', 'dev_johnson', 'S', 'QA-California', 'Game Library > User profile', '', '', 'Audio Engine', '', 'PS5 Slim', 'Tanaka', 'InProgress', 1, 1, '2026-01-11 03:00:00', '2026-01-13 09:00:00'),
('001202-' || hex(randomblob(8)), '001202', 35, '[PlayStation Store] Capture gallery issue when voice chat active', 'PlayStation Store', '## 问题描述
在使用 PlayStation Store 功能时，Capture gallery 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 PlayStation Store 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入PlayStation Store功能', 'B', 'Critical', 'Security', 'Once', 'v2.0.0-beta', 'IT', 'TC-PS5-9459', 'dev_tanaka', 'C', 'Dev-Test', 'PlayStation Store > Capture gallery', 'v1.2.0', '', 'Graphics Engine', 'Fixed in firmware update', 'PS5 Slim', 'Zhang Wei', 'Closed', 1, 1, '2026-01-18 20:00:00', '2026-01-20 15:00:00'),
('001203-' || hex(randomblob(8)), '001203', 35, '[Remote Play] Settings panel issue when resuming from suspend', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 resuming from suspend 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入Remote Play功能', 'C', 'Trivial', 'UI', 'Once', 'v1.0.0', 'UAT', 'TC-PS5-7530', 'dev_suzuki', 'S', 'Beta Tester', 'Remote Play > Settings panel', 'v1.1.0', '', 'Frontend', 'Code patch applied', 'PS5 Pro', 'Tanaka', 'Closed', 1, 1, '2025-12-06 13:00:00', '2025-12-06 21:00:00'),
('001204-' || hex(randomblob(8)), '001204', 35, '[VR Support] Capture gallery issue when downloading content', 'VR Support', '## 问题描述
在使用 VR Support 功能时，Capture gallery 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 VR Support 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.1.0', '重启主机或重新进入VR Support功能', 'B', 'Trivial', 'Security', 'Rarely', 'v2.0.0', 'ST', 'TC-PS5-1846', 'dev_wang', 'A', 'Dev-Test', 'VR Support > Capture gallery', 'v1.1.0', '', 'Backend', 'Configuration issue resolved', 'PS5 Pro', 'Li Ming', 'Closed', 1, 1, '2025-11-09 15:00:00', '2025-11-10 14:00:00');