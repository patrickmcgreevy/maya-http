package main

import (
	// "http/pkg/libclient"
	"bufio"
	"fmt"
	"net"
    "http/pkg/libhttp"
	// "os"
)

func SendTraffic() {
	messages := []string{
		"hey",
		"there",
		"sexy",
		"thing!",
	}

    conn, err := net.Dial("tcp", "127.0.0.1:8081")
    if err != nil {
        fmt.Printf("could not open a connection: %s\n", err.Error())
        return
    }
    defer conn.Close()
	for _, m := range messages {
        req := libhttp.NewRequest(m)
		fmt.Fprintf(conn, string(req.AsBytes()))
		message, err := bufio.NewReader(conn).ReadBytes('\n')
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        req = libhttp.FromBytes(message)
		fmt.Printf("Message from server: %+v\n", req)
	}
}

func main() {
	SendTraffic()
}
