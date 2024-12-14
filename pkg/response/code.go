package response

// 系统级错误码
const (
	SUCCESS = 200 // 成功码

	// 客户端错误码 (400-499)
	InvalidParams    = 400
	Unauthorized     = 401
	Forbidden        = 403
	NotFound         = 404
	MethodNotAllowed = 405

	// 服务端错误码 (500-599)
	ServerError        = 500
	ServiceUnavailable = 503

	// 自定义错误码 (1000-1999)
	TokenExpired     = 1001 // Token已过期
	TokenNotValidYet = 1002 // Token尚未生效
	InvalidToken     = 1003 // 无效的Token
	UserNotFound     = 1004 // 用户不存在
	PermissionDenied = 1005 // 权限不足
)

// ErrorCode 错误码结构
type ErrorCode struct {
	Code    int    // 错误码
	Message string // 错误信息
}

// 错误码映射
var systemMessages = map[int]string{
	SUCCESS:            "操作成功",
	InvalidParams:      "请求参数错误",
	Unauthorized:       "未授权",
	Forbidden:          "无权限",
	NotFound:           "资源不存在",
	MethodNotAllowed:   "方法不允许",
	ServerError:        "服务器内部错误",
	ServiceUnavailable: "服务不可用",

	// 自定义错误码
	TokenExpired:     "Token已过期",
	TokenNotValidYet: "Token尚未生效",
	InvalidToken:     "无效的Token",
	UserNotFound:     "用户不存在",
	PermissionDenied: "权限不足",
}

// 自定义错误码映射
var customMessages = make(map[int]string)

// RegisterCode 注册自定义错误码
func RegisterCode(codes ...ErrorCode) {
	for _, code := range codes {
		customMessages[code.Code] = code.Message
	}
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	// 优先查找自定义错误码
	if msg, ok := customMessages[code]; ok {
		return msg
	}
	// 然后查找系统错误码
	if msg, ok := systemMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
