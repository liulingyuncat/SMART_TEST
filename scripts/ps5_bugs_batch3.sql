-- Batch 3: Defects 1135 to 1144
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001135-' || hex(randomblob(8)), '001135', 35, '[Storage Management] Store page issue when voice chat active', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Store page 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.0.1', '重启主机或重新进入Storage Management功能', 'D', 'Minor', 'Compatibility', 'Once', 'v1.0.1', 'IT', 'TC-PS5-5275', 'dev_chen', 'C', 'QA-California', 'Storage Management > Store page', '', '', 'Graphics Engine', '', 'PS5 Pro', 'Tanaka', 'InProgress', 1, 1, '2025-11-24 15:00:00', '2025-11-25 01:00:00'),
('001136-' || hex(randomblob(8)), '001136', 35, '[Storage Management] Party interface issue when loading game', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 loading game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.0.0', '重启主机或重新进入Storage Management功能', 'C', 'Critical', 'Functional', 'Sometimes', 'v1.0.1', 'IT', 'TC-PS5-6663', 'dev_wang', 'C', 'Beta Tester', 'Storage Management > Party interface', '', '', 'Frontend', '', 'PS5 Slim', 'Kim', 'Reopened', 1, 1, '2025-11-17 12:00:00', '2025-11-19 11:00:00'),
('001137-' || hex(randomblob(8)), '001137', 35, '[Share/Capture] Party interface issue when starting application', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入Share/Capture功能', 'A', 'Major', 'Performance', 'Once', 'v2.0.0-beta', 'UAT', 'TC-PS5-9292', 'dev_suzuki', 'B', 'QA-Tokyo', 'Share/Capture > Party interface', '', '', 'Audio Engine', '', 'PS5 Pro', 'Tanaka', 'Reopened', 1, 1, '2025-11-25 18:00:00', '2025-11-28 05:00:00'),
('001138-' || hex(randomblob(8)), '001138', 35, '[Storage Management] Trophy list issue when voice chat active', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Trophy list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.1.0', '重启主机或重新进入Storage Management功能', 'D', 'Trivial', 'UI', 'Sometimes', 'v2.0.0-beta', 'IT', 'TC-PS5-1272', 'dev_johnson', 'B', 'QA-Tokyo', 'Storage Management > Trophy list', 'v1.1.0', '', 'Frontend', 'Configuration issue resolved', 'PS5 Standard Edition', 'Kim', 'Closed', 1, 1, '2026-01-12 15:00:00', '2026-01-15 03:00:00'),
('001139-' || hex(randomblob(8)), '001139', 35, '[Audio System] Capture gallery issue when loading game', 'Audio System', '## 问题描述
在使用 Audio System 功能时，Capture gallery 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Audio System 菜单
3. 执行 loading game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v2.0.0', '重启主机或重新进入Audio System功能', 'C', 'Major', 'UI', 'Sometimes', 'v2.0.0', 'IT', 'TC-PS5-2772', 'dev_kim', 'B', 'QA-Shanghai', 'Audio System > Capture gallery', 'v1.2.0', '', 'Frontend', 'Code patch applied', 'PS5 Standard Edition', 'Sato', 'Resolved', 1, 1, '2025-12-05 05:00:00', '2025-12-05 11:00:00'),
('001140-' || hex(randomblob(8)), '001140', 35, '[Game Library] Home screen issue when downloading content', 'Game Library', '## 问题描述
在使用 Game Library 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入Game Library功能', 'B', 'Minor', 'UIInteraction', 'Always', 'v2.0.0', 'Regression', 'TC-PS5-7904', 'dev_johnson', 'A', 'QA-Shanghai', 'Game Library > Home screen', '', '', 'Storage Driver', '', 'PS5 Slim', 'John Smith', 'New', 1, 1, '2026-01-15 22:00:00', '2026-01-16 04:00:00'),
('001141-' || hex(randomblob(8)), '001141', 35, '[Remote Play] Capture gallery issue when connecting to server', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Capture gallery 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 connecting to server 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.2.0', '重启主机或重新进入Remote Play功能', 'C', 'Major', 'Compatibility', 'Always', 'v1.2.0', 'Regression', 'TC-PS5-1743', 'dev_tanaka', 'A', 'Beta Tester', 'Remote Play > Capture gallery', '', '', 'Firmware', '', 'PS5 Standard Edition', 'Yamamoto', 'Reopened', 1, 1, '2025-12-25 20:00:00', '2025-12-27 17:00:00'),
('001142-' || hex(randomblob(8)), '001142', 35, '[Share/Capture] Friends list issue when voice chat active', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入Share/Capture功能', 'B', 'Minor', 'UI', 'Often', 'v1.0.1', 'UAT', 'TC-PS5-6273', 'dev_wang', 'C', 'Beta Tester', 'Share/Capture > Friends list', '', '', 'Firmware', '', 'PS5 Slim', 'Emily Chen', 'Reopened', 1, 1, '2025-12-16 03:00:00', '2025-12-18 19:00:00'),
('001143-' || hex(randomblob(8)), '001143', 35, '[Network] Game menu issue when starting application', 'Network', '## 问题描述
在使用 Network 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Network 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Network功能', 'D', 'Trivial', 'Performance', 'Once', 'v1.0.1', 'IT', 'TC-PS5-1617', 'dev_wang', 'B', 'QA-Tokyo', 'Network > Game menu', '', '', 'Firmware', '', 'PS5 Pro', 'Emily Chen', 'Confirmed', 1, 1, '2026-01-12 09:00:00', '2026-01-14 06:00:00'),
('001144-' || hex(randomblob(8)), '001144', 35, '[Storage Management] Game menu issue when updating system', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 updating system 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入Storage Management功能', 'B', 'Critical', 'Security', 'Once', 'v1.0.1', 'UT', 'TC-PS5-3110', 'dev_suzuki', 'S', 'QA-California', 'Storage Management > Game menu', '', '', 'Network Stack', '', 'PS5 Digital Edition', 'Sato', 'InProgress', 1, 1, '2025-12-05 01:00:00', '2025-12-06 20:00:00');