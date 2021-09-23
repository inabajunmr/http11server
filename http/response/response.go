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

type EchoResponse struct {
	Version      http.HTTPVersion
	StatusCode   int
	ReasonPhrase string
	Request      request.Request
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

	headerStrs := []string{}
	for _, h := range r.Request.Headers {
		headerStrs = append(headerStrs, h.ToString())
	}

	body, _ := json.Marshal(map[string]interface{}{
		"method":         r.Request.StartLine.Method.ToString(),
		"request_target": r.Request.StartLine.RequestTarget,
		"version":        r.Request.StartLine.Version.ToString(),
		"headers":        headerStrs,
		"body":           string(r.Request.Body),
	}) // TODO err
	return body
}

func (r EchoResponse) Response(conn net.Conn) {
	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
	conn.Write([]byte(r.Body()))
}
