package response

import (
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type pages struct {
	Count    int64       `json:"count"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     interface{} `json:"list"`
}

// Ok 通用响应
func Ok(c *gin.Context, err error, data interface{}) {
	code, httpStatus, msg := erroron.DecodeErr(err)
	r := Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	c.AbortWithStatusJSON(httpStatus, r)
}

// PageOk 列表响应
func PageOk(c *gin.Context, err error, data interface{}, count int64, page, pageSize int) {
	p := pages{
		Count:    count,
		Page:     page,
		PageSize: pageSize,
		List:     data,
	}
	Ok(c, err, p)
}

func Error(c *gin.Context, err error, data interface{}) {
	code, _, msg := erroron.DecodeErr(err)
	if code != 0 && code < 10000 {
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
	} else {
		c.JSON(code, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
	}
	c.Abort()
}
