package email

import "log"

type DummyMailer struct{}

func (DummyMailer) Send(to, cc, subject string, body []byte) error {
	if !AddressValid(to) {
		return ErrInvalidAddress
	}
	log.Println("------ dummy mailer ------")
	log.Printf("to: %s", to)
	if cc != "" {
		log.Printf("cc: %s", cc)
	}
	log.Printf("subject: %s", subject)
	log.Printf("body: %s", body)
	return nil
}
