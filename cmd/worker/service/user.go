package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// QueryUsernameList 查询size个用户名
func QueryUsernameList(size int) *[]string {
	users, err := selectUsernameList(size)
	if err != nil {
		return nil
	}
	return &users
}

// QueryByToken 通过token查询已登录的用户
// return user facade.User, err facade.BizErr
// user:	用户信息
// err:		异常，0-存在且成功，1-失败或未找到（无效token或登录已过期）
func QueryByToken(token string) (user facade.User, err int) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	username := getTokenUsername(token)
	if username == "" {
		return facade.User{}, 1
	}
	return QueryByUsername(username)
}

// Logon 用户注册
// return username string, err facade.BizError
// username:	注册成功后的用户名（如果传入的用户名带有前或后空格，则会去除空格后再注册并返回）
// err: 		异常，0-成功，1-用户名已存在，2-参数不符合要求
func Logon(request *facade.UserLogonRequest) (username string, err int) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	if request.Username == "" || request.Nickname == "" || request.Password == "" {
		// todo: 校验字段长度和字符规范
		return "", 2
	}
	// insert user
	_, errInsert := insertUser(*request)
	if errInsert != nil {
		// failure
		return "", 1
	}
	go func() {
		user, err := QueryByUsername(request.Username)
		if err == 0 {
			cacheUserInfoIntoRedis(user)
		}
	}()
	return request.Username, 0
}

// Login 用户登录
// return user facade.User, err int
// user: 	成功则返回用户信息
// token:	登录成功之后返回token
// err:  	异常，0-成功，1-账号不存在，2-密码错误
func Login(request *facade.UserLoginRequest) (user facade.User, token string, err int) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	// 校验账号密码是否正确
	// if not exist, query from db, and refresh into redis
	userT, daoErr := queryUserByUsername(request.Username)
	if daoErr != nil {
		logrus.Warnf("query user from db by username:%s failure, msg:%v \n", request.Username, err)
		return facade.User{}, "", 1
	}
	encoded := selectPassword(request.Password)
	if strings.Compare(encoded, userT.Password) != 0 { // 密码不正确
		logrus.Warnf("login failure, username:%s", request.Username)
		return facade.User{}, "", 2
	}
	token = generateToken(request.Username)
	go cacheTokenForUser(token, userT.Username)
	go cacheUserInfoIntoRedis(*userConvert(*userT))
	return *userConvert(*userT), token, 0
}

// QueryByUsername 根据username查询用户信息
// return user facade.User, err facade.BizError
// user: 查询的用户信息
// err: 0-成功, 1-用户名不存在
func QueryByUsername(username string) (facade.User, int) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	// query from redis
	cached, errNegligible := queryUserInfoFromRedis(username)
	if errNegligible == nil {
		return *cached, 0
	}
	// if not exist, query from db, and refresh into redis
	userT, err := queryUserByUsername(username)
	if err != nil {
		logrus.Warnf("query user from db by username:%s failure, msg:%v \n", username, err)
		return facade.User{}, 1
	}
	user := *userConvert(*userT)
	// cache into redis
	go cacheUserInfoIntoRedis(user)
	return user, 0
}

// UpdateUserProfile 更新用户头像
// return user facade.User，err facade.BizError
// user:	成功则返回更新后的用户信息
// err:		异常，0-成功，1-参数异常，2-用户名不存在
func UpdateUserProfile(request *facade.UserUpdateRequest) (user facade.User, err int) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	if request.ProfilePath == "" || request.Username == "" {
		return facade.User{}, 1
	}
	code, errUpdate := updateUserProfile(*request)
	if errUpdate != nil || code != 1 {
		return facade.User{}, 2
	}
	user, errQuery := QueryByUsername(request.Username)
	if errQuery != 0 { // 未找到
		return facade.User{}, 2
	}
	user.ProfilePath = request.ProfilePath
	go cacheUserInfoIntoRedis(user)
	return user, 0
}

// UpdateUserNick 更新用户昵称
// return user facade.User，err facade.BizError
// user:	成功则返回更新后的用户信息
// err:		异常，0-成功，1-参数异常，2-用户名不存在
func UpdateUserNick(request *facade.UserUpdateRequest) (user facade.User, err int) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	if request.Nickname == "" || request.Username == "" {
		return facade.User{}, 1
	}
	code, errUpdate := updateUserNick(*request)
	if errUpdate != nil || code != 1 {
		return facade.User{}, 2
	}
	user, errQuery := QueryByUsername(request.Username)
	if errQuery != 0 { // 未找到
		return facade.User{}, 2
	}
	user.Nickname = request.Nickname
	go cacheUserInfoIntoRedis(user)
	return user, 0
}

// getTokenUsername 根据token获取登录的用户名
func getTokenUsername(token string) string {
	key := keyUserToken(token)
	cmd := redisClient.Get(key)
	if cmd.Err() != nil {
		return ""
	}
	return cmd.Val()
}

// cacheTokenForUser 保存用户登录token与用户名的关系
func cacheTokenForUser(token, username string) {
	key := keyUserToken(token)
	redisClient.Set(key, username, 30*time.Minute)
	logrus.Debugf("cache user login token:%s -> %s", username, key)
}

func keyUserToken(token string) string {
	return "et_token::" + token
}

// cacheUserInfoIntoRedis 缓存用户信息到redis
// 缓存时间： 30分钟
func cacheUserInfoIntoRedis(user facade.User) {
	// set into redis
	key := userInfoCacheKey(user.Username)
	redisClient.HMSet(key, userInfo2Map(user))
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
func userInfo2Map(user facade.User) map[string]interface{} {
	userMap := make(map[string]interface{})
	userMap["id"] = user.Id
	userMap["un"] = user.Username
	userMap["nn"] = user.Nickname
	userMap["pp"] = user.ProfilePath
	return userMap
}

// userInfoCacheKey user info key in redis
func userInfoCacheKey(username string) string {
	return "et_user::" + username
}

//
func userConvert(userT TUser) *facade.User {
	return &facade.User{Id: userT.Id, Username: userT.Username, Nickname: userT.Nickname, ProfilePath: userT.ProfilePath}
}

func generateToken(username string) string {
	timestamp := time.Now().Format(time.RFC3339)
	mux := timestamp + strconv.Itoa(rand.Int()) + username
	h := md5.New()
	h.Write([]byte(mux))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
