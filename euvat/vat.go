// Package euvat models VAT rates from an an European Union point of view.
package euvat

import (
	"math"

	"github.com/dys2p/eco/countries"
)

type Rate string

const (
	RateZero         Rate = "zero" // don't use zero value (empty string) because tax stuff should be explicit
	RateStandard     Rate = "standard"
	RateReduced1     Rate = "reduced-1"
	RateReduced2     Rate = "reduced-2"
	RateSuperReduced Rate = "super-reduced"
	RateParking      Rate = "parking"
)

// Rates are typically the VAT rates of a country.
type Rates map[Rate]float64

// Gross returns the gross of the given net amount using the given VAT rate. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (rates Rates) Gross(net float64, rate Rate) (float64, bool) {
	rateVal, ok := rates.Get(rate)
	return net * (1.0 + rateVal), ok
}

// GrossInt returns the gross of the given net amount using the given VAT rate. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (rates Rates) GrossInt(net int, rate Rate) (int, bool) {
	rateVal, ok := rates.Get(rate)
	return int(math.Round(float64(net) * (1.0 + rateVal))), ok
}

// Gross returns the net of the given gross amount using the given VAT rate. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (rates Rates) Net(gross float64, rate Rate) (float64, bool) {
	rateVal, ok := rates.Get(rate)
	return gross / (1.0 + rateVal), ok
}

// Get returns the value of the given VAT rate. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (rates Rates) Get(rate Rate) (float64, bool) {
	if rate == RateZero { // Rates does not contain RateZero because the zero rate is always the same
		return 0, true
	}
	rateVal, ok := rates[rate]
	if !ok {
		for _, rv := range rates {
			rateVal = max(rateVal, rv) // return max rate
		}
	}
	return rateVal, ok
}

func Convert(value int, src countries.Country, srcRate Rate, dst countries.Country, dstRate Rate) int {
	if src == dst && srcRate == dstRate {
		return value
	}
	srcVal := float64(value)
	netVal, _ := Get(src).Net(srcVal, srcRate)
	dstVal, _ := Get(dst).Gross(netVal, dstRate)
	return int(math.Round(dstVal))
}

// VAT rates are from: https://europa.eu/youreurope/business/taxation/vat/vat-rules-rates/index_en.htm#shortcut-5.
// Note that the ISO 3166-1 code for Greece is "GR", but the VAT rate table uses its ISO 639-1 code "EL".
func Get(c countries.Country) Rates {
	switch c {
	case countries.AT:
		return map[Rate]float64{
			RateStandard: 0.20,
			RateReduced1: 0.10,
			RateReduced2: 0.13,
			RateParking:  0.13,
		}
	case countries.BE:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.06,
			RateReduced2: 0.12,
			RateParking:  0.12,
		}
	case countries.BG:
		return map[Rate]float64{
			RateStandard: 0.20,
			RateReduced1: 0.09,
		}
	case countries.CY:
		return map[Rate]float64{
			RateStandard: 0.19,
			RateReduced1: 0.05,
			RateReduced2: 0.09,
		}
	case countries.CZ:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.12,
		}
	case countries.DE:
		return map[Rate]float64{
			RateStandard: 0.19,
			RateReduced1: 0.07,
		}
	case countries.DK:
		return map[Rate]float64{
			RateStandard: 0.25,
		}
	case countries.EE:
		return map[Rate]float64{
			RateStandard: 0.24,
			RateReduced1: 0.09,
		}
	case countries.ES:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.10,
		}
	case countries.FI:
		return map[Rate]float64{
			RateStandard: 0.255,
			RateReduced1: 0.10,
			RateReduced2: 0.14,
		}
	case countries.FR:
		return map[Rate]float64{
			RateStandard:     0.20,
			RateReduced1:     0.055,
			RateReduced2:     0.10,
			RateSuperReduced: 0.021,
		}
	case countries.GR:
		return map[Rate]float64{
			RateStandard: 0.24,
			RateReduced1: 0.06,
			RateReduced2: 0.13,
		}
	case countries.HR:
		return map[Rate]float64{
			RateStandard: 0.25,
			RateReduced1: 0.05,
			RateReduced2: 0.13,
		}
	case countries.HU:
		return map[Rate]float64{
			RateStandard: 0.27,
			RateReduced1: 0.05,
			RateReduced2: 0.18,
		}
	case countries.IE:
		return map[Rate]float64{
			RateStandard:     0.23,
			RateReduced1:     0.09,
			RateReduced2:     0.135,
			RateSuperReduced: 0.048,
		}
	case countries.IT:
		return map[Rate]float64{
			RateStandard:     0.22,
			RateReduced1:     0.05,
			RateReduced2:     0.10,
			RateSuperReduced: 0.04,
		}
	case countries.LT:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.05,
			RateReduced2: 0.09,
		}
	case countries.LU:
		return map[Rate]float64{
			RateStandard:     0.17,
			RateReduced1:     0.08,
			RateSuperReduced: 0.03,
			RateParking:      0.14,
		}
	case countries.LV:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.05,
			RateReduced2: 0.12,
		}
	case countries.MT:
		return map[Rate]float64{
			RateStandard: 0.18,
			RateReduced1: 0.05,
			RateReduced2: 0.07,
		}
	case countries.NL:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.09,
		}
	case countries.PL:
		return map[Rate]float64{
			RateStandard: 0.23,
			RateReduced1: 0.05,
			RateReduced2: 0.08,
		}
	case countries.PT:
		return map[Rate]float64{
			RateStandard: 0.23,
			RateReduced1: 0.06,
			RateReduced2: 0.13,
			RateParking:  0.13,
		}
	case countries.RO:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.11,
		}
	case countries.SE:
		return map[Rate]float64{
			RateStandard: 0.25,
			RateReduced1: 0.06,
			RateReduced2: 0.12,
		}
	case countries.SI:
		return map[Rate]float64{
			RateStandard: 0.22,
			RateReduced1: 0.05,
			RateReduced2: 0.095,
		}
	case countries.SK:
		return map[Rate]float64{
			RateStandard:     0.23,
			RateReduced1:     0.19,
			RateSuperReduced: 0.05,
		}
	default:
		return map[Rate]float64{
			RateStandard: 0.0,
		}
	}
}
