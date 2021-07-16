package zerorpc

import (
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"testing"
	"time"
)

type User struct {
	Name string
	Nick string
	Age  int
}

var userRecord = make(map[string]User)

func init() {
	userRecord["zero"] = User{Name: "zero", Nick: "zero-zero", Age: 27}
	userRecord["one"] = User{Name: "one", Nick: "one-one", Age: 22}
	userRecord["two"] = User{Name: "two", Nick: "two-two", Age: 31}
}

func queryUser(name string) (User, error) {
	if u, ok := userRecord[name]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("user:%s not exists", name)
}

func TestRpcInvoke(t *testing.T) {
	gob.Register(User{})
	addr := "127.0.0.1:6666"
	srv := NewServer(addr)

	srv.Register("queryUser", queryUser)

	go srv.Run()

	time.Sleep(2 * time.Second)
	log.Info("continue test")

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}

	cli := NewClient(conn)

	log.Info("server connected...")

	var query func(string) (User, error)

	cli.callRpc("queryUser", &query)

	u, errQuery := query("zero")
	if errQuery != nil {
		t.Error(errQuery)
	}
	t.Logf("query success, u:%v \n", u)

	//_, errYes := query("not_exist")
	//if errYes != nil {
	//	t.Logf("query success, user should not exist.")
	//} else {
	//	t.Error("query show not exist.")
	//}
}
