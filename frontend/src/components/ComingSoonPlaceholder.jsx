import { Empty, Typography } from 'antd';
import { ClockCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

const ComingSoonPlaceholder = ({ feature }) => {
  const { t } = useTranslation();

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100%',
      minHeight: 400,
      background: '#fff',
      borderRadius: 8,
      padding: 48
    }}>
      <Empty
        image={<ClockCircleOutlined style={{ fontSize: 64, color: '#1890ff' }} />}
        imageStyle={{ height: 80 }}
        description={
          <div style={{ marginTop: 16 }}>
            <Text style={{ fontSize: 16, color: 'rgba(0, 0, 0, 0.65)' }}>
              {feature ? `${feature} - ${t('projectDetail.comingSoon')}` : t('projectDetail.comingSoon')}
            </Text>
          </div>
        }
      />
    </div>
  );
};

export default ComingSoonPlaceholder;
