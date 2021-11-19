package vul

import (
	"encoding/json"
	"fmt"
	"gogoat/utils"
	"html/template"
	"io/ioutil"
	"net/http"
)

func Userprofile(w http.ResponseWriter, r *http.Request) {

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

	// v, _ := json.Marshal(profile)
	t, _ := template.ParseFiles("views/profile.html")
	t.Execute(w, template.HTML(username))
}

func Comment(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/comment.html", "views/front.html")
		t.Execute(w, utils.GetUsername(r))
	} else if r.Method == "POST" {
		// username := utils.GetUsername(r)
		// r.ParseForm()
		// name := r.PostFormValue("name")
		// message := r.PostFormValue("message")
		js, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		mss := make([]utils.Message, 0)
		ms := &utils.Message{}
		fmt.Println(ms)
		err = json.Unmarshal(js, ms)
		if err != nil {
			fmt.Println(err)
			return
		}

		if ms.Name == "" || ms.Message == "" {
			return
		}
		_, err = utils.X.Insert(ms)
		if err != nil {
			fmt.Println(err)
			return
		}

		utils.X.Find(&mss)

		v, err := json.Marshal(mss)
		if err != nil {
			fmt.Println(err)
			return
		}
		// t, _ := template.ParseFiles("views/profile.html")
		// t.Execute(w, v)
		w.Write(v)
		return
	}
}
