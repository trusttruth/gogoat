package vul

import (
	"errors"
	"fmt"
	"gogoat/utils"
	"html/template"
	"net/http"
	"strconv"
)

type Money struct {
	Account  string  `json:"account"`
	Amount   float64 `json:"amount"`
	Username string  `json:"username"`
}

func ChangePass(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		username := utils.GetUsername(r)
		t, _ := template.ParseFiles("views/changepassword.html")
		t.Execute(w, username)
	} else {
		username := utils.GetUsername(r)
		r.ParseForm()
		pass01 := r.PostFormValue("newpass01")
		pass02 := r.PostFormValue("newpass02")

		if pass01 != pass02 || pass01 == "" || pass02 == "" {

			t, _ := template.ParseFiles("views/changepassword.html")
			t.Execute(w, username+" you two pass is not same or pass is nil")
		} else {
			pass := utils.GetUserpassFromDb(username)
			if utils.GetSha256(pass01) == pass {
				t, _ := template.ParseFiles("views/changepassword.html")
				t.Execute(w, username+" you pass should be different from now!")
			} else {
				utils.ChangePass(utils.GetSha256(pass01), username)
				w.Write([]byte("password change success"))

			}
		}

	}
}
func transfer(account string, amount float64, r *http.Request) error {
	to_user := &utils.User{}
	to_user.Account = account
	to_amount := utils.GetUserAmount(account)
	if to_amount == -1 {
		return errors.New("not find account " + account)
	}
	to_user.Amount = amount + to_amount

	from_username := utils.GetUsername(r)
	from_id := utils.GetUseridFromUsername(from_username)

	if from_id == -1 {
		return errors.New("not find username " + from_username)
	}
	from_user, err := utils.GetUserinfo(from_id)
	if err != nil {
		return err
	}
	if account == from_user.Account {
		return errors.New("the two account cannot be same")
	}

	if amount > from_user.Amount {
		fmt.Println("balance is not enough")
		return errors.New("balance is not enough")
	}
	from_user.Amount = from_user.Amount - amount

	session := utils.X.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return err
	}
	if _, err := session.ID(utils.GetUseridFromAccount(account)).Update(to_user); err != nil {
		session.Rollback()
		fmt.Println("balance add fail")
		return errors.New("internel error")
	}
	if _, err := session.ID(from_id).Update(from_user); err != nil {
		session.Rollback()
		fmt.Println("balance del fail")
		return errors.New("internel error")
	}

	if err = session.Commit(); err != nil {
		return err
	}
	return nil

}

func Transfer(w http.ResponseWriter, r *http.Request) {
	userm := &Money{}
	user := &utils.User{}
	username := utils.GetUsername(r)
	user.Username = username

	_, err := utils.X.Get(user)
	utils.TranStruct(user, userm)
	if err != nil {
		fmt.Println(err)
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/transfer.html")
		// fmt.Println(userm)
		err := t.Execute(w, *userm)
		if err != nil {
			fmt.Println(err)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		account := r.PostFormValue("account")
		a := r.PostFormValue("amount")

		if account == "" || a == "" {
			return
		} else {
			amount, err := strconv.Atoi(a)
			if err != nil {
				fmt.Println(err)
			}
			if err = transfer(account, float64(amount), r); err != nil {
				fmt.Println(err)
				w.Write([]byte("internel error"))
				return
			}
			t, _ := template.ParseFiles("views/transfer.html")
			// fmt.Println(userm)
			userm.Amount = utils.GetUserAmountFromUsername(username)
			err = t.Execute(w, *userm)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
