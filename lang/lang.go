// Package lang provides utilities for translating web services.
//
// The Path function and the Languages type help you to follow [Google's advice]
// to use  "different URLs for each language version of a page rather than using
// cookies or browser settings to adjust the content language on the page",
// using the "Subdirectories with gTLD" URL structure where localized URLs start
// with e.g. "/en/".
//
// [Google's advice]: https://developers.google.com/search/docs/specialty/international/managing-multi-regional-sites
package lang

import (
	"fmt"
	"net/http"
	"path"
	"slices"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Lang string

// Path formats a string and prepends l to it.
func (l Lang) Path(format string, a ...any) string {
	return path.Join("/", string(l), fmt.Sprintf(format, a...))
}

// String returns l as a string.
func (l Lang) String() string {
	return string(l)
}

// Tr translates the given input text.
//
// l may be any BCP47-like string because Tr calls MatchLanguage, which in turn calls MatchStrings.
func (l Lang) Tr(key message.Reference, a ...interface{}) string {
	return message.NewPrinter(message.MatchLanguage(string(l))).Sprintf(key, a...)
}

// Languages can be used to store a list of language url prefixes supported by your application.
type Languages []string

// getPrefix returns the first url segment and whether it is contained in langs.
func (langs Languages) getPrefix(reqpath string) (string, bool) {
	reqpath = strings.TrimLeft(reqpath, "/")
	prefix, _, _ := strings.Cut(reqpath, "/")
	return prefix, slices.Contains(langs, prefix)
}

// ByPath returns the language which matches the r.URL.Path prefix.
func (langs Languages) ByPath(r *http.Request) Lang {
	if langstr, ok := langs.getPrefix(r.URL.Path); ok {
		return Lang(langstr)
	}
	if len(langs) > 0 {
		return Lang(langs[0])
	}
	return Lang("en")
}

// Redirect redirects to the localized version of r.URL according to the Accept-Language header.
// If r.URL it is already localized, Redirect responds with a "not found" error.
// It is recommended to chain Redirect behind your http router.
func (langs Languages) Redirect(w http.ResponseWriter, r *http.Request) {
	if _, ok := langs.getPrefix(r.URL.Path); ok {
		http.NotFound(w, r) // url already starts with a supported language, prevent redirect loop
	} else {
		var supported = make([]language.Tag, len(langs))
		for i := range langs {
			supported[i] = language.Make(langs[i])
		}
		matcher := language.NewMatcher(supported)
		_, index := language.MatchStrings(matcher, r.Header.Get("Accept-Language"))
		prefix := langs[index]
		http.Redirect(w, r, path.Join("/", prefix, r.URL.Path), http.StatusSeeOther)
	}
}
