package erroron

var (
	OK                = &Errno{Code: 200, HttpStatus: 200, Msg: "OK"}
	ErrNotFound       = &Errno{Code: 404, HttpStatus: 404, Msg: "Page not found"}
	ErrInternalServer = &Errno{Code: 500, HttpStatus: 500, Msg: "服务器内部错误"}
	ErrNoPerm         = &Errno{Code: 403, HttpStatus: 403, Msg: "无访问权限"}
	ErrParameter      = &Errno{Code: 400, HttpStatus: 400, Msg: "请求参数无效"}
)
