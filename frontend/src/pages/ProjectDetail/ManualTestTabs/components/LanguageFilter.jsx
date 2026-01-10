import React from 'react';
import { Radio } from 'antd';
import './LanguageFilter.css';

/**
 * 语言筛选组件 - 用于筛选用例数据的语言版本
 * 与顶部导航栏的UI语言切换独立
 * @param {Object} props
 * @param {string} props.value - 当前选中的语言(受控模式)
 * @param {Function} props.onChange - 语言变更回调
 * @param {Function} props.onLanguageChange - 语言变更回调(向后兼容)
 * @param {string} props.defaultLanguage - 默认语言(非受控模式)
 */
const LanguageFilter = ({ value, onChange, onLanguageChange, defaultLanguage = '中文' }) => {
  // 优先使用受控模式
  const currentValue = value !== undefined ? value : defaultLanguage;
  const handleChange = onChange || onLanguageChange;

  // 语言切换处理
  const handleLanguageChange = (e) => {
    const newLanguage = e.target.value;
    if (handleChange) {
      handleChange(newLanguage);
    }
  };

  return (
    <div className="language-filter">
      <Radio.Group 
        value={currentValue} 
        onChange={handleLanguageChange}
        buttonStyle="solid"
      >
        <Radio.Button value="中文">中文(CN)</Radio.Button>
        <Radio.Button value="English">English(EN)</Radio.Button>
        <Radio.Button value="日本語">日本語(JP)</Radio.Button>
      </Radio.Group>
    </div>
  );
};

export default LanguageFilter;