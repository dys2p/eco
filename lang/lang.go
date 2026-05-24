// Package lang provides utilities for translating web services.
//
// The Path function and the Languages type help you to follow [Google's advice]
// to use "different URLs for each language version of a page rather than using
// cookies or browser settings to adjust the content language on the page",
// using the "Subdirectories with gTLD" URL structure where localized URLs start
// with e.g. "/en/".
//
// Generate translations with:
//
//	gotext-update-templates -srclang=en-US -lang=de-DE,en-US -out=catalog.go .
//
// Then use them in your code:
//
//	langs := lang.MakeLanguages(nil, "de", "en")
//	lang.Handle(http.DefaultServeMux, langs, "/", func(w http.ResponseWriter, r *http.Request) {
//		l, _, _ := langs.FromURL(r.URL)
//		l.Printer.Fprintf(w, "Hello World")
//	})
//
// As in the example, adding routes for each language is recommended over using route parameters with possibly conflicting rules.
//
// [Google's advice]: https://developers.google.com/search/docs/specialty/international/managing-multi-regional-sites
package lang

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type Lang struct {
	BCP47   string
	Prefix  string // URL prefix
	Printer *message.Printer
	Tag     language.Tag
}

// Collator creates a case-insensitive collator for l.Tag. The collator is not stored in Lang because it is not thread-safe (see https://github.com/golang/go/issues/57314).
func (l Lang) Collator() *collate.Collator {
	return collate.New(l.Tag, collate.IgnoreCase)
}

// String returns the l.Prefix to be used in URLs.
func (l Lang) String() string {
	return l.Prefix
}

func (l Lang) Tr(key message.Reference, a ...interface{}) string {
	return l.Printer.Sprintf(key, a...)
}

type Languages []Lang

// MakeLanguages takes a list of URL path prefixes used in your application (e. g.  "de", "en")
// in the alphabetical order of the dictionary keys in the catalog.
// If catalog is nil, then message.DefaultCatalog is used.
// MakeLanguages panics if len(prefixes) does not equal the number of languages in the catalog.
func MakeLanguages(catalog catalog.Catalog, prefixes ...string) Languages {
	if catalog == nil {
		catalog = message.DefaultCatalog
	}
	if len(prefixes) == 0 {
		panic("need at least one language prefix")
	}

	tags := catalog.Languages()
	if len(prefixes) != len(tags) {
		panic(fmt.Sprintf("got %d prefixes but catalog has %d languages", len(prefixes), len(tags)))
	}

	var langs = make(Languages, len(prefixes))
	for i, prefix := range prefixes {
		langs[i].BCP47 = tags[i].String()
		langs[i].Prefix = prefix
		langs[i].Printer = message.NewPrinter(tags[i], message.Catalog(catalog))
		langs[i].Tag = tags[i]
	}
	return langs
}

// ByPrefix returns the language whose prefix matches the given prefix.
func (langs Languages) ByPrefix(prefix string) (Lang, bool) {
	for _, l := range langs {
		if l.Prefix == prefix {
			return l, true
		}
	}
	return langs[0], false
}

// FromPath returns the language whose prefix matches the first segment of the path and the remaining path.
// If no language matches, it returns langs[0], the full path and false.
//
// Deprecated: Use ByPrefix or FromURL instead.
func (langs Languages) FromPath(path string) (Lang, string, bool) {
	path = strings.TrimLeft(path, "/")
	prefix, remainder, _ := strings.Cut(path, "/")
	for _, l := range langs {
		if l.Prefix == prefix {
			return l, remainder, true
		}
	}
	return langs[0], path, false
}

// FromURL returns the language whose prefix matches the first segment of the path, and the remaining path and query.
// If no language matches, it returns langs[0], the full path and query, and false.
func (langs Languages) FromURL(u *url.URL) (Lang, string, bool) {
	var path = strings.TrimLeft(u.Path, "/")
	prefix, remainder, _ := strings.Cut(path, "/")
	var query string
	if u.RawQuery != "" {
		query = "?" + u.RawQuery
	}
	for _, l := range langs {
		if l.Prefix == prefix {
			return l, remainder + query, true
		}
	}
	return langs[0], path + query, false
}
