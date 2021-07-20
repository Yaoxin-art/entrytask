package service

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

// TestPrepareUser 仅用于准备用户数据
func TestPrepareUser(t *testing.T) {
	if 1 > 0 { // just for init data, ignore it
		t.Logf("ignore")
		return
	}
	start := time.Now().Second()
	size := 10000000
	for i := 0; i < size; i++ {
		logrus.Infof("i:%d", i)
		TestRegisterUser1(t)
	}
	end := time.Now().Second()
	t.Logf("init %d users, spent %d seconds", size, end-start)
}
