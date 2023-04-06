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
	df := db.GetComicFreq()
	words := db.GetRawWords()
	model := nlp.TrainModel(words)
	fmt.Println("ENTER A QUERY:")
	// var query string
	// scanner := bufio.NewScanner(os.Stdin)
	// if scanner.Scan() {
	// 	query = scanner.Text()
	// }
	char, key, err := keyboard.GetSingleKey()
	query := ""
	currword := ""
	for key != keyboard.KeyEnter && key != keyboard.KeyEsc && err == nil && key != keyboard.KeyCtrlC {
		if key == keyboard.KeySpace {
			query += currword + " "
			fmt.Print(query)
			currword = ""
			char, key, err = keyboard.GetSingleKey()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			currword += string(char)
			if len(currword) >= 3 {
				auto, _ := model.Autocomplete(currword)
				fmt.Printf("\n%s suggestions: %v\n", query+currword, auto)
			} else {
				fmt.Print(string(char))
			}
			char, key, err = keyboard.GetSingleKey()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if key == keyboard.KeyEnter {
		query += currword
		rawQuery, stemQuery := nlp.CleanAndStem(query)
		hasTypo, autocorrectedRaw := nlp.Autocorect(model, rawQuery)
		if hasTypo {
			rawQuery = autocorrectedRaw
			rawQuery, stemQuery = nlp.CleanAndStem(rawQuery)
			fmt.Printf("SHOWING RESULTS FOR %s\n", rawQuery)

		}
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
}
