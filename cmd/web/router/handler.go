package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ping 服务健康检查
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "pong",
	})
}

//
//func query(c *gin.Context) {
//	offsetS := c.DefaultQuery("offset", "0")
//	limitS := c.DefaultQuery("limit", "10")
//	offset, offsetErr := strconv.Atoi(offsetS)
//	limit, limitErr := strconv.Atoi(limitS)
//	if offsetErr != nil || limitErr != nil {
//		logrus.Warnf("Invalid query param, offset:%s, limit:%s", offsetS, limitS)
//	}
//	res := nil
//	c.JSON(200, res)
//}

// logon 用户注册
func logon(c *gin.Context) {

}

// login 用户登录
func login(c *gin.Context) {

}

// logout 用户登出
func logout(c *gin.Context) {

}

// info 已登录用户查询用户信息
func info(c *gin.Context) {

}

// findByUsername 根据username查询用户信息
func findByUsername(c *gin.Context) {

}

// profileUpdate 修改用户头像
func profileUpdate(c *gin.Context) {

}

// nickUpdate 修改用户昵称
func nickUpdate(c *gin.Context) {

}
