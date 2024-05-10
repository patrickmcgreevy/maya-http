package main

import (
	"bufio"
	"fmt"
	"http/pkg/libhttp"
	"net"
)

func handleConnection(conn net.Conn) {
    var req libhttp.Request = libhttp.NewRequest()
    defer conn.Close()
	for {
        err := req.FromReader(bufio.NewReader(conn))
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        fmt.Println(req.String())

        resp := libhttp.NewResponse("WHAT IS THE MEANING OF THE WORLD???")
        _, err = conn.Write([]byte(resp.String()))
        if err != nil {
            fmt.Println(err.Error())
            return
        }
	}
}

func RecvTraffic() {
    fmt.Println("Starting server on localhost:8081")
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err.Error())
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("failed to open connection: %s\n", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func main() {
	RecvTraffic()
}
