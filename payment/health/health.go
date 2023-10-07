package health

import (
	"encoding/json"
	"net/http"

	"github.com/dys2p/btcpay"
)

type Server struct {
	BTCPay btcpay.Store
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

	responseData, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseData)
}
