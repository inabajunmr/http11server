package server

import (
	"bytes"
	"compress/gzip"
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

func TestGet_ConnectionClosed(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:80", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Connection", "close")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", "HOST: localhost:80", "ACCEPT-ENCODING: gzip", "CONNECTION: close")

}

func TestGet_KeepAlive(t *testing.T) {
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

	resp, err = http.Get("http://localhost:80")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", "HOST: localhost:80", "ACCEPT-ENCODING: gzip")

	resp, err = http.Get("http://localhost:80")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
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

func TestPost_ContentEncodingGzip(t *testing.T) {
	// create gziped request
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	writer.Write([]byte("hellohellohello"))
	writer.Close()
	gzip := buffer.Bytes()

	req, err := http.NewRequest("POST", "http://localhost:80",
		bytes.NewReader(gzip))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Encoding", "gzip")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	assertResponse(t, b, "hellohellohello", "POST", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", "CONTENT-LENGTH: 31",
		"HOST: localhost:80", "CONTENT-ENCODING: gzip", "ACCEPT-ENCODING: gzip")
}

func TestPost_ContentEncodingGzipGzip(t *testing.T) {
	// create gziped request
	var buffer1 bytes.Buffer
	writer := gzip.NewWriter(&buffer1)
	writer.Write([]byte("hellohellohello"))
	writer.Close()
	gziped := buffer1.Bytes()

	var buffer2 bytes.Buffer
	writer = gzip.NewWriter(&buffer2)
	writer.Write(gziped)
	writer.Close()
	gziped = buffer2.Bytes()

	req, err := http.NewRequest("POST", "http://localhost:80",
		bytes.NewReader(gziped))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Encoding", "gzip, gzip")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	assertResponse(t, b, "hellohellohello", "POST", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", "CONTENT-LENGTH: 50",
		"HOST: localhost:80", "CONTENT-ENCODING: gzip, gzip", "ACCEPT-ENCODING: gzip")
}

func TestPost_Chunkded(t *testing.T) {
	rd, wr := io.Pipe()
	defer rd.Close()

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
