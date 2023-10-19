// Package detect detects a customer's country options based on their Accept-Language and IP address.
package detect

import (
	"net/http"

	"github.com/dys2p/eco/countries"
	"golang.org/x/exp/maps"
)

type detector func(r *http.Request) ([]countries.Country, error)

// Countries returns all possible countries for a given HTTP request, based on the client's Accept-Language header and IP address.
//
// The result is a slice of European Union countries and a boolean value which indicates "non-EU".
func Countries(r *http.Request) ([]countries.Country, bool, error) {
	var eu = make(map[countries.Country]any)
	var nonEU = false
	for _, f := range []detector{acceptLanguage, ipAddress} {
		detectedCountries, err := f(r)
		if err != nil {
			return nil, false, err
		}
		// check if user can be anywhere
		if detectedCountries == nil {
			return countries.EuropeanUnion, true, nil
		}
		// merge EU countries
		for _, country := range detectedCountries {
			// important: check whether country is in the European Union
			if countries.InEuropeanUnion(country) {
				eu[country] = struct{}{}
			} else {
				nonEU = true
			}
		}
	}
	return maps.Keys(eu), nonEU, nil
}
