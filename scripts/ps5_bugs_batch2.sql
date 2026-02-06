-- Batch 2: Defects 1125 to 1134
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001125-' || hex(randomblob(8)), '001125', 35, '[Controller] Capture gallery issue when entering rest mode', 'Controller', '## 问题描述
在使用 Controller 功能时，Capture gallery 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Controller 菜单
3. 执行 entering rest mode 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.2.0', '重启主机或重新进入Controller功能', 'D', 'Major', 'Performance', 'Often', 'v1.0.1', 'UT', 'TC-PS5-2293', 'dev_chen', 'A', 'QA-California', 'Controller > Capture gallery', '', '', 'Graphics Engine', '', 'PS5 Slim', 'Tanaka', 'Confirmed', 1, 1, '2025-11-20 03:00:00', '2025-11-20 21:00:00'),
('001126-' || hex(randomblob(8)), '001126', 35, '[Storage Management] Download queue issue when streaming gameplay', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.1.0', '重启主机或重新进入Storage Management功能', 'D', 'Major', 'Functional', 'Often', 'v1.0.1', 'IT', 'TC-PS5-7090', 'dev_wang', 'A', 'QA-Tokyo', 'Storage Management > Download queue', '', '', 'Backend', '', 'PS5 Slim', 'Emily Chen', 'InProgress', 1, 1, '2025-11-27 00:00:00', '2025-11-27 06:00:00'),
('001127-' || hex(randomblob(8)), '001127', 35, '[Remote Play] Trophy list issue when resuming from suspend', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Trophy list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 resuming from suspend 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v2.0.0', '重启主机或重新进入Remote Play功能', 'B', 'Trivial', 'UI', 'Once', 'v1.0.0', 'Regression', 'TC-PS5-7636', 'dev_tanaka', 'C', 'QA-Shanghai', 'Remote Play > Trophy list', '', '', 'Storage Driver', '', 'PS5 Standard Edition', 'John Smith', 'Reopened', 1, 1, '2025-11-14 19:00:00', '2025-11-16 04:00:00'),
('001128-' || hex(randomblob(8)), '001128', 35, '[Audio System] User profile issue when capturing screenshot', 'Audio System', '## 问题描述
在使用 Audio System 功能时，User profile 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Audio System 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.0.1', '重启主机或重新进入Audio System功能', 'A', 'Minor', 'UI', 'Always', 'v2.0.0-beta', 'UT', 'TC-PS5-6493', 'dev_suzuki', 'C', 'QA-Tokyo', 'Audio System > User profile', 'v2.0.0', '', 'Frontend', 'Configuration issue resolved', 'PS5 Slim', 'John Smith', 'Resolved', 1, 1, '2026-01-02 09:00:00', '2026-01-03 17:00:00'),
('001129-' || hex(randomblob(8)), '001129', 35, '[User Account] Download queue issue when updating system', 'User Account', '## 问题描述
在使用 User Account 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 updating system 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v2.0.0', '重启主机或重新进入User Account功能', 'B', 'Minor', 'Performance', 'Rarely', 'v1.0.1', 'ST', 'TC-PS5-7027', 'dev_wang', 'B', 'Beta Tester', 'User Account > Download queue', '', '', 'Firmware', '', 'PS5 Pro', 'Yamamoto', 'Confirmed', 1, 1, '2025-11-05 03:00:00', '2025-11-06 17:00:00'),
('001130-' || hex(randomblob(8)), '001130', 35, '[User Account] Game menu issue when capturing screenshot', 'User Account', '## 问题描述
在使用 User Account 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 User Account 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入User Account功能', 'B', 'Minor', 'UIInteraction', 'Often', 'v1.1.0', 'Regression', 'TC-PS5-1127', 'dev_wang', 'S', 'QA-Tokyo', 'User Account > Game menu', 'v2.0.0-beta', '', 'Network Stack', 'Code patch applied', 'PS5 Digital Edition', 'Emily Chen', 'Closed', 1, 1, '2026-01-21 12:00:00', '2026-01-24 07:00:00'),
('001131-' || hex(randomblob(8)), '001131', 35, '[Audio System] Home screen issue when capturing screenshot', 'Audio System', '## 问题描述
在使用 Audio System 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Audio System 菜单
3. 执行 capturing screenshot 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.0.0', '重启主机或重新进入Audio System功能', 'B', 'Critical', 'Functional', 'Always', 'v1.1.0', 'IT', 'TC-PS5-4200', 'dev_wang', 'A', 'QA-California', 'Audio System > Home screen', 'v1.2.0', '', 'Backend', 'Fixed in firmware update', 'PS5 Digital Edition', 'Sato', 'Resolved', 1, 1, '2025-11-22 09:00:00', '2025-11-25 06:00:00'),
('001132-' || hex(randomblob(8)), '001132', 35, '[VR Support] Trophy list issue when switching user', 'VR Support', '## 问题描述
在使用 VR Support 功能时，Trophy list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 VR Support 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v2.0.0-beta', '重启主机或重新进入VR Support功能', 'D', 'Major', 'Compatibility', 'Always', 'v1.0.1', 'ST', 'TC-PS5-6667', 'dev_tanaka', 'S', 'QA-Shanghai', 'VR Support > Trophy list', '', '', 'Backend', '', 'PS5 Digital Edition', 'Zhang Wei', 'New', 1, 1, '2025-12-28 19:00:00', '2025-12-30 17:00:00'),
('001133-' || hex(randomblob(8)), '001133', 35, '[Game Library] Settings panel issue when voice chat active', 'Game Library', '## 问题描述
在使用 Game Library 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.0.0', '重启主机或重新进入Game Library功能', 'A', 'Critical', 'Compatibility', 'Sometimes', 'v1.0.1', 'UAT', 'TC-PS5-2241', 'dev_johnson', 'B', 'QA-Tokyo', 'Game Library > Settings panel', 'v2.0.0', '', 'Storage Driver', 'Code patch applied', 'PS5 Slim', 'Yamamoto', 'Resolved', 1, 1, '2025-12-07 01:00:00', '2025-12-09 03:00:00'),
('001134-' || hex(randomblob(8)), '001134', 35, '[VR Support] Home screen issue when entering rest mode', 'VR Support', '## 问题描述
在使用 VR Support 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 VR Support 菜单
3. 执行 entering rest mode 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.2.0', '重启主机或重新进入VR Support功能', 'D', 'Minor', 'Environment', 'Always', 'v1.2.0', 'IT', 'TC-PS5-8664', 'dev_johnson', 'S', 'Dev-Test', 'VR Support > Home screen', '', '', 'Firmware', '', 'PS5 Pro', 'Tanaka', 'Reopened', 1, 1, '2025-11-01 07:00:00', '2025-11-02 11:00:00');