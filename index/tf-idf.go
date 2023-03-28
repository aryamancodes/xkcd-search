// Functions to compute tf and idf for the data and rank queries via tf-idf

package index

import (
	"math"
	"regexp"
	"strings"

	"xkcd/model"
)

var termFreqChan = make(chan model.TermFreq, 250)

func computeTermFreq(terms []string, comic model.Comic) {
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
func ComputeAllTermFreq(comics []model.Comic) []model.TermFreq {
	termFreqs := make([]model.TermFreq, 0)
	for _, comic := range comics {
		terms := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(comic.Transcript, " ")
		termsList := strings.Split(terms, " ")
		go computeTermFreq(termsList, comic)
	}

	for range comics {
		termFreqs = append(termFreqs, <-termFreqChan)
	}
	//db.BatchStoreTermFreq(termFreqs)
	return termFreqs
}

// calculate the number of comics each term occurs in and the total number of comics
func ComputeAllComicFreq(comics []model.Comic, termFreqs []model.TermFreq) model.ComicFreq {
	comicFreq := make(map[string]int)
	for _, termFreq := range termFreqs {
		for term := range termFreq.TermInComicFreq {
			comicFreq[term]++
		}
	}

	result := model.ComicFreq{
		ComicsWithTermFreq: comicFreq,
		TotalComics:        len(comics),
	}
	//db.BatchStoreComicFreq(result)
	return result
}

func RankQuery(query string, allTerms []model.TermFreq, allComics model.ComicFreq) []model.RankedComic {
	tf := func(queryTerm string, currComicTerms model.TermFreq) float64 {
		queryTermInCurrComic := float64(currComicTerms.TermInComicFreq[queryTerm])
		if queryTermInCurrComic == 0 {
			return 0
		}

		// ln normalised tf to favour distinct query terms matching fewer times rather than the same query terms matching many times
		// ref: https://ecommons.cornell.edu/bitstream/handle/1813/7281/97-1626.pdf?sequence=1 (page 8)
		return 1 + math.Log(queryTermInCurrComic)
	}

	idf := func(queryTerm string) float64 {
		comicsWithQueryTerm := math.Max(float64(allComics.ComicsWithTermFreq[queryTerm]), 1)
		totalComics := float64(allComics.TotalComics)
		return math.Log10(totalComics / comicsWithQueryTerm)
	}

	rankings := make([]model.RankedComic, 0)
	queryTerms := strings.Split(query, " ")

	for i := 0; i < allComics.TotalComics; i++ {
		rank := 0.0
		for _, queryTerm := range queryTerms {
			rank += tf(queryTerm, allTerms[i]) * idf(queryTerm)
		}

		//only return comics whose rank isn't 0
		if rank > 0 {
			rankings = append(rankings, model.RankedComic{
				Comic: allTerms[i].Comic,
				Rank:  rank,
			})
		}
	}
	return rankings
}
