package router

import (
	"git.garena.com/zhenrong.zeng/entrytask/internal/web/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func ping(c *gin.Context) {
	path := c.Request.RequestURI
	log.Infof("Request path:%s", path)
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func query(c *gin.Context) {
	offsetS := c.DefaultQuery("offset", "0")
	limitS := c.DefaultQuery("limit", "10")
	offset, offsetErr := strconv.Atoi(offsetS)
	limit, limitErr := strconv.Atoi(limitS)
	if offsetErr != nil || limitErr != nil {
		log.Warnf("Invalid query param, offset:%s, limit:%s", offsetS, limitS)
	}
	res := service.QueryByPaging(offset, limit)
	c.JSON(200, res)
}
