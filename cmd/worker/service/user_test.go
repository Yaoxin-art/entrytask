package service

import (
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// TestQueryUserByUsername 测试根据username查询用户信息
// case 1：存在的用户查询
func TestQueryUserByUsername1(t *testing.T) {
	username := "zero1234"
	user, err := QueryByUsername(username)
	if err != 0 {
		t.Errorf("query user:%s empty, err:%v", username, err)
	}
	t.Logf("query user:%s success, user:%v", username, user)
}

// TestQueryUserByUsername 测试根据username查询用户信息
// case 2：不存在的用户查询
func TestQueryUserByUsername2(t *testing.T) {
	username := "zero.not.exist"
	_, err := QueryByUsername(username)
	if err == 0 {
		t.Errorf("query user:%s case:2 test failure", username)
	}
	t.Logf("query user:%s case:2 success", username)
}

var timestamp = time.Now().Unix()
var ops int64 = 10000000

// TestRegisterUser 测试用户注册
// case 1：正常信息注册
func TestRegisterUser1(t *testing.T) {
	atomic.AddInt64(&ops, 1)
	ctime := strconv.FormatInt(timestamp+ops, 16)
	// 随机生成username和nickname，固定前缀加随机自增数字
	username := "zero" + ctime
	nickname := "zero No." + ctime
	password := "123456"
	user := facade.UserLogonRequest{Username: username, Nickname: nickname, Password: password}
	name, err := Logon(&user)
	if err != 0 || name == "" {
		t.Errorf("register user:%v failure, name:%s", user, name)
	}
	t.Logf("register user:%v success", user)
}

// TestRegisterUser 测试用户注册
// case 2：重复注册（username重复）
func TestRegisterUser2(t *testing.T) {
	// 随机生成username和nickname，固定前缀加随机自增数字
	username := "zero1234"
	nickname := "zero No.1234"
	password := "123456"
	user := facade.UserLogonRequest{Username: username, Nickname: nickname, Password: password}
	_, err := Logon(&user)
	if err == 0 {
		t.Errorf("register user:%v case:2 test failure, code:%d", user, err)
	}
	if err != 1 {
		t.Errorf("register duplicate user:%v case test should throw \"Duplicate entry 'zero1234' for key 't_user.unique_idx_username'\", and code should be 2, but got:%d", user, err)
	}
	t.Logf("register user:%v case:2 success, msg:%v", user, err)
}

// todo
func TestLogin(t *testing.T) {
	request := facade.UserLoginRequest{Username: "zero1234", Password: "123456"}
	user, token, err := Login(&request)
	if err == 0 {
		t.Logf("login success, user:%v, token:%s", user, token)
	} else {
		t.Errorf("login failure, should err=0, but got:%d", err)
	}
}

// TestUpdateUserNick 测试更新用户昵称
// case 1：正常请求
func TestUpdateUserNick1(t *testing.T) {
	time.Sleep(2 * time.Second) // 睡眠2秒，用于update profile 测试用例时间差
	request := facade.UserUpdateRequest{Username: "zero1234", Nickname: "zero No.1234"}
	user, err := UpdateUserNick(&request)
	if err == 0 {
		t.Logf("update user nickname case:1 test success, updated user:%v", user)
	} else {
		t.Errorf("update user nickname case:1 test faliure, should err=0, but got:%d", err)
	}
}

// TestUpdateUserNick 测试更新用户昵称
// case 2：不存在的用户更新
func TestUpdateUserNick2(t *testing.T) {
	request := facade.UserUpdateRequest{Username: "zero.not.exist", Nickname: "zero No.1234"}
	_, err := UpdateUserNick(&request)
	if err == 2 {
		t.Log("update user nickname case:2 test success")
	} else {
		t.Errorf("update user nickname case:2 test failure, should out err=2 but got:%d", err)
	}
}

// TestUpdateUserNick 测试更新用户昵称
// case 3：参数不完整
func TestUpdateUserNick3(t *testing.T) {
	request := facade.UserUpdateRequest{Username: "zero1234"}
	_, err := UpdateUserNick(&request)
	if err == 1 {
		t.Log("update user nickname case:3 test success")
	} else {
		t.Errorf("update user nickname case:3 test failure, shold out err=1 but got:%d", err)
	}
}

// TestUpdateUserProfile 测试更新用户头像
// case 1：正常请求
func TestUpdateUserProfile1(t *testing.T) {
	time.Sleep(2 * time.Second) // 睡眠2秒，用于update nickname 测试用例时间差
	request := facade.UserUpdateRequest{Username: "zero1234", ProfilePath: "/profile/default.jpg"}
	user, err := UpdateUserProfile(&request)
	if err == 0 {
		t.Logf("update user profile case:1 test success, updated user:%v", user)
	} else {
		t.Errorf("update user profile case:1 test faliure, should err=0, but got:%d", err)
	}
}

// TestUpdateUserProfile 测试更新用户头像
// case 2：不存在的用户更新
func TestUpdateUserProfile2(t *testing.T) {
	request := facade.UserUpdateRequest{Username: "zero_not_exist", ProfilePath: "/profile/avatar.jpg"}
	_, err := UpdateUserProfile(&request)
	if err == 2 {
		t.Log("update user profile case:2 test success")
	} else {
		t.Errorf("update user profile case:2 test failure, should out err=2 but got:%d", err)
	}
}

// TestUpdateUserProfile 测试更新用户头像
// case 3：参数不完整
func TestUpdateUserProfile3(t *testing.T) {
	request := facade.UserUpdateRequest{Username: "zero1234", ProfilePath: ""}
	_, err := UpdateUserProfile(&request)
	if err == 1 {
		t.Log("update user profile case:3 test success")
	} else {
		t.Errorf("update user profile case:3 test failure, shold out err=1 but got:%d", err)
	}
}

func TestGenerateToken(t *testing.T) {
	t.Logf("generate token:%s", generateToken("zeng1234"))
}
