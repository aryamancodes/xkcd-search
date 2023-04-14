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

var comicFreq model.ComicFreq
var ginLambda *ginadapter.GinLambda

func AWSHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func Serve(local bool) {
	r := gin.Default()
	db.Connect()
	comicFreq = db.GetComicFreq()
	words := db.GetRawWords()
	nlp.TrainModel(words)

	r.GET("/", handleLoad)

	//handle suggestions: /suggest?q="incomple"
	r.GET("/suggest", handleSuggest)

	// handle search: /search?q="query"
	r.GET("/search", handleSearch)

	if local {
		r.Use(corsMiddleware())
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

func handleSuggest(c *gin.Context) {
	var request model.Suggest
	err := c.Bind(&request)
	if err != nil {
		c.JSON(400, "Sorry! You didn't call the api correctly")
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
		c.JSON(400, "Sorry! You didn't call the api correctly")
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
	return
}
