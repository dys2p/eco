package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dys2p/eco/httputil"
	"github.com/dys2p/eco/lang"
	"github.com/dys2p/go-btcpay"
)

func init() {
	log.Println(`Don't forget to set up the BTCPay webhook for your store: URL: "/payment/btcpay/webhook", events: "Invoice - Received Payment", "Invoice - Is Settled", "Invoice - Payment Settled"`)
}

var btcpayTmpl = template.Must(template.ParseFS(htmlfiles, "btcpay.html"))

type btcpayTmplData struct {
	lang.Lang
	DefaultLanguage string
	Reference       string
	Status          bool
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

	ErrCreateInvoice func(err error, msg string) http.Handler
	ErrWebhook       func(err error) http.Handler
	GetStatus        func() []btcpay.StatusItem
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
		Status:          b.GetStatus != nil,
	})
	return template.HTML(buf.String()), err
}

func (b BTCPay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch path.Base(r.URL.Path) {
	case "create-invoice":
		httputil.HandlerFunc(b.createInvoice).ServeHTTP(w, r)
	case "status":
		b.status(w, r)
	case "webhook":
		httputil.HandlerFunc(b.webhook).ServeHTTP(w, r)
	}
}

func (b BTCPay) createInvoice(w http.ResponseWriter, r *http.Request) http.Handler {
	if b.ErrCreateInvoice == nil {
		b.ErrCreateInvoice = func(err error, msg string) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf("error creating btcpay invoice: %v", err)
				w.Write([]byte(msg))
			})
		}
	}

	defaultLanguage := r.PostFormValue("default-language")
	purchaseID, paymentKey, _ := strings.Cut(r.PostFormValue("reference"), ":")

	// redirect to existing invoice if it is younger than 15 minutes
	if last, ok := lastInvoice[purchaseID+":"+paymentKey]; ok && time.Now().Unix()-last.Time < 15*60 {
		return http.RedirectHandler(b.Store.InvoiceCheckoutLink(last.ID, strings.HasSuffix(r.Host, ".onion") || strings.Contains(r.Host, ".onion:")), http.StatusSeeOther)
	}

	sumCents, err := b.Purchases.PurchaseSumCents(purchaseID, paymentKey)
	if err != nil {
		return b.ErrCreateInvoice(err, "Error getting purchase sum. We are already working on it.")
	}

	invoiceRequest := &btcpay.InvoiceRequest{
		Amount:   float64(sumCents) / 100.0,
		Currency: "EUR",
	}
	invoiceRequest.ExpirationMinutes = max(30, min(1440, b.ExpirationMinutes))
	invoiceRequest.DefaultLanguage = defaultLanguage
	invoiceRequest.OrderID = purchaseID + ":" + paymentKey // reference
	invoiceRequest.RedirectURL = absHost(r) + path.Join("/", b.RedirectPath)
	invoice, err := b.Store.CreateInvoice(invoiceRequest)
	if err != nil {
		return b.ErrCreateInvoice(err, "Error creating BTCPay invoice. We are already working on it.")
	}

	lastInvoice[purchaseID+":"+paymentKey] = createdInvoice{
		ID:   invoice.ID,
		Time: time.Now().Unix(),
	}

	return http.RedirectHandler(b.Store.InvoiceCheckoutLink(invoice.ID, strings.HasSuffix(r.Host, ".onion") || strings.Contains(r.Host, ".onion:")), http.StatusSeeOther)
}

func (b BTCPay) status(w http.ResponseWriter, r *http.Request) http.Handler {
	if b.GetStatus != nil {
		json.NewEncoder(w).Encode(b.GetStatus())
		return nil
	} else {
		return http.NotFoundHandler()
	}
}

func (b BTCPay) webhook(w http.ResponseWriter, r *http.Request) http.Handler {
	if b.ErrWebhook == nil {
		b.ErrWebhook = func(err error) http.Handler {
			log.Printf("error processing btcpay webhook: %v", err)
			return nil
		}
	}

	event, err := b.Store.ParseInvoiceWebhook(r)
	if err != nil {
		return b.ErrWebhook(fmt.Errorf("getting event: %w", err))
	}
	purchaseID, paymentKey, _ := strings.Cut(event.InvoiceMetadata.OrderID, ":")

	switch event.Type {
	case btcpay.EventInvoiceProcessing:
		if err := b.Purchases.SetPurchaseProcessing(purchaseID, paymentKey); err != nil {
			return b.ErrWebhook(fmt.Errorf("setting purchase %s processing: %w", purchaseID, err))
		}
		return nil
	case btcpay.EventInvoiceSettled:
		if err := b.Purchases.SetPurchasePaid(purchaseID, paymentKey, "BTCPay"); err != nil {
			return b.ErrWebhook(fmt.Errorf("setting purchase %s paid: %w", purchaseID, err))
		}
		return nil
	case btcpay.EventInvoicePaymentSettled:
		amountEuro, _ := strconv.ParseFloat(event.Payment.Value, 64)
		amountCents := int(math.Round(amountEuro * 100.0))
		if err := b.Purchases.PaymentSettled(purchaseID, paymentKey, "BTCPay", event.Payment.ID, amountCents); err != nil {
			return b.ErrWebhook(fmt.Errorf("setting purchase %s payment %s of %d: %w", purchaseID, event.Payment.ID, amountCents, err))
		}
		return nil
	default:
		return b.ErrWebhook(fmt.Errorf("unknown event type: %s", event.Type))
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
