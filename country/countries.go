// Package country contains European countries and their VAT rates.
package country

var (
	AT = Country{"AT", "Österreich", 0.20}
	BE = Country{"BE", "Belgien", 0.21}
	BG = Country{"BG", "Bulgarien", 0.20}
	CY = Country{"CY", "Zypern", 0.21}
	CZ = Country{"CZ", "Tschechische Republik", 0.21}
	DE = Country{"DE", "Deutschland", 0.19}
	DK = Country{"DK", "Dänemark", 0.25}
	EE = Country{"EE", "Estland", 0.20}
	ES = Country{"ES", "Spanien", 0.21}
	FI = Country{"FI", "Finnland", 0.24}
	FR = Country{"FR", "Frankreich", 0.20}
	GR = Country{"GR", "Griechenland", 0.24}
	HR = Country{"HR", "Kroatien", 0.25}
	HU = Country{"HU", "Ungarn", 0.27}
	IE = Country{"IE", "Irland", 0.23}
	IT = Country{"IT", "Italien", 0.22}
	LT = Country{"LT", "Litauen", 0.21}
	LU = Country{"LU", "Luxemburg", 0.17}
	LV = Country{"LV", "Lettland", 0.21}
	MT = Country{"MT", "Malta", 0.18}
	NL = Country{"NL", "Niederlande", 0.21}
	PL = Country{"PL", "Polen", 0.23}
	PT = Country{"PT", "Portugal", 0.23}
	RO = Country{"RO", "Rumänien", 0.19}
	SE = Country{"SE", "Schweden", 0.25}
	SI = Country{"SI", "Slowenien", 0.22}
	SK = Country{"SK", "Slowakei", 0.20}
)

var EuropeanUnion = Sort([]Country{AT, BE, BG, CY, CZ, DE, DK, EE, ES, FI, FR, GR, HR, HU, IE, IT, LT, LU, LV, MT, NL, PL, PT, RO, SE, SI, SK})

var (
	CH = Country{"CH", "Schweiz", 0}
	GB = Country{"GB", "Großbritannien", 0}
)

var All = Sort([]Country{AT, BE, BG, CH, CY, CZ, DE, DK, EE, ES, FI, FR, GB, GR, HR, HU, IE, IT, LT, LU, LV, MT, NL, PL, PT, RO, SE, SI, SK})
