// Package email implements email submission.
package email

import (
	"bytes"
	"fmt"
	"mime"
	"net/mail"
	"strings"
	"time"

	"github.com/dys2p/eco/id"
)

type Emailer interface {
	Send(to string, subject string, body []byte) error
}

func AddressValid(addr string) bool {
	_, err := mail.ParseAddress(addr)
	return err == nil
}

func getDomain(address string) (string, error) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", err
	}
	if at := strings.LastIndex(addr.Address, "@"); at >= 0 {
		return addr.Address[at+1:], nil
	} else {
		return "", fmt.Errorf("addr contains no @: %s", addr)
	}
}

// newMessageId creates a new RFC5322 compliant Message-Id with the given domain as "id-right".
func newMessageId(domain string) string {
	idLeft := id.New(16, id.AlphanumCaseSensitiveDigits) // RFC5322 "atext"
	// RFC 5322: The message identifier (msg-id) syntax is a limited version of the addr-spec construct enclosed in the angle bracket characters, "<" and ">".
	// mail.Address.String() encloses the result in angle brackets.
	return (&mail.Address{Address: idLeft + "@" + domain}).String()
}

func MakeEmail(from, to, subject string, body []byte) (*bytes.Buffer, error) {
	fromDomain, err := getDomain(from)
	if err != nil {
		return nil, err
	}

	msg := &bytes.Buffer{}
	msg.WriteString("MIME-Version: 1.0" + "\r\n")
	msg.WriteString("Content-Type: text/plain; charset=utf-8" + "\r\n")
	msg.WriteString("Date: " + time.Now().Format("2 Jan 2006 15:04:05 -0700") + "\r\n")
	msg.WriteString("Message-ID: " + newMessageId(fromDomain) + "\r\n") //
	msg.WriteString("From: " + mime.QEncoding.Encode("utf-8", from) + "\r\n")
	msg.WriteString("Subject: " + mime.QEncoding.Encode("utf-8", subject) + "\r\n")
	msg.WriteString("To: " + mime.QEncoding.Encode("utf-8", to) + "\r\n")
	msg.WriteString("\r\n")
	msg.Write(body)
	return msg, nil
}
