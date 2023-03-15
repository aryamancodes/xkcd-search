// Functions to fetch all the data to index -- including the latest comic number and all explanations (concurrently) upto the latest

package index

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
	"xkcd/model"

	"encoding/json"
)

const currentComicURL = "https://xkcd.com/info.0.json"
const explainationURL = "https://www.explainxkcd.com/wiki/api.php?action=parse&page=%d&prop=wikitext&format=json&redirects=1&origin=*"

var explainChan = make(chan model.ExplainWikiJson, 250)

// get the latest comics number directly from the xkcd api
func getCurrentComicNum() int {
	resp, err := http.Get(currentComicURL)
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

func fetchExplanation(num int) {
	explainURL := fmt.Sprintf(explainationURL, num)

	var resp *http.Response

	//try to fetch explanation from wiki
	for {
		tryResp, err := http.Get(explainURL)
		if err != nil {
			log.Fatalln(err)
		}

		//concurrent fetches may lead to >500 cloudflare errors
		if tryResp.StatusCode != http.StatusOK {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		} else {
			resp = tryResp
			break
		}
	}

	//try to unmarshal response into json
	var fetchedExplainWiki model.ExplainWikiJson
	for {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()

		json.Unmarshal([]byte(body), &fetchedExplainWiki)

		//even if we get status code 200, an internal wiki error might have occured -- needs refetching
		if fetchedExplainWiki.Parse.Title == "" {
			fetchExplanation(num)
		} else {
			break
		}
	}

	explainChan <- fetchedExplainWiki
}

// fetch all explanations based on the current cominc number
func FetchAllExplanations() []model.ExplainWikiJson {
	latestComicNumber := getCurrentComicNum()
	explanationsList := make([]model.ExplainWikiJson, latestComicNumber)

	for i := 0; i < latestComicNumber; i++ {
		go fetchExplanation(i + 1)
	}

	for i := 0; i < latestComicNumber; i++ {
		explanationsList[i] = <-explainChan
	}

	fmt.Printf("FINAL RESULT IS %d\n", len(explanationsList))
	return explanationsList
}
