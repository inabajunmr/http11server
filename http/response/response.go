package response

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/inabajunmr/http11server/http"
	"github.com/inabajunmr/http11server/http/header"
	"github.com/inabajunmr/http11server/http/request"
)

type Response interface {
	Response(conn net.Conn) error
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

type Echo struct {
	Method        string   `xml:"method"`
	RequestTarget string   `xml:"request_target"`
	Version       string   `xml:"version"`
	Headers       []string `xml:"headers"`
	Body          string   `xml:"body"`
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
	b, _ := r.Body()
	headers = append(headers, &header.Header{FieldName: "Content-Length", FieldValue: strconv.Itoa(len(b))})
	headers = append(headers, &header.Header{FieldName: "Date", FieldValue: time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")})
	for _, ae := range r.Request.Headers.GetAcceptEncodings() {
		if ae.Coding == header.CONTENT_CODING_GZIP {
			headers = append(headers, &header.Header{FieldName: "Content-Encoding", FieldValue: "gzip"})
			break
		} else if ae.Coding == header.CONTENT_CODING_IDENTITY {
			break
		}
	}
	return headers
}

func (r EchoResponse) Body() ([]byte, error) {
	return echoBody(r.Request)
}

func (r EchoResponse) Response(conn net.Conn) error {

	b, err := r.Body()
	if err != nil {
		return err
	}

	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
	conn.Write([]byte(b))
	return nil
}

func (r HeadResponse) StatusLine() string {
	return fmt.Sprintf("%v %v %v\n", r.Version.ToString(), r.StatusCode, r.ReasonPhrase)
}

func (r HeadResponse) Headers() header.Headers {
	headers := header.Headers{}
	b, _ := r.Body()
	headers = append(headers, &header.Header{FieldName: "Content-Length", FieldValue: strconv.Itoa(len(b))})
	headers = append(headers, &header.Header{FieldName: "Date", FieldValue: time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")})
	for _, ae := range r.Request.Headers.GetAcceptEncodings() {
		if ae.Coding == header.CONTENT_CODING_GZIP {
			headers = append(headers, &header.Header{FieldName: "Content-Encoding", FieldValue: "gzip"})
		}
	}
	return headers
}

func (r HeadResponse) Response(conn net.Conn) error {
	conn.Write([]byte(r.StatusLine()))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
	return nil
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

func (r HeadResponse) Body() ([]byte, error) {
	return echoBody(r.Request)
}

func echoBody(r request.Request) ([]byte, error) {

	headerStrs := []string{}
	for _, h := range r.Headers {
		headerStrs = append(headerStrs, h.ToString())
	}

	accepts := r.Headers.GetAccept()
	if len(accepts) == 0 {
		// default is json
		j, err := json.Marshal(map[string]interface{}{
			"method":         r.StartLine.Method.ToString(),
			"request_target": r.StartLine.RequestTarget,
			"version":        r.StartLine.Version.ToString(),
			"headers":        headerStrs,
			"body":           string(r.Body),
		})

		if err != nil {
			return nil, err
		}

		return compress(j, r.Headers.GetAcceptEncodings()), nil
	}

	for _, a := range accepts {
		if a.Type == "application" && a.SubType == "json" {
			j, err := json.Marshal(map[string]interface{}{
				"method":         r.StartLine.Method.ToString(),
				"request_target": r.StartLine.RequestTarget,
				"version":        r.StartLine.Version.ToString(),
				"headers":        headerStrs,
				"body":           string(r.Body),
			})

			if err != nil {
				return nil, err
			}

			return compress(j, r.Headers.GetAcceptEncodings()), nil
		} else if a.Type == "application" && a.SubType == "xml" {
			v := &Echo{Method: r.StartLine.Method.ToString(),
				RequestTarget: r.StartLine.RequestTarget,
				Version:       r.StartLine.Version.ToString(),
				Headers:       headerStrs,
				Body:          string(r.Body)}
			xml, err := xml.MarshalIndent(v, "", " ")
			if err != nil {
				return nil, err
			}
			return compress(xml, r.Headers.GetAcceptEncodings()), nil
		}
	}

	return nil, &http.HTTPError{Status: 406, Msg: "Not Acceptable"}
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
