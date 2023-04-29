// Functions to compute tf and idf for the data and rank queries via tf-idf

package index

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"xkcd/db"
	"xkcd/model"
)

var termFreqChan = make(chan model.TermFreq, 250)

func computeTermFreq(comic model.Comic, comicNum int) {
	termFreq := make(map[string]int)

	//map each stem term to a string of raw terms for title, alt and transcript sections only.
	//This uses the fact that raw and stemmed sections lengths are always the same
	stemToRawMap := make(map[string]string)
	var currRawSection []string

	//weight terms in the following order: title > alt > transcript > explain for both raw and stem maps
	currRawSection = strings.Fields(comic.TitleRaw)
	for i, titleTerm := range strings.Fields(comic.Title) {
		termFreq[titleTerm] += 4
		stemToRawMap[titleTerm] += " " + currRawSection[i] + " " + currRawSection[i] + " " + currRawSection[i]
	}

	currRawSection = strings.Fields(comic.AltTextRaw)
	for i, altTerm := range strings.Fields(comic.AltText) {
		termFreq[altTerm] += 3
		stemToRawMap[altTerm] += " " + currRawSection[i] + " " + currRawSection[i]
	}

	currRawSection = strings.Fields(comic.TranscriptRaw)
	for i, transcriptTerm := range strings.Fields(comic.Transcript) {
		termFreq[transcriptTerm] += 2
		stemToRawMap[transcriptTerm] += " " + currRawSection[i]
	}

	for _, explainTerm := range strings.Fields(comic.Explanation) {
		termFreq[explainTerm]++
	}

	termFreqChan <- model.TermFreq{
		ComicNum:        comicNum,
		TermInComicFreq: termFreq,
		StemToRawMap:    stemToRawMap,
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

// compare the term freqs of previously stored incomplete comics and any new revisions
// store the new tf and return just the new terms added to the comic
func ComputeIncompleteTermFreq(old []model.Comic, curr []model.Comic) []model.TermFreq {

	computeTf := func(comics []model.Comic) []model.TermFreq {
		termFreqs := make([]model.TermFreq, 0)
		for _, comic := range comics {
			go computeTermFreq(comic, comic.Num)
		}

		for range comics {
			termFreqs = append(termFreqs, <-termFreqChan)
		}
		return termFreqs
	}
	//compute and sort the tfs for old and new for comparision of new terms. store the new tf
	oldTfList := computeTf(old)
	currTfList := computeTf(curr)
	db.UpdateTermFreq(currTfList)

	sort.Slice(oldTfList, func(i, j int) bool {
		return oldTfList[i].ComicNum >= oldTfList[j].ComicNum
	})
	sort.Slice(currTfList, func(i, j int) bool {
		return currTfList[i].ComicNum >= currTfList[j].ComicNum
	})

	var revisedTfList []model.TermFreq
	for i, currTf := range currTfList {
		revisedTf := model.TermFreq{TermInComicFreq: make(map[string]int)}
		revisedTf.ComicNum = currTf.ComicNum
		revisedTf.StemToRawMap = currTf.StemToRawMap
		for currTerm, currFreq := range currTf.TermInComicFreq {
			if oldTfList[i].TermInComicFreq[currTerm] == 0 {
				revisedTf.TermInComicFreq[currTerm] = currFreq
			}
		}
		revisedTfList = append(revisedTfList, revisedTf)
	}

	return revisedTfList
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

func ComputeNewComicFreq(newTermFreqs []model.TermFreq, oldComicFreq model.ComicFreq) model.ComicFreq {
	newComicFreq := oldComicFreq
	for _, newTermFreq := range newTermFreqs {
		for newTerm := range newTermFreq.TermInComicFreq {
			newComicFreq.ComicsWithTermFreq[newTerm]++
		}
	}

	db.UpdateComicFreq(newComicFreq)
	return newComicFreq
}

func tf(stemTerm string, rawTerm string, currComicTerms model.TermFreq) float64 {
	exactMatches := strings.Count(currComicTerms.StemToRawMap[stemTerm], rawTerm)
	queryTermInCurrComic := float64(currComicTerms.TermInComicFreq[stemTerm] * (1 + exactMatches))
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

func findMatchSections(queryTerm []string, rawQueryTerms []string, comic model.Comic) []model.TermSection {
	var sections []model.TermSection
	for i, term := range queryTerm {
		var currTerm model.TermSection
		currTerm.Term = rawQueryTerms[i]
		currTerm.TitleCount = strings.Count(comic.Title, term)
		currTerm.AltCount = strings.Count(comic.AltText, term)
		currTerm.TranscriptCount = strings.Count(comic.Transcript, term)
		currTerm.ExplanationCount = strings.Count(comic.Explanation, term)
		sections = append(sections, currTerm)
	}
	return sections
}

func RankQuery(rawQuery string, stemQuery string, allComics []model.Comic, comicFreq model.ComicFreq) []model.RankedComic {
	rankings := make([]model.RankedComic, 0)
	stemQueryTerms := strings.Fields(stemQuery)
	// fetch only the tf of comics that contain the query terms
	// ie. map of [comic num (containing atleast one query term)] -> termFreq of comic
	queryTermFreq := db.GetTermFreq(stemQueryTerms)
	rawQueryTerms := strings.Fields(rawQuery)

	for i, currTermFreq := range queryTermFreq {
		rank := 0.0
		for j, stemTerm := range stemQueryTerms {
			rank += tf(stemTerm, rawQueryTerms[j], currTermFreq) * idf(stemTerm, comicFreq)
		}

		//only return comics whose rank isn't 0
		if rank > 0 {
			rankings = append(rankings, model.RankedComic{
				ComicNum:     i,
				Incomplete:   allComics[i-1].Incomplete,
				Interactive:  allComics[i-1].Interactive,
				Rank:         rank,
				TermSections: findMatchSections(stemQueryTerms, rawQueryTerms, allComics[currTermFreq.ComicNum-1]),
			})
		}
	}
	return rankings
}
