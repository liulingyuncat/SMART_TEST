# T31 行内编辑功能实现总结

## 1. 核心逻辑梳理

### 1.1 行编辑三大核心函数

```javascript
// 1. 开始编辑 - 设置表单初始值并进入编辑状态
const startEdit = useCallback((record) => {
  // 步骤1: 规范化数据（null/undefined → 空字符串）
  const getDisplayValue = (value) => {
    if (value === null || value === undefined) return '';
    return String(value);
  };
  
  // 步骤2: 根据用例类型和当前语言构造表单值
  const formValues = buildFormValues(record, caseType, language);
  
  // 步骤3: 设置表单值
  form.setFieldsValue(formValues);
  
  // 步骤4: 设置编辑键（标记当前编辑行）
  setEditingKey(record.case_id);
}, [form, caseType, language]);

// 2. 保存编辑 - 收集变更并更新数据
const saveEdit = useCallback(async (recordParam) => {
  try {
    // 步骤1: 验证表单
    const row = await form.validateFields();
    
    // 步骤2: 获取最新记录（从cases数组）
    const currentRecord = cases.find(c => c.case_id === recordParam.case_id);
    
    // 步骤3: 收集有变化的字段
    const updates = collectChangedFields(row, currentRecord);
    
    // 步骤4: 如果有变更，调用API更新
    if (Object.keys(updates).length > 0) {
      await updateCaseAPI(currentRecord.case_id, updates);
      
      // 步骤5: 更新本地数据（不重新加载）
      setCases(prevCases => 
        prevCases.map(c => 
          c.case_id === currentRecord.case_id 
            ? { ...c, ...updates }
            : c
        )
      );
      
      message.success('保存成功');
    }
    
    // 步骤6: 退出编辑状态
    setEditingKey('');
    
  } catch (error) {
    console.error('Failed to save edit:', error);
    message.error('保存失败');
    setEditingKey('');  // 失败也要退出编辑
  }
}, [form, projectId, pagination, fetchCases, caseType, language, cases]);

// 3. 取消编辑 - 恢复原始状态
const cancelEdit = useCallback(() => {
  setEditingKey('');      // 清空编辑键
  form.resetFields();     // 重置表单值
}, [form]);
```

### 1.2 字段映射逻辑

#### AI用例（单语言）

```javascript
// AI用例直接使用无后缀字段名
const formValues = {
  case_number: record.case_number,
  major_function: record.major_function,
  middle_function: record.middle_function,
  minor_function: record.minor_function,
  precondition: record.precondition,
  test_steps: record.test_steps,
  expected_result: record.expected_result,
  remark: record.remark
};
```

#### 整体/变更/受入用例（多语言）

```javascript
// 根据当前语言添加后缀
const langSuffix = language === '中文' ? '_cn' 
                 : language === 'English' ? '_en' 
                 : '_jp';

const formValues = {
  case_number: record.case_number,
  [`major_function${langSuffix}`]: record[`major_function${langSuffix}`],
  [`middle_function${langSuffix}`]: record[`middle_function${langSuffix}`],
  [`minor_function${langSuffix}`]: record[`minor_function${langSuffix}`],
  [`precondition${langSuffix}`]: record[`precondition${langSuffix}`],
  [`test_steps${langSuffix}`]: record[`test_steps${langSuffix}`],
  [`expected_result${langSuffix}`]: record[`expected_result${langSuffix}`],
  test_result: record.test_result,
  remark: record.remark
};
```

#### Role类型用例（自动化测试）

```javascript
// 使用screen和function字段，添加语言后缀
const langSuffix = language === '中文' ? '_cn' 
                 : language === 'English' ? '_en' 
                 : '_jp';

const formValues = {
  case_num: record.case_num,  // 注意：Role用例用case_num而非case_number
  [`screen${langSuffix}`]: record[`screen${langSuffix}`],
  [`function${langSuffix}`]: record[`function${langSuffix}`],
  [`precondition${langSuffix}`]: record[`precondition${langSuffix}`],
  [`test_steps${langSuffix}`]: record[`test_steps${langSuffix}`],
  [`expected_result${langSuffix}`]: record[`expected_result${langSuffix}`],
  test_result: record.test_result,
  remark: record.remark
};
```

### 1.3 变更检测算法

```javascript
// 规范化函数：统一处理空值和空格
const normalizeValue = (value) => {
  if (value === null || value === undefined) return '';
  return String(value).trim();
};

// 收集变更字段
const updates = {};
Object.keys(row).forEach(key => {
  const formValue = normalizeValue(row[key]);
  const recordValue = normalizeValue(currentRecord[key]);
  
  // 只有值不同时才添加到更新对象
  if (formValue !== recordValue) {
    // 保留原始值类型，但将null/undefined转为空字符串
    updates[key] = row[key] === undefined || row[key] === null 
      ? '' 
      : String(row[key]);
  }
});
```

**关键点**：
1. 使用trim()去除首尾空格进行比较
2. null/undefined统一视为空字符串
3. 只更新有变化的字段（减少网络传输）
4. 保存时保留表单的原始值

### 1.4 状态控制逻辑

#### 编辑状态标记

```javascript
// editingKey: 当前编辑行的case_id
// - 空字符串: 没有行在编辑
// - UUID字符串: 对应行在编辑状态

setEditingKey(record.case_id);  // 进入编辑
setEditingKey('');               // 退出编辑

// 判断某行是否在编辑
const isEditing = editingKey === record.case_id;
```

#### 插入状态标记

```javascript
// hasEditChanges: 是否有未保存的插入操作
// - false: 没有待保存的插入
// - true: 有插入操作待保存

setHasEditChanges(true);   // Above/Below插入时设置
setHasEditChanges(false);  // 保存后或重新加载后清除
```

#### 按钮禁用逻辑

```javascript
// Edit按钮
disabled={editingKey !== '' || hasEditChanges}
// 原因：编辑中或有插入操作时不能开始新的编辑

// Delete按钮
disabled={editingKey === record.case_id || hasEditChanges}
// 原因：当前行编辑中或有插入操作时不能删除

// Above/Below按钮
disabled={editingKey !== '' || loading}
// 原因：编辑中或加载中不能插入
// 注意：不受hasEditChanges限制（可多次插入）

// Save按钮（保存所有插入）
disabled={pendingInserts.length === 0}
// 原因：没有待保存的插入时禁用
```

## 2. 特殊场景处理

### 2.1 TestResult字段即时保存

```javascript
// TestResult在编辑状态下显示为Select
// 用户选择后立即保存，不需要点Save按钮

<Select 
  defaultValue={record.test_result || 'NR'}
  onChange={(value) => {
    if (value !== record.test_result) {
      // 立即调用API更新
      updateCaseAPI(record.case_id, { test_result: value })
        .then(() => {
          message.success('保存成功');
          
          // 更新本地数据，但不调用fetchCases
          // 这样可以保持编辑状态和hasEditChanges状态
          setCases(prevCases =>
            prevCases.map(c =>
              c.case_id === record.case_id
                ? { ...c, test_result: value }
                : c
            )
          );
        })
        .catch(error => {
          message.error('保存失败');
        });
    }
  }}
>
  <Select.Option value="NR">NR</Select.Option>
  <Select.Option value="OK">OK</Select.Option>
  <Select.Option value="NG">NG</Select.Option>
  <Select.Option value="Block">Block</Select.Option>
</Select>
```

**关键点**：
1. 使用`defaultValue`而非`value`（避免受控组件问题）
2. 在`onChange`中立即保存
3. 只在值改变时调用API
4. 更新本地数据但不调用fetchCases
5. 保持编辑状态和hasEditChanges状态

### 2.2 多语言字段编辑

#### 中文界面的特殊处理

```javascript
// 在中文界面下，多语言字段显示为蓝色链接
// 点击打开对话框，可同时编辑三种语言

const isChinese = language === '中文';

// 渲染函数
render: (text, record) => {
  const isEditing = editingKey === record.case_id;
  const fieldValue = record[`major_function${langSuffix}`];
  
  if (isChinese) {
    // 中文界面：显示蓝色链接
    return (
      <div
        style={{ cursor: 'pointer', color: '#1890ff' }}
        onClick={() => !isEditing && openMultiLangModal(record, 'major_function')}
      >
        {fieldValue || '-'}
      </div>
    );
  }
  
  // 非中文界面：正常显示
  return fieldValue || '-';
}

// onCell配置
onCell: (record) => ({
  editing: editingKey === record.case_id && !isChinese,  // 中文界面不显示编辑框
  dataIndex: `major_function${langSuffix}`,
  title: `Major${langSuffix}`,
  record,
  inputType: 'text',
})
```

**中文界面逻辑**：
1. 字段显示为蓝色链接（非编辑状态）
2. 点击打开MultiLangEditModal
3. 对话框中同时编辑CN/EN/JP三种语言
4. 保存时更新所有三种语言

**英文/日文界面逻辑**：
1. 进入编辑状态时显示输入框
2. 只编辑当前语言的字段
3. 保存时只更新当前语言

#### MultiLangEditModal实现

```javascript
const openMultiLangModal = useCallback((record, fieldName) => {
  setMultiLangModalConfig({
    visible: true,
    record,
    fieldName,
    fieldTitle: getFieldTitle(fieldName),
    initialValues: {
      cn: record[`${fieldName}_cn`] || '',
      en: record[`${fieldName}_en`] || '',
      jp: record[`${fieldName}_jp`] || ''
    }
  });
}, []);

// 保存多语言编辑
const handleMultiLangSave = async (values) => {
  const updates = {
    [`${fieldName}_cn`]: values.cn,
    [`${fieldName}_en`]: values.en,
    [`${fieldName}_jp`]: values.jp
  };
  
  await updateCaseAPI(record.case_id, updates);
  
  // 更新本地数据
  setCases(prevCases =>
    prevCases.map(c =>
      c.case_id === record.case_id
        ? { ...c, ...updates }
        : c
    )
  );
  
  setMultiLangModalConfig({ visible: false });
  message.success('保存成功');
};
```

### 2.3 默认空行创建

#### 创建函数实现

```javascript
const createDefaultEmptyRow = useCallback(async () => {
  try {
    // 1. 判断用例类型
    const isRoleType = caseType && caseType.startsWith('role');
    
    // 2. 选择对应的创建API
    const createAPI = isRoleType ? createAutoCase : createCase;
    
    // 3. 构造请求参数（注意字段名格式）
    const createData = {
      case_type: caseType,  // 后端期望下划线格式
    };
    
    // 手工用例需要language参数
    if (!isRoleType) {
      createData.language = language;
    }
    
    console.log('[createDefaultEmptyRow] Creating with data:', createData);
    
    // 4. 调用API创建记录
    const newCase = await createAPI(projectId, createData);
    console.log('[createDefaultEmptyRow] Created case:', newCase);
    
    return newCase;
  } catch (error) {
    console.error('[createDefaultEmptyRow] Failed:', error);
    console.error('[createDefaultEmptyRow] Error response:', error?.response?.data);
    throw error;
  }
}, [projectId, caseType, language]);
```

#### 调用场景

**场景1：空表格加载**

```javascript
// 在fetchCases函数中
let casesData = data.cases || [];
if (casesData.length === 0) {
  try {
    const newCase = await createDefaultEmptyRow();
    casesData = [newCase];
  } catch (error) {
    console.error('Failed to create default empty row:', error);
    casesData = [];  // 创建失败则显示空表格
  }
}
```

**场景2：清空AI用例**

```javascript
const handleClearAICases = async () => {
  const confirmed = window.confirm('确定要清空所有AI用例吗?');
  if (confirmed) {
    try {
      await clearAICases(projectId);
      message.success('清空成功');
      
      // 清空后创建默认空行
      try {
        const newCase = await createDefaultEmptyRow();
        setCases([newCase]);
        setPagination(prev => ({ ...prev, current: 1, total: 1 }));
      } catch (error) {
        setCases([]);
        setPagination(prev => ({ ...prev, current: 1, total: 0 }));
      }
      
      if (onRefreshMetadata) {
        onRefreshMetadata();
      }
    } catch (error) {
      message.error('清空失败');
    }
  }
};
```

**场景3：删除到空表格**

```javascript
const handleDelete = useCallback(async (record) => {
  const confirmed = window.confirm(`确定要删除用例吗?`);
  if (confirmed) {
    try {
      await deleteAPI(projectId, record.case_id);
      message.success('删除成功');
      
      // 检查是否删除后为空
      if (cases.length === 1) {
        try {
          const newCase = await createDefaultEmptyRow();
          setCases([newCase]);
          setPagination(prev => ({ ...prev, current: 1, total: 1 }));
        } catch (error) {
          setCases([]);
          setPagination(prev => ({ ...prev, current: 1, total: 0 }));
        }
      } else {
        await fetchCases(pagination.current);
      }
    } catch (error) {
      message.error('删除失败');
    }
  }
}, [projectId, cases.length, createDefaultEmptyRow]);
```

**场景4：批量删除到空表格**

```javascript
const handleBatchDelete = useCallback(async () => {
  const confirmed = window.confirm(`确认删除选中的${selectedRowKeys.length}条用例吗?`);
  if (confirmed) {
    try {
      await batchDeleteAPI(projectId, { caseIds: selectedRowKeys });
      message.success(`成功删除${selectedRowKeys.length}条用例`);
      setSelectedRowKeys([]);
      
      // 检查是否删除后为空
      if (cases.length === selectedRowKeys.length) {
        try {
          const newCase = await createDefaultEmptyRow();
          setCases([newCase]);
          setPagination(prev => ({ ...prev, current: 1, total: 1 }));
        } catch (error) {
          setCases([]);
          setPagination(prev => ({ ...prev, current: 1, total: 0 }));
        }
      } else {
        await fetchCases(pagination.current);
      }
    } catch (error) {
      message.error('批量删除失败');
    }
  }
}, [selectedRowKeys, cases.length, createDefaultEmptyRow]);
```

## 3. 性能优化策略

### 3.1 本地数据更新vs重新加载

**旧方案（性能较差）**：
```javascript
// 保存后重新加载整页数据
await updateCaseAPI(case_id, updates);
await fetchCases(pagination.current);  // ❌ 触发整页刷新
```

**新方案（性能优化）**：
```javascript
// 保存后只更新本地数据
await updateCaseAPI(case_id, updates);

// ✅ 使用函数式setState更新本地数据
setCases(prevCases => 
  prevCases.map(c => 
    c.case_id === currentRecord.case_id 
      ? { ...c, ...updates }  // 合并更新字段
      : c                      // 其他行保持不变
  )
);
```

**优势**：
1. 不触发网络请求，响应更快
2. 保留Above/Below插入的新行（这些行可能不在当前页范围内）
3. 保持hasEditChanges状态（避免触发useEffect重新加载）
4. 避免编辑状态被打断
5. 更好的用户体验

### 3.2 useCallback优化

```javascript
// 所有回调函数使用useCallback包裹
const startEdit = useCallback((record) => {
  // ...
}, [form, caseType, language]);  // 明确依赖项

const saveEdit = useCallback(async (recordParam) => {
  // ...
}, [form, projectId, pagination, fetchCases, caseType, language, cases]);

const cancelEdit = useCallback(() => {
  // ...
}, [form]);
```

**优势**：
1. 避免不必要的函数重新创建
2. 减少子组件的重新渲染
3. 提高整体性能

### 3.3 函数式setState

```javascript
// ✅ 使用函数式setState
setCases(prevCases => {
  return prevCases.map(c => 
    c.case_id === targetId 
      ? { ...c, ...updates }
      : c
  );
});

// ❌ 避免直接依赖外部state
setCases(cases.map(c => 
  c.case_id === targetId 
    ? { ...c, ...updates }
    : c
));
```

**优势**：
1. 确保使用最新状态值
2. 避免闭包陷阱
3. 更可靠的状态更新

## 4. 关键技术点总结

### 4.1 状态管理

```javascript
// 核心状态
const [editingKey, setEditingKey] = useState('');           // 编辑键
const [hasEditChanges, setHasEditChanges] = useState(false); // 插入标记
const [cases, setCases] = useState([]);                      // 用例数据
const [form] = Form.useForm();                               // 表单实例

// 状态优先级
// 1. editingKey - 控制编辑行
// 2. hasEditChanges - 控制插入操作
// 3. loading - 控制加载状态
```

### 4.2 字段映射规则

```javascript
// 1. AI用例：单语言字段（无后缀）
case_number, major_function, middle_function, ...

// 2. 整体/变更/受入用例：多语言字段（_cn/_en/_jp后缀）
case_number, major_function_cn, major_function_en, major_function_jp, ...

// 3. Role用例：多语言字段（_cn/_en/_jp后缀）
case_num, screen_cn, screen_en, screen_jp, ...

// 语言后缀计算
const langSuffix = language === '中文' ? '_cn' 
                 : language === 'English' ? '_en' 
                 : '_jp';
```

### 4.3 API参数格式

```javascript
// 创建用例请求（注意字段名格式）
{
  case_type: caseType,     // ✅ 下划线格式（后端期望）
  language: language       // 手工用例必填
}

// 更新用例请求
{
  case_number: "TC001",           // 驼峰格式（前端字段名）
  major_function_cn: "用户管理"   // 下划线格式（多语言字段）
}
```

### 4.4 错误处理原则

```javascript
// 1. 保存失败后仍退出编辑状态
try {
  await saveEdit(record);
} catch (error) {
  message.error('保存失败');
  setEditingKey('');  // ✅ 退出编辑
}

// 2. 创建默认行失败时降级处理
try {
  const newCase = await createDefaultEmptyRow();
  setCases([newCase]);
} catch (error) {
  setCases([]);  // ✅ 显示空表格
}

// 3. 所有错误都要有用户提示
catch (error) {
  console.error('Operation failed:', error);
  message.error('操作失败');  // ✅ 用户友好的提示
}
```

## 5. 常见问题与解决方案

### 5.1 保存后数据不更新

**问题**：点击Save后数据没有更新

**原因**：调用了fetchCases但hasEditChanges=true时被阻止

**解决**：使用本地数据更新
```javascript
setCases(prevCases => prevCases.map(...));
```

### 5.2 编辑框不显示输入

**问题**：点击Edit后输入框不显示

**原因**：hasEditChanges=true时EditableCell被阻止

**解决**：检查hasEditChanges状态，保存后清除
```javascript
editing && !hasEditChanges  // 两个条件都满足才显示
```

### 5.3 默认空行不能编辑

**问题**：空表格的默认行点击Edit后报404

**原因**：使用了临时ID而非真实UUID

**解决**：调用createDefaultEmptyRow创建真实记录
```javascript
const newCase = await createDefaultEmptyRow();  // 获取真实UUID
```

### 5.4 切换语言后编辑状态丢失

**问题**：切换语言后正在编辑的行退出编辑

**原因**：useEffect依赖language触发了fetchCases

**解决**：检查hasEditChanges
```javascript
useEffect(() => {
  if (!hasEditChanges) {
    fetchCases(pagination.current);
  }
}, [language]);
```

## 6. 最佳实践

1. **使用useCallback优化回调函数**
2. **使用函数式setState更新状态**
3. **保存后本地更新而非重新加载**
4. **统一处理空值（null/undefined → ''）**
5. **明确的错误处理和用户提示**
6. **正确的状态优先级控制**
7. **清晰的字段映射规则**
8. **合理的按钮禁用逻辑**

## 7. 代码检查清单

- [ ] startEdit正确设置表单初始值
- [ ] saveEdit使用本地更新而非fetchCases
- [ ] cancelEdit清空编辑键和表单
- [ ] createDefaultEmptyRow使用正确的参数格式
- [ ] 所有useCallback声明了正确的依赖项
- [ ] 使用函数式setState更新数据
- [ ] 空值统一处理为空字符串
- [ ] 按钮禁用逻辑正确
- [ ] 错误处理完整且有用户提示
- [ ] TestResult字段即时保存实现正确
- [ ] 多语言字段映射规则正确
