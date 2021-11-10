package vul

import (
	"encoding/json"
	"fmt"
	"gogoat/utils"
	"net/http"
	"strconv"
)

func isAdmin(r *http.Request) int {
	admin, err := r.Cookie("admin")
	if err != nil {
		return 0
	}
	isadmin, err := strconv.Atoi(admin.Value)
	if err != nil {
		return 0
	}
	if isadmin == 1 {
		return 1
	} else {
		return 0
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {

	var username string
	var usercookie *http.Cookie
	usercookie, err := r.Cookie("username")
	if err != nil {
		// fmt.Println(err.Error())
		username = utils.GetUsername(r)
	} else {
		username = usercookie.Value
	}
	isadmin := isAdmin(r)
	fmt.Println(isadmin)
	if isadmin == 0 {
		user := &utils.User{Username: username}
		ok, _ := utils.X.Get(user)
		if !ok {
			w.Write([]byte("not find user"))
			return
		}
		profile := &userprofile{Username: username}
		profile.Address = user.Address
		profile.Age = user.Age
		profile.Message = user.Message
		v, _ := json.Marshal(profile)
		w.Write([]byte(v))
		return
	}

	if isadmin == 1 {
		users := make([]utils.User, 0)
		profile := make([]userprofile, 0)
		err := utils.X.Find(&users)
		if err != nil {
			w.Write([]byte("system error"))
			return
		}

		err = utils.TranStruct(users, &profile)
		if err != nil {
			w.Write([]byte("system error"))
			fmt.Println(err.Error())
			return
		}
		v, _ := json.Marshal(profile)
		w.Write([]byte(v))

	}

}
