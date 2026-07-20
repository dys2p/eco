package email

import "log"

type DummyMailer struct{}

func (DummyMailer) Send(em Email) error {
	// same validity checks as in Email.Bytes
	if !AddressValid(em.To) {
		return ErrInvalidAddress
	}
	if em.Cc != "" && !AddressValid(em.Cc) {
		return ErrInvalidAddress
	}

	log.Println("------ dummy mailer ------")
	log.Printf("To: %s", em.To)
	if em.Cc != "" {
		log.Printf("Cc: %s", em.Cc)
	}
	log.Printf("Subject: %s", em.Subject)
	log.Printf("%s", em.Body)
	return nil
}
