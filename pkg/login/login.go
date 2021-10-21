package login

import (
	"context"
	"encoding/json"
	"fmt"
	"gogoat/utils"
	"html/template"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var ctx = context.Background()

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
		username := data["username"].(string)
		password := data["password"].(string)
		// fmt.Println(username, password)
		if username == "" {
			w.Write([]byte("pls login"))
			return
		}
		if password == "" {
			w.Write([]byte("pls login"))
			return
		}
		if checkPass(username, password) {
			sessionid := utils.GetRandString(32)
			sessionCookie := &http.Cookie{
				Name:     "sessionid",
				Value:    sessionid,
				HttpOnly: true,
			}
			http.SetCookie(w, sessionCookie)
			rediscon := utils.GetRedisConn()
			if rediscon == nil {
				panic("get redis conn err")
			}
			defer rediscon.Close()
			// rediscon.Do("set", sessionid, username, "EX", 3600*2)
			err := rediscon.Set(ctx, sessionid, username, 4*time.Hour).Err()
			if err != nil {
				fmt.Println(err)
			}
			// w.Write([]byte("well come " + username))
			http.Redirect(w, r, "/home", 200)
		} else {
			// w.Write([]byte("pls login,like http://127.0.0.1/login?username=admin&password=pass"))
			// http.Redirect(w, r, "/login", 200)
			t, _ := template.ParseFiles("views/jsonlogin.html")
			t.Execute(w, "pass or username is not correct!")
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
				panic("get redis conn err")
			}
			// defer rediscon.Close()
			// rediscon.Do("set", sessionid, username, "EX", 3600*2)
			err := rediscon.Set(ctx, sessionid, username, 4*time.Hour).Err()
			if err != nil {
				panic(err)
			}
			// w.Write([]byte("well come " + username))
			http.Redirect(w, r, "/home", http.StatusFound)
		} else {
			w.Write([]byte("pls login,like http://127.0.0.1/login?username=admin&password=pass"))
		}
	}
}
func Logout(w http.ResponseWriter, r *http.Request) {
	v, err := r.Cookie("sessionid")
	fmt.Println(v.Value)
	if err != nil {
		fmt.Println("no found session")
		return
	}
	rediscon := utils.GetRedisConn()
	if rediscon == nil {
		panic("get redis conn err")
	}
	defer rediscon.Close()
	err = rediscon.Del(ctx, v.Value).Err()
	if err != nil {
		panic(err)
	}
	// w.Write([]byte("well come " + username))
	http.Redirect(w, r, "/home", http.StatusFound)
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
				fmt.Println(redisconn)
				if redisconn == nil {
					return
				}
				// defer redisconn.Close()
				v, err := redisconn.Get(ctx, val.Value).Result()
				// v, err := redisconn.Do("get", val.Value)
				if err != nil {
					http.Redirect(w, r, "/login", http.StatusFound)
				}
				fmt.Println("v = : ", v)
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
