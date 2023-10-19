// Package detect detects a customer's country options based on their Accept-Language and IP address.
package detect

import (
	"net/http"
	"slices"

	"github.com/dys2p/eco/countries"
	"golang.org/x/exp/maps"
)

// Countries returns the union of possible countries for a given HTTP request, based on the client's Accept-Language header and IP address.
//
// The result can contain ISO 3166-1 country codes of European Union countries and the NonEU constant.
// A nil return value means that the client can be anywhere.
func Countries(r *http.Request) ([]countries.Country, error) {
	accept, err := acceptLanguage(r)
	if err != nil {
		return nil, err
	}
	if accept == nil {
		return nil, nil // anywhere
	}

	ip, err := ipAddress(r)
	if err != nil {
		return nil, err
	}
	if ip == nil {
		return nil, nil // anywhere
	}

	var available = make(map[countries.Country]any)
	for _, country := range append(accept, ip...) {
		// replace non-EU country IDs by NonEU
		if !slices.Contains(countries.EuropeanUnion, country) {
			country = countries.NonEU
		}
		available[country] = struct{}{}
	}
	return maps.Keys(available), nil
}

// Valid returns whether the given country is a European Union country or the NonEU constant.
func Valid(c countries.Country) bool {
	return c == countries.NonEU || slices.Contains(countries.EuropeanUnion, c)
}
