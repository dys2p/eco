package eco

import "testing"

func TestFmt(t *testing.T) {
	tests := []struct {
		got  string
		want string
	}{
		{FmtEuro(123), "1,23 €"},
		{string(FmtEuroHTML(123)), "1,23&nbsp;€"},
		{string(FmtEuroMinMaxHTML(123, 123)), "1,23&nbsp;€"},
		{string(FmtEuroMinMaxHTML(123, 456)), "1,23&nbsp;–&nbsp;4,56&nbsp;€"},
		{string(FmtEuroPlusMinusHTML(123)), "+1,23&nbsp;€"},
		{string(FmtEuroPlusMinusHTML(-123)), "−1,23&nbsp;€"},
		{string(FmtPercentHTML(0.12)), "12&nbsp;%"},
		{string(FmtPercentHTML(0.123)), "12,3&nbsp;%"},
	}

	for _, test := range tests {
		if test.got != test.want {
			t.Fatalf("got %v, want %v", test.got, test.want)
		}
	}
}
