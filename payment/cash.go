package payment

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/dys2p/eco/language"
)

var cashTmpl = template.Must(template.ParseFS(htmlfiles, "cash.html"))

type cashTmplData struct {
	language.Lang
	AddressHTML template.HTML
	PurchaseID  string
}

type Cash struct {
	AddressHTML string
}

func (Cash) ID() string {
	return "cash"
}

func (Cash) Name(r *http.Request) string {
	return language.Get(r).Tr("Cash")
}

func (cash Cash) PayHTML(r *http.Request, purchaseID string) (template.HTML, error) {
	buf := &bytes.Buffer{}
	err := cashTmpl.Execute(buf, cashTmplData{
		Lang:        language.Get(r),
		AddressHTML: template.HTML(cash.AddressHTML),
		PurchaseID:  purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (Cash) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (Cash) VerifiesAdult() bool {
	return false
}
