# gotext-update-templates

Command gotext-update-templates merge translations and generates a catalog.

Unlike gotext update, it also extracts messages for translation from HTML templates. For that purpose it accepts an additional flag `trfunc`, which defaults to `Tr`. It extracts strings from pipelines `.Tr` and `$.Tr`.

Templates are read recursively from the working directory. If you use go generate, note that "the generator is run in the package's source directory".

## Example

`hello.html`:

```
<p>{{.Tr "Hello World"}}</p>
```

`main.go`:

```
package main

//go:generate gotext-update-templates -srclang=en-US -lang=en-US,de-DE -out=catalog.go .

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

type helloData struct {
	lang.Lang
}

func main() {
	for _, l := range []string{"en", "de"} {
		buf := &bytes.Buffer{}
		_ = hello.Execute(buf, helloData{
			Lang: lang.Lang(l),
		})
		fmt.Println(buf.String())
	}
}
```

Run `go generate` (or run `gotext-update-templates` directly), then build your program and run it:

```
$ go generate
2009/11/10 23:00:00 de-DE: Missing entry for "Hello World".
$ go build
$ ./example
<p>Hello World</p>
<p>Hello World</p>
```

Now copy `locales/de-DE/out.gotext.json` to `locales/de-DE/messages.gotext.json` and insert the translation. Do not edit `out.gotext.json`. Do not remove `messages.gotext.json` afterwards.

Run `go generate`, build your program and run it again:

```
$ go generate
2009/11/10 23:00:00 de-DE: Missing entry for "Hello World".
$ go build
$ ./example
<p>Hello World</p>
<p>Hallo Welt</p>
```
