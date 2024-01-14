package ntfysh

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: time.Minute,
}

func Publish(addr, title, message string) error {
	req, err := http.NewRequest(http.MethodPost, addr, strings.NewReader(message))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Title", title)
	_, err = client.Do(req)
	return err
}

func ValidateAddress(addr string) string {
	addr = strings.TrimSpace(addr)
	if len(addr) == 0 {
		return ""
	}

	if validTopic(addr) {
		return "https://ntfy.sh/" + addr
	}

	u, err := url.Parse(addr)
	if err != nil {
		return ""
	}
	if u.Host == "" {
		u, err = url.Parse("https://" + addr)
		if err != nil {
			return ""
		}
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	if u.Opaque != "" {
		return ""
	}
	if u.User != nil {
		return ""
	}
	if u.Host == "" {
		return ""
	}
	if topic := strings.TrimPrefix(u.Path, "/"); !validTopic(topic) {
		return ""
	}
	if u.RawQuery != "" {
		return ""
	}
	if u.Fragment != "" {
		return ""
	}
	return u.String()
}

func validTopic(topic string) bool {
	if len(topic) == 0 {
		return false
	}
	for i := range topic {
		b := topic[i]
		valid := ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') || b == '-' || b == '_' // not found in docs, but apparently it is [a-zA-Z0-9-_]
		if !valid {
			return false
		}
	}
	return true
}
