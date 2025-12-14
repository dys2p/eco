// Package countries defines variables and translations for European Union countries and some others.
package countries

import (
	"fmt"
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
	MK Country = "MK"
)

var All = []Country{AT, BE, BG, CH, CY, CZ, DE, DK, EE, ES, FI, FR, GB, GR, HR, HU, IE, IT, LT, LU, LV, ME, MK, MT, NL, PL, PT, RO, SE, SI, SK}

var EuropeanUnion = []Country{AT, BE, BG, CY, CZ, DE, DK, EE, ES, FI, FR, GR, HR, HU, IE, IT, LT, LU, LV, MT, NL, PL, PT, RO, SE, SI, SK}

// Get returns the country with the given id from the slice, or ("", false).
func Get(cs []Country, id string) (Country, bool) {
	for _, c := range cs {
		if string(c) == id {
			return c, true
		}
	}
	return "", false
}

func (c Country) InEU() bool {
	return slices.Contains(EuropeanUnion, c)
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
	case MK:
		return l.Tr("North Macedonia")
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

// MarshalText implements encoding.TextMarshaler.
func (c Country) MarshalText() (text []byte, err error) {
	return []byte(c), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (c *Country) UnmarshalText(text []byte) error {
	// allow zero value
	if len(text) == 0 {
		*c = Country("")
		return nil
	}

	got, ok := Get(All, string(text))
	if !ok {
		return fmt.Errorf("unmarshalling country: not found: %s", string(text))
	}
	*c = got
	return nil
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
