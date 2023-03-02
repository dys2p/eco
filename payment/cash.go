package payment

import (
	"html/template"
	"net/http"
)

type Cash struct {
	PayHtml string
}

func (Cash) ID() string {
	return "cash"
}

func (Cash) Name() string {
	return "Barzahlung"
}

func (cash Cash) PayHTML(purchaseID string) (template.HTML, error) {
	return template.HTML(cash.PayHtml), nil
}

func (Cash) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (Cash) VerifiesAdult() bool {
	return false
}
