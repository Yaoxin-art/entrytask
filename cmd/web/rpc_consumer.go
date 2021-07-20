package main

import "git.garena.com/zhenrong.zeng/entrytask/internal/facade"


func configRPC() {
	client.Config(facade.Logon, rpcServerAddr, &facade.UserLogon)
	client.Config(facade.Login, rpcServerAddr, &facade.UserLogin)
	client.Config(facade.Query, rpcServerAddr, &facade.UserQuery)
	client.Config(facade.UpdateNick, rpcServerAddr, &facade.UserUpdateNick)
	client.Config(facade.UpdateProfile, rpcServerAddr, &facade.UserUpdateProfile)
}