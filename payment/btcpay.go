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
)

func init() {
	log.Println(`Don't forget to set up the BTCPay webhook for your store: URL: "/payment/btcpay/webhook", event: "An invoice has been settled"`)
}

type BTCPay struct {
	CustomName string
	ShopURL    string // for redirect URL, like "https://example.com"
	Store      btcpay.Store
	Purchases  PurchaseRepo
}

type createdInvoice struct {
	CheckoutLink string
	Time         int64
}

var lastInvoice = make(map[string]createdInvoice) // key: purchase ID

var btcpayTmpl = template.Must(template.New("").Parse(`
	<p>Bezahle den angegebenen Betrag in Monero (XMR) oder Bitcoin (BTC). Der Betrag muss innerhalb von 60 Minuten vollständig und als einzelne Transaktion auf der angegebenen Adresse eingehen.</p>
	<p>Hinweis: Das BTCPay-Fenster bestätigt deine Zahlung, sobald die Zahlung in der Blockchain sichtbar ist. Der Status deiner Bestellung wird jedoch erst einige Minuten später aktualisiert, wenn die Transaktion ausreichend Bestätigungen hat.</p>
	<form action="/payment/btcpay/create-invoice" method="post">
		<input type="hidden" name="purchase-id" value="{{.}}">
		<button type="submit" class="btn btn-success">Zahlungsaufforderung erzeugen</button>
	</form>
`))

func (BTCPay) ID() string {
	return "btcpay"
}

func (b BTCPay) Name() string {
	if b.CustomName == "" {
		return "BTCPay"
	} else {
		return b.CustomName
	}
}

func (b BTCPay) PayHTML(purchaseID string) (template.HTML, error) {
	buf := &bytes.Buffer{}
	err := btcpayTmpl.Execute(buf, purchaseID)
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

	invoiceRequest := &btcpay.InvoiceRequest{
		Amount:   float64(sumCents) / 100.0,
		Currency: "EUR",
	}
	invoiceRequest.ExpirationMinutes = 60
	invoiceRequest.DefaultLanguage = "de-DE"
	invoiceRequest.OrderID = purchaseID
	invoiceRequest.RedirectURL = fmt.Sprintf("%s/view", b.ShopURL)
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

	if event.Type == btcpay.EventInvoiceSettled {
		if err := b.Purchases.SetPurchasePaid(purchaseID); err == nil {
			return nil
		} else {
			return fmt.Errorf("setting purchase %s paid: %w", purchaseID, err)
		}
	} else {
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

func (BTCPay) VerifiesAdult() bool {
	return false
}
