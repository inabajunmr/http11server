package response

import (
	"fmt"
	"net"
	"time"

	"github.com/inabajunmr/http11server/http"
	"github.com/inabajunmr/http11server/http/header"
	"github.com/inabajunmr/http11server/http/request"
)

type OptionsResponse struct {
	Version http.HTTPVersion
	Request request.Request
}

func (r OptionsResponse) StatusLine() string {
	return fmt.Sprintf("%v %v %v\n", r.Version.ToString(), 204, "No Content")
}

func (r OptionsResponse) Headers() header.Headers {
	headers := header.Headers{}
	headers = append(headers, &header.Header{FieldName: "Allow", FieldValue: "GET, POST, HEAD, OPTIONS"})
	headers = append(headers, &header.Header{FieldName: "Date", FieldValue: time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")})
	return headers
}

func (r OptionsResponse) Response(conn net.Conn) error {
	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
	return nil
}
