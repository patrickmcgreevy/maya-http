package libhttp

import (
	"fmt"
	"strings"
)

type Request struct {
	Message string
}

func NewRequest(msg string) Request {
	return Request{
		Message: msg,
	}
}

func (r Request) AsBytes() []byte {
	b := make([]byte, 0)
	b = append(b, []byte(r.Message)...)
	b = append(b, '\n')

	return b
}

func FromBytes(buf []byte) Request {
	return Request{
		Message: string(buf),
	}
}

func (r Request) String() string {
	return string(r.AsBytes())
}

type Response struct {
	ResponseLine responseLine
	Headers      []string // Again will probably be its own type
	Body         string
}

type responseLine struct {
	httpVersion string
	status      StatusCode
}

func newResponseLine(status StatusCode) responseLine {
    return responseLine{
        httpVersion: "HTTP/1.1",
        status: status,
    }
}

func (line responseLine) String() string {
    return fmt.Sprintf("%s %d %s\n", line.httpVersion, int(line.status), line.status.String())
}

func NewResponse(body string) Response {
	return Response{
		ResponseLine: "HTTP/1.1 200 OK",
		Headers:      []string{"Server: Maya/0.1", "Content-type: text/html", "Content-Length: 13"},
		Body:         "Hello, world!",
	}
}

func (r Response) String() string {
	var builder strings.Builder

	builder.WriteString(r.ResponseLine + "\n")
	for _, h := range r.Headers {
		builder.WriteString(h + "\n")
	}
	builder.WriteString("\n" + r.Body)

	return builder.String()
}
