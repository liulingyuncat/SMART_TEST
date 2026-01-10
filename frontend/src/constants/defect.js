/**
 * 缺陷管理常量定义
 */

// 缺陷状态
export const DEFECT_STATUS = {
  NEW: 'New',
  ACTIVE: 'Active',
  RESOLVED: 'Resolved',
  CLOSED: 'Closed',
};

// 状态对应颜色
export const DEFECT_STATUS_COLORS = {
  [DEFECT_STATUS.NEW]: 'blue',
  [DEFECT_STATUS.ACTIVE]: 'orange',
  [DEFECT_STATUS.RESOLVED]: 'green',
  [DEFECT_STATUS.CLOSED]: 'default',
};

// 状态对应国际化key
export const DEFECT_STATUS_I18N_KEYS = {
  [DEFECT_STATUS.NEW]: 'defect.statusNew',
  [DEFECT_STATUS.ACTIVE]: 'defect.statusActive',
  [DEFECT_STATUS.RESOLVED]: 'defect.statusResolved',
  [DEFECT_STATUS.CLOSED]: 'defect.statusClosed',
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

// 缺陷严重程度（按需求文档要求：A/B/C/D）
export const DEFECT_SEVERITY = {
  A: 'A',
  B: 'B',
  C: 'C',
  D: 'D',
};

// 严重程度对应颜色
export const DEFECT_SEVERITY_COLORS = {
  [DEFECT_SEVERITY.A]: 'red',       // A - 最严重
  [DEFECT_SEVERITY.B]: 'orange',    // B - 严重
  [DEFECT_SEVERITY.C]: 'blue',      // C - 一般
  [DEFECT_SEVERITY.D]: 'default',   // D - 轻微
};

// 严重程度对应国际化key
export const DEFECT_SEVERITY_I18N_KEYS = {
  [DEFECT_SEVERITY.A]: 'defect.severityA',
  [DEFECT_SEVERITY.B]: 'defect.severityB',
  [DEFECT_SEVERITY.C]: 'defect.severityC',
  [DEFECT_SEVERITY.D]: 'defect.severityD',
};

// 状态流转规则
export const DEFECT_STATUS_TRANSITIONS = {
  [DEFECT_STATUS.NEW]: [DEFECT_STATUS.ACTIVE, DEFECT_STATUS.CLOSED],
  [DEFECT_STATUS.ACTIVE]: [DEFECT_STATUS.RESOLVED, DEFECT_STATUS.CLOSED],
  [DEFECT_STATUS.RESOLVED]: [DEFECT_STATUS.CLOSED, DEFECT_STATUS.ACTIVE],
  [DEFECT_STATUS.CLOSED]: [DEFECT_STATUS.ACTIVE],
};

// 获取状态选项列表
export const getStatusOptions = (t) => [
  { value: DEFECT_STATUS.NEW, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.NEW]) },
  { value: DEFECT_STATUS.ACTIVE, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.ACTIVE]) },
  { value: DEFECT_STATUS.RESOLVED, label: t(DEFECT_STATUS_I18N_KEYS[DEFECT_STATUS.RESOLVED]) },
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
  { value: DEFECT_SEVERITY.A, label: t(DEFECT_SEVERITY_I18N_KEYS[DEFECT_SEVERITY.A]) },
  { value: DEFECT_SEVERITY.B, label: t(DEFECT_SEVERITY_I18N_KEYS[DEFECT_SEVERITY.B]) },
  { value: DEFECT_SEVERITY.C, label: t(DEFECT_SEVERITY_I18N_KEYS[DEFECT_SEVERITY.C]) },
  { value: DEFECT_SEVERITY.D, label: t(DEFECT_SEVERITY_I18N_KEYS[DEFECT_SEVERITY.D]) },
];

// 根据当前状态获取可流转的状态选项
export const getAvailableStatusOptions = (currentStatus, t) => {
  const availableStatuses = DEFECT_STATUS_TRANSITIONS[currentStatus] || [];
  return availableStatuses.map(status => ({
    value: status,
    label: t(DEFECT_STATUS_I18N_KEYS[status]),
  }));
};
