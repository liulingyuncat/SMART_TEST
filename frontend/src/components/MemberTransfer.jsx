import React from 'react';
import { Transfer } from 'antd';
import PropTypes from 'prop-types';

/**
 * 成员穿梭框组件
 * @param {Array} dataSource - 所有该角色用户数组 [{user_id, username, nickname, role}]
 * @param {Array} targetKeys - 已分配用户ID数组 [user_id1, user_id2, ...]
 * @param {Array} lockedKeys - 锁定不可选择的用户ID数组 [user_id1, ...]
 * @param {Function} onChange - 回调函数 (newTargetKeys) => void
 * @param {Array} title - 穿梭框标题 [leftTitle, rightTitle]
 */
const MemberTransfer = ({ dataSource, targetKeys, lockedKeys, onChange, title }) => {
  // 添加调试日志
  console.log('[MemberTransfer] dataSource:', dataSource);
  console.log('[MemberTransfer] targetKeys:', targetKeys);
  console.log('[MemberTransfer] lockedKeys:', lockedKeys);

  // 转换数据源为Transfer要求的格式
  const transferDataSource = dataSource.map((user) => {
    const isLocked = lockedKeys.includes(user.user_id);
    console.log(`[MemberTransfer] 用户 ${user.username} (${user.user_id}) 是否锁定: ${isLocked}`);
    return {
      key: String(user.user_id), // Transfer要求key为字符串
      title: `${user.nickname} (${user.username})`, // 显示昵称和用户名
      disabled: isLocked, // 根据lockedKeys判断是否禁用
    };
  });

  // 处理穿梭框变化
  const handleChange = (newTargetKeys) => {
    // 将字符串key转换回数字ID
    const numericKeys = newTargetKeys.map((key) => parseInt(key, 10));
    
    // 确保锁定的用户始终保留在targetKeys中
    const lockedNumericKeys = lockedKeys.filter(id => targetKeys.includes(id));
    const finalKeys = [...new Set([...numericKeys, ...lockedNumericKeys])];
    
    console.log('[MemberTransfer] 原始newTargetKeys:', newTargetKeys);
    console.log('[MemberTransfer] 转换后numericKeys:', numericKeys);
    console.log('[MemberTransfer] 锁定的Keys:', lockedNumericKeys);
    console.log('[MemberTransfer] 最终Keys:', finalKeys);
    
    onChange(finalKeys);
  };

  // 转换targetKeys为字符串数组
  const stringTargetKeys = targetKeys.map((id) => String(id));

  return (
    <Transfer
      dataSource={transferDataSource}
      targetKeys={stringTargetKeys}
      onChange={handleChange}
      titles={title}
      render={(item) => item.title}
      listStyle={{
        width: 250,
        height: 400,
      }}
      showSearch
      filterOption={(inputValue, option) =>
        option.title.toLowerCase().includes(inputValue.toLowerCase())
      }
    />
  );
};

MemberTransfer.propTypes = {
  dataSource: PropTypes.arrayOf(
    PropTypes.shape({
      user_id: PropTypes.number.isRequired,
      username: PropTypes.string.isRequired,
      nickname: PropTypes.string.isRequired,
      role: PropTypes.string,
    })
  ).isRequired,
  targetKeys: PropTypes.arrayOf(PropTypes.number).isRequired,
  lockedKeys: PropTypes.arrayOf(PropTypes.number),
  onChange: PropTypes.func.isRequired,
  title: PropTypes.arrayOf(PropTypes.string),
};

MemberTransfer.defaultProps = {
  lockedKeys: [],
  title: ['可选用户', '已选用户'],
};

export default MemberTransfer;
