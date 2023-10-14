// Package countries contains European countries and their VAT rates from an European Union point of view.
package countries

import "github.com/dys2p/eco/lang"

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
	l := lang.Lang(langstr)
	switch c.ID {
	case AT.ID:
		return l.Tr("Austria")
	case BE.ID:
		return l.Tr("Belgium")
	case BG.ID:
		return l.Tr("Bulgaria")
	case CH.ID:
		return l.Tr("Switzerland")
	case CY.ID:
		return l.Tr("Cyprus")
	case CZ.ID:
		return l.Tr("Czechia")
	case DE.ID:
		return l.Tr("Germany")
	case DK.ID:
		return l.Tr("Denmark")
	case EE.ID:
		return l.Tr("Estonia")
	case ES.ID:
		return l.Tr("Spain")
	case FI.ID:
		return l.Tr("Finland")
	case FR.ID:
		return l.Tr("France")
	case GB.ID:
		return l.Tr("United Kingdom")
	case GR.ID:
		return l.Tr("Greece")
	case HR.ID:
		return l.Tr("Croatia")
	case HU.ID:
		return l.Tr("Hungary")
	case IE.ID:
		return l.Tr("Ireland")
	case IT.ID:
		return l.Tr("Italy")
	case LT.ID:
		return l.Tr("Lithuania")
	case LU.ID:
		return l.Tr("Luxembourg")
	case LV.ID:
		return l.Tr("Latvia")
	case MT.ID:
		return l.Tr("Malta")
	case NL.ID:
		return l.Tr("Netherlands")
	case PL.ID:
		return l.Tr("Poland")
	case PT.ID:
		return l.Tr("Portugal")
	case RO.ID:
		return l.Tr("Romania")
	case SE.ID:
		return l.Tr("Sweden")
	case SI.ID:
		return l.Tr("Slovenia")
	case SK.ID:
		return l.Tr("Slovakia")
	default:
		return ""
	}
}
