package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/ginHelper"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func main() {
	gin.SetMode(gin.DebugMode)
	s := Routers()
	s.Run(":22222")
}

func Routers() *gin.Engine {
	Router.GET("/", ginHelper.Handle(Home))
	return Router
}

func Home(c *ginHelper.GinCtx) {
	c.APIOutPut(0, "ok", "")
}
