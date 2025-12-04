---
plan_id: T35-dev-plan-002
task_id: task_1763971682
task_name: T35-手工测试用例库-上方下方功能追加
status: in-progress
created_at: 2025-01-24T12:30:00Z
updated_at: 2025-01-24T12:30:00Z
dependencies:
  - step-01完成后可启动step-02至step-05
  - step-02至step-04可并行开发,建议顺序执行避免合并冲突
  - step-05依赖step-02至step-04全部完成
  - step-06作为最终验证环节
priority: P0-critical
estimated_hours: 5h
---

# T35-手工测试用例库-上方下方功能追加 开发执行计划

## 原始用户请求

为T35任务生成可执行的开发分解提示,该任务旨在增强Web智能测试平台手工测试用例库的Above和Below插入功能,实现以下核心目标:自动复制相邻用例的大功能和中功能分类字段,减少重复输入工作;针对整体和变更用例自动复制CN、JP、EN三种语言的分类字段,确保多语言版本一致性;通过颜色差异化显示相同分类的用例,提升用例组织和审查效率;正确处理连续多次Above或Below操作,确保字段复制逻辑的准确性

## 证据材料汇总

### 设计文档关键内容 version 7

设计目标包括实现自动分类复制,在Above和Below操作中自动复制相邻用例的大功能和中功能分类字段,减少90%的重复输入工作;多语言同步,针对整体和变更用例自动复制CN、JP、EN三种语言的分类字段;视觉增强,通过颜色差异化显示相同分类的用例;连续操作支持,正确处理连续多次Above或Below操作

核心函数设计包括5个关键函数:copyClassificationFields函数负责根据用例类型从源用例复制大功能和中功能分类字段到目标对象;getCellStyle函数负责根据当前用例和上一条用例的分类字段值计算单元格背景色;createEmptyRowByModule函数需要修改以支持AI用例字段初始化;handleInsertAbove函数在调用createEmptyRowByModule后立即合并分类字段,推荐使用函数式setState提升并发安全性;handleInsertBelow函数逻辑类似,源用例为当前用例

Table列定义修改涵盖两种场景:AI用例列定义需要为major_function和middle_function增加render函数,编辑状态返回undefined让EditableCell接管,非编辑状态应用颜色样式;整体和变更用例列定义需要修改major_function_cn、jp、en等列的render函数,保留多语言对话框逻辑,在非编辑状态下应用颜色样式

风险分析指出连续快速点击可能导致状态不一致,现有代码使用useState管理cases状态,React 18的自动批处理会合并连续的setState调用,推荐使用函数式setState基于最新状态计算,备选方案是增加loading状态锁

### 需求文档关键内容 version 2

功能需求FR-01规定Above按钮分类字段自动复制,当用户点击某条测试用例的Above按钮时,系统在该用例的上方插入一条新用例,新用例的大功能分类和中功能分类字段应自动复制自原用例的上一条用例,边界条件是如果目标用例是第1条用例则分类字段设置为空字符串

功能需求FR-02规定Below按钮分类字段自动复制,当用户点击某条测试用例的Below按钮时,系统在该用例的下方插入一条新用例,新用例的大功能分类和中功能分类字段应自动复制自被点击Below按钮的用例本身

功能需求FR-03规定多语言字段复制规则,针对整体用例和变更用例,系统需要同时复制中文、日文、英文三种语言的大功能和中功能分类字段,无论当前用户界面显示的是哪种语言,Above和Below操作都必须复制所有三种语言的字段

功能需求FR-04规定AI用例单语言处理,针对AI用例,系统只需要复制单语言字段major_function和middle_function,不涉及多语言后缀

功能需求FR-05规定相同分类的UI视觉提示,当某条用例的大功能分类或中功能分类与其上一条用例相同时,系统应在界面上以不同颜色显示该分类字段,中文字段相同时显示浅蓝色背景#E6F7FF,英文和日文字段相同时显示浅灰色背景#F5F5F5,AI用例使用浅灰色#F5F5F5

功能需求FR-06规定连续插入操作的字段复制逻辑,当用户在短时间内连续多次点击Above或Below按钮时,系统需要正确处理新插入用例的分类字段复制逻辑,确保每次插入都能准确地从正确的源用例复制分类字段

非功能性需求包括性能要求,点击Above或Below按钮后新用例应在500毫秒内插入并显示,颜色样式计算不应显著影响页面渲染性能,每页50条用例的渲染时间不超过100毫秒;可维护性要求,分类字段复制逻辑封装在独立的函数中便于复用和测试,颜色计算逻辑提取为工具函数便于维护

### 现有代码结构分析

EditableTable.jsx组件位于frontend/src/pages/ProjectDetail/ManualTestTabs/components目录,总计2388行代码,使用React 18、Ant Design 5.x、Form、Table等组件

createEmptyRowByModule函数位于第415至475行,当前实现对于api-cases返回接口测试用例字段,对于非api-cases仅返回多语言字段,缺少AI用例的单语言字段major_function和middle_function等

handleInsertAbove函数位于第480至530行,当前逻辑是先取消正在编辑的行,设置loading和hasEditChanges状态,找到目标行,调用createEmptyRowByModule创建空行,在目标行上方splice插入空行,对于非api-cases重新计算后续用例的id,最后setCases更新状态

handleInsertBelow函数位于第535至585行,逻辑类似handleInsertAbove,区别是插入位置为targetIndex加1,重新计算id的起始索引为targetIndex加2

AI用例列定义位于第1220至1280行,major_function和middle_function列仅通过onCell配置行内编辑,没有render函数,EditableCell组件控制显示逻辑

整体和变更用例列定义位于第1520至1650行,major_function_cn、jp、en等列已有render函数,编辑状态且非中文模式返回undefined让EditableCell接管,中文模式返回可点击div打开多语言对话框,非编辑状态返回普通文本

## 分步开发计划

### 步骤一 新增copyClassificationFields和getCellStyle工具函数

- [ ] 在EditableTable.jsx顶部第100至150行之间新增copyClassificationFields函数,封装分类字段复制逻辑
- [ ] 函数接收sourceCase和caseType两个参数,返回classificationFields对象
- [ ] 如果caseType等于ai则返回包含major_function和middle_function的对象,使用sourceCase.major_function或空字符串的语法确保空值安全
- [ ] 如果caseType等于overall或change则返回包含major_function_cn、major_function_jp、major_function_en、middle_function_cn、middle_function_jp、middle_function_en的对象
- [ ] 如果caseType为未知值则返回空对象
- [ ] 在EditableTable.jsx顶部新增getCellStyle函数,封装颜色计算逻辑
- [ ] 函数接收currentCase、previousCase、fieldName三个参数,返回styleObject样式对象
- [ ] 如果previousCase为null则返回空对象处理第一行情况
- [ ] 获取currentValue和previousValue,如果当前值或上一条值为null或空字符串则返回空对象
- [ ] 如果当前值不等于上一条值则返回空对象
- [ ] 如果值相同则根据字段后缀判断颜色,如果fieldName以下划线cn结尾则返回backgroundColor为#E6F7FF的对象代表中文浅蓝色
- [ ] 如果fieldName以下划线jp或下划线en结尾则返回backgroundColor为#F5F5F5的对象代表英文日文浅灰色
- [ ] 如果fieldName为major_function或middle_function则返回backgroundColor为#F5F5F5的对象代表AI用例浅灰色,其他情况返回空对象
- [ ] 验收标准是两个函数能正确处理null、undefined、空字符串等边界情况
- [ ] copyClassificationFields能根据caseType返回正确的字段集合,getCellStyle能根据字段名返回正确的颜色代码
- [ ] 预估工时1小时,优先级P0-critical

### 步骤二 修改createEmptyRowByModule函数增加AI用例字段支持

- [ ] 定位第415至475行的createEmptyRowByModule函数,在非api-cases分支中增加字段判断逻辑
- [ ] 首先保留现有的baseRow对象定义包含case_id、id、display_id、project_id、case_type等公共字段
- [ ] 然后新增判断const isAICase等于targetRow.case_type等于ai
- [ ] 如果isAICase为真则返回展开baseRow加上major_function空字符串、middle_function空字符串、minor_function空字符串、precondition空字符串、test_steps空字符串、expected_result空字符串
- [ ] 否则保持原有多语言字段逻辑返回展开baseRow加上major_function_cn至expected_result_jp等所有多语言字段均为空字符串
- [ ] 确保函数签名保持不变仍为targetRow和targetNo两个参数,不增加新参数避免破坏现有调用方
- [ ] 验收标准是AI用例新增行包含major_function等单语言字段
- [ ] 整体和变更用例新增行包含多语言字段,接口测试用例api-cases仍正常工作不受影响
- [ ] 预估工时0.5小时,优先级P0-critical

### 步骤三 修改handleInsertAbove函数集成分类复制逻辑

- [ ] 定位第480至530行的handleInsertAbove函数,在调用createEmptyRowByModule获取emptyRow后立即增加分类复制逻辑
- [ ] 首先计算sourceIndex等于targetIndex减1
- [ ] 然后判断如果sourceIndex大于等于0则获取sourceCase等于cases数组中索引为sourceIndex的元素
- [ ] 调用copyClassificationFields传入sourceCase和targetRow.case_type获取classificationFields
- [ ] 使用Object.assign合并emptyRow和classificationFields
- [ ] 如果sourceIndex小于0则不执行任何操作分类字段保持为空处理第一条用例Above的边界情况
- [ ] 推荐使用函数式setState提升并发安全性,将setCases的参数改为箭头函数接收prevCases
- [ ] 在函数内部const newCases等于展开prevCases,执行newCases.splice在targetIndex位置插入0个删除插入emptyRow
- [ ] 然后判断如果apiModule不等于api-cases则遍历从targetIndex加1到newCases.length减1
- [ ] 将newCases数组中索引i的元素的id和display_id均加1,最后return newCases
- [ ] 验收标准是点击第N条用例的Above按钮新用例的分类与第N减1条相同
- [ ] 点击第1条用例的Above按钮新用例的分类为空
- [ ] 整体和变更用例的多语言字段CN、JP、EN全部复制
- [ ] 连续点击2次Above第2次插入的用例复制第1次插入用例的分类
- [ ] 预估工时1小时,优先级P0-critical

### 步骤四 修改handleInsertBelow函数集成分类复制逻辑

- [ ] 定位第535至585行的handleInsertBelow函数,在调用createEmptyRowByModule获取emptyRow后立即增加分类复制逻辑
- [ ] 由于Below操作的源用例为当前用例本身无需边界检查,直接获取sourceCase等于cases数组中索引为targetIndex的元素
- [ ] 调用copyClassificationFields传入sourceCase和targetRow.case_type获取classificationFields
- [ ] 使用Object.assign合并emptyRow和classificationFields
- [ ] 使用函数式setState,将setCases的参数改为箭头函数接收prevCases
- [ ] 在函数内部const newCases等于展开prevCases,执行newCases.splice在targetIndex加1位置插入0个删除插入emptyRow
- [ ] 然后判断如果apiModule不等于api-cases则遍历从targetIndex加2到newCases.length减1
- [ ] 将newCases数组中索引i的元素的id和display_id均加1,最后return newCases
- [ ] 验收标准是点击第N条用例的Below按钮新用例的分类与第N条相同
- [ ] 连续点击Below新用例的分类保持一致,多语言字段全部复制
- [ ] 预估工时0.5小时,优先级P1-high

### 步骤五 修改Table列定义增加颜色提示

- [ ] 首先修改AI用例列定义位于第1220至1280行,为major_function和middle_function列增加render函数
- [ ] render函数接收text、record、index三个参数
- [ ] 首先判断如果editingKey等于record.case_id则return undefined让EditableCell接管
- [ ] 否则定义previousCase等于如果index大于0则cases数组中索引为index减1的元素否则为null
- [ ] 调用getCellStyle传入record、previousCase、major_function获取style对象
- [ ] 返回span标签设置style属性为style内容为text或短横线,middle_function列应用相同逻辑
- [ ] 然后修改整体和变更用例列定义位于第1520至1650行,修改major_function_cn、jp、en等列的render函数
- [ ] 保留现有的编辑状态判断和中文模式多语言对话框逻辑
- [ ] 在非编辑状态且非中文模式的分支中,定义previousCase等于如果index大于0则cases数组中索引为index减1的元素否则为null
- [ ] 调用getCellStyle传入record、previousCase、拼接major_function加上langFieldSuffix获取style对象
- [ ] 返回span标签设置style属性为style内容为fieldValue或短横线,middle_function_cn、jp、en等列应用相同逻辑
- [ ] 验收标准是整体用例中文模式下点击字段仍能打开多语言对话框
- [ ] 非编辑状态下分类相同的连续用例显示相同颜色
- [ ] 中文字段相同时显示浅蓝色#E6F7FF,英文和日文字段相同时显示浅灰色#F5F5F5
- [ ] 第一条用例无颜色样式,AI用例分类相同时显示浅灰色#F5F5F5
- [ ] 预估工时1.5小时,优先级P1-high

### 步骤六 测试验证

- [ ] 手工测试必须完成以下测试案例
- [ ] TC-01整体用例Above操作验证复制CN、JP、EN字段
- [ ] TC-02整体用例第1条Above操作验证字段为空
- [ ] TC-03 AI用例Above操作验证复制单语言字段
- [ ] TC-05整体用例Below操作
- [ ] TC-07 AI用例Below操作
- [ ] TC-14中文字段相同时显示浅蓝色
- [ ] TC-15英文字段相同时显示浅灰色
- [ ] TC-20连续2次Above操作验证字段传播
- [ ] TC-21连续3次Below操作
- [ ] TC-25连续插入后保存成功
- [ ] 回归测试包括接口测试用例Above和Below功能不受影响
- [ ] 整体用例多语言对话框正常打开,行内编辑功能正常
- [ ] 验收标准是所有测试案例通过,无控制台错误
- [ ] 性能指标插入响应小于500毫秒渲染小于100毫秒
- [ ] 预估工时0.5小时,优先级P2-normal

## 执行计划总结

本开发任务共分解为6个步骤,总预估工时5小时,关键依赖关系是步骤一新增工具函数为基础必须先行完成,步骤二、步骤三、步骤四可并行开发但建议顺序执行避免合并冲突,步骤五依赖步骤二至步骤四全部完成,步骤六作为最终验证环节

优先级排序P0最高级包括步骤一新增工具函数作为基础设施、步骤二修复AI用例字段缺失、步骤三Above复制逻辑核心功能,P1高级包括步骤四Below复制逻辑、步骤五UI颜色提示增强用户体验,P2普通级包括步骤六测试验证确保质量

风险提示包括连续快速点击可能导致状态不一致,已通过函数式setState基于最新状态计算缓解,详见设计文档5.1风险1;大数据量下颜色计算可能影响性能,当前单页50条用例可接受,如需支持单页100条用例建议使用React.memo或useMemo优化

后续待决策项包括是否复制minor_function小功能分类,颜色样式是否支持用户自定义,Above和Below按钮是否增加Tooltip提示将自动复制分类字段

参考文档包括设计文档T35详细设计文档版本7,需求文档T35需求规格说明书版本2,代码位置frontend/src/pages/ProjectDetail/ManualTestTabs/components/EditableTable.jsx

---

生成时间: 2025-01-24T12:30:00Z
文档版本: 2.0
