package payment

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/dys2p/btcpay"
	"github.com/dys2p/eco/lang"
)

func init() {
	log.Println(`Don't forget to set up the BTCPay webhook for your store: URL: "/payment/btcpay/webhook", events: "An invoice is processing" and "An invoice has been settled"`)
}

var btcpayTmpl = template.Must(template.ParseFS(htmlfiles, "btcpay.html"))

type btcpayTmplData struct {
	lang.Lang
	DefaultLanguage string
	Reference       string
}

type createdInvoice struct {
	ID   string
	Time int64
}

var lastInvoice = make(map[string]createdInvoice) // key: purchase ID

type BTCPay struct {
	ExpirationMinutes int
	RedirectPath      string
	Store             btcpay.Store
	Purchases         PurchaseRepo
}

func (BTCPay) ID() string {
	return "btcpay"
}

func (BTCPay) Name(l lang.Lang) string {
	return l.Tr("Monero or Bitcoin")
}

func (b BTCPay) PayHTML(purchaseID, paymentKey string, l lang.Lang) (template.HTML, error) {
	buf := &bytes.Buffer{}
	err := btcpayTmpl.Execute(buf, btcpayTmplData{
		Lang:            l,
		DefaultLanguage: l.Prefix,
		Reference:       purchaseID + ":" + paymentKey,
	})
	return template.HTML(buf.String()), err
}

func (b BTCPay) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch path.Base(r.URL.Path) {
	case "create-invoice":
		if err := b.createInvoice(w, r); err != nil {
			log.Printf("error creating btcpay invoice: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "webhook":
		if err := b.webhook(w, r); err != nil {
			log.Printf("error processing btcpay webhook: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (b BTCPay) createInvoice(w http.ResponseWriter, r *http.Request) error {
	defaultLanguage := r.PostFormValue("default-language")
	purchaseID, paymentKey, _ := strings.Cut(r.PostFormValue("reference"), ":")

	// redirect to existing invoice if it is younger than 15 minutes
	if last, ok := lastInvoice[purchaseID+":"+paymentKey]; ok && time.Now().Unix()-last.Time < 15*60 {
		http.Redirect(w, r, b.checkoutLink(r, last.ID), http.StatusSeeOther)
		return nil
	}

	sumCents, err := b.Purchases.PurchaseSumCents(purchaseID, paymentKey)
	if err != nil {
		return fmt.Errorf("getting sum: %w", err)
	}

	invoiceRequest := &btcpay.InvoiceRequest{
		Amount:   float64(sumCents) / 100.0,
		Currency: "EUR",
	}
	invoiceRequest.ExpirationMinutes = b.expirationMinutes()
	invoiceRequest.DefaultLanguage = defaultLanguage
	invoiceRequest.OrderID = purchaseID + ":" + paymentKey // reference
	invoiceRequest.RedirectURL = absHost(r) + path.Join("/", b.RedirectPath)
	invoice, err := b.Store.CreateInvoice(invoiceRequest)
	if err != nil {
		return fmt.Errorf("querying store: %w", err)
	}

	lastInvoice[purchaseID+":"+paymentKey] = createdInvoice{
		ID:   invoice.ID,
		Time: time.Now().Unix(),
	}

	http.Redirect(w, r, b.checkoutLink(r, invoice.ID), http.StatusSeeOther)
	return nil
}

func (b BTCPay) checkoutLink(r *http.Request, invoiceID string) string {
	// ignore invoice.CheckoutLink in favor of the onion option
	link := b.Store.InvoiceCheckoutLink(invoiceID)
	if strings.HasSuffix(r.Host, ".onion") || strings.Contains(r.Host, ".onion:") {
		link = b.Store.InvoiceCheckoutLinkPreferOnion(invoiceID)
	}
	return link
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
	purchaseID, paymentKey, _ := strings.Cut(event.InvoiceMetadata.OrderID, ":")

	switch event.Type {
	case btcpay.EventInvoiceProcessing:
		if err := b.Purchases.SetPurchaseProcessing(purchaseID, paymentKey); err != nil {
			return fmt.Errorf("setting purchase %s processing: %w", purchaseID, err)
		}
		return nil
	case btcpay.EventInvoiceSettled:
		if err := b.Purchases.SetPurchasePaid(purchaseID, paymentKey); err != nil {
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

// absHost returns the scheme and host part of an HTTP request. It uses a heuristic for the scheme.
//
// If you use nginx as a reverse proxy, make sure you have set "proxy_set_header Host $host;" besides proxy_pass in your configuration.
func absHost(r *http.Request) string {
	var proto = "https"
	if strings.HasPrefix(r.Host, "127.0.") || strings.HasPrefix(r.Host, "[::1]") || strings.HasSuffix(r.Host, ".onion") { // if running locally or through TOR
		proto = "http"
	}
	return fmt.Sprintf("%s://%s", proto, r.Host)
}
