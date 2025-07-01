package delivery

import (
	"embed"
	"net/mail"

	"golang.org/x/text/message"
)

//go:embed template.html
var TemplateFS embed.FS // usage: t = template.Must(t.ParseFS(delivery.TemplateFS, "*"))

// ShippingAddressView is the data for template "shipping-address-view".
//
// It does not contain a country. It does not check Method.IsShipping.
type ShippingAddressView struct {
	*message.Printer
	Address
	AddressTypes []AddressType
}

// ShippingAddressFormElements is the data for template "shipping-address-form-elements".
//
// It does not contain a country. It does not check Method.IsShipping.
type ShippingAddressFormElements struct {
	*message.Printer
	Address             Address
	AddressElsewhere    bool
	AddressOptions      []AddressTypeOption // empty if there is just one option
	CheckErrors         bool
	SelectedAddressType AddressType
}

func (f ShippingAddressFormElements) ErrAddressCity() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.City == ""
}

func (f ShippingAddressFormElements) ErrAddressCustomerID() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.SelectedAddressType.CustomerIDRequired != "" && f.Address.CustomerID == ""
}

func (f ShippingAddressFormElements) ErrAddressEmail() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.SelectedAddressType.EmailRequired && !emailAddressValid(f.Address.Email)
}

func (f ShippingAddressFormElements) ErrAddressHouseNumber() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.HouseNumber == ""
}

func (f ShippingAddressFormElements) ErrAddressLastName() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.LastName == ""
}

func (f ShippingAddressFormElements) ErrAddressPostcode() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.Postcode == ""
}

func (f ShippingAddressFormElements) ErrAddressStreet() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.Street == ""
}

func (f ShippingAddressFormElements) HasErr() bool {
	return f.ErrAddressCity() || f.ErrAddressCustomerID() || f.ErrAddressEmail() || f.ErrAddressHouseNumber() || f.ErrAddressLastName() || f.ErrAddressPostcode() || f.ErrAddressStreet()
}

func emailAddressValid(addr string) bool {
	_, err := mail.ParseAddress(addr)
	return err == nil
}
