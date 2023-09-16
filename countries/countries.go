package countries

import "sort"

// Country Codes are ISO 3166-1.
// VAT rates are from: https://europa.eu/youreurope/business/taxation/vat/vat-rules-rates/index_en.htm#shortcut-5.
// Note that the ISO 3166-1 code for Greece is "GR", but the VAT rate table uses its ISO 639-1 code "EL".
var (
	AT = Country{"AT", "Österreich", map[string]float64{
		"standard":  0.20,
		"reduced-1": 0.10,
		"reduced-2": 0.13,
		"parking":   0.13,
	}}
	BE = Country{"BE", "Belgien", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.06,
		"reduced-2": 0.12,
		"parking":   0.12,
	}}
	BG = Country{"BG", "Bulgarien", map[string]float64{
		"standard":  0.20,
		"reduced-1": 0.09,
	}}
	CY = Country{"CY", "Zypern", map[string]float64{
		"standard":  0.19,
		"reduced-1": 0.05,
		"reduced-2": 0.09,
	}}
	CZ = Country{"CZ", "Tschechische Republik", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.10,
		"reduced-2": 0.15,
	}}
	DE = Country{"DE", "Deutschland", map[string]float64{
		"standard":  0.19,
		"reduced-1": 0.07,
	}}
	DK = Country{"DK", "Dänemark", map[string]float64{
		"standard": 0.25,
	}}
	EE = Country{"EE", "Estland", map[string]float64{
		"standard":  0.20,
		"reduced-1": 0.09,
	}}
	ES = Country{"ES", "Spanien", map[string]float64{
		"standard":      0.21,
		"reduced-1":     0.10,
		"super-reduced": 0.04,
	}}
	FI = Country{"FI", "Finnland", map[string]float64{
		"standard":  0.24,
		"reduced-1": 0.10,
		"reduced-2": 0.14,
	}}
	FR = Country{"FR", "Frankreich", map[string]float64{
		"standard":      0.20,
		"reduced-1":     0.055,
		"reduced-2":     0.10,
		"super-reduced": 0.021,
	}}
	GR = Country{"GR", "Griechenland", map[string]float64{
		"standard":  0.24,
		"reduced-1": 0.06,
		"reduced-2": 0.13,
	}}
	HR = Country{"HR", "Kroatien", map[string]float64{
		"standard":  0.25,
		"reduced-1": 0.05,
		"reduced-2": 0.13,
	}}
	HU = Country{"HU", "Ungarn", map[string]float64{
		"standard":  0.27,
		"reduced-1": 0.05,
		"reduced-2": 0.18,
	}}
	IE = Country{"IE", "Irland", map[string]float64{
		"standard":      0.23,
		"reduced-1":     0.09,
		"reduced-2":     0.135,
		"super-reduced": 0.048,
		"parking":       0.135,
	}}
	IT = Country{"IT", "Italien", map[string]float64{
		"standard":      0.22,
		"reduced-1":     0.05,
		"reduced-2":     0.10,
		"super-reduced": 0.04,
	}}
	LT = Country{"LT", "Litauen", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.05,
		"reduced-2": 0.09,
	}}
	LU = Country{"LU", "Luxemburg", map[string]float64{
		"standard":      0.17,
		"reduced-1":     0.08,
		"super-reduced": 0.03,
		"parking":       0.14,
	}}
	LV = Country{"LV", "Lettland", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.12,
		"reduced-2": 0.05,
	}}
	MT = Country{"MT", "Malta", map[string]float64{
		"standard":  0.18,
		"reduced-1": 0.05,
		"reduced-2": 0.07,
	}}
	NL = Country{"NL", "Niederlande", map[string]float64{
		"standard":  0.21,
		"reduced-1": 0.09,
	}}
	PL = Country{"PL", "Polen", map[string]float64{
		"standard":  0.23,
		"reduced-1": 0.05,
		"reduced-2": 0.08,
	}}
	PT = Country{"PT", "Portugal", map[string]float64{
		"standard":  0.23,
		"reduced-1": 0.06,
		"reduced-2": 0.13,
		"parking":   0.13,
	}}
	RO = Country{"RO", "Rumänien", map[string]float64{
		"standard":  0.19,
		"reduced-1": 0.05,
		"reduced-2": 0.09,
	}}
	SE = Country{"SE", "Schweden", map[string]float64{
		"standard":  0.25,
		"reduced-1": 0.06,
		"reduced-2": 0.12,
	}}
	SI = Country{"SI", "Slowenien", map[string]float64{
		"standard":  0.22,
		"reduced-1": 0.05,
		"reduced-2": 0.095,
	}}
	SK = Country{"SK", "Slowakei", map[string]float64{
		"standard":  0.20,
		"reduced-1": 0.10,
	}}
)

var EuropeanUnion = []Country{AT, BE, BG, CY, CZ, DE, DK, EE, ES, FI, FR, GR, HR, HU, IE, IT, LT, LU, LV, MT, NL, PL, PT, RO, SE, SI, SK}

var (
	CH = Country{"CH", "Schweiz", nil}
	GB = Country{"GB", "Großbritannien", nil}
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

func SortByName(countries []Country) []Country {
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})
	return countries
}
