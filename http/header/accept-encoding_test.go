package header

import "testing"

func assertAcceptEncoding(t *testing.T, actual AcceptEncoding, expectedContentCoding ContentCoding, expectedWeight float64) {
	if actual.Coding != expectedContentCoding {
		t.Errorf("Unexpected coding: %v.", actual.Coding)
	}
	if actual.Weight != expectedWeight {
		t.Errorf("Unexpected weight: %v.", actual.Weight)
	}
}
func TestParseAcceptEncoding_Multiple(t *testing.T) {
	actual := ParseAcceptEncoding("compress, gzip")

	if len(actual) != 2 {
		t.Errorf("Unexpected result: %v.", actual)
	}
	assertAcceptEncoding(t, actual[0], CONTENT_CODING_COMPRESS, 1)
	assertAcceptEncoding(t, actual[1], CONTENT_CODING_GZIP, 1)
}

func TestParseAcceptEncoding_Priority(t *testing.T) {
	actual := ParseAcceptEncoding("compress;q=0.5, gzip;q=1.0")

	if len(actual) != 2 {
		t.Errorf("Unexpected result: %v.", actual)
	}
	assertAcceptEncoding(t, actual[0], CONTENT_CODING_GZIP, 1)
	assertAcceptEncoding(t, actual[1], CONTENT_CODING_COMPRESS, 0.5)
}

func TestParseAcceptEncoding_WildCard(t *testing.T) {
	actual := ParseAcceptEncoding("*")

	if len(actual) != 4 {
		t.Errorf("Unexpected result: %v.", actual)
	}
	assertAcceptEncoding(t, actual[0], CONTENT_CODING_COMPRESS, 1)
	assertAcceptEncoding(t, actual[1], CONTENT_CODING_DEFLATE, 1)
	assertAcceptEncoding(t, actual[2], CONTENT_CODING_GZIP, 1)
	assertAcceptEncoding(t, actual[3], CONTENT_CODING_IDENTITY, 1)
}

func TestParseAcceptEncoding_Complex1(t *testing.T) {
	actual := ParseAcceptEncoding("gzip;q=1.0, identity; q=0.5, *;q=0")

	if len(actual) != 2 {
		t.Errorf("Unexpected result: %v.", actual)
	}
	assertAcceptEncoding(t, actual[0], CONTENT_CODING_GZIP, 1)
	assertAcceptEncoding(t, actual[1], CONTENT_CODING_IDENTITY, 0.5)
}

func TestParseAcceptEncoding_Complex2(t *testing.T) {
	actual := ParseAcceptEncoding("gzip;q=1.0, identity; q=0.5, compress; q=0, *;q=0.3")

	if len(actual) != 3 {
		t.Errorf("Unexpected result: %v.", actual)
	}
	assertAcceptEncoding(t, actual[0], CONTENT_CODING_GZIP, 1)
	assertAcceptEncoding(t, actual[1], CONTENT_CODING_IDENTITY, 0.5)
	assertAcceptEncoding(t, actual[2], CONTENT_CODING_DEFLATE, 0.3)
}

func TestParseAcceptEncoding_AllDeny(t *testing.T) {
	actual := ParseAcceptEncoding("gzip;q=0")

	if len(actual) != 1 {
		t.Errorf("Unexpected result: %v.", actual)
	}
	assertAcceptEncoding(t, actual[0], CONTENT_CODING_IDENTITY, 1)
}
