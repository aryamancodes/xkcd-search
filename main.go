package main

import (
	"fmt"
	"os"
	"sort"
	"xkcd/index"
)

func main() {
	comics := index.FetchAllExplanations()

	tf := index.ComputeAllTermFreq(comics)

	df := index.ComputeAllComicFreq(comics, tf)
	fmt.Fprintf(os.Stderr, "\nTERM FREQ LOOKS LIKE: %+v\n", tf)
	fmt.Fprintf(os.Stderr, "\nDOC FREQ LOOKS LIKE: %+v\n", df)

	for {
		fmt.Println("ENTER A QUERY TERM:")
		var query string

		fmt.Scanln(&query)

		rankings := index.RankQuery(query, tf, df)
		fmt.Fprintf(os.Stderr, "c\nRANKINGS FOR %s ARE: %+v\n", query, rankings)
		sort.Slice(rankings, func(i, j int) bool {
			return rankings[i].Rank >= rankings[j].Rank
		})

		top10 := rankings[0:11]
		for i, ranked := range top10 {
			fmt.Printf("\n%d) %s (rank: %f) \n", i, ranked.Comic.Parse.Title, ranked.Rank)
		}
	}
}
