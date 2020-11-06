package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jinzhu/gorm"
)

type DBServer struct {
	Host string  //主机
	Port int  //端口
	DBName string  //数据库名
	User string  //登陆用户
	Password string  //密码
}

func (db DBServer) ConnectString() string {
	if db.User == "" {
		return fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable",
			db.Host, db.Port, db.DBName)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.DBName)
}

func (db DBServer) NewGromDB() (DB *gorm.DB ,err error){
	return gorm.Open("postgres",db.ConnectString())
}

func (db DBServer) NewPostgresDB( ) (DB *sqlx.DB, err error) {
	return sqlx.Open("postgres",db.ConnectString())
}
