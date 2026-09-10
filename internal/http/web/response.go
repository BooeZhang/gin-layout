package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

func Error(c *gin.Context, err error) {
	descriptor := DecodeError(err)
	c.JSON(descriptor.HTTPStatus, Response{Code: descriptor.Code, Message: descriptor.Message})
}
