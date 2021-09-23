package server

import (
	"bufio"
	"fmt"
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
	for {
		log.Println("Listen")
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		reader := bufio.NewReader(conn)
		req, err := request.ParseRequest(reader)
		if err != nil {
			switch httpErr := err.(type) {
			case *http.HTTPError:
				res := &response.EchoResponse{Version: http.HTTP11, StatusCode: httpErr.Status, ReasonPhrase: httpErr.Msg}
				res.Response(conn)
				conn.Close()
			default:
				res := &response.EchoResponse{Version: http.HTTP11, StatusCode: 503, ReasonPhrase: "Service Unavailable"}
				res.Response(conn)
				conn.Close()
			}
			continue
		}

		log.Println(req.StartLine.ToString())
		log.Println(req.Headers.ToString())
		log.Println(string(req.Body))

		res := &response.EchoResponse{Version: http.HTTP11, StatusCode: 200, ReasonPhrase: "OK", Request: *req}
		res.Response(conn)
		conn.Close()
	}
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
