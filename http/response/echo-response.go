package response

import (
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

type EchoResponse struct {
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

func (r EchoResponse) StatusLine(code int) string {
	return fmt.Sprintf("%v %v %v\n", r.Version.ToString(), code, r.ReasonPhrase)
}

// body is not ranged body. If request has Range header, it's not splited yet.
func echoResponseHeader(r request.Request, body []byte) header.Headers {
	headers := header.Headers{}
	headers = append(headers, &header.Header{FieldName: "Date", FieldValue: time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")})
	headers = append(headers, &header.Header{FieldName: "Vary", FieldValue: "accept-encoding, accept"})
	headers = append(headers, &header.Header{FieldName: "Accept-Range", FieldValue: "bytes"})

	for _, ae := range r.Headers.GetAcceptEncodings() {
		if ae.Coding == header.CONTENT_CODING_GZIP {
			headers = append(headers, &header.Header{FieldName: "Content-Encoding", FieldValue: "gzip"})
			break
		} else if ae.Coding == header.CONTENT_CODING_IDENTITY {
			break
		}
	}

	ranges, _ := r.Headers.GetRanges()
	if len(ranges) >= 1 { // multiple ranges is unsuppored yet
		s, e := getRange(ranges[0].Start, ranges[0].End, body)
		if len(body) < e {
			// Range header specify bigger end value than body length
			headers = append(headers, &header.Header{FieldName: "Content-Range",
				FieldValue: fmt.Sprintf("bytes */%v", len(body))})
			headers = append(headers, &header.Header{FieldName: "Content-Length", FieldValue: "0"})
		} else {
			headers = append(headers, &header.Header{FieldName: "Content-Range",
				FieldValue: fmt.Sprintf("bytes %v-%v/%v", s, e, len(body))})
			headers = append(headers, &header.Header{FieldName: "Content-Length", FieldValue: strconv.Itoa(e - s)})
		}
	} else {
		headers = append(headers, &header.Header{FieldName: "Content-Length", FieldValue: strconv.Itoa(len(body))})
	}
	return headers
}

func getRange(start *int, end *int, body []byte) (int, int) {
	rs := start
	if rs == nil {
		s := len(body) - *end
		rs = &s
		length := len(body)
		return *rs, length
	}

	re := end
	if re == nil {
		length := len(body)
		re = &length
	}

	return *rs, *re
}

func (r EchoResponse) Headers() header.Headers {
	b, _ := r.Body()
	return echoResponseHeader(r.Request, b)
}

func (r EchoResponse) Body() ([]byte, error) {
	return echoBody(r.Request)
}

func (r EchoResponse) Response(conn net.Conn) error {

	fullBody, err := r.Body()
	if err != nil {
		return err
	}

	ranges, _ := r.Request.Headers.GetRanges()

	var b []byte
	status := 200
	if len(ranges) >= 1 {
		s, e := getRange(ranges[0].Start, ranges[0].End, fullBody)
		b = fullBody[s:e]

		if len(fullBody) < e {
			status = 416
			conn.Write([]byte(r.StatusLine(status)))
			conn.Write([]byte(r.Headers().ToString()))
			conn.Write([]byte("\n"))
			return nil
		}
		status = 206
	} else {
		b = fullBody
	}

	conn.Write([]byte(r.StatusLine(status)))
	conn.Write([]byte(r.Headers().ToString()))
	conn.Write([]byte("\n"))
	conn.Write(b)

	return nil
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
