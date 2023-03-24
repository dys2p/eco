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

type Method interface {
	http.Handler
	ID() string
	Name() string
	PayHTML(purchaseID string) (template.HTML, error)
	VerifiesAdult() bool
}

type PurchaseRepo interface {
	PurchaseSumCents(purchaseID string) (int, error)
	SetPurchasePaid(purchaseID string) error
	SetPurchaseProcessing(purchaseID string) error
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
