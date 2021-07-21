package router

import (
	"fmt"
	"git.garena.com/zhenrong.zeng/entrytask/internal/facade"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"
)

type user struct {
	username string // 从数据库中获取
	password string // 当前所有用户密码都为 "123456"
}

const (
	clientSize = 200
	userSize   = 200
)

var clients []*http.Client
var users []user

func init() {
	initUsers()
	initHttpClients()
}

func initUsers() {
	names := facade.QueryUsernameList(userSize)
	if names == nil {
		fmt.Errorf("init failure for init username list")
		return
	}
	for _, name := range *names {
		users = append(users, user{
			username: name,
			password: "123456",
		})
	}
	logrus.Infof("init user list success, user list size:%d", len(users))
}

func initHttpClients() {
	for i := 0; i < clientSize; i++ {
		client := getClient()
		// 登录
		clientLogin(client, users[i])
		clients = append(clients, client)
	}
}

// clientLogin 登录用户，使client具备登录后的token凭证
func clientLogin(client *http.Client, u user) {

}

func BenchmarkLogin(b *testing.B) {

}

func BenchmarkUpdateNick(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkUpdateProfile(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkInfo(b *testing.B) {

}

const (
	MaxConnsPerHost     int = 0
	MaxIdleConns        int = 0
	MaxIdleConnsPerHost int = 0
)

// getClient init http client
func getClient() *http.Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxConnsPerHost:     MaxConnsPerHost,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		},
		Jar: cookieJar,
	}
	return client
}
