package header

import (
	"sort"
	"strconv"
	"strings"
)

type Accept struct {
	Type             string
	SubType          string
	AcceptParameters []AcceptParameter
}

type AcceptParameter struct {
	Key   string
	Value string
}

func (a Accept) getWeight() float64 {
	for _, v := range a.AcceptParameters {
		if v.Key == "q" {
			q, err := strconv.ParseFloat(v.Value, 64)
			if err != nil {
				return 1
			}
			return q
		}
	}

	return 1
}

func ParseAccept(headerValue string) []Accept {

	accepts := []Accept{}

	for _, val := range strings.Split(headerValue, ",") {
		sp := strings.Split(val, ";")
		typ := strings.Split(sp[0], "/")

		mainType := strings.TrimSpace(typ[0])
		subType := strings.TrimSpace(typ[1])

		acceptParameters := []AcceptParameter{}
		for _, v := range sp[1:] {
			params := strings.Split(v, "=")
			if len(params) == 1 {
				acceptParameters = append(acceptParameters,
					AcceptParameter{Value: strings.TrimSpace(params[0])})
				continue
			}
			acceptParameters = append(acceptParameters,
				AcceptParameter{Key: strings.TrimSpace(params[0]), Value: strings.TrimSpace(params[1])})
		}
		accepts = append(accepts, Accept{Type: mainType, SubType: subType,
			AcceptParameters: acceptParameters})
	}

	sort.Slice(accepts, func(i, j int) bool {
		return accepts[i].getWeight() > accepts[j].getWeight()
	})

	return accepts
}
