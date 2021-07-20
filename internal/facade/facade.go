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

// User.String
func (u User) String() string {
	return fmt.Sprintf("Id:%d, Username:%s, Nickname:%s, Profile:%s", u.Id, u.Username, u.Nickname, u.ProfilePath)
}

// UserLogonRequest 用户注册时的请求对象
type UserLogonRequest struct {
	Username string
	Nickname string
	Password string
}

// UserLogonRequest.String 去除密码打印
func (ulr *UserLogonRequest) String() string {
	return fmt.Sprintf("UserLogonRequest{Username:%s, Nickname:%s}", ulr.Username, ulr.Nickname)
}

// UserLoginRequest 用户登录时的请求对象
type UserLoginRequest struct {
	Username string
	Password string
}

// UserLoginRequest.String 去除密码打印
func (ulr *UserLoginRequest) String() string {
	return fmt.Sprintf("UserLoginRequest{Username:%s}", ulr.Username)
}

// UserUpdateRequest 用户更新昵称和头像时的请求对象
// 更新昵称时，Username和Nickname不可为空
// 更新头像时，Username和ProfilePath不可为空
type UserUpdateRequest struct {
	Username    string
	Nickname    string
	ProfilePath string
}

// BizError 业务异常类型
type BizError int

// BizError
// 0-success, other then error
func (be *BizError) Error() string {
	return fmt.Sprintf("bisness error code:%d", be)
}

// UserLogon 用户注册暴露方法: "Logon"
// return username string, err BizError
// username:	注册成功后的用户名（如果传入的用户名带有前或后空格，则会去除空格后再注册并返回）
// err: 		异常，0-成功，1-用户名已存在，2-参数不符合要求
var UserLogon func(request UserLogonRequest) (username string, err BizError)

// UserLogin 用户登录暴露方法: "Login"
// return user User, err BizError
// user:	成功则返回更新后的用户信息
// err:		异常，0-成功，1-账号不存在，2-密码错误
var UserLogin func(request UserLoginRequest) (user User, err BizError)

// UserQuery 查询用户信息暴露方法: "QueryUser"
// return user User, err BizError
// user:	返回用户详细信息
// err:		异常，0-成功，1-用户名不存在
var UserQuery func(username string) (user User, err BizError)

// UserUpdateProfile 更新用户头像暴露方法: "UpdateNick"
// return user User, err BizError
// user:	成功则返回更新后的用户信息
// err:		异常，0-成功，1-参数异常，2-用户名不存在
var UserUpdateProfile func(request UserUpdateRequest) (user User, err BizError)

// UserUpdateNick 更新用户昵称暴露方法: "UpdateProfile"
// return user User, err BizError
// user:	成功则返回更新后的用户信息
// err:		异常，0-成功，1-参数异常，2-用户名不存在
var UserUpdateNick func(request UserUpdateRequest) (user User, err BizError)

const (
	Logon         = "Logon"
	Login         = "Login"
	Query         = "QueryUser"
	UpdateNick    = "UpdateNick"
	UpdateProfile = "UpdateProfile"
)
