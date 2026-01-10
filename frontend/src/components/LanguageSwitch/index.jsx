import React from 'react';
import PropTypes from 'prop-types';
import { Dropdown, Button, Select } from 'antd';
import { GlobalOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { languageOptions } from './config';

/**
 * 语言切换组件
 * @param {Object} props - 组件属性
 * @param {'dropdown'|'select'} props.variant - UI样式变体 (dropdown用于Header, select用于Login)
 * @param {boolean} props.showLabel - 是否显示"语言"文字
 * @param {string} props.className - 自定义样式类
 */
function LanguageSwitch({ variant, showLabel, className }) {
  const { t, i18n } = useTranslation();

  const handleChange = (languageKey) => {
    // 1. 切换i18n语言
    i18n.changeLanguage(languageKey);

    // 2. 保存到localStorage
    try {
      localStorage.setItem('user_language', languageKey);
    } catch (e) {
      console.warn('localStorage not available, language preference will not persist');
    }
  };

  if (variant === 'dropdown') {
    // Dropdown样式(用于Header)
    const menuItems = languageOptions.map((lang) => ({
      key: lang.key,
      label: lang.label,
    }));

    return (
      <Dropdown
        menu={{
          items: menuItems,
          onClick: ({ key }) => handleChange(key),
        }}
        className={className}
      >
        <Button type="text" icon={<GlobalOutlined />}>
          {showLabel && t('common.language')}
        </Button>
      </Dropdown>
    );
  } else {
    // Select样式(用于登录页)
    return (
      <Select
        value={i18n.language}
        onChange={handleChange}
        options={languageOptions.map((lang) => ({
          value: lang.key,
          label: lang.label,
        }))}
        style={{ width: 120 }}
        className={className}
      />
    );
  }
}

LanguageSwitch.propTypes = {
  variant: PropTypes.oneOf(['dropdown', 'select']),
  showLabel: PropTypes.bool,
  className: PropTypes.string,
};

LanguageSwitch.defaultProps = {
  variant: 'dropdown',
  showLabel: true,
  className: '',
};

export default LanguageSwitch;
