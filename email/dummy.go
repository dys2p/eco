package email

import "log"

type DummyMailer struct{}

func (DummyMailer) Send(to string, subject string, body []byte) error {
	log.Println("------ dummy mailer ------")
	log.Println("to: %s", to)
	log.Println("subject: %s", subject)
	log.Println("body: %s", body)
	return nil
}
