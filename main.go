package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"xkcd/db"
	"xkcd/index"
	"xkcd/nlp"

	"github.com/eiannone/keyboard"
)

func main() {
	db.Connect()
	db.GetComics()
	//comics := index.FetchAllComics()
	df := db.GetComicFreq()
	words := db.GetRawWords()
	model := nlp.TrainModel(words)

	fmt.Println("ENTER A QUERY:")

	char, key, err := keyboard.GetSingleKey()
	query := ""
	for key != keyboard.KeyEnter && err == nil {
		query += string(char)
		if len(query) >= 3 {
			auto, _ := model.Autocomplete(query)
			fmt.Printf("\n%s suggestions: %v\n", query, auto)
		} else {
			fmt.Print(string(char))
		}
		char, key, err = keyboard.GetSingleKey()
		if err != nil {
			log.Fatal(err)
		}
	}
	rawQuery, stemQuery := nlp.CleanAndStem(query)
	rankings := index.RankQuery(rawQuery, stemQuery, df)
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Rank >= rankings[j].Rank
	})

	if len(rankings) == 0 {
		fmt.Printf("\n No results found. Maybe there isn't an xkcd for everything!\n")
	}

	maxDisplayed := int(math.Min(float64(len(rankings)), 9))
	rankings = rankings[:maxDisplayed]
	for i, ranked := range rankings {
		fmt.Printf("%d) https://xkcd.com/%d (rank: %f) \n", i, ranked.ComicNum, ranked.Rank)
	}
}
