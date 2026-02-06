-- Batch 4: Defects 1145 to 1154
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('001145-' || hex(randomblob(8)), '001145', 35, '[Settings] Game menu issue when starting application', 'Settings', '## 问题描述
在使用 Settings 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Settings 菜单
3. 执行 starting application 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.1', '重启主机或重新进入Settings功能', 'B', 'Major', 'Compatibility', 'Often', 'v1.0.1', 'UT', 'TC-PS5-5860', 'dev_wang', 'S', 'Beta Tester', 'Settings > Game menu', 'v1.2.0', '', 'Frontend', 'Code patch applied', 'PS5 Pro', 'Kim', 'Resolved', 1, 1, '2026-01-01 15:00:00', '2026-01-03 09:00:00'),
('001146-' || hex(randomblob(8)), '001146', 35, '[Share/Capture] Friends list issue when switching user', 'Share/Capture', '## 问题描述
在使用 Share/Capture 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Share/Capture 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v2.0.0', '重启主机或重新进入Share/Capture功能', 'C', 'Trivial', 'Performance', 'Rarely', 'v1.2.0', 'UT', 'TC-PS5-7814', 'dev_chen', 'C', 'QA-California', 'Share/Capture > Friends list', '', '', 'Frontend', '', 'PS5 Slim', 'Emily Chen', 'Confirmed', 1, 1, '2026-01-14 08:00:00', '2026-01-15 20:00:00'),
('001147-' || hex(randomblob(8)), '001147', 35, '[VR Support] Game menu issue when streaming gameplay', 'VR Support', '## 问题描述
在使用 VR Support 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 VR Support 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0-beta', '重启主机或重新进入VR Support功能', 'A', 'Major', 'Environment', 'Rarely', 'v1.2.0', 'ST', 'TC-PS5-3205', 'dev_kim', 'A', 'QA-Tokyo', 'VR Support > Game menu', 'v1.0.1', '', 'Firmware', 'Code patch applied', 'PS5 Slim', 'John Smith', 'Resolved', 1, 1, '2025-11-25 08:00:00', '2025-11-27 05:00:00'),
('001148-' || hex(randomblob(8)), '001148', 35, '[Audio System] Download queue issue when streaming gameplay', 'Audio System', '## 问题描述
在使用 Audio System 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Audio System 菜单
3. 执行 streaming gameplay 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Digital Edition
- 固件版本: v1.1.0', '重启主机或重新进入Audio System功能', 'A', 'Major', 'Performance', 'Sometimes', 'v1.0.1', 'ST', 'TC-PS5-5806', 'dev_kim', 'S', 'QA-California', 'Audio System > Download queue', '', '', 'Backend', '', 'PS5 Standard Edition', 'Zhang Wei', 'InProgress', 1, 1, '2025-12-11 21:00:00', '2025-12-12 21:00:00'),
('001149-' || hex(randomblob(8)), '001149', 35, '[VR Support] Store page issue when voice chat active', 'VR Support', '## 问题描述
在使用 VR Support 功能时，Store page 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 VR Support 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0-beta', '重启主机或重新进入VR Support功能', 'C', 'Major', 'Environment', 'Often', 'v1.1.0', 'UT', 'TC-PS5-4608', 'dev_chen', 'A', 'QA-Tokyo', 'VR Support > Store page', '', '', 'Backend', '', 'PS5 Pro', 'Yamamoto', 'Reopened', 1, 1, '2025-11-11 18:00:00', '2025-11-12 08:00:00'),
('001150-' || hex(randomblob(8)), '001150', 35, '[VR Support] Friends list issue when voice chat active', 'VR Support', '## 问题描述
在使用 VR Support 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 VR Support 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Slim
- 固件版本: v1.2.0', '重启主机或重新进入VR Support功能', 'A', 'Minor', 'UIInteraction', 'Once', 'v1.1.0', 'UT', 'TC-PS5-5727', 'dev_wang', 'A', 'QA-Shanghai', 'VR Support > Friends list', '', '', 'Frontend', '', 'PS5 Digital Edition', 'John Smith', 'Confirmed', 1, 1, '2026-01-15 04:00:00', '2026-01-16 17:00:00'),
('001151-' || hex(randomblob(8)), '001151', 35, '[Network] Store page issue when downloading content', 'Network', '## 问题描述
在使用 Network 功能时，Store page 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Network 菜单
3. 执行 downloading content 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v2.0.0', '重启主机或重新进入Network功能', 'B', 'Trivial', 'Environment', 'Once', 'v1.0.0', 'ST', 'TC-PS5-4177', 'dev_chen', 'A', 'QA-California', 'Network > Store page', '', '', 'Graphics Engine', '', 'PS5 Slim', 'Zhang Wei', 'InProgress', 1, 1, '2025-11-10 04:00:00', '2025-11-11 06:00:00'),
('001152-' || hex(randomblob(8)), '001152', 35, '[Trophy System] Store page issue when voice chat active', 'Trophy System', '## 问题描述
在使用 Trophy System 功能时，Store page 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Trophy System 菜单
3. 执行 voice chat active 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v1.2.0', '重启主机或重新进入Trophy System功能', 'A', 'Minor', 'Performance', 'Once', 'v1.0.1', 'UAT', 'TC-PS5-3662', 'dev_kim', 'A', 'QA-Shanghai', 'Trophy System > Store page', '', '', 'Frontend', '', 'PS5 Pro', 'Kim', 'Reopened', 1, 1, '2025-11-26 14:00:00', '2025-11-27 13:00:00'),
('001153-' || hex(randomblob(8)), '001153', 35, '[PlayStation Store] Game menu issue when entering rest mode', 'PlayStation Store', '## 问题描述
在使用 PlayStation Store 功能时，Game menu 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 PlayStation Store 菜单
3. 执行 entering rest mode 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v2.0.0', '重启主机或重新进入PlayStation Store功能', 'B', 'Major', 'UI', 'Rarely', 'v2.0.0', 'UAT', 'TC-PS5-5009', 'dev_kim', 'B', 'QA-Shanghai', 'PlayStation Store > Game menu', '', '', 'Audio Engine', '', 'PS5 Slim', 'Zhang Wei', 'Confirmed', 1, 1, '2025-11-17 17:00:00', '2025-11-19 04:00:00'),
('001154-' || hex(randomblob(8)), '001154', 35, '[Video Output] Friends list issue when switching user', 'Video Output', '## 问题描述
在使用 Video Output 功能时，Friends list 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Video Output 菜单
3. 执行 switching user 操作
4. 观察到问题发生

## 测试环境
- 机型: PS5 Pro
- 固件版本: v2.0.0-beta', '重启主机或重新进入Video Output功能', 'D', 'Trivial', 'UIInteraction', 'Often', 'v1.0.0', 'IT', 'TC-PS5-3631', 'dev_kim', 'A', 'QA-California', 'Video Output > Friends list', '', '', 'Audio Engine', '', 'PS5 Pro', 'Zhang Wei', 'New', 1, 1, '2026-01-06 10:00:00', '2026-01-06 23:00:00');