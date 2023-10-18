package detect

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestAcceptLanguage(t *testing.T) {
	tests := map[string][]string{ // key: accept-language
		"bg":                        []string{"BG"},
		"de":                        []string{"AT", "DE", "non-EU"},
		"de-DE":                     []string{"AT", "DE", "non-EU"},
		"de, en-gb;q=0.8, en;q=0.7": []string{"AT", "DE", "non-EU"},
		"en-US,en;q=0.5":            []string{"non-EU", "AT", "BE", "BG", "CY", "CZ", "DE", "DK", "EE", "ES", "FI", "FR", "GR", "HR", "HU", "IE", "IT", "LT", "LU", "LV", "MT", "NL", "PL", "PT", "RO", "SE", "SI", "SK"},
		"fr":                        []string{"FR", "non-EU"},
	}

	for accept, want := range tests {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Language", accept)
		got, err := acceptLanguage(req)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(got, want) {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func TestIPAddress(t *testing.T) {
	tests := []struct {
		remoteAddr   string
		forwardedFor string
		want         []string
	}{
		// de-fra-ovpn-001.relays.mullvad.net
		{"185.213.155.66", "", []string{"DE"}},
		{"185.213.155.66:8080", "", []string{"DE"}},
		{"2a03:1b20:6:f011::1f", "", []string{"DE"}},
		{"127.0.0.1", "185.213.155.66", []string{"DE"}},
		{"127.0.0.1", "185.213.155.66:8080", []string{"DE"}},
		{"127.0.0.1", "2a03:1b20:6:f011::1f", []string{"DE"}},

		// undefined
		{"1.1.1.1", "", nil},
		{"127.0.0.1", "1.1.1.1", nil},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = test.remoteAddr
		req.Header.Set("X-Forwarded-For", test.forwardedFor)
		got, err := ipAddress(req)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(got, test.want) {
			t.Fatalf("got %s, want %s", got, test.want)
		}
	}
}
