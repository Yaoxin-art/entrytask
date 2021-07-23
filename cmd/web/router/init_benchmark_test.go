package router

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	driverName = "mysql"
	dataSource = "root:zeng123456789@tcp(127.0.0.1:3306)/entrytask"
)

func initDB() *sqlx.DB {
	mysql, err := sqlx.Connect(driverName, dataSource)
	if err != nil {
		logrus.Panic(err)
	}
	mysql.SetMaxIdleConns(5)
	mysql.SetMaxOpenConns(10)
	logrus.Info("init mysql client instance success.")
	return mysql
}

var mysqlDB = initDB()

// for test
var sqlSelectUsernameList = ` SELECT username FROM t_user LIMIT ? `

// selectUsernameList 查询出size个用户名
func selectUsernameList(size int) ([]string, error) {
	rows, err := mysqlDB.Queryx(sqlSelectUsernameList, size)
	if err != nil {
		return make([]string, 0, 0), err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)
	var results []string
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			logrus.Errorf("scan username row err:%v", err)
			return results, err
		}
		results = append(results, item)
	}
	return results, nil
}

// http client pool
type httpClient struct {
	client *http.Client
}

type httpClientFactory struct {
	login bool
}


func (f *httpClientFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	return pool.NewPooledObject(
			&httpClient{
				client: initHttpClients(f.login),
			}),
		nil
}

func (f *httpClientFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	// do destroy
	myObj := object.Object.(*httpClient)
	logrus.Debugf("sessoin in poll destroyed, ctx:%v", ctx)
	myObj.client.CloseIdleConnections()
	return nil
}

func (f *httpClientFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	// do validate
	logrus.Debugf("sessoin in pool destroyed, ctx:%v", ctx)
	return true
}

func (f *httpClientFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do activate
	logrus.Debugf("session in pool activate, ctx:%v", ctx)
	return nil
}

func (f *httpClientFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do passivate(put into idle list)
	logrus.Debugf("session in pool passivate, ctx:%v", ctx)
	return nil
}
