package main

import (
	"flag"
	"fmt"
	"gogoat/pkg/login"
	"gogoat/pkg/vul"
	"net/http"
)

func main() {

	port := flag.String("port", "8888", "listen port")
	flag.Parse()

	http.HandleFunc("/login", login.JsonLogin)
	http.HandleFunc("/logout", login.Logout)
	http.HandleFunc("/", login.JsonLogin)
	http.HandleFunc("/register", login.Register)
	http.Handle("/ssrf", login.IsLogin(vul.Ssrf))
	http.Handle("/home", login.IsLogin(login.Home))
	http.Handle("/ping", login.IsLogin(vul.Pingcmd))
	http.Handle("/sqlraw", login.IsLogin(vul.Sqlraw))
	http.Handle("/sqlxorm", login.IsLogin(vul.Sqlxorm))
	http.Handle("/upload", login.IsLogin(vul.Upload))
	http.Handle("/reflectXss", login.IsLogin(vul.ReflectXss))
	http.Handle("/csrf", login.IsLogin(vul.ChangePass))

	addr := "127.0.0.1:" + *port
	fmt.Println("App where be runing on " + addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		fmt.Println(err)
	}

}
