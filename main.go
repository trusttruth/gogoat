package main

import (
	"fmt"
	"gogoat/pkg/login"
	"gogoat/pkg/vul"
	"net/http"
)

func main() {

	stat := http.StatusText(200)
	fmt.Println(stat) //状态码200对应的状态OK

	stringtype := http.DetectContentType([]byte("test"))
	fmt.Println(stringtype) //text/plain; charset=utf-8

	// http.HandleFunc("/test", Test)
	http.HandleFunc("/login", login.Login)
	http.Handle("/ssrf", login.IsLogin(vul.Ssrf))
	http.HandleFunc("/ping", vul.Pingcmd)
	http.HandleFunc("/sql1", vul.Sqlxorm)
	http.HandleFunc("/sql2", vul.Sqlraw)
	// http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println(err)
	}

}
