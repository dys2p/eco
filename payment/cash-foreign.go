package payment

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/dys2p/eco/lang"
	"github.com/dys2p/eco/payment/rates"
)

var cashForeignTmpl = template.Must(template.ParseFS(htmlfiles, "cash-foreign.html"))

type cashForeignTmplData struct {
	lang.Lang
	AddressHTML     template.HTML
	CurrencyOptions []rates.Option
	PurchaseID      string
}

type CashForeign struct {
	AddressHTML string
	Purchases   PurchaseRepo
	History     *rates.History
}

func (CashForeign) ID() string {
	return "cash-foreign"
}

func (CashForeign) Name(langstr string) string {
	return lang.Lang(langstr).Tr("Cash in Foreign Currency")
}

func (cash CashForeign) PayHTML(purchaseID, paymentKey, langstr string) (template.HTML, error) {
	date, err := cash.Purchases.PurchaseCreationDate(purchaseID, paymentKey)
	if err != nil {
		log.Printf("error getting purchase creation date from database: %v", err)
		return template.HTML("Error getting purchase information from database"), nil
	}
	eurocents, err := cash.Purchases.PurchaseSumCents(purchaseID, paymentKey)
	if err != nil {
		log.Printf("error getting purchase sum from database: %v", err)
		return template.HTML("Error getting purchase information from database"), nil
	}
	euros := float64(eurocents) / 100.0
	currencyOptions, err := cash.History.Get(date, euros)
	if err != nil {
		log.Printf("error getting currency options: %v", err)
		return template.HTML("Error getting exchange rates. Please try again in a minute."), nil
	}

	buf := &bytes.Buffer{}
	err = cashForeignTmpl.Execute(buf, cashForeignTmplData{
		Lang:            lang.Lang(langstr),
		AddressHTML:     template.HTML(cash.AddressHTML),
		CurrencyOptions: currencyOptions,
		PurchaseID:      purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (CashForeign) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (CashForeign) VerifiesAdult() bool {
	return false
}
