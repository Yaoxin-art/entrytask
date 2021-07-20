package router

import "fmt"

// 以J
var strictStrRegex = "^[a-zA-Z][\\w]{5,19}$"
var strictStrRegexAllowBlank = "^[a-zA-Z][\\w ]{5,19}(^ )$"

func checkLogonParam(username, nickname, password string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if len(password) < 6 {
		return fmt.Errorf("password too short")
	}
	return nil
}
