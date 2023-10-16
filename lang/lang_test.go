package lang

import (
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestRedirect(t *testing.T) {
	var langs = Languages([]string{"en", "de"})

	tests := map[string]string{
		"de":                    "/de",
		"de-CH":                 "/de",
		"de-DE":                 "/de",
		"":                      "/en",
		"en":                    "/en",
		"en-GB;q=0.8, en;q=0.7": "/en",
	}

	for accept, want := range tests {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Language", accept)
		rec := httptest.NewRecorder()
		langs.Redirect(rec, req)
		result := rec.Result()
		if result.StatusCode != http.StatusSeeOther {
			t.Fatalf("got status %d, want %d", result.StatusCode, http.StatusSeeOther)
		}
		if got := result.Header.Get("Location"); got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
