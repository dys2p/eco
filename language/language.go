// Package language provides utilities for translating web services.
package language

import (
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/message"
)

var PathPrefixes = []string{"de", "en"}

type Lang string

// Get returns the first URL path item (if it is in PathPrefixes) or the HTTP "Accept-Language" header value.
func Get(r *http.Request) Lang {
	relpath := strings.TrimPrefix(r.URL.Path, "/")
	pathprefix, _, _ := strings.Cut(relpath, "/")
	if slices.Contains(PathPrefixes, pathprefix) {
		return Lang(pathprefix)
	}
	return Lang(r.Header.Get("Accept-Language"))
}

// Tr translates the given input text.
func (lang Lang) Tr(key message.Reference, a ...interface{}) string {
	return message.NewPrinter(message.MatchLanguage(string(lang))).Sprintf(key, a...)
}
