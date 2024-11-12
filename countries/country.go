// Package countries contains European countries and their VAT rates from an European Union point of view.
package countries

import (
	"slices"
	"sort"

	"github.com/dys2p/eco/lang"
)

type Country string // ISO 3166-1 code

const (
	// European Union
	AT Country = "AT"
	BE Country = "BE"
	BG Country = "BG"
	CY Country = "CY"
	CZ Country = "CZ"
	DE Country = "DE"
	DK Country = "DK"
	EE Country = "EE"
	ES Country = "ES"
	FI Country = "FI"
	FR Country = "FR"
	GR Country = "GR"
	HR Country = "HR"
	HU Country = "HU"
	IE Country = "IE"
	IT Country = "IT"
	LT Country = "LT"
	LU Country = "LU"
	LV Country = "LV"
	MT Country = "MT"
	NL Country = "NL"
	PL Country = "PL"
	PT Country = "PT"
	RO Country = "RO"
	SE Country = "SE"
	SI Country = "SI"
	SK Country = "SK"

	NonEU Country = "non-EU"

	// selected Non-EU countries
	CH Country = "CH"
	GB Country = "GB"
	ME Country = "ME"
)

var All = []Country{AT, BE, BG, CH, CY, CZ, DE, DK, EE, ES, FI, FR, GB, GR, HR, HU, IE, IT, LT, LU, LV, ME, MT, NL, PL, PT, RO, SE, SI, SK}

var EuropeanUnion = []Country{AT, BE, BG, CY, CZ, DE, DK, EE, ES, FI, FR, GR, HR, HU, IE, IT, LT, LU, LV, MT, NL, PL, PT, RO, SE, SI, SK}

func Get(cs []Country, id string) (Country, bool) {
	for _, c := range cs {
		if string(c) == id {
			return c, true
		}
	}
	return "", false
}

func InEuropeanUnion(country Country) bool {
	return slices.Contains(EuropeanUnion, country)
}

func (c Country) TranslateName(l lang.Lang) string {
	switch c {
	case AT:
		return l.Tr("Austria")
	case BE:
		return l.Tr("Belgium")
	case BG:
		return l.Tr("Bulgaria")
	case CH:
		return l.Tr("Switzerland")
	case CY:
		return l.Tr("Cyprus")
	case CZ:
		return l.Tr("Czechia")
	case DE:
		return l.Tr("Germany")
	case DK:
		return l.Tr("Denmark")
	case EE:
		return l.Tr("Estonia")
	case ES:
		return l.Tr("Spain")
	case FI:
		return l.Tr("Finland")
	case FR:
		return l.Tr("France")
	case GB:
		return l.Tr("United Kingdom")
	case GR:
		return l.Tr("Greece")
	case HR:
		return l.Tr("Croatia")
	case HU:
		return l.Tr("Hungary")
	case IE:
		return l.Tr("Ireland")
	case IT:
		return l.Tr("Italy")
	case LT:
		return l.Tr("Lithuania")
	case LU:
		return l.Tr("Luxembourg")
	case LV:
		return l.Tr("Latvia")
	case ME:
		return l.Tr("Montenegro")
	case MT:
		return l.Tr("Malta")
	case NL:
		return l.Tr("Netherlands")
	case PL:
		return l.Tr("Poland")
	case PT:
		return l.Tr("Portugal")
	case RO:
		return l.Tr("Romania")
	case SE:
		return l.Tr("Sweden")
	case SI:
		return l.Tr("Slovenia")
	case SK:
		return l.Tr("Slovakia")
	default:
		return string(c)
	}
}

// VAT rates are from: https://europa.eu/youreurope/business/taxation/vat/vat-rules-rates/index_en.htm#shortcut-5.
// Note that the ISO 3166-1 code for Greece is "GR", but the VAT rate table uses its ISO 639-1 code "EL".
func (c Country) VAT() VATRates {
	switch c {
	case AT:
		return map[Rate]float64{
			RateStandard: 0.20,
			RateReduced1: 0.10,
			RateReduced2: 0.13,
			RateParking:  0.13,
		}
	case BE:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.06,
			RateReduced2: 0.12,
			RateParking:  0.12,
		}
	case BG:
		return map[Rate]float64{
			RateStandard: 0.20,
			RateReduced1: 0.09,
		}
	case CY:
		return map[Rate]float64{
			RateStandard: 0.19,
			RateReduced1: 0.05,
			RateReduced2: 0.09,
		}
	case CZ:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.12,
			RateReduced2: 0.15,
		}
	case DE:
		return map[Rate]float64{
			RateStandard: 0.19,
			RateReduced1: 0.07,
		}
	case DK:
		return map[Rate]float64{
			RateStandard: 0.25,
		}
	case EE:
		return map[Rate]float64{
			RateStandard: 0.22,
			RateReduced1: 0.09,
		}
	case ES:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.10,
		}
	case FI:
		return map[Rate]float64{
			RateStandard: 0.255,
			RateReduced1: 0.10,
			RateReduced2: 0.14,
		}
	case FR:
		return map[Rate]float64{
			RateStandard:     0.20,
			RateReduced1:     0.055,
			RateReduced2:     0.10,
			RateSuperReduced: 0.021,
		}
	case GR:
		return map[Rate]float64{
			RateStandard: 0.24,
			RateReduced1: 0.06,
			RateReduced2: 0.13,
		}
	case HR:
		return map[Rate]float64{
			RateStandard: 0.25,
			RateReduced1: 0.05,
			RateReduced2: 0.13,
		}
	case HU:
		return map[Rate]float64{
			RateStandard: 0.27,
			RateReduced1: 0.05,
			RateReduced2: 0.18,
		}
	case IE:
		return map[Rate]float64{
			RateStandard: 0.23,
			RateReduced1: 0.09,
			RateReduced2: 0.135,
		}
	case IT:
		return map[Rate]float64{
			RateStandard:     0.22,
			RateReduced1:     0.05,
			RateReduced2:     0.10,
			RateSuperReduced: 0.04,
		}
	case LT:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.05,
			RateReduced2: 0.09,
		}
	case LU:
		return map[Rate]float64{
			RateStandard:     0.17,
			RateReduced1:     0.08,
			RateSuperReduced: 0.03,
			RateParking:      0.14,
		}
	case LV:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.05,
			RateReduced2: 0.12,
		}
	case MT:
		return map[Rate]float64{
			RateStandard: 0.18,
			RateReduced1: 0.05,
			RateReduced2: 0.07,
			RateParking:  0.12,
		}
	case NL:
		return map[Rate]float64{
			RateStandard: 0.21,
			RateReduced1: 0.09,
		}
	case PL:
		return map[Rate]float64{
			RateStandard: 0.23,
			RateReduced1: 0.05,
			RateReduced2: 0.08,
		}
	case PT:
		return map[Rate]float64{
			RateStandard: 0.23,
			RateReduced1: 0.06,
			RateReduced2: 0.13,
			RateParking:  0.13,
		}
	case RO:
		return map[Rate]float64{
			RateStandard: 0.19,
			RateReduced1: 0.05,
			RateReduced2: 0.09,
		}
	case SE:
		return map[Rate]float64{
			RateStandard: 0.25,
			RateReduced1: 0.06,
			RateReduced2: 0.12,
		}
	case SI:
		return map[Rate]float64{
			RateStandard: 0.22,
			RateReduced1: 0.05,
			RateReduced2: 0.095,
		}
	case SK:
		return map[Rate]float64{
			RateStandard: 0.20,
			RateReduced1: 0.10,
		}
	default:
		return nil
	}
}

type CountryOption struct {
	Country
	Name     string
	Selected bool
}

func TranslateAndSort(l lang.Lang, countries []Country, selected Country) []CountryOption {
	var result = make([]CountryOption, len(countries))
	for i := range countries {
		result[i] = CountryOption{
			Country:  countries[i],
			Name:     countries[i].TranslateName(l),
			Selected: countries[i] == selected,
		}
	}

	collator := l.Collator()
	sort.Slice(result, func(i, j int) bool {
		return collator.CompareString(result[i].Name, result[j].Name) < 0
	})
	return result
}
