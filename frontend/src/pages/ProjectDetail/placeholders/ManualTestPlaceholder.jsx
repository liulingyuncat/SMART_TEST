import React from 'react';
import { Empty } from 'antd';
import { useTranslation } from 'react-i18next';

const ManualTestPlaceholder = () => {
  const { t } = useTranslation();
  
  return (
    <Empty
      description={`${t('projectDetail.comingSoon')} (T10、T11、T12)`}
      style={{ padding: '60px 0' }}
    />
  );
};

export default ManualTestPlaceholder;
