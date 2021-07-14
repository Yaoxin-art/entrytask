package service

import "testing"

func TestQueryByPaging(t *testing.T) {
	userList := QueryByPaging(0, 10)
	if len(userList.Users) < 1 {
		t.Error("User list cloud not be empty...")
	}
}
