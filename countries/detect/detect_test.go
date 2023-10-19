package detect

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/dys2p/eco/countries"
)

func TestAcceptLanguage(t *testing.T) {
	tests := map[string][]countries.Country{ // key: accept-language
		"bg":                        []countries.Country{countries.BG},
		"de":                        []countries.Country{countries.AT, countries.BE, countries.DE, countries.IT, countries.LU, countries.NonEU},
		"de-DE":                     []countries.Country{countries.AT, countries.BE, countries.DE, countries.IT, countries.LU, countries.NonEU},
		"de, en-gb;q=0.8, en;q=0.7": []countries.Country{countries.AT, countries.BE, countries.DE, countries.IT, countries.LU, countries.NonEU},
		"en-US,en;q=0.5":            nil,
		"fr":                        []countries.Country{countries.BE, countries.FR, countries.IT, countries.LU, countries.NonEU},
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
		want         []countries.Country
	}{
		// de-fra-ovpn-001.relays.mullvad.net
		{"185.213.155.66", "", []countries.Country{countries.DE}},
		{"185.213.155.66:8080", "", []countries.Country{countries.DE}},
		{"2a03:1b20:6:f011::1f", "", []countries.Country{countries.DE}},
		{"127.0.0.1", "185.213.155.66", []countries.Country{countries.DE}},
		{"127.0.0.1", "185.213.155.66:8080", []countries.Country{countries.DE}},
		{"127.0.0.1", "2a03:1b20:6:f011::1f", []countries.Country{countries.DE}},

		// gr-ath-ovpn-101.relays.mullvad.net (must return "GR", not "EL")
		{"149.102.246.28", "", []countries.Country{countries.GR}},
		{"2a02:6ea0:f501:4::1f", "", []countries.Country{countries.GR}},
		{"192.168.1.1", "149.102.246.28", []countries.Country{countries.GR}},
		{"192.168.1.1", "2a02:6ea0:f501:4::1f", []countries.Country{countries.GR}},

		// undefined
		{"1.1.1.1", "", nil},
		{"127.0.0.1", "1.1.1.1", nil},
		{"127.0.0.1", "", nil},
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
