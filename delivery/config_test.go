package delivery

import (
	"slices"
	"testing"

	"github.com/dys2p/eco/countries"
)

func optionIDs(options []MethodOption) []string {
	var ids []string
	for _, option := range options {
		ids = append(ids, option.Method.ID)
	}
	slices.Sort(ids)
	return ids
}

func TestOptions(t *testing.T) {
	type Product string
	var (
		boringProduct     = Product("boring")
		dangerousProduct  = Product("dangerous")
		restrictedProduct = Product("restricted")
	)

	courier := Method{
		ID:   "courier",
		Name: "My Courier",
		Details: func(weight, netValue int, country countries.Country) (netPrice, minDays, maxDays int, supported bool) {
			netPrice = 599
			if weight > 20000 {
				netPrice = 1499
			}
			return netPrice, 1, 3, slices.Contains(countries.EuropeanUnion, country)
		},
		IsShipping: true,
	}
	store := Method{
		ID:   "store",
		Name: "Pick up at Store",
		Details: func(weight, netValue int, country countries.Country) (netPrice, minDays, maxDays int, supported bool) {
			return 0, 0, 0, country == countries.DE
		},
	}
	conf := Config[Product]{
		Methods: []Method{courier, store},
		ForbidCountry: func(country countries.Country, product Product) bool {
			return country != countries.DE && product == restrictedProduct
		},
		ForbidDelivery: func(country countries.Country, method Method, product Product) bool {
			return method.ID == "courier" && product == dangerousProduct
		},
	}

	tests := []struct {
		selected *Method
		country  countries.Country
		products []Product

		wantOptionIDs  []string
		wantNone       []Product
		wantPickupOnly []Product
	}{
		// DE: boring product okay
		{nil, countries.DE, []Product{boringProduct}, []string{"courier", "store"}, nil, nil},

		// FR: boring product okay
		{nil, countries.FR, []Product{boringProduct}, []string{"courier"}, nil, nil},

		// DE: restricted product okay
		{nil, countries.DE, []Product{boringProduct, restrictedProduct}, []string{"courier", "store"}, nil, nil},

		// FR: no restricted product
		{nil, countries.FR, []Product{boringProduct, restrictedProduct}, nil, []Product{restrictedProduct}, nil},

		// DE: dangerous product pickup only
		{nil, countries.DE, []Product{boringProduct, dangerousProduct, restrictedProduct}, []string{"store"}, nil, []Product{dangerousProduct}},
		{&courier, countries.DE, []Product{boringProduct, dangerousProduct, restrictedProduct}, []string{"store"}, nil, []Product{dangerousProduct}}, // same but with courier selected by user
		{&store, countries.DE, []Product{boringProduct, dangerousProduct, restrictedProduct}, []string{"store"}, nil, []Product{dangerousProduct}},   // same but with store selected by user

		// FR: no restricted product, no dangerous product because there is no pickup point in that country
		{nil, countries.FR, []Product{boringProduct, dangerousProduct, restrictedProduct}, nil, []Product{dangerousProduct, restrictedProduct}, nil},
	}

	for i, test := range tests {
		gotOptions, gotNone, gotPickupOnly := conf.Options(test.selected, 1000, 2500, test.country, slices.Values(test.products))
		gotOptionIDs := optionIDs(gotOptions)
		if !slices.Equal(gotOptionIDs, test.wantOptionIDs) {
			t.Fatalf("test %d: options: got %v, want %v", i, gotOptionIDs, test.wantOptionIDs)
		}
		if !slices.Equal(gotNone, test.wantNone) {
			t.Fatalf("test %d: none: got %v, want %v", i, gotNone, test.wantNone)
		}
		if !slices.Equal(gotPickupOnly, test.wantPickupOnly) {
			t.Fatalf("test %d: pickup only: got %v, want %v", i, gotPickupOnly, test.wantPickupOnly)
		}
	}
}
