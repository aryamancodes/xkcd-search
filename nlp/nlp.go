// Functions to stem words to their roots, decode special chars and add autocorrect functionality
package nlp

import (
	"log"
	"regexp"
	"strings"

	"github.com/kljensen/snowball/english"
	"github.com/mozillazg/go-unidecode"
	"github.com/sajari/fuzzy"
)

var language *fuzzy.Model

// return a cleaned and stemmed version of content
func CleanAndStem(content string) (string, string) {
	cleaned := ""
	stemmed := ""
	words := strings.Fields(content)
	for i, word := range words {
		word = removeSpecialCharacters(word)
		cleaned += word
		stemmed += stem(word)
		if i != len(words)-1 {
			cleaned += " "
			stemmed += " "
		}
	}
	return cleaned, stemmed
}

func stem(word string) string {
	word = english.Stem(string(word), true)
	return word
}

func removeSpecialCharacters(word string) string {
	word = unidecode.Unidecode(word)
	word = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(word, "")
	word = strings.ToLower(word)
	return word
}

func TrainModel(words []string) {
	language = fuzzy.NewModel()
	language.SetThreshold(5)
	language.SetDepth(3)
	language.Train(words)
}

func Autocorect(raw string) (bool, string) {
	corrected := ""
	rawWords := strings.Fields(raw)
	for i, term := range rawWords {
		corrected += language.SpellCheck(term)
		if i != len(rawWords)-1 {
			corrected += " "
		}
	}
	return raw != corrected, corrected
}

func Autocomplete(term string) []string {
	terms, err := language.Autocomplete(term)
	if err != nil {
		log.Fatal(err)
	}
	return terms
}
