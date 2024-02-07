// Command gotext-update-templates extracts and merges translations and generates a catalog.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template/parse"

	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message/pipeline"
	"golang.org/x/tools/go/packages"
)

type Config struct {
	Dir               string
	Lang              string
	Out               string
	Packages          []string
	SrcLang           string
	TranslateFuncName string
	Verbose           bool
}

func main() {
	// own FlagSet because the global one is already in use
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	dir := fs.String("dir", "locales", "default subdirectory to store translation files")
	lang := fs.String("lang", "en-US", "comma-separated list of languages to process")
	out := fs.String("out", "catalog.go", "output file to write to")
	srcLang := fs.String("srclang", "en-US", "the source-code language")
	trFunc := fs.String("trfunc", "Tr", "name of translate method which is used in templates")
	verbose := fs.Bool("v", false, "output list of processed template files")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <package>* [flags]\n", os.Args[0])
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[1:])

	config := Config{
		Dir:               *dir,
		Lang:              *lang,
		Out:               *out,
		Packages:          fs.Args(),
		SrcLang:           *srcLang,
		TranslateFuncName: *trFunc,
		Verbose:           *verbose,
	}
	if err := config.Run(); err != nil {
		log.Fatalln(err)
	}
}

func (config Config) Run() error {
	log.SetFlags(0)

	var templateMessages = []pipeline.Message{}
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedFiles | packages.NeedEmbedFiles}, config.Packages...)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return errors.New(pkg.Errors[0].Msg)
		}
		for _, embedFile := range pkg.EmbedFiles {
			if ext := filepath.Ext(embedFile); ext == ".html" || ext == ".txt" {
				if config.Verbose {
					log.Println(embedFile)
				}
				file, err := os.ReadFile(embedFile)
				if err != nil {
					return err
				}
				// similar to parse.Parse but with SkipFuncCheck
				trees := make(map[string]*parse.Tree)
				t := parse.New("name")
				t.Mode |= parse.SkipFuncCheck
				if _, err := t.Parse(string(file), "", "", trees); err != nil {
					return err
				}
				for _, tree := range trees {
					config.processNode(&templateMessages, tree.Root)
				}
			}
		}
	}

	supported := []language.Tag{}
	for _, l := range strings.FieldsFunc(config.Lang, func(r rune) bool { return r == ',' }) {
		supported = append(supported, language.Make(l))
	}

	pconf := &pipeline.Config{
		Supported:      supported,
		SourceLanguage: language.Make(config.SrcLang),
		Packages:       config.Packages,
		Dir:            config.Dir,
		GenFile:        config.Out,
	}

	// see https://cs.opensource.google/go/x/text/+/master:cmd/gotext/update.go
	state, err := pipeline.Extract(pconf)
	if err != nil {
		return err
	}
	state.Extracted.Messages = append(state.Extracted.Messages, templateMessages...)
	if err := state.Import(); err != nil {
		return err
	}
	if err := state.Merge(); err != nil {
		return err
	}
	if err := state.Export(); err != nil {
		return err
	}
	if err := state.Generate(); err != nil {
		return err
	}
	return nil
}

func (config Config) processNode(templateMessages *[]pipeline.Message, node parse.Node) {
	if node.Type() == parse.NodeList {
		if listNode, ok := node.(*parse.ListNode); ok {
			for _, childNode := range listNode.Nodes {
				config.processNode(templateMessages, childNode)
			}
		}
	}
	if node.Type() == parse.NodeIf {
		if ifNode, ok := node.(*parse.IfNode); ok {
			config.processNode(templateMessages, ifNode.List)
			if ifNode.ElseList != nil {
				config.processNode(templateMessages, ifNode.ElseList)
			}
		}
	}
	if node.Type() == parse.NodeWith {
		if withNode, ok := node.(*parse.WithNode); ok {
			config.processNode(templateMessages, withNode.List)
			if withNode.ElseList != nil {
				config.processNode(templateMessages, withNode.ElseList)
			}
		}
	}
	if node.Type() == parse.NodeRange {
		if rangeNode, ok := node.(*parse.RangeNode); ok {
			config.processNode(templateMessages, rangeNode.List)
			if rangeNode.ElseList != nil {
				config.processNode(templateMessages, rangeNode.ElseList)
			}
		}
	}
	if node.Type() == parse.NodeAction {
		if actionNode, ok := node.(*parse.ActionNode); ok {
			for _, cmd := range actionNode.Pipe.Cmds {
				if !containsIdentifier(cmd, config.TranslateFuncName) {
					continue
				}
				for _, arg := range cmd.Args {
					if arg.Type() == parse.NodeString {
						if stringNode, ok := arg.(*parse.StringNode); ok {
							text := stringNode.Text
							message := pipeline.Message{
								ID:  pipeline.IDList{text},
								Key: text,
								Message: pipeline.Text{
									Msg: text,
								},
							}
							*templateMessages = append(*templateMessages, message)
						}
					}
				}
			}
		}
	}
}

func containsIdentifier(cmd *parse.CommandNode, identifier string) bool {
	if len(cmd.Args) == 0 {
		return false
	}
	arg := cmd.Args[0]
	var identifiers []string
	switch arg.Type() {
	case parse.NodeField:
		identifiers = arg.(*parse.FieldNode).Ident
	case parse.NodeVariable:
		identifiers = arg.(*parse.VariableNode).Ident
	}
	return slices.Contains(identifiers, identifier)
}
