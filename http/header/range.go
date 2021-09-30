package header

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/inabajunmr/http11server/http"
)

type Range struct {
	Start *int
	End   *int
}

func ParseRange(val string) ([]Range, error) {
	trim := strings.TrimSpace(val)
	if !strings.HasPrefix(trim, "bytes=") {
		return nil, &http.HTTPError{Status: 400, Msg: "Only bytes bytes-unit supported."}
	}
	vals := strings.Split(strings.TrimLeft(trim, "bytes="), ",")
	ranges := []Range{}
	for _, v := range vals {
		trimV := strings.TrimSpace(v)
		if strings.HasPrefix(trimV, "-") {
			end, err := strconv.Atoi(strings.TrimLeft(trimV, "-"))
			if err != nil {
				return nil, &http.HTTPError{Status: 400, Msg: fmt.Sprintf("Invalid Range header: %v.", val)}
			}
			ranges = append(ranges, Range{End: &end})
			continue
		} else if strings.HasSuffix(trimV, "-") {
			start, err := strconv.Atoi(strings.TrimRight(trimV, "-"))
			if err != nil {
				return nil, &http.HTTPError{Status: 400, Msg: fmt.Sprintf("Invalid Range header: %v.", val)}
			}
			ranges = append(ranges, Range{Start: &start})
			continue
		} else {
			se := strings.Split(trimV, "-")
			if len(se) != 2 {
				return nil, &http.HTTPError{Status: 400, Msg: fmt.Sprintf("Invalid Range header: %v.", val)}
			}
			start, err := strconv.Atoi(se[0])
			if err != nil {
				return nil, &http.HTTPError{Status: 400, Msg: fmt.Sprintf("Invalid Range header: %v.", val)}
			}
			end, err := strconv.Atoi(se[1])
			if err != nil {
				return nil, &http.HTTPError{Status: 400, Msg: fmt.Sprintf("Invalid Range header: %v.", val)}
			}
			ranges = append(ranges, Range{Start: &start, End: &end})
			continue
		}
	}
	return ranges, nil
}
