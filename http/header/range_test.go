package header

import "testing"

func TestParseRange_Single(t *testing.T) {
	ranges, err := ParseRange("bytes=0-100")
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) != 1 {
		t.Errorf("Unexpected Ranges: %v", ranges)
	}

	assertRange(t, ranges[0], intPointer(0), intPointer(100))
}

func TestParseRange_Multiple(t *testing.T) {
	ranges, err := ParseRange("bytes=0-100, 101-200")
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) != 2 {
		t.Errorf("Unexpected Ranges: %v", ranges)
	}

	assertRange(t, ranges[0], intPointer(0), intPointer(100))
	assertRange(t, ranges[1], intPointer(101), intPointer(200))
}

func TestParseRange_Complex1(t *testing.T) {
	ranges, err := ParseRange("bytes=200-1000, 2000-6576, 19000-")
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) != 3 {
		t.Errorf("Unexpected Ranges: %v", ranges)
	}

	assertRange(t, ranges[0], intPointer(200), intPointer(1000))
	assertRange(t, ranges[1], intPointer(2000), intPointer(6576))
	assertRange(t, ranges[2], intPointer(19000), nil)
}

func TestParseRange_Complex2(t *testing.T) {
	ranges, err := ParseRange("bytes=0-499, -500")
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) != 2 {
		t.Errorf("Unexpected Ranges: %v", ranges)
	}

	assertRange(t, ranges[0], intPointer(0), intPointer(499))
	assertRange(t, ranges[1], nil, intPointer(500))
}

func assertRange(t *testing.T, r Range, start *int, end *int) {
	if start != nil && *r.Start != *start {
		t.Errorf("Unexpected Ranges: %v", r)
	}
	if end != nil && *r.End != *end {
		t.Errorf("Unexpected Ranges: %v", r)
	}
	if start == nil && r.Start != start {
		t.Errorf("Unexpected Ranges: %v", r)
	}
	if end == nil && r.End != end {
		t.Errorf("Unexpected Ranges: %v", r)
	}
}

func intPointer(val int) *int {
	return &val
}
