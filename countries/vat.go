package countries

type VATRates map[string]float64

// Gross returns the gross of the given net amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (vr VATRates) Gross(net float64, rateKey string) (float64, bool) {
	rate, ok := vr.Rate(rateKey)
	return net * (1.0 + rate), ok
}

// Gross returns the net of the given gross amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (vr VATRates) Net(gross float64, rateKey string) (float64, bool) {
	rate, ok := vr.Rate(rateKey)
	return gross / (1.0 + rate), ok
}

// VATRate returns the VAT rate with the given key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (vr VATRates) Rate(rateKey string) (float64, bool) {
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