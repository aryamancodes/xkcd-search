// Functions to stem words to their roots, decode special chars and add autocorrect functionality
package nlp

import (
	"regexp"
	"strings"

	"github.com/kljensen/snowball/english"
	"github.com/mozillazg/go-unidecode"
)

// return a cleaned and stemmed version of content
func CleanAndStem(content string) (string, string) {
	cleaned := ""
	stemmed := ""
	for _, word := range strings.Fields(content) {
		word = removeSpecialCharacters(word)
		cleaned += word + " "
		stemmed += stem(word) + " "
	}
	return cleaned, stemmed
}

func stem(word string) string {
	word = english.Stem(string(word), true)
	return word
}

func removeSpecialCharacters(word string) string {
	word = unidecode.Unidecode(word)
	word = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(word, "")
	word = strings.ToLower(word)
	return word
}
