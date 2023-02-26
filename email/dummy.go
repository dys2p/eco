package email

import "log"

type DummyMailer struct{}

func (DummyMailer) Send(to string, subject string, body []byte) error {
	log.Println("------ dummy mailer ------")
	log.Printf("to: %s", to)
	log.Printf("subject: %s", subject)
	log.Printf("body: %s", body)
	return nil
}
