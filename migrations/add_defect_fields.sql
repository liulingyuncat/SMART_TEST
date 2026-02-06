-- 添加缺陷管理新字段的数据库迁移脚本
-- 执行日期: 2026-02-04

-- 1. 修改Severity字段长度以支持新的枚举值
ALTER TABLE defects MODIFY COLUMN severity VARCHAR(20) DEFAULT 'Major';

-- 2. 添加新字段
ALTER TABLE defects ADD COLUMN IF NOT EXISTS type VARCHAR(30) COMMENT '缺陷类型';
ALTER TABLE defects ADD COLUMN IF NOT EXISTS detected_version VARCHAR(50) COMMENT '发现版本';
ALTER TABLE defects ADD COLUMN IF NOT EXISTS recovery_rank VARCHAR(50) COMMENT '恢复等级';
ALTER TABLE defects ADD COLUMN IF NOT EXISTS detection_team VARCHAR(100) COMMENT '检测团队';
ALTER TABLE defects ADD COLUMN IF NOT EXISTS location VARCHAR(200) COMMENT '位置';
ALTER TABLE defects ADD COLUMN IF NOT EXISTS fix_version VARCHAR(50) COMMENT '修复版本';
ALTER TABLE defects ADD COLUMN IF NOT EXISTS sqa_memo TEXT COMMENT 'SQA备注';
ALTER TABLE defects ADD COLUMN IF NOT EXISTS component VARCHAR(100) COMMENT '组件';

-- 3. 数据迁移：将旧的detected_in_release复制到detected_version
UPDATE defects SET detected_version = detected_in_release WHERE detected_version IS NULL OR detected_version = '';

-- 4. 数据迁移：将旧的严重程度值映射到新值
UPDATE defects SET severity = 'Critical' WHERE severity = 'A';
UPDATE defects SET severity = 'Major' WHERE severity = 'B';
UPDATE defects SET severity = 'Minor' WHERE severity = 'C';
UPDATE defects SET severity = 'Trivial' WHERE severity = 'D';

-- 5. 数据迁移：将Active状态映射到InProgress
UPDATE defects SET status = 'InProgress' WHERE status = 'Active';

-- 注意：detected_in_release字段保留用于向后兼容，不删除
