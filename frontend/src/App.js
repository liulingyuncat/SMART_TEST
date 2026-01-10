import React, { useEffect, useState } from 'react';
import { Provider, useDispatch } from 'react-redux';
import { ConfigProvider, App as AntApp } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import enUS from 'antd/locale/en_US';
import jaJP from 'antd/locale/ja_JP';
import { useTranslation } from 'react-i18next';
import store from './store';
import AppRouter from './router';
import { login } from './store/authSlice';
import './i18n';
import './App.css';

// Ant Design locale映射
const localeMap = {
  zh: zhCN,
  en: enUS,
  ja: jaJP,
};

// 包装器组件用于初始化状态
function AppInitializer({ children }) {
  const dispatch = useDispatch();
  const [isInitialized, setIsInitialized] = useState(false);

  useEffect(() => {
    // 从 localStorage 恢复用户状态
    const token = localStorage.getItem('auth_token');
    const userInfo = localStorage.getItem('user_info');

    if (token && userInfo) {
      try {
        const user = JSON.parse(userInfo);
        dispatch(login({ token, user }));
      } catch (error) {
        console.error('Failed to parse user info:', error);
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user_info');
      }
    }
    
    // 标记初始化完成
    setIsInitialized(true);
  }, [dispatch]);

  // 等待初始化完成
  if (!isInitialized) {
    return null; // 或者显示一个加载指示器
  }

  return children;
}

function App() {
  const { i18n } = useTranslation();
  const [locale, setLocale] = useState(localeMap[i18n.language] || zhCN);

  useEffect(() => {
    // 监听i18n的languageChanged事件
    const handleLanguageChange = (lng) => {
      setLocale(localeMap[lng] || zhCN);
    };

    i18n.on('languageChanged', handleLanguageChange);

    // 清理事件监听器
    return () => {
      i18n.off('languageChanged', handleLanguageChange);
    };
  }, [i18n]);

  return (
    <ConfigProvider locale={locale}>
      <Provider store={store}>
        <AppInitializer>
          <AntApp>
            <AppRouter />
          </AntApp>
        </AppInitializer>
      </Provider>
    </ConfigProvider>
  );
}

export default App;
