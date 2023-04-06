// Functions to compute tf and idf for the data and rank queries via tf-idf

package index

import (
	"fmt"
	"math"
	"strings"

	"xkcd/db"
	"xkcd/model"
)

var termFreqChan = make(chan model.TermFreq, 250)

func computeTermFreq(comic model.Comic, comicNum int) {
	termFreq := make(map[string]int)

	//weight terms in the following order: title > alt > transcript > explain
	for _, titleTerm := range strings.Fields(comic.Title) {
		termFreq[titleTerm] += 4
	}
	for _, altTerm := range strings.Fields(comic.AltText) {
		termFreq[altTerm] += 3
	}
	for _, transcriptTerm := range strings.Fields(comic.Transcript) {
		termFreq[transcriptTerm] += 2
	}
	for _, explainTerm := range strings.Fields(comic.Explanation) {
		termFreq[explainTerm]++
	}

	termFreqChan <- model.TermFreq{
		ComicNum:        comicNum,
		TermInComicFreq: termFreq,
		TotalTerms:      len(termFreq),
	}
}

// (concurrently) calculate the number of occurences of a term and the total number of terms in a comic, for all comics
func ComputeAllTermFreq(comics []model.Comic) []model.TermFreq {
	termFreqs := make([]model.TermFreq, 0)
	for _, comic := range comics {
		go computeTermFreq(comic, comic.Num)
	}

	for range comics {
		termFreqs = append(termFreqs, <-termFreqChan)
	}
	fmt.Println("FETCHED ALL TFS")
	db.BatchStoreTermFreq(termFreqs)
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
	fmt.Println("FETCHED ALL DFS")
	db.BatchStoreComicFreq(result)
	return result
}

func tf(queryTerm string, currComicTerms model.TermFreq) float64 {
	queryTermInCurrComic := float64(currComicTerms.TermInComicFreq[queryTerm])
	if queryTermInCurrComic == 0 {
		return 0
	}

	// ln normalised tf to favour distinct query terms matching fewer times rather than the same query terms matching many times
	// ref: https://ecommons.cornell.edu/bitstream/handle/1813/7281/97-1626.pdf?sequence=1 (page 8)
	return 1 + math.Log(queryTermInCurrComic)
}

func idf(queryTerm string, allComics model.ComicFreq) float64 {
	comicsWithQueryTerm := math.Max(float64(allComics.ComicsWithTermFreq[queryTerm]), 1)
	totalComics := float64(allComics.TotalComics)
	return math.Log10(totalComics / comicsWithQueryTerm)
}

func RankQuery(query string, allComics model.ComicFreq) []model.RankedComic {
	rankings := make([]model.RankedComic, 0)
	queryTerms := strings.Fields(query)
	// fetch only the tf of comics that contain the query terms
	// ie. map of [comic (containing atleast one query term)] -> termFreq of comic
	queryTermFreq := db.GetTermFreq(queryTerms)

	for i := 0; i < allComics.TotalComics; i++ {
		rank := 0.0
		for _, queryTerm := range queryTerms {
			rank += tf(queryTerm, queryTermFreq[i]) * idf(queryTerm, allComics)
		}

		//only return comics whose rank isn't 0
		if rank > 0 {
			rankings = append(rankings, model.RankedComic{
				ComicNum: i,
				Rank:     rank,
			})
		}
	}
	return rankings
}
