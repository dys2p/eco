package delivery

import (
	"iter"
	"math"
	"slices"

	"github.com/dys2p/eco/countries"
)

type Config[P any] struct {
	Methods        []Method
	ForbidCountry  func(countries.Country, P) bool         // optional, for country laws
	ForbidDelivery func(countries.Country, Method, P) bool // optional, for shipping company size limits and terms of service
}

func (conf Config[P]) Method(id string) *Method {
	for _, method := range conf.Methods {
		if method.ID == id {
			return &method
		}
	}
	return nil
}

// checks Details, ForbidCountry and ForbidDelivery
func (conf Config[P]) Valid(method Method, weightGrams, netPrice int, country countries.Country, products iter.Seq[P]) bool {
	if _, _, _, supported := method.Details(weightGrams, netPrice, country); !supported {
		return false
	}
	for product := range products {
		if conf.ForbidCountry != nil && conf.ForbidCountry(country, product) {
			return false
		}
		if conf.ForbidDelivery != nil && conf.ForbidDelivery(country, method, product) {
			return false
		}
	}
	return true
}

// for cart view, does not set MethodOption.GrossPrice
func (conf Config[P]) Options(selected *Method, weightGrams, netPrice int, country countries.Country, products iter.Seq[P]) (options []MethodOption, none []P, pickupOnly []P) {
	var selectedID string
	if selected != nil {
		selectedID = selected.ID
	}

	// check Details
	for _, method := range conf.Methods {
		netDeliveryPrice, minDays, maxDays, supported := method.Details(weightGrams, netPrice, country)
		if supported {
			options = append(options, MethodOption{
				Method:          method,
				NetPrice:        netDeliveryPrice,
				Selected:        method.ID == selectedID,
				ShippingDaysMax: maxDays,
				ShippingDaysMin: minDays,
			})
		}
	}

	// check ForbidDelivery, remove affected options
	var removeOptions = make(map[int]any) // don't remove instantly because we want to collect all responsible products (not just the first product)
	for product := range products {
		var pickup, shipping bool
		for i, option := range options {
			if conf.ForbidDelivery != nil && conf.ForbidDelivery(country, option.Method, product) {
				removeOptions[i] = struct{}{}
				continue
			}
			if option.Method.IsShipping {
				shipping = true
			} else {
				pickup = true
			}
		}
		if !pickup && !shipping {
			none = append(none, product)
		}
		if pickup && !shipping {
			pickupOnly = append(pickupOnly, product)
		}
	}
	for i := range removeOptions {
		options = slices.Delete(options, i, i+1)
	}

	// check ForbidCountry (do this at last because it will make options empty)
	for product := range products {
		if conf.ForbidCountry != nil && conf.ForbidCountry(country, product) {
			options = nil
			none = append(none, product)
		}
	}

	return
}

func (conf Config[P]) Product(country countries.Country, product P, weightGrams, netPrice int) (preview Preview) {
	if conf.ForbidCountry != nil && conf.ForbidCountry(country, product) {
		return
	}
	preview.ShippingDaysMin = math.MaxInt
	preview.ShippingNetMin = math.MaxInt
	for _, method := range conf.Methods {
		netPrice, minDays, maxDays, supported := method.Details(weightGrams, netPrice, country)
		if !supported {
			continue
		}
		if conf.ForbidDelivery != nil && conf.ForbidDelivery(country, method, product) {
			continue
		}

		if method.IsShipping {
			preview.Shipping = true
			preview.ShippingDaysMax = max(preview.ShippingDaysMax, maxDays)
			preview.ShippingDaysMin = min(preview.ShippingDaysMin, minDays)
			preview.ShippingNetMin = min(preview.ShippingNetMin, netPrice)
		} else {
			preview.Pickup = true
		}
	}
	return
}
