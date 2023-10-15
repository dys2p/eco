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

// Path formats a string and prepends l to it. Note that Path does not validate l.
func (l Lang) Path(format string, a ...any) string {
	return path.Join("/", string(l), fmt.Sprintf(format, a...))
}

// Tr translates the given input text.
//
// l may be any language-like value because Tr calls MatchLanguage, which in turn calls MatchStrings.
func (l Lang) Tr(key message.Reference, a ...interface{}) string {
	return message.NewPrinter(message.MatchLanguage(string(l))).Sprintf(key, a...)
}

// Languages can be used to store a list of language url prefixes supported by your application.
type Languages []string

// getPrefix returns the first url segment and whether it is contained in langs.
func (langs Languages) getPrefix(reqpath string) (string, bool) {
	reqpath = strings.TrimLeft(reqpath, "/")
	first, _, _ := strings.Cut(reqpath, "/")
	return first, slices.Contains(langs, first)
}

// Get returns the language which matches the r.URL prefix.
func (langs Languages) Get(r *http.Request) Lang {
	langstr, ok := langs.getPrefix(r.URL.Path)
	if !ok {
		langstr = ""
	}
	return langs.get(langstr)
}

func (langs Languages) get(langstr string) Lang {
	_, index := language.MatchStrings(message.DefaultCatalog.Matcher(), langstr)
	return Lang(langs[index])
}

// Redirect redirects to the localized version of r.URL according to the Accept-Language header.
// If r.URL it is already localized, Redirect responds with a "not found" error.
// It is recommended to chain Redirect behind your http router.
func (langs Languages) Redirect(w http.ResponseWriter, r *http.Request) {
	if _, ok := langs.getPrefix(r.URL.Path); ok {
		http.NotFound(w, r) // url already starts with a supported language, prevent redirect loop
	} else {
		l := langs.get(r.Header.Get("Accept-Language"))
		http.Redirect(w, r, l.Path(r.URL.Path), http.StatusSeeOther)
	}
}
