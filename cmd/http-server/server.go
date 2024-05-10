package main

import (
	"fmt"
	"http/pkg/libhttp"
)

type counter int
func (c *counter) ServeHTTP(w libhttp.ResponseWriter, req *libhttp.Request) {
    fmt.Println("in c.ServeHTTP")
    fmt.Println(*c)
    _, err := fmt.Fprint(&w, *c)
    if err != nil {
        fmt.Printf("ServeHTTP error: '%s'\n", err.Error())
        return
    }
    *c++
}

func main() {
    var server libhttp.SimpleServer
    var c counter = 12
    server.Handle("/counter", &c)
    server.ListenAndServe(":8081", nil)
}
