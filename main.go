package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"xkcd/db"
	"xkcd/model"
)

const baseURL = "https://xkcd.com/"
const suffixURL = "/info.0.json"

var currentNum int
var comicChan = make(chan model.Comic)

func main() {
	db.New()

	currComicURL := fmt.Sprintf("%s%s", baseURL, suffixURL)
	resp, err := http.Get(currComicURL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var currComic model.Comic
	json.Unmarshal([]byte(body), &currComic)
	currentNum = currComic.Number
	getAllComics()
}

func getComic(num int) {
	path := fmt.Sprintf("%s%d%s", baseURL, num, suffixURL)
	resp, err := http.Get(path)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var requestedComic model.Comic
	json.Unmarshal([]byte(body), &requestedComic)
	comicChan <- requestedComic
}

func getAllComics() []model.Comic {
	output := make([]model.Comic, currentNum)
	for i := 0; i < currentNum; i++ {
		go getComic(i + 1)
		output[i] = <-comicChan

		if output[i].Title == "" {
			output[i].Title = fmt.Sprintf("%d", i) // Some comics have no title (such as 404 LOL) so set the number as title
		} else {
			log.Printf("q%d: %s\n", i+1, output[i].Title)
		}
	}
	return output
}
