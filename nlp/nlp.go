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

func TrainModel(words []string) *fuzzy.Model {
	model := fuzzy.NewModel()
	model.SetThreshold(5)
	model.SetDepth(3)
	model.Train(words)
	return model
}

func Autocorect(model *fuzzy.Model, raw string) (bool, string) {
	corrected := ""
	for _, term := range strings.Fields(raw) {
		corrected += model.SpellCheck(term) + " "
	}
	return raw != corrected, corrected
}

func Autocomplete(model *fuzzy.Model, term string) []string {
	terms, err := model.Autocomplete(term)
	if err != nil {
		log.Fatal(err)
	}
	return terms
}
