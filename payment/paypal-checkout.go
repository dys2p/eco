package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/dys2p/eco/httputil"
	"github.com/dys2p/eco/lang"
	"github.com/dys2p/go-paypal"
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

	Err func(err error) http.Handler // should write an error message or error template to the ResponseWriter
}

func (p PayPal) Handler() http.Handler {
	if p.Err == nil {
		p.Err = func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf("error processing PayPal transaction: %v", err)
				w.Write([]byte("There was an error processing your PayPal purchase. We have been notified and will fix it soon. Sorry for the inconvenience."))
			})
		}
	}

	var mux = http.NewServeMux()
	mux.Handle("POST /payment/paypal-checkout/create-order", httputil.HandlerFunc(p.createTransaction))
	mux.Handle("POST /payment/paypal-checkout/capture-order", httputil.HandlerFunc(p.captureTransaction))
	return mux
}

func (PayPal) ID() string {
	return "paypal-checkout"
}

func (PayPal) Name(l lang.Lang) string {
	return "PayPal"
}

func (p PayPal) PayHTML(purchaseID, paymentKey string, l lang.Lang) (template.HTML, error) {
	b := &bytes.Buffer{}
	err := payPalTmpl.Execute(b, paypalTmplData{
		Lang:      l,
		ClientID:  p.Config.ClientID,
		Reference: purchaseID + ":" + paymentKey,
	})
	return template.HTML(b.String()), err
}

func (PayPal) VerifiesAdult() bool {
	return true
}

func (p PayPal) createTransaction(w http.ResponseWriter, r *http.Request) http.Handler {
	reference, _ := io.ReadAll(r.Body)
	purchaseID, paymentKey, _ := strings.Cut(string(reference), ":")

	sumCents, err := p.Purchases.PurchaseSumCents(purchaseID, paymentKey)
	if err != nil {
		return p.Err(fmt.Errorf("getting purchase sum: %w", err))
	}

	authResult, err := p.Config.Auth()
	if err != nil {
		return p.Err(err)
	}

	generateOrderResponse, err := p.Config.CreateOrder(authResult, "Purchase "+purchaseID, purchaseID, paymentKey, sumCents)
	if err != nil {
		return p.Err(err)
	}

	// 5. Return a successful response to the client with the order ID
	successResponse, err := json.Marshal(&paypal.SuccessResponse{OrderID: generateOrderResponse.ID})
	if err != nil {
		return p.Err(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(successResponse)
	return nil
}

type captureRequest struct {
	OrderID string `json:"orderID"`
}

// advantage over webhook: this works on localhost
func (p PayPal) captureTransaction(w http.ResponseWriter, r *http.Request) http.Handler {
	var captureReq captureRequest
	if err := json.NewDecoder(r.Body).Decode(&captureReq); err != nil {
		return p.Err(fmt.Errorf("decoding capture request: %w", err))
	}

	authResult, err := p.Config.Auth()
	if err != nil {
		return p.Err(fmt.Errorf("getting auth: %w", err))
	}

	// 2a. Get the order ID from the request body
	// 3. Call PayPal to capture the order
	captureResponse, err := p.Config.Capture(authResult, captureReq.OrderID)
	if err != nil {
		return p.Err(fmt.Errorf("capturing response: %w", err))
	}

	if len(captureResponse.PurchaseUnits) == 0 {
		return p.Err(errors.New("no purchase units"))
	}
	if len(captureResponse.PurchaseUnits[0].Payments.Captures) == 0 {
		return p.Err(errors.New("no captures"))
	}

	var (
		amountStr  = captureResponse.PurchaseUnits[0].Payments.Captures[0].Amount.Value
		captureID  = captureResponse.PurchaseUnits[0].Payments.Captures[0].ID
		paymentKey = captureResponse.PurchaseUnits[0].ReferenceID
		purchaseID = captureResponse.PurchaseUnits[0].Payments.Captures[0].InvoiceID
	)
	amountEuro, _ := strconv.ParseFloat(amountStr, 64)
	amountCents := int(math.Round(amountEuro * 100.0))

	log.Printf("[%s] captured transaction: order: %s, capture: %s", purchaseID+":"+paymentKey, captureReq.OrderID, captureID)

	if err := p.Purchases.PaymentSettled(purchaseID, paymentKey, "PayPal", captureID, amountCents); err != nil {
		return p.Err(err)
	}

	if err := p.Purchases.SetPurchasePaid(purchaseID, paymentKey, "PayPal"); err != nil {
		return p.Err(err)
	}

	// not mentioned in paypal docs: must return some json
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("true"))
	return nil
}
