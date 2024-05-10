package main

import (
	"http/pkg/libhttp"
)

func main() {
    var server libhttp.SimpleServer
    server.ListenAndServe(":8081", nil)
}
