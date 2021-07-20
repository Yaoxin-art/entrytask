package router

import (
	"crypto/md5"
	"fmt"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Response struct {
	Code 	int			`json:"code"`
	Msg 	string		`json:"msg"`
	Data 	interface{}	`json:"data"`
}

var (
	profilePath		= "/Users/zhenrong.zeng/Workspaces/Data/entrytask"
	profileURIPrefix= "https://127.0.0.1/entrytask/static"	// config at nginx
)

// ping 服务健康检查
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 1,
		Msg: "pong",
	})
}

// logon 用户注册
func logon(c *gin.Context) {
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	nickname := c.DefaultPostForm("nickname", "")
	// todo: check param
	request := facade.UserLogonRequest{
		Username: username,
		Nickname: nickname,
		Password: password,
	}
	regName, bizErr := facade.UserLogon(request)
	if bizErr == 0 {
		// success
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg: "success",
			Data: regName,
		})
	} else if bizErr == 1 {
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg: "duplicate username",
			Data: "",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg: "bad request param",
			Data: "",
		})
	}
}

// login 用户登录
func login(c *gin.Context) {
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	// todo: param check
	request := facade.UserLoginRequest{
		Username: username,
		Password: password,
	}
	user, token, bizErr := facade.UserLogin(request)
	if bizErr == 0 {
		// success
		c.Header("token", token)	// todo: generate token

		fillProfilePrefix(user)	// fillProfilePrefix
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg: token,
			Data: user,
		})
	} else if bizErr == 2 {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg: "failure",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg: "username not exists",
		})
	}
}

// logout 用户登出
func logout(c *gin.Context) {
	// todo：无要求
	// 删除redis中token记录即可
	c.JSON(http.StatusOK, Response{
		Code: 1,
		Msg: "not support",
	})
}

// info 已登录用户查询用户信息
func info(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		notLogin(c)
		return
	}
	user, bizErr := facade.UserQueryByToken(token)
	if bizErr == 0 {
		fillProfilePrefix(user)
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg: "success",
			Data: user,
		})
		return
	} else {
		// 登录已过期
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Msg: "login invalid",
		})
		return
	}
}

// findByUsername 根据username查询用户信息
func findByUsername(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	user, bizErr := facade.UserQuery(username)
	if bizErr == 0 {
		// success
		fillProfilePrefix(user)
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg: "success",
			Data: user,
		})
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg: "not exists",
		})
	}
}

// storageUploadFile 保存上传的文件，并返回数据库中的存储路径和访问路径
// return profile, url string
// profile:	数据库中的保存路径
// url:		前端访问路径
func storageUploadFile(c *gin.Context) (profile, url string) {
	file, header, err := c.Request.FormFile("profile")
	if err != nil {
		logrus.Errorf("upload file, get form file err:%v", err)
		c.JSON(http.StatusBadRequest, Response{
			Code: 0,
			Msg: "file not valid",
		})
		return "", ""
	}
	filename := header.Filename
	logrus.Infof("get file:%s", filename)
	buf := make([]byte, 0, header.Size)
	size, err := file.Read(buf)
	fmt.Printf("file header size:%d, buf size:%d, buf len:%d", header.Size, size, len(buf))
	sumByte := md5.Sum(buf)
	sum := *(*string)(unsafe.Pointer(&sumByte))
	pathAppend := "/profile/" + strconv.Itoa(time.Now().Year()) + sum
	var path = profilePath + pathAppend
	out, err := os.Create(path)
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			logrus.Errorf("upload file close err:%v", err)
		}
	}(file)
	_, err = io.Copy(out, file)
	if err != nil {
		logrus.Errorf("storage file err:%v", err)
		return "", ""
	}
	logrus.Infof("upload file success, filename:%s, pathAppend:%s, url:%s", filename, pathAppend, profileURIPrefix + pathAppend)
	// storage success
	return pathAppend, profileURIPrefix+pathAppend
}

// profileUpdate 修改用户头像
func profileUpdate(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {

	}
	username := c.DefaultPostForm("username", "")
	profile, url := storageUploadFile(c)
	if profile == "" {
		// file storage failure
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg: "bad request",
		})
		return
	}
	logrus.Debugf("upload file success, access profile url:%s", url)
	// todo: param check
	request := facade.UserUpdateRequest{
		Username: username,
		ProfilePath: profile,
	}
	user, bizErr := facade.UserUpdateProfile(request)
	if bizErr == 0 {
		// success
		fillProfilePrefix(user)
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg: "success",
			Data: user,
		})
		return
	} else if bizErr == 2 {
		// user not exists
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg: "not login",
		})
		return
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg: "bad request",
		})
		return
	}
}

// nickUpdate 修改用户昵称
func nickUpdate(c *gin.Context) {
	username := c.DefaultPostForm("username", "")
	nickname := c.DefaultPostForm("nickname", "")
	// todo: param check
	request := facade.UserUpdateRequest{
		Username: username,
		Nickname: nickname,
	}
	user, bizErr := facade.UserUpdateNick(request)
	if bizErr == 0 {
		// success
		fillProfilePrefix(user)
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg: "success",
			Data: user,
		})
		return
	} else if bizErr == 2 {
		// user not exists
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg: "not login",
		})
		return
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg: "bad request",
		})
		return
	}
}

func fillProfilePrefix(user facade.User) {
	if !strings.HasPrefix(user.ProfilePath, "http") {
		user.ProfilePath = profileURIPrefix + user.ProfilePath
	}
}


func notLogin(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg: "not login",
	})
}