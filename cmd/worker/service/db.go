package service

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
	mysql.SetMaxIdleConns(100)
	mysql.SetMaxOpenConns(200)
	logrus.Info("init mysql client instance success.")
	return mysql
}

var mysqlDB = initDB()
