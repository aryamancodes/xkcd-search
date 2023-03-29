package main

import (
	"xkcd/db"
	"xkcd/index"
)

func main() {
	db.Connect()
	comics := index.FetchAllComics()
	tf := index.ComputeAllTermFreq(comics)
	index.ComputeAllComicFreq(comics, tf)

	// for {
	// 	fmt.Println("ENTER A QUERY:")
	// 	var query string

	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	if scanner.Scan() {
	// 		query = scanner.Text()
	// 	}

	// 	rankings := index.RankQuery(query, tf, df)
	// 	fmt.Fprintf(os.Stderr, "c\nRANKINGS FOR %s ARE: %+v\n", query, rankings)
	// 	sort.Slice(rankings, func(i, j int) bool {
	// 		return rankings[i].Rank >= rankings[j].Rank
	// 	})

	// 	if len(rankings) == 0 {
	// 		fmt.Printf("\n No results found. Maybe there isn't an xkcd for everything!")
	// 	}

	// 	for i, ranked := range rankings {
	// 		fmt.Printf("\n%d) %s (rank: %f) \n", i, ranked.Comic.Transcript, ranked.Rank)
	// 	}
	//}
}
