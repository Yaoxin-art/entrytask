package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/sirupsen/logrus"
	"io"
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


var users []user
var clients *HttpClientPool

type HttpClientPool struct {
	clientPool *pool.ObjectPool
}
func NewHttpClientPool(login bool, size int) *HttpClientPool {
	ctx := context.Background()
	config := pool.ObjectPoolConfig{
		MaxTotal:           size,
		MaxIdle:            size,
		BlockWhenExhausted: true,
	}
	return &HttpClientPool{
		clientPool: pool.NewObjectPool(ctx, &httpClientFactory{login: login}, &config),
	}
}

func initForBenchmark(login bool) {
	clients = NewHttpClientPool(login, clientSize)
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

func initHttpClients(login bool) *http.Client {
	i := rand.Intn(userSize)
	client := getClient()
	// 登录
	if login {
		clientLogin(client, users[i])
	}
	return client
}
func destroyHttpClients() {
	ctx := context.Background()
	clients.clientPool.Close(ctx)
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
		logrus.Errorf("login request close err:%v", err)
	}
	err = req.Body.Close()
	if err != nil {
		logrus.Errorf("post login err:%v", err)
		return
	}
	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		logrus.Errorf("get body err:%v", err)
		return
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

	ctx := context.Background()
	b.ResetTimer()
	parallelism := clientSize
	b.SetParallelism(parallelism)

	defer fmt.Printf("benchmark login parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			client := obj.(*httpClient).client
			iu := rand.Intn(userSize)
			u := users[iu]
			data, errData := json.Marshal(u)
			if errData != nil {
				logrus.Panicf("json error:%v", u)
				b.Skipped()
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}
			reqUrl := httpServerAddr + "/user/login"
			req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(data))
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
			}
			err = req.Body.Close()
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}

			if res.StatusCode != http.StatusOK {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					_ = clients.clientPool.ReturnObject(ctx, obj)
					continue
				}
			} else {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					logrus.Errorf("close body err:%v", err)
				}
			}
			_ = clients.clientPool.ReturnObject(ctx, obj)
		}
	})
}

func BenchmarkUpdateNick(b *testing.B) {
	initForBenchmark(true)
	defer destroyHttpClients()

	ctx := context.Background()
	b.ResetTimer()
	parallelism := clientSize / 20
	b.SetParallelism(parallelism)
	fmt.Printf("benchmark update nickname parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			client := obj.(*httpClient).client
			iu := rand.Intn(userSize)
			u := users[iu]
			un := nicked{
				Username: u.Username,
				Nickname: "New Nick",
			}
			data, errData := json.Marshal(un)
			if errData != nil {
				logrus.Panicf("json error:%v", un)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}
			reqUrl := httpServerAddr + "/user/login"
			req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(data))
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
			}
			err = req.Body.Close()
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}

			if res.StatusCode != http.StatusOK {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					_ = clients.clientPool.ReturnObject(ctx, obj)
					continue
				}
			} else {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					logrus.Errorf("close body err:%v", err)
				}
			}
			_ = clients.clientPool.ReturnObject(ctx, obj)
		}
	})
}

func BenchmarkInfoFix(b *testing.B) {
	initForBenchmark(false)
	//defer destroyHttpClients()

	ctx := context.Background()
	b.ResetTimer()
	parallelism := clientSize
	b.SetParallelism(parallelism)
	fmt.Printf("benchmark info fix parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			client := obj.(*httpClient).client
			id := rand.Intn(clientSize)
			u := users[id]
			requestUrl := httpServerAddr + "/user/find?username=" + u.Username
			req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}

			if res.StatusCode != http.StatusOK {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					_ = clients.clientPool.ReturnObject(ctx, obj)
					continue
				}
			} else {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					logrus.Errorf("close body err:%v", err)
				}
			}
			_ = clients.clientPool.ReturnObject(ctx, obj)
		}
	})
}

func BenchmarkInfoRandom(b *testing.B) {
	initForBenchmark(false)
	defer destroyHttpClients()

	ctx := context.Background()
	b.ResetTimer()
	parallelism := clientSize
	b.SetParallelism(parallelism)
	fmt.Printf("benchmark info random parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			client := obj.(*httpClient).client
			uid := rand.Intn(userSize)
			u := users[uid]
			requestUrl := httpServerAddr + "/user/find?username=" + u.Username
			req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}

			if res.StatusCode != http.StatusOK {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					_ = clients.clientPool.ReturnObject(ctx, obj)
					continue
				}
			} else {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					logrus.Errorf("close body err:%v", err)
				}
			}
			_ = clients.clientPool.ReturnObject(ctx, obj)
		}
	})
}

const (
	MaxConnsPerHost     int = 1
	MaxIdleConns        int = 0
	MaxIdleConnsPerHost int = 0
)

// getClient init http client
func getClient() *http.Client {
	cookieJar, err := cookiejar.New(&cookiejar.Options{
	})
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   1 * time.Second,
				KeepAlive: 90 * time.Second,
			}).DialContext,
			MaxConnsPerHost:     MaxConnsPerHost,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		},
		Jar: cookieJar,
	}
	return client
}
