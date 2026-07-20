// Package email implements email submission.
package email

import (
	"bytes"
	"errors"
	"fmt"
	"mime"
	"net/mail"
	"strings"
	"time"

	"github.com/dys2p/eco/id"
)

var ErrInvalidAddress = errors.New("invalid address")

type Emailer interface {
	Send(em Email) error
}

// AddressValid returns true if addr is a well-formed email address, and if it exactly one email address and not a list.
// Use AddressValid to check the email address in your application.
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

type Email struct {
	To      string
	Cc      string
	Subject string
	Body    []byte
}

func (em Email) bytes(from string) (*bytes.Buffer, error) {
	fromDomain, err := getDomain(from)
	if err != nil {
		return nil, err
	}

	if !AddressValid(em.To) {
		return nil, ErrInvalidAddress
	}
	if em.Cc != "" && !AddressValid(em.Cc) {
		return nil, ErrInvalidAddress
	}

	msg := &bytes.Buffer{}
	msg.WriteString("MIME-Version: 1.0" + "\r\n")
	msg.WriteString("Content-Type: text/plain; charset=utf-8" + "\r\n")
	msg.WriteString("Date: " + time.Now().Format("2 Jan 2006 15:04:05 -0700") + "\r\n")
	msg.WriteString("Message-ID: " + newMessageId(fromDomain) + "\r\n") //
	msg.WriteString("From: " + mime.QEncoding.Encode("utf-8", from) + "\r\n")
	msg.WriteString("Subject: " + mime.QEncoding.Encode("utf-8", em.Subject) + "\r\n")
	msg.WriteString("To: " + mime.QEncoding.Encode("utf-8", em.To) + "\r\n")
	if em.Cc != "" {
		msg.WriteString("Cc: " + mime.QEncoding.Encode("utf-8", em.Cc) + "\r\n")
	}
	msg.WriteString("\r\n")
	msg.Write(em.Body)
	return msg, nil
}
