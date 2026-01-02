package rates

import (
	"math"
	"testing"
	"time"
)

func TestRates(t *testing.T) {
	history, err := MakeAndRun(
		"/tmp/rates-test.sqlite3",
		func() (map[string]float64, error) {
			return map[string]float64{"USD": 1.1, "GBP": 0.85}, nil
		},
	)
	if err != nil {
		t.Fatalf("opening db: %v", err)
	}

	time.Sleep(100 * time.Millisecond) // wait until database is created

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
