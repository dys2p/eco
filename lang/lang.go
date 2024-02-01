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
//	langs := lang.MakeLanguages("de", "en")
//	for _, l := range langs {
//		http.HandleFunc("/"+l.Prefix, func(w http.ResponseWriter, r *http.Request) {
//			l, _ := langs.FromPath(r)
//			l.Printer.Fprintf(w, "Hello World")
//		})
//	}
//
// Then generate the translations with:
//
//	gotext-update-templates -srclang=en-US -lang=de-DE,en-US -out=catalog.go .
//
// [Google's advice]: https://developers.google.com/search/docs/specialty/international/managing-multi-regional-sites
package lang

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	fallbackCollator = collate.New(language.English, collate.IgnoreCase)
	fallbackPrinter  = message.NewPrinter(language.English)
)

type Lang struct {
	Collator *collate.Collator // case-insensitive
	Prefix   string
	Printer  *message.Printer
}

func (l Lang) Tr(key message.Reference, a ...interface{}) string {
	return l.Printer.Sprintf(key, a...)
}

type Languages []Lang

// MakeLanguages takes a list of URL path prefixes used in your application (e. g.  "de", "en")
// in the alphabetical order of the dictionary keys in the default catalog (message.DefaultCatalog).
// MakeLanguages panics if len(prefixes) does not equal the number of languages in the default catalog.
func MakeLanguages(prefixes ...string) Languages {
	tags := message.DefaultCatalog.Languages()
	if len(prefixes) != len(tags) {
		panic(fmt.Sprintf("MakeLanguages got %d prefixes but catalog has %d languages", len(prefixes), len(tags)))
	}

	var langs = make(Languages, len(prefixes))
	for i, prefix := range prefixes {
		langs[i].Collator = collate.New(tags[i], collate.IgnoreCase)
		langs[i].Prefix = prefix
		langs[i].Printer = message.NewPrinter(tags[i])
	}
	return langs
}

// FromPath returns the language prefix of r.URL.Path and a message printer for it.
// The boolean return value indicates whether the path has a known prefix. If it has not, the returned prefix is the fallback prefix.
func (langs Languages) FromPath(path string) (Lang, bool) {
	path = strings.TrimLeft(path, "/")
	prefix, _, _ := strings.Cut(path, "/")
	for _, l := range langs {
		if l.Prefix == prefix {
			return l, true
		}
	}
	// fix prefix if possible
	if len(langs) > 0 {
		prefix = langs[0].Prefix
	}
	return Lang{
		Collator: fallbackCollator,
		Prefix:   prefix,
		Printer:  fallbackPrinter,
	}, false
}

// Redirect redirects to the localized version of r.URL according to the Accept-Language header.
// If r.URL it is already localized, Redirect responds with a "not found" error.
// It is recommended to chain Redirect behind your http router.
// Matching is done with message.DefaultCatalog.Matcher().
func (langs Languages) Redirect(w http.ResponseWriter, r *http.Request) {
	if _, ok := langs.FromPath(r.URL.Path); ok {
		// url already starts with a supported language, prevent redirect loop
		http.NotFound(w, r)
	} else {
		_, index := language.MatchStrings(message.DefaultCatalog.Matcher(), r.Header.Get("Accept-Language"))
		http.Redirect(w, r, path.Join("/", langs[index].Prefix, r.URL.Path), http.StatusSeeOther)
	}
}
