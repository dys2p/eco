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

	"github.com/dys2p/paypal"
)

// PayPal does the PayPal Standard Checkout described at https://developer.paypal.com/docs/checkout/standard/
type PayPal struct {
	Config    *paypal.Config
	Purchases PurchaseRepo
}

type paypalTmplData struct {
	ClientID   string
	PurchaseID string
}

var payPalTmpl = template.Must(template.New("").Parse(`
	<script src="https://www.paypal.com/sdk/js?client-id={{.ClientID}}&currency=EUR"></script>
	<!-- Set up a container element for the button -->
	<div id="paypal-button-container" style="text-align: center;"></div>
	<script>
		paypal.Buttons({
			// Order is created on the server and the order id is returned
			createOrder() {
				return fetch("/payment/paypal-checkout/create-order", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: "{{.PurchaseID}}"
				})
				.then((response) => response.json())
				.then((order) => order.id);
			},
			onApprove(data) {
				console.log(data);
				return fetch("/payment/paypal-checkout/capture-order", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify({
						orderID: data.orderID
					})
				})
				.then((response) => response.json())
				.then((orderData) => {
					window.location.href = "/view";
				});
			}
		}).render('#paypal-button-container');
	</script>

	<p>Wir übermitteln nur die Bestellnummer an PayPal. Deine bestellten Artikel sowie die Details zu Lieferung oder Abholung werden nicht an PayPal gesendet.</p>
	<p>Falls du TOR oder einen VPN benutzt: Die angezeigten Bezahlmöglichkeiten sind von der Länderzuordnung deiner IP-Adresse abhängig. Darüber hinaus blockiert PayPal manche TOR Exit Nodes. In dem Fall versuche es mit „New Circuit for this Site“.</p>
`))

func (PayPal) ID() string {
	return "paypal-checkout"
}

func (PayPal) Name() string {
	return "PayPal"
}

func (p PayPal) PayHTML(purchaseID string) (template.HTML, error) {
	b := &bytes.Buffer{}
	err := payPalTmpl.Execute(b, paypalTmplData{
		ClientID:   p.Config.ClientID,
		PurchaseID: purchaseID,
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
	purchaseIDBytes, _ := io.ReadAll(r.Body)
	purchaseID := string(purchaseIDBytes)

	sumCents, err := p.Purchases.PurchaseSumCents(purchaseID)
	if err != nil {
		return fmt.Errorf("getting sum: %w", err)
	}

	authResult, err := p.Config.Auth()
	if err != nil {
		return err
	}

	generateOrderResponse, err := p.Config.CreateOrder(authResult, purchaseID, sumCents)
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

	purchaseID := captureResponse.PurchaseUnits[0].ReferenceID

	log.Printf("[%s] captured transaction: order: %s, capture: %s", purchaseID, captureReq.OrderID, captureResponse.PurchaseUnits[0].Payments.Captures[0].ID)

	if err := p.Purchases.SetPurchasePaid(purchaseID); err != nil {
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
