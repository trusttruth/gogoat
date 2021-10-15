package login

import (
	"encoding/json"
	"fmt"
	"gogoat/utils"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func checkPass(user, pass string) bool {
	db := utils.GetRawDBcon()
	sql := "select username from user where username = ? and password = ?"
	passhash := utils.GetSha256(pass)
	re, err := db.Query(sql, user, passhash)
	if err != nil {
		panic("db query err")
	}

	var username string
	for re.Next() {
		err := re.Scan(&username)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		if user == username {
			fmt.Println("The user " + user + " login successed")
			return true
		}
	}
	fmt.Println("The user " + user + " login failed")
	return false
}

func JsonLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/jsonlogin.html")
		t.Execute(w, nil)
	} else {
		body := r.Body
		data := make(map[string]interface{})
		json.NewDecoder(body).Decode(&data)
		// for key, value := range data {
		// 	log.Println("key:", key, " => value :", value)
		// }
		username := data["username"].(string)
		password := data["password"].(string)
		fmt.Println(username, password)
		if username == "" {
			w.Write([]byte("pls login"))
			return
		}
		if password == "" {
			w.Write([]byte("pls login"))
			return
		}
		sessionid := utils.GetRandString(32)

		if checkPass(username, password) {
			sessionCookie := &http.Cookie{
				Name:     "sessionid",
				Value:    sessionid,
				HttpOnly: true,
			}
			http.SetCookie(w, sessionCookie)
			rediscon := utils.GetRedisConn()
			if rediscon == nil {
				fmt.Println("get redis conn err")
			}
			defer rediscon.Close()
			rediscon.Do("set", sessionid, username, "EX", 3600*2)
			// w.Write([]byte("well come " + username))
			http.Redirect(w, r, "/home", 200)
		} else {
			w.Write([]byte("pls login,like http://127.0.0.1/login?username=admin&password=pass"))
		}
	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if username == "" {
			w.Write([]byte("pls login"))
			return
		}
		if password == "" {
			w.Write([]byte("pls login"))
			return
		}
		sessionid := utils.GetRandString(32)

		if checkPass(username, password) {
			sessionCookie := &http.Cookie{
				Name:     "sessionid",
				Value:    sessionid,
				HttpOnly: true,
			}
			http.SetCookie(w, sessionCookie)
			rediscon := utils.GetRedisConn()
			if rediscon == nil {
				fmt.Println("get redis conn err")
			}
			defer rediscon.Close()
			rediscon.Do("set", sessionid, username, "EX", 3600*2)
			// w.Write([]byte("well come " + username))
			http.Redirect(w, r, "/home", http.StatusFound)
		} else {
			w.Write([]byte("pls login,like http://127.0.0.1/login?username=admin&password=pass"))
		}
	}
}

func IsLogin(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, err := r.Cookie("sessionid")
		if err != nil {
			// w.Write([]byte("pls login"))
			http.Redirect(w, r, "/login", http.StatusFound)
			// fmt.Println(err.Error())
		} else {
			if val.Value != "" {
				redisconn := utils.GetRedisConn()
				if redisconn == nil {
					return
				}
				defer redisconn.Close()
				v, err := redisconn.Do("get", val.Value)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(v)
				f(w, r)
			} else {
				// w.Write([]byte("pls login"))
				http.Redirect(w, r, "/login", http.StatusFound)

			}
		}
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/boothome.html")
	t.Execute(w, nil)
}
