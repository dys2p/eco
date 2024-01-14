package ntfysh

import "testing"

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"foo-bar", "https://ntfy.sh/foo-bar"},                 // dash is ok
		{"FOO_BAR", "https://ntfy.sh/FOO_BAR"},                 // underscore is ok
		{"1234567", "https://ntfy.sh/1234567"},                 // digits are ok
		{"example.com/foo", "https://example.com/foo"},         // add https
		{"http://example.com/foo", "http://example.com/foo"},   // keep http
		{"https://example.com/foo", "https://example.com/foo"}, // keep https
		{"ftp://example.com/foo", ""},                          // wrong scheme
		{"foo?bar", ""},                                        // invalid character
		{"example.com/foo?bar", ""},                            // invalid character
		{"blob:foo", ""},                                       // url with opaque data
	}

	for _, test := range tests {
		if got := ValidateAddress(test.input); got != test.want {
			t.Fatalf("got %s, want %s", got, test.want)
		}
	}
}
