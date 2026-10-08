package payment

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/dys2p/eco/lang"
)

var cashTmpl = template.Must(template.ParseFS(htmlfiles, "cash.html"))

type cashTmplData struct {
	lang.Lang
	AddressHTML template.HTML
	Amount      float64
	PurchaseID  string
}

type Cash struct {
	AddressHTML string
	Purchases   PurchaseRepo
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

func (cash Cash) PayHTML(purchaseID, paymentKey, redirectURL string, l lang.Lang) (template.HTML, error) {
	eurocents, err := cash.Purchases.PurchaseDueCents(purchaseID, paymentKey)
	if err != nil {
		log.Printf("error getting purchase due from database: %v", err)
		return template.HTML("Error getting purchase from database"), nil
	}

	buf := &bytes.Buffer{}
	err = cashTmpl.Execute(buf, cashTmplData{
		Lang:        l,
		AddressHTML: template.HTML(cash.AddressHTML),
		Amount:      float64(eurocents) / 100.0,
		PurchaseID:  purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (Cash) VerifiesAdult() bool {
	return false
}
