/**
 * 密码脱敏工具函数
 * 用于在前端显示时将密码替换为暗码
 * 不修改实际数据，只修改显示内容
 * 
 * 策略：优先使用精确密码替换（从元数据获取），其次使用模式匹配
 */

/**
 * 根据已知密码列表进行精确替换（推荐方式）
 * @param {string} text - 需要脱敏的文本
 * @param {Array<string>} knownPasswords - 已知的密码列表（从元数据获取）
 * @param {string} maskChar - 用于替换的字符，默认为'*'
 * @param {number} maskLength - 替换后的长度，默认为6
 * @returns {string} 脱敏后的文本
 */
export function maskKnownPasswords(text, knownPasswords = [], maskChar = '*', maskLength = 6) {
  if (!text || typeof text !== 'string') {
    return text;
  }
  
  if (!Array.isArray(knownPasswords) || knownPasswords.length === 0) {
    return text;
  }

  let maskedText = text;
  const mask = maskChar.repeat(maskLength);

  // 按密码长度降序排序，避免短密码被先替换导致长密码无法匹配
  const sortedPasswords = [...knownPasswords]
    .filter(p => p && typeof p === 'string' && p.length > 0)
    .sort((a, b) => b.length - a.length);

  sortedPasswords.forEach(password => {
    // 转义特殊正则字符
    const escapedPassword = password.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    const regex = new RegExp(escapedPassword, 'g');
    maskedText = maskedText.replace(regex, mask);
  });

  return maskedText;
}

/**
 * 通用模式匹配密码脱敏（备用方式，不推荐）
 * @param {string} text - 需要脱敏的文本
 * @param {string} maskChar - 用于替换的字符，默认为'*'
 * @param {number} maskLength - 替换后的长度，默认为6
 * @returns {string} 脱敏后的文本
 */
export function maskPasswords(text, maskChar = '*', maskLength = 6) {
  if (!text || typeof text !== 'string') {
    return text;
  }

  // 注意：此函数已废弃，建议使用 maskKnownPasswords
  return text;
}

/**
 * 对对象中的特定字段进行密码脱敏
 * @param {Object} obj - 需要脱敏的对象
 * @param {Array<string>} fields - 需要脱敏的字段名列表
 * @param {Array<string>} knownPasswords - 已知的密码列表
 * @param {string} maskChar - 用于替换的字符，默认为'*'
 * @param {number} maskLength - 替换后的长度，默认为6
 * @returns {Object} 脱敏后的对象（浅拷贝）
 */
export function maskObjectPasswords(obj, fields = [], knownPasswords = [], maskChar = '*', maskLength = 6) {
  if (!obj || typeof obj !== 'object') {
    return obj;
  }

  const maskedObj = { ...obj };

  fields.forEach(field => {
    if (maskedObj[field]) {
      maskedObj[field] = maskKnownPasswords(maskedObj[field], knownPasswords, maskChar, maskLength);
    }
  });

  return maskedObj;
}

/**
 * 预定义的需要脱敏的字段列表（用于Web/API用例）
 * 注意：script_code 不在列表中，保持脚本完整性
 */
export const WEB_CASE_PASSWORD_FIELDS = [
  'precondition_cn',
  'precondition_jp',
  'precondition_en',
  'test_steps_cn',
  'test_steps_jp',
  'test_steps_en',
  'expected_result_cn',
  'expected_result_jp',
  'expected_result_en',
];

export const API_CASE_PASSWORD_FIELDS = [
  'body',
  'header',
  'response',
  'url',
];

/**
 * 对Web用例数据进行密码脱敏
 * @param {Object} caseData - Web用例数据
 * @param {Array<string>} knownPasswords - 已知的密码列表
 * @returns {Object} 脱敏后的用例数据
 */
export function maskWebCasePasswords(caseData, knownPasswords = []) {
  return maskObjectPasswords(caseData, WEB_CASE_PASSWORD_FIELDS, knownPasswords);
}

/**
 * 对API用例数据进行密码脱敏
 * @param {Object} caseData - API用例数据
 * @param {Array<string>} knownPasswords - 已知的密码列表
 * @returns {Object} 脱敏后的用例数据
 */
export function maskApiCasePasswords(caseData, knownPasswords = []) {
  return maskObjectPasswords(caseData, API_CASE_PASSWORD_FIELDS, knownPasswords);
}

export default {
  maskPasswords,
  maskKnownPasswords,
  maskObjectPasswords,
  maskWebCasePasswords,
  maskApiCasePasswords,
  WEB_CASE_PASSWORD_FIELDS,
  API_CASE_PASSWORD_FIELDS,
};
