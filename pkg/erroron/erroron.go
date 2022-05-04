package erroron

import "net/http"

// Errno 错误定义
type Errno struct {
	HttpStatus int
	Code       int
	Msg        string
}

// Error 错误字符串返回
func (err Errno) Error() string {
	return err.Msg
}

// DecodeErr 解析错误信息
func DecodeErr(err error) (int, int, string) {
	if err == nil {
		return OK.Code, http.StatusOK, OK.Msg
	}

	switch typed := err.(type) {
	case *Errno:
		if typed.HttpStatus == 0 {
			typed.HttpStatus = http.StatusOK
		}
		return typed.Code, typed.HttpStatus, typed.Msg
	default:
		return 500, http.StatusInternalServerError, typed.Error()
	}
}
