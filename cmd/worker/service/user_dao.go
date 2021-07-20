package service

import (
	"fmt"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
	"github.com/sirupsen/logrus"
)

type TUser struct {
	Id          int64  `db:"id"`
	Username    string `db:"username""`
	Nickname    string `db:"nickname"`
	Password    string `db:"passwd"`
	ProfilePath string `db:"profile_path"`
}

// insert t_user
var sqlInsertUser = `INSERT INTO t_user(username, nickname, passwd, created_timestamp, modified_timestamp) 
	VALUES(?, ?, CONCAT('*', UPPER(SHA1(UNHEX(SHA1(?))))), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

var sqlUpdateProfile = ` UPDATE t_user SET profile_path=?, modified_timestamp=CURRENT_TIMESTAMP WHERE username=? `

var sqlUpdateNickname = ` UPDATE t_user SET nickname=?, modified_timestamp=CURRENT_TIMESTAMP WHERE username=? `

var sqlSelectByUsername = ` SELECT id,username,nickname,passwd,profile_path FROM t_user WHERE username=? LIMIT 1 `

// insertUser 插入用户记录
// return 返回状态码，1-成功，0-失败
func insertUser(request facade.UserLogonRequest) (int, error) {
	result, err := mysqlDB.Exec(sqlInsertUser, request.Username, request.Nickname, request.Password)
	if err != nil {
		// eg: Duplicate entry 'xxx' for key 't_user.unique_idx_username'
		logrus.Warnf("insert user err:%v, request:%s", err, request)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("insert user err: %v, sql:%s", err, sqlInsertUser)
		return 0, err
	}
	lastId, errId := result.LastInsertId()
	if errId != nil {
		logrus.Errorf("insert user err: %v, sql:%s", err, sqlInsertUser)
		return 0, errId
	}
	logrus.Infof("insert result, rows:%v, last id:%v \n", rows, lastId)
	return 1, nil
}

// queryUserByUsername 根据username查询用户信息
func queryUserByUsername(username string) (*TUser, error) {
	rows, err := mysqlDB.Queryx(sqlSelectByUsername, username)
	if err != nil {
		logrus.Errorf("query user by username:%s err:%v", username, err)
		return nil, fmt.Errorf("username:%s not exists", username)
	}
	// SQL限定只返回一行
	var user TUser
	if rows.Next() {
		errMapper := rows.StructScan(&user)
		if errMapper != nil {
			return nil, fmt.Errorf("query user:%s result mapper err:%v", username, errMapper)
		}
		logrus.Infof("query user:%s by username, get:%v ", username, user)
		return &user, nil
	} else {
		return nil, fmt.Errorf("query user:%s not exists", username)
	}
}

// updateUserProfile 更新用户头像
func updateUserProfile(request facade.UserUpdateRequest) (int, error) {
	res, err := mysqlDB.Exec(sqlUpdateProfile, request.ProfilePath, request.Username)
	if err != nil {
		logrus.Errorf("update user:%s profile failure:%v", request.Username, err)
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		logrus.Errorf("update user:%s profile failure, rows:%v, err:%v", request.Username, rows, err)
		return 0, err
	}
	logrus.Infof("update user:%s profile success", request.Username)
	return 1, nil
}

// updateUserNick 更新用户昵称
func updateUserNick(request facade.UserUpdateRequest) (int, error) {
	res, err := mysqlDB.Exec(sqlUpdateNickname, request.Nickname, request.Username)
	if err != nil {
		logrus.Errorf("update user:%s nickname failure:%v", request.Username, err)
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		logrus.Errorf("update user:%s nickname failure, rows:%v, err:%v", request.Username, rows, err)
		return 0, err
	}
	logrus.Infof("update user:%s nickname success", request.Username)
	return 1, nil
}
