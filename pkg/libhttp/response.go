package libhttp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)


type Response struct {
	ResponseLine  responseLine
	Headers       headers
	Body          string
	ContentLength int
}

type responseLine struct {
	httpVersion string
	status      StatusCode
}

type headerList []string
type headers map[string]headerList

func (h *headers) String() string {
	var str strings.Builder

	for k, v := range *h {
		str.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	return str.String()
}

func (hl headerList) String() string {
    return strings.Join(hl, ", ")
}

func newResponseLine(status StatusCode) responseLine {
	return responseLine{
		httpVersion: "HTTP/1.1",
		status:      status,
	}
}

func (line responseLine) String() string {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("%s %d %s\n", line.httpVersion, int(line.status), line.status.String()))
	return str.String()
}

func NewResponse(body string) Response {
    fmt.Printf("in NewResponse. body='%s'\n", body)
	return Response{
		ResponseLine: newResponseLine(OK),
		Headers: map[string]headerList{
			"Server":       {"Maya/0.1"},
			"Content-type": {"text/html"},
			// "Content-Length": {fmt.Sprint(len(body))},
		},
		Body:          body,
		ContentLength: len(body),
	}
}

func (r Response) String() string {
	var builder strings.Builder

	builder.WriteString(r.ResponseLine.String())
	for k, v := range r.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
    builder.WriteString(fmt.Sprintf("Content-Length: %d\n", r.ContentLength))
	builder.WriteString("\n" + r.Body)

	return builder.String()
}

func ReadHttpResponse(reader io.Reader) (Response, error) {
	var resp Response

	buf_reader := bufio.NewReader(reader)
	resp.Headers = make(map[string]headerList)

	response_line, err := buf_reader.ReadString('\n')
	if err != nil {
		return Response{}, err
	}
	response_line = response_line[0 : len(response_line)-1]
	fmt.Printf("response_line: '%s'\n", response_line)
	fields := strings.Fields(response_line)
	if len(fields) != 3 {
		return Response{}, fmt.Errorf("malformed response line: %s", response_line)
	}
	ver, code := fields[0], fields[1]
	status_code, err := strconv.Atoi(code)
	if err != nil {
		return Response{}, fmt.Errorf("malformed status code: %s", code)
	}
	resp.ResponseLine = responseLine{httpVersion: ver, status: StatusCode(status_code)}

	for {
		header, err := buf_reader.ReadString('\n')
		if err != nil {
			return Response{}, err
		}
		header = header[0 : len(header)-1]
		if header == "" {
			break
		}
		fmt.Printf("header: '%s'\n", header)
		parts := strings.Split(header, ":")
		if len(parts) != 2 {
			return Response{}, fmt.Errorf("malformed header: '%s'", header)
		}
		resp.Headers[parts[0]] = []string{parts[1][1:len(parts[1])]}
	}

	length, ok := resp.Headers["Content-Length"]
	if !ok {
		return Response{}, fmt.Errorf("missing 'Content-Length' header")
	}
	l, err := strconv.Atoi(length[0])
	if err != nil {
		return Response{}, fmt.Errorf("malformed 'Content-Length' header: %s", length)
	}

	buf := make([]byte, l)
	_, err = io.ReadFull(buf_reader, buf)
	if err != nil {
		return Response{}, err
	}

	resp.Body = string(buf)

	return resp, nil
}
