# gotext-update-templates

Command gotext-update-templates extracts and merges translations and generates a catalog.

Unlike `gotext update`, it also extracts messages for translation from HTML templates. For that purpose it accepts an additional flag `trfunc`, which defaults to `Tr`. It extracts strings from pipelines `.Tr` and `$.Tr`.

## Example

`hello.html`:

```
<p>{{.Tr "Hello World"}}</p>
```

`main.go`:

```
package main

//go:generate gotext-update-templates -srclang=en-US -lang=de-DE,en-US -out=catalog.go .

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/dys2p/eco/lang"
)

//go:embed hello.html
var fs embed.FS

var hello = template.Must(template.ParseFS(fs, "hello.html"))

func main() {
	langs := lang.MakeLanguages("de", "en")
	for _, l := range langs {
		buf := &bytes.Buffer{}
		hello.Execute(buf, l)
		fmt.Println(buf.String())
	}
}
```

Run `go generate`, then build your program and run it:

```
$ go generate
2009/11/10 23:00:00 de-DE: Missing entry for "Hello World".
$ go build
$ ./example
<p>Hello World</p>
<p>Hello World</p>
```

Now copy `locales/de-DE/out.gotext.json` to `locales/de-DE/messages.gotext.json` and insert the translation. Do not edit `out.gotext.json`. Do not remove `messages.gotext.json` afterwards.

Run `go generate` again, build your program and run it again:

```
$ go generate
2009/11/10 23:00:00 de-DE: Missing entry for "Hello World".
$ go build
$ ./example
<p>Hello World</p>
<p>Hallo Welt</p>
```

## Crossing module boundaries

The Go translation framework is currently using a global variable `golang.org/x/text/message.DefaultCatalog` which makes it nearly impossible to use multiple catalogs. We can, at least for local modules, work around this by extracting messages from other modules and defining their translations in our code:

`gotext-update-templates -srclang=en-US -lang=en-US,de-DE -out=catalog.go . ./html ../other-module`

This will fail with the error: `directory ../other-module outside main module or its selected dependencies`.

To fix this, create a [Go workspace](https://go.dev/blog/get-familiar-with-workspaces) which contains your main module and the other module, because "each module within a workspace is treated as a main module when resolving dependencies". For example:

```
go work init .
go work use ../other-module
```
