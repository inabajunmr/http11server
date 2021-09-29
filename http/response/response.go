package response

import (
	"bytes"
	"compress/gzip"
	"net"

	"github.com/inabajunmr/http11server/http"
	"github.com/inabajunmr/http11server/http/header"
	"github.com/inabajunmr/http11server/http/request"
)

type Response interface {
	Response(conn net.Conn) error
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

func compress(body []byte, acceptEncodings []header.AcceptEncoding) []byte {
	for _, ae := range acceptEncodings {
		if ae.Coding == header.CONTENT_CODING_GZIP {
			var b bytes.Buffer
			writer := gzip.NewWriter(&b)
			writer.Write(body) // TODO error
			writer.Flush()
			writer.Close()
			return b.Bytes()
		} else if ae.Coding == header.CONTENT_CODING_IDENTITY {
			return body
		}
	}
	return []byte{} // TODO error
}
