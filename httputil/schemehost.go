package httputil

import (
	"fmt"
	"net/http"
	"strings"
)

// SchemeHost returns the scheme plus host (including port, if available). The scheme is https, except when the host looks like localhost or a Tor hidden service.
//
// If you use nginx as a reverse proxy, make sure you have set "proxy_set_header Host $host;" besides proxy_pass in your configuration.
func SchemeHost(r *http.Request) string {
	var proto = "https"
	if strings.HasPrefix(r.Host, "127.0.") || strings.HasPrefix(r.Host, "[::1]") || strings.HasSuffix(r.Host, ".onion") || strings.Contains(r.Host, ".onion:") {
		proto = "http"
	}
	return fmt.Sprintf("%s://%s", proto, r.Host)
}
