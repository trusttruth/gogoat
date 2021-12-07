package vul

import (
	"net/http"
	"os/exec"
)

func Pingcmd(w http.ResponseWriter, r *http.Request) {
	ping := r.URL.Query()
	target := ping.Get("ip")
	testStr := "ping -c 2 " + target
	// fmt.Println(testStr)
	res := "the command is: " + testStr

	testres, _ := exec.Command("sh", "-c", testStr).Output()
	w.Write([]byte(res + "\n" + string(testres)))

}
