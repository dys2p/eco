// Package lang provides utilities for translating web services.
package lang

import (
	"golang.org/x/text/message"
)

type Lang string

// Tr translates the given input text.
func (lang Lang) Tr(key message.Reference, a ...interface{}) string {
	return message.NewPrinter(message.MatchLanguage(string(lang))).Sprintf(key, a...)
}
