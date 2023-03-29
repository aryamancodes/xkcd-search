package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"xkcd/db"
	"xkcd/index"
)

func main() {
	db.Connect()
	db.GetComics()
	df := db.GetComicFreq()

	for {
		fmt.Println("ENTER A QUERY:")
		var query string

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			query = scanner.Text()
		}

		rankings := index.RankQuery(query, df)
		sort.Slice(rankings, func(i, j int) bool {
			return rankings[i].Rank >= rankings[j].Rank
		})

		if len(rankings) == 0 {
			fmt.Printf("\n No results found. Maybe there isn't an xkcd for everything!\n")
		}

		for i, ranked := range rankings {
			fmt.Printf("\n%d) %d (rank: %f) \n", i, ranked.ComicNum, ranked.Rank)
		}
	}
}
