// Functions to fetch all the data to index -- including the latest comic number and all explanations (concurrently) upto the latest

package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"xkcd/db"
	"xkcd/model"
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
			if fetchedExplainWiki.Parse.Title == "" {
				continue
			} else {
				break
			}
		}
	}

	explainChan <- fetchedExplainWiki
}

// (concurrently) fetch all explanations based on the current comic number
func FetchAllExplanations() []model.ExplainWikiJson {
	latestComicNumber := getCurrentComicNum()
	test := int(math.Max(250, float64(latestComicNumber)))
	explanationsList := make([]model.ExplainWikiJson, 0)

	for i := 0; i < test; i++ {
		go fetchExplanation(i + 1)
	}

	for i := 0; i < test; i++ {
		explanation := <-explainChan
		db.StoreComics(i, explanation.Parse.Wikitext.Content)
		explanationsList = append(explanationsList, explanation)
	}

	return explanationsList
}
