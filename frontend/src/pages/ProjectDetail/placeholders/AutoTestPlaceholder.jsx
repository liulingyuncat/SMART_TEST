import React from 'react';
import { Empty } from 'antd';
import { useTranslation } from 'react-i18next';

const AutoTestPlaceholder = () => {
  const { t } = useTranslation();
  
  return (
    <Empty
      description={`${t('projectDetail.comingSoon')} (T13)`}
      style={{ padding: '60px 0' }}
    />
  );
};

export default AutoTestPlaceholder;
