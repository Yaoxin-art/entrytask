package service

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	redisNetwork = "tcp"
	redisAddr    = "localhost:6379"
)

func initRedis() *redis.Client {
	options := redis.Options{
		Network:            redisNetwork,
		Addr:               redisAddr,
		Dialer:             nil,
		OnConnect:          nil,
		Password:           "",
		DB:                 0,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolSize:           100,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
	}
	logrus.Info("init redis client instance success.")
	// 新建一个client
	return redis.NewClient(&options)
	// close
	// defer ClientRedis.Close()
}

var redisClient = initRedis()
