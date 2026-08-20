package web

import (
	"errors"
	"net/http"
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

type errorMapping struct {
	target     error
	descriptor ErrorDescriptor
}

func newErrorMapping(target error, httpStatus, code int, message string) errorMapping {
	return errorMapping{
		target: target,
		descriptor: ErrorDescriptor{
			HTTPStatus: httpStatus,
			Code:       code,
			Message:    message,
		},
	}
}

func matchError(err error, mappings []errorMapping) (ErrorDescriptor, bool) {
	for _, mapping := range mappings {
		if errors.Is(err, mapping.target) {
			return mapping.descriptor, true
		}
	}
	return ErrorDescriptor{}, false
}

func DecodeError(err error) ErrorDescriptor {
	if err == nil {
		return successDescriptor
	}

	for _, mappings := range allErrorMappings {
		if descriptor, ok := matchError(err, mappings); ok {
			return descriptor
		}
	}

	return internalServerError
}
