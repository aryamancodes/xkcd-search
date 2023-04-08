// Routes and handlers for the api

package api

import (
	"math"
	"sort"
	"strings"
	"xkcd/db"
	"xkcd/index"
	"xkcd/model"
	"xkcd/nlp"

	"github.com/gin-gonic/gin"
	"github.com/sajari/fuzzy"
)

var comicFreq model.ComicFreq
var isLoaded = false
var language *fuzzy.Model

func Serve() {
	r := gin.Default()
	db.Connect()
	comicFreq = db.GetComicFreq()
	words := db.GetRawWords()
	language = nlp.TrainModel(words)

	r.GET("/", handleLoad)

	//handle suggestions: /suggest?q="incomple"
	r.GET("/suggest", handleSuggest)

	// handle search: /search?q="query"
	r.GET("/search", handleSearch)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func handleLoad(c *gin.Context) {
	c.JSON(200, "Welcome to xkcd-search!")
}

func handleSuggest(c *gin.Context) {
	var request model.Suggest
	var terms []string

	if c.ShouldBind(&request) == nil {
		terms = strings.Fields(request.Query)
	}
	currTerm := terms[len(terms)-1]
	currTerm, _ = nlp.CleanAndStem(currTerm)

	termSuggestions := nlp.Autocomplete(language, currTerm)
	var querySuggestions []string
	for _, suggest := range termSuggestions {
		terms = append(terms[:len(terms)-1], suggest)
		querySuggest := strings.Join(terms, " ")
		querySuggestions = append(querySuggestions, querySuggest)
	}

	c.JSON(200, querySuggestions)

}

func handleSearch(c *gin.Context) {
	var request model.Search
	var query string

	if c.ShouldBind(&request) == nil {
		query = strings.Replace(request.Query, "+", " ", -1)
	}
	rawQuery, stemQuery := nlp.CleanAndStem(query)
	hasTypo, autocorrectedRaw := nlp.Autocorect(language, rawQuery)
	if hasTypo {
		rawQuery, stemQuery = nlp.CleanAndStem(autocorrectedRaw)
	}
	rankings := index.RankQuery(rawQuery, stemQuery, comicFreq)

	//if no comics are found, return 404
	if len(rankings) == 0 {
		c.JSON(404, "No results found. Maybe there isn't an xkcd for everything")
		return
	}

	//sort rankings and return at most 10 comics
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Rank >= rankings[j].Rank
	})
	length := int(math.Min(float64(len(rankings)), 10))
	rankings = rankings[:length]

	//if there was a typo, return the autocorrected version and ranking
	if hasTypo {
		c.JSON(200, gin.H{
			"autocorrect": autocorrectedRaw,
			"rankings":    rankings,
		})
		return
	}

	//return just the ranking if there was no typo
	c.JSON(200, rankings)
}
