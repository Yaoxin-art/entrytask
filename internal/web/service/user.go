package service

import (
	"fmt"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
	"github.com/sirupsen/logrus"
)

func QueryUserByUsername(username string) (*facade.User, error) {
	// query from redis
	cached, errNegligible := queryUserInfoFromRedis(username)
	if errNegligible == nil {
		return cached, nil
	}
	// if not exist, query from db, and refresh into redis
	userT, err := queryUserByUsername(username)
	if err != nil {
		logrus.Warnf("query user from db by username:%s failure, msg:%v \n", username, err)
		return nil, err
	}
	user := facade.User{Id: userT.Id, Username: userT.Username, Nickname: userT.Nickname, ProfilePath: userT.ProfilePath}
	// cache into redis
	go cacheUserInfoIntoRedis(user)
	return &user, nil
}

func RegisterUser(request facade.UserLogonRequest) (int, error) {
	if request.Username == "" || request.Nickname == "" || request.Password == "" {
		// todo: 校验字段长度和字符规范
		return 0, fmt.Errorf("register with invalid param:%v", request)
	}
	// insert user
	code, errInsert := insertUser(request)
	if errInsert != nil {
		// failure
		return 2, errInsert
	}
	go func() {
		user, err := QueryUserByUsername(request.Username)
		if err == nil {
			cacheUserInfoIntoRedis(*user)
		}
	}()
	return code, nil
}

// UpdateUserProfile 更新用户头像
// return *facade.User，error
// user：返回更新后的用户信息
// error：更新失败的情况下返回失败异常
func UpdateUserProfile(request facade.UserUpdateRequest) (*facade.User, error) {
	if request.ProfilePath == "" || request.Username == "" {
		return nil, fmt.Errorf("request param not valid, request:%v", request)
	}
	code, errUpdate := updateUserProfile(request)
	if errUpdate != nil {
		return nil, errUpdate
	}
	if code != 1 {
		return nil, fmt.Errorf("update profile, user:%s not exists", request.Username)
	}
	user, errQuery := QueryUserByUsername(request.Username)
	if errQuery != nil {
		return nil, errQuery
	}
	go cacheUserInfoIntoRedis(*user)
	return user, nil
}

// UpdateUserNick 更新用户昵称
// return *facade.User，error
// user：返回更新后的用户信息
// error：更新失败的情况下返回失败异常
func UpdateUserNick(request facade.UserUpdateRequest) (user *facade.User, err error) {
	if request.Nickname == "" || request.Username == "" {
		return nil, fmt.Errorf("request param not valid, request:%v", request)
	}
	code, errUpdate := updateUserNick(request)
	if errUpdate != nil {
		return nil, errUpdate
	}
	if code != 1 {
		return nil, fmt.Errorf("update nickname, user:%s not exists", request.Username)
	}
	user, errQuery := QueryUserByUsername(request.Username)
	if errQuery != nil {
		return nil, errQuery
	}
	go cacheUserInfoIntoRedis(*user)
	return user, nil
}

// cacheUserInfoIntoRedis 缓存用户信息到redis
func cacheUserInfoIntoRedis(user facade.User) error {
	// todo: redis op

	logrus.Infof("cached user:%v", user)
	return nil
}

// queryUserInfoFromRedis 根据username从redis中查询用户信息
// return *facade.User, error
// user：获取到的用户信息
// err： 不存在或异常情况，可忽略的异常
func queryUserInfoFromRedis(username string) (user *facade.User, err error) {
	return nil, fmt.Errorf("not exist")
}
