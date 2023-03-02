package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"xkcd/model"

	"encoding/json"
)

const baseURL = "https://xkcd.com/"
const suffixURL = "/info.0.json"

var currentNum int

var comicChan = make(chan model.Comic)
var termFreqChan = make(chan map[string]int)

func getTermFreq(title string, alt string, transcript *string) {
	termFreq := make(map[string]int)
	allTerms := fmt.Sprintf("%s %s", title, alt)
	if transcript != nil {
		allTerms = fmt.Sprintf("%s %s", allTerms, *transcript)
	}
	allTerms = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(allTerms, "")

	for _, word := range strings.Fields(allTerms) {
		currFreq := termFreq[word]
		termFreq[word] = currFreq + 1
	}
	fmt.Printf("ADDED FREQ FOR: %s", title)
	fmt.Println("THE FREQ MAP IS ", termFreq)

	termFreqChan <- termFreq
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

	// Some comics have no title (such as 404 LOL) so set the number as title
	if requestedComic.Title == "" {
		requestedComic.Title = "Not Found"
	}
	comicChan <- requestedComic
}

func GetAllComics() ([]model.Comic, map[int](map[string]int)) {
	getCurrentComicNumber()
	comicList := make([]model.Comic, currentNum)
	termFreqList := make(map[int](map[string]int))
	for i := 0; i < currentNum; i++ {
		go getComic(i + 1)
		comicList[i] = <-comicChan
		fmt.Printf("ADDED COMIC: %s\n", comicList[i].Title)

		go getTermFreq(comicList[i].Title, comicList[i].AltTitle, comicList[i].Transcript)

		termFreqList[i] = <-termFreqChan

		if comicList[i].Title == "" {
			comicList[i].Title = "Not Found"
		} else {
			log.Printf("q%d: %s\n", i+1, comicList[i].Title)
		}
	}
	return comicList, termFreqList
}

func getCurrentComicNumber() {
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
}
