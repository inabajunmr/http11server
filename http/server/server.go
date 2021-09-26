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

var PORT int

func Serve(port int) {

	service := fmt.Sprintf(":%v", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err = net.ListenTCP("tcp4", tcpAddr)
	checkError(err)
	PORT = listener.Addr().(*net.TCPAddr).Port
	log.Printf("LISTEN PORT:%v", PORT)

	for {
		conn, _ := listener.AcceptTCP()
		reader := bufio.NewReader(conn)
		go processRequest(conn, reader)
	}
}

func processRequest(conn net.Conn, reader *bufio.Reader) {
	for {
		req, err := request.ParseRequest(reader)
		if err != nil {
			if handleError(conn, err) {
				log.Println("Close")
				conn.Close()
				return
			}
			continue
		}

		log.Println(req.StartLine.ToString())
		log.Println(req.Headers.ToString())
		log.Println(string(req.Body))

		err = response.GetResponse(*req).Response(conn)
		if err != nil {
			if handleError(conn, err) {
				log.Println("Close")
				conn.Close()
				return
			}
		}

		if req.Headers.IsConnectionClose() {
			log.Println("Close")
			conn.Close()
			return
		} else {
			continue
		}
	}
}

func Stop() {
	listener.Close()
}

func handleError(conn net.Conn, err error) bool {
	switch httpErr := err.(type) {
	case *http.HTTPError:
		res := &response.EchoResponse{Version: http.HTTP11, StatusCode: httpErr.Status, ReasonPhrase: httpErr.Msg}
		res.Response(conn)
	case *http.WaitRequestError:
		return false
	default:
		if err == io.EOF {
			return true
		}
		res := &response.EchoResponse{Version: http.HTTP11, StatusCode: 503, ReasonPhrase: "Service Unavailable"}
		res.Response(conn)
	}
	return false
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
