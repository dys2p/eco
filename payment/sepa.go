package payment

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/dys2p/eco/language"
)

var sepaTmpl = template.Must(template.ParseFS(htmlfiles, "sepa.html"))

type sepaTmplData struct {
	language.Lang
	Account SEPAAccount
	Amount  float64
	Purpose string
}

type SEPAAccount struct {
	Holder   string
	IBAN     string
	BIC      string
	BankName string
}

type SEPA struct {
	Account   SEPAAccount
	Purchases PurchaseRepo
}

func (SEPA) ID() string {
	return "sepa"
}

func (SEPA) Name(r *http.Request) string {
	return language.Get(r).Tr("SEPA Bank Transfer")
}

func (sepa SEPA) PayHTML(r *http.Request, purchaseID string) (template.HTML, error) {
	eurocents, err := sepa.Purchases.PurchaseSumCents(purchaseID)
	if err != nil {
		log.Printf("error getting purchase sum from database: %v", err)
		return template.HTML("Error getting purchase information from database"), nil
	}

	buf := &bytes.Buffer{}
	err = sepaTmpl.Execute(buf, sepaTmplData{
		Lang:    language.Get(r),
		Account: sepa.Account,
		Amount:  float64(eurocents) / 100.0,
		Purpose: purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (SEPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (SEPA) VerifiesAdult() bool {
	return false
}
