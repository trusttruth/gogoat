package login

import (
	"fmt"
	"gogoat/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func isExist(name string) bool {

	db := utils.GetRawDBcon()
	defer db.Close()
	// sql := "select username from user where username = 'yy'"
	re, err := db.Query("select username from user where username = ?", name)
	// re, err := db.Query(sql)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}
	defer re.Close()

	var tmp string
	for re.Next() {
		err := re.Scan(&tmp)
		// fmt.Println(tmp)
		if err != nil {
			return false
		}
		if tmp == name {
			return true
		}
	}
	return false
}
func registerRedriect(w http.ResponseWriter, message interface{}) {
	t, _ := template.ParseFiles("views/register.html")
	t.Execute(w, message)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// t, _ := template.ParseFiles("views/register.html")
		// t.Execute(w, nil)
		registerRedriect(w, nil)
	} else {
		r.ParseForm()
		username := r.PostFormValue("username")
		password1 := r.PostFormValue("password1")
		password2 := r.PostFormValue("password2")
		address := r.PostFormValue("address")
		message := r.PostFormValue("message")
		age := r.PostFormValue("age")

		if age == "" || username == "" || password1 == "" || address == "" || message == "" {
			w.Write([]byte("Please fill in all the information"))
			return
		}

		// if !utils.Validate(username, "username") || !utils.Validate(password1, "password") {
		// 	w.Write([]byte("username or password not match!username need 6 and password need 8-12!"))
		// 	return
		// }

		if password1 != password2 {
			registerRedriect(w, "the two pass not match")
			return
		}
		if isExist(username) {
			registerRedriect(w, username+" is exist")
			return
		}
		hashpass := utils.GetSha256(password1)

		db := utils.GetRawDBcon()
		defer db.Close()

		// sql := "insert into user (username,salary,password) values (?,?,?)"
		// _, err := db.Exec(sql, username, 1000000, hashpass)
		Age, err := strconv.Atoi(age)
		// fmt.Println(Age)
		if err != nil {
			registerRedriect(w, "age is not a number")
			return
		}
		user := new(utils.User)
		user.Address = address
		user.Password = hashpass
		user.Message = message
		user.Username = username
		user.Salary = 10000
		user.Age = int64(Age)

		_, err = utils.X.Insert(user)

		if err != nil {
			fmt.Println(err.Error())
			t, _ := template.ParseFiles("views/register.html")
			t.Execute(w, "register fail")
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
