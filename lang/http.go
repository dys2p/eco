package lang

import (
	"net/http"
	"path"
	"strings"

	"golang.org/x/text/language"
)

// Handle is a shortcut for HandleFunc(mux, langs, pattern, handler.ServeHTTP).
func Handle(mux *http.ServeMux, langs Languages, pattern string, handler http.Handler) {
	HandleFunc(mux, langs, pattern, handler.ServeHTTP)
}

// HandleFunc registers /lang/pattern for each language, and a redirect handler for /pattern.
//
// It registers each route explicitly rather than redirecting 404s to /default-lang/pattern, so we don't mess with chained handlers.
func HandleFunc(mux *http.ServeMux, langs Languages, pattern string, handler http.HandlerFunc) {
	// no language support if langs is empty
	if len(langs) == 0 {
		mux.HandleFunc(pattern, handler)
		return
	}

	// handle /lang/pattern (insert lang at the right place: "In general, a pattern looks like [METHOD ][HOST]/[PATH]")
	insertAt := strings.Index(pattern, "/")
	if insertAt < 0 {
		insertAt = len(pattern)
	}
	for _, l := range langs {
		mux.HandleFunc(pattern[:insertAt]+"/"+l.Prefix+pattern[insertAt:], handler)
	}

	// handle /pattern (panics if pattern is empty, which is correct behavior)
	var tags = make([]language.Tag, len(langs))
	for i := range tags {
		tags[i] = langs[i].Tag
	}
	matcher := language.NewMatcher(tags)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		_, index := language.MatchStrings(matcher, r.Header.Get("Accept-Language"))
		var u = *r.URL // copy
		u.Path = path.Join("/", langs[index].Prefix, u.Path)
		http.Redirect(w, r, u.String(), http.StatusSeeOther)
	})
}
