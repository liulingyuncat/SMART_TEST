# 国际化(i18n)使用指南

## 概述
本项目使用 react-i18next 实现国际化,支持中文(zh)、英文(en)、日语(ja)三种语言。

## 翻译资源文件结构

翻译资源位于 `frontend/src/i18n/index.js`,按模块组织:

```javascript
resources: {
  zh: {
    translation: {
      common: {},      // 通用文本
      login: {},       // 登录页面
      menu: {},        // 菜单项
      project: {},     // 项目相关
      user: {},        // 用户相关
      profile: {},     // 个人信息
      validation: {},  // 表单验证
      message: {}      // 系统消息
    }
  },
  en: { ... },
  ja: { ... }
}
```

## 在组件中使用翻译

### 1. 函数组件中使用

```javascript
import { useTranslation } from 'react-i18next';

const MyComponent = () => {
  const { t } = useTranslation();
  
  return (
    <div>
      <h1>{t('common.title')}</h1>
      <button>{t('common.submit')}</button>
    </div>
  );
};
```

### 2. 带参数的翻译

```javascript
// 翻译资源定义
message: {
  userDeleted: '用户 {{username}} 已删除'
}

// 使用
message.success(t('message.userDeleted', { username: 'admin' }));
// 输出: "用户 admin 已删除"
```

### 3. 在非组件模块中使用

对于 API 客户端、工具函数等非组件模块,直接导入 i18n 实例:

```javascript
import i18n from '../i18n';

// 使用 i18n.t() 而不是 t()
message.error(i18n.t('message.networkError'));
```

## 添加新的翻译键

1. 在 `frontend/src/i18n/index.js` 中找到对应模块
2. 同时在 zh、en、ja 三个语言的相同位置添加翻译键
3. 确保键名和结构完全一致

```javascript
// 示例:添加新的项目相关翻译
resources: {
  zh: {
    translation: {
      project: {
        title: '项目管理',
        createButton: '创建项目',
        newKey: '新增的文本'  // 新增
      }
    }
  },
  en: {
    translation: {
      project: {
        title: 'Project Management',
        createButton: 'Create Project',
        newKey: 'New Text'  // 对应英文
      }
    }
  },
  ja: {
    translation: {
      project: {
        title: 'プロジェクト管理',
        createButton: 'プロジェクト作成',
        newKey: '新しいテキスト'  // 对应日文
      }
    }
  }
}
```

## 语言切换

### 用户端切换

- **登录页面**: 右上角 Select 下拉选择器
- **应用内**: Header 右上角地球图标按钮

### 程序化切换

```javascript
import { useTranslation } from 'react-i18next';

const { i18n } = useTranslation();

// 切换到英文
i18n.changeLanguage('en');

// 切换到日文
i18n.changeLanguage('ja');
```

## 持久化机制

语言偏好自动保存到 localStorage:
- 键名: `user_language`
- 允许值: `zh`, `en`, `ja`
- 默认语言: `zh`
- 无效值自动回退到默认语言并清除

## 调试

开发环境已启用 i18n debug 模式,打开浏览器控制台可看到:
- 翻译键查找日志
- 缺失的翻译键警告(格式: `i18next:: key "xxx" not found`)
- 语言切换事件

## 最佳实践

1. **避免硬编码文本**: 所有用户可见的文本都应使用 `t()` 函数
2. **命名规范**: 翻译键使用小驼峰命名,如 `createButton`, `userDeleted`
3. **模块化**: 按功能模块组织翻译资源,避免单一大对象
4. **参数化**: 对于动态内容,使用参数而非字符串拼接
5. **完整性**: 添加新键时必须同步更新所有语言版本
6. **测试**: 切换语言后检查页面,确保无硬编码文本残留

## Ant Design 组件国际化

Ant Design 组件(DatePicker、Pagination、Modal 等)的内置文本通过 `ConfigProvider` 自动切换:

```javascript
// App.js 中已配置
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import enUS from 'antd/locale/en_US';
import jaJP from 'antd/locale/ja_JP';

<ConfigProvider locale={locale}>
  {/* 应用内容 */}
</ConfigProvider>
```

无需手动翻译 Ant Design 组件的内置文本。

## 常见问题

### Q: 添加新翻译后页面没有更新?
A: 刷新浏览器,i18n 资源在应用启动时加载。

### Q: 如何检查翻译是否完整?
A: 1. 启用 debug 模式查看控制台警告 2. 切换语言遍历所有页面检查

### Q: 翻译文本过长导致布局错乱?
A: 使用 CSS 处理文本溢出,或调整翻译内容使其更简洁。

### Q: 如何支持更多语言?
A: 1. 在 `i18n/index.js` 的 resources 中添加新语言对象 2. 在 SUPPORTED_LANGUAGES 中添加语言代码 3. 在 LanguageSwitch 的 languageOptions 中添加选项 4. 在 App.js 的 localeMap 中添加映射

## 参考资源

- [react-i18next 官方文档](https://react.i18next.com/)
- [Ant Design 国际化](https://ant.design/docs/react/i18n)
- [项目翻译资源文件](./frontend/src/i18n/index.js)
