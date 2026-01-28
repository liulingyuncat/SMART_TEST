import { useState, useEffect, useMemo, Fragment } from 'react';
import { Form, Input, Select, Button, Card, Spin, message, Space, Row, Col, Statistic, Progress, Tag, Empty, Typography } from 'antd';
import { EditOutlined, SaveOutlined, CloseOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { getProjectById, updateProject } from '../api/project';
import { getUsers } from '../api/user';
import { getExecutionTasks } from '../api/executionTask';
import { getExecutionStatistics } from '../api/executionCaseResult';
import DefectTrendChart from './DefectTrendChart';

const { TextArea } = Input;
const { Option } = Select;
const { Text } = Typography;

const ProjectInfoTab = ({ projectId }) => {
  const { t } = useTranslation();
  const { user } = useSelector(state => state.auth);
  const [form] = Form.useForm();

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [editing, setEditing] = useState(false);
  const [projectData, setProjectData] = useState(null);
  const [projectManagers, setProjectManagers] = useState([]);
  const [canEdit, setCanEdit] = useState(false);
  const [ownerName, setOwnerName] = useState('');

  // 任务相关状态
  const [tasks, setTasks] = useState([]);
  const [taskLoading, setTaskLoading] = useState(false);
  const [statsMap, setStatsMap] = useState({});
  // 判断是否可编辑 (仅项目管理员)
  useEffect(() => {
    setCanEdit(user?.role === 'project_manager');
  }, [user]);

  // 加载项目管理员列表
  useEffect(() => {
    const fetchManagers = async () => {
      try {
        const data = await getUsers();
        const usersList = Array.isArray(data) ? data : [];
        const managers = usersList.filter(u => u.role === 'project_manager');
        setProjectManagers(managers);
      } catch (error) {
        console.error('[ProjectInfoTab] Failed to load managers:', error);
      }
    };

    fetchManagers();
  }, []);

  // 加载项目详情
  useEffect(() => {
    const fetchProject = async () => {
      try {
        setLoading(true);
        const project = await getProjectById(projectId);
        console.log('[ProjectInfoTab] Project data received:', project);
        console.log('[ProjectInfoTab] owner_name:', project.owner_name);
        setProjectData(project);

        // 直接使用后端返回的owner_name
        setOwnerName(project.owner_name || '');

        form.setFieldsValue({
          name: project.name,
          description: project.description || '',
          status: project.status || 'pending',
          owner_id: project.owner_id || undefined,
        });
      } catch (error) {
        console.error('[ProjectInfoTab] Failed to load project:', error);
        message.error(t('projectDetail.loadFailed'));
      } finally {
        setLoading(false);
      }
    };

    fetchProject();
  }, [projectId, form, t]);

  // 加载项目任务列表
  useEffect(() => {
    const fetchTasks = async () => {
      try {
        setTaskLoading(true);
        const data = await getExecutionTasks(projectId);
        setTasks(Array.isArray(data) ? data : []);
      } catch (error) {
        console.error('[ProjectInfoTab] Failed to load tasks:', error);
        message.error(t('projectInfo.taskLoadFailed'));
      } finally {
        setTaskLoading(false);
      }
    };

    if (projectId) {
      fetchTasks();
    }
  }, [projectId, t]);

  // 加载进行中任务的统计数据
  useEffect(() => {
    const fetchStatistics = async () => {
      const inProgressTasks = tasks.filter(task => task.task_status === 'in_progress');
      if (inProgressTasks.length === 0) return;

      const statsPromises = inProgressTasks.map(async (task) => {
        try {
          const stats = await getExecutionStatistics(task.task_uuid);
          return { uuid: task.task_uuid, stats };
        } catch (error) {
          console.error(`[ProjectInfoTab] Failed to load stats for task ${task.task_uuid}:`, error);
          return { uuid: task.task_uuid, stats: null };
        }
      });

      const results = await Promise.all(statsPromises);
      const newStatsMap = {};
      results.forEach(({ uuid, stats }) => {
        if (stats) {
          newStatsMap[uuid] = stats;
        }
      });
      setStatsMap(newStatsMap);
    };

    if (tasks.length > 0) {
      fetchStatistics();
    }
  }, [tasks]);

  // 计算任务统计数据
  const { inProgressCount, pendingCount, completedCount, inProgressTasks } = useMemo(() => {
    const inProgress = tasks.filter(task => task.task_status === 'in_progress');
    const pending = tasks.filter(task => task.task_status === 'pending');
    const completed = tasks.filter(task => task.task_status === 'completed');
    return {
      inProgressCount: inProgress.length,
      pendingCount: pending.length,
      completedCount: completed.length,
      inProgressTasks: inProgress,
    };
  }, [tasks]);

  // 计算所有进行中任务的总统计
  const totalStats = useMemo(() => {
    let total_ok = 0, total_ng = 0, total_block = 0, total_nr = 0;

    inProgressTasks.forEach(task => {
      const stats = statsMap[task.task_uuid];
      if (stats) {
        total_ok += stats.ok_count || 0;
        total_ng += stats.ng_count || 0;
        total_block += stats.block_count || 0;
        total_nr += stats.nr_count || 0;
      }
    });

    const total = total_ok + total_ng + total_block + total_nr;
    const executed = total - total_nr;
    const progress = total > 0 ? Math.round((executed / total) * 100) : 0;
    const passRate = executed > 0 ? Math.round((total_ok / executed) * 100) : 0;

    return {
      ok: total_ok,
      ng: total_ng,
      block: total_block,
      nr: total_nr,
      total,
      progress,
      passRate,
    };
  }, [inProgressTasks, statsMap]);

  // 开始编辑
  const handleEdit = () => {
    setEditing(true);
  };

  // 取消编辑
  const handleCancel = () => {
    setEditing(false);
    // 恢复原始值
    form.setFieldsValue({
      name: projectData.name,
      description: projectData.description || '',
      status: projectData.status || 'pending',
      owner_id: projectData.owner_id || undefined,
    });
  };

  // 保存项目信息
  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);

      const updates = {
        name: values.name,
        description: values.description,
        status: values.status,
        owner_id: values.owner_id || 0, // 0 表示清空负责人
      };

      await updateProject(projectId, updates);

      // 更新本地数据
      const updatedProject = await getProjectById(projectId);
      setProjectData(updatedProject);

      message.success(t('project.updateSuccess') || '保存成功');
      setEditing(false);
    } catch (error) {
      console.error('[ProjectInfoTab] Failed to save project:', error);
      if (error.errorFields) {
        // 表单验证错误
        return;
      }
      message.error(t('message.updateFailed'));
    } finally {
      setSaving(false);
    }
  };

  // 格式化日期
  const formatDate = (dateString) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    });
  };

  // 获取状态显示文本
  const getStatusText = (status) => {
    const statusMap = {
      'pending': t('project.statusPending'),
      'in-progress': t('project.statusInProgress'),
      'completed': t('project.statusCompleted')
    };
    return statusMap[status] || status;
  };

  // 获取负责人名称
  const getOwnerName = () => {
    // 优先使用后端直接返回的owner_name（创建人昵称）
    if (ownerName) return ownerName;
    // 降级方案：从projectManagers列表查找
    if (!projectData?.owner_id) return '-';
    const owner = projectManagers.find(m => m.id === projectData.owner_id);
    return owner ? (owner.nickname || owner.username) : '-';
  };

  // 计算执行进度
  const calculateExecutionProgress = (stats) => {
    if (!stats) return 0;
    const { ok_count = 0, ng_count = 0, block_count = 0, nr_count = 0 } = stats;
    const total = ok_count + ng_count + block_count + nr_count;
    if (total === 0) return 0;
    return Math.round(((ok_count + ng_count + block_count) / total) * 100);
  };

  // 计算通过率
  const calculatePassRate = (stats) => {
    if (!stats) return 0;
    const { ok_count = 0, ng_count = 0, block_count = 0, nr_count = 0 } = stats;
    const total = ok_count + ng_count + block_count + nr_count;
    const executed = total - nr_count;
    if (executed === 0) return 0;
    return Math.round((ok_count / executed) * 100);
  };

  // 渲染总统计区域
  const renderTotalStatistics = () => {
    if (inProgressTasks.length === 0 || totalStats.total === 0) return null;

    return (
      <div style={{
        background: 'linear-gradient(135deg, #f6f9fc 0%, #eef2f7 100%)',
        borderRadius: '8px',
        padding: '16px 20px',
        marginBottom: '16px',
        border: '1px solid #e3e8ef',
      }}>
        <div style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          flexWrap: 'wrap',
          gap: '12px',
        }}>
          {/* 左侧：总统计标题 */}
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
          }}>
            <div style={{
              width: '4px',
              height: '20px',
              background: 'linear-gradient(180deg, #1890ff 0%, #096dd9 100%)',
              borderRadius: '2px',
            }} />
            <span style={{
              fontSize: '14px',
              fontWeight: 600,
              color: '#1f2937',
            }}>
              {t('projectInfo.totalStatistics')}
            </span>
            <span style={{
              fontSize: '12px',
              color: 'rgba(0,0,0,0.45)',
              marginLeft: '4px',
            }}>
              ({t('common.total')} {totalStats.total} {t('projectInfo.cases')})
            </span>
          </div>

          {/* 中间：测试结果统计 */}
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '24px',
          }}>
            {/* OK */}
            <div style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              minWidth: '50px',
            }}>
              <span style={{ fontSize: '11px', color: 'rgba(0,0,0,0.45)', marginBottom: '2px' }}>OK</span>
              <span style={{
                fontSize: '18px',
                fontWeight: 700,
                color: '#52c41a',
                lineHeight: 1,
              }}>
                {totalStats.ok}
              </span>
            </div>

            {/* NG */}
            <div style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              minWidth: '50px',
            }}>
              <span style={{ fontSize: '11px', color: 'rgba(0,0,0,0.45)', marginBottom: '2px' }}>NG</span>
              <span style={{
                fontSize: '18px',
                fontWeight: 700,
                color: '#ff4d4f',
                lineHeight: 1,
              }}>
                {totalStats.ng}
              </span>
            </div>

            {/* Block */}
            <div style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              minWidth: '50px',
            }}>
              <span style={{ fontSize: '11px', color: 'rgba(0,0,0,0.45)', marginBottom: '2px' }}>Block</span>
              <span style={{
                fontSize: '18px',
                fontWeight: 700,
                color: '#faad14',
                lineHeight: 1,
              }}>
                {totalStats.block}
              </span>
            </div>

            {/* NR */}
            <div style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              minWidth: '50px',
            }}>
              <span style={{ fontSize: '11px', color: 'rgba(0,0,0,0.45)', marginBottom: '2px' }}>NR</span>
              <span style={{
                fontSize: '18px',
                fontWeight: 700,
                color: '#8c8c8c',
                lineHeight: 1,
              }}>
                {totalStats.nr}
              </span>
            </div>
          </div>

          {/* 右侧：进度和通过率 */}
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '20px',
          }}>
            {/* 执行进度 */}
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
            }}>
              <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)' }}>
                {t('testExecution.statistics.progress')}
              </span>
              <Progress
                percent={totalStats.progress}
                size="small"
                style={{ width: 80, margin: 0 }}
                strokeColor="#1890ff"
                showInfo={false}
              />
              <span style={{
                fontSize: '13px',
                fontWeight: 600,
                color: '#1890ff',
                minWidth: '36px',
                textAlign: 'right',
              }}>
                {totalStats.progress}%
              </span>
            </div>

            {/* 通过率 */}
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
            }}>
              <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)' }}>
                {t('testExecution.statistics.passRate')}
              </span>
              <Progress
                percent={totalStats.passRate}
                size="small"
                style={{ width: 80, margin: 0 }}
                strokeColor={totalStats.passRate >= 80 ? '#52c41a' : totalStats.passRate >= 60 ? '#faad14' : '#ff4d4f'}
                showInfo={false}
              />
              <span style={{
                fontSize: '13px',
                fontWeight: 600,
                color: totalStats.passRate >= 80 ? '#52c41a' : totalStats.passRate >= 60 ? '#faad14' : '#ff4d4f',
                minWidth: '36px',
                textAlign: 'right',
              }}>
                {totalStats.passRate}%
              </span>
            </div>
          </div>
        </div>
      </div>
    );
  };

  // 渲染任务卡片
  const renderTaskCard = () => (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: '24px' }}>
          <span style={{ fontSize: '16px', fontWeight: 600 }}>{t('projectInfo.taskSection')}</span>
          <div style={{ display: 'flex', alignItems: 'center', gap: '20px', fontSize: '13px' }}>
            <span>
              <span style={{ color: 'rgba(0,0,0,0.45)' }}>{t('testExecution.taskList.inProgress')}:</span>
              <span style={{ color: '#1890ff', fontWeight: 600, marginLeft: '6px' }}>{inProgressCount}</span>
            </span>
            <span>
              <span style={{ color: 'rgba(0,0,0,0.45)' }}>{t('testExecution.taskList.pending')}:</span>
              <span style={{ color: '#faad14', fontWeight: 600, marginLeft: '6px' }}>{pendingCount}</span>
            </span>
            <span>
              <span style={{ color: 'rgba(0,0,0,0.45)' }}>{t('testExecution.taskList.completed')}:</span>
              <span style={{ color: '#52c41a', fontWeight: 600, marginLeft: '6px' }}>{completedCount}</span>
            </span>
          </div>
        </div>
      }
      style={{ marginBottom: 16 }}
      bodyStyle={{ padding: '20px 24px' }}
    >
      {taskLoading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: 24 }}>
          <Spin />
        </div>
      ) : (
        <>
          {/* 总统计区域 */}
          {renderTotalStatistics()}

          {/* 进行中任务列表 */}
          {inProgressTasks.length === 0 ? (
            <Empty description={t('projectInfo.noInProgressTasks')} style={{ padding: '24px 0' }} />
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {inProgressTasks.map(task => {
                const stats = statsMap[task.task_uuid];
                const executionProgress = calculateExecutionProgress(stats);
                const passRate = calculatePassRate(stats);

                return (
                  <div
                    key={task.task_uuid}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: 10,
                      padding: '10px 14px',
                      background: '#fff',
                      border: '1px solid #f0f0f0',
                      borderRadius: 4,
                      fontSize: '13px',
                    }}
                  >
                    {/* 任务名称 */}
                    <Text
                      strong
                      style={{
                        width: 120,
                        flexShrink: 0,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        fontSize: '13px'
                      }}
                      title={task.task_name}
                    >
                      {task.task_name}
                    </Text>
                    <Tag color="processing" style={{ flexShrink: 0, margin: 0, fontSize: '12px' }}>
                      {t('testExecution.taskList.inProgress')}
                    </Tag>

                    {stats ? (
                      <>
                        {/* 执行结果统计 */}
                        <span style={{ whiteSpace: 'nowrap', flexShrink: 0, fontSize: '12px' }}>
                          <span style={{ color: 'rgba(0,0,0,0.45)' }}>OK:</span> <Text style={{ color: '#52c41a', fontWeight: 600 }}>{stats.ok_count || 0}</Text>
                        </span>
                        <span style={{ whiteSpace: 'nowrap', flexShrink: 0, fontSize: '12px' }}>
                          <span style={{ color: 'rgba(0,0,0,0.45)' }}>NG:</span> <Text style={{ color: '#ff4d4f', fontWeight: 600 }}>{stats.ng_count || 0}</Text>
                        </span>
                        <span style={{ whiteSpace: 'nowrap', flexShrink: 0, fontSize: '12px' }}>
                          <span style={{ color: 'rgba(0,0,0,0.45)' }}>Block:</span> <Text style={{ color: '#faad14', fontWeight: 600 }}>{stats.block_count || 0}</Text>
                        </span>
                        <span style={{ whiteSpace: 'nowrap', flexShrink: 0, fontSize: '12px', color: 'rgba(0,0,0,0.65)' }}>
                          <span style={{ color: 'rgba(0,0,0,0.45)' }}>NR:</span> {stats.nr_count || 0}
                        </span>

                        {/* 进度条 */}
                        <div style={{ display: 'flex', alignItems: 'center', gap: 6, flexShrink: 0 }}>
                          <span style={{ whiteSpace: 'nowrap', fontSize: '12px', color: 'rgba(0,0,0,0.45)' }}>{t('testExecution.statistics.progress')}</span>
                          <Progress percent={executionProgress} size="small" style={{ width: 60, margin: 0 }} showInfo={false} />
                          <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', width: '32px', textAlign: 'right' }}>{executionProgress}%</span>
                        </div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 6, flexShrink: 0 }}>
                          <span style={{ whiteSpace: 'nowrap', fontSize: '12px', color: 'rgba(0,0,0,0.45)' }}>{t('testExecution.statistics.passRate')}</span>
                          <Progress
                            percent={passRate}
                            size="small"
                            style={{ width: 60, margin: 0 }}
                            strokeColor={passRate >= 80 ? '#52c41a' : passRate >= 60 ? '#faad14' : '#ff4d4f'}
                            showInfo={false}
                          />
                          <span style={{ fontSize: '12px', color: 'rgba(0,0,0,0.65)', width: '32px', textAlign: 'right' }}>{passRate}%</span>
                        </div>
                      </>
                    ) : (
                      <Text type="secondary" style={{ fontSize: '12px' }}>-</Text>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </>
      )}
    </Card>
  );

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: 48 }}>
        <Spin size="large" />
      </div>
    );
  }

  // 只读模式显示
  if (!editing) {
    return (
      <Fragment>
        <Card
          title={<span style={{ fontSize: '16px', fontWeight: 600 }}>{t('projectDetail.projectInfo')}</span>}
          extra={
            canEdit && (
              <Button type="primary" icon={<EditOutlined />} onClick={handleEdit} size="small">
                {t('common.edit')}
              </Button>
            )
          }
          style={{ marginBottom: 16 }}
          bodyStyle={{ padding: '20px 24px' }}
        >
          <div style={{ background: '#fafafa', padding: '16px', borderRadius: '4px' }}>
            <div style={{ display: 'grid', gridTemplateColumns: '100px 1fr 100px 1fr', gap: '12px 20px', marginBottom: '12px', fontSize: '13px' }}>
              <div style={{ fontWeight: 600, color: 'rgba(0, 0, 0, 0.65)' }}>{t('project.name')}</div>
              <div style={{ color: 'rgba(0, 0, 0, 0.85)' }}>{projectData.name}</div>
              <div style={{ fontWeight: 600, color: 'rgba(0, 0, 0, 0.65)' }}>{t('project.status')}</div>
              <div style={{ color: 'rgba(0, 0, 0, 0.85)' }}>{getStatusText(projectData.status)}</div>

              <div style={{ fontWeight: 600, color: 'rgba(0, 0, 0, 0.65)' }}>{t('project.owner')}</div>
              <div style={{ color: 'rgba(0, 0, 0, 0.85)' }}>{getOwnerName()}</div>
              <div style={{ fontWeight: 600, color: 'rgba(0, 0, 0, 0.65)' }}>{t('project.createdAt')}</div>
              <div style={{ color: 'rgba(0, 0, 0, 0.85)' }}>{formatDate(projectData.created_at)}</div>
            </div>

            <div style={{ borderTop: '1px solid #e8e8e8', paddingTop: '12px', marginTop: '4px' }}>
              <div style={{ display: 'grid', gridTemplateColumns: '100px 1fr', gap: '12px 20px', fontSize: '13px' }}>
                <div style={{ fontWeight: 600, color: 'rgba(0, 0, 0, 0.65)' }}>{t('project.introduction')}</div>
                <div style={{ whiteSpace: 'pre-wrap', color: 'rgba(0, 0, 0, 0.85)', lineHeight: '1.6' }}>{projectData.description || '-'}</div>
              </div>
            </div>
          </div>
        </Card>

        {renderTaskCard()}

        {/* 缺陷趋势图 */}
        <DefectTrendChart projectId={projectId} />
      </Fragment>
    );
  }

  // 编辑模式
  return (
    <Fragment>
      <Card
        title={t('projectDetail.projectInfo')}
        extra={
          <Space>
            <Button icon={<CloseOutlined />} onClick={handleCancel}>
              {t('common.cancel')}
            </Button>
            <Button type="primary" icon={<SaveOutlined />} onClick={handleSave} loading={saving}>
              {t('common.save')}
            </Button>
          </Space>
        }
      >
        <Form
          form={form}
          layout="vertical"
        >
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0 16px' }}>
            <Form.Item
              label={t('project.name')}
              name="name"
              rules={[
                { required: true, message: t('project.nameRequired') },
                { max: 100, message: t('project.nameTooLong') },
              ]}
            >
              <Input placeholder={t('project.namePlaceholder')} disabled />
            </Form.Item>

            <Form.Item
              label={t('project.status')}
              name="status"
            >
              <Select placeholder={t('project.status')}>
                <Option value="pending">{t('project.statusPending')}</Option>
                <Option value="in-progress">{t('project.statusInProgress')}</Option>
                <Option value="completed">{t('project.statusCompleted')}</Option>
              </Select>
            </Form.Item>

            <Form.Item
              label={t('project.owner')}
            >
              <Input value={ownerName || '-'} disabled />
            </Form.Item>

            <Form.Item label={t('project.createdAt')}>
              <Input value={formatDate(projectData.created_at)} disabled />
            </Form.Item>
          </div>

          <Form.Item
            label={t('project.introduction')}
            name="description"
            rules={[
              { max: 500, message: t('project.descTooLong') },
            ]}
          >
            <TextArea
              rows={4}
              placeholder={t('project.descPlaceholder')}
            />
          </Form.Item>
        </Form>
      </Card>
      {renderTaskCard()}
    </Fragment>
  );
};

export default ProjectInfoTab;
