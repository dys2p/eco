package rates

import (
	"math"
	"testing"
	"time"
)

func TestRates(t *testing.T) {
	db, err := OpenDB("/tmp/rates-test.sqlite3")
	if err != nil {
		t.Fatalf("opening db: %v", err)
	}
	history := History{
		Database: db,
		GetBuyRates: func(lastUpdateDate string) (map[string]float64, error) {
			return map[string]float64{"USD": 1.1, "GBP": 0.85}, nil
		},
	}

	go history.RunDaemon()
	time.Sleep(100 * time.Millisecond)

	got, err := history.Options(time.Now().AddDate(0, 0, 3).Format("2006-01-02"), 100.0) // three days in the future
	if err != nil {
		t.Fatalf("getting rate: %v", err)
	}

	want := []Option{
		{"GBP", 85},
		{"USD", 110},
	}
	for i := range got {
		if got[i].Currency != want[i].Currency {
			t.Fatalf("got %s, want %s", got[i].Currency, want[i].Currency)
		}
		epsilon := math.Abs(got[i].Price - want[i].Price)
		if epsilon > 0.001 {
			t.Fatalf("got %f, want %f", got[i].Price, want[i].Price)
		}
	}
}
