// Routes and handlers for the api

package api

import (
	"context"
	"math"
	"sort"
	"strings"
	"xkcd/db"
	"xkcd/index"
	"xkcd/model"
	"xkcd/nlp"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var comics []model.Comic
var comicFreq model.ComicFreq
var stats model.Stats
var ginLambda *ginadapter.GinLambda

func AWSHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func Serve(local bool) {
	r := gin.Default()
	r.Use(corsMiddleware())
	db.Connect()
	comics = db.GetComics()
	comicFreq = db.GetComicFreq()
	comicFreq.TotalComics = len(comics)
	stats.LastIndexedComic = len(comics)
	stats.LastCreatedComic = index.GetCurrentComicNum()
	stats = db.GetStats(stats)
	words := db.GetRawWords()
	nlp.TrainModel(words)

	r.GET("/", handleLoad)

	//handle suggestions: /suggest?q="incomple"
	r.GET("/suggest", handleSuggest)

	// handle search: /search?q="query"
	r.GET("/search", handleSearch)

	//handle stats: /stats
	r.GET("/stats", handleStats)

	if local {
		r.Run()
	} else {
		gin.SetMode(gin.ReleaseMode)
		ginLambda = ginadapter.New(r)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func handleLoad(c *gin.Context) {
	c.JSON(200, "Welcome to xkcd-search!")
}

func handleStats(c *gin.Context) {
	c.JSON(200, stats)
}

func handleSuggest(c *gin.Context) {
	var request model.Suggest
	err := c.Bind(&request)
	if err != nil {
		c.JSON(400, "Sorry! You didn't call the api correctly. Expected route is /suggest?q=incom ")
		return
	}
	terms := strings.Fields(request.Query)
	currTerm := terms[len(terms)-1]
	currTerm, _ = nlp.CleanAndStem(currTerm)

	termSuggestions := nlp.Autocomplete(currTerm)
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
	var autocorrectedRaw string
	hasTypo := false
	err := c.Bind(&request)
	if err != nil {
		c.JSON(400, "Sorry! You didn't call the api correctly. Expected route is /search?q=test&autocomplete=true")
		return
	}
	query := request.Query
	autocorrect := request.Autocorrect

	rawQuery, stemQuery := nlp.CleanAndStem(query)
	if autocorrect {
		hasTypo, autocorrectedRaw = nlp.Autocorect(rawQuery)
		if hasTypo {
			rawQuery, stemQuery = nlp.CleanAndStem(autocorrectedRaw)
		}
	}
	rankings := index.RankQuery(rawQuery, stemQuery, comics, comicFreq)

	//if no comics are found, return 404
	if len(rankings) == 0 {
		c.JSON(404, rankings)
		return
	}

	//sort rankings and return at the correct page
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Rank >= rankings[j].Rank
	})

	totalRanked := len(rankings)

	if len(rankings) > 10 {
		pageCount := int(math.Min(10.0, float64(len(rankings))))                         //total number of comics on one page
		start := int(math.Min(float64(request.Start), float64(len(rankings)-pageCount))) //starting index of the page, defaults to 0
		rankings = rankings[start : start+pageCount]
	}

	//if there was a typo, return the autocorrected version and ranking
	if hasTypo {
		c.JSON(200, gin.H{
			"autocorrect": autocorrectedRaw,
			"rankings":    rankings,
			"totalRanked": totalRanked,
		})
		return
	}

	//return just the ranking if there was no typo
	c.JSON(200, rankings)
	return
}
