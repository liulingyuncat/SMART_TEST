import React, { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Tabs, Badge } from 'antd';
import {
  DEFECT_STATUS,
  DEFECT_STATUS_COLORS,
} from '../../../../constants/defect';

const { TabPane } = Tabs;

/**
 * 缺陷状态分栏标签
 * 显示各状态的缺陷数量，支持切换筛选
 */
const DefectStatusTabs = ({
  statusCounts = {},
  currentStatus,
  onStatusChange,
}) => {
  const { t, i18n } = useTranslation();

  // 使用 useMemo 缓存翻译标签，只在语言变化时重新计算
  const statusLabels = useMemo(() => ({
    all: t('common.all', '全部'),
    [DEFECT_STATUS.NEW]: t('defect.statusNew', '新建'),
    [DEFECT_STATUS.ACTIVE]: t('defect.statusActive', '处理中'),
    [DEFECT_STATUS.RESOLVED]: t('defect.statusResolved', '已解决'),
    [DEFECT_STATUS.CLOSED]: t('defect.statusClosed', '已关闭'),
  }), [t, i18n.language]);

  // 计算全部数量
  const totalCount = Object.values(statusCounts).reduce(
    (sum, count) => sum + count,
    0
  );

  // 状态Tab配置
  const statusTabs = useMemo(() => [
    { key: '', label: statusLabels.all, count: totalCount },
    ...Object.values(DEFECT_STATUS).map((status) => ({
      key: status,
      label: statusLabels[status],
      color: DEFECT_STATUS_COLORS[status],
      count: statusCounts[status] || 0,
    })),
  ], [statusLabels, statusCounts, totalCount]);

  return (
    <div className="defect-status-tabs" style={{ marginBottom: 16 }}>
      <Tabs
        activeKey={currentStatus}
        onChange={onStatusChange}
        type="card"
        size="small"
      >
        {statusTabs.map((tab) => (
          <TabPane
            key={tab.key}
            tab={
              <span>
                {tab.label}
                <Badge
                  count={tab.count}
                  style={{
                    marginLeft: 8,
                    backgroundColor: tab.color || '#999',
                  }}
                  overflowCount={999}
                />
              </span>
            }
          />
        ))}
      </Tabs>
    </div>
  );
};

export default DefectStatusTabs;
