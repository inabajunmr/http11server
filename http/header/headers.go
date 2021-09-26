package header

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/inabajunmr/http11server/http"
)

type Headers []*Header

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
	exp := h.filter("EXPECT")
	if len(exp) != 0 && strings.ToUpper(exp[0].FieldValue) != "100-CONTINUE" {
		return &http.HTTPError{Status: 417, Msg: "Expectation Failed"}
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

func (h Headers) GetContentType() ContentType {
	c := h.filter("CONTENT-TYPE")
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

func (h Headers) GetContentLocation() *string {
	c := h.filter("CONTENT-LOCATION")
	if len(c) == 0 {
		return nil
	}

	return &c[0].FieldValue
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
