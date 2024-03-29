// Functions to fetch all the data to index -- including the latest comic number and all explanations (concurrently) upto the latest

package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"xkcd/db"
	"xkcd/model"
)

const CURR_COMIC_URL = "https://xkcd.com/info.0.json"
const EXPLAIN_URL = "https://www.explainxkcd.com/wiki/api.php?action=parse&page=%d&prop=wikitext&format=json&redirects=1&origin=*"

var comicChan = make(chan model.Comic, 250)

// get the latest comic's number directly from the xkcd api
func GetCurrentComicNum() int {
	resp, err := http.Get(CURR_COMIC_URL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var currentComic model.CurrentComicJson
	json.Unmarshal([]byte(body), &currentComic)
	return currentComic.Number
}

func fetchComic(num int) {
	explainURL := fmt.Sprintf(EXPLAIN_URL, num)
	var fetchedExplainWiki model.ExplainWikiJson
	for {
		//try to fetch explanation from wiki
		resp, err := http.Get(explainURL)
		if err != nil {
			log.Fatalln(err)
		}
		//concurrent fetches may lead to >500 cloudflare errors
		if resp.StatusCode != http.StatusOK {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
			continue
		} else {
			//try to unmarshal response into json
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
			json.Unmarshal([]byte(body), &fetchedExplainWiki)
			//even if we get status code 200, an internal wiki error might have occured -- needs refetching
			if fetchedExplainWiki.Parse.Title != "" {
				break
			}
		}
	}
	comicChan <- Parse(fetchedExplainWiki)
}

// (concurrently) fetch all comics+explanations based on the current comic number
func FetchAllComics() []model.Comic {
	latestComicNumber := GetCurrentComicNum()
	comicList := make([]model.Comic, 0, latestComicNumber)

	for i := 0; i < latestComicNumber; i++ {
		go fetchComic(i + 1)
	}

	for i := 0; i < latestComicNumber; i++ {
		comicList = append(comicList, <-comicChan)
	}

	db.BatchStoreComics(comicList)
	return comicList
}

// fetch new comics when updating
func FetchNewComics() []model.Comic {
	lastStoredComic := db.GetLastStoredComicNum()
	latestComicNumber := GetCurrentComicNum()
	if lastStoredComic == latestComicNumber {
		return nil
	}
	comicList := make([]model.Comic, 0)
	for i := lastStoredComic; i < latestComicNumber; i++ {
		go fetchComic(i + 1)
	}

	for i := lastStoredComic; i < latestComicNumber; i++ {
		comicList = append(comicList, <-comicChan)
	}
	db.BatchStoreComics(comicList)
	return comicList
}

func FetchIncompleteComics(incompleteComics []model.Comic) []model.Comic {
	incompleteList := make([]model.Comic, 0)
	for _, comic := range incompleteComics {
		go fetchComic(comic.Num)
	}

	for i := 0; i < len(incompleteComics); i++ {
		incompleteList = append(incompleteList, <-comicChan)
	}

	db.UpdateIncompleteComics(incompleteList)
	return incompleteList
}
