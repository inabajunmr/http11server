package header

import (
	"fmt"
	"strings"

	"github.com/inabajunmr/http11server/http"
)

type Header struct {
	FieldName  string
	FieldValue string
}

func (h Header) ToString() string {
	return fmt.Sprintf("%v: %v", h.FieldName, h.FieldValue)
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
