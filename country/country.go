// Package country contains European countries and their VAT rates from an European Union point of view.
package country

import "sort"

type Country struct {
	ID       string // ISO 3166-1 code
	Name     string // German
	VATRates map[string]float64
}

// maxVATRate returns the highest VAT rate. It can be useful to get a fail-safe rate if a key is not present in the VATRates map.
func (c Country) maxVATRate() float64 {
	var max float64 = 0
	for _, rate := range c.VATRates {
		if max < rate {
			max = rate
		}
	}
	return max
}

// Gross returns the gross of the given net amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (c Country) Gross(net int, rateKey string) (int, bool) {
	rate, ok := c.VATRates[rateKey]
	if ok {
		return int(float64(net) * (1.0 + rate)), ok
	} else {
		return int(float64(net) * (1.0 + c.maxVATRate())), ok
	}
}

// Gross returns the net of the given gross amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (c Country) Net(gross int, rateKey string) (int, bool) {
	rate, ok := c.VATRates[rateKey]
	if ok {
		return int(float64(gross) / (1.0 + rate)), ok
	} else {
		return int(float64(gross) / (1.0 + c.maxVATRate())), ok
	}
}

func Get(id string) (Country, bool) {
	for _, c := range All {
		if c.ID == id {
			return c, true
		}
	}
	return Country{}, false
}

func SortByName(countries []Country) []Country {
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})
	return countries
}
