package vul

import (
	"gogoat/utils"
	"html/template"
	"net/http"
)

func ChangePass(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		username := utils.GetUsername(r)
		t, _ := template.ParseFiles("views/profile.html")
		t.Execute(w, username)
	} else {
		username := utils.GetUsername(r)
		r.ParseForm()
		pass01 := r.PostFormValue("newpass01")
		pass02 := r.PostFormValue("newpass02")

		if pass01 != pass02 || pass01 == "" || pass02 == "" {

			t, _ := template.ParseFiles("views/profile.html")
			t.Execute(w, username+" you two pass is not same or pass is nil")
		} else {
			pass := utils.GetUserpassFromDb(username)
			if utils.GetSha256(pass01) == pass {
				t, _ := template.ParseFiles("views/profile.html")
				t.Execute(w, username+" you pass should be different from now!")
			} else {
				utils.ChangePass(utils.GetSha256(pass01), username)
				w.Write([]byte("password change success"))

			}
		}

	}
}
