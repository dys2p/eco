package euvat

import (
	"math"
	"testing"

	"github.com/dys2p/eco/countries"
)

const epsilon = 1e-9

func TestGrossNet(t *testing.T) {
	tests := []struct {
		country countries.Country
		rate    Rate
		net     float64
		gross   int
		ok      bool
	}{
		{countries.DE, RateReduced1, 100, 107, true},
		{countries.DE, RateStandard, 100, 119, true},
		{countries.DE, RateSuperReduced, 100, 100, false}, // rate does not exist in DE
		{countries.DE, Rate("unknown"), 100, 100, false},
		{countries.IE, RateReduced1, 100, 109, true},
		{countries.IE, RateStandard, 100, 123, true},
		{countries.IE, Rate("unknown"), 100, 100, false},
	}

	for _, test := range tests {
		if gross, ok := Gross(test.country, test.net, test.rate); math.Abs(float64(gross-test.gross)) > epsilon || ok != test.ok {
			t.Fatalf("gross: got (%d, %v), want (%d, %v)", gross, ok, test.gross, test.ok)
		}
		if net, ok := Net(test.country, test.gross, test.rate); math.Abs(net-test.net) > epsilon || ok != test.ok {
			t.Fatalf("net: got (%f, %v), want (%f, %v)", net, ok, test.net, test.ok)
		}
		if netInt, ok := NetInt(test.country, test.gross, test.rate); math.Abs(float64(netInt)-test.net) > epsilon || ok != test.ok {
			t.Fatalf("net: got (%d, %v), want (%f, %v)", netInt, ok, test.net, test.ok)
		}
	}
}

func TestValue(t *testing.T) {
	tests := []struct {
		country   countries.Country
		rate      Rate
		wantValue float64
		wantOk    bool
	}{
		{countries.DE, RateReduced1, 0.07, true},
		{countries.DE, RateStandard, 0.19, true},
		{countries.DE, RateSuperReduced, 0.0, false},
		{countries.DE, Rate("unknown"), 0.0, false},
		{countries.IE, RateReduced1, 0.09, true},
		{countries.IE, RateStandard, 0.23, true},
		{countries.IE, Rate("unknown"), 0.0, false},
	}

	for _, test := range tests {
		if value, ok := Value(test.country, test.rate); value != test.wantValue || ok != test.wantOk {
			t.Fatalf("VATRate: got (%f, %v), want (%f, %v)", value, ok, test.wantValue, test.wantOk)
		}
	}
}
