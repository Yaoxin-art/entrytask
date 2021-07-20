package router

import (
	"github.com/gin-gonic/gin"
)

//initRouter
func initRouter(engine *gin.Engine) {
	engine.POST("/user/login", login)
	engine.POST("/user/logon", logon)
	engine.POST("/user/logout", logout)
	engine.GET("/user/info", info)
	engine.GET("/user/find", findByUsername)

	engine.POST("/user/profile", profileUpdate)
	engine.GET("/user/nick", nickUpdate)
	engine.GET("/ping", ping)
}
