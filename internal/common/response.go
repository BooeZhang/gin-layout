package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type CodedError interface {
	error
	ResponseCode() int
	ResponseHTTPStatus() int
	ResponseMessage() string
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

func Error(c *gin.Context, err error) {
	httpStatus, code, message := DecodeError(err)
	c.JSON(httpStatus, Response{Code: code, Message: message})
}

func DecodeError(err error) (httpStatus int, bizCode int, message string) {
	if err == nil {
		return http.StatusOK, 0, "success"
	}

	if codedErr, ok := errors.AsType[CodedError](err); ok {
		return codedErr.ResponseHTTPStatus(), codedErr.ResponseCode(), codedErr.ResponseMessage()
	}

	return http.StatusInternalServerError, 50001, "internal server error"
}
