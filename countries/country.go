// Package countries contains European countries and their VAT rates from an European Union point of view.
package countries

type Country struct {
	ID       string // ISO 3166-1 code
	Name     string // German
	VATRates map[string]float64
}

// Gross returns the gross of the given net amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (c Country) Gross(net float64, rateKey string) (float64, bool) {
	rate, ok := c.VATRate(rateKey)
	return float64(net) * (1.0 + rate), ok
}

// Gross returns the net of the given gross amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (c Country) Net(gross float64, rateKey string) (float64, bool) {
	rate, ok := c.VATRate(rateKey)
	return float64(gross) / (1.0 + rate), ok
}

// VATRate returns the VAT rate with the given key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (c Country) VATRate(key string) (float64, bool) {
	rate, ok := c.VATRates[key]
	if !ok {
		for _, r := range c.VATRates {
			if rate < r {
				rate = r
			}
		}
	}
	return rate, ok
}
