package payment

import (
	"fmt"
	"html/template"
	"net/http"
)

type Wire struct {
	PayHtml string // must contain %s as a placeholder for the order id
}

func (Wire) ID() string {
	return "wire"
}

func (Wire) Name() string {
	return "Banküberweisung"
}

func (wire Wire) PayHTML(purchaseID string) (template.HTML, error) {
	return template.HTML(fmt.Sprintf(wire.PayHtml, purchaseID)), nil
}

func (Wire) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (Wire) VerifiesAdult() bool {
	return false
}
