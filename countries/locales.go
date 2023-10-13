package countries

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

//go:generate gotext-update-templates -srclang=en-US -lang=en-US,de-DE -out=catalog.go .

// languages should match the -lang arguments above
var matcher = language.NewMatcher([]language.Tag{
	language.English,
	language.MustParse("de-DE"), // does not work with language.German
})

// don't rely on message.DefaultCatalog which may have been overwritten
var cat = func() catalog.Catalog {
	dict := map[string]catalog.Dictionary{
		"de_DE": &dictionary{index: de_DEIndex, data: de_DEData},
		"en_US": &dictionary{index: en_USIndex, data: en_USData},
	}
	fallback := language.MustParse("en-US")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	return cat
}()
