package lang

import (
	"net/http"
	"net/http/httptest"
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

	langs = MakeLanguages(nil, "en", "de")
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
		langs.RedirectHandler()(w, r)
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
	for _, l := range langs {
		http.HandleFunc("/"+l.Prefix, func(w http.ResponseWriter, r *http.Request) {
			l, _ := langs.FromPath(r.URL.Path)
			l.Printer.Fprintf(w, "Hello World")
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
