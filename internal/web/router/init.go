package router

import (
	"flag"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"regexp"
)

const (
	DefaultPort = ":7777"
)

var port = flag.String("httpPort", DefaultPort, "Http port, like ':7777'")
var env = flag.String("env", "debug", "Run mode, like '<debug|test|release'")
var HttpEngine *gin.Engine

func InitGin() {
	r := gin.Default()
	r.Use() // 全局handler
	validConf()
	gin.SetMode(*env) // run mode
	initRouter(r)
	log.Info("Gin init router finished...")
	r.Run(*port) // listen and serve on 0.0.0.0:7777 (default "localhost:7777")
	log.Info("Gin running...")
	HttpEngine = r
}

func validConf() {
	matchEnv, _ := regexp.Match("H(debug|test|release)", []byte(*env))
	if !matchEnv {
		*env = gin.DebugMode
		log.Warnf("Invalid env config, set with default: %s", gin.DebugMode)
	}
	matchPort, _ := regexp.Match("H(^\\:\\d+)", []byte(*port))
	if !matchPort {
		*port = DefaultPort
		log.Warnf("Invalid port config, set with default: %s", DefaultPort)
	}
}
