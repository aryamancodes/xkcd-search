/*
Fetch all current comics, explanations and calculate tf-idf statistics.
*/

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

	"github.com/tidwall/gjson"
)

const comicBaseURL = "https://xkcd.com/"
const comicSuffixURL = "/info.0.json"
const explainBaseURL = "https://www.explainxkcd.com/wiki/api.php?action=parse&page=%d:_%s&prop=wikitext&format=json"

var currentNum int

var comicChan = make(chan model.Comic)
var explainChan = make(chan string)
var termFreqChan = make(chan map[string]int)

func getTermFreq(title string, alt string, explain string, transcript *string) {
	termFreq := make(map[string]int)
	allTerms := fmt.Sprintf("%s %s %s", title, alt, explain)
	if transcript != nil {
		allTerms = fmt.Sprintf("%s %s", allTerms, *transcript)
	}
	allTerms = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(allTerms, "")

	for _, word := range strings.Fields(allTerms) {
		currFreq := termFreq[word]
		termFreq[word] = currFreq + 1
	}
	termFreqChan <- termFreq
}

func getDocFreq() {
	//TODO!!
}

func getComic(num int) {
	path := fmt.Sprintf("%s%d%s", comicBaseURL, num, comicSuffixURL)
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
	explainList := make([]string, currentNum)
	termFreqList := make(map[int](map[string]int))
	for i := 0; i < currentNum; i++ {
		go getComic(i + 1)
		latestComic := <-comicChan

		go getExplanation(latestComic.Number, latestComic.Title)
		latestExplain := <-explainChan

		go getTermFreq(latestComic.Title, latestComic.AltTitle, latestExplain, latestComic.Transcript)
		termFreqList[i] = <-termFreqChan

		comicList[i] = latestComic
		explainList[i] = latestExplain
		//fmt.Printf("%d: %s %v %s\n", comicList[i].Number, comicList[i].Title, termFreqList[i], explainList[i])
	}
	return comicList, termFreqList
}

func getCurrentComicNumber() {
	currComicURL := fmt.Sprintf("%s%s", comicBaseURL, comicSuffixURL)
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

func getExplanation(num int, title string) {
	title = strings.ReplaceAll(title, " ", "_")
	explainURL := fmt.Sprintf(explainBaseURL, num, title)
	resp, err := http.Get(explainURL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	explanation := gjson.Get(string(body), "parse.wikitext.*").String()

	//store the explanation result as is for now
	explanation = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(explanation, "")
	explainChan <- explanation
}
