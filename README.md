http-auth
=========

Authentication HTTP middleware for Go applications.

## Usage

Here is a ready to use example with [Negroni](https://github.com/codegangsta/negroni):

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/codegangsta/negroni"
  "github.com/K-Phoen/http-negotiate/negotiate"
)

func main() {
  mux := http.NewServeMux()

  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
    fmt.Fprintf(w, "The negotiated format is: " + w.Header().Get("Content-Type"))
  })

  authOptions := &AuthOptions{
    Realm: "Restricted",
    AuthenticationMethod: func(login, password string) bool {
      return login == "test" && password == "tata"
    },
  }

  n := negroni.Classic()
  n.UseHandler(BasicAuth(authOptions))
  n.UseHandler(mux)
  n.Run(":3000")
}
```

Or without [Negroni](https://github.com/codegangsta/negroni) and only for one
handler:

```go
package main

import (
  "github.com/K-Phoen/http-auth/auth"
  "log"
  "net/http"
)

func hello_action(w http.ResponseWriter, req *http.Request) {
  w.Write([]byte("Hello!"))
}

func main() {
  mux := http.NewServeMux()

  basicAuth := auth.BasicAuth(&auth.AuthOptions{
    Realm: "Restricted",
    AuthenticationMethod: func(login, password string) bool {
      return login == "test" && password == "toto"
    },
  })

  mux.Handle("/hello", basicAuth.Wrap(hello_action))

  log.Println("Listening...")
  http.ListenAndServe(":3000", mux)
}
```


## ToDo

  * write tests
  * implement other authentication types

## License

This library is released under the MIT License. See the bundled LICENSE file for
details.
