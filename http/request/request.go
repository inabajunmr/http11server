package request

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"strconv"
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

	l, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	if *l == "" {
		return nil, &http.WaitRequestError{Msg: "Response no request yet."}
	}

	startLine, err := ParseStartLine(*l)
	if err != nil {
		return nil, err
	}
	headers, err := readHeaders(reader)
	if err == io.EOF {
		return &Request{StartLine: *startLine, Headers: *headers, Body: nil}, nil
	}
	if headers.Validate() != nil {
		return nil, err
	}

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
		line, err := readLine(reader)
		if err == io.EOF {
			return &headers, err
		}
		if err != nil {
			return nil, err
		}
		if *line == "" {
			// next is request body...
			return &headers, nil
		}
		h, err := header.ParseHeader(*line)
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
		l, _ := reader.Read(body) // TODO if body is shorter than length?
		if l != length {
			return nil, &http.HTTPError{Msg: "Content-Length and real body size are different.", Status: 400}
		}
		return decompress(body, headers), nil
	} else {
		if headers.IsChunkedTransferEncoding() {
			// TODO trailer
			// TODO compress
			b := parseChunkBody(reader)
			t := headers.GetCompressType()
			log.Println(t)
			switch t {
			case header.TRANSFER_ENCODING_GZIP:
				// TODO untested because I can't find HTTP Client send 'Transfer-Encoding: gzip, chunked
				gr, _ := gzip.NewReader(bytes.NewReader(b)) // TODO
				unzip, _ := ioutil.ReadAll(gr)              // TODO
				return unzip, nil
			default:
				return decompress(b, headers), nil

			}

		} else {
			return nil, &http.HTTPError{Msg: "Transfer-Encoding is invalid.", Status: 400}
		}
	}
}

func decompress(b []byte, headers header.Headers) []byte {
	ces := headers.GetContentEncodings()
	log.Println(ces)
	for _, ce := range ces {
		switch ce {
		case header.CONTENT_CODING_GZIP: // TODO defrate, compress
			br := bytes.NewReader(b)
			gr, _ := gzip.NewReader(br)
			b, _ = ioutil.ReadAll(gr)
			log.Println(string(b))
		case header.CONTENT_CODING_IDENTITY:
			// NOP
		}
	}
	return b
}

func parseChunkBody(reader *bufio.Reader) []byte {
	chunks := []byte{}
	for {
		l, err := readLine(reader) // TODO error
		if err == io.EOF {
			return chunks

		}

		chunkSize := ParseChunkSize(*l)
		if chunkSize == 0 {
			return chunks
		}
		var chunk = make([]byte, chunkSize)
		reader.Read(chunk)
		readLine(reader) // skip to next line
		chunks = append(chunks, chunk...)
	}
}

func ParseChunkSize(line string) int64 {
	v, _ := strconv.ParseInt(line, 16, 64) // TODO error
	// TODO chunk-ext
	return v
}

func readLine(reader *bufio.Reader) (*string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	t := strings.Trim(line, "\r\n")
	return &t, nil
}
