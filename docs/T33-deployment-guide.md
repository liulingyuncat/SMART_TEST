# T33 自动化测试用例版本管理 - 部署指南

## 部署前检查清单

### 1. 环境要求
- ✅ Go 1.21+
- ✅ Node.js 16+
- ✅ SQLite 3
- ✅ 磁盘空间：至少500MB用于版本文件存储

### 2. 后端部署步骤

#### 2.1 创建存储目录
```powershell
# Windows PowerShell
cd D:\VSCode\webtest\backend
New-Item -ItemType Directory -Path "storage\versions\auto-cases" -Force

# Linux/Mac
cd /path/to/webtest/backend
mkdir -p storage/versions/auto-cases
chmod 755 storage/versions/auto-cases
```

#### 2.2 执行数据库迁移

**方式1：使用SQLite迁移脚本（推荐）**
```powershell
# Windows
cd D:\VSCode\webtest\backend
Get-Content migrations\009_create_auto_test_case_versions_table_sqlite.sql | sqlite3 webtest.db

# Linux/Mac
cd /path/to/webtest/backend
sqlite3 webtest.db < migrations/009_create_auto_test_case_versions_table_sqlite.sql
```

**方式2：GORM自动迁移（已包含在启动代码中）**
```go
// cmd/server/main.go 中已包含
db.AutoMigrate(&models.AutoTestCaseVersion{})
```

#### 2.3 验证数据库表创建
```powershell
sqlite3 webtest.db "SELECT name FROM sqlite_master WHERE type='table' AND name='auto_test_case_versions';"
# 预期输出: auto_test_case_versions

sqlite3 webtest.db "SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='auto_test_case_versions';"
# 预期输出: 
# idx_auto_versions_project_version
# idx_auto_versions_created
# idx_auto_versions_role
```

#### 2.4 编译后端
```powershell
# Windows
cd D:\VSCode\webtest\backend
go build -o server.exe cmd/server/main.go

# Linux/Mac
cd /path/to/webtest/backend
go build -o server cmd/server/main.go
```

#### 2.5 启动后端服务
```powershell
# Windows - 前台运行
.\server.exe

# Windows - 后台运行
Start-Process -FilePath ".\server.exe" -WindowStyle Hidden

# Linux/Mac - 前台运行
./server

# Linux/Mac - 后台运行
nohup ./server > server.log 2>&1 &
```

#### 2.6 验证服务启动
```powershell
# 检查端口占用
netstat -ano | findstr :8080  # Windows
netstat -tuln | grep :8080    # Linux/Mac

# 测试健康检查
curl http://localhost:8080/api/v1/projects
```

### 3. 前端部署步骤

#### 3.1 安装依赖
```powershell
cd D:\VSCode\webtest\frontend
npm install
```

#### 3.2 开发模式启动（测试用）
```powershell
npm start
# 访问 http://localhost:3000
```

#### 3.3 生产构建
```powershell
npm run build
# 构建产物在 build/ 目录
```

#### 3.4 部署到Web服务器

**Nginx配置示例：**
```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    # 前端静态文件
    location / {
        root /var/www/webtest/frontend/build;
        try_files $uri /index.html;
    }
    
    # 后端API代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Apache配置示例：**
```apache
<VirtualHost *:80>
    ServerName yourdomain.com
    DocumentRoot /var/www/webtest/frontend/build
    
    <Directory /var/www/webtest/frontend/build>
        Options -Indexes +FollowSymLinks
        AllowOverride All
        Require all granted
        RewriteEngine On
        RewriteCond %{REQUEST_FILENAME} !-f
        RewriteCond %{REQUEST_FILENAME} !-d
        RewriteRule . /index.html [L]
    </Directory>
    
    ProxyPass /api/ http://localhost:8080/api/
    ProxyPassReverse /api/ http://localhost:8080/api/
</VirtualHost>
```

### 4. 功能验证测试

#### 4.1 后端API测试
```powershell
# 1. 登录获取token
$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
    -Method POST `
    -ContentType "application/json" `
    -Body '{"username":"admin","password":"admin123"}'

$token = $loginResponse.data.token

# 2. 测试保存版本
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/projects/1/auto-cases/versions" `
    -Method POST `
    -Headers @{ "Authorization" = "Bearer $token" }

# 3. 获取版本列表
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/projects/1/auto-cases/versions?page=1&size=10" `
    -Method GET `
    -Headers @{ "Authorization" = "Bearer $token" }
```

#### 4.2 前端E2E测试
参考文档：`docs/T33-auto-version-e2e-test.md`

关键测试场景：
1. ✅ 版本保存与自动跳转
2. ✅ 版本列表展示
3. ✅ 下载版本压缩包
4. ✅ 备注编辑
5. ✅ 删除版本

### 5. 性能优化建议

#### 5.1 数据库优化
```sql
-- 定期清理旧版本（保留最近180天）
DELETE FROM auto_test_case_versions 
WHERE created_at < datetime('now', '-180 days');

-- 重建索引
REINDEX idx_auto_versions_project_version;
REINDEX idx_auto_versions_created;
```

#### 5.2 存储空间管理
```powershell
# 查看存储使用情况
Get-ChildItem -Path "storage\versions\auto-cases" -Recurse | 
    Measure-Object -Property Length -Sum | 
    Select-Object @{Name="TotalSizeMB";Expression={[math]::Round($_.Sum/1MB, 2)}}

# 清理孤儿文件（数据库中不存在的文件）
# 建议定期执行
```

#### 5.3 并发性能
- 批量保存版本使用goroutine并发处理4个ROLE
- 预期性能：1000条用例 < 5秒
- 如需提升，可调整worker数量或使用消息队列异步处理

### 6. 监控和日志

#### 6.1 后端日志
```go
// 日志位置：stdout（可重定向到文件）
// 关键日志：
// - 版本保存成功/失败
// - 文件生成错误
// - 数据库操作错误
```

#### 6.2 监控指标
建议监控：
- 版本保存成功率
- 平均保存时间
- 存储空间使用率
- API响应时间（/auto-cases/versions）

### 7. 回滚方案

#### 7.1 数据库回滚
```sql
-- 删除版本表
DROP TABLE IF EXISTS auto_test_case_versions;

-- 删除索引
DROP INDEX IF EXISTS idx_auto_versions_project_version;
DROP INDEX IF EXISTS idx_auto_versions_created;
DROP INDEX IF EXISTS idx_auto_versions_role;
```

#### 7.2 代码回滚
```bash
# 回退到上一个Git提交
git revert HEAD

# 或回退到特定版本
git reset --hard <commit-hash>

# 重新编译
go build -o server.exe cmd/server/main.go
```

### 8. 常见问题排查

#### 问题1：编译错误 "s.repo.GetDB undefined"
**原因**：service未传递db参数  
**解决**：确保`NewAutoTestCaseService(repo, projectService, db)`

#### 问题2：迁移失败 "SERIAL not supported"
**原因**：使用了PostgreSQL语法但数据库是SQLite  
**解决**：使用`009_create_auto_test_case_versions_table_sqlite.sql`

#### 问题3：存储目录权限错误
**原因**：目录不存在或无写权限  
**解决**：
```powershell
New-Item -ItemType Directory -Path "storage\versions\auto-cases" -Force
icacls "storage\versions\auto-cases" /grant Users:F
```

#### 问题4：前端下载失败 "Network Error"
**原因**：CORS配置或后端未启动  
**解决**：检查后端CORS设置，确保允许前端域名

#### 问题5：zip文件为空
**原因**：文件路径不正确或文件不存在  
**解决**：检查`file_path`字段，验证物理文件存在

### 9. 安全检查

#### 9.1 路径遍历防护
```go
// 已实现：validateFilePath函数
// 防止../等路径遍历攻击
```

#### 9.2 权限控制
- ✅ 所有API需要PM或PM Member权限
- ✅ JWT token验证
- ✅ 项目成员验证

#### 9.3 文件大小限制
建议在Nginx/Apache配置：
```nginx
client_max_body_size 50M;
```

### 10. 维护计划

#### 日常维护
- 每周检查存储空间使用情况
- 每月清理180天前的旧版本

#### 定期任务
- 每季度审查权限配置
- 每半年备份重要版本文件

#### 紧急响应
- 磁盘空间不足：立即清理旧版本或扩容
- 服务崩溃：查看日志，重启服务

---

## 部署验证清单

部署完成后，请逐项验证：

- [ ] 后端编译成功，无错误
- [ ] 数据库表`auto_test_case_versions`已创建
- [ ] 3个索引已创建
- [ ] 存储目录`storage/versions/auto-cases`已创建
- [ ] 后端服务启动成功（端口8080）
- [ ] 5个版本管理API路由已注册
- [ ] 前端构建成功，无错误
- [ ] 前端可访问，组件渲染正常
- [ ] 登录功能正常
- [ ] 版本保存功能测试通过
- [ ] 版本列表显示正常
- [ ] 版本下载功能正常
- [ ] 备注编辑功能正常
- [ ] 版本删除功能正常
- [ ] 日志输出正常
- [ ] 性能测试达标（1000用例<5秒）

---

**部署完成日期**: ___________  
**部署人员**: ___________  
**验证人员**: ___________  
**备注**: ___________
