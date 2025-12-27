package ssg_test

import (
	"os"

	"github.com/dys2p/eco/lang"
	"github.com/dys2p/eco/ssg"
)

//go:generate gotext-update-templates -srclang=en-US -lang=de-DE,en-US -out=catalog.go -d . -d ./example.com

func ExampleWebsite_WriteFiles() {
	langs := lang.MakeLanguages(nil, "de", "en")
	ws, err := ssg.MakeWebsite(os.DirFS("./example.com"), nil, langs, nil)
	if err != nil {
		panic(err)
	}
	err = ws.WriteFiles("/tmp/build/example.com")
	if err != nil {
		panic(err)
	}

	ssg.ListenAndServe("/tmp/build/example.com")
}
