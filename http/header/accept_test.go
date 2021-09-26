package header

import "testing"

func TestParseAccept1(t *testing.T) {
	actual := ParseAccept("audio/*; q=0.2, audio/basic")

	assertType(t, actual[0], "audio", "basic") // high priority
	assertNoParameter(t, actual[0].AcceptParameters)

	assertType(t, actual[1], "audio", "*")
	assertParameterSize(t, actual[1].AcceptParameters, 1)
	assertParameter(t, actual[1].AcceptParameters[0], "q", "0.2")
}

func TestParseAccept2(t *testing.T) {
	actual := ParseAccept("text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c")

	assertType(t, actual[0], "text", "html")
	assertNoParameter(t, actual[0].AcceptParameters)

	assertType(t, actual[1], "text", "x-c")
	assertNoParameter(t, actual[1].AcceptParameters)

	assertType(t, actual[2], "text", "x-dvi")
	assertParameterSize(t, actual[2].AcceptParameters, 1)
	assertParameter(t, actual[2].AcceptParameters[0], "q", "0.8")

	assertType(t, actual[3], "text", "plain")
	assertParameterSize(t, actual[3].AcceptParameters, 1)
	assertParameter(t, actual[3].AcceptParameters[0], "q", "0.5")
}

func assertType(t *testing.T, accept Accept, typ string, subType string) {
	if accept.Type != typ {
		t.Errorf("Unexpected type: %v.", accept.Type)
	}
	if accept.SubType != subType {
		t.Errorf("Unexpected type: %v.", accept.SubType)
	}
}

func assertParameter(t *testing.T, param AcceptParameter, key string, value string) {
	if param.Key != key {
		t.Errorf("Unexpected parameter: %v.", param.Key)
	}
	if param.Value != value {
		t.Errorf("Unexpected parameter: %v.", param.Value)
	}
}

func assertParameterSize(t *testing.T, params []AcceptParameter, expected int) {
	if len(params) != expected {
		t.Errorf("Unexpected parameter: %v.", params)
	}
}

func assertNoParameter(t *testing.T, params []AcceptParameter) {
	if len(params) != 0 {
		t.Errorf("Unexpected parameter: %v.", params)
	}
}
