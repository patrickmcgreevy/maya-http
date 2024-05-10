package libhttp

import (
	"bufio"
	"fmt"
	"net"
)

type Server interface {
	ListenAndServe(address string, handler *Handler)
	HandleFunc(path string, handler *Handler)
}

type Handler func(ResponseWriter, *Request)

type ResponseWriter struct {
}

type SimpleServer struct {
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (w ResponseWriter) Close() {
}

func (s *SimpleServer) ListenAndServe(address string, handler *Handler) {
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

		go s.handleConnection(conn)
	}
}

func (s *SimpleServer) HandleFunc(path string, handler *Handler) {
}

func (s *SimpleServer) handleConnection(conn net.Conn) {
	var req Request = NewRequest()
	defer conn.Close()
	for {
		err := req.FromReader(bufio.NewReader(conn))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(req.String())

		resp := NewResponse("WHAT IS THE MEANING OF THE WORLD???")
		_, err = conn.Write([]byte(resp.String()))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

}
