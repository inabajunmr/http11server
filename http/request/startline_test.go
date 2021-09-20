package request

import (
	"testing"

	"github.com/inabajunmr/http11server/http"
)

func TestParseStartLine(t *testing.T) {
	arg := "GET /aaa HTTP/1.1"
	result, err := ParseStartLine(arg)
	if err != nil {
		t.Errorf("Unexpected error %v.", err.Error())
	}
	if result.Method != GET {
		t.Errorf("Expected method is GET but %v.", result.Method.ToString())
	}
	if result.RequestTarget != "/aaa" {
		t.Errorf("Expected request target is /aaa but %v.", result.RequestTarget)
	}
	if result.Version != http.HTTP11 {
		t.Errorf("Expected request target is HTTP/1.1 but %v.", result.Version.ToString())
	}
}

func TestParseStartLine_HTTP09(t *testing.T) {
	arg := "YEAH /aaa"
	_, err := ParseStartLine(arg)
	if err == nil {
		t.Error("this request is not for HTTP/1.1")
	}
	if err.Error() != "this request is not for HTTP/1.1" {
		t.Errorf("Unexpected error: %v", err.Error())
	}
}

func TestParseStartLine_HTTP10(t *testing.T) {
	arg := "POST /aaa HTTP/1.0"
	_, err := ParseStartLine(arg)
	if err == nil {
		t.Error("Unexpected success.")
	}
	if err.Error() != "HTTP/1.0 is not supported HTTP version" {
		t.Errorf("Unexpected error: %v", err.Error())
	}
}

func TestParseStartLine_UnknownMethod(t *testing.T) {
	arg := "YEAH /aaa HTTP/1.1"
	_, err := ParseStartLine(arg)
	if err == nil {
		t.Error("Unexpected success.")
	}
	if err.Error() != "HTTP method YEAH is not implemented" {
		t.Errorf("Unexpected error: %v", err.Error())
	}
}
