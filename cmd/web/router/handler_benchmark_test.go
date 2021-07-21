package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"runtime"
	"testing"
	"time"
)

type user struct {
	Username string `json:"username"` // 从数据库中获取
	Password string `json:"password"`// 当前所有用户密码都为 "123456"
}

const httpServerAddr = "http://localhost:7777"

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
	names, err := selectUsernameList(userSize)
	if err != nil || names == nil {
		logrus.Errorf("init failure for init username list")
		return
	}
	for _, name := range names {
		users = append(users, user{
			Username: name,
			Password: "123456",
		})
	}
	logrus.Infof("init user list success, user list size:%d, users:%v", len(users), users)
}

func initHttpClients() {
	for i := 0; i < clientSize; i++ {
		client := getClient()
		// 登录
		clientLogin(client, users[i])
		clients = append(clients, client)
	}
}
func destroyHttpClients() {
	for _, client := range clients {
		client.CloseIdleConnections()
	}
}

// clientLogin 登录用户，使client具备登录后的token凭证
func clientLogin(client *http.Client, u user) {
	data, errData := json.Marshal(u)
	if errData != nil {
		logrus.Panicf("json error:%v", u)
		return
	}
	var err error
	logrus.Infof("login user:%v", u)
	logrus.Infof("login user:%v", string(data))
	reqUrl := httpServerAddr + "/user/login"
	req, err := http.NewRequest(http.MethodPost, reqUrl,  bytes.NewBuffer(data))
	//res, err := client.Post(reqUrl, "application/json", bytes.NewBuffer(data))
	res, err := client.Do(req)
	if err != nil {
		logrus.Errorf("post login err:%v", err)
	}
	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		logrus.Errorf("get body err:%v", err)
	}
	logrus.Infof("post login response:%v", string(body[:]))
}

func BenchmarkLogin(b *testing.B) {
	defer destroyHttpClients()
	b.ResetTimer()
	b.SetParallelism(clientSize / 15)
	fmt.Printf("parallelism:%d", clientSize / runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		defer func() {
			logrus.Infof("test ...")
		}()
		for pb.Next() {
			id := rand.Intn(clientSize)
			client := clients[id]
			requestUrl := httpServerAddr + "/user/info"
			req, errReq:= http.NewRequest(http.MethodGet, requestUrl, nil)
			if errReq != nil {
				b.Error(errReq)
				continue
			}
			res, errRes := client.Do(req)
			if errRes != nil {
				b.Error(errRes)
				continue
			}
			if res.StatusCode != http.StatusOK {
				resBody, _ := ioutil.ReadAll(res.Body)
				b.Error(string(resBody))
				err := res.Body.Close()
				if err != nil {
					b.Error(err)
				}
				continue
			}
			_, err := ioutil.ReadAll(res.Body)
			if err != nil {
				b.Error(err)
				continue
			}
			errRes = res.Body.Close()
			if errRes != nil {
				b.Error(errRes)
			}

		}
	})
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
