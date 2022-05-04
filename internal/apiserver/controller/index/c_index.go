package index

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/service/index"
	"github.com/BooeZhang/gin-layout/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	srv index.IIndex
}

// NewIndexController index控制器
func NewIndexController(store store.Factory, cache store.Cache) *Controller {
	return &Controller{
		srv: index.NewIndexService(store, cache),
	}
}

func (i *Controller) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (i *Controller) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
