package libhttp

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type Request struct {
	/*
	   The HTTP method. Empty defaults to 'GET'
	*/
	Method string
	/*
	   The URL to access.
	*/
	URL *url.URL
	/*
	   Headers sent by the client
	*/
	Headers headers
	/*
	   The body of the request
	*/
	Body io.ReadCloser
	// ContentLength records the length of the associated content.
	// The value -1 indicates that the length is unknown.
	// Values >= 0 indicate that the given number of bytes may
	// be read from Body.
	ContentLength int64
}

func NewRequest() Request {
	return Request{
		Headers: make(headers),
	}
}

func (r *Request) FromReader(in io.Reader) error {
	buf := bufio.NewReader(in)
	status_line, err := buf.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Println(status_line)
	status_line = status_line[0 : len(status_line)-1]
	parts := strings.Split(status_line, " ")
	if len(parts) != 3 {
		return fmt.Errorf("malformed status line: %s", status_line)
	}
	r.Method = parts[0]
	r.URL, err = url.Parse(parts[1])
	if err != nil {
		return err
	}
	for {
		header, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
		header = strings.TrimSpace(header)
		fmt.Printf("header: '%s'\n", header)
		if header == "" {
			fmt.Println("headers done")
			break
		}
		key, value, found := strings.Cut(header, ":")
		if !found {
			return fmt.Errorf("malformed header: '%s'", header)
		}
		r.Headers[key] = []string{value[1:]}
	}

	return nil
}

func (r Request) AsBytes() []byte {
	return []byte(r.String())
}

func (r Request) String() string {
	var str strings.Builder

	str.WriteString(fmt.Sprintf("%s %s HTTP/1.1\n", r.Method, r.URL.String()))
	for k, v := range r.Headers {
		str.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	return str.String()
}
