package main

import (
	"git.garena.com/zhenrong.zeng/entrytask/internal/web/router"
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func init() {
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	log.Formatter = &logrus.TextFormatter{}
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	log.SetOutput(os.Stdout)
	//设置最低loglevel
	log.Level = logrus.DebugLevel
}

func main() {
	log.Info("web app starting...")
	router.InitGin()
	log.Info("web app started...")
}
