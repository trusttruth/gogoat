package vul

import (
	"encoding/json"
	"fmt"
	"gogoat/utils"
	"net/http"
)

type userprofile struct {
	Username string
	Message  string
	Age      int64
	Address  string
}

func Profile(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", h)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json")
	// if r.Method == "OPTIONS" {
	// 	// fmt.Println("OPTIONS")
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }

	username := utils.GetUsername(r)
	user := &utils.User{Username: username}
	ok, _ := utils.X.Get(user)
	if ok {
		fmt.Println(user)
	}
	profile := &userprofile{Username: username}
	profile.Address = user.Address
	profile.Age = user.Age
	profile.Message = user.Message
	v, _ := json.Marshal(profile)

	w.Write([]byte(v))
}
