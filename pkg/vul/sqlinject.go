package vul

import (
	"fmt"
	"gogoat/utils"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func Sqlraw(w http.ResponseWriter, r *http.Request) {
	para := r.URL.Query()
	key := para.Get("key")
	db := utils.GetRawDBcon()
	// var sql2 = "select  name from user where 	Name= ?" //no vulable
	var sql2 = fmt.Sprintf("select  name  from user where Name=\"%s\" ", key) //vulable
	// res, err1 := db.Query(sql2, key)
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
			log.Fatal(err)
		}
		fmt.Println(name)
		s += name
	}
	w.Write([]byte(s))
}

func Sqlxorm(w http.ResponseWriter, r *http.Request) {
	para := r.URL.Query()
	key := para.Get("key")
	// var user []utils.User
	// var users = make([]utils.User, 0)
	var users = new(utils.User)
	// err1 := x.Table("user1").Where("name = ?", key).Find(&user)  novulable!
	// err2 := utils.X.Table("user").Where("name = \"" + key + "\"").Find(&user) //vulable
	_, err2 := utils.X.Table("user").Sum(users, key) //vulable
	// _, err2 := utils.X.Table("user").Sum(users, "User_id`),0)from`user`where`User_id`>'0'#")
	if err2 != nil {
		fmt.Println(err2)
		w.Write([]byte(err2.Error()))
	}
	s := "you name is : \n"
	// for i := 0; i < len(user); i++ {
	// 	fmt.Println(user[i].Name)
	// 	s += user[i].Name
	// 	s += "\n"
	// }
	w.Write([]byte(s))
}
