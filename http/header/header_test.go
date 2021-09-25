package header

import (
	"testing"
)

func TestParseHeader(t *testing.T) {
	type expected struct {
		fieldName  string
		fieldValue string
	}

	tests := []struct {
		name string
		args string
		want expected
	}{
		{
			name: "Location",
			args: "Location: example.com",
			want: expected{fieldName: "LOCATION", fieldValue: "example.com"},
		},
		{
			name: "Content-Type",
			args: "Content-Type: text/html; charset=utf-8",
			want: expected{fieldName: "CONTENT-TYPE", fieldValue: "text/html; charset=utf-8"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := ParseHeader(tt.args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err.Error())
			}
			if h.FieldName != tt.want.fieldName {
				t.Errorf("Unexpected field name: %v", h.FieldName)
			}
			if h.FieldValue != tt.want.fieldValue {
				t.Errorf("Unexpected field value: %v", h.FieldValue)
			}

		})
	}
}

func TestValidate_NoHostHeader(t *testing.T) {
	h, err := ParseHeader("AAA: aaa")
	if err != nil {
		t.Errorf("Unexpected error: %v.", err)
	}
	hs := Headers{h}
	httpError := hs.Validate()
	if httpError == nil {
		t.Error("Unexpected success")
	}
	if httpError.Error() != "Request require only one Host header." {
		t.Errorf("Unexpected error: %v.", httpError.Error())
	}
}

func TestValidate_MultipleHostHeader(t *testing.T) {
	h1, err := ParseHeader("Host: example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v.", err)
	}
	h2, err := ParseHeader("Host: example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v.", err)
	}
	hs := Headers{h1, h2}
	httpError := hs.Validate()
	if httpError == nil {
		t.Error("Unexpected success")
	}
	if httpError.Error() != "Request require only one Host header." {
		t.Errorf("Unexpected error: %v.", httpError.Error())
	}
}
