package detect

import (
	"net/http"

	"github.com/dys2p/eco/countries"
	"golang.org/x/text/language"
)

// langCountries contains all language tags which are available in golang
//
// Noteworthy official languages in the European Union:
//
// * AT: German, Croatian, Slovenian, Hungarian
// * BE: Dutch, French, German
// * CY: Greek, Turkish
// * FI: Finnish⁠, Swedish
// * IE: Irish (not available in golang), English
// * IT: Italian, German, French, Slovenian⁠
// * LU: Luxembourgish (not available in golang),⁠ French, German
// * MT: Maltese (not available in golang), English
// * RO: Romanian, Hungarian
// * SE: Swedish, Finnish
// * SI: Slovenian, Italian, Hungarian
//
// Nil means that the country can't determined.
var langCountries = []struct {
	tag       language.Tag
	countries []countries.Country
}{
	// English can be anywhere because e.g. Mullvad Browser and Tor Browser hide the user language by default and set "Accept-Language: en-US,en;q=0.5".
	{language.AmericanEnglish, nil},
	{language.English, nil},

	// official languages of the European Union
	{language.Bulgarian, []countries.Country{countries.BG}},
	{language.Catalan, []countries.Country{countries.ES, countries.NonEU}}, // non-EU: Andorra
	{language.Croatian, []countries.Country{countries.AT, countries.HR}},
	{language.Czech, []countries.Country{countries.CZ}},
	{language.Danish, []countries.Country{countries.DK}},
	{language.Dutch, []countries.Country{countries.BE, countries.NL}},
	{language.Estonian, []countries.Country{countries.EE}},
	{language.EuropeanPortuguese, []countries.Country{countries.PT}},
	{language.EuropeanSpanish, []countries.Country{countries.ES}},
	{language.Finnish, []countries.Country{countries.FI, countries.SE}},
	{language.French, []countries.Country{countries.BE, countries.FR, countries.IT, countries.LU, countries.NonEU}},               // non-EU: world language
	{language.German, []countries.Country{countries.AT, countries.BE, countries.DE, countries.IT, countries.LU, countries.NonEU}}, // non-EU: Liechtenstein, Switzerland
	{language.Greek, []countries.Country{countries.CY, countries.GR}},
	{language.Hungarian, []countries.Country{countries.AT, countries.HU, countries.RO, countries.SI}},
	{language.Italian, []countries.Country{countries.IT, countries.SI}},
	{language.Latvian, []countries.Country{countries.LV}},
	{language.Lithuanian, []countries.Country{countries.LT}},
	{language.Polish, []countries.Country{countries.PL}},
	{language.Portuguese, []countries.Country{countries.PT, countries.NonEU}}, // non-EU: world language
	{language.Romanian, []countries.Country{countries.RO}},
	{language.Slovak, []countries.Country{countries.SK}},
	{language.Slovenian, []countries.Country{countries.AT, countries.IT, countries.SI}},
	{language.Spanish, []countries.Country{countries.ES, countries.NonEU}}, // non-EU: world language
	{language.Swedish, []countries.Country{countries.FI, countries.SE}},
	{language.Turkish, []countries.Country{countries.CY, countries.NonEU}}, // non-EU: Turkey

	// languages which are no official languages of the European Union
	{language.Afrikaans, []countries.Country{countries.NonEU}},
	{language.Albanian, []countries.Country{countries.NonEU}},
	{language.Amharic, []countries.Country{countries.NonEU}},
	{language.Arabic, []countries.Country{countries.NonEU}},
	{language.Armenian, []countries.Country{countries.NonEU}},
	{language.Azerbaijani, []countries.Country{countries.NonEU}},
	{language.Bengali, []countries.Country{countries.NonEU}},
	{language.BrazilianPortuguese, []countries.Country{countries.NonEU}},
	{language.BritishEnglish, []countries.Country{countries.NonEU}},
	{language.Burmese, []countries.Country{countries.NonEU}},
	{language.CanadianFrench, []countries.Country{countries.NonEU}},
	{language.Chinese, []countries.Country{countries.NonEU}},
	{language.Filipino, []countries.Country{countries.NonEU}},
	{language.Georgian, []countries.Country{countries.NonEU}},
	{language.Gujarati, []countries.Country{countries.NonEU}},
	{language.Hebrew, []countries.Country{countries.NonEU}},
	{language.Hindi, []countries.Country{countries.NonEU}},
	{language.Icelandic, []countries.Country{countries.NonEU}},
	{language.Indonesian, []countries.Country{countries.NonEU}},
	{language.Japanese, []countries.Country{countries.NonEU}},
	{language.Kannada, []countries.Country{countries.NonEU}},
	{language.Kazakh, []countries.Country{countries.NonEU}},
	{language.Khmer, []countries.Country{countries.NonEU}},
	{language.Kirghiz, []countries.Country{countries.NonEU}},
	{language.Korean, []countries.Country{countries.NonEU}},
	{language.Lao, []countries.Country{countries.NonEU}},
	{language.LatinAmericanSpanish, []countries.Country{countries.NonEU}},
	{language.Macedonian, []countries.Country{countries.NonEU}},
	{language.Malayalam, []countries.Country{countries.NonEU}},
	{language.Malay, []countries.Country{countries.NonEU}},
	{language.Marathi, []countries.Country{countries.NonEU}},
	{language.ModernStandardArabic, []countries.Country{countries.NonEU}},
	{language.Mongolian, []countries.Country{countries.NonEU}},
	{language.Nepali, []countries.Country{countries.NonEU}},
	{language.Norwegian, []countries.Country{countries.NonEU}},
	{language.Persian, []countries.Country{countries.NonEU}},
	{language.Punjabi, []countries.Country{countries.NonEU}},
	{language.Russian, []countries.Country{countries.NonEU}},
	{language.SerbianLatin, []countries.Country{countries.NonEU}},
	{language.Serbian, []countries.Country{countries.NonEU}},
	{language.SimplifiedChinese, []countries.Country{countries.NonEU}},
	{language.Sinhala, []countries.Country{countries.NonEU}},
	{language.Swahili, []countries.Country{countries.NonEU}},
	{language.Tamil, []countries.Country{countries.NonEU}},
	{language.Telugu, []countries.Country{countries.NonEU}},
	{language.Thai, []countries.Country{countries.NonEU}},
	{language.TraditionalChinese, []countries.Country{countries.NonEU}},
	{language.Ukrainian, []countries.Country{countries.NonEU}},
	{language.Urdu, []countries.Country{countries.NonEU}},
	{language.Uzbek, []countries.Country{countries.NonEU}},
	{language.Vietnamese, []countries.Country{countries.NonEU}},
	{language.Zulu, []countries.Country{countries.NonEU}},
}

var matcher language.Matcher

func init() {
	var tags []language.Tag
	for _, lc := range langCountries {
		tags = append(tags, lc.tag)
	}
	matcher = language.NewMatcher(tags)
}

func acceptLanguage(r *http.Request) ([]countries.Country, error) {
	tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		return nil, err
	}
	_, index, _ := matcher.Match(tags...) // single best match for accept-languages
	return langCountries[index].countries, nil
}
