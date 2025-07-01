// Package delivery provides structs for delivery addresses and methods and their frontend options.
package delivery

import (
	"fmt"
	"strings"

	"github.com/dys2p/eco/countries"
)

// Address is a shipping address. It does not contain a country as we store the country separately, because it must be preserved for tax reasons. Address contains separate fields for first and last name and street and house number because some shippers may require it and because people easily forget to enter it else.
type Address struct {
	FirstName   string
	LastName    string
	Supplement  string
	CustomerID  string // e. g. DHL PostNumber
	Street      string
	HouseNumber string
	Postcode    string
	City        string
	Email       string
	Phone       string
}

func (a Address) NotEmpty() bool {
	return a.FirstName != "" || a.LastName != "" || a.Supplement != "" || a.CustomerID != "" || a.Street != "" || a.HouseNumber != "" || a.Postcode != "" || a.City != "" || a.Email != "" || a.Phone != ""
}

func (a Address) lines() []string {
	var lines []string
	if name := strings.TrimSpace(a.FirstName + " " + a.LastName); name != "" {
		lines = append(lines, name)
	}
	if a.Supplement != "" {
		lines = append(lines, a.Supplement)
	}
	if a.CustomerID != "" {
		lines = append(lines, a.CustomerID)
	}
	if street := strings.TrimSpace(a.Street + " " + a.HouseNumber); street != "" {
		lines = append(lines, street)
	}
	if city := strings.TrimSpace(a.Postcode + " " + a.City); city != "" {
		lines = append(lines, city)
	}
	if a.Email != "" {
		lines = append(lines, a.Email)
	}
	if a.Phone != "" {
		lines = append(lines, a.Phone)
	}
	return lines
}

func (a Address) String() string {
	return strings.Join(a.lines(), "\n")
}

func (a Address) StringComma() string {
	return strings.Join(a.lines(), ", ")
}

type AddressType struct {
	// ID and Name are only required if the DeliveryMethod has more than one AddressType
	ID   string
	Name string

	CustomerIDRequired string // name, e. g. "DHL PostNumber"
	Email              bool
	EmailRequired      bool
	Phone              bool
	StreetName         string
	Supplement         bool
}

type AddressTypeOption struct {
	AddressType
	Selected bool
}

type Method struct {
	ID          string
	Name        string
	Description string

	AddressTypes []AddressType // empty slice means one zero-valued AddressType
	Details      func(weight, netValue int, country countries.Country) (netPrice, minDays, maxDays int, supported bool)
	IsShipping   bool
	TrackingLink string // use %s placeholder for tracking number
	WarnBackend  string
}

func (method *Method) AddressTypeOptions(selectedID string) (options []AddressTypeOption, selected AddressType) {
	if method == nil {
		return
	}
	// return empty options slice if there is exactly one option
	if len(method.AddressTypes) == 1 {
		return nil, method.AddressTypes[0]
	}
	// collect options
	for _, t := range method.AddressTypes {
		if t.ID == selectedID {
			selected = t
		}
		options = append(options, AddressTypeOption{
			AddressType: t,
			Selected:    t.ID == selectedID,
		})
	}
	// if none is selected, then select first
	if selected == (AddressType{}) && len(options) > 0 {
		options[0].Selected = true
		selected = options[0].AddressType
	}
	return
}

func (method *Method) FmtTrackingLink(id string) string {
	id = strings.TrimSpace(id)
	if strings.HasPrefix(id, "https://") {
		return id
	} else {
		return fmt.Sprintf(method.TrackingLink, id)
	}
}

// GetAddressTypes returns method.AddressTypes or nil if method is nil. It is just a shortcut so you don't have to check method != nil.
func (method *Method) GetAddressTypes() []AddressType {
	if method == nil {
		return nil
	}
	return method.AddressTypes
}

type MethodOption struct {
	Method
	Gross           int
	Selected        bool
	ShippingDaysMax int
	ShippingDaysMin int
}

// Support represents whether a given product can be delivered to a given country. ShippingDaysMax and ShippingDaysMin should be across all shipping methods (not just the fastest method, which might not be the cheapest).
type Support struct {
	Pickup          bool
	Shipping        bool
	ShippingDaysMax int
	ShippingDaysMin int
	ShippingNetMin  int
}

func (support Support) Any() bool {
	return support.Pickup || support.Shipping
}

func (support Support) None() bool {
	return !support.Pickup && !support.Shipping
}

func (support Support) PickupOnly() bool {
	return support.Pickup && !support.Shipping
}
