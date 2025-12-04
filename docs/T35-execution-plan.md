---
plan_id: "35f7b8c1-d4e2-4a9f-b3c5-8a7e9f1d2c4b"
task_id: "task_1763971682"
status: "Pending Approval"
created_at: "2025-11-24T09:10:00Z"
updated_at: "2025-11-24T09:10:00Z"
dependencies:
  - { source: 'step-02', target: 'step-01' }
  - { source: 'step-03', target: 'step-01' }
  - { source: 'step-04', target: 'step-02' }
  - { source: 'step-04', target: 'step-03' }
  - { source: 'step-05', target: 'step-04' }
  - { source: 'step-06', target: 'step-05' }
---
- [ ] step-01: 在EditableTable.jsx文件顶部声明区域（约第50行后）新增两个工具函数，第一个函数名为copyClassificationFields，接收sourceCase对象和caseType字符串两个参数，函数内部通过if判断caseType值，当值为ai时返回包含major_function和middle_function两个字段的对象（使用或运算符确保空值转换为空字符串），当值为overall或change时返回包含六个多语言字段major_function_cn、major_function_jp、major_function_en、middle_function_cn、middle_function_jp、middle_function_en的对象（均使用或运算符处理空值），其他情况返回空对象。第二个函数名为getCellStyle，接收currentCase对象、previousCase对象或null、fieldName字符串三个参数，函数内部首先判断previousCase是否为null则返回空对象，然后获取currentValue和previousValue分别为currentCase和previousCase的fieldName属性值，如果任一值为null或空字符串或两值不相等则返回空对象，否则根据fieldName是否以_cn结尾返回backgroundColor为#E6F7FF的对象，或根据fieldName是否以_jp或_en结尾或等于major_function或middle_function返回backgroundColor为#F5F5F5的对象，其他情况返回空对象。 priority:high
- [ ] step-02: 修改EditableTable.jsx中的createEmptyRowByModule函数（约第415行），在函数内部非api-cases分支的返回语句前增加if判断，检查targetRow.case_type是否等于ai，如果等于则返回包含baseRow展开和ai用例专属字段的对象，专属字段包括major_function、middle_function、minor_function、precondition、test_steps、expected_result且均初始化为空字符串，否则继续执行原有的多语言字段返回逻辑。此修改确保AI用例和整体变更用例的字段初始化正确。 priority:high
- [ ] step-03: 修改EditableTable.jsx中的handleInsertAbove函数（约第480-550行），在调用createEmptyRowByModule创建emptyRow后、使用splice插入newCases数组前，增加分类字段复制逻辑，具体为计算sourceIndex等于targetIndex减1，使用if判断sourceIndex是否大于等于0，若成立则从cases数组获取sourceCase即cases[sourceIndex]，然后调用copyClassificationFields函数传入sourceCase和targetRow.case_type获取classificationFields对象，最后使用对象展开语法将classificationFields合并到emptyRow即emptyRow等于展开emptyRow再展开classificationFields的新对象。同时将原有的setCases调用修改为函数式setState，即setCases接收prevCases箭头函数参数，函数内部对prevCases进行数组展开、splice插入、id重新计算等操作后返回newCases，确保基于最新状态计算。此修改实现Above按钮的智能分类复制和连续点击防护。 priority:high
- [ ] step-04: 修改EditableTable.jsx中的handleInsertBelow函数（约第550-620行），在调用createEmptyRowByModule创建emptyRow后、使用splice插入newCases数组前，增加分类字段复制逻辑，具体为直接从cases数组获取sourceCase即cases[targetIndex]（Below操作源用例为当前用例），然后调用copyClassificationFields函数传入sourceCase和targetRow.case_type获取classificationFields对象，最后使用对象展开语法将classificationFields合并到emptyRow。同时将原有的setCases调用修改为函数式setState，逻辑同step-03。此修改实现Below按钮的智能分类复制和连续点击防护。 priority:high
- [ ] step-05: 修改EditableTable.jsx中的columns定义函数（约第1220-1650行），分两个方案处理，方案一针对AI用例的major_function和middle_function列（约第1220-1250行），在现有列定义对象中增加render函数属性，render函数接收text、record、index三个参数，函数内部首先判断editingKey是否等于record.case_id则返回undefined（让EditableCell接管编辑状态），否则计算previousCase为index大于0时的cases[index减1]或null，调用getCellStyle传入record、previousCase、对应字段名获取style对象，最后返回span标签包裹text或短横线并应用style样式。方案二针对整体变更用例的major_function_cn/jp/en和middle_function_cn/jp/en列（约第1552-1650行），在现有render函数的第三个分支（非编辑状态返回普通文本的位置），修改返回语句，在返回前计算previousCase为index大于0时的cases[index减1]或null，调用getCellStyle传入record、previousCase、字段名加langFieldSuffix获取style对象，将原return语句的文本包裹在span标签内并应用style样式。此修改实现相同分类的颜色视觉提示功能。 priority:medium
- [ ] step-06: 执行端到端测试验证所有功能需求FR-01至FR-06的25个测试案例，重点测试Above按钮复制逻辑TC-01至TC-04、Below按钮复制逻辑TC-05至TC-07、多语言字段复制TC-08至TC-10、AI用例处理TC-11至TC-13、颜色提示TC-14至TC-19、连续插入TC-20至TC-25，确保每个测试案例通过后标记为完成，如发现问题则回退到对应步骤修复代码并重新测试，所有测试通过后在项目根目录创建T35-implementation-report.md文档记录实施结果包括修改的代码行数、测试通过率、遗留问题（如有）。 priority:medium
