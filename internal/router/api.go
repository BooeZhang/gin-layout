package router

import (
	"github.com/gin-gonic/gin"
)

type api struct {
}

var Api = new(api)

func (*api) Load(g *gin.Engine) {

}
