package main

import (
	"bufio"
	"fmt"
	"os"
	"xkcd/db"
	"xkcd/nlp"
)

func main() {
	db.Connect()
	words := db.GetRawWords()
	model := nlp.TrainModel(words)
	for {
		fmt.Println("ENTER A QUERY:")
		var query string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			query = scanner.Text()
		}
		fmt.Printf("CORRECTED TO %s\n", model.SpellCheckSuggestions(query, 5))

	}
	// tf := index.ComputeAllTermFreq(comics)
	// index.ComputeAllComicFreq(comics, tf)
	// for {
	// 	fmt.Println("ENTER A QUERY:")
	// 	var query string

	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	if scanner.Scan() {
	// 		query = scanner.Text()
	// 	}

	// 	rankings := index.RankQuery(query, df)
	// 	sort.Slice(rankings, func(i, j int) bool {
	// 		return rankings[i].Rank >= rankings[j].Rank
	// 	})

	// 	if len(rankings) == 0 {
	// 		fmt.Printf("\n No results found. Maybe there isn't an xkcd for everything!\n")
	// 	}

	// 	for i, ranked := range rankings {
	// 		fmt.Printf("\n%d) %d (rank: %f) \n", i, ranked.ComicNum, ranked.Rank)
	// 	}
	// }
}
