// Command gotext-update-templates is like gotext, but it extracts messages for translation from HTML templates.
// It reads from the working directory. If you use go generate, note that "the generator is run in the package's source directory".
package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template/parse"

	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message/pipeline"
)

type Config struct {
	Lang              string
	Out               string
	SrcLang           string
	TranslateFuncName string
}

func main() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	lang := flagSet.String("lang", "de-DE,en-US", "")
	out := flagSet.String("out", "catalog.go", "")
	srcLang := flagSet.String("srclang", "de-DE", "")
	trFunc := flagSet.String("trfunc", "Tr", "")
	flagSet.Parse(os.Args[1:])

	config := Config{
		Lang:              *lang,
		Out:               *out,
		SrcLang:           *srcLang,
		TranslateFuncName: *trFunc,
	}
	if err := config.Run(); err != nil {
		log.Fatalln(err)
	}
}

func (config Config) Run() error {
	var messages = []pipeline.Message{}
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if ext := filepath.Ext(info.Name()); ext == ".html" || ext == ".txt" {
			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// similar to parse.Parse but wirh SkipFuncCheck
			trees := make(map[string]*parse.Tree)
			t := parse.New("name")
			t.Mode |= parse.SkipFuncCheck
			if _, err := t.Parse(string(file), "", "", trees); err != nil {
				return err
			}
			// nodes are in linear order, not nested
			for _, tree := range trees {
				for _, node := range tree.Root.Nodes {
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
											messages = append(messages, message)
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	supported := []language.Tag{}
	for _, l := range strings.Split(config.Lang, ",") {
		supported = append(supported, language.Make(l))
	}

	state := pipeline.State{
		Extracted: pipeline.Messages{
			Language: language.Make(config.SrcLang),
			Messages: messages,
		},
		Config: pipeline.Config{
			Supported:      supported,
			SourceLanguage: language.Make(config.SrcLang),
			GenFile:        config.Out,
		},
		Translations: nil,
	}
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
