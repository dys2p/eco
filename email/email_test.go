package email

import "testing"

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
