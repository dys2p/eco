package lang

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestHandle(t *testing.T) {
	var mux = http.NewServeMux()

	HandleFunc(mux, langs, "GET /{$}", respond("GET /"))
	HandleFunc(mux, langs, "GET /foo", respond("GET /foo"))
	HandleFunc(mux, langs, "POST /bar", respond("POST /bar"))
	HandleFunc(mux, langs, "example.net/", respond("example.net/"))
	HandleFunc(mux, langs, "example.net/baz", respond("example.net/baz"))

	tests := map[string]string{
		// httptest uses example.com by default
		"/":                          `<a href="/en">See Other</a>.`,
		"/de/":                       "GET /",
		"/en/":                       "GET /",
		"/foo":                       `<a href="/en/foo">See Other</a>.`,
		"/de/foo":                    `GET /foo`,
		"/en/foo":                    `GET /foo`,
		"/bar":                       `Method Not Allowed`,
		"/de/bar":                    `Method Not Allowed`,
		"/en/bar":                    `Method Not Allowed`,
		"https://example.net/":       `<a href="https://example.net/en">See Other</a>.`,
		"https://example.net/de/":    `example.net/`,
		"https://example.net/en/":    `example.net/`,
		"https://example.net/baz":    `<a href="https://example.net/en/baz">See Other</a>.`,
		"https://example.net/de/baz": `example.net/baz`,
		"https://example.net/en/baz": `example.net/baz`,
	}

	for reqpath, want := range tests {
		r := httptest.NewRequest(http.MethodGet, reqpath, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		gotBytes, _ := io.ReadAll(w.Result().Body)
		got := strings.TrimSpace(string(gotBytes))
		if got != want {
			t.Fatalf("[%s] got %s, want %s", reqpath, got, want)
		}
	}
}
