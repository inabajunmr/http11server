package response

import (
	"fmt"
	"net"

	"github.com/inabajunmr/http11server/http"
	"github.com/inabajunmr/http11server/http/header"
	"github.com/inabajunmr/http11server/http/request"
)

type HeadResponse struct {
	Version      http.HTTPVersion
	StatusCode   int
	ReasonPhrase string
	Request      request.Request
}

func (r HeadResponse) StatusLine() string {
	return fmt.Sprintf("%v %v %v\n", r.Version.ToString(), r.StatusCode, r.ReasonPhrase)
}

func (r HeadResponse) Headers() header.Headers {
	return echoResponseHeader(r.Request, r.Request.Body)
}

func (r HeadResponse) Response(conn net.Conn) error {
	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
	return nil
}

func (r HeadResponse) Body() ([]byte, error) {
	return echoBody(r.Request)
}
