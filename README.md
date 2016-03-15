## Secure Cookie for "net/http"(and other web framework)

Inspired by Python tornado web framework

## Short Example
```golang
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

```

Open browse, and visit http://localhost:8001/hello, you will see "hello, visitors", if you visit again, you will see "hello, jack".

## Cookie
When you write
```
securecookie.SetSecureCookie(w, "u", "jack")
```
the truth cookie persist is
```
u=2|1:0|10:1451987075|1:u|8:amFjaw==|d0095993d18e39e5e59149ba4cfa066cb7f03a5e32e35149d838213c12de8170
```
When you write
```
securecookie.GetSecureCookie(req, "u")
```
you will get "jack"

## Thanks
tornado (python web framework)
