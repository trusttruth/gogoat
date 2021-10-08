package vul

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var x *xorm.Engine

type User struct {
	User_id    int64  `xorm:"pk autoincr"` //指定主键并自增
	Name       string `xorm:"unique"`      //唯一的
	Balance    float64
	Time       int64 `xorm:"updated"` //修改后自动更新时间
	Creat_time int64 `xorm:"created"` //创建时间
	//Version    int     `xorm:"version"` //乐观锁
}

func init() {
	var err error
	x, err = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/xorm?charset=utf8")
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	if err := x.Sync(new(User)); err != nil {
		log.Fatal("数据表同步失败:", err)
	}
	x.ShowSQL(true)
}

func Sqlraw(w http.ResponseWriter, r *http.Request) {
	para := r.URL.Query()
	key := para.Get("key")
	// order := para.Get("order")
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	// var sql2 = fmt.Sprintf("select  username  from users where userName= ?")
	// var sql2 = fmt.Sprintf("select  username  from users where userName= ? order by ?")
	// var sql2 = fmt.Sprintf("select  username  from users where userName=\"%s\" order by \"%s\"", key, order)
	var sql2 = fmt.Sprintf("select  username  from users where userName like '%s'  and 1=1 ", strings.TrimSpace(key))
	// var sql2 = fmt.Sprintf("select  username  from users where userName=\"%s\" ", key)
	fmt.Println(sql2)
	// res, err1 := db.Query(sql2, key, order)
	res, err1 := db.Query(sql2)
	fmt.Println(res)
	if err1 != nil {
		fmt.Println(err1.Error())
		w.Write([]byte(err1.Error()))
		return
	}
	s := "you name is : "
	var name string
	for res.Next() {
		err := res.Scan(&name)
		if err != nil {
			// log.Fatal(err)
		}
		fmt.Println(name)
		s += name
	}
	w.Write([]byte(s))

}

func Sqlxorm(w http.ResponseWriter, r *http.Request) {
	para := r.URL.Query()
	key := para.Get("key")
	var user []User
	// err1 := x.Table("user1").Where("name = ?", key).Find(&user)
	err2 := x.Table("user1").Where("name = \"" + key + "\"").Find(&user)
	if err2 != nil {
		fmt.Println(err2)
		w.Write([]byte(err2.Error()))
	}
	s := "you name is : \n"
	for i := 0; i < len(user); i++ {

		fmt.Println(user[i].Name)
		s += user[i].Name
		s += "\n"
	}
	w.Write([]byte(s))

}
