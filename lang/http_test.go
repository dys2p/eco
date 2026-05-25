package lang

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAcceptLanguage(t *testing.T) {
	var mux = http.NewServeMux()
	HandleFunc(mux, langs, "/", func(http.ResponseWriter, *http.Request) {})

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
		mux.ServeHTTP(w, r)
		result := w.Result()
		if result.StatusCode != http.StatusSeeOther {
			t.Fatalf("got status %d, want %d", result.StatusCode, http.StatusSeeOther)
		}
		if got := result.Header.Get("Location"); got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func respond(s string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(s))
	}
}

func TestInsert(t *testing.T) {
	tests := map[string]string{
		"/":                       `/en/`,
		"/{$}":                    `/en/{$}`,
		"/foo":                    `/en/foo`,
		"/foo/bar":                `/en/foo/bar`,
		"GET /":                   "GET /en/",
		"GET /{$}":                "GET /en/{$}",
		"GET /foo":                `GET /en/foo`,
		"GET /foo/bar":            `GET /en/foo/bar`,
		"example.com/":            `example.com/en/`,
		"example.com/{$}":         `example.com/en/{$}`,
		"example.com/foo":         `example.com/en/foo`,
		"example.com/foo/bar":     `example.com/en/foo/bar`,
		"GET example.com/":        "GET example.com/en/",
		"GET example.com/{$}":     "GET example.com/en/{$}",
		"GET example.com/foo":     `GET example.com/en/foo`,
		"GET example.com/foo/bar": `GET example.com/en/foo/bar`,
	}

	for pattern, want := range tests {
		got := insert(pattern, "en")
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}
