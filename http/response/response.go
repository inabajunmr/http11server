package response

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/inabajunmr/http11server/http"
	"github.com/inabajunmr/http11server/http/header"
	"github.com/inabajunmr/http11server/http/request"
)

type Response interface {
	Response(conn net.Conn)
}

type EchoResponse struct {
	Version      http.HTTPVersion
	StatusCode   int
	ReasonPhrase string
	Request      request.Request
}

type HeadResponse struct {
	Version      http.HTTPVersion
	StatusCode   int
	ReasonPhrase string
	Request      request.Request
}

type OptionsResponse struct {
	Version http.HTTPVersion
	Request request.Request
}

func GetResponse(req request.Request) Response {
	switch req.StartLine.Method {
	case request.HEAD:
		return HeadResponse{Version: http.HTTP11, StatusCode: 200, ReasonPhrase: "OK", Request: req}
	case request.OPTIONS:
		return OptionsResponse{Version: http.HTTP11, Request: req}
	default:
		return EchoResponse{Version: http.HTTP11, StatusCode: 200, ReasonPhrase: "OK", Request: req}

	}
}

func (r EchoResponse) StatusLine() string {
	return fmt.Sprintf("%v %v %v\n", r.Version.ToString(), r.StatusCode, r.ReasonPhrase)
}

func (r EchoResponse) Headers() header.Headers {
	headers := header.Headers{}
	headers = append(headers, &header.Header{FieldName: "Content-Length", FieldValue: strconv.Itoa(len(r.Body()))})
	return headers
}

func (r EchoResponse) Body() []byte {
	return echoBody(r.Request)
}

func (r EchoResponse) Response(conn net.Conn) {
	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
	conn.Write([]byte(r.Body()))
}

func (r HeadResponse) StatusLine() string {
	return fmt.Sprintf("%v %v %v\n", r.Version.ToString(), r.StatusCode, r.ReasonPhrase)
}

func (r HeadResponse) Headers() header.Headers {
	headers := header.Headers{}
	headers = append(headers, &header.Header{FieldName: "Content-Length", FieldValue: strconv.Itoa(len(r.Body()))})
	return headers
}

func (r HeadResponse) Response(conn net.Conn) {
	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
}

func (r OptionsResponse) StatusLine() string {
	return fmt.Sprintf("%v %v %v\n", r.Version.ToString(), 204, "No Content")
}

func (r OptionsResponse) Headers() header.Headers {
	headers := header.Headers{}
	headers = append(headers, &header.Header{FieldName: "Allow", FieldValue: "GET, POST, HEAD, OPTIONS"})
	return headers
}

func (r OptionsResponse) Response(conn net.Conn) {
	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
}

func (r HeadResponse) Body() []byte {
	return echoBody(r.Request)
}

func echoBody(r request.Request) []byte {

	headerStrs := []string{}
	for _, h := range r.Headers {
		headerStrs = append(headerStrs, h.ToString())
	}

	body, _ := json.Marshal(map[string]interface{}{
		"method":         r.StartLine.Method.ToString(),
		"request_target": r.StartLine.RequestTarget,
		"version":        r.StartLine.Version.ToString(),
		"headers":        headerStrs,
		"body":           string(r.Body),
	}) // TODO err
	return body
}
