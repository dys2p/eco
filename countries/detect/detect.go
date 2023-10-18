package detect

import (
	"fmt"
	"net/http"
	"slices"

	"golang.org/x/text/language"
)

var all = append([]string{"non-EU"}, eu...)

var eu = []string{"AT", "BE", "BG", "CY", "CZ", "DE", "DK", "EE", "ES", "FI", "FR", "GR", "HR", "HU", "IE", "IT", "LT", "LU", "LV", "MT", "NL", "PL", "PT", "RO", "SE", "SI", "SK"}

// Countries returns the union of possible countries for a given HTTP request, based on the client's Accept-Language header and IP address.
func Countries(r *http.Request) (map[string]any, error) {

	// assume that these methods return ISO 3166-1 codes
	accept, err := acceptLanguage(r)
	if err != nil {
		return nil, err
	}
	ip, err := ipAddress(r)
	if err != nil {
		return nil, err
	}

	var available = make(map[string]any)
	for _, country := range append(accept, ip...) {
		if !slices.Contains(all, country) {
			return nil, fmt.Errorf("invalid country code detected: %s", country)
		}
		available[country] = struct{}{}
	}
	return available, nil
}

type LangCountry struct {
	Tag     language.Tag
	Country []string
}

// all language tags
var langCountries = []LangCountry{
	LangCountry{language.Afrikaans, []string{"non-EU"}},
	LangCountry{language.Albanian, []string{"non-EU"}},
	LangCountry{language.AmericanEnglish, all}, // Mullvad Browser and Tor Browser send "Accept-Language: en-US,en;q=0.5"
	LangCountry{language.Amharic, []string{"non-EU"}},
	LangCountry{language.Arabic, []string{"non-EU"}},
	LangCountry{language.Armenian, []string{"non-EU"}},
	LangCountry{language.Azerbaijani, []string{"non-EU"}},
	LangCountry{language.Bengali, []string{"non-EU"}},
	LangCountry{language.BrazilianPortuguese, []string{"non-EU"}},
	LangCountry{language.BritishEnglish, []string{"non-EU"}},
	LangCountry{language.Bulgarian, []string{"BG"}},
	LangCountry{language.Burmese, []string{"non-EU"}},
	LangCountry{language.CanadianFrench, []string{"non-EU"}},
	LangCountry{language.Catalan, []string{"ES"}},
	LangCountry{language.Chinese, []string{"non-EU"}},
	LangCountry{language.Croatian, []string{"HR"}},
	LangCountry{language.Czech, []string{"CZ"}},
	LangCountry{language.Danish, []string{"DK"}},
	LangCountry{language.Dutch, []string{"NL"}},
	LangCountry{language.English, all}, // Mullvad Browser and Tor Browser send "Accept-Language: en-US,en;q=0.5"
	LangCountry{language.Estonian, []string{"EE"}},
	LangCountry{language.EuropeanPortuguese, []string{"PT"}},
	LangCountry{language.EuropeanSpanish, []string{"ES"}},
	LangCountry{language.Filipino, []string{"non-EU"}},
	LangCountry{language.Finnish, []string{"FI"}},
	LangCountry{language.French, []string{"FR", "non-EU"}},
	LangCountry{language.Georgian, []string{"non-EU"}},
	LangCountry{language.German, []string{"AT", "DE", "non-EU"}},
	LangCountry{language.Greek, []string{"GR"}},
	LangCountry{language.Gujarati, []string{"non-EU"}},
	LangCountry{language.Hebrew, []string{"non-EU"}},
	LangCountry{language.Hindi, []string{"non-EU"}},
	LangCountry{language.Hungarian, []string{"HU"}},
	LangCountry{language.Icelandic, []string{"non-EU"}},
	LangCountry{language.Indonesian, []string{"non-EU"}},
	LangCountry{language.Italian, []string{"IT"}},
	LangCountry{language.Japanese, []string{"non-EU"}},
	LangCountry{language.Kannada, []string{"non-EU"}},
	LangCountry{language.Kazakh, []string{"non-EU"}},
	LangCountry{language.Khmer, []string{"non-EU"}},
	LangCountry{language.Kirghiz, []string{"non-EU"}},
	LangCountry{language.Korean, []string{"non-EU"}},
	LangCountry{language.Lao, []string{"non-EU"}},
	LangCountry{language.LatinAmericanSpanish, []string{"non-EU"}},
	LangCountry{language.Latvian, []string{"LV"}},
	LangCountry{language.Lithuanian, []string{"LT"}},
	LangCountry{language.Macedonian, []string{"non-EU"}},
	LangCountry{language.Malayalam, []string{"non-EU"}},
	LangCountry{language.Malay, []string{"non-EU"}},
	LangCountry{language.Marathi, []string{"non-EU"}},
	LangCountry{language.ModernStandardArabic, []string{"non-EU"}},
	LangCountry{language.Mongolian, []string{"non-EU"}},
	LangCountry{language.Nepali, []string{"non-EU"}},
	LangCountry{language.Norwegian, []string{"non-EU"}},
	LangCountry{language.Persian, []string{"non-EU"}},
	LangCountry{language.Polish, []string{"PL"}},
	LangCountry{language.Portuguese, []string{"PT", "non-EU"}},
	LangCountry{language.Punjabi, []string{"non-EU"}},
	LangCountry{language.Romanian, []string{"RO"}},
	LangCountry{language.Russian, []string{"non-EU"}},
	LangCountry{language.SerbianLatin, []string{"non-EU"}},
	LangCountry{language.Serbian, []string{"non-EU"}},
	LangCountry{language.SimplifiedChinese, []string{"non-EU"}},
	LangCountry{language.Sinhala, []string{"non-EU"}},
	LangCountry{language.Slovak, []string{"SK"}},
	LangCountry{language.Slovenian, []string{"SI"}},
	LangCountry{language.Spanish, []string{"ES", "non-EU"}},
	LangCountry{language.Swahili, []string{"non-EU"}},
	LangCountry{language.Swedish, []string{"SE"}},
	LangCountry{language.Tamil, []string{"non-EU"}},
	LangCountry{language.Telugu, []string{"non-EU"}},
	LangCountry{language.Thai, []string{"non-EU"}},
	LangCountry{language.TraditionalChinese, []string{"non-EU"}},
	LangCountry{language.Turkish, []string{"non-EU"}},
	LangCountry{language.Ukrainian, []string{"non-EU"}},
	LangCountry{language.Urdu, []string{"non-EU"}},
	LangCountry{language.Uzbek, []string{"non-EU"}},
	LangCountry{language.Vietnamese, []string{"non-EU"}},
	LangCountry{language.Zulu, []string{"non-EU"}},
}

var langCountriesMatcher = func() language.Matcher {
	var tags []language.Tag
	for _, lc := range langCountries {
		tags = append(tags, lc.Tag)
	}
	return language.NewMatcher(tags)
}()
