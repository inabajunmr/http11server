package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/inabajunmr/http11server/http/response"
)

func TestMain(m *testing.M) {
	go Serve(0)
	m.Run()
}

func addr() string {
	log.Printf("http://localhost:%v", PORT)
	return fmt.Sprintf("http://localhost:%v", PORT)
}

func TestGet(t *testing.T) {
	resp, err := http.Get(addr())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertJsonResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip")

	if resp.Header.Get("Date") == "" {
		t.Errorf("Missing Date header.")
	}
	if resp.Header.Get("Vary") != "accept-encoding, accept" {
		t.Errorf("Unexpected Vary Header: %v.", resp.Header.Get("Vary"))
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected status: %v.", resp.StatusCode)
	}
}

func TestGet_ConnectionClosed(t *testing.T) {
	req, err := http.NewRequest("GET", addr(), nil)
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

	assertJsonResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip", "CONNECTION: close")
}

func TestGet_KeepAlive(t *testing.T) {
	resp, err := http.Get(addr())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertJsonResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip")

	resp, err = http.Get(addr())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertJsonResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip")

	resp, err = http.Get(addr())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assertJsonResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip")
}

func TestGet_AcceptJson(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assertJsonResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip", "ACCEPT: application/json")
}

func TestGet_AcceptXml(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Accept", "application/xml")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println(string(b))
	assertXmlResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip", "ACCEPT: application/xml")
}

func TestGet_AcceptXmlAndJson(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Accept", "application/json; q=0.5, application/xml")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println(string(b))
	assertXmlResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip", "ACCEPT: application/json; q=0.5, application/xml")
}

func TestGet_AccceptEncodingGzip(t *testing.T) {
	req, err := http.NewRequest("GET", addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Accept-Encoding", "gzip")

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

	gzipReader, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	unzip, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		t.Fatal(err)
	}

	assertJsonResponse(t, unzip, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip")
}

func TestGet_AccceptEncodingIdentityGzip(t *testing.T) {
	req, err := http.NewRequest("GET", addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Accept-Encoding", "gzip; q=0.5, identity")

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

	assertJsonResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip; q=0.5, identity")
}

func TestGet_Range1(t *testing.T) {
	req, err := http.NewRequest("GET", addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Range", "bytes=0-")

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

	assertJsonResponse(t, b, "", "GET", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", fmt.Sprintf("HOST: localhost:%v", PORT), "RANGE: bytes=0-")
	if resp.StatusCode != 206 {
		t.Errorf("Unexpected status: %v.", resp.StatusCode)
	}
}

func TestGet_Range2(t *testing.T) {
	req, err := http.NewRequest("GET", addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Range", "bytes=1-10")

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

	if string(b) != "\"body\":\"\"" {
		t.Errorf("Unexpected Body: %v", string(b))
	}
	if resp.Header.Get("Content-Range") != "bytes 1-10/157" {
		t.Errorf("Unexpected Content-Range: %v", resp.Header.Get("Content-Range"))
	}
	if resp.Header.Get("Content-Length") != strconv.Itoa(len(b)) {
		t.Errorf("Unexpected Content-Length: %v", resp.Header.Get("Content-Length"))
	}
	if resp.StatusCode != 206 {
		t.Errorf("Unexpected status: %v.", resp.StatusCode)
	}
}

func TestGet_Range3(t *testing.T) {
	req, err := http.NewRequest("GET", addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Range", "bytes=-10")

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

	if string(b) != "HTTP/1.1\"}" {
		t.Errorf("Unexpected Body: %v", string(b))
	}
	if resp.Header.Get("Content-Range") != "bytes 146-156/156" {
		t.Errorf("Unexpected Content-Range: %v", resp.Header.Get("Content-Range"))
	}
	if resp.Header.Get("Content-Length") != strconv.Itoa(len(b)) {
		t.Errorf("Unexpected Content-Length: %v", resp.Header.Get("Content-Length"))
	}
	if resp.StatusCode != 206 {
		t.Errorf("Unexpected status: %v.", resp.StatusCode)
	}
}

func TestGet_Range4(t *testing.T) {
	req, err := http.NewRequest("GET", addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Range", "bytes=0-159")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Header.Get("Content-Range") != "bytes */158" {
		t.Errorf("Unexpected Content-Range: %v", resp.Header.Get("Content-Range"))
	}
	if resp.Header.Get("Content-Length") != "0" {
		t.Errorf("Unexpected Content-Length: %v", resp.Header.Get("Content-Length"))
	}
	if resp.StatusCode != 416 {
		t.Errorf("Unexpected status: %v.", resp.StatusCode)
	}
}

func TestHead(t *testing.T) {
	resp, err := http.Head(addr())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 0 {
		t.Error("Unexpected body.")
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected status: %v", resp.StatusCode)
	}
}

func TestOptions(t *testing.T) {
	req, err := http.NewRequest(http.MethodOptions, addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if len(b) != 0 {
		t.Error("Unexpected body.")
	}
	if resp.Header.Get("Allow") != "GET, POST, HEAD, OPTIONS" {
		t.Errorf("Unexpected header: %v.", resp.Header.Get("Allow"))
	}
	if resp.StatusCode != 204 {
		t.Errorf("Unexpected status: %v", resp.StatusCode)
	}
}

func TestPost(t *testing.T) {
	resp, err := http.Post(addr(),
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

	assertJsonResponse(t, b, "aaaaabbbbbccccc", "POST", "/", "HTTP/1.1",
		"CONTENT-TYPE: application/x-www-form-urlencoded", "USER-AGENT: Go-http-client/1.1",
		"CONTENT-LENGTH: 15", fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip")
}

func TestPost_ContentLocation(t *testing.T) {
	req, err := http.NewRequest("POST", addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Location", "http://example.com")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	if strings.HasPrefix(string(b), "<") {
		t.Errorf("Unexpected body: %v", string(b))
	}
}

func TestPost_ContentEncodingGzip(t *testing.T) {
	// create gziped request
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	writer.Write([]byte("hellohellohello"))
	writer.Close()
	gzip := buffer.Bytes()

	req, err := http.NewRequest("POST", addr(),
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

	assertJsonResponse(t, b, "hellohellohello", "POST", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", "CONTENT-LENGTH: 31",
		fmt.Sprintf("HOST: localhost:%v", PORT), "CONTENT-ENCODING: gzip", "ACCEPT-ENCODING: gzip")
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

	req, err := http.NewRequest("POST", addr(),
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

	assertJsonResponse(t, b, "hellohellohello", "POST", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1", "CONTENT-LENGTH: 50",
		fmt.Sprintf("HOST: localhost:%v", PORT), "CONTENT-ENCODING: gzip, gzip", "ACCEPT-ENCODING: gzip")
}

func TestPost_Chunkded(t *testing.T) {
	rd, wr := io.Pipe()
	defer rd.Close()

	req, err := http.NewRequest("POST", addr(), rd)
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

	assertJsonResponse(t, b, "hellohellohello", "POST", "/", "HTTP/1.1",
		"USER-AGENT: Go-http-client/1.1",
		fmt.Sprintf("HOST: localhost:%v", PORT), "ACCEPT-ENCODING: gzip",
		"TRANSFER-ENCODING: chunked")
}

func assertJsonResponse(t *testing.T, response []byte, expectedBody string, expectedMethod string, expectedRequestTarget string, expectedVersion string, expectedHeaders ...string) {
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

func assertXmlResponse(t *testing.T, resp []byte, expectedBody string, expectedMethod string, expectedRequestTarget string, expectedVersion string, expectedHeaders ...string) {
	res := response.Echo{}
	xml.Unmarshal(resp, &res)
	if res.Body != expectedBody {
		t.Errorf("Unexpected body:%v.", res.Body)
	}
	s := make([]interface{}, len(res.Headers))
	for i, v := range res.Headers {
		s[i] = v
	}
	assertHeaders(t, s, expectedHeaders...)
	if res.Method != expectedMethod {
		t.Errorf("Unexpected method:%v.", res.Method)
	}
	if res.RequestTarget != expectedRequestTarget {
		t.Errorf("Unexpected request_target:%v.", res.RequestTarget)
	}
	if res.Version != expectedVersion {
		t.Errorf("Unexpected version:%v.", res.Version)
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
