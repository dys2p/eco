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

// gets method and checks Details, ForbidCountry and ForbidDelivery
func (conf Config[P]) Checkout(methodID string, weightGrams, goodsNetPrice int, country countries.Country, products iter.Seq[P]) (method *Method, deliveryNetPrice, minDays, maxDays int) {
	method = conf.Method(methodID)
	if method == nil {
		return nil, 0, 0, 0
	}
	deliveryNetPrice, minDays, maxDays, supported := method.Details(weightGrams, goodsNetPrice, country)
	if !supported {
		return nil, 0, 0, 0
	}
	for product := range products {
		if conf.ForbidCountry != nil && conf.ForbidCountry(country, product) {
			return nil, 0, 0, 0
		}
		if conf.ForbidDelivery != nil && conf.ForbidDelivery(country, *method, product) {
			return nil, 0, 0, 0
		}
	}
	return method, deliveryNetPrice, minDays, maxDays
}

// for cart view, does not set MethodOption.GrossPrice
func (conf Config[P]) Options(selected *Method, weightGrams, goodsNetPrice int, country countries.Country, products iter.Seq[P]) (options []MethodOption, none []P, pickupOnly []P) {
	var selectedID string
	if selected != nil {
		selectedID = selected.ID
	}

	// check Details
	for _, method := range conf.Methods {
		deliveryNetPrice, minDays, maxDays, supported := method.Details(weightGrams, goodsNetPrice, country)
		if supported {
			options = append(options, MethodOption{
				Method:          method,
				NetPrice:        deliveryNetPrice,
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

// for product view
func (conf Config[P]) Product(country countries.Country, product P, weightGrams, goodsNetPrice int) (preview Preview) {
	if conf.ForbidCountry != nil && conf.ForbidCountry(country, product) {
		return
	}
	preview.ShippingDaysMin = math.MaxInt
	preview.ShippingNetMin = math.MaxInt
	for _, method := range conf.Methods {
		deliveryNetPrice, minDays, maxDays, supported := method.Details(weightGrams, goodsNetPrice, country)
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
			preview.ShippingNetMin = min(preview.ShippingNetMin, deliveryNetPrice)
		} else {
			preview.Pickup = true
		}
	}
	return
}
