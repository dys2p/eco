package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dys2p/eco/lang"
	"github.com/dys2p/paypal"
)

var payPalTmpl = template.Must(template.ParseFS(htmlfiles, "paypal-checkout.html"))

type paypalTmplData struct {
	lang.Lang
	ClientID  string
	Reference string
}

// PayPal does the PayPal Standard Checkout described at https://developer.paypal.com/docs/checkout/standard/
type PayPal struct {
	Config    *paypal.Config
	Purchases PurchaseRepo
}

func (PayPal) ID() string {
	return "paypal-checkout"
}

func (PayPal) Name(r *http.Request) string {
	return "PayPal"
}

func (p PayPal) PayHTML(r *http.Request, purchaseID, paymentKey string) (template.HTML, error) {
	b := &bytes.Buffer{}
	err := payPalTmpl.Execute(b, paypalTmplData{
		Lang:      lang.Get(r),
		ClientID:  p.Config.ClientID,
		Reference: purchaseID + ":" + paymentKey,
	})
	return template.HTML(b.String()), err
}

func (p PayPal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	switch r.URL.Path {
	case "/payment/paypal-checkout/create-order":
		if err := p.createTransaction(w, r); err != nil {
			log.Printf("error creating PayPal transaction: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/payment/paypal-checkout/capture-order":
		if err := p.captureTransaction(w, r); err != nil {
			log.Printf("error capturing PayPal transaction: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (p PayPal) createTransaction(w http.ResponseWriter, r *http.Request) error {
	reference, _ := io.ReadAll(r.Body)
	purchaseID, paymentKey, _ := strings.Cut(string(reference), ":")

	sumCents, err := p.Purchases.PurchaseSumCents(purchaseID, paymentKey)
	if err != nil {
		return fmt.Errorf("getting sum: %w", err)
	}

	authResult, err := p.Config.Auth()
	if err != nil {
		return err
	}

	generateOrderResponse, err := p.Config.CreateOrder(authResult, "Purchase "+purchaseID, purchaseID, paymentKey, sumCents)
	if err != nil {
		return err
	}

	// 5. Return a successful response to the client with the order ID
	successResponse, err := json.Marshal(&paypal.SuccessResponse{OrderID: generateOrderResponse.ID})
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(successResponse)
	return nil
}

type captureRequest struct {
	OrderID string `json:"orderID"`
}

// advantage over webhook: this works on localhost
func (p PayPal) captureTransaction(w http.ResponseWriter, r *http.Request) error {
	var captureReq captureRequest
	if err := json.NewDecoder(r.Body).Decode(&captureReq); err != nil {
		return fmt.Errorf("decoding capture request: %w", err)
	}

	authResult, err := p.Config.Auth()
	if err != nil {
		return fmt.Errorf("getting auth: %w", err)
	}

	// 2a. Get the order ID from the request body
	// 3. Call PayPal to capture the order
	captureResponse, err := p.Config.Capture(authResult, captureReq.OrderID)
	if err != nil {
		return fmt.Errorf("capturing response: %w", err)
	}

	purchaseID, paymentKey, _ := strings.Cut(captureResponse.PurchaseUnits[0].ReferenceID, ":")

	log.Printf("[%s] captured transaction: order: %s, capture: %s", purchaseID+":"+paymentKey, captureReq.OrderID, captureResponse.PurchaseUnits[0].Payments.Captures[0].ID)

	if err := p.Purchases.SetPurchasePaid(purchaseID, paymentKey); err != nil {
		return err
	}

	// not in paypal docs: must return some json
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("true"))

	return nil
}

func (PayPal) VerifiesAdult() bool {
	return true
}
