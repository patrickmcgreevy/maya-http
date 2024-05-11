package libhttp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Response struct {
	responseLine  responseLine
	headers       headers
	body          string
	contentLength int
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

func (h *headers) Add(key, value string) {
	if v, ok := (*h)[key]; ok {
		(*h)[key] = append(v, value)
		return
	}
	(*h)[key] = []string{value}
}

func (h *headers) Del(key string) {
	delete(*h, key)
}

func (h *headers) Get(key string) string {
	// return (*h)[key][0]
	if v, ok := (*h)[key]; ok {
		return v[0]
	}

	return ""
}

func (h *headers) Set(key, value string) {
	(*h)[key] = []string{value}
}

func (h *headers) Values(key string) headerList {
	if v, ok := (*h)[key]; ok {
		return v
	}

	return []string{}
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
		responseLine: newResponseLine(OK),
		headers: map[string]headerList{
			"Server":       {"Maya/0.1"},
			"Content-type": {"text/html"},
		},
		body:          body,
		contentLength: len(body),
	}
}

func (r Response) Header() headers {
    return r.headers
}

func (r Response) WriteHeader(statusCode int) {
    r.responseLine = newResponseLine(StatusCode(statusCode))
}

func (r Response) String() string {
	var builder strings.Builder

	builder.WriteString(r.responseLine.String())
	for k, v := range r.headers {
		builder.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	builder.WriteString(fmt.Sprintf("Content-Length: %d\n", r.contentLength))
	builder.WriteString("\n" + r.body)

	return builder.String()
}

func ReadHttpResponse(reader io.Reader) (Response, error) {
	var resp Response

	buf_reader := bufio.NewReader(reader)
	resp.headers = make(map[string]headerList)

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
	resp.responseLine = responseLine{httpVersion: ver, status: StatusCode(status_code)}

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
		resp.headers[parts[0]] = []string{parts[1][1:len(parts[1])]}
	}

	length, ok := resp.headers["Content-Length"]
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

	resp.body = string(buf)

	return resp, nil
}
