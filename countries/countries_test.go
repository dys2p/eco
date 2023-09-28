package countries

import (
	"math"
	"testing"
)

const epsilon = 1e-9

func TestGrossNet(t *testing.T) {
	tests := []struct {
		country string
		rate    string
		net     float64
		gross   float64
	}{
		{"DE", "reduced-1", 100, 107},
		{"DE", "standard", 100, 119},
		{"DE", "unknown", 100, 119},
		{"IE", "reduced-1", 100, 109},
		{"IE", "standard", 100, 123},
		{"IE", "unknown", 100, 123},
	}

	for _, test := range tests {
		country, ok := Get(test.country)
		if !ok {
			t.Fatalf("country not found: %s", test.country)
		}
		if gross, _ := country.Gross(test.net, test.rate); math.Abs(gross-test.gross) > epsilon {
			t.Fatalf("gross: got %f, want %f", gross, test.gross)
		}
		if net, _ := country.Net(test.gross, test.rate); math.Abs(net-test.net) > epsilon {
			t.Fatalf("net: got %f, want %f", net, test.net)
		}
	}
}

func TestVATRate(t *testing.T) {
	tests := []struct {
		country string
		rate    string
		want    float64
	}{
		{"DE", "reduced-1", 0.07},
		{"DE", "standard", 0.19},
		{"DE", "unknown", 0.19},
		{"IE", "reduced-1", 0.09},
		{"IE", "standard", 0.23},
		{"IE", "unknown", 0.23},
	}

	for _, test := range tests {
		country, ok := Get(test.country)
		if !ok {
			t.Fatalf("country not found: %s", test.country)
		}
		if got, _ := country.VATRate(test.rate); got != test.want {
			t.Fatalf("VATRate: got %f, want %f", got, test.want)
		}
	}
}
