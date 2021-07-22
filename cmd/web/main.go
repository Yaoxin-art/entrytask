package main

import (
	"git.garena.com/zhenrong.zeng/entrytask/cmd/web/router"
	"git.garena.com/zhenrong.zeng/entrytask/pkg/zerorpc"
	_ "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path"
	"runtime"
	"time"
)

func initLog() {
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, path.Base(frame.File)
		},
	})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logrus.SetOutput(os.Stdout)
	//设置最低loglevel
	logrus.SetLevel(logrus.InfoLevel)
	//
	logrus.SetReportCaller(true)
}

var (
	httpPort      = ":7777"
	rpcServerAddr = "127.0.0.1:9999"
)
var client *zerorpc.GrettyClient

func main() {
	initLog()

	client = zerorpc.NewGrettyClient(rpcServerAddr, 50, 100) // init rpc client instance
	configRPC()                                              // config for rpc client

	httpServer := router.InitGin() // init http server with gin
	go func() {
		err := httpServer.Run(httpPort)
		if err != nil {
			logrus.Warnf("Gin run err:%v", err)
		}
	}() // listen and serve on 0.0.0.0:7777

	logrus.Infof("web app started, listen at:%s", httpPort)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		<-c
		logrus.Info("web app graceful shutdown...")
		return
	}
}
