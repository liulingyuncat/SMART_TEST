-- Batch 1: Defects 1115 to 1124
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001115-' || hex(randomblob(8)), '001115', 35, '[Remote Play] Friends list issue when installing game', 'Remote Play', '## 问题描述
在使用 Remote Play 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Remote Play 菜单
3. 执行 installing game 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.2.0', '重启主机或重新进入Remote Play功能', 'B', 'Critical', 'Functional', 'Always', 'v1.2.0', 'UT', 'TC-PS5-5996', 'dev_suzuki', 'S', 'Dev-Test', 'Remote Play > Friends list', 'v1.2.0', '', 'Firmware', 'Configuration issue resolved', 'PS5 Pro', 'Yamamoto', 'Resolved', 1, 1, '2026-01-24 22:00:00', '2026-01-26 09:00:00'),
('001116-' || hex(randomblob(8)), '001116', 35, '[Network] Home screen issue when streaming gameplay', 'Network', '## 问题描述
在使用 Network 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Network 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0-beta', '重启主机或重新进入Network功能', 'D', 'Trivial', 'Functional', 'Sometimes', 'v1.1.0', 'ST', 'TC-PS5-5571', 'dev_tanaka', 'S', 'Beta Tester', 'Network > Home screen', 'v1.2.0', '', 'Network Stack', 'Code patch applied', 'PS5 Slim', 'Kim', 'Resolved', 1, 1, '2026-01-02 08:00:00', '2026-01-05 01:00:00'),
('001117-' || hex(randomblob(8)), '001117', 35, '[Storage Management] Download queue issue when downloading content', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.0.0', '重启主机或重新进入Storage Management功能', 'D', 'Major', 'Environment', 'Once', 'v2.0.0-beta', 'Regression', 'TC-PS5-5042', 'dev_kim', 'S', 'QA-California', 'Storage Management > Download queue', '', '', 'Storage Driver', '', 'PS5 Digital Edition', 'Li Ming', 'Confirmed', 1, 1, '2025-11-20 00:00:00', '2025-11-21 23:00:00'),
('001118-' || hex(randomblob(8)), '001118', 35, '[Share/Capture] User profile issue when resuming from suspend', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，User profile 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 resuming from suspend 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0', '重启主机或重新进入Share/Capture功能', 'D', 'Major', 'Environment', 'Rarely', 'v1.0.0', 'UT', 'TC-PS5-3932', 'dev_tanaka', 'B', 'Beta Tester', 'Share/Capture > User profile', '', '', 'Graphics Engine', '', 'PS5 Slim', 'Li Ming', 'InProgress', 1, 1, '2026-01-24 14:00:00', '2026-01-24 23:00:00'),
('001119-' || hex(randomblob(8)), '001119', 35, '[Network] Capture gallery issue when updating system', 'Network', '## 问题描述
在使用 Network 功能时，Capture gallery 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Network 菜单
3. 执行 updating system 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.0.0', '重启主机或重新进入Network功能', 'B', 'Major', 'UI', 'Rarely', 'v2.0.0-beta', 'UAT', 'TC-PS5-9286', 'dev_chen', 'S', 'QA-Shanghai', 'Network > Capture gallery', 'v2.0.0', '', 'Frontend', 'Configuration issue resolved', 'PS5 Standard Edition', 'Li Ming', 'Closed', 1, 1, '2025-11-06 00:00:00', '2025-11-06 20:00:00'),
('001120-' || hex(randomblob(8)), '001120', 35, '[Party Chat] Download queue issue when switching user', 'Party Chat', '## 问题描述
在使用 Party Chat 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Party Chat 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.1.0', '重启主机或重新进入Party Chat功能', 'C', 'Trivial', 'Security', 'Once', 'v1.2.0', 'IT', 'TC-PS5-6954', 'dev_johnson', 'A', 'QA-California', 'Party Chat > Download queue', '', '', 'Storage Driver', '', 'PS5 Pro', 'Sato', 'New', 1, 1, '2025-11-04 02:00:00', '2025-11-04 07:00:00'),
('001121-' || hex(randomblob(8)), '001121', 35, '[Game Library] Capture gallery issue when entering rest mode', 'Game Library', '## 问题描述
在使用 Game Library 功能时，Capture gallery 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 entering rest mode 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0', '重启主机或重新进入Game Library功能', 'D', 'Trivial', 'Security', 'Often', 'v1.0.0', 'UT', 'TC-PS5-5297', 'dev_kim', 'C', 'Dev-Test', 'Game Library > Capture gallery', 'v1.0.0', '', 'Graphics Engine', 'Configuration issue resolved', 'PS5 Slim', 'Zhang Wei', 'Resolved', 1, 1, '2026-01-01 03:00:00', '2026-01-03 04:00:00'),
('001122-' || hex(randomblob(8)), '001122', 35, '[Network] Game menu issue when streaming gameplay', 'Network', '## 问题描述
在使用 Network 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Network 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.1.0', '重启主机或重新进入Network功能', 'A', 'Major', 'Security', 'Sometimes', 'v2.0.0', 'UAT', 'TC-PS5-7510', 'dev_wang', 'B', 'QA-California', 'Network > Game menu', '', '', 'Graphics Engine', '', 'PS5 Standard Edition', 'Li Ming', 'Confirmed', 1, 1, '2025-11-02 20:00:00', '2025-11-03 16:00:00'),
('001123-' || hex(randomblob(8)), '001123', 35, '[Storage Management] Settings panel issue when streaming gameplay', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，Settings panel 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Storage Management 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.0.1', '重启主机或重新进入Storage Management功能', 'D', 'Minor', 'Compatibility', 'Often', 'v1.0.0', 'IT', 'TC-PS5-5456', 'dev_tanaka', 'A', 'QA-California', 'Storage Management > Settings panel', '', '', 'Graphics Engine', '', 'PS5 Standard Edition', 'Sato', 'InProgress', 1, 1, '2025-12-18 08:00:00', '2025-12-18 17:00:00'),
('001124-' || hex(randomblob(8)), '001124', 35, '[Controller] Party interface issue when downloading content', 'Controller', '## 问题描述
在使用 Controller 功能时，Party interface 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Controller 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Controller功能', 'C', 'Trivial', 'Performance', 'Once', 'v1.2.0', 'UAT', 'TC-PS5-9539', 'dev_tanaka', 'B', 'Beta Tester', 'Controller > Party interface', '', '', 'Network Stack', '', 'PS5 Pro', 'Li Ming', 'Reopened', 1, 1, '2026-01-07 18:00:00', '2026-01-10 02:00:00');