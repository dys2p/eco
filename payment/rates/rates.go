// Package rates retrieves and stores daily exchange rates.
package rates

import (
	"cmp"
	"errors"
	"log"
	"slices"
	"time"

	"github.com/dys2p/eco/lang"
)

type History struct {
	Currencies  []string // ISO 4217
	GetBuyRates func(currencies []string, lastUpdateDate string) (map[string]float64, error)
	Repository  *SQLiteDB
}

// RunDaemon starts a loop which fetches the rates every hour and inserts them into the repository. The function blocks.
func (h *History) RunDaemon() error {
	for ; true; <-time.Tick(time.Hour) {
		lastUpdateDate, err := h.Repository.LatestDate()
		if err != nil {
			log.Printf("\033[31m"+"error getting latest date: %v"+"\033[0m", err)
			continue
		}
		if lastUpdateDate == time.Now().Format("2006-01-02") {
			continue // already updated today
		}
		buyRates, err := h.GetBuyRates(h.Currencies, lastUpdateDate)
		if err != nil {
			log.Printf("\033[31m"+"error getting rates: %v"+"\033[0m", err)
			continue
		}
		if len(buyRates) == 0 {
			continue // nothing to insert
		}
		if err := h.Repository.Insert(time.Now().Format("2006-01-02"), buyRates); err != nil {
			log.Printf("\033[31m"+"error inserting rates: %v"+"\033[0m", err)
		}
	}
	return nil
}

// Get tries the given date and three previous days.
func (h *History) Get(date string, value float64) ([]Option, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	for _, days := range []int{0, -1, -2, -3} {
		if rates, err := h.Repository.Get(t.AddDate(0, 0, days).Format("2006-01-02")); err == nil {
			var options []Option
			for currency, rate := range rates {
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
	}

	return nil, errors.New("no rates found")
}

// Synced returns whether rates have been updated since three days ago.
func (h *History) Synced() bool {
	lastUpdateDate, err := h.Repository.LatestDate()
	if err != nil {
		return false
	}
	min := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	return lastUpdateDate >= min
}

type Option struct {
	Currency string // ISO 4217
	Price    float64
}

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
