# MCP Tools 完整列表（39个工具）

本文档列出当前已注册的所有MCP工具及其详细说明。所有工具均按照MCP Protocol实现，可通过标准的JSON-RPC调用。

## 目录

1. [用户与项目信息](#用户与项目信息) - 2个工具
2. [原始文档](#原始文档) - 2个工具
3. [需求条目](#需求条目) - 4个工具
4. [测试观点](#测试观点) - 4个工具
5. [用例集与手工用例](#用例集与手工用例) - 6个工具
6. [Web自动化用例](#web自动化用例) - 6个工具
7. [API接口用例](#api接口用例) - 6个工具
8. [用例评审](#用例评审) - 1个工具
9. [执行任务](#执行任务) - 3个工具
10. [缺陷管理](#缺陷管理) - 2个工具
11. [AI报告](#ai报告) - 2个工具

---

## 用户与项目信息

### get_current_user_info

获取当前Token对应的用户信息，包括user_id、username、nickname和role

**参数**：无

**返回示例**：

```json
{
  "user_id": 1,
  "username": "admin",
  "nickname": "管理员",
  "role": "system_admin"
}
```

### get_current_project_name

获取当前用户选择的当前项目名称和ID

**参数**：无

**返回示例**：

```json
{
  "project_id": 1,
  "project_name": "Test Project"
}
```

---

## 原始文档

### list_raw_documents

获取项目中的所有原始转换文档列表

**参数**：

- `project_id` (integer, required): 项目ID

**返回**：原始文档列表（包含文档名称、类型、转换时间等）

### get_converted_document

获取单个原始转换文档的详细内容

**参数**：

- `project_id` (integer, required): 项目ID
- `document_id` (integer, required): 文档ID

**返回**：完整的文档内容

---

## 需求条目

### list_requirement_items

获取项目中的AI需求文档列表

**参数**：

- `project_id` (integer, required): 项目ID

**返回**：需求文档列表（元数据）

### get_requirement_item

获取单个AI需求文档的详细内容

**参数**：

- `project_id` (integer, required): 项目ID
- `requirement_id` (integer, required): 需求文档ID

**返回**：完整的需求文档内容

### create_requirement_item

创建AI需求文档

**参数**：

- `project_id` (integer, required): 项目ID
- `name` (string, required): 需求名称
- `content` (string, required): 需求内容
- `parent_id` (integer, optional): 父需求ID（支持层级结构）

**返回**：新创建的需求文档ID和详细信息

### update_requirement_item

更新AI需求文档

**参数**：

- `project_id` (integer, required): 项目ID
- `requirement_id` (integer, required): 需求文档ID
- `name` (string, optional): 需求名称
- `content` (string, optional): 需求内容

**返回**：更新后的需求文档信息

---

## 测试观点

### list_viewpoint_items

获取项目中的AI观点文档列表

**参数**：

- `project_id` (integer, required): 项目ID

**返回**：观点文档列表（元数据）

### get_viewpoint_item

获取单个AI观点文档的详细内容

**参数**：

- `project_id` (integer, required): 项目ID
- `viewpoint_id` (integer, required): 观点文档ID

**返回**：完整的观点文档内容

### create_viewpoint_item

创建AI观点文档

**参数**：

- `project_id` (integer, required): 项目ID
- `name` (string, required): 观点名称
- `content` (string, required): 观点内容（Markdown格式）
- `requirement_id` (integer, optional): 关联的需求ID

**返回**：新创建的观点文档ID和详细信息

### update_viewpoint_item

更新AI观点文档

**参数**：

- `project_id` (integer, required): 项目ID
- `viewpoint_id` (integer, required): 观点文档ID
- `name` (string, optional): 观点名称
- `content` (string, optional): 观点内容

**返回**：更新后的观点文档信息

---

## 用例集与手工用例

### list_manual_groups

获取项目的手工测试用例集列表

**参数**：

- `project_id` (integer, required): 项目ID

**返回**：用例集列表

### list_manual_cases

获取用例集中的手工测试用例列表

**参数**：

- `project_id` (integer, required): 项目ID
- `case_group_id` (integer, required): 用例集ID
- `all_fields` (boolean, optional): 是否返回所有字段包括CN、JP、EN (默认false)

**返回**：用例列表

### create_case_group

创建测试用例集

**参数**：

- `project_id` (integer, required): 项目ID
- `name` (string, required): 用例集名称
- `type` (string, optional): 用例类型（默认: overall）
- `description` (string, optional): 用例集描述

**返回**：新创建的用例集ID和详细信息

### create_manual_cases

创建手工测试用例（支持批量创建）

**参数**：

- `project_id` (integer, required): 项目ID
- `case_group_id` (integer, required): 用例集ID
- `cases` (array, required): 用例数据数组，每个元素包含各语言字段
- `continue_on_error` (boolean, optional): 失败是否继续处理（默认: true）

**返回**：创建结果列表，包含每个用例的ID和状态

### update_manual_case

更新单个手工测试用例

**参数**：

- `project_id` (integer, required): 项目ID
- `case_id` (integer, required): 用例ID
- `data` (object, required): 要更新的用例数据

**返回**：更新后的用例信息

### update_manual_cases

批量更新手工测试用例

**参数**：

- `project_id` (integer, required): 项目ID
- `case_group_id` (integer, optional): 用例集ID（推荐提供）
- `cases` (array, required): 用例数据数组，每个元素必须包含id字段（整数型）
- `continue_on_error` (boolean, optional): 失败是否继续处理（默认: true）

**返回**：更新结果列表

---

## Web自动化用例

### list_web_groups

获取Web自动化用例集列表

**参数**：

- `project_id` (integer, required): 项目ID

**返回**：Web用例集列表

### get_web_group_metadata

获取Web用例集的元数据信息

**参数**：

- `project_id` (integer, required): 项目ID
- `group_id` (integer, required): 用例集ID

**返回**：用例集元数据（名称、描述、创建时间等）

### list_web_cases

获取AIWeb测试用例列表

**参数**：

- `project_id` (integer, required): 项目ID
- `group_id` (integer, required): 用例集ID

**返回**：Web用例列表

### create_web_group

创建Web自动化用例集

**参数**：

- `project_id` (integer, required): 项目ID
- `name` (string, required): 用例集名称
- `description` (string, optional): 用例集描述

**返回**：新创建的用例集ID和详细信息

### create_web_cases

创建AIWeb测试用例（支持批量创建）

**参数**：

- `project_id` (integer, required): 项目ID
- `group_id` (integer, required): 用例集ID
- `cases` (array, required): 用例数据数组（包含页面、操作等Web特定字段）
- `continue_on_error` (boolean, optional): 失败是否继续处理

**返回**：创建结果列表

### update_web_cases

批量更新AIWeb测试用例

**参数**：

- `project_id` (integer, required): 项目ID
- `cases` (array, required): 用例数据数组，每个元素必须包含id字段
- `continue_on_error` (boolean, optional): 失败是否继续处理

**返回**：更新结果列表

---

## API接口用例

### list_api_groups

获取API接口用例集列表

**参数**：

- `project_id` (integer, required): 项目ID

**返回**：API用例集列表

### get_api_group_metadata

获取API用例集的元数据信息

**参数**：

- `project_id` (integer, required): 项目ID
- `group_id` (integer, required): 用例集ID

**返回**：用例集元数据

### list_api_cases

获取AI接口测试用例列表

**参数**：

- `project_id` (integer, required): 项目ID
- `group_id` (integer, required): 用例集ID

**返回**：API用例列表

### create_api_group

创建API接口用例集

**参数**：

- `project_id` (integer, required): 项目ID
- `name` (string, required): 用例集名称
- `description` (string, optional): 用例集描述

**返回**：新创建的用例集ID和详细信息

### create_api_case

创建AI接口测试用例

**参数**：

- `project_id` (integer, required): 项目ID
- `group_id` (integer, required): 用例集ID
- `case_data` (object, required): 用例数据（包含请求URL、方法、参数、期望响应等）

**返回**：新创建的用例ID和详细信息

### update_api_case

更新AI接口测试用例

**参数**：

- `project_id` (integer, required): 项目ID
- `case_id` (integer, required): 用例ID
- `case_data` (object, required): 更新的用例数据

**返回**：更新后的用例信息

---

## 用例评审

### create_review_item

创建用例评审条目

**参数**：

- `project_id` (integer, required): 项目ID
- `name` (string, required): 审阅条目名称
- `content` (string, optional): 评审内容（Markdown格式）

**返回**：新创建的评审条目ID和详细信息

---

## 执行任务

### list_execution_tasks

获取项目的执行任务列表

**参数**：

- `project_id` (integer, required): 项目ID

**返回**：执行任务列表（包含任务ID、名称、创建时间等）

### get_execution_task_cases

获取执行任务关联的用例列表

**参数**：

- `project_id` (integer, required): 项目ID
- `task_id` (string, required): 执行任务ID

**返回**：用例列表及其执行结果

### update_execution_case_result

更新执行任务中单个用例的执行结果

**参数**：

- `project_id` (integer, required): 项目ID
- `case_id` (string, required): 用例ID
- `result` (string, required): 执行结果（BLOCK、OK、NG、NR）
- `comment` (string, optional): 执行备注

**返回**：更新后的结果信息

---

## 缺陷管理

### list_defects

获取项目的缺陷列表

**参数**：

- `project_id` (integer, required): 项目ID
- `page` (integer, optional): 分页页码（默认1）
- `page_size` (integer, optional): 每页数量（默认50）

**返回**：缺陷列表和总数

### update_defect

更新缺陷信息

**参数**：

- `project_id` (integer, required): 项目ID
- `defect_id` (string, required): 缺陷ID
- `status` (string, optional): 缺陷状态（Open、Fixed、Closed等）
- `comment` (string, optional): 更新备注

**返回**：更新后的缺陷信息

---

## AI报告

### create_ai_report

创建AI测试报告

**参数**：

- `project_id` (integer, required): 项目ID
- `title` (string, required): 报告标题
- `content` (string, required): 报告内容（Markdown格式）

**返回**：新创建的报告ID和详细信息

### update_ai_report

更新AI测试报告

**参数**：

- `project_id` (integer, required): 项目ID
- `report_id` (string, optional): 报告ID（字符串格式，如 report_xxx）
- `report_name` (string, optional): 报告名称（用于查找报告，支持精确匹配或模糊匹配）
- `content` (string, optional): 报告内容（Markdown格式）
- `new_name` (string, optional): 报告新名称（用于重命名）

**返回**：更新后的报告信息

---

## 工具统计

- **总工具数**：39个
- **分类总数**：11个
- **最多工具分类**：Web自动化用例、API接口用例、用例集与手工用例（各6个）
- **最少工具分类**：用例评审（1个）

## 使用说明

### 调用示例

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_current_user_info",
    "arguments": {}
  }
}
```

### 返回格式

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "..."
      }
    ]
  }
}
```

## 更新历史

- **2025-12-27**：首次完整文档化，共39个工具，按功能分为11个分类
