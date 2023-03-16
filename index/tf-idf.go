// Functions to compute tf and idf for the data and rank queries via tf-idf

package index

import (
	"regexp"
	"strings"

	"xkcd/model"
)

var termFreqChan = make(chan model.TermFreq, 250)

func computeTermFreq(terms []string) {
	termFreq := make(map[string]int)
	for _, term := range terms {
		termFreq[term]++
	}

	termFreqChan <- model.TermFreq{
		TermInComicFreq: termFreq,
		TotalTerms:      len(terms),
	}
}

// (concurrently) calculate the number of occurences of a term and the total number of terms in a comic, for all comics
func ComputeAllTermFreq(comics []model.ExplainWikiJson) []model.TermFreq {
	termFreqs := make([]model.TermFreq, len(comics))
	for _, comic := range comics {
		terms := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(comic.Parse.Title+" "+comic.Parse.Wikitext.Content, " ")
		termsList := strings.Split(terms, " ")
		go computeTermFreq(termsList)
	}

	for i := range comics {
		termFreqs[i] = <-termFreqChan
	}

	return termFreqs
}

// calculate the number of comics each term occurs in and the total number of comics
func ComputeAllComicFreq(comics []model.ExplainWikiJson, termFreqs []model.TermFreq) model.ComicFreq {
	comicFreq := make(map[string]int, len(comics))
	for _, termFreq := range termFreqs {
		for term := range termFreq.TermInComicFreq {
			comicFreq[term]++
		}
	}

	return model.ComicFreq{
		ComicsWithTermFreq: comicFreq,
		TotalComics:        len(comics),
	}

}

func RankQuery() {

}
