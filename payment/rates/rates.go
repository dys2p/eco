// Package rates retrieves and stores daily exchange rates.
package rates

import (
	"cmp"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"time"

	"github.com/dys2p/eco/lang"
)

type History struct {
	Database    *SQLiteDB
	GetBuyRates func(lastUpdateDate string) (map[string]float64, error)
}

// RunDaemon starts a loop which fetches the rates each day (after 9:00 AM) and inserts them into the database. RunDaemon blocks and cannot be stopped.
func (h *History) RunDaemon() error {
	for ; true; time.Sleep(time.Duration(45*int64(time.Minute) + rand.Int63n(15*int64(time.Minute)))) {
		if time.Now().Hour() < 9 {
			continue // too early in the morning, today's rates are probably not available yet
		}
		today := time.Now().Format("2006-01-02")
		lastUpdateDate, err := h.Database.LatestDate(today)
		if err != nil {
			log.Printf("\033[31m"+"error getting latest date: %v"+"\033[0m", err)
			continue
		}
		if lastUpdateDate == today {
			continue // already updated today
		}
		buyRates, err := h.GetBuyRates(lastUpdateDate)
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
	return nil
}

func (h *History) Options(maxDate string, value float64) ([]Option, error) {
	date, err := h.Database.LatestDate(maxDate)
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

// Synced returns whether rates have been fetched today or yesterday.
func (h *History) Synced() bool {
	lastUpdateDate, err := h.Database.LatestDate(time.Now().Format("2006-01-02"))
	if err != nil {
		return false
	}
	min := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return lastUpdateDate >= min
}

type Option struct {
	Currency string // from GetBuyRates result
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
