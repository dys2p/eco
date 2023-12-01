package countries

import "math"

type Rate string

const (
	RateStandard     Rate = "standard"
	RateReduced1     Rate = "reduced-1"
	RateReduced2     Rate = "reduced-2"
	RateSuperReduced Rate = "super-reduced"
	RateParking      Rate = "parking"
)

type VATRates map[Rate]float64

func Convert(value int, src Country, srcRate Rate, dst Country, dstRate Rate) int {
	if src == dst && srcRate == dstRate {
		return value
	}
	srcVal := float64(value)
	netVal, _ := src.VAT().Net(srcVal, srcRate)
	dstVal, _ := dst.VAT().Gross(netVal, dstRate)
	return int(math.Round(dstVal))
}

// Gross returns the gross of the given net amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (vr VATRates) Gross(net float64, rateKey Rate) (float64, bool) {
	rate, ok := vr.Rate(rateKey)
	return net * (1.0 + rate), ok
}

// Gross returns the net of the given gross amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (vr VATRates) Net(gross float64, rateKey Rate) (float64, bool) {
	rate, ok := vr.Rate(rateKey)
	return gross / (1.0 + rate), ok
}

// VATRate returns the VAT rate with the given key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (vr VATRates) Rate(rateKey Rate) (float64, bool) {
	if vr == nil {
		return 0, true
	}

	rate, ok := vr[rateKey]
	if !ok {
		// return max rate
		for _, r := range vr {
			if rate < r {
				rate = r
			}
		}
	}
	return rate, ok
}
