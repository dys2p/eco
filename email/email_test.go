package email

import (
	"regexp"
	"testing"
	"time"
)

var emailers = []Emailer{
	DummyMailer{},
	Sendmail{
		From: "mail@example.com",
	},
	SMTP{
		From:     "mail@example.com",
		Username: "mail@example.com",
		Password: "change-me",
		Host:     "example.com",
	},
}

func TestInvalidAddresses(t *testing.T) {
	var addrs = []string{
		"",
		"test",
		"test@",
		"example.com",
		"@example.com",
		"test1@example.com test2@example.com",
		"test1@example.com,test2@example.com",
		"test1@example.com;test2@example.com",
		"test+1@example.com;test+2@example.com",
	}
	for _, emailer := range emailers {
		for _, addr := range addrs {
			if err := emailer.Send(addr, "Subject", nil); err != ErrInvalidAddress {
				t.Fatalf("got %v, want %v", err, ErrInvalidAddress)
			}
		}
	}
}

var messageID = regexp.MustCompile("[a-zA-Z0-9]{16}")

func TestMakeEmail(t *testing.T) {
	buf, _ := MakeEmail("alice@example.com", "bob@example.com", "Hello World", []byte("This is an example email."))
	got := buf.String()
	got = messageID.ReplaceAllString(got, "0123456789ABCDEF")

	var want string
	want += "MIME-Version: 1.0" + "\r\n"
	want += "Content-Type: text/plain; charset=utf-8" + "\r\n"
	want += "Date: " + time.Now().Format("02 Jan 2006 15:04:05 -0700") + "\r\n"
	want += "Message-ID: <0123456789ABCDEF@example.com>" + "\r\n"
	want += "From: alice@example.com" + "\r\n"
	want += "Subject: Hello World" + "\r\n"
	want += "To: bob@example.com" + "\r\n"
	want += "\r\n"
	want += "This is an example email."

	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
