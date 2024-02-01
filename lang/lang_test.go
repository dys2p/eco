package lang

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

var langs Languages

func init() {
	english := language.MustParse("en-US")
	german := language.MustParse("de-DE")
	var b = catalog.NewBuilder(catalog.Fallback(english))
	b.SetString(english, "Hello World", "Hello World")
	b.SetString(german, "Hello World", "Hallo Welt")
	message.DefaultCatalog = b

	langs = MakeLanguages("en", "de")
}

func TestPrefixes(t *testing.T) {
	want := []string{"en", "de"}
	got := langs.Prefixes()
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRedirect(t *testing.T) {
	tests := map[string]string{
		"de":                    "/de",
		"de-CH":                 "/de",
		"de-DE":                 "/de",
		"":                      "/en",
		"en":                    "/en",
		"en-GB;q=0.8, en;q=0.7": "/en",
	}

	for accept, want := range tests {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Accept-Language", accept)
		w := httptest.NewRecorder()
		langs.Redirect(w, r)
		result := w.Result()
		if result.StatusCode != http.StatusSeeOther {
			t.Fatalf("got status %d, want %d", result.StatusCode, http.StatusSeeOther)
		}
		if got := result.Header.Get("Location"); got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func TestTranslate(t *testing.T) {
	for _, prefix := range langs.Prefixes() {
		http.HandleFunc("/"+prefix, func(w http.ResponseWriter, r *http.Request) {
			_, printer, _ := langs.FromPath(r)
			printer.Fprintf(w, "Hello World")
		})
	}

	tests := map[string]string{
		"de": "Hallo Welt",
		"en": "Hello World",
	}

	for prefix, want := range tests {
		r := httptest.NewRequest(http.MethodGet, "/"+prefix, nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, r)
		got := w.Body.String()

		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
