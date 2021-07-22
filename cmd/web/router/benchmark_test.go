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
	"testing"
	"time"
)

type user struct {
	Username string `json:"username"` // 从数据库中获取
	Password string `json:"password"` // 当前所有用户密码都为 "123456"
}

type nicked struct {
	Username string `json:"username"` // 从数据库中获取
	Password string `json:"password"` // 当前所有用户密码都为 "123456"
	Nickname string `json:"nickname"` // 昵称
}

const httpServerAddr = "http://127.0.0.1:7777"

const (
	clientSize = 200
	userSize   = 20000
)

var clients []*http.Client
var users []user

func initForBenchmark(login bool) {
	clients = make([]*http.Client, 0)
	users = make([]user, 0)

	initUsers()
	initHttpClients(login)
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
	logrus.Debugf("init user list success, user list size:%d, users:%v", len(users), users)
}

func initHttpClients(login bool) {
	for i := 0; i < clientSize; i++ {
		client := getClient()
		// 登录
		if login {
			clientLogin(client, users[i])
		}
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
	reqUrl := httpServerAddr + "/user/login"
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(data))
	//res, err := client.Post(reqUrl, "application/json", bytes.NewBuffer(data))
	res, err := client.Do(req)
	if err != nil {
		logrus.Errorf("post login err:%v", err)
	}
	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		logrus.Errorf("get body err:%v", err)
	}
	err = res.Body.Close()
	if err != nil {
		logrus.Errorf("close body err:%v", err)
	}
	logrus.Debugf("post login response:%v", string(body[:]))
}

func BenchmarkLogin(b *testing.B) {
	initForBenchmark(false)
	defer destroyHttpClients()

	b.ResetTimer()
	parallelism := clientSize
	b.SetParallelism(parallelism)

	defer fmt.Printf("parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			id := rand.Intn(clientSize)
			client := clients[id]
			iu := rand.Intn(userSize)
			u := users[iu]
			data, errData := json.Marshal(u)
			if errData != nil {
				logrus.Panicf("json error:%v", u)
				return
			}
			reqUrl := httpServerAddr + "/user/login"
			req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(data))
			if err != nil {
				b.Error(err)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
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
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logrus.Errorf("login get body err:%v", err)
			}
			err = res.Body.Close()
			if err != nil {
				logrus.Errorf("login close body err:%v", err)
			}
			logrus.Debugf("login response:%v", string(body[:]))
		}
	})
}

func BenchmarkUpdateNick(b *testing.B) {
	initForBenchmark(true)
	defer destroyHttpClients()

	b.ResetTimer()
	parallelism := clientSize / 20
	b.SetParallelism(parallelism)
	fmt.Printf("parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			id := rand.Intn(clientSize)
			client := clients[id]
			iu := rand.Intn(userSize)
			u := users[iu]
			un := nicked{
				Username: u.Username,
				Nickname: "New Nick",
			}
			data, errData := json.Marshal(un)
			if errData != nil {
				logrus.Panicf("json error:%v", un)
				return
			}
			reqUrl := httpServerAddr + "/user/login"
			req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(data))
			if err != nil {
				b.Error(err)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
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
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logrus.Errorf("login get body err:%v", err)
			}
			err = res.Body.Close()
			if err != nil {
				logrus.Errorf("login close body err:%v", err)
			}
			logrus.Debugf("login response:%v", string(body[:]))
		}
	})
}

func BenchmarkInfoFix(b *testing.B) {
	initForBenchmark(false)
	defer destroyHttpClients()

	b.ResetTimer()
	parallelism := clientSize
	b.SetParallelism(parallelism)
	fmt.Printf("parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			id := rand.Intn(clientSize)
			client := clients[id]
			uid := rand.Intn(id)
			user := users[uid]
			requestUrl := httpServerAddr + "/user/find?username=" + user.Username
			req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
			if err != nil {
				b.Error(err)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
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
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logrus.Errorf("get body err:%v", err)
			}
			err = res.Body.Close()
			if err != nil {
				logrus.Errorf("close body err:%v", err)
			}
			logrus.Debugf("find response:%v", string(body[:]))

		}
	})
}


func BenchmarkInfoRandom(b *testing.B) {
	initForBenchmark(false)
	defer destroyHttpClients()

	b.ResetTimer()
	parallelism := clientSize
	b.SetParallelism(parallelism)
	fmt.Printf("parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			id := rand.Intn(clientSize)
			client := clients[id]
			uid := rand.Intn(userSize)
			user := users[uid]
			requestUrl := httpServerAddr + "/user/find?username=" + user.Username
			req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
			if err != nil {
				b.Error(err)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
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
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logrus.Errorf("get body err:%v", err)
			}
			err = res.Body.Close()
			if err != nil {
				logrus.Errorf("close body err:%v", err)
			}
			logrus.Debugf("info response:%v", string(body[:]))

		}
	})
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
