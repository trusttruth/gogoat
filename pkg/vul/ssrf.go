package vul

import (
	"fmt"
	"net/http"

	curl "github.com/andelf/go-curl"
)

func Ssrf(w http.ResponseWriter, r *http.Request) {
	var para = r.URL.Query()
	target := para.Get("url")
	easy := curl.EasyInit()
	defer easy.Cleanup()
	easy.Setopt(curl.OPT_URL, target)

	fooTest := func(buf []byte, userdata interface{}) bool {
		println("DEBUG: size=>", len(buf))
		println("DEBUG: content=>", string(buf))
		w.Write((buf))
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

}
