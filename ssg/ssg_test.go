package ssg_test

import (
	"os"

	"github.com/dys2p/eco/lang"
	"github.com/dys2p/eco/ssg"
)

//go:generate gotext-update-templates -srclang=en-US -lang=de-DE,en-US -out=catalog.go -d . -d ./example.com

func ExampleWebsite_StaticHTML() {
	langs := lang.MakeLanguages(nil, "de", "en")
	ws := ssg.Must(ssg.MakeWebsite(os.DirFS("./example.com"), nil, langs))
	ws.StaticHTML("/tmp/build/example.com", false)

	ssg.ListenAndServe("/tmp/build/example.com")
}
