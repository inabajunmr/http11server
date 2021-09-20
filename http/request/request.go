package request

import (
	"bufio"
	"strings"

	"github.com/inabajunmr/http11server/http"
	"github.com/inabajunmr/http11server/http/header"
)

type Request struct {
	StartLine StartLine
	Headers   header.Headers
	Body      []byte
}

func ParseRequest(reader *bufio.Reader) (*Request, error) {
	startLine, err := ParseStartLine(readLine(reader))
	if err != nil {
		return nil, err
	}
	headers, err := readHeaders(reader)
	if err != nil {
		return nil, err
	}
	body, err := readBody(reader, *headers)
	if err != nil {
		return nil, err
	}

	return &Request{StartLine: *startLine, Headers: *headers, Body: body}, nil
}

func readHeaders(reader *bufio.Reader) (*header.Headers, error) {
	headers := header.Headers{}
	for {
		line := readLine(reader)
		if line == "" {
			// next is request body...
			return &headers, nil
		}
		h, err := header.ParseHeader(line)
		if err != nil {
			switch err.(type) {
			case *http.HTTPError:
				return nil, err
			default:
				continue
			}
		}

		headers = append(headers, h)
	}

}

func readBody(reader *bufio.Reader, headers header.Headers) ([]byte, error) {
	if len(headers.GetTransferEncodings()) == 0 {
		length, err := headers.GetContentLength()
		if err != nil {
			return nil, err
		}
		var body = make([]byte, length)
		l, err := reader.Read(body) // TODO if body is shorter than length?
		if l != length {
			return nil, &http.HTTPError{Msg: "Content-Length and real body size are different.", Status: 400}
		}
		return body, nil
	} else {
		if headers.IsChunkedTransferEncoding() {
			return nil, nil // TODO
		} else {
			return nil, &http.HTTPError{Msg: "Transfer-Encoding is invalid.", Status: 400}
		}
	}
}

func readLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.Trim(line, "\r\n")
}
