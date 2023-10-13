// Package countries contains European countries and their VAT rates from an European Union point of view.
package countries

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Country struct {
	ID       string // ISO 3166-1 code
	VATRates map[string]float64
}

// Gross returns the gross of the given net amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (c Country) Gross(net float64, rateKey string) (float64, bool) {
	rate, ok := c.VATRate(rateKey)
	return net * (1.0 + rate), ok
}

// Gross returns the net of the given gross amount using the given VAT rate key. The boolean return value indicates if the rate has been found. If it is not found, the maximum rate is used.
func (c Country) Net(gross float64, rateKey string) (float64, bool) {
	rate, ok := c.VATRate(rateKey)
	return gross / (1.0 + rate), ok
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

func (c Country) TranslateName(langstr string) string {
	tag, _ := language.MatchStrings(matcher, langstr)
	printer := message.NewPrinter(tag, message.Catalog(cat))
	switch c.ID {
	case AT.ID:
		return printer.Sprintf("Austria")
	case BE.ID:
		return printer.Sprintf("Belgium")
	case BG.ID:
		return printer.Sprintf("Bulgaria")
	case CH.ID:
		return printer.Sprintf("Switzerland")
	case CY.ID:
		return printer.Sprintf("Cyprus")
	case CZ.ID:
		return printer.Sprintf("Czechia")
	case DE.ID:
		return printer.Sprintf("Germany")
	case DK.ID:
		return printer.Sprintf("Denmark")
	case EE.ID:
		return printer.Sprintf("Estonia")
	case ES.ID:
		return printer.Sprintf("Spain")
	case FI.ID:
		return printer.Sprintf("Finland")
	case FR.ID:
		return printer.Sprintf("France")
	case GB.ID:
		return printer.Sprintf("United Kingdom")
	case GR.ID:
		return printer.Sprintf("Greece")
	case HR.ID:
		return printer.Sprintf("Croatia")
	case HU.ID:
		return printer.Sprintf("Hungary")
	case IE.ID:
		return printer.Sprintf("Ireland")
	case IT.ID:
		return printer.Sprintf("Italy")
	case LT.ID:
		return printer.Sprintf("Lithuania")
	case LU.ID:
		return printer.Sprintf("Luxembourg")
	case LV.ID:
		return printer.Sprintf("Latvia")
	case MT.ID:
		return printer.Sprintf("Malta")
	case NL.ID:
		return printer.Sprintf("Netherlands")
	case PL.ID:
		return printer.Sprintf("Poland")
	case PT.ID:
		return printer.Sprintf("Portugal")
	case RO.ID:
		return printer.Sprintf("Romania")
	case SE.ID:
		return printer.Sprintf("Sweden")
	case SI.ID:
		return printer.Sprintf("Slovenia")
	case SK.ID:
		return printer.Sprintf("Slovakia")
	default:
		return ""
	}
}
