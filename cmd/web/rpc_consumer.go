package main

import (
	"encoding/gob"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
)

func configRPC() {
	client.ConfigRemoteMethod(facade.Logon, &facade.UserLogon)
	client.ConfigRemoteMethod(facade.Login, &facade.UserLogin)
	client.ConfigRemoteMethod(facade.Query, &facade.UserQuery)
	client.ConfigRemoteMethod(facade.QueryToken, &facade.UserQueryByToken)
	client.ConfigRemoteMethod(facade.UpdateNick, &facade.UserUpdateNick)
	client.ConfigRemoteMethod(facade.UpdateProfile, &facade.UserUpdateProfile)

	client.ConfigRemoteMethod(facade.QueryList, &facade.QueryUsernameList)

	registerPojo()
}

func registerPojo() {
	gob.Register(&facade.User{})
	gob.Register(&facade.UserUpdateRequest{})
	gob.Register(&facade.UserLogonRequest{})
	gob.Register(&facade.UserLoginRequest{})
}
