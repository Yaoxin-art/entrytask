package main

import (
	"git.garena.com/zhenrong.zeng/entrytask/pkg/zerorpc"
	logger "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path"
	"runtime"
	"time"
)

func initLog() {
	//
	logger.SetReportCaller(true)
	//设置输出样式，自带的只有两种样式logger.JSONFormatter{}和logger.TextFormatter{}
	logger.SetFormatter(&logger.TextFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, path.Base(frame.File)
		},
	})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logger.SetOutput(os.Stdout)
	//设置最低loglevel
	logger.SetLevel(logger.InfoLevel)
}

const serverAddr = ":9999"
var server *zerorpc.Server

func main() {
	initLog()
	server = zerorpc.NewServer(serverAddr)
	registerRPC() 		// register rpc service
	go server.Run()		// run rpc server

	logger.Infof("worker app started, rpc server listen on[%s]", serverAddr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		<-c
		server.Close()	// rpc server close while interrupt, or use "def server.Close()"
		logger.Info("work app graceful shutdown...")
		return
	}
}
