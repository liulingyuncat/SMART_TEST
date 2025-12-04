-- 创建缺陷说明表
CREATE TABLE IF NOT EXISTS defect_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    defect_id VARCHAR(36) NOT NULL,
    content TEXT NOT NULL,
    created_by INTEGER NOT NULL,
    updated_by INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 创建索引
CREATE INDEX idx_defect_comments_defect_id ON defect_comments(defect_id);
CREATE INDEX idx_defect_comments_deleted_at ON defect_comments(deleted_at);

-- 外键约束（SQLite 3.6.19+）
-- 注：SQLite默认不启用外键，需在连接时设置 PRAGMA foreign_keys = ON;
-- 如果使用PostgreSQL，可启用外键约束：
-- ALTER TABLE defect_comments ADD CONSTRAINT fk_defect_comments_defect_id FOREIGN KEY (defect_id) REFERENCES defects(id) ON DELETE CASCADE;
-- ALTER TABLE defect_comments ADD CONSTRAINT fk_defect_comments_created_by FOREIGN KEY (created_by) REFERENCES users(id);
