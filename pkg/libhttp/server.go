package libhttp

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Server interface {
	ListenAndServe(address string, handler Handler)
	Handle(string, Handler)
	HandleFunc(path string, handler func(ResponseWriter, *Request))
}

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type ResponseWriter struct {
	buf    *[]byte
	writer io.WriteCloser
}

type SimpleServer struct {
    handlers map[string]Handler
}

func NewResponseWriter(w io.WriteCloser) ResponseWriter {
    buf := make([]byte, 0, 100)
    
    return ResponseWriter{
        buf: &buf,
        writer: w,
    }
}
func (rw *ResponseWriter) Write(b []byte) (int, error) {
    fmt.Printf("in ResponseWriter.Write. b='%s'\n", string(b))
	if rw.buf == nil {
        return 0, fmt.Errorf("cannot use a zero-value ResponseWriter")
	}
	*rw.buf = append(*rw.buf, b...)
    fmt.Println("buf:", string(*rw.buf))

	return len(b), nil
}

func (rw *ResponseWriter) Close() {
	defer rw.writer.Close()
    fmt.Printf("in ResponseWriter.Close. buf='%s'\n", string(*rw.buf))
    resp := NewResponse(string(*rw.buf))
    fmt.Println("in ResponseWriter.Close.", "resp=", resp)
    _, err := rw.writer.Write([]byte(resp.String()))
    if err != nil {
        fmt.Printf("error: '%s'", err.Error())
    }
}

var _ Server = &SimpleServer{}

func (s *SimpleServer) ListenAndServe(address string, handler Handler) {
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

func (s *SimpleServer) Handle(path string, h Handler) {
    if s.handlers == nil {
        s.handlers = make(map[string]Handler)
    }
    s.handlers[path] = h
}

func (s *SimpleServer) HandleFunc(path string, handler func(ResponseWriter, *Request)) {
}

func (s *SimpleServer) handleConnection(conn net.Conn) {
	var req Request = NewRequest()
    var handler Handler
	var respWriter ResponseWriter = NewResponseWriter(conn)
	defer respWriter.Close()
    err := req.FromReader(bufio.NewReader(conn))
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(req.String())

    // call the correct handler and pass in respWriter and req
    handler, err = s.getHandler(req.URL.Path)
    if err != nil {
        // TODO: ResponseWriter needs to have methods for setting status codes and headers
        fmt.Println(err.Error())
        return
    }
    handler.ServeHTTP(respWriter, &req)
}

func (s *SimpleServer) getHandler(path string) (Handler, error) {
    fmt.Printf("gethandler: '%s'\n", path)
    if h, ok := s.handlers[path]; ok {
        return h, nil
    }

    return nil, fmt.Errorf("no handler registered for path '%s'", path)
}
