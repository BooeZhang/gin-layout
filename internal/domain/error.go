package domain

import "errors"

// DomainError 是具有业务错误码和 HTTP 状态码的业务级错误。
type DomainError struct {
	Code       int    // 业务错误编号
	HTTPStatus int    // HTTP 状态码（如 400, 401, 403, 404, 409, 422, 500）
	Message    string // 详细错误消息
}

func (e *DomainError) Error() string { return e.Message }

func (e *DomainError) ResponseCode() int { return e.Code }

func (e *DomainError) ResponseHTTPStatus() int { return e.HTTPStatus }

func (e *DomainError) ResponseMessage() string { return e.Message }

func (e *DomainError) Is(target error) bool {
	var t *DomainError
	if !errors.As(target, &t) {
		return false
	}
	return e.Code == t.Code
}

func NewDomainError(code int, httpStatus int, message string) *DomainError {
	return &DomainError{Code: code, HTTPStatus: httpStatus, Message: message}
}
