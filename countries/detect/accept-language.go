package detect

import (
	"net/http"

	"golang.org/x/text/language"
)

func acceptLanguage(r *http.Request) ([]string, error) {
	tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		return nil, err
	}
	_, index, _ := langCountriesMatcher.Match(tags...) // one best match for accept-languages
	return langCountries[index].Country, nil
}
