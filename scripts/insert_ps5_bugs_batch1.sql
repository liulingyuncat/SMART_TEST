-- Batch 1: 10 PS5 Defects
INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES
('d1-' || hex(randomblob(16)), '000001', 35, '[Dashboard UI] Home screen crashes when loading game', 'Dashboard UI', '## 问题描述
在使用 Dashboard UI 功能时，Home screen 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Dashboard UI 菜单
3. 执行 loading game 操作
4. 观察到问题发生

## 预期结果
功能应正常工作，无错误提示。

## 实际结果
出现异常行为，影响用户体验。

## 测试环境
- 机型: PS5 Standard Edition
- 固件版本: v1.0.0
- 网络: 有线连接', '重启主机或重新进入Dashboard UI功能', 'A', 'Critical', 'Functional', 'Always', 'v1.0.0', 'UT', 'TC-PS5-1234', 'dev_tanaka', 'S', 'QA-Tokyo', 'Dashboard UI > Home screen', '', '', 'Frontend', '', 'PS5 Standard Edition', 'Tanaka', 'New', 1, 1, '2025-11-15 10:30:00', '2025-11-16 12:00:00'),

('d2-' || hex(randomblob(16)), '000002', 35, '[Game Library] Cannot download content in offline mode', 'Game Library', '## 问题描述
在使用 Game Library 功能时，Download queue 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 Game Library 菜单
3. 执行 downloading content 操作
4. 观察到问题发生', '重启主机或重新进入Game Library功能', 'B', 'Major', 'UI', 'Often', 'v1.0.1', 'IT', 'TC-PS5-2345', 'dev_suzuki', 'A', 'QA-Shanghai', 'Game Library > Download queue', 'v1.1.0', 'Fixed in firmware update', 'Backend', 'Code patch applied', 'PS5 Digital Edition', 'Zhang Wei', 'Resolved', 1, 1, '2025-11-20 14:20:00', '2025-11-22 16:30:00'),

('d3-' || hex(randomblob(16)), '000003', 35, '[PlayStation Store] Store page display issue on browsing store', 'PlayStation Store', '## 问题描述
在使用 PlayStation Store 功能时，Store page 发生异常。', '重启主机或重新进入PlayStation Store功能', 'C', 'Minor', 'UIInteraction', 'Sometimes', 'v1.1.0', 'ST', 'TC-PS5-3456', 'dev_wang', 'B', 'QA-California', 'PlayStation Store > Store page', '', 'Needs further testing', 'Frontend', '', 'PS5 Slim', 'Li Ming', 'InProgress', 1, 1, '2025-12-01 09:15:00', '2025-12-02 11:00:00'),

('d4-' || hex(randomblob(16)), '000004', 35, '[Settings] Performance drop when updating system', 'Settings', '## 问题描述
在使用 Settings 功能时，Settings panel 发生异常。', '重启主机或重新进入Settings功能', 'A', 'Critical', 'Performance', 'Rarely', 'v1.2.0', 'UAT', 'TC-PS5-4567', 'dev_chen', 'S', 'Dev-Test', 'Settings > Settings panel', '', '', 'Backend', '', 'PS5 Pro', 'John Smith', 'Confirmed', 1, 1, '2025-12-10 13:45:00', '2025-12-11 15:20:00'),

('d5-' || hex(randomblob(16)), '000005', 35, '[Network] Network error when connecting to server', 'Network', '## 问题描述
在使用 Network 功能时，Friends list 发生异常。', '重启主机或重新进入Network功能', 'B', 'Major', 'Functional', 'Always', 'v2.0.0-beta', 'Regression', 'TC-PS5-5678', 'dev_johnson', 'A', 'Beta Tester', 'Network > Friends list', 'v2.0.0', 'Configuration issue resolved', 'Network Stack', 'Network protocol updated', 'PS5 Standard Edition', 'Emily Chen', 'Closed', 1, 1, '2025-12-15 08:30:00', '2025-12-18 10:00:00'),

('d6-' || hex(randomblob(16)), '000006', 35, '[Controller] Controller button not responding', 'Controller', '## 问题描述
在使用 Controller 功能时，Game menu 发生异常。', '重启主机或重新进入Controller功能', 'D', 'Trivial', 'Compatibility', 'Once', 'v1.0.0', 'UT', 'TC-PS5-6789', 'dev_kim', 'C', 'QA-Tokyo', 'Controller > Game menu', '', 'Customer reported issue', 'Firmware', '', 'PS5 Digital Edition', 'Sato', 'Reopened', 1, 1, '2025-12-20 16:00:00', '2025-12-21 18:30:00'),

('d7-' || hex(randomblob(16)), '000007', 35, '[Audio System] Audio glitch during voice chat active', 'Audio System', '## 问题描述
在使用 Audio System 功能时，Party interface 发生异常。', '重启主机或重新进入Audio System功能', 'A', 'Major', 'Performance', 'Often', 'v1.0.1', 'IT', 'TC-PS5-7890', 'dev_tanaka', 'S', 'QA-Shanghai', 'Audio System > Party interface', '', '', 'Audio Engine', '', 'PS5 Slim', 'Kim', 'New', 1, 1, '2026-01-05 11:20:00', '2026-01-06 13:00:00'),

('d8-' || hex(randomblob(16)), '000008', 35, '[Video Output] Trophy list freezes after syncing trophies', 'Video Output', '## 问题描述
在使用 Video Output 功能时，Trophy list 发生异常。', '重启主机或重新进入Video Output功能', 'C', 'Minor', 'UI', 'Sometimes', 'v1.1.0', 'ST', 'TC-PS5-8901', 'dev_suzuki', 'B', 'QA-California', 'Video Output > Trophy list', 'v1.2.0', 'Fixed in firmware update', 'Graphics Engine', 'Code patch applied', 'PS5 Pro', 'Tanaka', 'Resolved', 1, 1, '2026-01-10 14:50:00', '2026-01-12 16:20:00'),

('d9-' || hex(randomblob(16)), '000009', 35, '[Storage Management] Incorrect User profile in capture gallery screen', 'Storage Management', '## 问题描述
在使用 Storage Management 功能时，User profile 发生异常。', '重启主机或重新进入Storage Management功能', 'B', 'Major', 'Functional', 'Rarely', 'v1.2.0', 'UAT', 'TC-PS5-9012', 'dev_wang', 'A', 'Dev-Test', 'Storage Management > User profile', '', 'Priority escalated by PM', 'Storage Driver', '', 'PS5 Standard Edition', 'Yamamoto', 'InProgress', 1, 1, '2026-01-15 09:30:00', '2026-01-16 11:00:00'),

('d10-' || hex(randomblob(16)), '000010', 35, '[User Account] Memory leak in notification system', 'User Account', '## 问题描述
在使用 User Account 功能时，Capture gallery 发生异常。', '重启主机或重新进入User Account功能', 'A', 'Critical', 'Performance', 'Always', 'v2.0.0', 'Regression', 'TC-PS5-0123', 'dev_chen', 'S', 'Beta Tester', 'User Account > Capture gallery', '', '', 'Backend', '', 'PS5 Digital Edition', 'Zhang Wei', 'Confirmed', 1, 1, '2026-01-20 12:15:00', '2026-01-21 14:45:00');
