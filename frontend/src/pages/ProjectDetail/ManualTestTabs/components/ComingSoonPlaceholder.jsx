import React from 'react';
import { useTranslation } from 'react-i18next';
import { Empty } from 'antd';

/**
 * 占位组件 - 用于未实现的Tab页面
 * @param {Object} props
 * @param {string} props.taskId - 任务ID (T11/T12)
 */
const ComingSoonPlaceholder = ({ taskId = '' }) => {
  const { t } = useTranslation();

  const descriptionKey = taskId === 'T11' 
    ? 'manualTest.comingSoonT11' 
    : taskId === 'T12'
    ? 'manualTest.comingSoonT12'
    : 'projectDetail.comingSoon';

  return (
    <div style={{ padding: '60px 0', textAlign: 'center' }}>
      <Empty
        description={t(descriptionKey)}
        image={Empty.PRESENTED_IMAGE_SIMPLE}
      />
    </div>
  );
};

export default ComingSoonPlaceholder;
