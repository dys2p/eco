package eco

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

func FmtEuro(cents int) string {
	return fmtCurrency(cents) + " €"
}

func FmtEuroHTML(cents int) template.HTML {
	return template.HTML(fmtCurrency(cents) + "&nbsp;€")
}

func fmtCurrency(cents int) string {
	s := fmt.Sprintf("%.2f", float64(cents)/100.0)
	s = strings.Replace(s, ".", ",", 1)
	s = strings.Replace(s, "-", "−", 1)
	return s
}

func FmtEuroMinMaxHTML(minCents, maxCents int) template.HTML {
	if minCents == 0 && maxCents == 0 {
		return template.HTML("")
	}
	if minCents == maxCents {
		return FmtEuroHTML(minCents)
	}
	return template.HTML(fmt.Sprintf("%s&nbsp;–&nbsp;%s&nbsp;€", fmtCurrency(minCents), fmtCurrency(maxCents)))
}

func FmtEuroPlusMinusHTML(cents int) template.HTML {
	if cents < 0 {
		return FmtEuroHTML(cents) // FmtEuroHTML will add the minus sign
	} else {
		return "+" + FmtEuroHTML(cents)
	}
}

func FmtPercentHTML(f float64) template.HTML {
	s := strconv.FormatFloat(f*100.0, 'f', 2, 64) // precision 2
	s = strings.TrimRight(s, "0")                 // remove trailing zeroes
	s = strings.TrimSuffix(s, ".")                // remove trailing dot
	s = strings.Replace(s, ".", ",", 1)
	return template.HTML(s + "&nbsp;%")
}
