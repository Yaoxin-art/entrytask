package router

import (
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
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	profilePath      = "/Users/zhenrong.zeng/Workspaces/Test/golang/EntryTask/htmls"
	profileURIPrefix = "http://127.0.0.1" // config at nginx
)

// ping 服务健康检查
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 1,
		Msg:  "pong",
	})
}

// logon 用户注册
func logon(c *gin.Context) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	// todo: check param
	request := facade.UserLogonRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "bad request param",
			Data: "",
		})
		return
	}
	regName, bizErr := facade.UserLogon(&request)
	if bizErr == 0 {
		// success
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg:  "success",
			Data: regName,
		})
	} else if bizErr == 1 {
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg:  "duplicate username",
			Data: "",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "bad request param",
			Data: "",
		})
	}
}

// login 用户登录
// code: 1-成功，2-密码错误，3-用户名不存在
func login(c *gin.Context) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	request := facade.UserLoginRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "bad request param",
			Data: "",
		})
		return
	}
	if request.Password == "" || request.Username == "" {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "request param not valid",
		})
		return
	}

	user, token, bizErr := facade.UserLogin(&request)
	if bizErr == 0 {
		// success
		fillProfilePrefix(user)                                          // fillProfilePrefix
		c.SetCookie("token", token, 1800, "/", "localhost", false, true) // set cookie
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg:  token,
			Data: user,
		})
	} else if bizErr == 2 {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "failure",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg:  "username not exists",
		})
	}
}

// logout 用户登出
func logout(c *gin.Context) {
	// todo：无要求
	// 删除redis中token记录即可
	c.JSON(http.StatusOK, Response{
		Code: 1,
		Msg:  "not support",
	})
}

// info 已登录用户查询用户信息
// code: 1-成功，0-未登录
func info(c *gin.Context) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		notLogin(c)
		return
	}
	c.Header("Access-Control-Allow-Origin", "localhost")
	user, bizErr := facade.UserQueryByToken(token)
	if bizErr == 0 {
		fillProfilePrefix(user)
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg:  "success",
			Data: user,
		})
		return
	} else {
		// 登录已过期
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Msg:  "login invalid",
		})
		return
	}
}

// findByUsername 根据username查询用户信息
func findByUsername(c *gin.Context) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	username := c.DefaultQuery("username", "")
	if username == "" {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "request param username is empty",
		})
		return
	}
	user, bizErr := facade.UserQuery(username)
	if bizErr == 0 {
		// success
		fillProfilePrefix(user)
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg:  "success",
			Data: user,
		})
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg:  "not exists",
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
			Msg:  "file not valid",
		})
		return "", ""
	}
	filename := header.Filename
	suffix := filename[strings.LastIndex(filename, "."):]

	sum := Md5UploadFile(file)
	pathAppend := "/profile/" + strconv.Itoa(time.Now().Year()) + "_" + sum + suffix
	var path = profilePath + pathAppend
	out, err := os.Create(path)
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			logrus.Errorf("upload file close err:%v", err)
		}
	}(file)

	_, errSeek := file.Seek(0, 0) // 重置文件指针
	if errSeek != nil {
		logrus.Errorf("seek file err:%v", errSeek)
		return "", ""
	}
	_, err = io.Copy(out, file)
	if err != nil {
		if os.IsExist(err) {
			logrus.Errorf("file exist, file name:%s", out.Name())
		} else {
			logrus.Errorf("storage file err:%v", err)
			return "", ""
		}
	}
	logrus.Infof("upload file success, filename:%s, pathAppend:%s, url:%s", filename, pathAppend, profileURIPrefix+pathAppend)
	// storage success
	return pathAppend, profileURIPrefix + pathAppend
}

// profileUpdate 修改用户头像
func profileUpdate(c *gin.Context) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		notLogin(c)
		return
	}
	user, bizErr := facade.UserQueryByToken(token)
	if bizErr != 0 {
		notLogin(c)
		return
	}
	username := user.Username
	profile, url := storageUploadFile(c)
	if profile == "" {
		// file storage failure
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "bad request",
		})
		return
	}
	logrus.Debugf("upload file success, access profile url:%s", url)
	// todo: param check
	request := facade.UserUpdateRequest{
		Username:    username,
		ProfilePath: profile,
	}
	userNew, bizErrUpdate := facade.UserUpdateProfile(&request)
	if bizErrUpdate == 0 {
		// success
		fillProfilePrefix(userNew)
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg:  "success",
			Data: userNew,
		})
		return
	} else if bizErrUpdate == 2 {
		// user not exists
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg:  "not login",
		})
		return
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "bad request",
		})
		return
	}
}

// nickUpdate 修改用户昵称
func nickUpdate(c *gin.Context) {
	start := time.Now().UnixNano()
	defer func() {
		end := time.Now().UnixNano()
		logrus.Infof("login spend time:%d ns", end - start)
	}()
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		notLogin(c)
		return
	}
	user, bizErr := facade.UserQueryByToken(token)
	if bizErr != 0 {
		notLogin(c)
		return
	}
	request := facade.UserUpdateRequest{}
	errReq := c.BindJSON(&request)
	if errReq != nil {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "bad request param",
			Data: "",
		})
		return
	}
	request.Username = user.Username
	// todo: param check
	if request.Nickname == "" {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "bad request param",
		})
		return
	}

	userNew, bizErrUpdate := facade.UserUpdateNick(&request)
	if bizErrUpdate == 0 {
		// success
		fillProfilePrefix(userNew)
		c.JSON(http.StatusOK, Response{
			Code: 1,
			Msg:  "success",
			Data: userNew,
		})
		return
	} else if bizErrUpdate == 2 {
		// user not exists
		c.JSON(http.StatusOK, Response{
			Code: 2,
			Msg:  "not login",
		})
		return
	} else {
		c.JSON(http.StatusOK, Response{
			Code: 3,
			Msg:  "bad request",
		})
		return
	}
}

func fillProfilePrefix(user *facade.User) {
	if user == nil || user.ProfilePath == "" {
		return
	}
	if !strings.HasPrefix(user.ProfilePath, "http") {
		user.ProfilePath = profileURIPrefix + user.ProfilePath
	}
}

func notLogin(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "not login",
	})
}
