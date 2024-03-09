package erroron

var (
	OK                = NewError(200, "OK")
	ErrNotFound       = NewError(404, "Page not found")
	ErrInternalServer = NewError(500, "服务器内部错误")
	ErrNoPerm         = NewError(403, "无访问权限")
	ErrParameter      = NewError(400, "请求参数无效")
	ErrNotLogin       = NewError(10000, "未登录或非法访问")
	ErrTokenInvalid   = NewError(8001, "token 无效")
	ErrTokenExpired   = NewError(8002, "token 过期")
	ErrTokenNotActive = NewError(8003, "token 不是活跃的")
	ErrNotFoundUser   = NewError(8004, "未找到用户")
	ErrUserNameOrPwd  = NewError(8005, "用户名或密码错误")
)
