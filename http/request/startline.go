package request

import (
	"fmt"
	"strings"

	"github.com/inabajunmr/http11server/http"
)

type HTTPMethod int

const (
	GET HTTPMethod = iota
	HEAD
	POST
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
)

type StartLine struct {
	Method        HTTPMethod
	RequestTarget string
	Version       http.HTTPVersion
}

func (s StartLine) ToString() string {
	return fmt.Sprintf("%v %v %v", s.Method.ToString(), s.RequestTarget, s.Version.ToString()) // TODO

}

func ParseStartLine(line string) (*StartLine, error) {
	s := strings.Split(line, " ")
	if len(s) != 3 {
		// not HTTP/1.1 request
		return nil, &http.HTTPError{Msg: "this request is not for HTTP/1.1", Status: 400}
	}

	method, err := getHTTPMethod(s[0])
	if err != nil {
		return nil, &http.HTTPError{Msg: fmt.Sprintf("HTTP method %v is not implemented", s[0]), Status: 400}
	}
	requestTarget := s[1] // TODO parse
	httpVersion := s[2]
	if httpVersion != "HTTP/1.1" {
		return nil, &http.HTTPError{Msg: fmt.Sprintf("%v is not supported HTTP version", s[2]), Status: 400}
	}

	return &StartLine{method, requestTarget, http.HTTP11}, nil

}

func (m HTTPMethod) ToString() string {
	switch m {
	case GET:
		return "GET"
	case HEAD:
		return "HEAD"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case CONNECT:
		return "CONNECT"
	case OPTIONS:
		return "OPTIONS"
	case TRACE:
		return "TRACE"
	default:
		return ""
	}
}

func getHTTPMethod(method string) (HTTPMethod, error) {
	switch method {
	case "GET":
		return GET, nil
	case "HEAD":
		return HEAD, nil
	case "POST":
		return POST, nil
	case "PUT":
		return PUT, nil
	case "DELETE":
		return DELETE, nil
	case "CONNECT":
		return CONNECT, nil
	case "OPTIONS":
		return OPTIONS, nil
	case "TRACE":
		return TRACE, nil
	default:
		return 0, fmt.Errorf("unsupported method %v", method)
	}
}
