package erroron

import (
	"errors"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"net/http"
)

// Errno 错误定义
type Errno struct {
	Code int
	Msg  string
}

var codes = map[int]string{}

func NewError(code int, msg string) *Errno {
	if _, ok := codes[code]; ok {
		log.Panicf("错误码 %d 已经存在，请更换一个", code)
	}

	codes[code] = msg
	return &Errno{
		Code: code,
		Msg:  msg,
	}
}

// Error 错误字符串返回
func (err Errno) Error() string {
	return err.Msg
}

// StatusCode 特殊 http statusCode 处理
func (err Errno) StatusCode() int {
	switch err.Code {
	case OK.Code:
		return http.StatusOK
	default:
		return http.StatusOK
	}

}

// DecodeErr 解析错误信息
func DecodeErr(err error) (int, int, string) {
	if err == nil {
		return OK.Code, http.StatusOK, OK.Msg
	}

	var _errorOn *Errno
	if errors.As(err, _errorOn) {
		httpCode := _errorOn.StatusCode()
		return _errorOn.Code, httpCode, _errorOn.Msg
	} else {
		log.Errorf("[error]--> %s", err)
		return 500, 200, "服务器内部错误"
	}
}
