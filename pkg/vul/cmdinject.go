package vul

import (
	"fmt"
	"net/http"
	"os/exec"
)

func Pingcmd(w http.ResponseWriter, r *http.Request) {
	ping := r.URL.Query()
	target := ping.Get("ip")
	testStr := "ping -c 2 " + target
	fmt.Println(testStr)
	testres, _ := exec.Command("sh", "-c", testStr).Output()
	w.Write(testres)

}
