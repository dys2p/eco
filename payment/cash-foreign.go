package payment

import (
	"bytes"
	"encoding/json"
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

func (cash CashForeign) Handler() http.Handler {
	var mux = http.NewServeMux()
	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cash.History.Synced)
	})
	return mux
}

func (CashForeign) ID() string {
	return "cash-foreign"
}

func (CashForeign) Name(l lang.Lang) string {
	return l.Tr("Cash in Foreign Currency")
}

func (cash CashForeign) PayHTML(purchaseID, paymentKey string, l lang.Lang) (template.HTML, error) {
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
	currencyOptions, err := cash.History.Options(date, euros)
	if err != nil {
		log.Printf("error getting currency options: %v", err)
		return template.HTML("Error getting exchange rates. Please try again in a minute."), nil
	}

	buf := &bytes.Buffer{}
	err = cashForeignTmpl.Execute(buf, cashForeignTmplData{
		Lang:            l,
		AddressHTML:     template.HTML(cash.AddressHTML),
		CurrencyOptions: currencyOptions,
		PurchaseID:      purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (CashForeign) VerifiesAdult() bool {
	return false
}
