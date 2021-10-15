package vul

import (
	"net/http"
)

func ReflectXss(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	key := url.Get("key")

	if key == "anything" {
		w.Write([]byte("ok,you get"))
	} else {
		w.Write([]byte("<html>oops! not fount: " + key + "</html>"))

	}

}
