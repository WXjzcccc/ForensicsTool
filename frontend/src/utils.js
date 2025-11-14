// 工具函数集合

/**
 * 生成普通文本输出
 * @param {string} text - 文本内容
 * @param {string} color - 文本颜色
 * @returns {string} 带样式的HTML字符串
 */
export function generateNormalTextOutput(text, color = 'black') {
  return `<span style="color: ${color}">${text}</span><br>`;
}

/**
 * 生成成功文本输出
 * @param {string} label - 标签文本
 * @param {string} value - 值文本
 * @returns {string} 带样式的HTML字符串
 */
export function generateSuccessTextOutput(label, value) {
  return `<span style="color: green">${label}: ${value}</span><br>`;
}

/**
 * 生成错误文本输出
 * @param {string} text - 文本内容
 * @returns {string} 带样式的HTML字符串
 */
export function generateErrorTextOutput(text) {
  return `<span style="color: red">${text}</span><br>`;
}

/**
 * 生成警告文本输出
 * @param {string} text - 文本内容
 * @returns {string} 带样式的HTML字符串
 */
export function generateWarningTextOutput(text) {
  return `<span style="color: orange">${text}</span><br>`;
}

/**
 * 生成信息文本输出
 * @param {string} text - 文本内容
 * @returns {string} 带样式的HTML字符串
 */
export function generateInfoTextOutput(text) {
  return `<span style="color: blue">${text}</span><br>`;
}