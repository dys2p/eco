package payment

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/dys2p/eco/language"
)

var sepaTmpl = template.Must(template.ParseFS(htmlfiles, "sepa.html"))

type sepaTmplData struct {
	language.Lang
	Account SEPAAccount
	Purpose string
}

type SEPAAccount struct {
	Holder   string
	IBAN     string
	BIC      string
	BankName string
}

type SEPA struct {
	Account SEPAAccount
}

func (SEPA) ID() string {
	return "sepa"
}

func (SEPA) Name(r *http.Request) string {
	return language.Get(r).Tr("SEPA-Banküberweisung")
}

func (sepa SEPA) PayHTML(r *http.Request, purchaseID string) (template.HTML, error) {
	buf := &bytes.Buffer{}
	err := sepaTmpl.Execute(buf, sepaTmplData{
		Lang:    language.Get(r),
		Account: sepa.Account,
		Purpose: purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (SEPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (SEPA) VerifiesAdult() bool {
	return false
}
