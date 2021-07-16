package main

import (
	"flag"
	"git.garena.com/zhenrong.zeng/entrytask/internal/web/router"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
)

const (
	DefaultPort = ":7777"
)

var (
	port       string
	env        string
	HttpEngine *gin.Engine
)

func init() {
	// config
	flag.StringVar(&port, "port", DefaultPort, "Http port, like ':7777'")
	flag.Parse()
	validConf()

	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	logrus.SetFormatter(&logrus.TextFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logrus.SetOutput(os.Stdout)
	//设置最低loglevel
	logrus.SetLevel(logrus.InfoLevel)
}

func validConf() {
	matchPort, _ := regexp.MatchString("(^\\:\\d+)", port)
	if !matchPort {
		port = DefaultPort
		logrus.Warnf("Invalid port config, set with default: %s", DefaultPort)
	}
}

func main() {
	httpServer := router.InitGin()
	logrus.Infof("web app started, listen at:%s", port)
	// listen and serve on 0.0.0.0:7777 (default "localhost:7777")
	httpServer.Run(port)
}
