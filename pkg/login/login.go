package login

import (
	"fmt"
	"net/http"
)

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Checking to see if Authorized header set...")

		if val, ok := r.Header["Authorized"]; ok {
			fmt.Println(val)
			if val[0] == "true" {
				fmt.Println("Header is set! We can serve content!")
				endpoint(w, r)
			}
		} else {
			fmt.Println("Not Authorized!!")
			fmt.Fprintf(w, "Not Authorized!!")
		}
	})
}

func checkPass(user, pass string) bool {
	return true

}

func Login(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	username := url.Get("username")
	password := url.Get("password")
	if ok := checkPass(username, password); ok {

	} else {
		w.Write([]byte("pls login,like http://127.0.0.1/login?username=admin&password=pass"))
	}

}

func IsLogin(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := r.Header["Authorized"]; ok {
			if val[0] == "true" {
				f(w, r)
			}
		} else {
			fmt.Println("pls login")
		}
	})
}
