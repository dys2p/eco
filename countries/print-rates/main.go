// Command print-rates prints the VAT rates of European Union countries in a format similar to https://europa.eu/youreurope/business/taxation/vat/vat-rules-rates/index_en.htm#shortcut-5 so we can easily diff it.
package main

import (
	"fmt"
	"strconv"

	"github.com/dys2p/eco/countries"
)

func main() {
	for _, c := range countries.EuropeanUnion {
		// Country code
		fmt.Print(c, "\t")
		// Standard rate
		fmt.Print(fmtPercent(c.VAT()["standard"]), "\t")
		// Reduced rate
		if r1 := c.VAT()["reduced-1"]; r1 > 0 {
			fmt.Print(fmtPercent(r1))
		} else {
			fmt.Print("-")
		}
		if r2 := c.VAT()["reduced-2"]; r2 > 0 {
			fmt.Print(" / ", fmtPercent(r2))
		}
		fmt.Print("\t")
		// Super reduced rate
		if sr := c.VAT()["super-reduced"]; sr > 0 {
			fmt.Print(fmtPercent(sr))
		} else {
			fmt.Print("-")
		}
		fmt.Print("\t")
		// Parking rate
		if pr := c.VAT()["parking"]; pr > 0 {
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
