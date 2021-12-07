package vul

import (
	"net/http"
)

func ReflectXss(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	key := url.Get("key")

	// t, _ := template.New("foo").Parse(`<html>Hello, {{.}}!</html>`)
	// t.Execute(w, ("<script>alert('you have been pwned')</script>"))
	// return
	// t.Execute(w, template.html(v))

	if key == "anything" {
		w.Write([]byte("ok,you get"))
		return
	} else {
		w.Write([]byte("<html>oops! not fount: " + key + "</html>"))
		return
	}

}
