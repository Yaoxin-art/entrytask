package main

import "git.garena.com/zhenrong.zeng/entrytask/cmd/worker/service"

func registerRPC() {
	server.Register("Logon", service.Logon)
	server.Register("Login", service.Login)
	server.Register("QueryUser", service.QueryByUsername)
	server.Register("UpdateNick", service.UpdateUserNick)
	server.Register("UpdateProfile", service.UpdateUserProfile)
}
