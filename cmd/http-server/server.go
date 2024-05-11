package main

import (
	"fmt"
	"http/pkg/libhttp"
	"io"
)

type fileReader io.ReaderFrom

type Stringer interface {
	String() string
}

type counter int

func (c *counter) ServeHTTP(w libhttp.ResponseWriter, req *libhttp.Request) {
	fmt.Println("in c.ServeHTTP")
	fmt.Println(*c)
	_, err := fmt.Fprint(w, *c)
	if err != nil {
		fmt.Printf("ServeHTTP error: '%s'\n", err.Error())
		return
	}
	*c++
}

type raw_string string

func (s raw_string) ServeHTTP(w libhttp.ResponseWriter, req *libhttp.Request) {
	fmt.Fprint(w, s)
}

func main() {
	var server *libhttp.SimpleServer = &libhttp.SimpleServer{}
	var c counter = 12
    var maya  = libhttp.FileHandler{
        Path: "index.html",
    }
	server.Handle("/counter", &c)
	server.Handle("/maya", maya)
	server.Handle("/", maya)
	server.ListenAndServe(":8081", nil)
}
