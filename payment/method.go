// Package payment implements payment methods.
//
// Register the http.Handler for POST requests under /id:
//
//	router.Handler(http.MethodPost, fmt.Sprintf("/payment/%s/*path", paymentMethod.ID()), paymentMethod)
//
// The handler will be publicly available.
package payment

import (
	"errors"
	"html/template"
	"net/http"
)

// Method is the interface that wraps payment methods.
//
// PayHTML takes a purchase ID and an optional payment key.
// The purchase ID should be unique. It is usually short and suitable for bank
// transfer forms, bookkeeping etc. It must not contain a colon.
// The optional payment key should have a high entropy. It can prevent the loss
// of goods if an purchase ID is accidentally or maliciously issued twice.
// Even if a payment key is used, payment methods should store the purchase ID
// because tax accounting may require a connection between payment and purchase.
type Method interface {
	ID() string
	Name(langstr string) string
	PayHTML(purchaseID, paymentKey, langstr string) (template.HTML, error)
	ServeHTTP(w http.ResponseWriter, r *http.Request, langstr string)
	VerifiesAdult() bool
}

func Get(methods []Method, id string) (Method, error) {
	for _, m := range methods {
		if m.ID() == id {
			return m, nil
		}
	}
	if len(methods) > 0 {
		return methods[0], nil
	}
	return nil, errors.New("no payment methods found")
}

type PurchaseRepo interface {
	PurchaseCreationDate(purchaseID, paymentKey string) (string, error) // yyyy-mm-dd
	PurchaseSumCents(purchaseID, paymentKey string) (int, error)
	SetPurchasePaid(purchaseID, paymentKey string) error
	SetPurchaseProcessing(purchaseID, paymentKey string) error
}
