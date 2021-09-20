package header

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/inabajunmr/http11server/http"
)

type Header struct {
	FieldName  string
	FieldValue string
}

type Headers []*Header

type TransferEncoding int

const (
	CHUNKED TransferEncoding = iota
	COMPRESS
	DEFLATE
	GZIP
	IDENTITY
)

func (hs Headers) ToString() string {
	v := ""
	for _, h := range hs {
		v += fmt.Sprintf("%v: %v\n", h.FieldName, h.FieldValue)
	}
	return v
}

func (h Headers) GetContentLength() (int, error) {
	filtered := h.filter("CONTENT-LENGTH")
	if len(filtered) >= 2 {
		return 0, &http.HTTPError{Msg: "Multiple Content-Length is not allowed."}
	}
	if len(filtered) == 0 {
		return 0, nil
	}
	length, err := strconv.Atoi(filtered[0].FieldValue)
	if err != nil {
		return 0, &http.HTTPError{Msg: fmt.Sprintf("Content-Length:%v is not number.", filtered[0].FieldName)}
	}
	return length, nil
}

func (h Headers) GetTransferEncodings() []TransferEncoding {
	var tes = []TransferEncoding{}
	filtered := h.filter("TRANSFER-ENCODING")
	for _, header := range filtered {
		for _, v := range strings.Split(header.FieldName, ",") {
			tes = append(tes, getTransferEncoding(v))
		}
	}
	return tes
}

func (h Headers) IsChunkedTransferEncoding() bool {
	for _, te := range h.GetTransferEncodings() {
		if te == CHUNKED {
			return true
		}
	}
	return false
}

func (h Headers) filter(key string) Headers {
	var headers = Headers{}
	for _, header := range h {
		if header.FieldName == key {
			headers = append(headers, header)
		}
	}
	return headers
}

func getTransferEncoding(v string) TransferEncoding {
	switch v {
	case "chunked":
		return CHUNKED
	case "compress":
		return COMPRESS
	case "deflate":
		return DEFLATE
	case "gzip":
		return GZIP
	case "identity":
		return IDENTITY
	}
	return 0 // TODO
}

func ParseHeader(line string) (*Header, error) {
	// TODO this method don't allow obs-fold if it's in message/http container.
	l := strings.SplitN(line, ":", 2)
	if len(l) != 2 {
		return nil, &HeaderParserError{Msg: fmt.Sprintf("Header line:%v don't has ':'.", line)}
	}

	if strings.HasSuffix(l[0], " ") || strings.HasSuffix(l[0], "\t") {
		return nil, &http.HTTPError{Msg: "Header field name don't allow space before colon.", Status: 400}

	}

	if !validateFieldName(l[0]) {
		return nil, &HeaderParserError{Msg: fmt.Sprintf("Header line:%v don't has ':'.", line)}
	}
	if len(l[1]) == 0 {
		return nil, &HeaderParserError{Msg: fmt.Sprintf("Header line:%v don't has ':'.", line)}
	}

	return &Header{strings.ToUpper(l[0]), strings.TrimSpace(l[1])}, nil

}

func validateFieldName(fieldName string) bool {
	if len(fieldName) == 0 {
		return false
	}
	return isVisibleASCII(fieldName) // l[0] never conains ':'

}

func isVisibleASCII(s string) bool {
	for _, c := range s {
		if c < '\u0021' || c > '\u007E' {
			return false
		}
	}
	return true
}
