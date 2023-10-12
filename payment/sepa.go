package payment

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/dys2p/eco/lang"
	qrcode "github.com/skip2/go-qrcode"
)

var sepaTmpl = template.Must(template.ParseFS(htmlfiles, "sepa.html"))

type sepaTmplData struct {
	lang.Lang
	Account     SEPAAccount
	Amount      float64
	EPCImageSrc string
	Purpose     string
}

type SEPAAccount struct {
	Holder   string
	IBAN     string
	BIC      string
	BankName string
}

type SEPA struct {
	Account   SEPAAccount
	Purchases PurchaseRepo
}

func (SEPA) ID() string {
	return "sepa"
}

func (SEPA) Name(r *http.Request) string {
	return lang.Get(r).Tr("SEPA Bank Transfer")
}

func (sepa SEPA) PayHTML(r *http.Request, purchaseID, paymentKey string) (template.HTML, error) {
	eurocents, err := sepa.Purchases.PurchaseSumCents(purchaseID, paymentKey)
	if err != nil {
		log.Printf("error getting purchase sum from database: %v", err)
		return template.HTML("Error getting purchase information from database"), nil
	}

	epcString := `BCD
001
1
SCT
` + removeWhitespaces(sepa.Account.BIC) + `
` + strings.TrimSpace(sepa.Account.Holder) + `
` + removeWhitespaces(sepa.Account.IBAN) + `
EUR` + fmt.Sprintf("%.2f", float64(eurocents)/100.0) + `
GDSV

` + purchaseID + `
SEPA payment for purchase` // GDSV = Purchase & Sale of Goods and Services

	epcPNG, err := qrcode.Encode(epcString, qrcode.Medium, 256)
	if err != nil {
		log.Printf("error creating EPC QR code: %v", err) // don't exit
	}

	buf := &bytes.Buffer{}
	err = sepaTmpl.Execute(buf, sepaTmplData{
		Lang:        lang.Get(r),
		Account:     sepa.Account,
		Amount:      float64(eurocents) / 100.0,
		EPCImageSrc: base64.StdEncoding.EncodeToString(epcPNG),
		Purpose:     purchaseID,
	})
	return template.HTML(buf.String()), err
}

func (SEPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (SEPA) VerifiesAdult() bool {
	return false
}

func removeWhitespaces(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			result.WriteRune(ch)
		}
	}
	return result.String()
}
