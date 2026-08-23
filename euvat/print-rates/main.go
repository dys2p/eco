// Command print-rates prints the VAT rates of European Union countries in a format similar to https://europa.eu/youreurope/business/taxation/vat/vat-rules-rates/index_en.htm#shortcut-5 so we can easily diff it.
package main

import (
	"fmt"
	"strconv"

	"github.com/dys2p/eco/countries"
	"github.com/dys2p/eco/euvat"
)

func main() {
	for _, c := range countries.EuropeanUnion {
		// Country code
		fmt.Print(c, "\t")
		// Standard rate
		s, _ := euvat.Value(c, euvat.RateStandard)
		fmt.Print(fmtPercent(s), "\t")
		// Reduced rate
		if r1, ok := euvat.Value(c, euvat.RateReduced1); ok {
			fmt.Print(fmtPercent(r1))
		} else {
			fmt.Print("-")
		}
		if r2, ok := euvat.Value(c, euvat.RateReduced2); ok {
			fmt.Print(" / ", fmtPercent(r2))
		}
		fmt.Print("\t")
		// Super reduced rate
		if sr, ok := euvat.Value(c, euvat.RateSuperReduced); ok {
			fmt.Print(fmtPercent(sr))
		} else {
			fmt.Print("-")
		}
		fmt.Print("\t")
		// Parking rate
		if pr, ok := euvat.Value(c, euvat.RateParking); ok {
			fmt.Print(fmtPercent(pr))
		} else {
			fmt.Print("-")
		}
		fmt.Print("\n")
	}
}

func fmtPercent(f float64) string {
	return strconv.FormatFloat(f*100.0, 'g', 3, 64)
}
