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

var currentNum int = 1

func main() {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	db.New()

	var currentComic model.Comic
	json.Unmarshal([]byte(body), &currentComic)
	currentNum = currentComic.Number
	getAllComics()
}

func getComic(num int) model.Comic {
	path := fmt.Sprintf("%s%d%s", "https://xkcd.com/", num, "/info.0.json")
	resp, err := http.Get(path)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var output model.Comic
	json.Unmarshal([]byte(body), &output)
	return output
}

func getAllComics() []model.Comic {
	output := make([]model.Comic, currentNum)
	for i := 0; i < currentNum; i++ {
		output[i] = getComic(i + 1)
		if output[i].Title == "" {
			log.Printf("WARNING: THIS COMIC HAS NO TITLE %d\n", i+1)
		} else {
			log.Printf("%d: %s\n", i+1, output[i].Title)
		}
	}
	return output
}
