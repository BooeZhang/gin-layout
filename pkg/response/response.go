package response

import (
	"errors"
	"fmt"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Ok 通用响应
func Ok(c *gin.Context, err error, data interface{}) {
	code, httpCode, msg := erroron.DecodeErr(err)
	r := Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	c.AbortWithStatusJSON(httpCode, r)
}

type validationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
}

func Error(c *gin.Context, err error, data interface{}) {
	var validatorErrs validator.ValidationErrors
	if errors.As(err, &validatorErrs) {
		errs := make([]validationError, 0, len(validatorErrs))
		for i := range validatorErrs {
			validationErr := validatorErrs[i]
			errs = append(errs, validationError{
				ActualTag: validationErr.ActualTag(),
				Namespace: validationErr.Namespace(),
				Kind:      validationErr.Kind().String(),
				Type:      validationErr.Type().String(),
				Value:     fmt.Sprintf("%v", validationErr.Value()),
				Param:     validationErr.Param(),
			})
		}
		log.Errorf("%+v", errs)
		Ok(c, erroron.ErrParameter, nil)
		return
	}
	code, httpCode, msg := erroron.DecodeErr(err)
	c.AbortWithStatusJSON(httpCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
