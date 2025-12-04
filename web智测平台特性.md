# Web智能测试平台 - 特性文档

## 文档版本

- **版本**: v1.0
- **日期**: 2025年10月30日
- **状态**: 草案

---

## 一、技术架构

### 1.1 技术栈

| 类别    | 技术选型                                         |
| ----- | -------------------------------------------- |
| 后端语言  | Go                                           |
| Web框架 | Gin                                          |
| ORM框架 | GORM                                         |
| 前端框架  | React + Ant Design                           |
| 数据库   | MongoDB (默认) / PostgreSQL (可选) / SQLite (可选) |
| 认证方式  | JWT (JSON Web Token)                         |

### 1.2 系统架构特点

- 前后端分离架构
- RESTful API设计
- 基于角色的访问控制 (RBAC)
- 多数据库支持，灵活配置

---

## 二、核心特性

### Feature 1: 用户认证与授权系统

#### 2.1.1 多语言支持

- **特性ID**: F-AUTH-001
- **描述**: 支持中文(CN)、英文(EN)、日文(JP)三种语言切换
- **默认语言**: 简体中文(CN)
- **实现方式**: 前端国际化(i18n)

#### 2.1.2 JWT令牌认证

- **特性ID**: F-AUTH-002
- **描述**: 基于JWT的无状态认证机制
- **Token内容**: 包含用户身份、权限范围(scope)
- **安全特性**:
  - Token过期机制
  - 密码加密存储
  - 权限验证中间件

#### 2.1.3 用户角色体系

- **特性ID**: F-AUTH-003
- **角色类型**:
  - **系统管理员** (System Administrator)
    - 预置账号: `admin` (密码: `admin123!`), `root` (密码: `root123`)
    - 权限: 用户管理、全局配置
  - **项目管理员** (Project Manager)
    - 初始密码: `admin!123`
    - 权限: 项目管理、成员管理、人员分配
  - **项目成员** (Project Member)
    - 初始密码: `user!123`
    - 权限: 用例编辑、统计查看

#### 2.1.4 导航系统

- **特性ID**: F-AUTH-004
- **布局**: 侧边栏常驻导航
- **导航菜单** (根据角色显示):
  - 项目管理 (项目管理员、项目成员)
  - 项目统计 (项目管理员、项目成员)
  - 人员管理 (系统管理员、项目管理员)
  - 人员分配 (项目管理员)
  - 我的信息 (所有角色)

---

### Feature 2: 项目管理系统

#### 2.2.1 项目列表管理

- **特性ID**: F-PROJ-001
- **功能**:
  - 显示当前用户所属的项目列表
  - 项目创建、删除、重命名
  - 点击单个项目后,进入当前项目详细画面
- **权限控制**:
  - 项目管理员: 可创建、删除、编辑项目
  - 系统管理员: 不显示项目管理画面
  - 项目成员: 只读访问
  - 项目名唯一性
- **自动关联**: 创建项目时自动将创建者加入项目
- **项目详细画面布局**:
  - **上栏**: 显示当前项目名称
  - **中栏**: 横向导航栏菜单
    - 需求管理 (2.2.2)
    - 功能测试用例管理 (手工测试) (2.2.3)
    - 功能测试用例管理 (自动化测试) (2.2.4)
    - 接口测试用例管理 (2.2.5)
  - **下栏**: 根据所选导航菜单显示相应模块内容

#### 2.2.2 需求管理模块

- **特性ID**: F-PROJ-002
- **文档格式**: Markdown
- **功能**:
  - 在线编辑器支持
  - 实时预览
  - 文档下载 (.md格式)
  - 服务端持久化存储
- **编辑器特性**:
  - Markdown语法高亮
  - 工具栏快捷操作
  - 保存按钮持久化

#### 2.2.3 功能测试用例管理 (手工测试)

- **特性ID**: F-PROJ-003

- **界面布局**: 上下分栏

- **上栏元数据**:
  
  ```
  - 测试版本 (可编辑文本框)
  - 测试环境 (可编辑文本框)
  - 测试日期 (可编辑文本框)
  - 执行人 (可编辑文本框)
  ```

- **下栏表格字段**:
  
  | 字段   | 类型   | 可编辑 | 说明                    |
  | ---- | ---- | --- | --------------------- |
  | ID   | 自增序号 | 否   | 自动生成，支持重排             |
  | 用例编号 | 文本   | 是   | 用户自定义编号               |
  | 画面   | 文本   | 是   | 功能所属页面                |
  | 功能   | 文本   | 是   | 测试功能描述                |
  | 前置条件 | 文本   | 是   | 执行前提条件                |
  | 测试步骤 | 文本   | 是   | 详细操作步骤                |
  | 期待值  | 文本   | 是   | 预期结果                  |
  | 测试结果 | 下拉选择 | 是   | OK/NG/Block/NR (默认NR) |
  | 备考   | 文本   | 是   | 备注信息                  |

- **表格操作**:
  
  - 插入空行
  - 删除行
  - ID重新生成 (升序排列)
  - 保存编辑的数据
  - 分页显示 (50条/页)

- **导出功能**: 导出元数据 + 所有用例数据

#### 2.2.4 功能测试用例管理 (自动化测试)

- **特性ID**: F-PROJ-004

- **界面布局**: 上下分栏

- **上栏元数据**:
  
  ```
  - 测试版本 (可编辑文本框)
  - 测试日期 (可编辑文本框)
  ```

- **下栏表格字段**:
  
  | 字段   | 类型   | 可编辑 | 说明   |
  | ---- | ---- | --- | ---- |
  | ID   | 自增序号 | 否   | 自动变化 |
  | 用例编号 | 文本   | 是   | -    |
  | 画面   | 文本   | 是   | -    |
  | 功能   | 文本   | 是   | -    |
  | 前置条件 | 文本   | 是   | -    |
  | 测试步骤 | 文本   | 是   | -    |
  | 期待值  | 文本   | 是   | -    |
  | 测试结果 | 下拉选择 | 是   | 默认NR |
  | 备考   | 文本   | 是   | -    |

- **特性**:
  
  - ID重新生成
  - 支持插入、删除行
  - 分页显示 (50条/页)
  - 数据导出功能
  - 保存编辑的数据

#### 2.2.5 接口测试用例管理

- **特性ID**: F-PROJ-005

- **界面布局**: 上下分栏

- **上栏元数据**:
  
  ```
  - 测试版本 (可编辑文本框)
  - 测试日期 (可编辑文本框)
  ```

- **下栏表格字段**:
  
  | 字段       | 类型   | 可编辑 | 说明                 |
  | -------- | ---- | --- | ------------------ |
  | ID       | 自增序号 | 否   | 自动变化               |
  | 用例编号     | 文本   | 是   | -                  |
  | 画面       | 文本   | 是   | 接口所属模块             |
  | URL      | 文本   | 是   | 接口地址               |
  | Header   | 文本   | 是   | 请求头                |
  | Method   | 文本   | 是   | HTTP方法 (GET/POST等) |
  | Body     | 文本   | 是   | 请求体                |
  | Response | 文本   | 是   | 响应结果               |
  | 测试结果     | 下拉选择 | 是   | 默认NR               |
  | 备考       | 文本   | 是   | -                  |

- **特性**:
  
  - ID重新生成
  - 支持增删行
  - 分页显示 (50条/页)
  - 数据导出功能
  - 保存编辑的数据

---

### Feature 3: 项目统计与数据可视化

#### 2.3.1 统计数据看板

- **特性ID**: F-STAT-001
- **访问方式**: 项目列表 → 项目统计画面
- **可视化类型**: 饼图 (Pie Chart)

#### 2.3.2 手工测试用例统计

- **特性ID**: F-STAT-002
- **统计指标**:
  
  ```
  - 总用例数
  - OK数量
  - NG数量
  - Block数量
  - NR数量
  - 实施进度 = (OK + NG + Block) / 总数 × 100%
  - 通过率 = OK / (OK + NG + Block) × 100%
  ```

#### 2.3.3 自动化测试用例统计

- **特性ID**: F-STAT-003
- **统计指标**:
  
  ```
  - 总用例数
  - OK数量
  - NG数量
  - NR数量
  - 实施进度 = (OK + NG) / 总数 × 100%
  - 通过率 = OK / (OK + NG) × 100%
  ```

#### 2.3.4 接口测试用例统计

- **特性ID**: F-STAT-004
- **统计指标**:
  
  ```
  - 总用例数
  - OK数量
  - NG数量
  - NR数量
  - 实施进度 = (OK + NG) / 总数 × 100%
  - 通过率 = OK / (OK + NG) × 100%
  ```

#### 2.3.5 数据展示

- **特性ID**: F-STAT-005
- **布局**: 三个饼图并列显示
- **交互**: 实时数据更新，饼图下方显示详细数值

---

### Feature 4: 人员管理系统

#### 2.4.1 用户唯一性约束

- **特性ID**: F-USER-001
- **约束规则**:
  - 用户名全局唯一
  - 昵称全局唯一
  - 任一重复即创建失败

#### 2.4.2 系统管理员视图

- **特性ID**: F-USER-002

- **界面布局**: 双栏布局 (项目管理员栏 + 项目成员栏)

- **项目管理员管理**:
  
  - 创建项目管理员 (用户名、昵称)
  - 创建后的默认密码(→ `admin!123`)
  - 表格显示: 用户名 | 昵称 | 权限 | 操作
  - 操作: 密码重置 (→ `admin!123`) | 删除
  - 昵称可内联编辑，唯一性验证

- **项目成员管理**:
  
  - 创建项目成员 (用户名、昵称)
  - 创建后的默认密码(→ `user!123`)
  - 表格显示: 用户名 | 昵称 | 权限 | 操作
  - 操作: 密码重置 (→ `user!123`) | 删除
  - 昵称可内联编辑，唯一性验证

#### 2.4.3 项目管理员视图

- **特性ID**: F-USER-003
- **界面布局**: 单栏布局 (仅项目成员栏)
- **功能**: 与系统管理员的项目成员管理相同

#### 2.4.4 项目成员视图

- **特性ID**: F-USER-004
- **权限**: 不显示人员管理画面

---

### Feature 5: 人员分配系统

#### 2.5.1 访问控制

- **特性ID**: F-ASSIGN-001
- **可见性**: 仅项目管理员

#### 2.5.2 三栏式分配界面

- **特性ID**: F-ASSIGN-002

##### 第一栏: 项目选择器

- 下拉列表显示当前登录管理员所在的所有项目
- 选择项目后联动更新第二、三栏

##### 第二栏: 项目管理员分配

- **布局**: 左右分栏
  - 左侧: 当前项目的项目管理员
  - 右侧: 其他可用的项目管理员
- **操作**:
  - 支持双向移动
  - 当前登录用户不可移动 (锁定)

##### 第三栏: 项目成员分配

- **布局**: 左右分栏
  - 左侧: 当前项目的项目成员
  - 右侧: 其他可用的项目成员
- **操作**: 支持双向移动

#### 2.5.3 数据持久化

- **特性ID**: F-ASSIGN-003
- **保存机制**: 点击保存按钮后，批量更新项目成员关系

---

### Feature 6: 个人信息管理

#### 2.6.1 信息展示

- **特性ID**: F-PROFILE-001
- **字段**:
  - 用户名 (只读)
  - 昵称 (可编辑)，唯一性确认
  - 权限 (只读)

#### 2.6.2 昵称变更

- **特性ID**: F-PROFILE-002
- **流程**:
  1. 点击编辑按钮打开弹窗
  2. 输入新昵称
  3. 唯一性验证
  4. 保存成功/失败提示

#### 2.6.3 密码变更

- **特性ID**: F-PROFILE-003
- **流程**:
  1. 点击密码变更按钮打开弹窗
  2. 输入新密码
  3. 确认新密码
  4. 两次输入一致性验证
  5. 加密存储新密码

---

## 三、数据模型

### 3.1 核心实体

#### User (用户)

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"unique;not null"`
    Nickname  string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    Role      string    `gorm:"not null"` // system_admin, project_manager, project_member
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### Project (项目)

```go
type Project struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`
    Description string
    CreatedBy   uint      // User ID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### ProjectMember (项目成员关联)

```go
type ProjectMember struct {
    ID        uint `gorm:"primaryKey"`
    ProjectID uint `gorm:"not null"`
    UserID    uint `gorm:"not null"`
    Role      string // project_manager, project_member
}
```

#### Requirement (需求文档)

```go
type Requirement struct {
    ID        uint      `gorm:"primaryKey"`
    ProjectID uint      `gorm:"not null"`
    Content   string    `gorm:"type:text"`
    UpdatedAt time.Time
}
```

#### ManualTestCase (手工测试用例)

```go
type ManualTestCase struct {
    ID             uint   `gorm:"primaryKey"`
    ProjectID      uint   `gorm:"not null"`
    TestVersion    string
    TestEnv        string
    TestDate       string
    Executor       string
    CaseNumber     string
    Screen         string
    Function       string
    Precondition   string
    TestSteps      string `gorm:"type:text"`
    ExpectedResult string `gorm:"type:text"`
    TestResult     string // OK, NG, Block, NR
    Remark         string
}
```

#### AutoTestCase (自动化测试用例)

```go
type AutoTestCase struct {
    ID             uint   `gorm:"primaryKey"`
    ProjectID      uint   `gorm:"not null"`
    TestVersion    string
    TestDate       string
    CaseNumber     string
    Screen         string
    Function       string
    Precondition   string
    TestSteps      string `gorm:"type:text"`
    ExpectedResult string `gorm:"type:text"`
    TestResult     string // NR, OK, NG
    Remark         string
}
```

#### APITestCase (接口测试用例)

```go
type APITestCase struct {
    ID          uint   `gorm:"primaryKey"`
    ProjectID   uint   `gorm:"not null"`
    TestVersion string
    TestDate    string
    CaseNumber  string
    Screen      string
    URL         string
    Header      string `gorm:"type:text"`
    Method      string
    Body        string `gorm:"type:text"`
    Response    string `gorm:"type:text"`
    TestResult  string // NR, OK, NG
    Remark      string
}
```

---

## 四、API设计规范

### 4.1 认证相关

```
POST   /api/v1/auth/login          # 用户登录
POST   /api/v1/auth/logout         # 用户登出
POST   /api/v1/auth/refresh        # 刷新Token
```

### 4.2 项目管理

```
GET    /api/v1/projects            # 获取项目列表
POST   /api/v1/projects            # 创建项目
PUT    /api/v1/projects/:id        # 更新项目
DELETE /api/v1/projects/:id        # 删除项目
GET    /api/v1/projects/:id        # 获取项目详情
```

### 4.3 需求管理

```
GET    /api/v1/projects/:id/requirement        # 获取需求文档
PUT    /api/v1/projects/:id/requirement        # 更新需求文档
GET    /api/v1/projects/:id/requirement/export # 下载需求文档
```

### 4.4 测试用例管理

```
# 手工测试用例
GET    /api/v1/projects/:id/manual-cases       # 获取列表
POST   /api/v1/projects/:id/manual-cases       # 创建
PUT    /api/v1/projects/:id/manual-cases/:cid  # 更新
DELETE /api/v1/projects/:id/manual-cases/:cid  # 删除
GET    /api/v1/projects/:id/manual-cases/export # 导出

# 自动化测试用例
GET    /api/v1/projects/:id/auto-cases
POST   /api/v1/projects/:id/auto-cases
PUT    /api/v1/projects/:id/auto-cases/:cid
DELETE /api/v1/projects/:id/auto-cases/:cid
GET    /api/v1/projects/:id/auto-cases/export

# 接口测试用例
GET    /api/v1/projects/:id/api-cases
POST   /api/v1/projects/:id/api-cases
PUT    /api/v1/projects/:id/api-cases/:cid
DELETE /api/v1/projects/:id/api-cases/:cid
GET    /api/v1/projects/:id/api-cases/export
```

### 4.5 统计数据

```
GET    /api/v1/projects/:id/statistics         # 获取项目统计数据
```

### 4.6 用户管理

```
GET    /api/v1/users                   # 获取用户列表
POST   /api/v1/users                   # 创建用户
PUT    /api/v1/users/:id               # 更新用户
DELETE /api/v1/users/:id               # 删除用户
POST   /api/v1/users/:id/reset-password # 重置密码
```

### 4.7 人员分配

```
GET    /api/v1/projects/:id/members    # 获取项目成员
PUT    /api/v1/projects/:id/members    # 更新项目成员
```

### 4.8 个人信息

```
GET    /api/v1/profile                 # 获取个人信息
PUT    /api/v1/profile/nickname        # 更新昵称
PUT    /api/v1/profile/password        # 修改密码
```

---

## 五、前端技术特性

### 5.1 组件库选型

- **UI框架**: Ant Design (antd)
- **路由**: React Router
- **状态管理**: Redux / Context API
- **HTTP客户端**: Axios
- **国际化**: react-i18next

### 5.2 关键组件

#### 5.2.1 可编辑表格 (EditableTable)

- 支持内联编辑
- 自动保存机制
- ID自动排序
- 分页组件集成

#### 5.2.2 Markdown编辑器

- 实时预览
- 工具栏
- 代码高亮

#### 5.2.3 统计图表 (Charts)

- 饼图组件 (基于 ECharts / Recharts)
- 响应式设计

#### 5.2.4 人员穿梭框 (Transfer)

- 双向移动
- 搜索过滤
- 禁用项控制

---

## 六、安全特性

### 6.1 认证安全

- JWT Token有效期控制
- 密码加密存储 (bcrypt)
- HTTPS传输

### 6.2 权限控制

- 基于角色的访问控制 (RBAC)
- API级别权限验证
- 前端路由守卫

### 6.3 数据验证

- 输入参数验证
- SQL注入防护 (GORM参数化查询)
- XSS防护

---

## 七、非功能特性

### 7.1 性能要求

- 页面加载时间 < 2秒
- API响应时间 < 500ms
- 支持1000+并发用户

### 7.2 可用性

- 7×24小时服务
- 数据自动备份
- 错误日志记录

### 7.3 兼容性

- 浏览器支持: Chrome, Firefox, Safari, Edge (最新两个版本)
- 响应式设计: 支持桌面端和平板

### 7.4 可扩展性

- 微服务架构预留
- 多数据库切换支持
- 插件化测试引擎接口

---

## 八、开发规范

### 8.1 代码规范

- Go: 遵循 Effective Go
- React: ESLint + Prettier
- Git提交: Conventional Commits

### 8.2 测试要求

- 单元测试覆盖率 > 70%
- 集成测试覆盖核心流程
- E2E测试覆盖主要用户场景

### 8.3 文档要求

- API文档: Swagger/OpenAPI
- 代码注释: GoDoc标准
- 部署文档: Docker化

---

## 九、部署架构

### 9.1 推荐部署方案

```
[Nginx] → [Go Backend] → [MongoDB/PostgreSQL]
            ↓
        [React SPA]
```

### 9.2 容器化

- Docker镜像构建
- docker-compose编排
- 环境变量配置

### 9.3 CI/CD

- 自动化测试
- 自动化部署
- 版本管理

---

## 十、待扩展特性

### 10.1 短期规划

- [ ] 测试用例导入功能 (Excel, CSV)
- [ ] 操作日志审计
- [ ] 邮件通知功能

### 10.2 中期规划

- [ ] 自动化测试引擎集成
- [ ] 缺陷管理模块
- [ ] 测试报告生成器

### 10.3 长期规划

- [ ] AI辅助用例生成
- [ ] 性能测试模块
- [ ] 移动端支持

---

## 附录

### A. 术语表

| 术语   | 说明                                  |
| ---- | ----------------------------------- |
| JWT  | JSON Web Token，无状态认证令牌              |
| RBAC | Role-Based Access Control，基于角色的访问控制 |
| ORM  | Object-Relational Mapping，对象关系映射    |
| SPA  | Single Page Application，单页应用        |

### B. 参考文档

- Gin框架文档: https://gin-gonic.com/
- GORM文档: https://gorm.io/
- Ant Design: https://ant.design/
- React文档: https://react.dev/

---

**文档结束**
