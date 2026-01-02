// Package rates retrieves and stores daily exchange rates.
package rates

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/dys2p/eco/lang"
)

type History struct {
	Database    *SQLiteDB
	GetBuyRates func() (map[string]float64, error)
	Synced      bool // updated today or yesterday
}

// MakeAndRun starts a goroutine which calls GetBuyRates every 45-60 minutes. If GetBuyRates returns rates, they are inserted into the database and GetBuyRates is not called until the next day.
func MakeAndRun(sqlitePath string, getBuyRates func() (map[string]float64, error)) (*History, error) {
	db, err := OpenDB(sqlitePath)
	if err != nil {
		return nil, err
	}

	h := &History{
		Database:    db,
		GetBuyRates: getBuyRates,
	}

	go func() {
		for ; true; time.Sleep(time.Duration(45*int64(time.Minute) + rand.Int63n(15*int64(time.Minute)))) {
			now := time.Now()
			today := now.Format("2006-01-02")
			yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

			lastUpdateDate, err := h.Database.LatestDate(today)
			if err != nil {
				log.Printf("\033[31m"+"error getting latest date from database: %v"+"\033[0m", err)
				h.Synced = false // database error means no good for sync status
				continue
			}

			h.Synced = lastUpdateDate == today || lastUpdateDate == yesterday
			if lastUpdateDate == today {
				continue // already updated today
			}

			buyRates, err := h.GetBuyRates()
			if err != nil {
				log.Printf("\033[31m"+"error getting rates: %v"+"\033[0m", err)
				continue
			}
			if len(buyRates) == 0 {
				continue // nothing to insert
			}
			if err := h.Database.Insert(today, buyRates); err == nil {
				log.Println("\033[32m" + "updated foreign cash rates" + "\033[0m")
			} else {
				log.Printf("\033[31m"+"error inserting rates: %v"+"\033[0m", err)
			}
		}
	}()

	return h, nil
}

func (h *History) Options(effectiveDate string, value float64) ([]Option, error) {
	date, err := h.Database.LatestDate(effectiveDate)
	if err != nil {
		return nil, fmt.Errorf("getting latest date: %w", err)
	}

	rs, err := h.Database.Get(date)
	if err != nil {
		return nil, fmt.Errorf("getting rates for %s from database: %w", date, err) // unlikely because we got the date from the database
	}
	var options []Option
	for currency, rate := range rs {
		options = append(options, Option{
			Currency: currency,
			Price:    value * rate,
		})
	}
	slices.SortFunc(options, func(a, b Option) int {
		return cmp.Compare(a.Currency, b.Currency)
	})
	return options, nil
}

// SyncedHandler writes JSON true or false. The returned handler can be queried extensively because it just reads a variable.
func (h *History) SyncedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.Synced)
	}
}

type Option struct {
	Currency string // from GetBuyRates
	Price    float64
}

// ISO 4217
func (opt Option) Tr(l lang.Lang) string {
	switch opt.Currency {
	case "AUD":
		return l.Tr("Australian dollars")
	case "BGN":
		return l.Tr("Bulgarian lev")
	case "CAD":
		return l.Tr("Canadian dollars")
	case "CHF":
		return l.Tr("Swiss francs")
	case "CNY":
		return l.Tr("Chinese renminbi")
	case "CZK":
		return l.Tr("Czech koruna")
	case "DKK":
		return l.Tr("Danish krone")
	case "GBP":
		return l.Tr("Pound sterling")
	case "ISK":
		return l.Tr("Icelandic króna")
	case "JPY":
		return l.Tr("Japanese yen")
	case "ILS":
		return l.Tr("New Israeli shekel (NIS)")
	case "NOK":
		return l.Tr("Norwegian krone")
	case "NZD":
		return l.Tr("New Zealand dollars")
	case "PLN":
		return l.Tr("Polish złoty")
	case "RON":
		return l.Tr("Romanian leu")
	case "RSD":
		return l.Tr("Serbian dinar")
	case "SEK":
		return l.Tr("Swedish krona")
	case "TWD":
		return l.Tr("New Taiwan dollars")
	case "USD":
		return l.Tr("United States dollars")
	default:
		return ""
	}
}
