package router

import (
	"github.com/gin-gonic/gin"
)

//initRouter
func initRouter(engine *gin.Engine) {
	engine.GET("/user/list", query)
	engine.GET("/ping", ping)
}
