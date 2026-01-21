/**
 * 脚本测试 API
 * 用于在用例详情页面直接测试脚本
 */
import client from './client';

/**
 * 测试脚本
 * @param {number} projectId - 项目ID
 * @param {Object} params - 测试参数
 * @param {string} params.script_code - 脚本代码
 * @param {number} [params.group_id] - 用例集ID（用于获取变量）
 * @param {string} [params.group_type] - 用例集类型 ('web' | 'api')
 * @returns {Promise<Object>} 测试结果
 * {
 *   success: boolean,
 *   output: string,
 *   error_message?: string,
 *   response_time: number,
 *   executed_at: string
 * }
 */
export const testScript = async (projectId, params) => {
    // client 拦截器已经返回 response.data，所以这里直接返回即可
    // 设置120秒超时，因为Web脚本测试可能需要较长时间
    const result = await client.post(`/projects/${projectId}/script-test`, params, {
        timeout: 120000, // 120秒超时
    });
    return result;
};

export default {
    testScript,
};
