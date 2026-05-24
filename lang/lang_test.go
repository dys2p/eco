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

func TestTranslate(t *testing.T) {
	for _, l := range langs {
		http.HandleFunc("/"+l.Prefix, func(w http.ResponseWriter, r *http.Request) {
			l, _, _ := langs.FromURL(r.URL)
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
