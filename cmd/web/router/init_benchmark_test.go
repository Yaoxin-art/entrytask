package router

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
	defer rows.Close()
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
