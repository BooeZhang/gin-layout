package erroron

var (
	OK                = &Errno{Code: 200, HttpStatus: 200, Msg: "OK"}
	ErrNotFound       = &Errno{Code: 404, HttpStatus: 404, Msg: "Page not found"}
	ErrNoPerm         = &Errno{Code: 403, HttpStatus: 403, Msg: "无访问权限"}
	ErrParameter      = &Errno{Code: 400, HttpStatus: 400, Msg: "请求参数无效"}
	ErrInternalServer = &Errno{Code: 500, HttpStatus: 500, Msg: "服务器内部错误"}
)

var (
	ErrNotLogin       = &Errno{Code: 10000, HttpStatus: 200, Msg: "未登录或非法访问"}
	ErrTokenInvalid   = &Errno{Code: 8001, HttpStatus: 200, Msg: "token 无效"}
	ErrTokenExpired   = &Errno{Code: 8002, HttpStatus: 200, Msg: "token 过期"}
	ErrTokenNotActive = &Errno{Code: 8003, HttpStatus: 200, Msg: "token 不是活跃的"}
	ErrNotFoundUser   = &Errno{Code: 3, HttpStatus: 200, Msg: "未找到用户"}
	ErrUserNameOrPwd  = &Errno{Code: 4, HttpStatus: 200, Msg: "用户名或密码错误"}
)
