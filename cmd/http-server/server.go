package main

import (
	"bufio"
	"fmt"
	"http/pkg/libhttp"
	"net"
)

func handleConnection(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			fmt.Println(err.Error())
			return
		}
        req := libhttp.FromBytes(message)
		fmt.Print("Message Received:", req.String())
        // req.Message = strings.ToUpper(req.Message)
		// _, err = conn.Write(req.AsBytes())
        resp := libhttp.NewResponse(req.Message)
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
		fmt.Printf("failed tostart server: %s\n", err.Error())
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
