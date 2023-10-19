package detect

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/abh/geoip"
	"github.com/dys2p/eco/countries"
)

// geoip requires different instances for IPv4 and IPv6, else Open returns the error "Invalid database type GeoIP Country V6 Edition, expected GeoIP Country Edition"
var geoIPv4 *geoip.GeoIP
var geoIPv6 *geoip.GeoIP

func init() {
	var err error
	geoIPv4, err = geoip.Open("/usr/share/GeoIP/GeoIP.dat")
	if err != nil {
		panic("Could not open GeoIP database\n")
	}
	geoIPv6, err = geoip.Open("/usr/share/GeoIP/GeoIPv6.dat")
	if err != nil {
		panic("Could not open GeoIP database\n")
	}
}

// ipAddress returns zero or one country. The country can be outside of the European Union.
func ipAddress(r *http.Request) ([]countries.Country, error) {
	// TOR
	if strings.HasSuffix(r.Host, ".onion") || strings.Contains(r.Host, ".onion:") {
		return nil, nil
	}

	// first X-Forwarded-For header value overrides http remote address
	clientAddr := r.RemoteAddr // RemoteAddr is "IP:port"
	if forwardedFor := strings.FieldsFunc(r.Header.Get("X-Forwarded-For"), func(r rune) bool { return r == ',' || r == ' ' }); len(forwardedFor) > 0 {
		clientAddr = forwardedFor[0]
	}

	// parse addresses with and without a port number
	var ip net.IP
	if tcpAddr, err := net.ResolveTCPAddr("tcp", clientAddr); err == nil { // ResolveTCPAddr requires port number
		ip = tcpAddr.IP
	} else {
		if ipAddr, err := net.ResolveIPAddr("ip", clientAddr); err == nil { // ResolveIPAddr is without port number
			ip = ipAddr.IP
		} else {
			return nil, fmt.Errorf("resolving tcp address: %w", err)
		}
	}

	var country string
	if ipv4 := ip.To4(); ipv4 != nil {
		country, _ = geoIPv4.GetCountry(ipv4.String())
	} else {
		country, _ = geoIPv6.GetCountry_v6(ip.String()) // "be sure to load a database with IPv6 data to get any results"
	}
	if country != "" {
		return []countries.Country{countries.Country(country)}, nil
	} else {
		return nil, nil
	}
}
