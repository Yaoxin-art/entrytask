package router

import "fmt"

func checkLoginParam(username, password string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if len(password) < 6 {
		return fmt.Errorf("password too short")
	}
	return nil
}
