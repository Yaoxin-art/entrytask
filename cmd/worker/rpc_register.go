package main

import (
	"encoding/gob"
	"git.garena.com/zhenrong.zeng/entrytask/cmd/worker/service"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
)

func registerRPC() {
	server.Register(facade.Logon, service.Logon)
	server.Register(facade.Login, service.Login)
	server.Register(facade.Query, service.QueryByUsername)
	server.Register(facade.QueryToken, service.QueryByToken)
	server.Register(facade.UpdateNick, service.UpdateUserNick)
	server.Register(facade.UpdateProfile, service.UpdateUserProfile)

	registerPojo()
}

func registerPojo() {
	gob.Register(&facade.User{})
	gob.Register(&facade.UserUpdateRequest{})
	gob.Register(&facade.UserLogonRequest{})
	gob.Register(&facade.UserLoginRequest{})
}
