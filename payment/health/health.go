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
	"time"

	"github.com/dys2p/btcpay"
	"github.com/dys2p/eco/payment/rates"
)

type Server struct {
	BTCPay btcpay.Store
	Rates  *rates.History
	status []Item
}

func (srv *Server) Run() {
	for ; true; <-time.Tick(10 * time.Second) {
		var status []Item
		if srv.BTCPay != nil {
			if serverStatus, err := srv.BTCPay.GetServerStatus(); err == nil {
				for _, syncStatus := range serverStatus.SyncStatuses {
					status = append(status, Item{
						Name:   syncStatus.CryptoCode,
						Synced: syncStatus.ChainHeight == syncStatus.SyncHeight,
					})
				}
			}
		}
		if srv.Rates != nil {
			status = append(status, Item{
				Name:   "Foreign Cash",
				Synced: srv.Rates.Synced(),
			})
		}
		srv.status = status
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseData, _ := json.Marshal(srv.status)
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseData)
}
