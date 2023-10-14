package payment

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/dys2p/eco/lang"
)

var cashTmpl = template.Must(template.ParseFS(htmlfiles, "cash.html"))

type cashTmplData struct {
	lang.Lang
	AddressHTML template.HTML
	PurchaseID  string
}

type Cash struct {
	AddressHTML string
}

func (Cash) ID() string {
	return "cash"
}

func (Cash) Name(langstr string) string {
	return lang.Lang(langstr).Tr("Cash")
}

func (cash Cash) PayHTML(purchaseID, paymentKey, langstr string) (template.HTML, error) {
	buf := &bytes.Buffer{}
	err := cashTmpl.Execute(buf, cashTmplData{
		Lang:        lang.Lang(langstr),
		AddressHTML: template.HTML(cash.AddressHTML),
		PurchaseID:  purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (Cash) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (Cash) VerifiesAdult() bool {
	return false
}
