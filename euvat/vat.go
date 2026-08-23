// Package euvat models VAT rates from an an European Union point of view.
//
// Gross values are represented as int because they are usually explicit.
// Net values are represented as float64 because they are usually intermediate.
//
// Note that the ISO 3166-1 code for Greece is "GR", but the VAT rate table uses its ISO 639-1 code "EL".
package euvat

import (
	"math"

	"github.com/dys2p/eco/countries"
)

type Rate string

const (
	RateAuxiliary    Rate = "" // taxation depends on the main product
	RateZero         Rate = "zero"
	RateStandard     Rate = "standard"
	RateReduced1     Rate = "reduced-1"
	RateReduced2     Rate = "reduced-2"
	RateSuperReduced Rate = "super-reduced"
	RateParking      Rate = "parking"
)

// Gross returns the gross of the given net amount using the given VAT rate. The boolean return value indicates if the country and rate have been found.
func Gross(c countries.Country, net float64, rate Rate) (int, bool) {
	val, ok := Value(c, rate)
	return int(math.Round(net * (1.0 + val))), ok
}

// Net returns the net of the given gross amount using the given VAT rate. The boolean return value indicates if the country and rate have been found.
func Net(c countries.Country, gross int, rate Rate) (float64, bool) {
	val, ok := Value(c, rate)
	return float64(gross) / (1.0 + val), ok
}

func NetInt(c countries.Country, gross int, rate Rate) (int, bool) {
	val, ok := Value(c, rate)
	return int(math.Round(float64(gross) / (1.0 + val))), ok
}

// VAT rates are from: https://europa.eu/youreurope/business/taxation/vat/vat-rules-rates/index_en.htm#shortcut-5.
func Value(c countries.Country, r Rate) (v float64, ok bool) {
	if r == RateZero {
		return 0, true
	}

	switch c {
	case countries.AT:
		v, ok = map[Rate]float64{
			RateStandard: 0.20,
			RateReduced1: 0.10,
			RateReduced2: 0.13,
			RateParking:  0.13,
		}[r]
	case countries.BE:
		v, ok = map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.06,
			RateReduced2: 0.12,
			RateParking:  0.12,
		}[r]
	case countries.BG:
		v, ok = map[Rate]float64{
			RateStandard: 0.20,
			RateReduced1: 0.09,
		}[r]
	case countries.CY:
		v, ok = map[Rate]float64{
			RateStandard: 0.19,
			RateReduced1: 0.05,
			RateReduced2: 0.09,
		}[r]
	case countries.CZ:
		v, ok = map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.12,
		}[r]
	case countries.DE:
		v, ok = map[Rate]float64{
			RateStandard: 0.19,
			RateReduced1: 0.07,
		}[r]
	case countries.DK:
		v, ok = map[Rate]float64{
			RateStandard: 0.25,
		}[r]
	case countries.EE:
		v, ok = map[Rate]float64{
			RateStandard: 0.24,
			RateReduced1: 0.09,
		}[r]
	case countries.ES:
		v, ok = map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.10,
		}[r]
	case countries.FI:
		v, ok = map[Rate]float64{
			RateStandard: 0.255,
			RateReduced1: 0.10,
			RateReduced2: 0.14,
		}[r]
	case countries.FR:
		v, ok = map[Rate]float64{
			RateStandard:     0.20,
			RateReduced1:     0.055,
			RateReduced2:     0.10,
			RateSuperReduced: 0.021,
		}[r]
	case countries.GR:
		v, ok = map[Rate]float64{
			RateStandard: 0.24,
			RateReduced1: 0.06,
			RateReduced2: 0.13,
		}[r]
	case countries.HR:
		v, ok = map[Rate]float64{
			RateStandard: 0.25,
			RateReduced1: 0.05,
			RateReduced2: 0.13,
		}[r]
	case countries.HU:
		v, ok = map[Rate]float64{
			RateStandard: 0.27,
			RateReduced1: 0.05,
			RateReduced2: 0.18,
		}[r]
	case countries.IE:
		v, ok = map[Rate]float64{
			RateStandard:     0.23,
			RateReduced1:     0.09,
			RateReduced2:     0.135,
			RateSuperReduced: 0.048,
		}[r]
	case countries.IT:
		v, ok = map[Rate]float64{
			RateStandard:     0.22,
			RateReduced1:     0.05,
			RateReduced2:     0.10,
			RateSuperReduced: 0.04,
		}[r]
	case countries.LT:
		v, ok = map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.05,
			RateReduced2: 0.09,
		}[r]
	case countries.LU:
		v, ok = map[Rate]float64{
			RateStandard:     0.17,
			RateReduced1:     0.08,
			RateSuperReduced: 0.03,
			RateParking:      0.14,
		}[r]
	case countries.LV:
		v, ok = map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.05,
			RateReduced2: 0.12,
		}[r]
	case countries.MT:
		v, ok = map[Rate]float64{
			RateStandard: 0.18,
			RateReduced1: 0.05,
			RateReduced2: 0.07,
		}[r]
	case countries.NL:
		v, ok = map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.09,
		}[r]
	case countries.PL:
		v, ok = map[Rate]float64{
			RateStandard: 0.23,
			RateReduced1: 0.05,
			RateReduced2: 0.08,
		}[r]
	case countries.PT:
		v, ok = map[Rate]float64{
			RateStandard: 0.23,
			RateReduced1: 0.06,
			RateReduced2: 0.13,
			RateParking:  0.13,
		}[r]
	case countries.RO:
		v, ok = map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.11,
		}[r]
	case countries.SE:
		v, ok = map[Rate]float64{
			RateStandard: 0.25,
			RateReduced1: 0.06,
			RateReduced2: 0.12,
		}[r]
	case countries.SI:
		v, ok = map[Rate]float64{
			RateStandard: 0.22,
			RateReduced1: 0.05,
			RateReduced2: 0.095,
		}[r]
	case countries.SK:
		v, ok = map[Rate]float64{
			RateStandard:     0.23,
			RateReduced1:     0.19,
			RateSuperReduced: 0.05,
		}[r]
	}

	return v, ok
}
