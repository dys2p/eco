// Package health provides a widget which displays the synchronization status of payment methods. So far it only supports BTCPay Server.
//
// Register the health handler in your HTTP router:
//
//	router.Handler(http.MethodGet, "/payment-health", health.Server{btcpayStore})
//
// Parse the health template string along with your HTML templates:
//
//	t = template.Must(t.Parse(health.TemplateString))
//
// Execute the template:
//
//	{{template "health"}}
package health

import (
	"encoding/json"
	"net/http"

	"github.com/dys2p/btcpay"
	"github.com/dys2p/eco/payment/rates"
)

type Server struct {
	BTCPay btcpay.Store
	Rates  *rates.History
}

func (srv Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var response []Item

	if srv.BTCPay != nil {
		status, err := srv.BTCPay.GetServerStatus()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, syncStatus := range status.SyncStatuses {
			response = append(response, Item{
				Name:   syncStatus.CryptoCode,
				Synced: syncStatus.ChainHeight == syncStatus.SyncHeight,
			})
		}
	}

	if srv.Rates != nil {
		response = append(response, Item{
			Name:   "Foreign Cash",
			Synced: srv.Rates.Synced(),
		})
	}

	responseData, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseData)
}
