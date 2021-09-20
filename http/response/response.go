package response

import (
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
	headers = append(headers, &header.Header{FieldName: "Connection", FieldValue: "close"})
	headers = append(headers, &header.Header{FieldName: "Content-Length", FieldValue: strconv.Itoa(len(r.Body()))})
	return headers
}

func (r EchoResponse) Body() []byte {
	return []byte(fmt.Sprintf(`<html><head><title>HTTP/1.1</title></head>
<body>
    <h1>Echo</h1>
    <h2>Method</h2>
	<span>%v<span>
    <h2>Target</h2>
	<span>%v<span>
    <h2>Version</h2>
	<span>%v<span>
    <h2>Headers</h2>
	<span>%v<span>
	<h2>Body</h2>
	<span>%v<span>
	<h2>POST<h2>
	<form action="/" method="post">
	<textarea name="text" rows="5">Yeah</textarea>
	<div>
	  <button>Post</button>
	</div>
  </form>	
</body>
<htmL>`, r.Request.StartLine.Method.ToString(),
		r.Request.StartLine.RequestTarget,
		r.Request.StartLine.Version.ToString(),
		r.Request.Headers.ToString(),
		string(r.Request.Body)))
}

func (r EchoResponse) Response(conn net.Conn) {
	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
	conn.Write([]byte(r.Body()))
}
