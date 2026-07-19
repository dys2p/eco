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

	// handle /pattern
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

	// handle /lang/pattern
	for _, l := range langs {
		mux.HandleFunc(insert(pattern, l.Prefix), handler)
	}
}

// insert inserts pathPrefix into pattern ("In general, a pattern looks like [METHOD ][HOST]/[PATH]").
// pathPrefix must not be empty.
func insert(pattern, pathPrefix string) string {
	firstSlash := strings.Index(pattern, "/")
	if firstSlash < 0 {
		return pattern // pattern is invalid, return unchanged
	}
	methodHost, path := pattern[:firstSlash], pattern[firstSlash:]
	if path == "/{$}" {
		path = "" // see TestInsert
	}
	return methodHost + "/" + pathPrefix + path
}
