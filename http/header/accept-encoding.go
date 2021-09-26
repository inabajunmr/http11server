package header

import (
	"sort"
	"strconv"
	"strings"
)

type AcceptEncoding struct {
	Coding ContentCoding
	Weight float64
}

func ParseAcceptEncoding(headerValue string) []AcceptEncoding {

	acceptEncodings := []AcceptEncoding{}
	wildCardQ := 0.0

	for _, val := range strings.Split(headerValue, ",") {
		sp := strings.Split(val, ";")
		ccstr := strings.TrimSpace(sp[0])
		cc := getContentCoding(ccstr)
		if len(sp) == 1 {
			if ccstr == "*" {
				wildCardQ = 1
				continue
			}
			acceptEncodings = append(acceptEncodings, AcceptEncoding{Coding: cc, Weight: 1})
			continue
		}

		params := strings.Split(sp[1], "=")
		q := 0.0
		if len(params) == 0 {
			q = 1
		} else {
			q, _ = strconv.ParseFloat(strings.TrimSpace(params[1]), 64)
		}
		if ccstr == "*" {
			wildCardQ = q
			continue
		}
		acceptEncodings = append(acceptEncodings, AcceptEncoding{Coding: cc, Weight: q})
	}

	if wildCardQ != 0 {
		allCodings := []ContentCoding{CONTENT_CODING_COMPRESS, CONTENT_CODING_DEFLATE, CONTENT_CODING_GZIP, CONTENT_CODING_IDENTITY}
		for _, c := range allCodings {
			if !containsCoding(acceptEncodings, c) {
				acceptEncodings = append(acceptEncodings, AcceptEncoding{Coding: c, Weight: wildCardQ})
			}
		}
	}

	sort.Slice(acceptEncodings, func(i, j int) bool {
		return acceptEncodings[i].Weight > acceptEncodings[j].Weight
	})

	without0 := []AcceptEncoding{}
	for _, a := range acceptEncodings {
		if a.Weight != 0 {
			without0 = append(without0, a)
		}
	}

	if len(without0) == 0 {
		return []AcceptEncoding{{Coding: CONTENT_CODING_IDENTITY, Weight: 1}}
	}

	return without0
}

func containsCoding(aes []AcceptEncoding, c ContentCoding) bool {
	for _, ae := range aes {
		if ae.Coding == c {
			return true
		}
	}

	return false
}
