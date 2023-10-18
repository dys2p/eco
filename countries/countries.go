package countries

import (
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// Country Codes are ISO 3166-1.
// VAT rates are from: https://europa.eu/youreurope/business/taxation/vat/vat-rules-rates/index_en.htm#shortcut-5.
// Note that the ISO 3166-1 code for Greece is "GR", but the VAT rate table uses its ISO 639-1 code "EL".
var (
	AT = Country{"AT", map[string]float64{
		"standard":  0.20,
		"reduced-1": 0.10,
		"reduced-2": 0.13,
		"parking":   0.13,
	}}
	BE = Country{"BE", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.06,
		"reduced-2": 0.12,
		"parking":   0.12,
	}}
	BG = Country{"BG", map[string]float64{
		"standard":  0.20,
		"reduced-1": 0.09,
	}}
	CY = Country{"CY", map[string]float64{
		"standard":  0.19,
		"reduced-1": 0.05,
		"reduced-2": 0.09,
	}}
	CZ = Country{"CZ", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.10,
		"reduced-2": 0.15,
	}}
	DE = Country{"DE", map[string]float64{
		"standard":  0.19,
		"reduced-1": 0.07,
	}}
	DK = Country{"DK", map[string]float64{
		"standard": 0.25,
	}}
	EE = Country{"EE", map[string]float64{
		"standard":  0.20,
		"reduced-1": 0.09,
	}}
	ES = Country{"ES", map[string]float64{
		"standard":      0.21,
		"reduced-1":     0.10,
		"super-reduced": 0.04,
	}}
	FI = Country{"FI", map[string]float64{
		"standard":  0.24,
		"reduced-1": 0.10,
		"reduced-2": 0.14,
	}}
	FR = Country{"FR", map[string]float64{
		"standard":      0.20,
		"reduced-1":     0.055,
		"reduced-2":     0.10,
		"super-reduced": 0.021,
	}}
	GR = Country{"GR", map[string]float64{
		"standard":  0.24,
		"reduced-1": 0.06,
		"reduced-2": 0.13,
	}}
	HR = Country{"HR", map[string]float64{
		"standard":  0.25,
		"reduced-1": 0.05,
		"reduced-2": 0.13,
	}}
	HU = Country{"HU", map[string]float64{
		"standard":  0.27,
		"reduced-1": 0.05,
		"reduced-2": 0.18,
	}}
	IE = Country{"IE", map[string]float64{
		"standard":      0.23,
		"reduced-1":     0.09,
		"reduced-2":     0.135,
		"super-reduced": 0.048,
		"parking":       0.135,
	}}
	IT = Country{"IT", map[string]float64{
		"standard":      0.22,
		"reduced-1":     0.05,
		"reduced-2":     0.10,
		"super-reduced": 0.04,
	}}
	LT = Country{"LT", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.05,
		"reduced-2": 0.09,
	}}
	LU = Country{"LU", map[string]float64{
		"standard":      0.17,
		"reduced-1":     0.08,
		"super-reduced": 0.03,
		"parking":       0.14,
	}}
	LV = Country{"LV", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.12,
		"reduced-2": 0.05,
	}}
	MT = Country{"MT", map[string]float64{
		"standard":  0.18,
		"reduced-1": 0.05,
		"reduced-2": 0.07,
	}}
	NL = Country{"NL", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.09,
	}}
	PL = Country{"PL", map[string]float64{
		"standard":  0.23,
		"reduced-1": 0.05,
		"reduced-2": 0.08,
	}}
	PT = Country{"PT", map[string]float64{
		"standard":  0.23,
		"reduced-1": 0.06,
		"reduced-2": 0.13,
		"parking":   0.13,
	}}
	RO = Country{"RO", map[string]float64{
		"standard":  0.19,
		"reduced-1": 0.05,
		"reduced-2": 0.09,
	}}
	SE = Country{"SE", map[string]float64{
		"standard":  0.25,
		"reduced-1": 0.06,
		"reduced-2": 0.12,
	}}
	SI = Country{"SI", map[string]float64{
		"standard":  0.22,
		"reduced-1": 0.05,
		"reduced-2": 0.095,
	}}
	SK = Country{"SK", map[string]float64{
		"standard":  0.20,
		"reduced-1": 0.10,
	}}
)

var EuropeanUnion = []Country{AT, BE, BG, CY, CZ, DE, DK, EE, ES, FI, FR, GR, HR, HU, IE, IT, LT, LU, LV, MT, NL, PL, PT, RO, SE, SI, SK}

var (
	CH = Country{"CH", nil}
	GB = Country{"GB", nil}
)

var All = []Country{AT, BE, BG, CH, CY, CZ, DE, DK, EE, ES, FI, FR, GB, GR, HR, HU, IE, IT, LT, LU, LV, MT, NL, PL, PT, RO, SE, SI, SK}

func Get(id string) (Country, bool) {
	for _, c := range All {
		if c.ID == id {
			return c, true
		}
	}
	return Country{}, false
}

func TranslateAndSort(langstr string, countries []Country) []CountryWithName {
	var result = make([]CountryWithName, len(countries))
	for i := range countries {
		result[i] = CountryWithName{
			Country: countries[i],
			Name:    countries[i].TranslateName(langstr),
		}
	}

	collator := collate.New(language.Make(langstr), collate.IgnoreCase)
	sort.Slice(result, func(i, j int) bool {
		return collator.CompareString(result[i].Name, result[j].Name) < 0
	})
	return result
}

type CountryWithName struct {
	Country
	Name string
}
