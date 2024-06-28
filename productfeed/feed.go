// Package productfeed creates an Atom feed according to Google product data specification. See https://support.google.com/merchants/answer/7052112
package productfeed

import (
	"bytes"
	"encoding/xml"
	"strings"

	"golang.org/x/net/html"
)

type Feed struct {
	XMLName   xml.Name  `xml:"http://www.w3.org/2005/Atom feed"`
	Namespace string    `xml:"xmlns:g,attr"` // hack for adding xmlns:g, see https://stackoverflow.com/q/72804320
	ID        string    `xml:"id"`           // "If you have a long-term, renewable lease on your Internet domain name, then you can feel free to use your website's address."
	Title     string    `xml:"title"`        // "Contains a human readable title for the feed. Often the same as the title of the associated website."
	Link      []Link    `xml:"link"`         // "Recommended feed element"
	Updated   string    `xml:"updated"`      // ISO8601
	Products  []Product `xml:"entry"`
}

func (feed Feed) Bytes() ([]byte, error) {
	feed.Namespace = "http://base.google.com/ns/1.0"

	var buf bytes.Buffer
	buf.Write([]byte(xml.Header))
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "\t")
	if err := enc.Encode(feed); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// subset of https://pkg.go.dev/google.golang.org/api/content/v2#Product but with xml
type Product struct {
	Adult            string `xml:"g:adult,omitempty"`        // "yes" or "no", default is "no"
	Availability     string `xml:"g:availability,omitempty"` // "in stock", "out of stock", "preorder"
	Brand            string `xml:"g:brand,omitempty"`
	Condition        string `xml:"g:condition,omitempty"`
	Description      string `xml:"g:description,omitempty"`
	Gtin             string `xml:"g:gtin,omitempty"`
	Id               string `xml:"g:id,omitempty"`
	IdentifierExists bool   `xml:"g:identifier_exists],omitempty"`
	ImageLink        string `xml:"g:image_link,omitempty"`
	ItemGroupId      string `xml:"g:item_group_id,omitempty"`
	Link             string `xml:"g:link,omitempty"`
	Mpn              string `xml:"g:mpn,omitempty"`
	Price            string `xml:"g:price,omitempty"`
	Title            string `xml:"g:title,omitempty"`
}

// from https://pkg.go.dev/golang.org/x/tools/blog/atom#Link so users don't have to import that
type Link struct {
	Rel      string `xml:"rel,attr,omitempty"`
	Href     string `xml:"href,attr"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Length   uint   `xml:"length,attr,omitempty"`
}

func HTMLtoText(htm string) string {
	var result strings.Builder
	tokenizer := html.NewTokenizerFragment(strings.NewReader(htm), "div")
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.TextToken {
			result.Write(tokenizer.Text())
		}
	}
	return result.String()
}
