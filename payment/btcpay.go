package payment

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dys2p/btcpay"
	"github.com/dys2p/eco/language"
)

func init() {
	log.Println(`Don't forget to set up the BTCPay webhook for your store: URL: "/payment/btcpay/webhook", events: "An invoice is processing" and "An invoice has been settled"`)
}

var btcpayTmpl = template.Must(template.ParseFS(htmlfiles, "btcpay.html"))

type btcpayTmplData struct {
	language.Lang
	PurchaseID string
}

type createdInvoice struct {
	CheckoutLink string
	Time         int64
}

var lastInvoice = make(map[string]createdInvoice) // key: purchase ID

type BTCPay struct {
	ExpirationMinutes int
	RedirectURL       string
	Store             btcpay.Store
	Purchases         PurchaseRepo
}

func (BTCPay) ID() string {
	return "btcpay"
}

func (BTCPay) Name(r *http.Request) string {
	return language.Get(r).Tr("Monero oder Bitcoin")
}

func (b BTCPay) PayHTML(r *http.Request, purchaseID string) (template.HTML, error) {
	buf := &bytes.Buffer{}
	err := btcpayTmpl.Execute(buf, btcpayTmplData{
		Lang:       language.Get(r),
		PurchaseID: purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (b BTCPay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	switch r.URL.Path {
	case "/payment/btcpay/create-invoice":
		if err := b.createInvoice(w, r); err != nil {
			log.Printf("error creating btcpay invoice: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/payment/btcpay/webhook":
		if err := b.webhook(w, r); err != nil {
			log.Printf("error processing btcpay webhook: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (b BTCPay) createInvoice(w http.ResponseWriter, r *http.Request) error {
	purchaseID := r.PostFormValue("purchase-id")

	// redirect to existing invoice if it is younger than 15 minutes
	if last, ok := lastInvoice[purchaseID]; ok && time.Now().Unix()-last.Time < 15*60 {
		http.Redirect(w, r, last.CheckoutLink, http.StatusSeeOther)
		return nil
	}

	sumCents, err := b.Purchases.PurchaseSumCents(purchaseID)
	if err != nil {
		return fmt.Errorf("getting sum: %w", err)
	}

	defaultLanguage := strings.TrimPrefix(language.Get(r).Tr("btcpay:de-DE"), "btcpay:") // see https://github.com/btcpayserver/btcpayserver/tree/master/BTCPayServer/wwwroot/locales

	invoiceRequest := &btcpay.InvoiceRequest{
		Amount:   float64(sumCents) / 100.0,
		Currency: "EUR",
	}
	invoiceRequest.ExpirationMinutes = b.expirationMinutes()
	invoiceRequest.DefaultLanguage = defaultLanguage
	invoiceRequest.OrderID = purchaseID
	invoiceRequest.RedirectURL = b.RedirectURL
	invoice, err := b.Store.CreateInvoice(invoiceRequest)
	if err != nil {
		return fmt.Errorf("querying store: %w", err)
	}

	lastInvoice[purchaseID] = createdInvoice{
		CheckoutLink: invoice.CheckoutLink,
		Time:         time.Now().Unix(),
	}

	http.Redirect(w, r, invoice.CheckoutLink, http.StatusSeeOther)
	return nil
}

func (b BTCPay) expirationMinutes() int {
	if b.ExpirationMinutes == 0 {
		return 60 // default
	}
	if b.ExpirationMinutes < 30 {
		return 30
	}
	if b.ExpirationMinutes > 1440 {
		return 1440
	}
	return b.ExpirationMinutes
}

func (b BTCPay) webhook(w http.ResponseWriter, r *http.Request) error {
	event, err := b.Store.ProcessWebhook(r)
	if err != nil {
		return fmt.Errorf("getting event: %w", err)
	}
	invoice, err := b.Store.GetInvoice(event.InvoiceID)
	if err != nil {
		return fmt.Errorf("getting invoice %s: %w", event.InvoiceID, err)
	}
	purchaseID := invoice.InvoiceMetadata.OrderID

	switch event.Type {
	case btcpay.EventInvoiceProcessing:
		if err := b.Purchases.SetPurchaseProcessing(purchaseID); err != nil {
			return fmt.Errorf("setting purchase %s processing: %w", purchaseID, err)
		}
		return nil
	case btcpay.EventInvoiceSettled:
		if err := b.Purchases.SetPurchasePaid(purchaseID); err != nil {
			return fmt.Errorf("setting purchase %s paid: %w", purchaseID, err)
		}
		return nil
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

func (BTCPay) VerifiesAdult() bool {
	return false
}
