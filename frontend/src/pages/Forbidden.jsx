import { Result, Button } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const Forbidden = () => {
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <Result
      status="403"
      title="403"
      subTitle={t('message.forbidden')}
      extra={
        <Button type="primary" onClick={() => navigate('/')}>
          {t('common.ok')}
        </Button>
      }
    />
  );
};

export default Forbidden;
