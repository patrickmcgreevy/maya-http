package libhttp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type Server interface {
	Handle(string, Handler)
	HandleFunc(path string, handler func(ResponseWriter, *Request))
	ListenAndServe(address string, handler Handler)
}

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type ResponseWriter interface {
	Close()
    Header() headers // TODO: exertnal interface should be map[str]str??
	Write(b []byte) (int, error)
    WriteHeader(statusCode int)
}

type SimpleResponseWriter struct {
    *Response
	buf    *[]byte
	writer io.WriteCloser
}

type SimpleServer struct {
	handlers map[string]Handler
}

type FileHandler struct {
	Path string
}

func NewResponseWriter(w io.WriteCloser) SimpleResponseWriter {
	buf := make([]byte, 0, 100)

	return SimpleResponseWriter{
        Response: NewResponse(""),
		buf:    &buf,
		writer: w,
	}
}

func (rw *SimpleResponseWriter) Write(b []byte) (int, error) {
	fmt.Printf("in ResponseWriter.Write. b='%s'\n", string(b))
    rw.setBody(string(b))
    fmt.Println("body:", rw.body)

	return len(b), nil
}

func (rw *SimpleResponseWriter) Close() {
	defer rw.writer.Close()
	fmt.Printf("in ResponseWriter.Close. body='%s'\n", rw.body)
    // rw.Response.body = string(*rw.buf)
	// resp := NewResponse(string(*rw.buf))
	fmt.Println("in ResponseWriter.Close.", "resp=", rw.Response)
	_, err := rw.writer.Write([]byte(rw.Response.String()))
	if err != nil {
		fmt.Printf("error: '%s'", err.Error())
	}
}

func (rw *SimpleResponseWriter) Header() headers {
	return nil
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
	var respWriter SimpleResponseWriter = NewResponseWriter(conn)
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
        error404(&respWriter, &req)
		fmt.Println(err.Error())
		return
	}
	handler.ServeHTTP(&respWriter, &req)
}

func (s *SimpleServer) getHandler(path string) (Handler, error) {
	fmt.Printf("gethandler: '%s'\n", path)
	if h, ok := s.handlers[path]; ok {
		return h, nil
	}

	return nil, fmt.Errorf("no handler registered for path '%s'", path)
}

func (h FileHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	dat, err := os.ReadFile(h.Path)
	if err != nil {
		// TODO 500 error handling etc
		return
	}
	rw.Write(dat)
}

func error404(rw ResponseWriter, req *Request) {
    fmt.Println("in error404")
    rw.WriteHeader(404)
    fmt.Fprintf(rw, "Error 404: %s not found", req.URL.Path)
}
