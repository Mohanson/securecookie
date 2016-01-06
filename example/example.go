package main

import (
	"fmt"
	"net/http"

	"github.com/mohanson/securecookie"
)

func init() {
	securecookie.Config.SecureKey = "my srcret key"
	securecookie.Config.CacheDays = 10
}

func SayHello(w http.ResponseWriter, req *http.Request) {
	securecookie.SetSecureCookie(w, "u", "mohanson")
	u, err := securecookie.GetSecureCookie(req, "u")
	fmt.Println(u, err)
	if err != nil {
		w.Write([]byte("Hello, visitors"))
	} else {
		w.Write([]byte("Hello, " + u))
	}
}

func main() {
	http.HandleFunc("/hello", SayHello)
	http.ListenAndServe(":8001", nil)
}
