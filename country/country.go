package country

import "sort"

type Country struct {
	ID         string  // ISO-3166-1 code
	Name       string  // German
	VATRegular float64 // e.g. 0.19
}

// Gross returns the gross of the given net amount using the country's regular VAT rate.
func (c Country) Gross(net int) int {
	return int(float64(net) * c.vatRegularFactor())
}

// Gross returns the net of the given gross amount using the country's regular VAT rate.
func (c Country) Net(gross int) int {
	return int(float64(gross) / c.vatRegularFactor())
}

func (c Country) vatRegularFactor() float64 {
	return 1.0 + c.VATRegular
}

func Get(id string) (Country, bool) {
	for _, c := range All {
		if c.ID == id {
			return c, true
		}
	}
	return Country{}, false
}

func Sort(countries []Country) []Country {
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})
	return countries
}
