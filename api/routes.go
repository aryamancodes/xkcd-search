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

	// handle search: /search?q="query"
	r.GET("/search", handleSearch)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func handleLoad(c *gin.Context) {
	c.JSON(200, "Welcome to xkcd-search!")
}

func handleSearch(c *gin.Context) {
	var request model.Search
	var query string

	if c.ShouldBind(&request) == nil {
		query = strings.Replace(request.Query, "+", " ", -1)
		rawQuery, stemQuery := nlp.CleanAndStem(query)
		hasTypo, autocorrectedRaw := nlp.Autocorect(language, rawQuery)
		if hasTypo {
			rawQuery, stemQuery = nlp.CleanAndStem(autocorrectedRaw)
		}
		rankings := index.RankQuery(rawQuery, stemQuery, comicFreq)
		if len(rankings) == 0 {
			c.JSON(404, "No results found. Maybe there isn't an xkcd for everything")
		} else {
			sort.Slice(rankings, func(i, j int) bool {
				return rankings[i].Rank >= rankings[j].Rank
			})

			length := int(math.Min(float64(len(rankings)), 9))
			rankings = rankings[:length]
			if hasTypo {
				c.JSON(200, gin.H{
					"autocorrect": autocorrectedRaw,
					"rankings":    rankings,
				})
			} else {
				c.JSON(200, rankings)
			}
		}
	}
}
