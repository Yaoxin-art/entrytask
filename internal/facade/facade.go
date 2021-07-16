package facade

import "fmt"

// User
// 返回用户信息对象
type User struct {
	Id          int64
	Username    string
	Nickname    string
	ProfilePath string
}

// UserLogonRequest
// 用户注册时的请求对象
type UserLogonRequest struct {
	Username string
	Nickname string
	Password string
}

// UserLoginRequest
// 用户登录时的请求对象
type UserLoginRequest struct {
	Username string
	Password string
}

// UserUpdateRequest
// 用户更新昵称和头像时的请求对象
// 更新昵称时，Username和Nickname不可为空
// 更新头像时，Username和ProfilePath不可为空
type UserUpdateRequest struct {
	Username    string
	Nickname    string
	ProfilePath string
}

// User.String
func (u User) String() string {
	return fmt.Sprintf("Id:%d, Username:%s, Nickname:%s, Profile:%s", u.Id, u.Username, u.Nickname, u.ProfilePath)
}

// UserLogon 用户注册暴露方法
// return code 返回状态码：1-成功，2-用户名已存在
var UserLogon func(request UserLogonRequest) (code int)

// UserLogin 用户登录暴露方法
// return code 返回状态码：1-成功，2-用户名不存在，3-密码错误
var UserLogin func(request UserLoginRequest) (code int)

// UserQuery 查询用户信息
// return user 返回用户详细信息
var UserQuery func(username string) (user User)

// UserUpdateProfile 更新用户头像
// return code 返回状态码：1-成功，2-失败
var UserUpdateProfile func(request UserUpdateRequest) (code int)

// UserUpdateNick 更新用户昵称
// return code 返回状态码：1-成功，2-失败
var UserUpdateNick func(request UserUpdateRequest) (code int)
