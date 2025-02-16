package jsonld

type ItemList struct {
	Context  string            `json:"@context"`
	Type     string            `json:"@type"`
	Elements []ItemListElement `json:"itemListElement"`
}

type ItemListElement struct {
	Type     string `json:"@type"`
	Position int    `json:"position"` // starting with 1
	Name     string `json:"name"`
	Item     string `json:"item,omitempty"` // URL
}

type Breadcrumb interface {
	Name() string
	URLPath() string
}

// See https://developers.google.com/search/docs/appearance/structured-data/breadcrumb
//
// BreadcrumbList uses generics so we don't have to convert the slice from []SomeType to []Breadcrumb.
func BreadcrumbList[T Breadcrumb](urlprefix string, breadcrumbs []T) *ItemList {
	var elements []ItemListElement
	for i, breadcrumb := range breadcrumbs {
		var url string
		if path := breadcrumb.URLPath(); path != "" {
			url = urlprefix + path
		}

		elements = append(elements, ItemListElement{
			Type:     "ListItem",
			Position: i + 1,
			Name:     breadcrumb.Name(),
			Item:     url,
		})
	}
	if len(elements) == 0 {
		return nil
	}
	return &ItemList{
		Context:  "https://schema.org",
		Type:     "BreadcrumbList",
		Elements: elements,
	}
}
