package euvat

import (
	"math"
	"testing"

	"github.com/dys2p/eco/countries"
)

const epsilon = 1e-9

func TestConvert(t *testing.T) {
	tests := []struct {
		value   int
		src     countries.Country
		srcRate Rate
		dst     countries.Country
		dstRate Rate
		want    int
	}{
		{100, countries.DE, RateStandard, countries.DE, RateStandard, 100},
		{100, countries.DE, RateReduced1, countries.DK, RateStandard, 117},
	}

	for _, test := range tests {
		if got := Convert(test.value, test.src, test.srcRate, test.dst, test.dstRate); got != test.want {
			t.Fatalf("convert: got %d, want %d", got, test.want)
		}
	}
}

func TestGrossNet(t *testing.T) {
	tests := []struct {
		country countries.Country
		rate    Rate
		net     float64
		gross   float64
	}{
		{countries.DE, RateReduced1, 100, 107},
		{countries.DE, RateStandard, 100, 119},
		{countries.DE, Rate("unknown"), 100, 119},
		{countries.IE, RateReduced1, 100, 109},
		{countries.IE, RateStandard, 100, 123},
		{countries.IE, Rate("unknown"), 100, 123},
	}

	for _, test := range tests {
		if gross, _ := Get(test.country).Gross(test.net, test.rate); math.Abs(gross-test.gross) > epsilon {
			t.Fatalf("gross: got %f, want %f", gross, test.gross)
		}
		if net, _ := Get(test.country).Net(test.gross, test.rate); math.Abs(net-test.net) > epsilon {
			t.Fatalf("net: got %f, want %f", net, test.net)
		}
	}
}

func TestVATRate(t *testing.T) {
	tests := []struct {
		country countries.Country
		rate    Rate
		want    float64
	}{
		{countries.DE, RateReduced1, 0.07},
		{countries.DE, RateStandard, 0.19},
		{countries.DE, Rate("unknown"), 0.19},
		{countries.IE, RateReduced1, 0.09},
		{countries.IE, RateStandard, 0.23},
		{countries.IE, Rate("unknown"), 0.23},
	}

	for _, test := range tests {
		if got, _ := Get(test.country).Get(test.rate); got != test.want {
			t.Fatalf("VATRate: got %f, want %f", got, test.want)
		}
	}
}
