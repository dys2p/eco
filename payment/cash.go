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

func (Cash) Handler() http.Handler {
	return http.NewServeMux()
}

func (Cash) ID() string {
	return "cash"
}

func (Cash) Name(l lang.Lang) string {
	return l.Tr("Cash")
}

func (cash Cash) PayHTML(purchaseID, paymentKey string, l lang.Lang) (template.HTML, error) {
	buf := &bytes.Buffer{}
	err := cashTmpl.Execute(buf, cashTmplData{
		Lang:        l,
		AddressHTML: template.HTML(cash.AddressHTML),
		PurchaseID:  purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (Cash) VerifiesAdult() bool {
	return false
}
