package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/configor"
)

var X *xorm.Engine

type User struct {
	User_id  int64 `xorm:"pk autoincr"`
	Name     string
	Username string
	Balance  float64
	Salary   float64
	Intime   int64
	Password string
}

var Config = struct {
	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     string `default:"3306"`
		Host     string `default:"127.0.0.1"`
	}
}{}

func init() {
	var err error
	X, err = xorm.NewEngine("mysql", GetDBconnStr())
	if err != nil {
		log.Fatal("database connect err:", err)
	}
	if err := X.Sync(new(User)); err != nil {
		log.Fatal("database sync err:", err)
	}
	X.ShowSQL(true)
}

func GetDBconnStr() string {
	configor.Load(&Config, "pkg/vul/config.yaml")
	dbcon := Config.DB.User + ":" + Config.DB.Password + "@tcp(" + Config.DB.Host + ":" + Config.DB.Port + ")/" + Config.DB.Name + "?charset=utf8"
	return dbcon
}

func GetRawDBcon() *sql.DB {
	db, err := sql.Open("mysql", GetDBconnStr())
	if err != nil {
		fmt.Println("raw connect db err")
		panic("raw connect db err")
	}
	return db

}
