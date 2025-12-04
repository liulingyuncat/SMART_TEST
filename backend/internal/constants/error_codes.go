package constants

// 错误码定义
const (
	// 客户端错误 (40xxx)
	ErrInvalidInput         = 40000 // 无效输入参数
	ErrUnauthorized         = 40001 // 无权限访问项目
	ErrInvalidFileFormat    = 40002 // 文件格式不支持
	ErrFileSizeExceeded     = 40003 // 文件大小超限
	ErrExcelParseFailed     = 40004 // Excel解析失败
	ErrMissingRequiredField = 40005 // 必填字段缺失
	ErrVersionNotFound      = 40006 // 版本不存在
	ErrFileRequired         = 40007 // 文件未提供
	ErrFileNotFound         = 40008 // 文件不存在

	// 服务端错误 (50xxx)
	ErrExportFailed        = 50000 // 导出失败
	ErrExcelGenerateFailed = 50001 // Excel生成失败
	ErrFileStoreFailed     = 50002 // 文件存储失败
	ErrImportFailed        = 50003 // 导入失败
)

// 错误消息映射
var ErrorMessages = map[int]string{
	ErrInvalidInput:         "无效输入参数",
	ErrUnauthorized:         "无权限访问项目",
	ErrInvalidFileFormat:    "文件格式不支持",
	ErrFileSizeExceeded:     "文件大小超限",
	ErrExcelParseFailed:     "Excel解析失败",
	ErrMissingRequiredField: "必填字段缺失",
	ErrVersionNotFound:      "版本不存在",
	ErrFileRequired:         "文件未提供",
	ErrFileNotFound:         "文件不存在",
	ErrExportFailed:         "导出失败",
	ErrExcelGenerateFailed:  "Excel生成失败",
	ErrFileStoreFailed:      "文件存储失败",
	ErrImportFailed:         "导入失败",
}

// GetErrorMessage 获取错误消息
func GetErrorMessage(code int) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
