# Web 智能测试平台 - 用户登录功能手动测试指南

## 前置条件

1. 后端服务已启动(端口 8080)
2. 前端服务已启动(端口 3000)
3. 数据库初始化完成,管理员账号已创建

## 测试账号

| 用户名 | 密码     | 角色  |
|--------|----------|-------|
| admin  | admin123 | admin |
| root   | root123  | admin |

## 测试用例

### TC-01: 正常登录流程

**步骤:**
1. 访问 http://localhost:3000/login
2. 输入用户名: `admin`
3. 输入密码: `admin123`
4. 点击"登录"按钮

**预期结果:**
- ✅ 显示"登录成功"提示消息
- ✅ 自动跳转到首页 `/`
- ✅ localStorage 中存储 `auth_token`
- ✅ Redux 状态 `isAuthenticated` 为 true

**验证方法:**
```javascript
// 在浏览器控制台执行
console.log(localStorage.getItem('auth_token'));
```

---

### TC-02: 用户名密码错误

**步骤:**
1. 访问登录页面
2. 输入用户名: `admin`
3. 输入密码: `wrongpassword`
4. 点击"登录"按钮

**预期结果:**
- ✅ 显示"用户名或密码错误"错误提示
- ✅ 不跳转页面,保持在登录页
- ✅ localStorage 中无 token

---

### TC-03: 必填字段验证

**步骤:**
1. 访问登录页面
2. 不输入任何内容
3. 直接点击"登录"按钮

**预期结果:**
- ✅ 用户名输入框下方显示"请输入用户名"
- ✅ 密码输入框下方显示"请输入密码"
- ✅ 不发送 API 请求

---

### TC-04: 用户名长度验证

**步骤:**
1. 输入用户名: `ab` (少于3个字符)
2. 输入密码: `admin123`
3. 点击"登录"

**预期结果:**
- ✅ 显示"用户名长度为 3-50 个字符"

---

### TC-05: 密码长度验证

**步骤:**
1. 输入用户名: `admin`
2. 输入密码: `12345` (少于6个字符)
3. 点击"登录"

**预期结果:**
- ✅ 显示"密码长度为 6-50 个字符"

---

### TC-06: 语言切换

**步骤:**
1. 访问登录页面(默认中文)
2. 点击右上角语言选择器
3. 选择 "English"

**预期结果:**
- ✅ 页面标题变为 "Login"
- ✅ 所有中文文本切换为英文
- ✅ 表单验证消息也显示为英文

---

### TC-07: Token 持久化

**步骤:**
1. 成功登录
2. 刷新页面
3. 访问受保护路由

**预期结果:**
- ✅ 刷新后仍然保持登录状态
- ✅ 可以正常访问受保护路由

---

### TC-08: 登出功能

**步骤:**
1. 成功登录后
2. 调用登出操作

**验证方法:**
```javascript
// 在浏览器控制台
import { logout } from './store/authSlice';
// 或直接清除
localStorage.removeItem('auth_token');
window.location.reload();
```

**预期结果:**
- ✅ localStorage 中 token 被清除
- ✅ 访问受保护路由时重定向到登录页

---

### TC-09: 受保护路由访问控制

**步骤:**
1. 未登录状态下访问 http://localhost:3000/

**预期结果:**
- ✅ 自动重定向到 `/login`

---

### TC-10: API 请求验证(后端测试)

**使用 curl 测试:**
```bash
# 正确登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 预期响应:
# {"code":0,"message":"success","data":{"token":"eyJhbGc..."}}

# 错误密码
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrong"}'

# 预期响应:
# {"code":401,"message":"invalid username or password"}
```

---

### TC-11: Token 认证测试

**步骤:**
```bash
# 先获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.token')

# 使用 token 访问受保护接口
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $TOKEN"

# 预期响应:
# {"username":"admin","message":"authenticated user profile"}
```

---

## 回归测试检查清单

- [ ] 正常登录成功
- [ ] 错误密码被拒绝
- [ ] 表单验证生效
- [ ] 语言切换正常
- [ ] Token 持久化
- [ ] 登出清除状态
- [ ] 路由守卫生效
- [ ] API 返回正确状态码
- [ ] 网络错误处理
- [ ] 响应式布局适配

## 浏览器兼容性测试

- [ ] Chrome 90+
- [ ] Firefox 88+
- [ ] Safari 14+
- [ ] Edge 90+

## 性能测试

- [ ] 登录请求响应时间 < 200ms
- [ ] 首屏渲染时间 < 2s
- [ ] 无内存泄漏

## 安全测试

- [ ] 密码输入框类型为 password(隐藏字符)
- [ ] Token 使用 HTTPS 传输(生产环境)
- [ ] XSS 防护(React 默认)
- [ ] CSRF 防护(Token 验证)
