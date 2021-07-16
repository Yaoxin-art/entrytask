package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	HttpEngine *gin.Engine
)

func InitGin() *gin.Engine {
	//r.Use() // 全局handler
	gin.SetMode(gin.ReleaseMode) // run mode
	r := gin.New()
	initRouter(r)
	logrus.Info("Gin init router finished...")
	HttpEngine = r
	return r
}
