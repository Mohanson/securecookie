package main

import (
	"fmt"
	"net/http"

	"github.com/mohanson/secretcookie"
)

func init() {
	secretcookie.Config.SecretKey = "my srcret key"
	secretcookie.Config.CacheDays = 10
}

func SayHello(w http.ResponseWriter, req *http.Request) {
	secretcookie.SetSecretCookie(w, "u", "jack")
	u, err := secretcookie.GetSecretCookie(req, "u")
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
