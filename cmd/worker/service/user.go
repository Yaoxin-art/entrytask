package service

import (
	"fmt"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// Logon 用户注册
// return username string, err facade.BizError
// username:	注册成功后的用户名（如果传入的用户名带有前或后空格，则会去除空格后再注册并返回）
// err: 		异常，0-成功，1-用户名已存在，2-参数不符合要求
func Logon(request facade.UserLogonRequest) (username string, err facade.BizError) {
	if request.Username == "" || request.Nickname == "" || request.Password == "" {
		// todo: 校验字段长度和字符规范
		return "", facade.BizError(2)
	}
	// insert user
	_, errInsert := insertUser(request)
	if errInsert != nil {
		// failure
		return "", facade.BizError(1)
	}
	go func() {
		user, err := QueryByUsername(request.Username)
		if err == 0 {
			cacheUserInfoIntoRedis(user)
		}
	}()
	return request.Username, facade.BizError(0)
}

// Login 用户登录
// return user facade.User, err facade.BizError
// user: 成功则返回用户信息
// err:  异常，0-成功，1-账号不存在，2-密码错误
func Login(request facade.UserLoginRequest) (user facade.User, err facade.BizError) {
	// TODO
	return facade.User{}, facade.BizError(0)
}

// QueryByUsername 根据username查询用户信息
// return user facade.User, err facade.BizError
// user: 查询的用户信息
// err: 0-成功, 1-用户名不存在
func QueryByUsername(username string) (facade.User, facade.BizError) {
	// query from redis
	cached, errNegligible := queryUserInfoFromRedis(username)
	if errNegligible == nil {
		return *cached, facade.BizError(0)
	}
	// if not exist, query from db, and refresh into redis
	userT, err := queryUserByUsername(username)
	if err != nil {
		logrus.Warnf("query user from db by username:%s failure, msg:%v \n", username, err)
		return facade.User{}, facade.BizError(1)
	}
	user := facade.User{Id: userT.Id, Username: userT.Username, Nickname: userT.Nickname, ProfilePath: userT.ProfilePath}
	// cache into redis
	go cacheUserInfoIntoRedis(user)
	return user, facade.BizError(0)
}

// UpdateUserProfile 更新用户头像
// return user facade.User，err facade.BizError
// user:	成功则返回更新后的用户信息
// err:		异常，0-成功，1-参数异常，2-用户名不存在
func UpdateUserProfile(request facade.UserUpdateRequest) (user facade.User, err facade.BizError) {
	if request.ProfilePath == "" || request.Username == "" {
		return facade.User{}, facade.BizError(1)
	}
	code, errUpdate := updateUserProfile(request)
	if errUpdate != nil || code != 1 {
		return facade.User{}, facade.BizError(2)
	}
	user, errQuery := QueryByUsername(request.Username)
	if errQuery != 0 { // 未找到
		return facade.User{}, facade.BizError(2)
	}
	go cacheUserInfoIntoRedis(user)
	return user, facade.BizError(0)
}

// UpdateUserNick 更新用户昵称
// return user facade.User，err facade.BizError
// user:	成功则返回更新后的用户信息
// err:		异常，0-成功，1-参数异常，2-用户名不存在
func UpdateUserNick(request facade.UserUpdateRequest) (user facade.User, err facade.BizError) {
	if request.Nickname == "" || request.Username == "" {
		return facade.User{}, facade.BizError(1)
	}
	code, errUpdate := updateUserNick(request)
	if errUpdate != nil || code != 1 {
		return facade.User{}, facade.BizError(2)
	}
	user, errQuery := QueryByUsername(request.Username)
	if errQuery != 0 { // 未找到
		return facade.User{}, facade.BizError(2)
	}
	go cacheUserInfoIntoRedis(user)
	return user, facade.BizError(0)
}

// cacheUserInfoIntoRedis 缓存用户信息到redis
// 缓存时间： 30分钟
func cacheUserInfoIntoRedis(user facade.User) {
	// set into redis
	key := userInfoCacheKey(user.Username)
	redisClient.HMSet(key, *userInfo2Map(user))
	redisClient.Expire(key, 30*time.Minute)
	logrus.Debugf("cached user:%v", user)
}

// queryUserInfoFromRedis 根据username从redis中查询用户信息
// return *facade.User, error
// user：获取到的用户信息
// err： 不存在或异常情况，可忽略的异常
func queryUserInfoFromRedis(username string) (user *facade.User, err error) {
	cmd := redisClient.HMGet(userInfoCacheKey(username), userFields...)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("err while hmget for key:%s", username)
	}
	val := cmd.Val()
	if len(val) < 4 || val[0] == nil || val[1] == nil || val[2] == nil || val[3] == nil {
		return nil, fmt.Errorf("not exist")
	}
	id, err := strconv.Atoi(val[0].(string))
	if err != nil {
		return nil, fmt.Errorf("query user from redis with invalid id, %v", err)
	}
	fromRedis := facade.User{
		Id:          int64(id),
		Username:    val[1].(string),
		Nickname:    val[2].(string),
		ProfilePath: val[3].(string),
	}
	return &fromRedis, nil
}

var userFields = []string{"id", "un", "nn", "pp"}

// userInfo2Map 将用户对象转化成map(serialize in redis)
func userInfo2Map(user facade.User) *map[string]interface{} {
	userMap := make(map[string]interface{})
	userMap["id"] = user.Id
	userMap["un"] = user.Username
	userMap["nn"] = user.Nickname
	userMap["pp"] = user.ProfilePath
	return &userMap
}

// userInfoCacheKey user info key in redis
func userInfoCacheKey(username string) string {
	return "et_user::" + username
}
