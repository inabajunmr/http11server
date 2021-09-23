package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	go Serve()
	m.Run()
	Stop()
}

func TestGet(t *testing.T) {
	resp, err := http.Get("http://localhost:80")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", "HOST: localhost:80", "ACCEPT-ENCODING: gzip")
}

func TestPost(t *testing.T) {
	resp, err := http.Post("http://localhost:80",
		"application/x-www-form-urlencoded",
		strings.NewReader("aaaaabbbbbccccc"))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertResponse(t, b, "aaaaabbbbbccccc", "POST", "/", "HTTP/1.1",
		"CONTENT-TYPE: application/x-www-form-urlencoded", "USER-AGENT: Go-http-client/1.1",
		"CONTENT-LENGTH: 15", "HOST: localhost:80", "ACCEPT-ENCODING: gzip")
}

func TestPost_Chunkded(t *testing.T) {
	rd, wr := io.Pipe()

	req, err := http.NewRequest("POST", "http://localhost:80", rd)
	if err != nil {
		t.Fatal(err)
	}
	req.TransferEncoding = []string{"chunked"}

	go func() {
		wr.Write([]byte("hello"))
		wr.Write([]byte("hello"))
		wr.Write([]byte("hello"))
		wr.Close()
	}()
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	assertResponse(t, b, "hellohellohello", "POST", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1",
		"HOST: localhost:80", "ACCEPT-ENCODING: gzip",
		"TRANSFER-ENCODING: chunked")
}

func assertResponse(t *testing.T, response []byte, expectedBody string, expectedMethod string, expectedRequestTarget string, expectedVersion string, expectedHeaders ...string) {
	res := map[string]interface{}{}
	json.Unmarshal(response, &res)

	if res["body"] != expectedBody {
		t.Errorf("Unexpected body:%v.", res["body"])
	}
	assertHeaders(t, res["headers"].([]interface{}), expectedHeaders...,
	)
	if res["method"] != expectedMethod {
		t.Errorf("Unexpected method:%v.", res["method"])
	}
	if res["request_target"] != expectedRequestTarget {
		t.Errorf("Unexpected request_target:%v.", res["request_target"])
	}
	if res["version"] != expectedVersion {
		t.Errorf("Unexpected version:%v.", res["version"])
	}
}

func assertHeaders(t *testing.T, actual []interface{}, expected ...string) {
	if len(expected) != len(actual) {
		t.Errorf("Unexpected header:%v.", actual)
	}
	for _, e := range expected {
		if !containsHeader(e, actual) {
			t.Errorf("Unexpected header:%v. expected:%v.", actual, e)
		}
	}

}

func containsHeader(header string, headers []interface{}) bool {
	for _, h := range headers {
		if h == header {
			return true
		}
	}
	return false
}
