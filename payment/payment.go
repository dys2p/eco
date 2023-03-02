// Package payment implements payment methods.
//
// Register the http.Handler for POST requests under /id:
//
//	router.Handler(http.MethodPost, fmt.Sprintf("/%s/*path", paymentMethod.ID()), paymentMethod)
package payment

import (
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
}
