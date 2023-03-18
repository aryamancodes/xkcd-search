// Functions to compute tf and idf for the data and rank queries via tf-idf

package index

import (
	"math"
	"regexp"
	"strings"

	"xkcd/model"
)

var termFreqChan = make(chan model.TermFreq, 250)

func computeTermFreq(terms []string, comic model.ExplainWikiJson) {
	termFreq := make(map[string]int)
	for _, term := range terms {
		termFreq[term]++
	}

	termFreqChan <- model.TermFreq{
		Comic:           comic,
		TermInComicFreq: termFreq,
		TotalTerms:      len(terms),
	}
}

// (concurrently) calculate the number of occurences of a term and the total number of terms in a comic, for all comics
func ComputeAllTermFreq(comics []model.ExplainWikiJson) []model.TermFreq {
	termFreqs := make([]model.TermFreq, 0, len(comics))
	for _, comic := range comics {
		terms := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(comic.Parse.Title+" "+comic.Parse.Wikitext.Content, " ")
		termsList := strings.Split(terms, " ")
		go computeTermFreq(termsList, comic)
	}

	for range comics {
		termFreqs = append(termFreqs, <-termFreqChan)
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

func RankQuery(query string, allTerms []model.TermFreq, allComics model.ComicFreq) []model.RankedComic {

	tf := func(queryTerm string, currComicTerms model.TermFreq) float64 {
		queryTermInCurrComic := currComicTerms.TermInComicFreq[queryTerm]
		totalTerms := currComicTerms.TotalTerms
		return float64(queryTermInCurrComic) / float64(totalTerms)
	}

	idf := func(queryTerm string) float64 {
		queryTermInAllComics := math.Max(float64(allComics.ComicsWithTermFreq[queryTerm]), 1)
		totalComics := float64(allComics.TotalComics)
		return math.Log10(totalComics / queryTermInAllComics)
	}

	rankings := make([]model.RankedComic, allComics.TotalComics)
	queryTerms := strings.Split(query, " ")

	for i := 0; i < allComics.TotalComics; i++ {
		rank := float64(0.0)
		for _, queryTerm := range queryTerms {
			rank += tf(queryTerm, allTerms[i]) + idf(queryTerm)
		}
		rankings = append(rankings, model.RankedComic{
			Comic: allTerms[i].Comic,
			Rank:  rank,
		})
	}
	return rankings
}
