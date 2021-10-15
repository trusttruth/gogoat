package vul

import (
	"fmt"
	"gogoat/utils"
	"html/template"
	"net/http"
)

func ChangePass(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println(utils.GetUsername(r))
		t, _ := template.ParseFiles("views/profile.html")
		t.Execute(w, nil)
	} else {
		fmt.Println(utils.GetUsername(r))
	}
}
