# MCP 系统提示词书写规则

> 简洁参考指南，用于编写 `.prompt.md` 文件

---

## 文件结构

```markdown
---
name: 提示词名称
description: 简短描述
version: 1.0
arguments:
  - name: 参数名
    description: 参数说明
    required: true/false
---

[Markdown正文内容]
```

---

## 参数定义 (arguments)

### 必填参数
```yaml
arguments:
  - name: document_name
    description: 文档名称 (Document Name / ドキュメント名)
    required: true
```

### 可选参数
```yaml
arguments:
  - name: output_format
    description: 输出格式 (可选，默认markdown)
    required: false
```

### 多参数示例
```yaml
arguments:
  - name: document_name
    description: 文档名称
    required: true
  - name: language
    description: 目标语言 (可选)
    required: false
  - name: max_items
    description: 最大条目数 (可选)
    required: false
```

---

## 参数引用

在 Markdown 正文中使用 `{{参数名}}` 引用参数：

```markdown
请处理文档：**{{document_name}}**
输出语言：{{language}}
```

---

## 字段说明

| 字段 | 必填 | 类型 | 说明 |
|------|------|------|------|
| `name` | ✅ | string | 提示词唯一标识 |
| `description` | ❌ | string | 提示词描述 |
| `version` | ❌ | string | 版本号 |
| `arguments` | ❌ | array | 参数列表 |
| `arguments[].name` | ✅ | string | 参数名 |
| `arguments[].description` | ❌ | string | 参数描述 |
| `arguments[].required` | ❌ | bool | 是否必填，默认 false |

---

## 三语描述格式（简洁版）

```yaml
description: 中文说明 (English / 日本語)
```

示例：
```yaml
description: 转换文档名 (Document Name / ドキュメント名)
description: 用例数量 (Case Count / ケース数)
description: 输出格式 (Output Format / 出力形式)
```

---

## 命名规范

- **文件名**: `S{序号}_{功能名}.prompt.md`
- **参数名**: 使用 `snake_case`，如 `document_name`、`output_format`
- **版本号**: 使用语义化版本，如 `1.0`、`1.1`、`2.0`
