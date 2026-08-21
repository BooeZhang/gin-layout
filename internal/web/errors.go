package web

import (
	"errors"
	"net/http"

	"gin-layout/internal/apperror"
)

type ErrorDescriptor struct {
	HTTPStatus int
	Code       int
	Message    string
}

var (
	successDescriptor   = ErrorDescriptor{HTTPStatus: http.StatusOK, Code: 0, Message: "success"}
	internalServerError = ErrorDescriptor{
		HTTPStatus: http.StatusInternalServerError,
		Code:       50001,
		Message:    "internal server error",
	}
)

// DecodeError 将应用错误转换为统一的 Web 响应描述。
func DecodeError(err error) ErrorDescriptor {
	if err == nil {
		return successDescriptor
	}

	var appErr *apperror.Error
	if errors.As(err, &appErr) && appErr != nil {
		return ErrorDescriptor{
			HTTPStatus: httpStatusFor(appErr.Kind()),
			Code:       appErr.Code(),
			Message:    appErr.Message(),
		}
	}

	return internalServerError
}

// 将应用错误分类映射为 HTTP 状态码。
func httpStatusFor(kind apperror.Kind) int {
	switch kind {
	case apperror.InvalidInput:
		return http.StatusUnprocessableEntity
	case apperror.NotFound:
		return http.StatusNotFound
	case apperror.Conflict:
		return http.StatusConflict
	case apperror.Unauthenticated:
		return http.StatusUnauthorized
	case apperror.Forbidden:
		return http.StatusForbidden
	case apperror.BusinessResult:
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
