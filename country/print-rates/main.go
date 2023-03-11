// Command print-rates prints the VAT rates of European Union countries in a format similar to https://europa.eu/youreurope/business/taxation/vat/vat-rules-rates/index_en.htm#shortcut-5 so we can easily diff it.
package main

import (
	"fmt"
	"strconv"

	"github.com/dys2p/eco/country"
)

func main() {
	for _, c := range country.EuropeanUnion {
		// Country code
		fmt.Print(c.ID, "\t")
		// Standard rate
		fmt.Print(fmtPercent(c.VATRates["standard"]), "\t")
		// Reduced rate
		if r1 := c.VATRates["reduced-1"]; r1 > 0 {
			fmt.Print(fmtPercent(r1))
		} else {
			fmt.Print("-")
		}
		if r2 := c.VATRates["reduced-2"]; r2 > 0 {
			fmt.Print(" / ", fmtPercent(r2))
		}
		fmt.Print("\t")
		// Super reduced rate
		if sr := c.VATRates["super-reduced"]; sr > 0 {
			fmt.Print(fmtPercent(sr))
		} else {
			fmt.Print("-")
		}
		fmt.Print("\t")
		// Parking rate
		if pr := c.VATRates["parking"]; pr > 0 {
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
