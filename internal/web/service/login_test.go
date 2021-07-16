package service

import (
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
	"testing"
)

func TestLogin(t *testing.T) {
	// todo
	request := facade.UserLoginRequest{Username: "zero1234", Password: "123456"}
	code := Login(request)
	if code != 1 {
		t.Errorf("login failure")
	}
	t.Logf("login success")
}
