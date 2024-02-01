// Package lang provides utilities for translating web services.
//
// The Path function and the Languages type help you to follow [Google's advice]
// to use "different URLs for each language version of a page rather than using
// cookies or browser settings to adjust the content language on the page",
// using the "Subdirectories with gTLD" URL structure where localized URLs start
// with e.g. "/en/".
//
// Adding routes for each language is recommended over using route parameters with possibly conflicting rules. Example:
//
//	for _, prefix := range []string{"en", "de"} {
//		http.HandleFunc("/"+prefix, func(w http.ResponseWriter, r *http.Request) {
//			_, printer, _ := langs.FromPath(r)
//			printer.Fprintf(w, "Hello World")
//		})
//	}
//
// Then generate the translations with:
//
//	gotext-update-templates -srclang=en-US -lang=en-US,de-DE -out=catalog.go .
//
// [Google's advice]: https://developers.google.com/search/docs/specialty/international/managing-multi-regional-sites
package lang

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var fallback = message.NewPrinter(language.English)

type Languages []struct {
	Prefix  string
	Printer *message.Printer
}

// MakeLanguages takes a list of URL path prefixes used in your application (e. g. "en", "de").
// Their order must match the default catalog (message.DefaultCatalog).
// MakeLanguages panics if len(prefixes) does not equal the number of languages in the default catalog.
func MakeLanguages(prefixes ...string) Languages {
	tags := message.DefaultCatalog.Languages()
	if len(prefixes) != len(tags) {
		panic(fmt.Sprintf("MakeLanguages got %d prefixes but catalog has %d languages", len(prefixes), len(tags)))
	}

	var langs = make(Languages, len(prefixes))
	for i, prefix := range prefixes {
		langs[i].Prefix = prefix
		langs[i].Printer = message.NewPrinter(tags[i])
	}
	return langs
}

// FromPath returns the language prefix of r.URL.Path and a message printer for it.
// The boolean return value indicates whether the path has a known prefix. If it has not, the returned prefix is the fallback prefix.
func (langs Languages) FromPath(r *http.Request) (string, *message.Printer, bool) {
	reqpath := strings.TrimLeft(r.URL.Path, "/")
	prefix, _, _ := strings.Cut(reqpath, "/")
	for _, l := range langs {
		if l.Prefix == prefix {
			return prefix, l.Printer, true
		}
	}
	// fix prefix if possible
	if len(langs) > 0 {
		prefix = langs[0].Prefix
	}
	return prefix, fallback, false
}

// Redirect redirects to the localized version of r.URL according to the Accept-Language header.
// If r.URL it is already localized, Redirect responds with a "not found" error.
// It is recommended to chain Redirect behind your http router.
// Matching is done with message.DefaultCatalog.Matcher().
func (langs Languages) Redirect(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := langs.FromPath(r); ok {
		// url already starts with a supported language, prevent redirect loop
		http.NotFound(w, r)
	} else {
		_, index := language.MatchStrings(message.DefaultCatalog.Matcher(), r.Header.Get("Accept-Language"))
		http.Redirect(w, r, path.Join("/", langs[index].Prefix, r.URL.Path), http.StatusSeeOther)
	}
}
