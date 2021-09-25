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

func (h Header) ToString() string {
	return fmt.Sprintf("%v: %v", h.FieldName, h.FieldValue)
}

func (hs Headers) ToString() string {
	v := ""
	for _, h := range hs {
		v += h.ToString() + "\n"
	}
	return v
}

func (h Headers) Validate() error {
	if len(h.filter("HOST")) != 1 {
		return &http.HTTPError{Status: 400, Msg: "Request require only one Host header."}
	}
	return nil
}

func (h Headers) IsConnectionClose() bool {
	filtered := h.filter("CONNECTION")
	if len(filtered) == 0 {
		return false
	}
	return filtered[0].FieldValue == "close"
}

func (h Headers) GetContentLength() (int, error) {
	filtered := h.filter("CONTENT-LENGTH")
	if len(filtered) >= 2 {
		return 0, &http.HTTPError{Status: 400, Msg: "Multiple Content-Length is not allowed."}
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

func (h Headers) GetContentEncodings() []ContentCoding {
	ces := []ContentCoding{}
	filtered := h.filter("CONTENT-ENCODING")
	if len(filtered) == 0 {
		ces = append(ces, CONTENT_CODING_IDENTITY)
		return ces
	}

	vs := strings.Split(filtered[0].FieldValue, ",")
	for _, v := range vs {
		ces = append(ces, getContentCoding(strings.TrimSpace(v)))
	}
	return ces
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
		if te == TRANSFER_ENCODING_CHUNKED {
			return true
		}
	}
	return false
}

func (h Headers) GetCompressType() TransferEncoding {
	for _, te := range h.GetTransferEncodings() {
		if te != TRANSFER_ENCODING_CHUNKED {
			return te
		}
	}
	return TRANSFER_ENCODING_IDENTITY
}

type ContentType struct {
	contentType string
	subtype     string
}

func (h Headers) GetContentType() ContentType {
	c := h.filter("Content-Type")
	if len(c) == 0 {
		return ContentType{contentType: "application/octet-stream"}
	}

	s := strings.Split(c[0].FieldValue, ";")
	if len(s) >= 2 {
		return ContentType{contentType: s[0], subtype: s[1]}
	} else {
		return ContentType{contentType: s[0]}
	}
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
