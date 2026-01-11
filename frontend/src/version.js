// 版本号管理
// 生产环境：从构建时注入的环境变量读取
// 本地开发：使用 package.json 的版本或默认值

export const VERSION = process.env.REACT_APP_VERSION || 'dev';

export default VERSION;
