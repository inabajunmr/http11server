package request

import (
	"bufio"
	"strings"
	"testing"

	"github.com/inabajunmr/http11server/http"
)

func TestParseRequest_Get(t *testing.T) {
	request :=
		`GET / HTTP/1.1
Header1: aaa
Header2: bbb ccc
`

	result, err := ParseRequest(bufio.NewReader((strings.NewReader(request))))
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if result.StartLine.Method != GET {
		t.Errorf("Unexpected method: %v", result.StartLine.Method)
	}
	if result.StartLine.Version != http.HTTP11 {
		t.Errorf("Unexpected version: %v", result.StartLine.Version)
	}
	if result.StartLine.RequestTarget != "/" {
		t.Errorf("Unexpected request target: %v", result.StartLine.RequestTarget)
	}
	if result.Headers[0].FieldName != "HEADER1" {
		t.Errorf("Unexpected header name: %v", result.Headers[0].FieldName)
	}
	if result.Headers[0].FieldValue != "aaa" {
		t.Errorf("Unexpected header name: %v", result.Headers[0].FieldValue)
	}
	if result.Headers[1].FieldName != "HEADER2" {
		t.Errorf("Unexpected header name: %v", result.Headers[1].FieldName)
	}
	if result.Headers[1].FieldValue != "bbb ccc" {
		t.Errorf("Unexpected header name: %v", result.Headers[1].FieldValue)
	}
	if len(result.Body) != 0 {
		t.Errorf("Unexpected body: %v", result.Body)
	}
}

func TestParseRequest_Post(t *testing.T) {
	request :=
		`POST / HTTP/1.1
Header1: aaa
Header2: bbb ccc
Content-Length: 12

aaaaa
bbbbb
`

	result, err := ParseRequest(bufio.NewReader((strings.NewReader(request))))
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if result.StartLine.Method != POST {
		t.Errorf("Unexpected method: %v", result.StartLine.Method)
	}
	if result.StartLine.Version != http.HTTP11 {
		t.Errorf("Unexpected version: %v", result.StartLine.Version)
	}
	if result.StartLine.RequestTarget != "/" {
		t.Errorf("Unexpected request target: %v", result.StartLine.RequestTarget)
	}
	if result.Headers[0].FieldName != "HEADER1" {
		t.Errorf("Unexpected header name: %v", result.Headers[0].FieldName)
	}
	if result.Headers[0].FieldValue != "aaa" {
		t.Errorf("Unexpected header name: %v", result.Headers[0].FieldValue)
	}
	if result.Headers[1].FieldName != "HEADER2" {
		t.Errorf("Unexpected header name: %v", result.Headers[1].FieldName)
	}
	if result.Headers[1].FieldValue != "bbb ccc" {
		t.Errorf("Unexpected header name: %v", result.Headers[1].FieldValue)
	}
	if string(result.Body) != "aaaaa\nbbbbb\n" {
		t.Errorf("Unexpected body: %v", string(result.Body))
	}
}

func TestParseRequest_Post_InvalidContentLength(t *testing.T) {
	request :=
		`POST / HTTP/1.1
Header1: aaa
Header2: bbb ccc
Content-Length: 14

aaaaa
bbbbb
`

	_, err := ParseRequest(bufio.NewReader((strings.NewReader(request))))
	if err.Error() != "Content-Length and real body size are different." {
		t.Errorf("Unexpected err: %v", err.Error())

	}
}
