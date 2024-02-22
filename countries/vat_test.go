package countries

import (
	"math"
	"testing"
)

const epsilon = 1e-9

func TestConvert(t *testing.T) {
	tests := []struct {
		value   int
		src     Country
		srcRate Rate
		dst     Country
		dstRate Rate
		want    int
	}{
		{100, DE, RateStandard, DE, RateStandard, 100},
		{100, DE, RateReduced1, DK, RateStandard, 117},
	}

	for _, test := range tests {
		if got := Convert(test.value, test.src, test.srcRate, test.dst, test.dstRate); got != test.want {
			t.Fatalf("convert: got %d, want %d", got, test.want)
		}
	}
}

func TestGrossNet(t *testing.T) {
	tests := []struct {
		country Country
		rate    Rate
		net     float64
		gross   float64
	}{
		{DE, RateReduced1, 100, 107},
		{DE, RateStandard, 100, 119},
		{DE, Rate("unknown"), 100, 119},
		{IE, RateReduced1, 100, 109},
		{IE, RateStandard, 100, 123},
		{IE, Rate("unknown"), 100, 123},
	}

	for _, test := range tests {
		if gross, _ := test.country.VAT().Gross(test.net, test.rate); math.Abs(gross-test.gross) > epsilon {
			t.Fatalf("gross: got %f, want %f", gross, test.gross)
		}
		if net, _ := test.country.VAT().Net(test.gross, test.rate); math.Abs(net-test.net) > epsilon {
			t.Fatalf("net: got %f, want %f", net, test.net)
		}
	}
}

func TestVATRate(t *testing.T) {
	tests := []struct {
		country Country
		rate    Rate
		want    float64
	}{
		{DE, RateReduced1, 0.07},
		{DE, RateStandard, 0.19},
		{DE, Rate("unknown"), 0.19},
		{IE, RateReduced1, 0.09},
		{IE, RateStandard, 0.23},
		{IE, Rate("unknown"), 0.23},
	}

	for _, test := range tests {
		if got, _ := test.country.VAT().Rate(test.rate); got != test.want {
			t.Fatalf("VATRate: got %f, want %f", got, test.want)
		}
	}
}
