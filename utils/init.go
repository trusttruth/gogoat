package utils

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	// "github.com/go-xorm/xorm"
)

var X *xorm.Engine

type User struct {
	Id       int64  `xorm:"pk autoincr"`
	Username string `xorm:"unique" json:"username"`
	Salary   float64
	Password string
	Message  string  `json:"message"`
	Address  string  `json:"address"`
	Age      int64   `json:"age"`
	IsAdmin  bool    `xorm:"default false"`
	Amount   float64 `json:"amount"`
	Account  string  `json:"account"`
}

type Message struct {
	Username string
	Name     string
	Message  string
}

//login info,if login success ,  the field of success will be 1 else will be 0
type LoginStatus struct {
	Success bool
	Err     string
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
	if err := X.Sync(new(Message)); err != nil {
		log.Fatal("database sync err:", err)
	}
	X.ShowSQL(true)
}
