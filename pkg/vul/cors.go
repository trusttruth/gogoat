package vul

import (
	"fmt"
	"net/http"
)

type userprofile struct {
	username string
	message  string
	age      int
	address  string
}

func Profile(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", h)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "accessToken,appKey,User-Agent,DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	if r.Method == "OPTIONS" {
		fmt.Println("OPTIONS")
		w.WriteHeader(http.StatusOK)
		return
	}
	v := "{\"username\":\"test\"}"
	w.Write([]byte(v))

}
