package delivery

import (
	"embed"
	"net/mail"
)

//go:embed template.html
var TemplateFS embed.FS // usage: t = template.Must(t.ParseFS(delivery.TemplateFS, "*"))

type ShippingForm struct {
	Address             Address
	AddressElsewhere    bool
	AddressOptions      []AddressTypeOption // empty if there is just one option
	CheckErrors         bool
	SelectedAddressType AddressType
}

func (f ShippingForm) ErrAddressCity() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.City == ""
}

func (f ShippingForm) ErrAddressCustomerID() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.SelectedAddressType.CustomerIDRequired != "" && f.Address.CustomerID == ""
}

func (f ShippingForm) ErrAddressEmail() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.SelectedAddressType.EmailRequired && !emailAddressValid(f.Address.Email)
}

func (f ShippingForm) ErrAddressHouseNumber() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.HouseNumber == ""
}

func (f ShippingForm) ErrAddressLastName() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.LastName == ""
}

func (f ShippingForm) ErrAddressPostcode() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.Postcode == ""
}

func (f ShippingForm) ErrAddressStreet() bool {
	return f.CheckErrors && !f.AddressElsewhere && f.Address.Street == ""
}

func (form ShippingForm) HasErr() bool {
	return form.ErrAddressCity() || form.ErrAddressCustomerID() || form.ErrAddressEmail() || form.ErrAddressHouseNumber() || form.ErrAddressLastName() || form.ErrAddressPostcode() || form.ErrAddressStreet()
}

func emailAddressValid(addr string) bool {
	_, err := mail.ParseAddress(addr)
	return err == nil
}
