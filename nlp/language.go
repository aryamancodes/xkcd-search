// Functions to stem words to their roots, decode special chars
// and find synonyms for words
package nlp

import (
	"regexp"
	"strings"

	"github.com/kljensen/snowball/english"
	"github.com/mozillazg/go-unidecode"
)

func CleanContent(content string) string {
	cleaned := ""
	for _, word := range strings.Split(content, " ") {
		word = removeSpecialCharacters(word)
		cleaned += stem(word) + " "
	}
	return cleaned
}

func stem(word string) string {
	word = english.Stem(string(word), true)
	return word
}

func removeSpecialCharacters(word string) string {
	word = unidecode.Unidecode(word)
	word = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(word, "")
	return word
}
