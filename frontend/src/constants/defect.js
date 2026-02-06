/**
 * 缺陷管理常量定义
 */

// 缺陷状态
export const DEFECT_STATUS = {
  NEW: 'New',
  IN_PROGRESS: 'InProgress',  // 变更：Active → InProgress
  RESOLVED: 'Resolved',
  CLOSED: 'Closed',
  CONFIRMED: 'Confirmed',      // 新增
  REOPENED: 'Reopened',        // 新增
  REJECTED: 'Rejected',        // 新增
  // 向后兼容
  ACTIVE: 'Active',            // deprecated: use IN_PROGRESS
};

// 状态对应颜色
export const DEFECT_STATUS_COLORS = {
  [DEFECT_STATUS.NEW]: 'blue',
  [DEFECT_STATUS.IN_PROGRESS]: 'orange',
  [DEFECT_STATUS.ACTIVE]: 'orange',  // 向后兼容
  [DEFECT_STATUS.RESOLVED]: 'green',
  [DEFECT_STATUS.CLOSED]: 'default',
  [DEFECT_STATUS.CONFIRMED]: 'purple',
  [DEFECT_STATUS.REOPENED]: 'gold',
  [DEFECT_STATUS.REJECTED]: 'red',
};

// 状态对应国际化key
export const DEFECT_STATUS_I18N_KEYS = {
  [DEFECT_STATUS.NEW]: 'defect.statusNew',
  [DEFECT_STATUS.IN_PROGRESS]: 'defect.statusInProgress',
  [DEFECT_STATUS.ACTIVE]: 'defect.statusInProgress',  // 向后兼容
  [DEFECT_STATUS.RESOLVED]: 'defect.statusResolved',
  [DEFECT_STATUS.CLOSED]: 'defect.statusClosed',
  [DEFECT_STATUS.CONFIRMED]: 'defect.statusConfirmed',
  [DEFECT_STATUS.REOPENED]: 'defect.statusReopened',
  [DEFECT_STATUS.REJECTED]: 'defect.statusRejected',
};

// 缺陷优先级（按需求文档要求：A/B/C/D）
export const DEFECT_PRIORITY = {
  A: 'A',
  B: 'B',
  C: 'C',
  D: 'D',
};

// 优先级对应颜色
export const DEFECT_PRIORITY_COLORS = {
  [DEFECT_PRIORITY.A]: 'red',       // A - 最高优先级
  [DEFECT_PRIORITY.B]: 'orange',    // B - 高优先级
  [DEFECT_PRIORITY.C]: 'blue',      // C - 中优先级
  [DEFECT_PRIORITY.D]: 'default',   // D - 低优先级
};

// 优先级对应国际化key
export const DEFECT_PRIORITY_I18N_KEYS = {
  [DEFECT_PRIORITY.A]: 'defect.priorityA',
  [DEFECT_PRIORITY.B]: 'defect.priorityB',
  [DEFECT_PRIORITY.C]: 'defect.priorityC',
  [DEFECT_PRIORITY.D]: 'defect.priorityD',
};

// 缺陷严重程度（变更：ABCD → Critical/Major/Minor/Trivial）
export const DEFECT_SEVERITY = {
  CRITICAL: 'Critical',  // 变更：A → Critical
  MAJOR: 'Major',        // 变更：B → Major
  MINOR: 'Minor',        // 变更：C → Minor
  TRIVIAL: 'Trivial',    // 变更：D → Trivial
  // 向后兼容
  A: 'A',                // deprecated: use CRITICAL
  B: 'B',                // deprecated: use MAJOR
  C: 'C',                // deprecated: use MINOR
  D: 'D',                // deprecated: use TRIVIAL
};

// 严重程度对应颜色
export const DEFECT_SEVERITY_COLORS = {
  [DEFECT_SEVERITY.CRITICAL]: 'red',       // Critical - 最严重
  [DEFECT_SEVERITY.MAJOR]: 'orange',       // Major - 严重
  [DEFECT_SEVERITY.MINOR]: 'blue',         // Minor - 一般
  [DEFECT_SEVERITY.TRIVIAL]: 'default',    // Trivial - 轻微
  // 向后兼容
  [DEFECT_SEVERITY.A]: 'red',
  [DEFECT_SEVERITY.B]: 'orange',
  [DEFECT_SEVERITY.C]: 'blue',
  [DEFECT_SEVERITY.D]: 'default',
};

// 严重程度对应国际化key
export const DEFECT_SEVERITY_I18N_KEYS = {
  [DEFECT_SEVERITY.CRITICAL]: 'defect.severityCritical',
  [DEFECT_SEVERITY.MAJOR]: 'defect.severityMajor',
  [DEFECT_SEVERITY.MINOR]: 'defect.severityMinor',
  [DEFECT_SEVERITY.TRIVIAL]: 'defect.severityTrivial',
  // 向后兼容
  [DEFECT_SEVERITY.A]: 'defect.severityCritical',
  [DEFECT_SEVERITY.B]: 'defect.severityMajor',
  [DEFECT_SEVERITY.C]: 'defect.severityMinor',
  [DEFECT_SEVERITY.D]: 'defect.severityTrivial',
};

// 缺陷类型（新增）
export const DEFECT_TYPE = {
  FUNCTIONAL: 'Functional',
  UI: 'UI',
  UI_INTERACTION: 'UIInteraction',
  COMPATIBILITY: 'Compatibility',
  BROWSER_SPECIFIC: 'BrowserSpecific',
  PERFORMANCE: 'Performance',
  SECURITY: 'Security',
  ENVIRONMENT: 'Environment',
  USER_ERROR: 'UserError',
};

// 缺陷类型对应颜色
export const DEFECT_TYPE_COLORS = {
  [DEFECT_TYPE.FUNCTIONAL]: 'blue',
  [DEFECT_TYPE.UI]: 'cyan',
  [DEFECT_TYPE.UI_INTERACTION]: 'geekblue',
  [DEFECT_TYPE.COMPATIBILITY]: 'purple',
  [DEFECT_TYPE.BROWSER_SPECIFIC]: 'magenta',
  [DEFECT_TYPE.PERFORMANCE]: 'orange',
  [DEFECT_TYPE.SECURITY]: 'red',
  [DEFECT_TYPE.ENVIRONMENT]: 'green',
  [DEFECT_TYPE.USER_ERROR]: 'default',
};

// 缺陷类型对应国际化key
export const DEFECT_TYPE_I18N_KEYS = {
  [DEFECT_TYPE.FUNCTIONAL]: 'defect.typeFunctional',
  [DEFECT_TYPE.UI]: 'defect.typeUI',
  [DEFECT_TYPE.UI_INTERACTION]: 'defect.typeUIInteraction',
  [DEFECT_TYPE.COMPATIBILITY]: 'defect.typeCompatibility',
  [DEFECT_TYPE.BROWSER_SPECIFIC]: 'defect.typeBrowserSpecific',
  [DEFECT_TYPE.PERFORMANCE]: 'defect.typePerformance',
  [DEFECT_TYPE.SECURITY]: 'defect.typeSecurity',
  [DEFECT_TYPE.ENVIRONMENT]: 'defect.typeEnvironment',
  [DEFECT_TYPE.USER_ERROR]: 'defect.typeUserError',
};

// 状态流转规则
export const DEFECT_STATUS_TRANSITIONS = {
  [DEFECT_STATUS.NEW]: [DEFECT_STATUS.IN_PROGRESS, DEFECT_STATUS.CONFIRMED, DEFECT_STATUS.REJECTED, DEFECT_STATUS.CLOSED],
  [DEFECT_STATUS.IN_PROGRESS]: [DEFECT_STATUS.RESOLVED, DEFECT_STATUS.CLOSED],
  [DEFECT_STATUS.ACTIVE]: [DEFECT_STATUS.RESOLVED, DEFECT_STATUS.CLOSED],  // 向后兼容
  [DEFECT_STATUS.CONFIRMED]: [DEFECT_STATUS.IN_PROGRESS, DEFECT_STATUS.CLOSED],
  [DEFECT_STATUS.RESOLVED]: [DEFECT_STATUS.CLOSED, DEFECT_STATUS.REOPENED],
  [DEFECT_STATUS.REOPENED]: [DEFECT_STATUS.IN_PROGRESS, DEFECT_STATUS.CLOSED],
  [DEFECT_STATUS.REJECTED]: [DEFECT_STATUS.NEW, DEFECT_STATUS.CLOSED],
  [DEFECT_STATUS.CLOSED]: [DEFECT_STATUS.REOPENED],
};

// 获取状态选项列表
export const getStatusOptions = (t) => [
  { value: DEFECT_STATUS.NEW, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.NEW]) },
  { value: DEFECT_STATUS.IN_PROGRESS, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.IN_PROGRESS]) },
  { value: DEFECT_STATUS.CONFIRMED, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.CONFIRMED]) },
  { value: DEFECT_STATUS.RESOLVED, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.RESOLVED]) },
  { value: DEFECT_STATUS.REOPENED, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.REOPENED]) },
  { value: DEFECT_STATUS.REJECTED, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.REJECTED]) },
  { value: DEFECT_STATUS.CLOSED, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.CLOSED]) },
];

// 获取优先级选项列表
export const getPriorityOptions = (t) => [
  { value: DEFECT_PRIORITY.A, label: t(DEFECT_PRIORITY_I18N_KEYS[DEFECT_PRIORITY.A]) },
  { value: DEFECT_PRIORITY.B, label: t(DEFECT_PRIORITY_I18N_KEYS[DEFECT_PRIORITY.B]) },
  { value: DEFECT_PRIORITY.C, label: t(DEFECT_PRIORITY_I18N_KEYS[DEFECT_PRIORITY.C]) },
  { value: DEFECT_PRIORITY.D, label: t(DEFECT_PRIORITY_I18N_KEYS[DEFECT_PRIORITY.D]) },
];

// 获取严重程度选项列表
export const getSeverityOptions = (t) => [
  { value: DEFECT_SEVERITY.CRITICAL, label: t(DEFECT_SEVERITY_I18N_KEYS[DEFECT_SEVERITY.CRITICAL]) },
  { value: DEFECT_SEVERITY.MAJOR, label: t(DEFECT_SEVERITY_I18N_KEYS[DEFECT_SEVERITY.MAJOR]) },
  { value: DEFECT_SEVERITY.MINOR, label: t(DEFECT_SEVERITY_I18N_KEYS[DEFECT_SEVERITY.MINOR]) },
  { value: DEFECT_SEVERITY.TRIVIAL, label: t(DEFECT_SEVERITY_I18N_KEYS[DEFECT_SEVERITY.TRIVIAL]) },
];

// 获取缺陷类型选项列表（新增）
export const getTypeOptions = (t) => [
  { value: DEFECT_TYPE.FUNCTIONAL, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.FUNCTIONAL]) },
  { value: DEFECT_TYPE.UI, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.UI]) },
  { value: DEFECT_TYPE.UI_INTERACTION, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.UI_INTERACTION]) },
  { value: DEFECT_TYPE.COMPATIBILITY, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.COMPATIBILITY]) },
  { value: DEFECT_TYPE.BROWSER_SPECIFIC, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.BROWSER_SPECIFIC]) },
  { value: DEFECT_TYPE.PERFORMANCE, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.PERFORMANCE]) },
  { value: DEFECT_TYPE.SECURITY, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.SECURITY]) },
  { value: DEFECT_TYPE.ENVIRONMENT, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.ENVIRONMENT]) },
  { value: DEFECT_TYPE.USER_ERROR, label: t(DEFECT_TYPE_I18N_KEYS[DEFECT_TYPE.USER_ERROR]) },
];

// 根据当前状态获取可流转的状态选项
export const getAvailableStatusOptions = (currentStatus, t) => {
  const availableStatuses = DEFECT_STATUS_TRANSITIONS[currentStatus] || [];
  return availableStatuses.map(status => ({
    value: status,
    label: t(DEFECT_STATUS_I18N_KEYS[status]),
  }));
};
