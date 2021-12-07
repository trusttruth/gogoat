package vul

import (
	"encoding/json"
	"fmt"
	"gogoat/utils"
	"net/http"
)

type userprofile struct {
	Username string  `json:"username"`
	Message  string  `json:"message"`
	Age      int64   `json:"age"`
	Address  string  `json:"address"`
	Account  string  `json:"account"`
	Amount   float64 `json:"amount"`
}

type Usert struct {
	Username string
	Address  string
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("Origin")

	w.Header().Set("Access-Control-Allow-Origin", h)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json")
	username := utils.GetUsername(r)
	user := &utils.User{Username: username}
	ok, err := utils.X.Get(user)
	if !ok {
		w.Write([]byte("not find user"))
		return
	}
	if err != nil {
		w.Write([]byte("not find user"))
		return
	}

	profile := &userprofile{}

	err = utils.TranStruct(user, profile)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("system err,pls contact admin !"))
		return
	}

	v, _ := json.Marshal(profile)
	w.Write([]byte(v))
}
