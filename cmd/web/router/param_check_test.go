package router

import (
	"fmt"
	"regexp"
	"testing"
)

func TestCheckStrict(t *testing.T) {
	matchStrict("abcdefghi")
	matchStrict("1bvsda8s9d")
	matchStrict("_88sd sda")
	matchStrict("ssddasdsadsasdsdsadshduassd")
	matchStrict("_88sd s2")
	matchStrict("(_8_8)8s")
	matchStrict("d(_8_8)8s")
	t.Logf("ignore")
}

func matchStrict(username string) {
	if ok, _ := regexp.MatchString(strictStrRegex, username); !ok {
		fmt.Println("match false")
	} else {
		fmt.Println("match true")
	}
}