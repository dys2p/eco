package jsonld

import (
	"encoding/json"
	"testing"
)

type breadcrumb struct {
	name    string
	urlpath string
}

func (b breadcrumb) Name() string {
	return b.name
}

func (b breadcrumb) URLPath() string {
	return b.urlpath
}

func TestBreadcrumbList(t *testing.T) {
	got, _ := json.Marshal(BreadcrumbList("https://example.com", []Breadcrumb{
		breadcrumb{"Books", "/books"},
		breadcrumb{"Science Fiction", "/books/sciencefiction"},
		breadcrumb{"Award Winners", ""},
	}))
	want := `{"@context":"https://schema.org","@type":"BreadcrumbList","itemListElement":[{"@type":"ListItem","position":1,"name":"Books","item":"https://example.com/books"},{"@type":"ListItem","position":2,"name":"Science Fiction","item":"https://example.com/books/sciencefiction"},{"@type":"ListItem","position":3,"name":"Award Winners"}]}`

	if string(got) != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestBreadcrumbListEmpty(t *testing.T) {
	got, _ := json.Marshal(BreadcrumbList("https://example.com", []breadcrumb{}))

	if string(got) != "null" {
		t.Fatalf("got %s, want null", got)
	}
}
