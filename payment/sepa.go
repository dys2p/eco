package payment

import (
	"fmt"
	"html/template"
	"net/http"
)

type SEPA struct {
	PayHtml string // must contain %s as a placeholder for the order id
}

func (SEPA) ID() string {
	return "sepa"
}

func (SEPA) Name() string {
	return "SEPA-Banküberweisung"
}

func (sepa SEPA) PayHTML(purchaseID string) (template.HTML, error) {
	return template.HTML(fmt.Sprintf(sepa.PayHtml, purchaseID)), nil
}

func (SEPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (SEPA) VerifiesAdult() bool {
	return false
}
