package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/inabajunmr/http11server/http"
	"github.com/inabajunmr/http11server/http/request"
	"github.com/inabajunmr/http11server/http/response"
)

var listener *net.TCPListener

func Serve() {

	service := ":80"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err = net.ListenTCP("tcp4", tcpAddr)
	checkError(err)
	conn, _ := listener.AcceptTCP()
	reader := bufio.NewReader(conn)

	for {
		if processRequest(conn, reader) {
			log.Println("Close")
			conn.Close()
			conn, _ = listener.AcceptTCP()
			reader = bufio.NewReader(conn)
		}
	}
}

func processRequest(conn net.Conn, reader *bufio.Reader) bool {
	req, err := request.ParseRequest(reader)
	if err != nil {
		switch httpErr := err.(type) {
		case *http.HTTPError:
			res := &response.EchoResponse{Version: http.HTTP11, StatusCode: httpErr.Status, ReasonPhrase: httpErr.Msg}
			res.Response(conn)
		default:
			if err == io.EOF {
				return true
			}
			res := &response.EchoResponse{Version: http.HTTP11, StatusCode: 503, ReasonPhrase: "Service Unavailable"}
			res.Response(conn)
		}
		return false
	}

	log.Println(req.StartLine.ToString())
	log.Println(req.Headers.ToString())
	log.Println(string(req.Body))

	res := &response.EchoResponse{Version: http.HTTP11, StatusCode: 200, ReasonPhrase: "OK", Request: *req}
	res.Response(conn)
	return req.Headers.IsConnectionClose()
}

func Stop() {
	listener.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
