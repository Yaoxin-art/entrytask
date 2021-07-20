package main

import (
	"encoding/gob"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
)

func configRPC() {
	client.Config(facade.Logon, rpcServerAddr, &facade.UserLogon)
	client.Config(facade.Login, rpcServerAddr, &facade.UserLogin)
	client.Config(facade.Query, rpcServerAddr, &facade.UserQuery)
	client.Config(facade.QueryToken, rpcServerAddr, &facade.UserQueryByToken)
	client.Config(facade.UpdateNick, rpcServerAddr, &facade.UserUpdateNick)
	client.Config(facade.UpdateProfile, rpcServerAddr, &facade.UserUpdateProfile)

	registerPojo()
}

func registerPojo() {
	gob.Register(&facade.User{})
	gob.Register(&facade.UserUpdateRequest{})
	gob.Register(&facade.UserLogonRequest{})
	gob.Register(&facade.UserLoginRequest{})
}
