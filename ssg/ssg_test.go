package ssg_test

import (
	"os"

	"github.com/dys2p/eco/ssg"
)

//go:generate gotext-update-templates -srclang=en-US -lang=de-DE,en-US -out=catalog.go -d . -d ./example.com

func ExampleStaticHTML() {
	ssg.Must(ssg.MakeWebsite(os.DirFS("./example.com"), nil, "de", "en")).StaticHTML("/tmp/build/example.com")
	ssg.ListenAndServe("/tmp/build/example.com")
}
