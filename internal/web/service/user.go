package service

import (
	"fmt"
)

type UserList struct {
	Total       int    // total size
	PageIdx     int    // current page number
	CurrentSize int    // current page element size
	Users       []User // user element list
}

func (ul UserList) String() string {
	return fmt.Sprintf("Total:%d, PageIdx:%d, Current Page Size:%d, User List:%v", ul.Total, ul.PageIdx, ul.CurrentSize, ul.Users)
}

type User struct {
	Id          int
	Username    string
	Nickname    string
	ProfilePath string
}

func (u User) String() string {
	return fmt.Sprintf("Id:%d, Username:%s, Nickname:%s, Profile:%s", u.Id, u.Username, u.Nickname, u.ProfilePath)
}

func QueryByPaging(offset, limit int) UserList {
	userList := make([]User, 2)
	userList[0] = User{1, "zero", "nickname1", "https://avatars.githubusercontent.com/u/8939585?s=120&v=4"}
	userList[1] = User{2, "one", "nickname2", "https://avatars.githubusercontent.com/u/8939585?s=120&v=4"}
	return UserList{2, 1, 1, userList}
}
