-- 清理旧的需求管理相关表
-- 注意：执行前请确保数据已迁移到新表(requirement_items, viewpoint_items)

-- 1. 删除旧的 requirements 表
DROP TABLE IF EXISTS requirements;

-- 说明：
-- 旧表 requirements 包含以下字段（已废弃）:
--   - overall_requirements: 整体需求 (已迁移到 requirement_items)
--   - overall_test_viewpoint: 整体测试观点 (已迁移到 viewpoint_items)
--   - change_requirements: 变更需求 (已迁移到 requirement_items)
--   - change_test_viewpoint: 变更测试观点 (已迁移到 viewpoint_items)
--
-- 新架构使用：
--   - requirement_items: 动态需求条目列表
--   - viewpoint_items: 动态AI观点条目列表
--   - versions: 统一版本管理（通过 item_type 字段区分类型）
